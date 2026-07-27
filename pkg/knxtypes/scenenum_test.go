package knxtypes

import "testing"

func TestDPTSceneNumberMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value uint8
	}

	testCases := []testCase{
		{name: "Min", value: 1},
		{name: "Max", value: 64},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTSceneNumber{Value: tc.value}
			data, err := dpt.MarshalBinary()
			if err != nil {
				t.Fatalf("%s", err.Error())
			}
			other := DPTSceneNumber{}
			if err = other.UnmarshalBinary(data); err != nil {
				t.Fatalf("%s", err.Error())
			}
			if other.Value != dpt.Value {
				t.Errorf("got value=%d, want value=%d", other.Value, dpt.Value)
			}
		})
	}
}

func TestDPTSceneNumberBadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTSceneNumber](t)
}

func TestDPTSceneNumberSet(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTSceneNumber
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Min", input: "1", expected: DPTSceneNumber{Value: 1}, errorExpected: false},
		{name: "Max", input: "64", expected: DPTSceneNumber{Value: 64}, errorExpected: false},
		{name: "Mid", input: "32", expected: DPTSceneNumber{Value: 32}, errorExpected: false},
		{name: "Zero", input: "0", errorExpected: true},
		{name: "Overflow", input: "65", errorExpected: true},
		{name: "Bad", input: "abc", errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTSceneNumber{}
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

func FuzzDPTSceneNumber(f *testing.F) {
	dpt := DPTSceneNumber{Value: 32}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTSceneNumber{}
		_ = other.UnmarshalBinary(data)
	})
}
