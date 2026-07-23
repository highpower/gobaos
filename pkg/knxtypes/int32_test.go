package knxtypes

import (
	"math"
	"strconv"
	"testing"
)

func TestDPTInt32Marshal(t *testing.T) {
	type testCase struct {
		name  string
		value int32
	}

	testCases := []testCase{
		{name: "Zero", value: 0},
		{name: "Max", value: math.MaxInt32},
		{name: "Min", value: math.MinInt32},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTInt32{Value: tc.value}
			data, err := dpt.MarshalBinary()
			if err != nil {
				t.Fatalf("%s", err.Error())
			}
			other := DPTInt32{}
			if err = other.UnmarshalBinary(data); err != nil {
				t.Fatalf("%s", err.Error())
			}
			if other.Value != dpt.Value {
				t.Errorf("got value=%d, want value=%d", other.Value, dpt.Value)
			}
		})
	}
}

func TestDPTInt32BadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTInt32](t)
}

func TestDPTInt32Set(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTInt32
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Zero", input: "0", expected: DPTInt32{Value: 0}, errorExpected: false},
		{name: "Max", input: strconv.Itoa(math.MaxInt32), expected: DPTInt32{Value: math.MaxInt32}, errorExpected: false},
		{name: "Min", input: strconv.Itoa(math.MinInt32), expected: DPTInt32{Value: math.MinInt32}, errorExpected: false},
		{name: "Bad", input: "abc", errorExpected: true},
		{name: "Overflow", input: strconv.FormatInt(int64(math.MaxInt32)+1, 10), errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTInt32{}
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

func FuzzDPTInt32(f *testing.F) {
	dpt := DPTInt32{Value: -123456789}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTInt32{}
		_ = other.UnmarshalBinary(data)
	})
}
