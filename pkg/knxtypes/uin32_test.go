package knxtypes

import (
	"math"
	"strconv"
	"testing"
)

func TestDPTUint32Marshal(t *testing.T) {
	type testCase struct {
		name  string
		value uint32
	}

	testCases := []testCase{
		{name: "Zero", value: 0},
		{name: "Max", value: math.MaxUint32},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTUint32{Value: tc.value}
			data, err := dpt.MarshalBinary()
			if err != nil {
				t.Fatalf("%s", err.Error())
			}
			other := DPTUint32{}
			if err = other.UnmarshalBinary(data); err != nil {
				t.Fatalf("%s", err.Error())
			}
			if other.Value != dpt.Value {
				t.Errorf("got value=%d, want value=%d", other.Value, dpt.Value)
			}
		})
	}
}

func TestDPTUint32BadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTUint32](t)
}

func TestDPTUint32Set(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTUint32
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Zero", input: "0", expected: DPTUint32{Value: 0}, errorExpected: false},
		{name: "Max", input: strconv.FormatUint(math.MaxUint32, 10), expected: DPTUint32{Value: math.MaxUint32},
			errorExpected: false},
		{name: "Bad", input: "abc", errorExpected: true},
		{name: "Negative", input: "-1", errorExpected: true},
		{name: "Overflow", input: strconv.FormatInt(int64(math.MaxUint32)+1, 10), errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTUint32{}
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

func FuzzDPTUint32(f *testing.F) {
	dpt := DPTUint32{Value: 123456789}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTUint32{}
		_ = other.UnmarshalBinary(data)
	})
}
