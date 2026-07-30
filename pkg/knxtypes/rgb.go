package knxtypes

import (
	"fmt"
	"strconv"
	"strings"
)

// DPTRedGreenBlue represents KNX DPT 232.600 (RGB color).
type DPTRedGreenBlue struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

// Get returns the color components as a comma-separated string.
func (d *DPTRedGreenBlue) Get() string {
	return fmt.Sprintf("%d,%d,%d", d.Red, d.Green, d.Blue)
}

// Set parses a "R,G,B" string and stores the color components.
func (d *DPTRedGreenBlue) Set(value string) error {
	parts := strings.Split(value, ",")
	if len(parts) != 3 {
		return fmt.Errorf("%w redGreenBlue '%s'", ErrBadValue, value)
	}
	numbers := [3]uint8{0, 0, 0}
	for i, part := range parts {
		v, err := strconv.ParseUint(part, 10, 8)
		if err != nil {
			return err
		}
		numbers[i] = uint8(v)
	}
	d.Red, d.Green, d.Blue = numbers[0], numbers[1], numbers[2]
	return nil
}

// TypeCode returns the KNX type code for DPTRedGreenBlue.
func (d *DPTRedGreenBlue) TypeCode() TypeCode {
	return RedGreenBlueCode
}

// String returns a human-readable representation of DPTRedGreenBlue.
func (d *DPTRedGreenBlue) String() string {
	return fmt.Sprintf("dptRedGreenBlue[%d,%d,%d]", d.Red, d.Green, d.Blue)
}

// MarshalBinary encodes DPTRedGreenBlue into a 3-byte payload.
func (d *DPTRedGreenBlue) MarshalBinary() ([]byte, error) {
	return []byte{d.Red, d.Green, d.Blue}, nil
}

// UnmarshalBinary decodes a 3-byte payload into DPTRedGreenBlue.
func (d *DPTRedGreenBlue) UnmarshalBinary(data []byte) error {
	if len(data) != 3 {
		return fmt.Errorf("%w: type %s requires 3 bytes, got %d", ErrUnmarshal, RedGreenBlueCode.String(), len(data))
	}
	d.Red = data[0]
	d.Green = data[1]
	d.Blue = data[2]
	return nil
}
