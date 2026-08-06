package knxtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDPTStringMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value string
	}

	testCases := []testCase{
		{name: "Empty", value: ""},
		{name: "Short", value: "hello"},
		{name: "Max", value: "12345678901234"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTString{Value: tc.value}
			data, err := dpt.MarshalBinary()
			require.NoError(t, err)

			other := DPTString{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.Equal(t, dpt.Value, other.Value)
		})
	}
}

func TestDPTStringBadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTString](t)
}

func TestDPTStringSet(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTString
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Empty", input: "", expected: DPTString{Value: ""}},
		{name: "Short", input: "hello", expected: DPTString{Value: "hello"}},
		{name: "Max", input: "12345678901234", expected: DPTString{Value: "12345678901234"}},
		{name: "TooLong", input: "123456789012345", errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTString{}
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

func FuzzDPTString(f *testing.F) {
	dpt := DPTString{Value: "hello"}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTString{}
		_ = other.UnmarshalBinary(data)
	})
}
