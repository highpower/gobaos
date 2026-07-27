package knxtypes

import "testing"

func TestDPTEnumMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value uint8
	}

	testCases := []testCase{
		{name: "Zero", value: 0},
		{name: "Max", value: 255},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTEnum{Value: tc.value}
			data, err := dpt.MarshalBinary()
			if err != nil {
				t.Fatalf("%s", err.Error())
			}
			other := DPTEnum{}
			if err = other.UnmarshalBinary(data); err != nil {
				t.Fatalf("%s", err.Error())
			}
			if other.Value != dpt.Value {
				t.Errorf("got value=%d, want value=%d", other.Value, dpt.Value)
			}
		})
	}
}

func TestDPTEnumBadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTEnum](t)
}

func TestDPTEnumSet(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTEnum
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Zero", input: "0", expected: DPTEnum{Value: 0}, errorExpected: false},
		{name: "Max", input: "255", expected: DPTEnum{Value: 255}, errorExpected: false},
		{name: "Mid", input: "128", expected: DPTEnum{Value: 128}, errorExpected: false},
		{name: "Bad", input: "abc", errorExpected: true},
		{name: "Overflow", input: "256", errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTEnum{}
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

func FuzzDPTEnum(f *testing.F) {
	dpt := DPTEnum{Value: 128}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTEnum{}
		_ = other.UnmarshalBinary(data)
	})
}
