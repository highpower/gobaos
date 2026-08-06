package knxtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	case RedGreenBlueWhiteCode:
		return 6
	case DateTime8ByteCode:
		return 8
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
	sizes := map[string]int{"Zero": 0, "Larger": length + 1}
	if length != 1 {
		sizes["Smaller"] = length - 1
	}
	for name, size := range sizes {
		t.Run(name, func(t *testing.T) {
			data := make([]byte, size)
			for i := range data {
				data[i] = byte(i)
			}
			assert.Error(t, W(&value).UnmarshalBinary(data))
		})
	}
}
