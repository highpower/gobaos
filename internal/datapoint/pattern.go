package datapoint

import (
	"bytes"
	"fmt"
	"strings"
)

type pattern struct {
	index int
	parts []string
}

func (p *pattern) apply(str string) string {
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

func parsePattern(str string) (pattern, error) {
	if str == "" {
		return pattern{}, fmt.Errorf("%w: %s", ErrBadPattern, str)
	}
	index := -1
	parts := strings.Split(str, "/")
	for i, part := range parts {
		switch {
		case part == "+" && index == -1:
			index = i
			continue
		case part == "+":
			return pattern{}, fmt.Errorf("%w: %s contains more than 1 '+' token", ErrBadPattern, str)
		}
	}
	if index == -1 {
		return pattern{}, fmt.Errorf("%w: '%s' contains no wildcard", ErrBadPattern, str)
	}
	result := pattern{index: index, parts: parts}
	return result, nil
}
