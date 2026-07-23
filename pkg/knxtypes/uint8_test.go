package knxtypes

import (
	"math"
	"strconv"
	"testing"
)

func TestDPTUInt8Marshal(t *testing.T) {
	type testCase struct {
		name  string
		value uint8
	}

	testCases := []testCase{
		{name: "Zero", value: 0},
		{name: "Max", value: math.MaxUint8},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTUInt8{Value: tc.value}
			data, err := dpt.MarshalBinary()
			if err != nil {
				t.Fatalf("%s", err.Error())
			}
			other := DPTUInt8{}
			if err = other.UnmarshalBinary(data); err != nil {
				t.Fatalf("%s", err.Error())
			}
			if other.Value != dpt.Value {
				t.Errorf("got value=%d, want value=%d", other.Value, dpt.Value)
			}
		})
	}
}

func TestDPTUInt8BadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTUInt8](t)
}

func TestDPTUInt8Set(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTUInt8
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Zero", input: "0", expected: DPTUInt8{Value: 0}, errorExpected: false},
		{name: "Max", input: strconv.Itoa(math.MaxUint8), expected: DPTUInt8{Value: 255}, errorExpected: false},
		{name: "Bad", input: "abc", errorExpected: true},
		{name: "Negative", input: "-1", errorExpected: true},
		{name: "Overflow", input: strconv.FormatInt(int64(math.MaxUint8)+1, 10), errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTUInt8{}
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

func FuzzDPTUInt8(f *testing.F) {
	dpt := DPTUInt8{Value: 42}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTUInt8{}
		_ = other.UnmarshalBinary(data)
	})
}
