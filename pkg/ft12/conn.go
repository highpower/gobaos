package ft12

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"syscall"
)

type frame []byte

// Conn is an FT1.2 serial connection.
type Conn struct {
	needReset    bool
	file         *os.File
	incoming     bytes.Buffer
	incomingEven bool
	outgoing     []frame
	outgoingEven bool
}

// Poll reads incoming FT1.2 frames from the serial connection and buffers the payload for later retrieval via Read.
func (c *Conn) Poll() error {
	if err := c.resetIfNeeded(); err != nil {
		return err
	}
	buffer := [260]byte{}
	reading := true
	for reading {
		n, err := c.readFrameToBuffer(buffer[:])
		switch {
		case err != nil:
			return err
		case n == 0:
			reading = false
		}
	}
	for _, f := range c.outgoing {
		if _, err := c.file.Write(f); err != nil {
			return err
		}
	}
	c.outgoing = c.outgoing[:0]
	return nil
}

func (c *Conn) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if len(p) < c.incoming.Len() {
		return 0, fmt.Errorf("%w: %d, expected %d", ErrNoSpace, len(p), c.incoming.Len())
	}
	result, err := c.incoming.Read(p)
	if err != nil {
		return 0, err
	}
	c.incoming.Reset()
	return result, err
}

func (c *Conn) Write(p []byte) (int, error) {
	for len(p) != 0 {
		length := min(len(p), 255)
		fr := make(frame, length+7) // sizeof(header) + cr + checksum + 0x16
		copy(fr[:4], []byte{0x68, byte(length), byte(length), 0x68})
		if c.outgoingEven {
			fr[4] = 0x53
		} else {
			fr[4] = 0x73
		}
		c.outgoingEven = !c.outgoingEven
		copy(fr[5:length+5], p[:length])
		fr[len(fr)-2] = checksum(p[:length])
		fr[len(fr)-1] = 0x16
		c.outgoing = append(c.outgoing, fr)
		p = p[length:]
	}
	return len(p), nil
}

// Close closes the underlying serial device.
func (c *Conn) Close() error {
	return c.file.Close()
}

// Available returns the number of unread payload bytes in the internal buffer.
func (c *Conn) Available() int {
	return c.incoming.Len()
}

func (c *Conn) resetIfNeeded() error {
	if !c.needReset {
		return nil
	}
	if err := c.writeResetRequest(); err != nil {
		return err
	}
	c.needReset = false
	return nil
}

func (c *Conn) writeResetRequest() error {
	if _, err := c.file.Write([]byte{0x10, 0x40, 0x40, 0x16}); err != nil {
		return err
	}
	return c.readAck()
}

func (c *Conn) readFrameToBuffer(buffer []byte) (int, error) {
	hd := header{}
	n, err := io.ReadFull(c.file, hd[:])
	switch {
	case err != nil:
		return 0, err
	case n == 0:
		return 0, nil
	case hd.isResetIndication():
		c.needReset = true
		return 0, nil
	case hd.valid() != nil:
		return 0, fmt.Errorf("%w header: %x", ErrInvalid, hd)
	}
	if _, err = io.ReadFull(c.file, buffer[:hd.len()+2]); err != nil {
		return 0, err
	}
	if err := c.writeAck(); err != nil {
		return 0, err
	}
	data, err := c.unmarshalFrame(buffer[:hd.len()+2])
	if err != nil {
		return 0, err
	}
	c.incoming.Write(data)
	return 0, nil
}

func (c *Conn) unmarshalFrame(data []byte) ([]byte, error) {
	if len(data) < 3 {
		return nil, fmt.Errorf("%w frame length", ErrInvalid)
	}
	if (c.incomingEven && data[0] != 0xD3) || (!c.incomingEven && data[0] != 0xF3) {
		return nil, fmt.Errorf("%w control byte: %d", ErrInvalid, data[0])
	}
	c.incomingEven = !c.incomingEven
	calculated := checksum(data[:len(data)-2])
	if cs := data[len(data)-2]; cs != calculated {
		return nil, fmt.Errorf("%w checksum: %d, expected %d", ErrInvalid, calculated, cs)
	}
	if data[len(data)-1] != 0x16 {
		return nil, fmt.Errorf("%w last byte: %d", ErrInvalid, data[len(data)-1])
	}
	return data[1 : len(data)-2], nil
}

func (c *Conn) readAck() error {
	buffer := [1]byte{}
	if _, err := io.ReadFull(c.file, buffer[:]); err != nil {
		return err
	}
	switch buffer[0] {
	case 0xE5:
		return nil
	case 0x10:
		return c.readResetIndicationRest()
	}
	return fmt.Errorf("%w ack: %02x", ErrInvalid, buffer[0])
}

func (c *Conn) readResetIndicationRest() error {
	buffer := [3]byte{}
	if _, err := io.ReadFull(c.file, buffer[:]); err != nil {
		return err
	}
	if !bytes.Equal(buffer[:], []byte{0xC0, 0xC0, 0x16}) {
		return fmt.Errorf("%w reset indication", ErrInvalid)
	}
	c.needReset = true
	return nil
}

func (c *Conn) writeAck() error {
	buffer := [1]byte{0xE5}
	if _, err := c.file.Write(buffer[:]); err != nil {
		return err
	}
	return nil
}

// NewConn opens an FT1.2 connection on the named TTY device.
func NewConn(ttyName string) (*Conn, error) {
	f, err := os.OpenFile(ttyName, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	fd := int(f.Fd())
	if err := syscall.SetNonblock(fd, true); err != nil {
		_ = f.Close()
		return nil, err
	}
	if err := flush(fd); err != nil {
		_ = f.Close()
		return nil, err
	}
	return &Conn{needReset: true, file: f, incoming: bytes.Buffer{}, incomingEven: false,
		outgoing: nil, outgoingEven: false}, nil
}

func checksum(data []byte) byte {
	checksum := 0
	for _, b := range data {
		checksum += int(b)
	}
	return byte(checksum % 256)
}
