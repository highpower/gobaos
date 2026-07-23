package knxtypes

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
)

type DPTFloat struct {
	Value float64
}

func (f *DPTFloat) Get() string {
	return fmt.Sprintf("%0.4f", f.Value)
}

func (f *DPTFloat) Set(s string) error {
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	f.Value = value
	return nil
}

func (f *DPTFloat) TypeCode() TypeCode {
	return Float
}

func (f *DPTFloat) BytesRequired() int {
	return 2
}

func (f *DPTFloat) String() string {
	return fmt.Sprintf("dptFloat[%0.4f]", f.Value)
}

func (f *DPTFloat) MarshalBinary() ([]byte, error) {
	if math.IsNaN(f.Value) || math.IsInf(f.Value, 0) {
		return nil, fmt.Errorf("%w: cannot encode NaN or Inf", ErrUnmarshal)
	}
	mantissa, exp, negative := floatExpansion(f.Value)
	if mantissa > 2047 {
		return nil, fmt.Errorf("%w: too big %f", ErrUnmarshal, f.Value)
	}
	compressed := uint16((exp&0x0F)<<11) | uint16(mantissa&0x7FF)
	if negative {
		compressed |= 0x800
	}
	result := make([]byte, 2)
	binary.BigEndian.PutUint16(result, compressed)
	return result, nil

}

func (f *DPTFloat) UnmarshalBinary(data []byte) error {
	if len(data) != 2 {
		return fmt.Errorf("%w: expected 2 bytes, got %d", ErrUnmarshal, len(data))
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
