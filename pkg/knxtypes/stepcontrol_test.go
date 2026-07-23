package knxtypes

import "testing"

func TestDPTStepControlMarshal(t *testing.T) {
	type testCase struct {
		name          string
		control       bool
		stepCode      uint8
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Zero", control: false, stepCode: 0, errorExpected: false},
		{name: "MaxStepCode", control: true, stepCode: 7, errorExpected: false},
		{name: "StepCodeTooLarge", control: false, stepCode: 8, errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dpt := DPTStepControl{Control: tc.control, StepCode: tc.stepCode}
			data, err := dpt.MarshalBinary()
			switch {
			case err == nil && tc.errorExpected:
				t.Fatalf("error expected but did not happen")
			case err != nil && !tc.errorExpected:
				t.Fatalf("%s", err.Error())
			case err != nil:
				return
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
	t.Parallel()

	err := (&DPTStepControl{}).UnmarshalBinary([]byte{0x00, 0x00})
	if err == nil {
		t.Fatal("expected error for invalid length")
	}
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
		{name: "ValidFalse", input: "false:0", expected: DPTStepControl{false, 0}, errorExpected: false},
		{name: "WrongParts", input: "true", errorExpected: true},
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
