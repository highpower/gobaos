package knxtypes

import "testing"

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
			if err != nil {
				t.Fatalf("%s", err.Error())
			}
			other := DPTControl{}
			if err = other.UnmarshalBinary(data); err != nil {
				t.Fatalf("%s", err.Error())
			}
			if other.Value != dpt.Value || other.Control != dpt.Control {
				t.Errorf("got value=%t control=%t, want value=%t control=%t", other.Value, other.Control, dpt.Value, dpt.Control)
			}
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
		{name: "Valid", input: "true:true", expected: DPTControl{Value: true, Control: true}, errorExpected: false},
		{name: "WrongParts", input: "true:true:true", errorExpected: true},
		{name: "BadValue", input: "maybe:true", errorExpected: true},
		{name: "BadControl", input: "true:maybe", errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTControl{}
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

func FuzzDPTControl(f *testing.F) {
	dpt := DPTControl{Control: true, Value: true}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTControl{}
		_ = other.UnmarshalBinary(data)
	})
}
