package knxtypes

import "testing"

func TestDPTInt8Marshal(t *testing.T) {
	type testCase struct {
		name  string
		value int8
	}

	testCases := []testCase{
		{name: "Zero", value: 0},
		{name: "Positive", value: 127},
		{name: "Negative", value: -128},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTInt8{Value: tc.value}
			data, err := dpt.MarshalBinary()
			if err != nil {
				t.Fatalf("%s", err.Error())
			}
			other := DPTInt8{}
			if err = other.UnmarshalBinary(data); err != nil {
				t.Fatalf("%s", err.Error())
			}
			if other.Value != dpt.Value {
				t.Errorf("got value=%d, want value=%d", other.Value, dpt.Value)
			}
		})
	}
}

func TestDPTInt8BadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTInt8](t)
}

func TestDPTInt8Set(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTInt8
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Zero", input: "0", expected: DPTInt8{Value: 0}, errorExpected: false},
		{name: "Max", input: "127", expected: DPTInt8{Value: 127}, errorExpected: false},
		{name: "Min", input: "-128", expected: DPTInt8{Value: -128}, errorExpected: false},
		{name: "Bad", input: "abc", errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTInt8{}
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
