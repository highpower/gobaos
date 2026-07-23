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

func TestDPTControlUnmarshal(t *testing.T) {
	err := (&DPTControl{}).UnmarshalBinary([]byte{0x00, 0x00})
	if err == nil {
		t.Fatal("expected error for invalid length")
	}
}

func TestDPTControlSet(t *testing.T) {

	type testCase struct {
		name     string
		input    string
		expected DPTControl
		err      bool
	}

	testCases := []testCase{
		{name: "Valid", input: "true:true", expected: DPTControl{Value: true, Control: true}, err: false},
		{name: "ValidFalse", input: "false:false", expected: DPTControl{Value: false, Control: false}, err: false},
		{name: "WrongParts", input: "true", err: true},
		{name: "BadBool", input: "maybe:true", err: true},
		{name: "BadControl", input: "true:maybe", err: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dpt := DPTControl{}
			err := dpt.Set(tc.input)
			switch {
			case err == nil && tc.err:
				t.Fatalf("error expected but did not happen")
			case err != nil && !tc.err:
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
