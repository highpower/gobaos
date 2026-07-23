package knxtypes

import (
	"encoding"
	"fmt"
	"strings"
)

type TypeCode int

type DPT interface {
	Get() string
	Set(value string) error
	TypeCode() TypeCode
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

const (
	Bool        TypeCode = 1
	Control     TypeCode = 2
	StepControl TypeCode = 3
	Char        TypeCode = 4
	UInt8       TypeCode = 5
	Int8        TypeCode = 6
	UInt16      TypeCode = 7
	Int16       TypeCode = 8
	Float       TypeCode = 9
	TimeOfDay   TypeCode = 10
	Date        TypeCode = 11
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
	case UInt8:
		return "UInt8"
	case UInt16:
		return "UInt16"
	case Int16:
		return "Int16"
	case Int8:
		return "Int8"
	case Float:
		return "Float"
	case TimeOfDay:
		return "TimeOfDay"
	case Date:
		return "Date"
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
	case "uint8":
		return UInt8, nil
	case "int8":
		return Int8, nil
	case "uint16":
		return UInt16, nil
	case "int16":
		return Int16, nil
	case "float":
		return Float, nil
	case "timeofday":
		return TimeOfDay, nil
	case "date":
		return Date, nil
	}
	return TypeCode(0), fmt.Errorf("%w: %s", ErrUnknownType, s)
}
