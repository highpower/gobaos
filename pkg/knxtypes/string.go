package knxtypes

import (
	"fmt"
	"strings"
)

// DPTString represents KNX DPT 16.xxx (14-byte string).
type DPTString struct {
	Value string
}

// Get returns the current string value.
func (s *DPTString) Get() string {
	return s.Value
}

// Set validates and sets the string value.
func (s *DPTString) Set(value string) error {
	if err := validString(value); err != nil {
		return err
	}
	s.Value = value
	return nil
}

// TypeCode returns the KNX type code for DPTString.
func (s *DPTString) TypeCode() TypeCode {
	return StringCode
}

// String returns a human-readable representation of DPTString.
func (s *DPTString) String() string {
	return fmt.Sprintf("dptString['%s']", s.Value)
}

// Valid checks whether the current value is valid.
func (s *DPTString) Valid() error {
	return validString(s.Value)
}

// MarshalBinary encodes DPTString into a 14-byte payload, zero-padded.
func (s *DPTString) MarshalBinary() ([]byte, error) {
	b := make([]byte, 14)
	copy(b, s.Value)
	return b, nil
}

// UnmarshalBinary decodes a 14-byte payload into DPTString, trimming trailing nulls.
func (s *DPTString) UnmarshalBinary(data []byte) error {
	if len(data) != 14 {
		return fmt.Errorf("%w: type %s requires 14 bytes, got %d", ErrUnmarshal, StringCode.String(), len(data))
	}
	s.Value = strings.TrimRightFunc(string(data), func(r rune) bool {
		return r == 0
	})
	return nil
}

func validString(s string) error {
	if len(s) > 14 {
		return fmt.Errorf("%w: string length must be at most 14 bytes, got %d", ErrBadValue, len(s))
	}
	return nil
}
