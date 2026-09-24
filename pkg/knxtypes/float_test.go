package knxtypes

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDPTFloatMarshal(t *testing.T) {
	type testCase struct {
		name          string
		value         float64
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Zero", value: 0.0},
		{name: "Nonzero", value: 22.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dpt := DPTFloat{tc.value}
			data, err := dpt.MarshalBinary()
			if tc.errorExpected {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			other := DPTFloat{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.InDelta(t, dpt.Value, other.Value, math.SmallestNonzeroFloat64)
		})
	}
}

func TestDPTFloatBadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTFloat](t)
}

func TestDPTFloatSet(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTFloat
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Zero", input: "0", expected: DPTFloat{Value: 0}},
		{name: "Positive", input: "22.0", expected: DPTFloat{Value: 22.0}},
		{name: "Negative", input: "-3.5", expected: DPTFloat{Value: -3.5}},
		{name: "TooLarge", input: "2.25e18", errorExpected: true},
		{name: "Bad", input: "abc", errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTFloat{}
			err := dpt.Set(tc.input)
			if tc.errorExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, dpt)
		})
	}
}

func FuzzDPTFloat(f *testing.F) {
	dpt := DPTFloat{Value: 22.0}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTFloat{}
		_ = other.UnmarshalBinary(data)
	})
}
