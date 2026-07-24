package knxtypes

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

// DPTUint16 represents KNX DPT 7.xxx (unsigned 16-bit integer).
type DPTUint16 struct {
	Value uint16
}

// Get returns the value as a decimal string.
func (d *DPTUint16) Get() string {
	return fmt.Sprintf("%d", d.Value)
}

// Set parses a decimal string and stores the value.
func (d *DPTUint16) Set(s string) error {
	value, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return err
	}
	d.Value = uint16(value)
	return nil
}

// TypeCode returns the type code for DPTUint16.
func (d *DPTUint16) TypeCode() TypeCode {
	return Uint16
}

// String returns a human-readable representation of DPTUint16.
func (d *DPTUint16) String() string {
	return fmt.Sprintf("dptUInt16[%d]", d.Value)
}

// MarshalBinary encodes DPTUint16 into a 2-byte big-endian payload.
func (d *DPTUint16) MarshalBinary() ([]byte, error) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, d.Value)
	return b, nil
}

// UnmarshalBinary decodes a 2-byte big-endian payload into DPTUint16.
func (d *DPTUint16) UnmarshalBinary(data []byte) error {
	if len(data) != 2 {
		return fmt.Errorf("%w: type %s requires 2 bytes, got %d", ErrUnmarshal, Uint16.String(), len(data))
	}
	d.Value = binary.BigEndian.Uint16(data)
	return nil
}
