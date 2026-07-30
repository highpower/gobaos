package knxtypes

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

// DPTInt16 represents KNX DPT 8.xxx (signed 16-bit integer).
type DPTInt16 struct {
	Value int16
}

// Get returns the value as a decimal string.
func (d *DPTInt16) Get() string {
	return fmt.Sprintf("%d", d.Value)
}

// Set parses a decimal string and stores the value.
func (d *DPTInt16) Set(s string) error {
	value, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		return err
	}
	d.Value = int16(value)
	return nil
}

// TypeCode returns the KNX type code for DPTInt16.
func (d *DPTInt16) TypeCode() TypeCode {
	return Int16Code
}

// String returns a human-readable representation of DPTInt16.
func (d *DPTInt16) String() string {
	return fmt.Sprintf("dptInt16[%d]", d.Value)
}

// MarshalBinary encodes DPTInt16 into a 2-byte big-endian payload.
func (d *DPTInt16) MarshalBinary() ([]byte, error) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(d.Value))
	return b, nil
}

// UnmarshalBinary decodes a 2-byte big-endian payload into DPTInt16.
func (d *DPTInt16) UnmarshalBinary(data []byte) error {
	if len(data) != 2 {
		return fmt.Errorf("%w: type %s requires 2 bytes, got %d", ErrUnmarshal, Int16Code.String(), len(data))
	}
	d.Value = int16(binary.BigEndian.Uint16(data))
	return nil
}
