package knxtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDPTSceneNumberMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value uint8
	}

	testCases := []testCase{
		{name: "Min", value: 1},
		{name: "Max", value: 64},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTSceneNumber{Value: tc.value}
			data, err := dpt.MarshalBinary()
			require.NoError(t, err)

			other := DPTSceneNumber{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.Equal(t, dpt.Value, other.Value)
		})
	}
}

func TestDPTSceneNumberBadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTSceneNumber](t)
}

func TestDPTSceneNumberSet(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTSceneNumber
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Min", input: "1", expected: DPTSceneNumber{Value: 1}},
		{name: "Max", input: "64", expected: DPTSceneNumber{Value: 64}},
		{name: "Mid", input: "32", expected: DPTSceneNumber{Value: 32}},
		{name: "Zero", input: "0", errorExpected: true},
		{name: "Overflow", input: "65", errorExpected: true},
		{name: "Bad", input: "abc", errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTSceneNumber{}
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

func FuzzDPTSceneNumber(f *testing.F) {
	dpt := DPTSceneNumber{Value: 32}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTSceneNumber{}
		_ = other.UnmarshalBinary(data)
	})
}
