package knxtypes

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDPTUInt8Marshal(t *testing.T) {
	type testCase struct {
		name  string
		value uint8
	}

	testCases := []testCase{
		{name: "Zero", value: 0},
		{name: "Max", value: math.MaxUint8},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTUint8{Value: tc.value}
			data, err := dpt.MarshalBinary()
			require.NoError(t, err)

			other := DPTUint8{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.Equal(t, dpt.Value, other.Value)
		})
	}
}

func TestDPTUInt8BadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTUint8](t)
}

func TestDPTUInt8Set(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTUint8
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Zero", input: "0", expected: DPTUint8{Value: 0}},
		{name: "Max", input: strconv.Itoa(math.MaxUint8), expected: DPTUint8{Value: 255}},
		{name: "Bad", input: "abc", errorExpected: true},
		{name: "Negative", input: "-1", errorExpected: true},
		{name: "Overflow", input: strconv.FormatInt(int64(math.MaxUint8)+1, 10), errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTUint8{}
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

func FuzzDPTUInt8(f *testing.F) {
	dpt := DPTUint8{Value: 42}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTUint8{}
		_ = other.UnmarshalBinary(data)
	})
}
