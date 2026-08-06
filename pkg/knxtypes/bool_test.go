package knxtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDPTBoolMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value bool
	}

	testCases := []testCase{
		{name: "False", value: false},
		{name: "True", value: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTBool{Value: tc.value}
			data, err := dpt.MarshalBinary()
			require.NoError(t, err)

			other := DPTBool{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.Equal(t, dpt.Value, other.Value)
		})
	}
}

func TestDPTBoolBadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTBool](t)
}

func TestDPTBoolSet(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTBool
		errorExpected bool
	}

	testCases := []testCase{
		{name: "True", input: "true", expected: DPTBool{Value: true}},
		{name: "False", input: "false", expected: DPTBool{Value: false}},
		{name: "One", input: "1", expected: DPTBool{Value: true}},
		{name: "Zero", input: "0", expected: DPTBool{Value: false}},
		{name: "Bad", input: "maybe", errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTBool{}
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

func FuzzDPTBool(f *testing.F) {
	dpt := DPTBool{Value: true}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTBool{}
		_ = other.UnmarshalBinary(data)
	})
}
