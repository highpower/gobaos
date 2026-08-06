package knxtypes

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		{name: "MaxFloat32", value: math.MaxFloat32},
		{name: "SmallestNonzeroFloat32", value: math.SmallestNonzeroFloat32},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTFloat4Byte{Value: tc.value}
			data, err := dpt.MarshalBinary()
			require.NoError(t, err)

			other := DPTFloat4Byte{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.InDelta(t, dpt.Value, other.Value, math.SmallestNonzeroFloat64)
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
		{name: "Zero", input: "0", expected: DPTFloat4Byte{Value: 0}},
		{name: "Positive", input: "22.0", expected: DPTFloat4Byte{Value: 22.0}},
		{name: "Negative", input: "-3.5", expected: DPTFloat4Byte{Value: -3.5}},
		{name: "Bad", input: "abc", errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTFloat4Byte{}
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

func FuzzDPTFloat4Byte(f *testing.F) {
	dpt := DPTFloat4Byte{Value: 22.0}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTFloat4Byte{}
		_ = other.UnmarshalBinary(data)
	})
}
