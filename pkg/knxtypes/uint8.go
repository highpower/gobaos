package knxtypes

import (
	"fmt"
	"strconv"
)

type DPTUint8 struct {
	Value uint8
}

func (u *DPTUint8) Get() string {
	return fmt.Sprintf("%d", u.Value)
}

func (u *DPTUint8) Set(s string) error {
	value, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return err
	}
	u.Value = uint8(value)
	return nil
}

func (u *DPTUint8) TypeCode() TypeCode {
	return Uint8
}

func (u *DPTUint8) String() string {
	return fmt.Sprintf("dptUInt8[%d]", u.Value)
}

func (u *DPTUint8) MarshalBinary() ([]byte, error) {
	return []byte{u.Value}, nil
}

func (u *DPTUint8) UnmarshalBinary(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("%w: type %s requires 1 byte, got %d", ErrUnmarshal, Uint8.String(), len(data))
	}
	u.Value = data[0]
	return nil
}
