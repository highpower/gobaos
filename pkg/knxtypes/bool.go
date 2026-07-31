package knxtypes

import (
	"fmt"
	"strconv"
	"strings"
)

// DPTBool represents KNX DPT 1.xxx (boolean).
type DPTBool struct {
	Value bool
}

// Get returns the boolean value as a string.
func (b *DPTBool) Get() string {
	return fmt.Sprintf("%t", b.Value)
}

// Set parses a boolean string and stores the value.
func (b *DPTBool) Set(s string) error {
	switch strings.ToLower(s) {
	case "on", "enable":
		b.Value = true
	case "off", "disable":
		b.Value = false
	default:
		value, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		b.Value = value
	}
	return nil
}

// TypeCode returns the KNX type code for DPTBool.
func (b *DPTBool) TypeCode() TypeCode {
	return BoolCode
}

// String returns a human-readable representation of DPTBool.
func (b *DPTBool) String() string {
	return fmt.Sprintf("dptBool[%t]", b.Value)
}

// MarshalBinary encodes DPTBool into a 1-byte payload.
func (b *DPTBool) MarshalBinary() ([]byte, error) {
	result := []byte{0}
	if b.Value {
		result[0] = 1
	}
	return result, nil
}

// UnmarshalBinary decodes a 1-byte payload into DPTBool.
func (b *DPTBool) UnmarshalBinary(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("%w: type %s requires 1 byte, got %d", ErrUnmarshal, BoolCode.String(), len(data))
	}
	b.Value = data[0]&0x01 != 0
	return nil
}
