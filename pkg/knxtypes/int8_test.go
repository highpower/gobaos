package knxtypes

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDPTInt8Marshal(t *testing.T) {
	type testCase struct {
		name  string
		value int8
	}

	testCases := []testCase{
		{name: "Zero", value: 0},
		{name: "Positive", value: math.MaxInt8},
		{name: "Negative", value: math.MinInt8},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTInt8{Value: tc.value}
			data, err := dpt.MarshalBinary()
			require.NoError(t, err)

			other := DPTInt8{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.Equal(t, dpt.Value, other.Value)
		})
	}
}

func TestDPTInt8BadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTInt8](t)
}

func TestDPTInt8Set(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTInt8
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Zero", input: "0", expected: DPTInt8{Value: 0}},
		{name: "Max", input: "127", expected: DPTInt8{Value: 127}},
		{name: "Min", input: "-128", expected: DPTInt8{Value: -128}},
		{name: "Bad", input: "abc", errorExpected: true},
		{name: "Overflow", input: strconv.FormatInt(int64(math.MaxInt8)+1, 10), errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTInt8{}
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

func FuzzDPTInt8(f *testing.F) {
	dpt := DPTInt8{Value: -42}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTInt8{}
		_ = other.UnmarshalBinary(data)
	})
}
