package knxtypes

import "testing"

func TestDPTCharMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value rune
	}

	testCases := []testCase{
		{name: "A", value: 'A'},
		{name: "Zero", value: 0},
		{name: "Max", value: 255},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTChar{Value: tc.value}
			data, err := dpt.MarshalBinary()
			if err != nil {
				t.Fatalf("%s", err.Error())
			}
			other := DPTChar{}
			if err = other.UnmarshalBinary(data); err != nil {
				t.Fatalf("%s", err.Error())
			}
			if other.Value != dpt.Value {
				t.Errorf("got value=%c, want value=%c", other.Value, dpt.Value)
			}
		})
	}
}

func TestDPTCharBadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTChar](t)
}

func TestDPTCharSet(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTChar
		errorExpected bool
	}

	testCases := []testCase{
		{name: "A", input: "A", expected: DPTChar{Value: 'A'}, errorExpected: false},
		{name: "Empty", input: "", errorExpected: true},
		{name: "TooLong", input: "AB", errorExpected: true},
		{name: "MultiByte", input: "世", errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTChar{}
			err := dpt.Set(tc.input)
			switch {
			case err == nil && tc.errorExpected:
				t.Fatalf("error expected but did not happen")
			case err != nil && !tc.errorExpected:
				t.Fatalf("%s", err.Error())
			case err != nil:
				return
			}
			if dpt != tc.expected {
				t.Errorf("got dpt=%s, expected %s", dpt.String(), tc.expected.String())
			}
		})
	}
}
