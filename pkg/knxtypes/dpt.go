package knxtypes

import (
	"encoding"
	"fmt"
	"strings"
)

// TypeCode identifies a KNX datapoint type.
type TypeCode int

// DPT is the common interface for all KNX datapoint types.
type DPT interface {
	Get() string
	Set(value string) error
	TypeCode() TypeCode
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

// Supported KNX datapoint type codes.
const (
	Bool         TypeCode = 1
	Control      TypeCode = 2
	StepControl  TypeCode = 3
	Char         TypeCode = 4
	Uint8        TypeCode = 5
	Int8         TypeCode = 6
	Uint16       TypeCode = 7
	Int16        TypeCode = 8
	Float        TypeCode = 9
	TimeOfDay    TypeCode = 10
	Date         TypeCode = 11
	Uint32       TypeCode = 12
	Int32        TypeCode = 13
	Float4Byte   TypeCode = 14
	RedGreenBlue TypeCode = 232
)

// String returns the textual name of the type code.
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
	case Uint8:
		return "Uint8"
	case Uint16:
		return "Uint16"
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
	case Uint32:
		return "Uint32"
	case Int32:
		return "Int32"
	case Float4Byte:
		return "Float4Byte"
	case RedGreenBlue:
		return "RedGreenBlue"
	}
	return "Unknown"
}

// ParseTypeCode converts a string name to a TypeCode.
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
		return Uint8, nil
	case "int8":
		return Int8, nil
	case "uint16":
		return Uint16, nil
	case "int16":
		return Int16, nil
	case "float":
		return Float, nil
	case "timeofday":
		return TimeOfDay, nil
	case "date":
		return Date, nil
	case "uint32":
		return Uint32, nil
	case "int32":
		return Int32, nil
	case "float4byte":
		return Float4Byte, nil
	case "rgb", "redgreenblue":
		return RedGreenBlue, nil
	}
	return TypeCode(0), fmt.Errorf("%w: %s", ErrUnknownType, s)
}
