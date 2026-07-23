package knxtypes

import (
	"fmt"
	"strconv"
)

type DPTBool struct {
	Value bool
}

func (b *DPTBool) Get() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *DPTBool) Set(s string) error {
	value, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	b.Value = value
	return nil
}

func (b *DPTBool) TypeCode() TypeCode {
	return Bool
}

func (b *DPTBool) BytesRequired() int {
	return 1
}

func (b *DPTBool) String() string {
	return fmt.Sprintf("dptBool[%t]", b.Value)
}

func (b *DPTBool) MarshalBinary() ([]byte, error) {
	result := []byte{0}
	if b.Value {
		result[0] = 1
	}
	return result, nil
}

func (b *DPTBool) UnmarshalBinary(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("%w: type %s requires 1 byte, got %b", ErrUnmarshal, Bool.String(), len(data))
	}
	b.Value = data[0]&0x01 != 0
	return nil
}
