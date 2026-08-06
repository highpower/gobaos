package knxtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDPTControlMarshal(t *testing.T) {
	type testCase struct {
		name    string
		value   bool
		control bool
	}

	testCases := []testCase{
		{name: "Zero", value: false, control: false},
		{name: "Both", value: true, control: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTControl{Value: tc.value, Control: tc.control}
			data, err := dpt.MarshalBinary()
			require.NoError(t, err)

			other := DPTControl{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.Equal(t, dpt.Value, other.Value)
			assert.Equal(t, dpt.Control, other.Control)
		})
	}
}

func TestDPTControlBadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTControl](t)
}

func TestDPTControlSet(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTControl
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Valid", input: "true:true", expected: DPTControl{Value: true, Control: true}},
		{name: "WrongParts", input: "true:true:true", errorExpected: true},
		{name: "BadValue", input: "maybe:true", errorExpected: true},
		{name: "BadControl", input: "true:maybe", errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTControl{}
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

func FuzzDPTControl(f *testing.F) {
	dpt := DPTControl{Control: true, Value: true}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTControl{}
		_ = other.UnmarshalBinary(data)
	})
}
