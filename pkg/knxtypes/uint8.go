package knxtypes

import (
	"fmt"
	"strconv"
)

type DPTUInt8 struct {
	Value uint8
}

func (u *DPTUInt8) Get() string {
	return fmt.Sprintf("%d", u.Value)
}

func (u *DPTUInt8) Set(s string) error {
	value, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return err
	}
	u.Value = uint8(value)
	return nil
}

func (u *DPTUInt8) TypeCode() TypeCode {
	return UInt8
}

func (u *DPTUInt8) String() string {
	return fmt.Sprintf("dptUInt8[%d]", u.Value)
}

func (u *DPTUInt8) MarshalBinary() ([]byte, error) {
	return []byte{u.Value}, nil
}

func (u *DPTUInt8) UnmarshalBinary(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("%w: type %s requires 1 byte, got %d", ErrUnmarshal, UInt8.String(), len(data))
	}
	u.Value = data[0]
	return nil
}
