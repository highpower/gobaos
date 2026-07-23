package knxtypes

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

type DPTInt16 struct {
	Value int16
}

func (d *DPTInt16) Get() string {
	return fmt.Sprintf("%d", d.Value)
}

func (d *DPTInt16) Set(s string) error {
	value, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		return err
	}
	d.Value = int16(value)
	return nil
}

func (d *DPTInt16) TypeCode() TypeCode {
	return Int16
}

func (d *DPTInt16) String() string {
	return fmt.Sprintf("int16[%d]", d.Value)
}

func (d *DPTInt16) MarshalBinary() ([]byte, error) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(d.Value))
	return b, nil
}

func (d *DPTInt16) UnmarshalBinary(data []byte) error {
	if len(data) != 2 {
		return fmt.Errorf("%w: type %s requires 2 bytes, got %d", ErrUnmarshal, Int16.String(), len(data))
	}
	d.Value = int16(binary.BigEndian.Uint16(data))
	return nil
}
