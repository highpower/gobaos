package knxtypes

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

type DPTUint32 struct {
	Value uint32
}

func (d *DPTUint32) Get() string {
	return fmt.Sprintf("%d", d.Value)
}

func (d *DPTUint32) Set(s string) error {
	value, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return err
	}
	d.Value = uint32(value)
	return nil
}

func (d *DPTUint32) TypeCode() TypeCode {
	return UInt32
}

func (d *DPTUint32) String() string {
	return fmt.Sprintf("dptUInt32[%d]", d.Value)
}

func (d *DPTUint32) MarshalBinary() ([]byte, error) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, d.Value)
	return b, nil
}

func (d *DPTUint32) UnmarshalBinary(data []byte) error {
	if len(data) != 4 {
		return fmt.Errorf("%w: type %s requires 4 bytes, got %d", ErrUnmarshal, UInt32.String(), len(data))
	}
	d.Value = binary.BigEndian.Uint32(data)
	return nil
}
