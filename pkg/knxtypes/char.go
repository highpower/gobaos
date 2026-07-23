package knxtypes

import (
	"fmt"
	"unicode/utf8"
)

type DPTChar struct {
	Value rune
}

func (c *DPTChar) Get() string {
	return fmt.Sprintf("%c", c.Value)
}

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

func (c *DPTChar) TypeCode() TypeCode {
	return Char
}

func (c *DPTChar) Valid() error {
	return validChar(c.Value)
}

func (c *DPTChar) String() string {
	return fmt.Sprintf("dptChar[%c]", c.Value)
}

func (c *DPTChar) MarshalBinary() (data []byte, err error) {
	return []byte{byte(c.Value)}, nil
}

func (c *DPTChar) UnmarshalBinary(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("%w: type %s requires 1 byte, got %d", ErrUnmarshal, Char.String(), len(data))
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
