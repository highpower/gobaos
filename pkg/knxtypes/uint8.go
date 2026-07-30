package knxtypes

import (
	"fmt"
	"strconv"
)

// DPTUint8 represents KNX DPT 5.xxx (unsigned 8-bit integer).
type DPTUint8 struct {
	Value uint8
}

// Get returns the value as a decimal string.
func (u *DPTUint8) Get() string {
	return fmt.Sprintf("%d", u.Value)
}

// Set parses a decimal string and stores the value.
func (u *DPTUint8) Set(s string) error {
	value, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return err
	}
	u.Value = uint8(value)
	return nil
}

// TypeCode returns the KNX type code for DPTUint8.
func (u *DPTUint8) TypeCode() TypeCode {
	return Uint8Code
}

// String returns a human-readable representation of DPTUint8.
func (u *DPTUint8) String() string {
	return fmt.Sprintf("dptUInt8[%d]", u.Value)
}

// MarshalBinary encodes DPTUint8 into a 1-byte payload.
func (u *DPTUint8) MarshalBinary() ([]byte, error) {
	return []byte{u.Value}, nil
}

// UnmarshalBinary decodes a 1-byte payload into DPTUint8.
func (u *DPTUint8) UnmarshalBinary(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("%w: type %s requires 1 byte, got %d", ErrUnmarshal, Uint8Code.String(), len(data))
	}
	u.Value = data[0]
	return nil
}
