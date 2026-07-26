package ft12

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"syscall"
)

// Conn is an FT1.2 serial connection.
type Conn struct {
	needReset  bool
	serverEven bool
	hostEven   bool
	file       *os.File
	data       bytes.Buffer
}

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
	return nil
}

func (c *Conn) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if len(p) < c.data.Len() {
		return 0, fmt.Errorf("%w data: %d, expected %d", ErrNoSpace, len(p), c.data.Len())
	}
	result, err := c.data.Read(p)
	if err != nil {
		return 0, err
	}
	c.data.Reset()
	return result, err
}

func (c *Conn) Write(p []byte) (int, error) {
	if err := c.resetIfNeeded(); err != nil {
		return 0, err
	}

	return 0, nil
}

// Close closes the underlying serial device.
func (c *Conn) Close() error {
	return c.file.Close()
}

// Available returns the number of unread payload bytes in the internal buffer.
func (c *Conn) Available() int {
	return c.data.Len()
}

func (c *Conn) resetIfNeeded() error {
	if c.needReset {
		err := c.writeResetRequest()
		c.needReset = false
		return err
	}
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
	c.data.Write(data)
	return 0, nil
}

func (c *Conn) unmarshalFrame(data []byte) ([]byte, error) {
	if len(data) < 3 {
		return nil, fmt.Errorf("%w frame length", ErrInvalid)
	}
	if !validControl(c.serverEven, data[0]) {
		return nil, fmt.Errorf("%w control byte: %d", ErrInvalid, data[0])
	}
	c.serverEven = !c.serverEven
	if data[len(data)-1] != 0x16 {
		return nil, fmt.Errorf("%w last byte: %d", ErrInvalid, data[len(data)-1])
	}
	checksum := 0
	for i := 0; i < len(data)-2; i++ {
		checksum += int(data[i])
	}
	if cs := int(data[len(data)-2]); cs != checksum%256 {
		return nil, fmt.Errorf("%w checksum: %d, expected %d", ErrInvalid, checksum%256, cs)
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
	return &Conn{needReset: true, serverEven: false, hostEven: false, file: f, data: bytes.Buffer{}}, nil
}

func validControl(even bool, seen byte) bool {
	if even {
		return seen == 0xD3
	}
	return seen == 0xF3
}
