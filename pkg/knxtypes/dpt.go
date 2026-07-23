package knxtypes

import (
	"encoding"
	"strings"
)

type TypeCode int

type DPT interface {
	Get() string
	Set(value string) error
	TypeCode() TypeCode
	BytesRequired() int
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

const (
	Bool        TypeCode = 1
	Control     TypeCode = 2
	StepControl TypeCode = 3
	Char        TypeCode = 4
	Float       TypeCode = 9
)

func (t TypeCode) String() string {
	switch t {
	case Bool:
		return "Bool"
	case Control:
		return "Control"
	case StepControl:
		return "StepControl"
	case Char:
		return "Char"
	case Float:
		return "Float"
	}
	return "Unknown"
}

func ParseTypeCode(s string) (TypeCode, error) {
	switch strings.ToLower(s) {
	case "bool":
		return Bool, nil
	case "control":
		return Control, nil
	case "stepcontrol":
		return StepControl, nil
	case "char":
		return Char, nil
	case "float":
		return Float, nil
	}
	return TypeCode(0), nil
}
