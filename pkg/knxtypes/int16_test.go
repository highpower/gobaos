package knxtypes

import (
	"math"
	"strconv"
	"testing"
)

func TestDPTInt16Marshal(t *testing.T) {
	type testCase struct {
		name  string
		value int16
	}

	testCases := []testCase{
		{name: "Zero", value: 0},
		{name: "Positive", value: math.MaxInt16},
		{name: "Negative", value: math.MinInt16},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTInt16{Value: tc.value}
			data, err := dpt.MarshalBinary()
			if err != nil {
				t.Fatalf("%s", err.Error())
			}
			other := DPTInt16{}
			if err = other.UnmarshalBinary(data); err != nil {
				t.Fatalf("%s", err.Error())
			}
			if other.Value != dpt.Value {
				t.Errorf("got value=%d, want value=%d", other.Value, dpt.Value)
			}
		})
	}
}

func TestDPTInt16BadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTInt16](t)
}

func TestDPTInt16Set(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTInt16
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Zero", input: "0", expected: DPTInt16{Value: 0}, errorExpected: false},
		{name: "Max", input: strconv.Itoa(math.MaxInt16), expected: DPTInt16{Value: math.MaxInt16}, errorExpected: false},
		{name: "Min", input: strconv.Itoa(math.MinInt16), expected: DPTInt16{Value: math.MinInt16}, errorExpected: false},
		{name: "Bad", input: "abc", errorExpected: true},
		{name: "Overflow", input: strconv.FormatInt(int64(math.MaxInt16)+1, 10), errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTInt16{}
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

func FuzzDPTInt16(f *testing.F) {
	dpt := DPTInt16{Value: -12345}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTInt16{}
		_ = other.UnmarshalBinary(data)
	})
}
