package knxtypes

import (
	"math"
	"testing"
)

func TestDPTFloatMarshal(t *testing.T) {
	type testCase struct {
		name          string
		value         float64
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Zero", value: 0.0, errorExpected: false},
		{name: "Nonzero", value: 22.0, errorExpected: false},
		{name: "TooLarge", value: 2.2e4, errorExpected: true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dpt := DPTFloat{tc.value}
			data, err := dpt.MarshalBinary()
			switch {
			case err == nil && tc.errorExpected:
				t.Fatalf("error expected but did not happen")
			case err != nil && !tc.errorExpected:
				t.Fatalf("%s", err.Error())
			case err != nil:
				return
			}
			other := DPTFloat{}
			if err = other.UnmarshalBinary(data); err != nil {
				t.Fatalf("%s", err.Error())
			}
			if math.Abs(other.Value-dpt.Value) > math.SmallestNonzeroFloat64 {
				t.Errorf("%04f != %04f", dpt.Value, other.Value)
			}
		})
	}
}
