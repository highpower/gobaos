package knxtypes

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

type DPTUint16 struct {
	Value uint16
}

func (d *DPTUint16) Get() string {
	return fmt.Sprintf("%d", d.Value)
}

func (d *DPTUint16) Set(s string) error {
	value, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return err
	}
	d.Value = uint16(value)
	return nil
}

func (d *DPTUint16) TypeCode() TypeCode {
	return UInt16
}

func (d *DPTUint16) String() string {
	return fmt.Sprintf("dptUInt16[%d]", d.Value)
}

func (d *DPTUint16) MarshalBinary() ([]byte, error) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, d.Value)
	return b, nil
}

func (d *DPTUint16) UnmarshalBinary(data []byte) error {
	if len(data) != 2 {
		return fmt.Errorf("%w: type %s requires 2 bytes, got %d", ErrUnmarshal, UInt16.String(), len(data))
	}
	d.Value = binary.BigEndian.Uint16(data)
	return nil
}
