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
	case Bool, Control, StepControl, Char, Uint8, Int8:
		return 1
	case Uint16, Int16, Float:
		return 2
	case TimeOfDay, Date, RedGreenBlue:
		return 3
	case Uint32, Int32, Float4Byte:
		return 4
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
