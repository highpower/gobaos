package knxtypes

import (
	"math"
	"testing"
)

func TestDPTFloat4ByteMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value float64
	}

	testCases := []testCase{
		{name: "Zero", value: 0.0},
		{name: "Positive", value: 22.0},
		{name: "Negative", value: -3.5},
		{name: "MaxFloat32", value: float64(math.MaxFloat32)},
		{name: "SmallestNonzeroFloat32", value: float64(math.SmallestNonzeroFloat32)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTFloat4Byte{Value: tc.value}
			data, err := dpt.MarshalBinary()
			if err != nil {
				t.Fatalf("%s", err.Error())
			}
			other := DPTFloat4Byte{}
			if err = other.UnmarshalBinary(data); err != nil {
				t.Fatalf("%s", err.Error())
			}
			if float32(other.Value) != float32(dpt.Value) {
				t.Errorf("got value=%f, want value=%f", other.Value, dpt.Value)
			}
		})
	}
}

func TestDPTFloat4ByteBadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTFloat4Byte](t)
}

func TestDPTFloat4ByteSet(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTFloat4Byte
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Zero", input: "0", expected: DPTFloat4Byte{Value: 0}, errorExpected: false},
		{name: "Positive", input: "22.0", expected: DPTFloat4Byte{Value: 22.0}, errorExpected: false},
		{name: "Negative", input: "-3.5", expected: DPTFloat4Byte{Value: -3.5}, errorExpected: false},
		{name: "Bad", input: "abc", errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTFloat4Byte{}
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

func FuzzDPTFloat4Byte(f *testing.F) {
	dpt := DPTFloat4Byte{Value: 22.0}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTFloat4Byte{}
		_ = other.UnmarshalBinary(data)
	})
}