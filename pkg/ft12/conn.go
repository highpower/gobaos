package ft12

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"syscall"
)

type Conn struct {
	needReset bool
	file      *os.File
	data      bytes.Buffer
}

func (c *Conn) Read(p []byte) (int, error) {
	if p == nil && c.data.Len() != 0 {
		return c.data.Read(p)
	}
	if err := c.readToBuffer(); err != nil {
		return 0, err
	}
	return 0, nil
}

func (c *Conn) Write(p []byte) (int, error) {
	if err := c.writeResetIfNeeded(); err != nil {
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

func (c *Conn) readToBuffer() error {
	if err := c.writeResetIfNeeded(); err != nil {
		return err
	}
	buffer := [256]byte{}
	reading := true
	for reading {
		n, err := c.readFrame(buffer[:])
		switch {
		case err != nil:
			return err
		case n == 0:
			reading = false
		}
	}
	return nil
}

func (c *Conn) readFrame(buffer []byte) (int, error) {
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
		return 0, fmt.Errorf("invalid header: %x", hd)
	}
	if _, err = io.ReadFull(c.file, buffer[:hd.len()]); err != nil {
		return 0, err
	}
	if err := c.writeAck(); err != nil {
		return 0, err
	}
	data, err := unmarshal(buffer[:hd.len()])
	if err != nil {
		return 0, err
	}
	c.data.Write(data)
	return 0, nil
}

func (c *Conn) readAck() error {
	buffer := [1]byte{}
	if _, err := io.ReadFull(c.file, buffer[:]); err != nil {
		return err
	}
	if buffer[0] != 0xE5 {
		return fmt.Errorf("invalid ft12 ack")
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

func (c *Conn) writeResetIfNeeded() error {
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

func NewConn(ttyName string) (*Conn, error) {
	f, err := os.OpenFile(ttyName, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	if err := syscall.SetNonblock(int(f.Fd()), true); err != nil {
		_ = f.Close()
		return nil, err
	}
	return &Conn{needReset: true, file: f, data: bytes.Buffer{}}, nil
}
