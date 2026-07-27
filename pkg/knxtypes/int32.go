package knxtypes

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

// DPTInt32 represents KNX DPT 13.xxx (signed 32-bit integer).
type DPTInt32 struct {
	Value int32
}

// Get returns the value as a decimal string.
func (d *DPTInt32) Get() string {
	return fmt.Sprintf("%d", d.Value)
}

// Set parses a decimal string and stores the value.
func (d *DPTInt32) Set(s string) error {
	value, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return err
	}
	d.Value = int32(value)
	return nil
}

// TypeCode returns the type code for DPTInt32.
func (d *DPTInt32) TypeCode() TypeCode {
	return Int32Code
}

// String returns a human-readable representation of DPTInt32.
func (d *DPTInt32) String() string {
	return fmt.Sprintf("dptInt32[%d]", d.Value)
}

// MarshalBinary encodes DPTInt32 into a 4-byte big-endian payload.
func (d *DPTInt32) MarshalBinary() ([]byte, error) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(d.Value))
	return b, nil
}

// UnmarshalBinary decodes a 4-byte big-endian payload into DPTInt32.
func (d *DPTInt32) UnmarshalBinary(data []byte) error {
	if len(data) != 4 {
		return fmt.Errorf("%w: type %s requires 4 bytes, got %d", ErrUnmarshal, Int32Code.String(), len(data))
	}
	d.Value = int32(binary.BigEndian.Uint32(data))
	return nil
}
