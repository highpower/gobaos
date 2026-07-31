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
	BoolCode         TypeCode = 1
	ControlCode      TypeCode = 2
	StepControlCode  TypeCode = 3
	CharCode         TypeCode = 4
	Uint8Code        TypeCode = 5
	Int8Code         TypeCode = 6
	Uint16Code       TypeCode = 7
	Int16Code        TypeCode = 8
	FloatCode        TypeCode = 9
	TimeOfDayCode    TypeCode = 10
	DateCode         TypeCode = 11
	Uint32Code       TypeCode = 12
	Int32Code        TypeCode = 13
	Float4ByteCode     TypeCode = 14
	AccessCode         TypeCode = 15
	StringCode         TypeCode = 16
	SceneNumberCode    TypeCode = 17
	SceneControlCode   TypeCode = 18
	DateTime8ByteCode  TypeCode = 19
	EnumCode           TypeCode = 20
	RedGreenBlueCode   TypeCode = 232
)

// String returns the textual name of the type code.
func (t TypeCode) String() string {
	switch t {
	case BoolCode:
		return "Bool"
	case ControlCode:
		return "Control"
	case StepControlCode:
		return "StepControl"
	case CharCode:
		return "Char"
	case Uint8Code:
		return "Uint8"
	case Uint16Code:
		return "Uint16"
	case Int16Code:
		return "Int16"
	case Int8Code:
		return "Int8"
	case FloatCode:
		return "Float"
	case TimeOfDayCode:
		return "TimeOfDay"
	case DateCode:
		return "Date"
	case Uint32Code:
		return "Uint32"
	case Int32Code:
		return "Int32"
	case Float4ByteCode:
		return "Float4Byte"
	case AccessCode:
		return "Access"
	case StringCode:
		return "String"
	case SceneNumberCode:
		return "SceneNumber"
	case SceneControlCode:
		return "SceneControl"
	case DateTime8ByteCode:
		return "DateTime8Byte"
	case EnumCode:
		return "Enum"
	case RedGreenBlueCode:
		return "RedGreenBlue"
	}
	return "Unknown"
}

// ParseTypeCode converts a string name to a TypeCode.
func ParseTypeCode(s string) (TypeCode, error) {
	switch strings.ToLower(s) {
	case "bool":
		return BoolCode, nil
	case "control":
		return ControlCode, nil
	case "stepcontrol":
		return StepControlCode, nil
	case "char":
		return CharCode, nil
	case "uint8":
		return Uint8Code, nil
	case "int8":
		return Int8Code, nil
	case "uint16":
		return Uint16Code, nil
	case "int16":
		return Int16Code, nil
	case "float":
		return FloatCode, nil
	case "timeofday":
		return TimeOfDayCode, nil
	case "date":
		return DateCode, nil
	case "uint32":
		return Uint32Code, nil
	case "int32":
		return Int32Code, nil
	case "float4byte":
		return Float4ByteCode, nil
	case "access":
		return AccessCode, nil
	case "string":
		return StringCode, nil
	case "scenenumber":
		return SceneNumberCode, nil
	case "scenecontrol":
		return SceneControlCode, nil
	case "datetime8byte":
		return DateTime8ByteCode, nil
	case "enum":
		return EnumCode, nil
	case "rgb", "redgreenblue":
		return RedGreenBlueCode, nil
	}
	return TypeCode(0), fmt.Errorf("%w: %s", ErrUnknownType, s)
}
