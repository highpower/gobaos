package knxtypes

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
)

// DPTFloat4Byte represents KNX DPT 14.xxx (4-byte IEEE 754 float).
type DPTFloat4Byte struct {
	Value float64
}

// Get returns the value as a decimal string.
func (d *DPTFloat4Byte) Get() string {
	return fmt.Sprintf("%f", d.Value)
}

// Set parses a decimal string and stores the value.
func (d *DPTFloat4Byte) Set(s string) error {
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	d.Value = value
	return nil
}

// TypeCode returns the type code for DPTFloat4Byte.
func (d *DPTFloat4Byte) TypeCode() TypeCode {
	return Float4ByteCode
}

// String returns a human-readable representation of DPTFloat4Byte.
func (d *DPTFloat4Byte) String() string {
	return fmt.Sprintf("dptFloat4Byte[%f]", d.Value)
}

// MarshalBinary encodes DPTFloat4Byte into a 4-byte IEEE 754 payload.
func (d *DPTFloat4Byte) MarshalBinary() ([]byte, error) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, math.Float32bits(float32(d.Value)))
	return b, nil
}

// UnmarshalBinary decodes a 4-byte IEEE 754 payload into DPTFloat4Byte.
func (d *DPTFloat4Byte) UnmarshalBinary(data []byte) error {
	if len(data) != 4 {
		return fmt.Errorf("%w: type %s requires 4 bytes, got %d", ErrUnmarshal, Float4ByteCode.String(), len(data))
	}
	d.Value = float64(math.Float32frombits(binary.BigEndian.Uint32(data)))
	return nil
}
