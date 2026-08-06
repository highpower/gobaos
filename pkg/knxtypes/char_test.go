package knxtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDPTCharMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value rune
	}

	testCases := []testCase{
		{name: "A", value: 'A'},
		{name: "Zero", value: 0},
		{name: "Max", value: 255},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTChar{Value: tc.value}
			data, err := dpt.MarshalBinary()
			require.NoError(t, err)

			other := DPTChar{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.Equal(t, dpt.Value, other.Value)
		})
	}
}

func TestDPTCharBadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTChar](t)
}

func TestDPTCharSet(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTChar
		errorExpected bool
	}

	testCases := []testCase{
		{name: "A", input: "A", expected: DPTChar{Value: 'A'}},
		{name: "Empty", input: "", errorExpected: true},
		{name: "TooLong", input: "AB", errorExpected: true},
		{name: "MultiByte", input: "世", errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTChar{}
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

func FuzzDPTChar(f *testing.F) {
	dpt := DPTChar{Value: 'A'}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTChar{}
		_ = other.UnmarshalBinary(data)
	})
}
