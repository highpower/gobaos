package knxtypes

import (
	"encoding/binary"
	"fmt"
	"math"
)

type DPTFloat struct {
	Value float64
}

//goland:noinspection GoMixedReceiverTypes
func (d DPTFloat) String() string {
	return fmt.Sprintf("dptFloat[%0.4f]", d.Value)
}

//goland:noinspection GoMixedReceiverTypes
func (d DPTFloat) MarshalBinary() ([]byte, error) {
	if math.IsNaN(d.Value) || math.IsInf(d.Value, 0) {
		return nil, fmt.Errorf("%w: cannot encode NaN or Inf", ErrUnmarshal)
	}
	negative := false
	val := int(d.Value * 100)
	if val < 0 {
		val = -val
		negative = true
	}
	exp := 0
	for val > 2047 && exp < 7 {
		val >>= 1
		exp++
	}
	compressed := uint16((exp&0x0F)<<11) | uint16(val&0x7FF)
	if negative {
		compressed |= 0x800
	}
	result := make([]byte, 2)
	binary.BigEndian.PutUint16(result, compressed)
	return result, nil

}

//goland:noinspection GoMixedReceiverTypes
func (d *DPTFloat) UnmarshalBinary(data []byte) error {
	if len(data) != 2 {
		return fmt.Errorf("%w: expected 2 bytes, got %d", ErrUnmarshal, len(data))
	}
	compressed := binary.BigEndian.Uint16(data)
	sign := compressed & 0x8000
	exp := int((compressed & 0x7800) >> 11)
	mantissa := int(compressed & 0x07FF)

	d.Value = 0.01 * float64(mantissa) * math.Pow(2, float64(exp))
	if sign == 0x8000 {
		d.Value = -d.Value
	}
	return nil
}

func NewDPTFloat() DPTFloat {
	return DPTFloat{}
}
