package knxtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDPTSceneControlMarshal(t *testing.T) {
	type testCase struct {
		name  string
		scene uint8
		save  bool
	}

	testCases := []testCase{
		{name: "ExecuteMin", scene: 1, save: false},
		{name: "SaveMax", scene: 64, save: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTSceneControl{Scene: tc.scene, Save: tc.save}
			data, err := dpt.MarshalBinary()
			require.NoError(t, err)

			other := DPTSceneControl{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.Equal(t, dpt.Scene, other.Scene)
			assert.Equal(t, dpt.Save, other.Save)
		})
	}
}

func TestDPTSceneControlBadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTSceneControl](t)
}

func TestDPTSceneControlSet(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTSceneControl
		errorExpected bool
	}

	testCases := []testCase{
		{name: "ExecuteMin", input: "1,false", expected: DPTSceneControl{Scene: 1, Save: false}},
		{name: "SaveMax", input: "64,true", expected: DPTSceneControl{Scene: 64, Save: true}},
		{name: "ExecuteMid", input: "32,false", expected: DPTSceneControl{Scene: 32, Save: false}},
		{name: "Zero", input: "0,false", errorExpected: true},
		{name: "Overflow", input: "65,true", errorExpected: true},
		{name: "WrongParts", input: "1", errorExpected: true},
		{name: "BadScene", input: "abc,true", errorExpected: true},
		{name: "BadSave", input: "1,maybe", errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTSceneControl{}
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

func FuzzDPTSceneControl(f *testing.F) {
	dpt := DPTSceneControl{Scene: 32, Save: true}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTSceneControl{}
		_ = other.UnmarshalBinary(data)
	})
}
