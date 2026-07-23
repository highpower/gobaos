package knxtypes

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
)

type DPTFloat4Byte struct {
	Value float64
}

func (d *DPTFloat4Byte) Get() string {
	return fmt.Sprintf("%f", d.Value)
}

func (d *DPTFloat4Byte) Set(s string) error {
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	d.Value = value
	return nil
}

func (d *DPTFloat4Byte) TypeCode() TypeCode {
	return Float4Byte
}

func (d *DPTFloat4Byte) String() string {
	return fmt.Sprintf("dptFloat4Byte[%f]", d.Value)
}

func (d *DPTFloat4Byte) MarshalBinary() ([]byte, error) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, math.Float32bits(float32(d.Value)))
	return b, nil
}

func (d *DPTFloat4Byte) UnmarshalBinary(data []byte) error {
	if len(data) != 4 {
		return fmt.Errorf("%w: type %s requires 4 bytes, got %d", ErrUnmarshal, Float4Byte.String(), len(data))
	}
	d.Value = float64(math.Float32frombits(binary.BigEndian.Uint32(data)))
	return nil
}
