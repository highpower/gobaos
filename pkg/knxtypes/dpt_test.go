package knxtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidTypeCode(t *testing.T) {
	type testCase struct {
		name          string
		code          TypeCode
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Bool", code: BoolCode},
		{name: "Control", code: ControlCode},
		{name: "StepControl", code: StepControlCode},
		{name: "Char", code: CharCode},
		{name: "Uint8", code: Uint8Code},
		{name: "Int8", code: Int8Code},
		{name: "Uint16", code: Uint16Code},
		{name: "Int16", code: Int16Code},
		{name: "Float", code: FloatCode},
		{name: "TimeOfDay", code: TimeOfDayCode},
		{name: "Date", code: DateCode},
		{name: "Uint32", code: Uint32Code},
		{name: "Int32", code: Int32Code},
		{name: "Float4Byte", code: Float4ByteCode},
		{name: "Access", code: AccessCode},
		{name: "String", code: StringCode},
		{name: "SceneNumber", code: SceneNumberCode},
		{name: "SceneControl", code: SceneControlCode},
		{name: "DateTime8Byte", code: DateTime8ByteCode},
		{name: "Enum", code: EnumCode},
		{name: "RedGreenBlue", code: RedGreenBlueCode},
		{name: "RedGreenBlueWhite", code: RedGreenBlueWhiteCode},
		{name: "Zero", code: TypeCode(0), errorExpected: true},
		{name: "Negative", code: TypeCode(-1), errorExpected: true},
		{name: "Invalid", code: TypeCode(999), errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.code.Valid()
			if tc.errorExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestParseTypeCode(t *testing.T) {
	type testCase struct {
		name          string
		value         string
		result        TypeCode
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Bool", value: "bool", result: BoolCode},
		{name: "BoolUpper", value: "BOOL", result: BoolCode},
		{name: "Control", value: "control", result: ControlCode},
		{name: "StepControl", value: "stepcontrol", result: StepControlCode},
		{name: "Char", value: "char", result: CharCode},
		{name: "Uint8", value: "uint8", result: Uint8Code},
		{name: "Int8", value: "int8", result: Int8Code},
		{name: "Uint16", value: "uint16", result: Uint16Code},
		{name: "Int16", value: "int16", result: Int16Code},
		{name: "Float", value: "float", result: FloatCode},
		{name: "TimeOfDay", value: "timeofday", result: TimeOfDayCode},
		{name: "Date", value: "date", result: DateCode},
		{name: "Uint32", value: "uint32", result: Uint32Code},
		{name: "Int32", value: "int32", result: Int32Code},
		{name: "Float4Byte", value: "float4byte", result: Float4ByteCode},
		{name: "Access", value: "access", result: AccessCode},
		{name: "String", value: "string", result: StringCode},
		{name: "SceneNumber", value: "scenenumber", result: SceneNumberCode},
		{name: "SceneControl", value: "scenecontrol", result: SceneControlCode},
		{name: "DateTime8Byte", value: "datetime8byte", result: DateTime8ByteCode},
		{name: "Enum", value: "enum", result: EnumCode},
		{name: "RGB", value: "rgb", result: RedGreenBlueCode},
		{name: "RedGreenBlue", value: "redgreenblue", result: RedGreenBlueCode},
		{name: "RGBW", value: "rgbw", result: RedGreenBlueWhiteCode},
		{name: "RedGreenBlueWhite", value: "redgreenbluewhite", result: RedGreenBlueWhiteCode},
		{name: "Empty", value: "", errorExpected: true},
		{name: "Invalid", value: "unknown", errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res, err := ParseTypeCode(tc.value)
			if tc.errorExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.result, res)
		})
	}
}
