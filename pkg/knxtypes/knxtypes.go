package knxtypes

import (
	"fmt"
	"strings"
)

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
