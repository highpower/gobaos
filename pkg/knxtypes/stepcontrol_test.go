package knxtypes

import "testing"

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
			if err != nil {
				t.Fatalf("%s", err.Error())
			}
			other := DPTStepControl{}
			if err = other.UnmarshalBinary(data); err != nil {
				t.Fatalf("%s", err.Error())
			}
			if other.Control != dpt.Control || other.StepCode != dpt.StepCode {
				t.Errorf("got control=%t stepCode=%d, want control=%t stepCode=%d", other.Control, other.StepCode, dpt.Control, dpt.StepCode)
			}
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
		{name: "Valid", input: "true:3", expected: DPTStepControl{true, 3}, errorExpected: false},
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

func FuzzDPTStepControl(f *testing.F) {
	dpt := DPTStepControl{Control: true, StepCode: 3}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTStepControl{}
		_ = other.UnmarshalBinary(data)
	})
}
