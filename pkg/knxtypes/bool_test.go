package knxtypes

import "testing"

func TestDPTBoolMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value bool
	}

	testCases := []testCase{
		{name: "False", value: false},
		{name: "True", value: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTBool{Value: tc.value}
			data, err := dpt.MarshalBinary()
			if err != nil {
				t.Fatalf("%s", err.Error())
			}
			other := DPTBool{}
			if err = other.UnmarshalBinary(data); err != nil {
				t.Fatalf("%s", err.Error())
			}
			if other.Value != dpt.Value {
				t.Errorf("got value=%t, want value=%t", other.Value, dpt.Value)
			}
		})
	}
}

func TestDPTBoolBadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTBool](t)
}

func TestDPTBoolSet(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTBool
		errorExpected bool
	}

	testCases := []testCase{
		{name: "True", input: "true", expected: DPTBool{Value: true}, errorExpected: false},
		{name: "False", input: "false", expected: DPTBool{Value: false}, errorExpected: false},
		{name: "One", input: "1", expected: DPTBool{Value: true}, errorExpected: false},
		{name: "Zero", input: "0", expected: DPTBool{Value: false}, errorExpected: false},
		{name: "Bad", input: "maybe", errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTBool{}
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

func FuzzDPTBool(f *testing.F) {
	dpt := DPTBool{Value: true}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTBool{}
		_ = other.UnmarshalBinary(data)
	})

}
