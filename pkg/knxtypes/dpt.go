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
	BoolCode              TypeCode = 1
	ControlCode           TypeCode = 2
	StepControlCode       TypeCode = 3
	CharCode              TypeCode = 4
	Uint8Code             TypeCode = 5
	Int8Code              TypeCode = 6
	Uint16Code            TypeCode = 7
	Int16Code             TypeCode = 8
	FloatCode             TypeCode = 9
	TimeOfDayCode         TypeCode = 10
	DateCode              TypeCode = 11
	Uint32Code            TypeCode = 12
	Int32Code             TypeCode = 13
	Float4ByteCode        TypeCode = 14
	AccessCode            TypeCode = 15
	StringCode            TypeCode = 16
	SceneNumberCode       TypeCode = 17
	SceneControlCode      TypeCode = 18
	DateTime8ByteCode     TypeCode = 19
	EnumCode              TypeCode = 20
	RedGreenBlueCode      TypeCode = 232
	RedGreenBlueWhiteCode TypeCode = 251
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
	case RedGreenBlueWhiteCode:
		return "RedGreenBlueWhite"
	}
	return "Unknown"
}

// Valid returns nil if code is a valid KNX DPT type,
// otherwise it returns ErrUnknownType.
func (t TypeCode) Valid() error {
	switch t {
	case BoolCode, ControlCode, StepControlCode, CharCode,
		Uint8Code, Int8Code, Uint16Code, Int16Code,
		FloatCode, TimeOfDayCode, DateCode,
		Uint32Code, Int32Code, Float4ByteCode,
		AccessCode, StringCode, SceneNumberCode,
		SceneControlCode, DateTime8ByteCode, EnumCode,
		RedGreenBlueCode, RedGreenBlueWhiteCode:
		return nil
	}
	return fmt.Errorf("%w: %d", ErrUnknownType, t)
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
	case "rgbw", "redgreenbluewhite":
		return RedGreenBlueWhiteCode, nil
	}
	return TypeCode(0), fmt.Errorf("%w: %s", ErrUnknownType, s)
}

func NewDPT(code TypeCode) (DPT, error) {
	switch code {
	case BoolCode:
		return &DPTBool{}, nil
	case ControlCode:
		return &DPTControl{}, nil
	case StepControlCode:
		return &DPTStepControl{}, nil
	case CharCode:
		return &DPTChar{}, nil
	case Uint8Code:
		return &DPTUint8{}, nil
	case Int8Code:
		return &DPTInt8{}, nil
	case Uint16Code:
		return &DPTUint16{}, nil
	case Int16Code:
		return &DPTInt16{}, nil
	case FloatCode:
		return &DPTFloat{}, nil
	case TimeOfDayCode:
		return &DPTTimeOfDay{}, nil
	case DateCode:
		return &DPTDate{}, nil
	case Uint32Code:
		return &DPTUint32{}, nil
	case Int32Code:
		return &DPTInt32{}, nil
	case Float4ByteCode:
		return &DPTFloat4Byte{}, nil
	case AccessCode:
		return &DPTAccess{}, nil
	case StringCode:
		return &DPTString{}, nil
	case SceneNumberCode:
		return &DPTSceneNumber{}, nil
	case SceneControlCode:
		return &DPTSceneControl{}, nil
	case DateTime8ByteCode:
		return &DPTDateTime8Byte{}, nil
	case EnumCode:
		return &DPTEnum{}, nil
	case RedGreenBlueCode:
		return &DPTRedGreenBlue{}, nil
	case RedGreenBlueWhiteCode:
		return &DPTRedGreenBlueWhite{}, nil
	}
	return nil, fmt.Errorf("%w: %d", ErrUnknownType, code)
}
