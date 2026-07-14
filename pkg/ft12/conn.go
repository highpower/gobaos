package ft12

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"syscall"
)

type Conn struct {
	needReset  bool
	serverEven bool
	hostEven   bool
	file       *os.File
	data       bytes.Buffer
}

func (c *Conn) Read(p []byte) (int, error) {
	if len(p) != 0 && c.data.Len() != 0 {
		return c.data.Read(p)
	}
	if err := c.readToBuffer(); err != nil {
		return 0, err
	}
	return 0, nil
}

func (c *Conn) Write(p []byte) (int, error) {
	if err := c.resetIfNeeded(); err != nil {
		return 0, err
	}
	return 0, nil
}

func (c *Conn) Close() error {
	return c.file.Close()
}

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

func (c *Conn) readToBuffer() error {
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

func (c *Conn) readFrameToBuffer(buffer []byte) (int, error) {
	hd := header{}
	n, err := io.ReadFull(c.file, hd[:])
	switch {
	case err != nil:
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
	for i := 0; i < len(data)-3; i++ {
		checksum += int(data[i])
	}
	if cs := int(data[len(data)-2]); cs != checksum%256 {
		return nil, fmt.Errorf("%w checksum: %d, expected %d", ErrInvalid, checksum, cs)
	}
	return data[1 : len(data)-2], nil
}

func (c *Conn) readAck() error {
	buffer := [1]byte{}
	if _, err := io.ReadFull(c.file, buffer[:]); err != nil {
		return err
	}
	if buffer[0] != 0xE5 {
		return fmt.Errorf("%w ack: %02x", ErrInvalid, buffer[0])
	}
	return nil
}

func (c *Conn) writeAck() error {
	buffer := [1]byte{0xE5}
	if _, err := c.file.Write(buffer[:]); err != nil {
		return err
	}
	return nil
}

func NewConn(ttyName string) (*Conn, error) {
	f, err := os.OpenFile(ttyName, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	if err := syscall.SetNonblock(int(f.Fd()), true); err != nil {
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
