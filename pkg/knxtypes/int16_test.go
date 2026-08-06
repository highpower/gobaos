package knxtypes

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDPTInt16Marshal(t *testing.T) {
	type testCase struct {
		name  string
		value int16
	}

	testCases := []testCase{
		{name: "Zero", value: 0},
		{name: "Positive", value: math.MaxInt16},
		{name: "Negative", value: math.MinInt16},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTInt16{Value: tc.value}
			data, err := dpt.MarshalBinary()
			require.NoError(t, err)

			other := DPTInt16{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.Equal(t, dpt.Value, other.Value)
		})
	}
}

func TestDPTInt16BadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTInt16](t)
}

func TestDPTInt16Set(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTInt16
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Zero", input: "0", expected: DPTInt16{Value: 0}},
		{name: "Max", input: strconv.Itoa(math.MaxInt16), expected: DPTInt16{Value: math.MaxInt16}},
		{name: "Min", input: strconv.Itoa(math.MinInt16), expected: DPTInt16{Value: math.MinInt16}},
		{name: "Bad", input: "abc", errorExpected: true},
		{name: "Overflow", input: strconv.FormatInt(int64(math.MaxInt16)+1, 10), errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTInt16{}
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

func FuzzDPTInt16(f *testing.F) {
	dpt := DPTInt16{Value: -12345}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTInt16{}
		_ = other.UnmarshalBinary(data)
	})
}
