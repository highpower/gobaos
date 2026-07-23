package knxtypes

import (
	"fmt"
	"strconv"
)

type DPTInt8 struct {
	Value int8
}

func (d *DPTInt8) Get() string {
	return fmt.Sprintf("%d", d.Value)
}

func (d *DPTInt8) Set(s string) error {
	value, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		return err
	}
	d.Value = int8(value)
	return nil
}

func (d *DPTInt8) TypeCode() TypeCode {
	return Int8
}

func (d *DPTInt8) String() string {
	return fmt.Sprintf("int8[%d]", d.Value)
}

func (d *DPTInt8) MarshalBinary() ([]byte, error) {
	return []byte{byte(d.Value)}, nil
}

func (d *DPTInt8) UnmarshalBinary(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("%w: type %s requires 1 byte, got %d", ErrUnmarshal, Int8.String(), len(data))
	}
	d.Value = int8(data[0])
	return nil
}
