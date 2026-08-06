package knxtypes

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDPTUint16Marshal(t *testing.T) {
	type testCase struct {
		name  string
		value uint16
	}

	testCases := []testCase{
		{name: "Zero", value: 0},
		{name: "Max", value: math.MaxUint16},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTUint16{Value: tc.value}
			data, err := dpt.MarshalBinary()
			require.NoError(t, err)

			other := DPTUint16{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.Equal(t, dpt.Value, other.Value)
		})
	}
}

func TestDPTUint16BadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTUint16](t)
}

func TestDPTUint16Set(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTUint16
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Zero", input: "0", expected: DPTUint16{Value: 0}},
		{name: "Max", input: "65535", expected: DPTUint16{Value: 65535}},
		{name: "Bad", input: "abc", errorExpected: true},
		{name: "Negative", input: "-1", errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTUint16{}
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

func FuzzDPTUint16(f *testing.F) {
	dpt := DPTUint16{Value: 12345}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTUint16{}
		_ = other.UnmarshalBinary(data)
	})
}
