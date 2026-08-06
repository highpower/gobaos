package knxtypes

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDPTUint32Marshal(t *testing.T) {
	type testCase struct {
		name  string
		value uint32
	}

	testCases := []testCase{
		{name: "Zero", value: 0},
		{name: "Max", value: math.MaxUint32},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTUint32{Value: tc.value}
			data, err := dpt.MarshalBinary()
			require.NoError(t, err)

			other := DPTUint32{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.Equal(t, dpt.Value, other.Value)
		})
	}
}

func TestDPTUint32BadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTUint32](t)
}

func TestDPTUint32Set(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTUint32
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Zero", input: "0", expected: DPTUint32{Value: 0}},
		{name: "Max", input: strconv.FormatUint(math.MaxUint32, 10), expected: DPTUint32{Value: math.MaxUint32}},
		{name: "Bad", input: "abc", errorExpected: true},
		{name: "Negative", input: "-1", errorExpected: true},
		{name: "Overflow", input: strconv.FormatInt(int64(math.MaxUint32)+1, 10), errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTUint32{}
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

func FuzzDPTUint32(f *testing.F) {
	dpt := DPTUint32{Value: 123456789}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTUint32{}
		_ = other.UnmarshalBinary(data)
	})
}
