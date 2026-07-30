package knxtypes

import (
	"testing"
)

type dptWrapper[G any] interface {
	*G
	DPT
}

func bytesRequired(code TypeCode) int {
	switch code {
	case BoolCode, ControlCode, StepControlCode, CharCode, Uint8Code, Int8Code, SceneNumberCode, SceneControlCode, EnumCode:
		return 1
	case Uint16Code, Int16Code, FloatCode:
		return 2
	case TimeOfDayCode, DateCode, RedGreenBlueCode:
		return 3
	case Uint32Code, Int32Code, Float4ByteCode, AccessCode:
		return 4
	case StringCode:
		return 14
	}
	return -1
}

func testBadUnmarshal[G any, W dptWrapper[G]](t *testing.T) {
	var value G
	code := W(&value).TypeCode()
	length := bytesRequired(code)
	if length < 0 {
		t.Fatalf("unknown length for type %s", code.String())
	}
	data := make([]byte, length+1)
	for i := range data {
		data[i] = byte(i)
	}
	if err := W(&value).UnmarshalBinary(data); err == nil {
		t.Errorf("error expected during unmarshal of %s", W(&value).TypeCode().String())
	}
}
