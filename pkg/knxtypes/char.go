package knxtypes

import "fmt"

type DPTChar struct {
	Value rune
}

func (c *DPTChar) Get() string {
	return fmt.Sprintf("%c", c.Value)
}

func (c *DPTChar) Set(s string) error {
	//TODO implement me
	panic("implement me")
}

func (c *DPTChar) TypeCode() TypeCode {
	return Char
}

func (c *DPTChar) BytesRequired() int {
	return 1
}

func (c *DPTChar) MarshalBinary() (data []byte, err error) {
	if c.Value < 0 || c.Value > 255 {
		return nil, fmt.Errorf("%w of char: %c", ErrBadValue, c.Value)
	}
	return []byte{byte(c.Value)}, nil
}

func (c *DPTChar) UnmarshalBinary(data []byte) error {
	//TODO implement me
	panic("implement me")
}
