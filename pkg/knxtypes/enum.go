package knxtypes

import (
	"fmt"
	"strconv"
)

// DPTEnum represents KNX DPT 20.xxx (1-byte enumeration).
type DPTEnum struct {
	Value uint8
}

// Get returns the enumeration value as a decimal string.
func (e *DPTEnum) Get() string {
	return fmt.Sprintf("%d", e.Value)
}

// Set parses a decimal string and stores the value as an unsigned 8-bit integer.
func (e *DPTEnum) Set(value string) error {
	v, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		return err
	}
	e.Value = uint8(v)
	return nil
}

// TypeCode returns the KNX type code for DPTEnum.
func (e *DPTEnum) TypeCode() TypeCode {
	return EnumCode
}

// String returns a human-readable representation of DPTEnum.
func (e *DPTEnum) String() string {
	return fmt.Sprintf("dptEnum[%d]", e.Value)
}

// MarshalBinary encodes DPTEnum into a 1-byte payload.
func (e *DPTEnum) MarshalBinary() ([]byte, error) {
	return []byte{e.Value}, nil
}

// UnmarshalBinary decodes a 1-byte payload into DPTEnum.
func (e *DPTEnum) UnmarshalBinary(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("%w: type %s requires 1 byte, got %d", ErrUnmarshal, EnumCode.String(), len(data))
	}
	e.Value = data[0]
	return nil
}
