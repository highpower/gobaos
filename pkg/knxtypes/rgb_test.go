package knxtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDPTRedGreenBlueMarshal(t *testing.T) {
	type testCase struct {
		name  string
		red   uint8
		green uint8
		blue  uint8
	}

	testCases := []testCase{
		{name: "Zero", red: 0, green: 0, blue: 0},
		{name: "Max", red: 255, green: 255, blue: 255},
		{name: "Mixed", red: 128, green: 64, blue: 32},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTRedGreenBlue{Red: tc.red, Green: tc.green, Blue: tc.blue}
			data, err := dpt.MarshalBinary()
			require.NoError(t, err)

			other := DPTRedGreenBlue{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.Equal(t, dpt.Red, other.Red)
			assert.Equal(t, dpt.Green, other.Green)
			assert.Equal(t, dpt.Blue, other.Blue)
		})
	}
}

func TestDPTRedGreenBlueBadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTRedGreenBlue](t)
}

func TestDPTRedGreenBlueSet(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTRedGreenBlue
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Valid", input: "255,128,0", expected: DPTRedGreenBlue{Red: 255, Green: 128, Blue: 0}},
		{name: "Zero", input: "0,0,0", expected: DPTRedGreenBlue{Red: 0, Green: 0, Blue: 0}},
		{name: "WrongParts", input: "255,128", errorExpected: true},
		{name: "BadNumber", input: "255,abc,0", errorExpected: true},
		{name: "Overflow", input: "255,128,256", errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTRedGreenBlue{}
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

func FuzzDPTRedGreenBlue(f *testing.F) {
	dpt := DPTRedGreenBlue{Red: 128, Green: 64, Blue: 32}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTRedGreenBlue{}
		_ = other.UnmarshalBinary(data)
	})
}
