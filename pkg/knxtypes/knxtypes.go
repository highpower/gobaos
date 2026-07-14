package knxtypes

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

type DPTControl struct {
	Control bool
	Value   bool
}

type DPTStepControl struct {
	Control  bool
	StepCode uint8
}

//goland:noinspection GoMixedReceiverTypes
func (d DPTControl) MarshalBinary() ([]byte, error) {
	result := []byte{0}
	if d.Control {
		result[0] |= 0x02
	}
	if d.Value {
		result[0] |= 0x01
	}
	return result, nil
}

//goland:noinspection GoMixedReceiverTypes
func (d *DPTControl) UnmarshalBinary(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("DPTControl %w: expected 1 byte, got %d", ErrUnmarshal, len(data))
	}
	d.Control = data[0]&0x02 != 0
	d.Value = data[0]&0x01 != 0
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (d DPTStepControl) MarshalBinary() ([]byte, error) {
	if d.StepCode > 7 {
		return nil, fmt.Errorf("DPTStepControl %w: step code %d out of range [0,7]", ErrUnmarshal, d.StepCode)
	}
	var b byte
	if d.Control {
		b |= 0x08
	}
	b |= d.StepCode & 0x07
	return []byte{b}, nil
}

//goland:noinspection GoMixedReceiverTypes
func (d *DPTStepControl) UnmarshalBinary(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("DPTStepControl %w: expected 1 byte, got %d", ErrUnmarshal, len(data))
	}
	d.Control = data[0]&0x08 != 0
	d.StepCode = data[0] & 0x07
	return nil
}

// DPT5 represents KNX DPT 5.xxx (1-byte unsigned).
type DPT5 uint8

// Marshal encodes DPT5 into a 1-byte payload.
func (d DPT5) Marshal() ([]byte, error) {
	return []byte{byte(d)}, nil
}

// Unmarshal decodes a 1-byte payload into DPT5.
func (d *DPT5) Unmarshal(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("DPT5 %w: expected 1 byte, got %d", ErrUnmarshal, len(data))
	}
	*d = DPT5(data[0])
	return nil
}

// DPT6 represents KNX DPT 6.xxx (1-byte signed).
type DPT6 int8

// Marshal encodes DPT6 into a 1-byte payload.
func (d DPT6) Marshal() ([]byte, error) {
	return []byte{byte(int8(d))}, nil
}

// Unmarshal decodes a 1-byte payload into DPT6.
func (d *DPT6) Unmarshal(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("DPT6 %w: expected 1 byte, got %d", ErrUnmarshal, len(data))
	}
	*d = DPT6(int8(data[0]))
	return nil
}

// DPT7 represents KNX DPT 7.xxx (2-byte unsigned).
type DPT7 uint16

// Marshal encodes DPT7 into a 2-byte big-endian payload.
func (d DPT7) Marshal() ([]byte, error) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(d))
	return b, nil
}

// Unmarshal decodes a 2-byte big-endian payload into DPT7.
func (d *DPT7) Unmarshal(data []byte) error {
	if len(data) != 2 {
		return fmt.Errorf("DPT7 %w: expected 2 bytes, got %d", ErrUnmarshal, len(data))
	}
	*d = DPT7(binary.BigEndian.Uint16(data))
	return nil
}

// DPT8 represents KNX DPT 8.xxx (2-byte signed).
type DPT8 int16

// Marshal encodes DPT8 into a 2-byte big-endian payload.
func (d DPT8) Marshal() ([]byte, error) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(int16(d)))
	return b, nil
}

// Unmarshal decodes a 2-byte big-endian payload into DPT8.
func (d *DPT8) Unmarshal(data []byte) error {
	if len(data) != 2 {
		return fmt.Errorf("DPT8 %w: expected 2 bytes, got %d", ErrUnmarshal, len(data))
	}
	*d = DPT8(int16(binary.BigEndian.Uint16(data)))
	return nil
}

// DPT9 represents KNX DPT 9.xxx (2-byte float, E5 mantissa).
type DPT9 struct {
	Value float64
}

// Marshal encodes DPT9 into a 2-byte float payload.
func (d *DPT9) Marshal() ([]byte, error) {
	if math.IsNaN(d.Value) || math.IsInf(d.Value, 0) {
		return nil, fmt.Errorf("DPT9 %w: cannot encode NaN or Inf", ErrUnmarshal)
	}

	v := d.Value
	exp := 0

	// Scale value into 11-bit signed mantissa range (-2048..2047).
	for v < -2048.0 || v > 2047.0 {
		if exp >= 15 {
			return nil, fmt.Errorf("DPT9 %w: value %v out of range", ErrUnmarshal, d.Value)
		}
		v /= 2.0
		exp++
	}
	for v != 0 && v > -1.0 && v < 1.0 {
		if exp <= -15 {
			return []byte{0, 0}, nil
		}
		v *= 2.0
		exp--
	}

	mant := int16(math.Round(v))
	if mant < -2048 || mant > 2047 {
		return nil, fmt.Errorf("DPT9 %w: mantissa %d out of range", ErrUnmarshal, mant)
	}

	var raw uint16
	if mant < 0 {
		raw = uint16(mant + 2048)
		raw |= 0x8000
	} else {
		raw = uint16(mant)
	}
	raw |= (uint16(exp) & 0x0F) << 11

	return []byte{byte(raw >> 8), byte(raw)}, nil
}

// Unmarshal decodes a 2-byte float payload into DPT9.
func (d *DPT9) Unmarshal(data []byte) error {
	if len(data) != 2 {
		return fmt.Errorf("DPT9 %w: expected 2 bytes, got %d", ErrUnmarshal, len(data))
	}

	raw := uint16(data[0])<<8 | uint16(data[1])
	sign := raw & 0x8000
	exp := int((raw & 0x7800) >> 11)
	mant := int(raw & 0x07FF)

	if sign != 0 {
		mant = mant - 2048
	}

	d.Value = float64(mant) * math.Pow(2, float64(exp))
	return nil
}

// DPT10 represents KNX DPT 10.xxx (Time of day).
type DPT10 struct {
	Weekday uint8 // 0=no day, 1=Monday..7=Sunday
	Hour    uint8 // 0..23
	Minute  uint8 // 0..59
	Second  uint8 // 0..59
}

// Marshal encodes DPT10 into a 3-byte payload.
func (d *DPT10) Marshal() ([]byte, error) {
	if d.Weekday > 7 {
		return nil, fmt.Errorf("DPT10 %w: weekday %d out of range [0,7]", ErrUnmarshal, d.Weekday)
	}
	if d.Hour > 23 {
		return nil, fmt.Errorf("DPT10 %w: hour %d out of range [0,23]", ErrUnmarshal, d.Hour)
	}
	if d.Minute > 59 {
		return nil, fmt.Errorf("DPT10 %w: minute %d out of range [0,59]", ErrUnmarshal, d.Minute)
	}
	if d.Second > 59 {
		return nil, fmt.Errorf("DPT10 %w: second %d out of range [0,59]", ErrUnmarshal, d.Second)
	}
	return []byte{
		(d.Weekday&0x07)<<5 | (d.Hour & 0x1F),
		d.Minute & 0x3F,
		d.Second & 0x3F,
	}, nil
}

// Unmarshal decodes a 3-byte payload into DPT10.
func (d *DPT10) Unmarshal(data []byte) error {
	if len(data) != 3 {
		return fmt.Errorf("DPT10 %w: expected 3 bytes, got %d", ErrUnmarshal, len(data))
	}
	d.Weekday = (data[0] >> 5) & 0x07
	d.Hour = data[0] & 0x1F
	d.Minute = data[1] & 0x3F
	d.Second = data[2] & 0x3F
	return nil
}

// DPT11 represents KNX DPT 11.xxx (Date).
type DPT11 struct {
	Day   uint8 // 1..31
	Month uint8 // 1..12
	Year  uint8 // 0..99
}

// Marshal encodes DPT11 into a 3-byte payload.
func (d *DPT11) Marshal() ([]byte, error) {
	if d.Day < 1 || d.Day > 31 {
		return nil, fmt.Errorf("DPT11 %w: day %d out of range [1,31]", ErrUnmarshal, d.Day)
	}
	if d.Month < 1 || d.Month > 12 {
		return nil, fmt.Errorf("DPT11 %w: month %d out of range [1,12]", ErrUnmarshal, d.Month)
	}
	return []byte{d.Day, d.Month, d.Year}, nil
}

// Unmarshal decodes a 3-byte payload into DPT11.
func (d *DPT11) Unmarshal(data []byte) error {
	if len(data) != 3 {
		return fmt.Errorf("DPT11 %w: expected 3 bytes, got %d", ErrUnmarshal, len(data))
	}
	d.Day = data[0]
	d.Month = data[1]
	d.Year = data[2]
	return nil
}

// DPT12 represents KNX DPT 12.xxx (4-byte unsigned).
type DPT12 uint32

// Marshal encodes DPT12 into a 4-byte big-endian payload.
func (d DPT12) Marshal() ([]byte, error) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(d))
	return b, nil
}

// Unmarshal decodes a 4-byte big-endian payload into DPT12.
func (d *DPT12) Unmarshal(data []byte) error {
	if len(data) != 4 {
		return fmt.Errorf("DPT12 %w: expected 4 bytes, got %d", ErrUnmarshal, len(data))
	}
	*d = DPT12(binary.BigEndian.Uint32(data))
	return nil
}

// DPT13 represents KNX DPT 13.xxx (4-byte signed).
type DPT13 int32

// Marshal encodes DPT13 into a 4-byte big-endian payload.
func (d DPT13) Marshal() ([]byte, error) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(int32(d)))
	return b, nil
}

// Unmarshal decodes a 4-byte big-endian payload into DPT13.
func (d *DPT13) Unmarshal(data []byte) error {
	if len(data) != 4 {
		return fmt.Errorf("DPT13 %w: expected 4 bytes, got %d", ErrUnmarshal, len(data))
	}
	*d = DPT13(int32(binary.BigEndian.Uint32(data)))
	return nil
}

// DPT14 represents KNX DPT 14.xxx (4-byte float, IEEE 754).
type DPT14 float64

// Marshal encodes DPT14 into a 4-byte IEEE 754 big-endian payload.
func (d DPT14) Marshal() ([]byte, error) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, math.Float32bits(float32(d)))
	return b, nil
}

// Unmarshal decodes a 4-byte IEEE 754 big-endian payload into DPT14.
func (d *DPT14) Unmarshal(data []byte) error {
	if len(data) != 4 {
		return fmt.Errorf("DPT14 %w: expected 4 bytes, got %d", ErrUnmarshal, len(data))
	}
	*d = DPT14(math.Float32frombits(binary.BigEndian.Uint32(data)))
	return nil
}

// DPT16 represents KNX DPT 16.xxx (14-byte string).
type DPT16 string

// Marshal encodes DPT16 into a 14-byte payload, zero-padded.
func (d DPT16) Marshal() ([]byte, error) {
	s := string(d)
	if len(s) > 14 {
		return nil, fmt.Errorf("DPT16 %w: string length %d exceeds 14 bytes", ErrUnmarshal, len(s))
	}
	b := make([]byte, 14)
	copy(b, s)
	return b, nil
}

// Unmarshal decodes a 14-byte payload into DPT16, trimming trailing nulls.
func (d *DPT16) Unmarshal(data []byte) error {
	if len(data) != 14 {
		return fmt.Errorf("DPT16 %w: expected 14 bytes, got %d", ErrUnmarshal, len(data))
	}
	*d = DPT16(strings.TrimRight(string(data), "\x00"))
	return nil
}

// DPT17 represents KNX DPT 17.xxx (scene number, 1..64).
type DPT17 uint8

// Marshal encodes DPT17 into a 1-byte payload.
func (d DPT17) Marshal() ([]byte, error) {
	if d < 1 || d > 64 {
		return nil, fmt.Errorf("DPT17 %w: scene number %d out of range [1,64]", ErrUnmarshal, d)
	}
	return []byte{byte(d - 1)}, nil
}

// Unmarshal decodes a 1-byte payload into DPT17.
func (d *DPT17) Unmarshal(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("DPT17 %w: expected 1 byte, got %d", ErrUnmarshal, len(data))
	}
	v := data[0] & 0x3F // mask reserved bits
	*d = DPT17(v + 1)
	return nil
}

// DPT18 represents KNX DPT 18.xxx (scene control).
type DPT18 struct {
	Save  bool  // true = learn/save, false = execute/recall
	Scene uint8 // 1..64
}

// Marshal encodes DPT18 into a 1-byte payload.
func (d *DPT18) Marshal() ([]byte, error) {
	if d.Scene < 1 || d.Scene > 64 {
		return nil, fmt.Errorf("DPT18 %w: scene number %d out of range [1,64]", ErrUnmarshal, d.Scene)
	}
	var b byte
	if d.Save {
		b |= 0x80
	}
	b |= (d.Scene - 1) & 0x3F
	return []byte{b}, nil
}

// Unmarshal decodes a 1-byte payload into DPT18.
func (d *DPT18) Unmarshal(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("DPT18 %w: expected 1 byte, got %d", ErrUnmarshal, len(data))
	}
	d.Save = data[0]&0x80 != 0
	d.Scene = (data[0] & 0x3F) + 1
	return nil
}

// DPT20 represents KNX DPT 20.xxx (1-byte enumeration).
type DPT20 uint8

// Marshal encodes DPT20 into a 1-byte payload.
func (d DPT20) Marshal() ([]byte, error) {
	return []byte{byte(d)}, nil
}

// Unmarshal decodes a 1-byte payload into DPT20.
func (d *DPT20) Unmarshal(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("DPT20 %w: expected 1 byte, got %d", ErrUnmarshal, len(data))
	}
	*d = DPT20(data[0])
	return nil
}

// DPT232 represents KNX DPT 232.600 (RGB color).
type DPT232 struct {
	R uint8
	G uint8
	B uint8
}

// Marshal encodes DPT232 into a 3-byte payload.
func (d *DPT232) Marshal() ([]byte, error) {
	return []byte{d.R, d.G, d.B}, nil
}

// Unmarshal decodes a 3-byte payload into DPT232.
func (d *DPT232) Unmarshal(data []byte) error {
	if len(data) != 3 {
		return fmt.Errorf("DPT232 %w: expected 3 bytes, got %d", ErrUnmarshal, len(data))
	}
	d.R = data[0]
	d.G = data[1]
	d.B = data[2]
	return nil
}
