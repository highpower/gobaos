package knxtypes

import (
	"fmt"
	"strconv"
)

// DPTSceneNumber represents KNX DPT 17.xxx (scene number, 1..64).
type DPTSceneNumber struct {
	Value uint8
}

// Get returns the scene number as a decimal string.
func (n *DPTSceneNumber) Get() string {
	return fmt.Sprintf("%d", n.Value)
}

// Set parses a decimal string and sets the scene number after validating the range.
func (n *DPTSceneNumber) Set(s string) error {
	value, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return err
	}
	if err := validSceneNumber(uint8(value)); err != nil {
		return err
	}
	n.Value = uint8(value)
	return nil
}

// TypeCode returns the KNX type code for DPTSceneNumber.
func (n *DPTSceneNumber) TypeCode() TypeCode {
	return SceneNumberCode
}

// String returns a human-readable representation of DPTSceneNumber.
func (n *DPTSceneNumber) String() string {
	return fmt.Sprintf("dptSceneNumber[%d]", n.Value)
}

// MarshalBinary encodes DPTSceneNumber into a 1-byte payload.
func (n *DPTSceneNumber) MarshalBinary() ([]byte, error) {
	return []byte{n.Value - 1}, nil
}

// UnmarshalBinary decodes a 1-byte payload into DPTSceneNumber.
func (n *DPTSceneNumber) UnmarshalBinary(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("%w: type %s requires 1 byte, got %d", ErrUnmarshal, SceneNumberCode.String(), len(data))
	}
	v := data[0] & 0x3F // mask reserved bits
	n.Value = v + 1
	return nil
}

func validSceneNumber(value uint8) error {
	if value == 0 || value > 64 {
		return fmt.Errorf("%w: scene number %d out of range [1,64]", ErrBadValue, value)
	}
	return nil
}
