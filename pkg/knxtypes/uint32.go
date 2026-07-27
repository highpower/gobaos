package knxtypes

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

// DPTUint32 represents KNX DPT 12.xxx (unsigned 32-bit integer).
type DPTUint32 struct {
	Value uint32
}

// Get returns the value as a decimal string.
func (d *DPTUint32) Get() string {
	return fmt.Sprintf("%d", d.Value)
}

// Set parses a decimal string and stores the value.
func (d *DPTUint32) Set(s string) error {
	value, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return err
	}
	d.Value = uint32(value)
	return nil
}

// TypeCode returns the type code for DPTUint32.
func (d *DPTUint32) TypeCode() TypeCode {
	return Uint32Code
}

// String returns a human-readable representation of DPTUint32.
func (d *DPTUint32) String() string {
	return fmt.Sprintf("dptUInt32[%d]", d.Value)
}

// MarshalBinary encodes DPTUint32 into a 4-byte big-endian payload.
func (d *DPTUint32) MarshalBinary() ([]byte, error) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, d.Value)
	return b, nil
}

// UnmarshalBinary decodes a 4-byte big-endian payload into DPTUint32.
func (d *DPTUint32) UnmarshalBinary(data []byte) error {
	if len(data) != 4 {
		return fmt.Errorf("%w: type %s requires 4 bytes, got %d", ErrUnmarshal, Uint32Code.String(), len(data))
	}
	d.Value = binary.BigEndian.Uint32(data)
	return nil
}
