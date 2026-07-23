package knxtypes

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

type DPTInt32 struct {
	Value int32
}

func (d *DPTInt32) Get() string {
	return fmt.Sprintf("%d", d.Value)
}

func (d *DPTInt32) Set(s string) error {
	value, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return err
	}
	d.Value = int32(value)
	return nil
}

func (d *DPTInt32) TypeCode() TypeCode {
	return Int32
}

func (d *DPTInt32) String() string {
	return fmt.Sprintf("dptInt32[%d]", d.Value)
}

func (d *DPTInt32) MarshalBinary() ([]byte, error) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(d.Value))
	return b, nil
}

func (d *DPTInt32) UnmarshalBinary(data []byte) error {
	if len(data) != 4 {
		return fmt.Errorf("%w: type %s requires 4 bytes, got %d", ErrUnmarshal, Int32.String(), len(data))
	}
	d.Value = int32(binary.BigEndian.Uint32(data))
	return nil
}
