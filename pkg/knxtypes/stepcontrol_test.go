package knxtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDPTStepControlMarshal(t *testing.T) {
	type testCase struct {
		name     string
		control  bool
		stepCode uint8
	}

	testCases := []testCase{
		{name: "Zero", control: false, stepCode: 0},
		{name: "MaxStepCode", control: true, stepCode: 7},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTStepControl{Control: tc.control, StepCode: tc.stepCode}
			data, err := dpt.MarshalBinary()
			require.NoError(t, err)

			other := DPTStepControl{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.Equal(t, dpt.Control, other.Control)
			assert.Equal(t, dpt.StepCode, other.StepCode)
		})
	}
}

func TestDPTStepControlUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTStepControl](t)
}

func TestDPTStepControlSet(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTStepControl
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Valid", input: "true:3", expected: DPTStepControl{true, 3}},
		{name: "WrongParts", input: "true:1:*", errorExpected: true},
		{name: "BadBool", input: "maybe:1", errorExpected: true},
		{name: "BadNumber", input: "true:abc", errorExpected: true},
		{name: "OutOfRange", input: "true:8", errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTStepControl{}
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

func FuzzDPTStepControl(f *testing.F) {
	dpt := DPTStepControl{Control: true, StepCode: 3}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTStepControl{}
		_ = other.UnmarshalBinary(data)
	})
}
