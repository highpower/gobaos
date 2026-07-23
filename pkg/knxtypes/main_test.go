package knxtypes

import "testing"

type dptWrapper[G any] interface {
	*G
	DPT
}

func bytesRequired(code TypeCode) int {
	switch code {
	case Bool, Control, StepControl, Char, UInt8, Int8:
		return 1
	case UInt16, Int16, Float:
		return 2
	case TimeOfDay, Date:
		return 3
	}
	return 0
}

func testBadUnmarshal[G any, W dptWrapper[G]](t *testing.T) {
	var value G
	data := make([]byte, bytesRequired(W(&value).TypeCode())+1)
	for i := range data {
		data[i] = byte(i)
	}
	if err := W(&value).UnmarshalBinary(data); err == nil {
		t.Fatalf("error expected during unmarshal of %s", W(&value).TypeCode().String())
	}
}
