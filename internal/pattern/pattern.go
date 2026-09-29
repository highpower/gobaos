package pattern

import (
	"bytes"
	"fmt"
	"strings"
)

type Parsed struct {
	index int
	parts []string
}

func (p *Parsed) Apply(str string) string {
	buffer := bytes.Buffer{}
	for i, part := range p.parts {
		if i == p.index {
			buffer.WriteString(str)
		} else {
			buffer.WriteString(part)
		}
		if i != len(p.parts)-1 {
			buffer.WriteString("/")
		}
	}
	return buffer.String()
}

func Parse(str string) (Parsed, error) {
	if str == "" {
		return Parsed{}, fmt.Errorf("%w: %s", ErrInvalid, str)
	}
	index := -1
	parts := strings.Split(str, "/")
	for i, part := range parts {
		switch {
		case part == "+" && index == -1:
			index = i
			continue
		case part == "+":
			return Parsed{}, fmt.Errorf("%w: %s contains more than 1 '+' token", ErrInvalid, str)
		}
	}
	if index == -1 {
		return Parsed{}, fmt.Errorf("%w: '%s' contains no wildcard", ErrInvalid, str)
	}
	result := Parsed{index: index, parts: parts}
	return result, nil
}
