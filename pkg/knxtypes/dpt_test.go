package knxtypes

import (
	"testing"
)

func TestValidTypeCode(t *testing.T) {
	type testCase struct {
		name          string
		code          TypeCode
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Bool", code: BoolCode, errorExpected: false},
		{name: "Control", code: ControlCode, errorExpected: false},
		{name: "StepControl", code: StepControlCode, errorExpected: false},
		{name: "Char", code: CharCode, errorExpected: false},
		{name: "Uint8", code: Uint8Code, errorExpected: false},
		{name: "Int8", code: Int8Code, errorExpected: false},
		{name: "Uint16", code: Uint16Code, errorExpected: false},
		{name: "Int16", code: Int16Code, errorExpected: false},
		{name: "Float", code: FloatCode, errorExpected: false},
		{name: "TimeOfDay", code: TimeOfDayCode, errorExpected: false},
		{name: "Date", code: DateCode, errorExpected: false},
		{name: "Uint32", code: Uint32Code, errorExpected: false},
		{name: "Int32", code: Int32Code, errorExpected: false},
		{name: "Float4Byte", code: Float4ByteCode, errorExpected: false},
		{name: "Access", code: AccessCode, errorExpected: false},
		{name: "String", code: StringCode, errorExpected: false},
		{name: "SceneNumber", code: SceneNumberCode, errorExpected: false},
		{name: "SceneControl", code: SceneControlCode, errorExpected: false},
		{name: "DateTime8Byte", code: DateTime8ByteCode, errorExpected: false},
		{name: "Enum", code: EnumCode, errorExpected: false},
		{name: "RedGreenBlue", code: RedGreenBlueCode, errorExpected: false},
		{name: "RedGreenBlueWhite", code: RedGreenBlueWhiteCode, errorExpected: false},
		{name: "Zero", code: TypeCode(0), errorExpected: true},
		{name: "Negative", code: TypeCode(-1), errorExpected: true},
		{name: "Invalid", code: TypeCode(999), errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.code.Valid()
			switch {
			case err == nil && tc.errorExpected:
				t.Errorf("expected error for code %d", tc.code)
			case err != nil && !tc.errorExpected:
				t.Errorf("unexpected error for code %d: %s", tc.code, err.Error())
			}
		})
	}
}
