package knxtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDPTEnumMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value uint8
	}

	testCases := []testCase{
		{name: "Zero", value: 0},
		{name: "Max", value: 255},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTEnum{Value: tc.value}
			data, err := dpt.MarshalBinary()
			require.NoError(t, err)

			other := DPTEnum{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.Equal(t, dpt.Value, other.Value)
		})
	}
}

func TestDPTEnumBadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTEnum](t)
}

func TestDPTEnumSet(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTEnum
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Zero", input: "0", expected: DPTEnum{Value: 0}},
		{name: "Max", input: "255", expected: DPTEnum{Value: 255}},
		{name: "Mid", input: "128", expected: DPTEnum{Value: 128}},
		{name: "Bad", input: "abc", errorExpected: true},
		{name: "Overflow", input: "256", errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTEnum{}
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

func FuzzDPTEnum(f *testing.F) {
	dpt := DPTEnum{Value: 128}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTEnum{}
		_ = other.UnmarshalBinary(data)
	})
}
