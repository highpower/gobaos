package knxtypes

import "testing"

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
			if err != nil {
				t.Fatalf("%s", err.Error())
			}
			other := DPTSceneControl{}
			if err = other.UnmarshalBinary(data); err != nil {
				t.Fatalf("%s", err.Error())
			}
			if other.Scene != dpt.Scene || other.Save != dpt.Save {
				t.Errorf("got scene=%d save=%t, want scene=%d save=%t", other.Scene, other.Save, dpt.Scene, dpt.Save)
			}
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
		{name: "ExecuteMin", input: "1,false", expected: DPTSceneControl{Scene: 1, Save: false}, errorExpected: false},
		{name: "SaveMax", input: "64,true", expected: DPTSceneControl{Scene: 64, Save: true}, errorExpected: false},
		{name: "ExecuteMid", input: "32,false", expected: DPTSceneControl{Scene: 32, Save: false}, errorExpected: false},
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
			switch {
			case err == nil && tc.errorExpected:
				t.Fatalf("error expected but did not happen")
			case err != nil && !tc.errorExpected:
				t.Fatalf("%s", err.Error())
			case err != nil:
				return
			}
			if dpt != tc.expected {
				t.Errorf("got dpt=%s, expected %s", dpt.String(), tc.expected.String())
			}
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
