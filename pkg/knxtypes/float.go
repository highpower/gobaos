package knxtypes

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
)

// DPTFloat represents KNX DPT 9.xxx (2-byte float).
type DPTFloat struct {
	Value float64
}

// Get returns the value as a formatted decimal string.
func (f *DPTFloat) Get() string {
	return fmt.Sprintf("%0.4f", f.Value)
}

// Set parses a decimal string, validates it, and stores the value.
func (f *DPTFloat) Set(s string) error {
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	if err := validFloat(value); err != nil {
		return err
	}
	f.Value = value
	return nil
}

// TypeCode returns the type code for DPTFloat.
func (f *DPTFloat) TypeCode() TypeCode {
	return FloatCode
}

// String returns a human-readable representation of DPTFloat.
func (f *DPTFloat) String() string {
	return fmt.Sprintf("dptFloat[%0.4f]", f.Value)
}

// Valid checks whether the float value can be encoded.
func (f *DPTFloat) Valid() error {
	return validFloat(f.Value)
}

// MarshalBinary encodes DPTFloat into a 2-byte payload.
func (f *DPTFloat) MarshalBinary() ([]byte, error) {
	mantissa, exp, negative := floatExpansion(f.Value)
	compressed := uint16((exp&0x0F)<<11) | uint16(mantissa&0x7FF)
	if negative {
		compressed |= 0x800
	}
	result := make([]byte, 2)
	binary.BigEndian.PutUint16(result, compressed)
	return result, nil
}

// UnmarshalBinary decodes a 2-byte payload into DPTFloat.
func (f *DPTFloat) UnmarshalBinary(data []byte) error {
	if len(data) != 2 {
		return fmt.Errorf("%w: type %s requires 2 bytes, got %d", ErrUnmarshal, FloatCode.String(), len(data))
	}
	compressed := binary.BigEndian.Uint16(data)
	sign := compressed & 0x8000
	exp := int((compressed & 0x7800) >> 11)
	mantissa := int(compressed & 0x07FF)

	f.Value = 0.01 * float64(mantissa) * math.Pow(2, float64(exp))
	if sign == 0x8000 {
		f.Value = -f.Value
	}
	return nil
}

func floatExpansion(number float64) (int, int, bool) {
	negative := false
	mantissa := int(number * 100)
	if mantissa < 0 {
		mantissa = -mantissa
		negative = true
	}
	exp := 0
	for mantissa > 2047 && exp < 7 {
		mantissa >>= 1
		exp++
	}
	return mantissa, exp, negative
}

func validFloat(value float64) error {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return fmt.Errorf("%w: cannot encode NaN or Inf", ErrBadValue)
	}
	mantissa, _, _ := floatExpansion(value)
	if mantissa > 2047 {
		return fmt.Errorf("%w: too big %f", ErrBadValue, value)
	}
	return nil
}
