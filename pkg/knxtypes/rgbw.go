package knxtypes

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

// DPTRedGreenBlueWhite represents KNX DPT 251.xxx (RGBW color).
type DPTRedGreenBlueWhite struct {
	Red   uint8
	Green uint8
	Blue  uint8
	White uint8
	Flags uint16
}

// Get returns the color components as a comma-separated string.
func (d *DPTRedGreenBlueWhite) Get() string {
	return fmt.Sprintf("%d,%d,%d,%d,%04x", d.Red, d.Green, d.Blue, d.White, d.Flags)
}

// Set parses a "R,G,B,W" string and stores the color components.
func (d *DPTRedGreenBlueWhite) Set(value string) error {
	parts := strings.Split(value, ",")
	if len(parts) != 5 {
		return fmt.Errorf("%w redGreenBlueWhite '%s'", ErrBadValue, value)
	}
	numbers := [4]uint8{0, 0, 0, 0}
	for i := 0; i < 4; i++ {
		v, err := strconv.ParseUint(parts[i], 10, 8)
		if err != nil {
			return err
		}
		numbers[i] = uint8(v)
	}
	flags, err := strconv.ParseUint(parts[4], 16, 16)
	if err != nil {
		return err
	}
	d.Red, d.Green, d.Blue, d.White, d.Flags = numbers[0], numbers[1], numbers[2], numbers[3], uint16(flags)
	return nil
}

// TypeCode returns the KNX type code for DPTRedGreenBlueWhite.
func (d *DPTRedGreenBlueWhite) TypeCode() TypeCode {
	return RedGreenBlueWhiteCode
}

// String returns a human-readable representation of DPTRedGreenBlueWhite.
func (d *DPTRedGreenBlueWhite) String() string {
	return fmt.Sprintf("dptRedGreenBlueWhite[%d,%d,%d,%d,%04x]", d.Red, d.Green, d.Blue, d.White, d.Flags)
}

// MarshalBinary encodes DPTRedGreenBlueWhite into a 4-byte payload.
func (d *DPTRedGreenBlueWhite) MarshalBinary() ([]byte, error) {
	result := []byte{d.Red, d.Green, d.Blue, d.White, 0, 0}
	binary.BigEndian.PutUint16(result[4:6], d.Flags)
	return result, nil
}

// UnmarshalBinary decodes a 4-byte payload into DPTRedGreenBlueWhite.
func (d *DPTRedGreenBlueWhite) UnmarshalBinary(data []byte) error {
	if len(data) != 6 {
		return fmt.Errorf("%w: type %s requires 6 bytes, got %d", ErrUnmarshal, RedGreenBlueWhiteCode.String(),
			len(data))
	}
	d.Red = data[0]
	d.Green = data[1]
	d.Blue = data[2]
	d.White = data[3]
	d.Flags = binary.BigEndian.Uint16(data[4:6])
	return nil
}
