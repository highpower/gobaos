package knxtypes

import "testing"

func TestDPTStringMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value string
	}

	testCases := []testCase{
		{name: "Empty", value: ""},
		{name: "Short", value: "hello"},
		{name: "Max", value: "12345678901234"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTString{Value: tc.value}
			data, err := dpt.MarshalBinary()
			if err != nil {
				t.Fatalf("%s", err.Error())
			}
			other := DPTString{}
			if err = other.UnmarshalBinary(data); err != nil {
				t.Fatalf("%s", err.Error())
			}
			if other.Value != dpt.Value {
				t.Errorf("got value=%q, want value=%q", other.Value, dpt.Value)
			}
		})
	}
}

func TestDPTStringBadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTString](t)
}

func TestDPTStringSet(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTString
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Empty", input: "", expected: DPTString{Value: ""}, errorExpected: false},
		{name: "Short", input: "hello", expected: DPTString{Value: "hello"}, errorExpected: false},
		{name: "Max", input: "12345678901234", expected: DPTString{Value: "12345678901234"}, errorExpected: false},
		{name: "TooLong", input: "123456789012345", errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTString{}
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

func FuzzDPTString(f *testing.F) {
	dpt := DPTString{Value: "hello"}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTString{}
		_ = other.UnmarshalBinary(data)
	})
}
