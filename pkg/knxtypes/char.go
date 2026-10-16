package knxtypes

import (
	"fmt"
	"unicode/utf8"
)

// DPTChar represents KNX DPT 4.xxx (ASCII character).
type DPTChar struct {
	Value rune
}

// Get returns the character as a string.
func (c *DPTChar) Get() string {
	return fmt.Sprintf("%c", c.Value)
}

// Set parses a single-character string and stores it.
func (c *DPTChar) Set(s string) error {
	if utf8.RuneCountInString(s) != 1 {
		return fmt.Errorf("%w: char value '%s'", ErrBadValue, s)
	}
	r, _ := utf8.DecodeRuneInString(s)
	if err := validChar(r); err != nil {
		return err
	}
	c.Value = r
	return nil
}

// Valid reports whether the character value is in the valid ASCII range.
func (c *DPTChar) Valid() error {
	return validChar(c.Value)
}

// TypeCode returns the KNX type code for DPTChar.
func (c *DPTChar) TypeCode() TypeCode {
	return CharCode
}

// String returns a human-readable representation of DPTChar.
func (c *DPTChar) String() string {
	return fmt.Sprintf("dptChar[%c]", c.Value)
}

// MarshalBinary encodes DPTChar into a 1-byte payload.
func (c *DPTChar) MarshalBinary() ([]byte, error) {
	return []byte{byte(c.Value)}, nil
}

// UnmarshalBinary decodes a 1-byte payload into DPTChar.
func (c *DPTChar) UnmarshalBinary(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("%w: type %s requires 1 byte, got %d", ErrUnmarshal, CharCode.String(), len(data))
	}
	c.Value = rune(data[0])
	return nil
}

func validChar(value rune) error {
	if value < 0 || value > 255 {
		return fmt.Errorf("%w of char: %c", ErrBadValue, value)
	}
	return nil
}
