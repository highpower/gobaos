package knxtypes

import (
	"fmt"
	"strconv"
)

// DPTInt8 represents KNX DPT 6.xxx (signed 8-bit integer).
type DPTInt8 struct {
	Value int8
}

// Get returns the value as a decimal string.
func (d *DPTInt8) Get() string {
	return fmt.Sprintf("%d", d.Value)
}

// Set parses a decimal string and stores the value.
func (d *DPTInt8) Set(s string) error {
	value, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		return err
	}
	d.Value = int8(value)
	return nil
}

// TypeCode returns the type code for DPTInt8.
func (d *DPTInt8) TypeCode() TypeCode {
	return Int8
}

// String returns a human-readable representation of DPTInt8.
func (d *DPTInt8) String() string {
	return fmt.Sprintf("dptInt8[%d]", d.Value)
}

// MarshalBinary encodes DPTInt8 into a 1-byte payload.
func (d *DPTInt8) MarshalBinary() ([]byte, error) {
	return []byte{byte(d.Value)}, nil
}

// UnmarshalBinary decodes a 1-byte payload into DPTInt8.
func (d *DPTInt8) UnmarshalBinary(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("%w: type %s requires 1 byte, got %d", ErrUnmarshal, Int8.String(), len(data))
	}
	d.Value = int8(data[0])
	return nil
}
