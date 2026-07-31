package knxtypes

import (
	"testing"
)

func TestDPTRedGreenBlueWhiteMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value DPTRedGreenBlueWhite
	}

	testCases := []testCase{
		{name: "Simple", value: DPTRedGreenBlueWhite{Red: 10, Green: 20, Blue: 30, White: 40, Flags: 0x2552}},
		{name: "Zero", value: DPTRedGreenBlueWhite{Red: 0, Green: 0, Blue: 0, White: 0, Flags: 0}},
		{name: "Max", value: DPTRedGreenBlueWhite{Red: 255, Green: 255, Blue: 255, White: 255, Flags: 0xFFFF}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			if err != nil {
				t.Fatalf("marshal: %s", err.Error())
			}
			other := DPTRedGreenBlueWhite{}
			if err := other.UnmarshalBinary(data); err != nil {
				t.Fatalf("unmarshal: %s", err.Error())
			}
			if other != tc.value {
				t.Errorf("round-trip failed: got %v, want %v", other, tc.value)
			}
		})
	}
}

func TestDPTRedGreenBlueWhiteBadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTRedGreenBlueWhite](t)
}

func TestDPTRedGreenBlueWhiteSet(t *testing.T) {
	type testCase struct {
		name          string
		value         string
		result        DPTRedGreenBlueWhite
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Valid", value: "10,20,30,40,2552",
			result:        DPTRedGreenBlueWhite{Red: 10, Green: 20, Blue: 30, White: 40, Flags: 0x2552},
			errorExpected: false},
		{name: "Min", value: "0,0,0,0,0",
			result:        DPTRedGreenBlueWhite{Red: 0, Green: 0, Blue: 0, White: 0, Flags: 0},
			errorExpected: false},
		{name: "Max", value: "255,255,255,255,FFFF",
			result:        DPTRedGreenBlueWhite{Red: 255, Green: 255, Blue: 255, White: 255, Flags: 0xFFFF},
			errorExpected: false},
		{name: "TooFewParts", value: "1,2,3", result: DPTRedGreenBlueWhite{}, errorExpected: true},
		{name: "TooManyParts", value: "1,2,3,4,5,7", result: DPTRedGreenBlueWhite{}, errorExpected: true},
		{name: "InvalidRed", value: "abc,0,0,0,2552", result: DPTRedGreenBlueWhite{}, errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res := DPTRedGreenBlueWhite{}
			err := res.Set(tc.value)
			switch {
			case err != nil && !tc.errorExpected:
				t.Errorf("error %s", err.Error())
			case err == nil && tc.errorExpected:
				t.Errorf("error expected")
			case err == nil && !tc.errorExpected:
				if res != tc.result {
					t.Errorf("incorrect result %v, expected %v", res, tc.result)
				}
			}
		})
	}
}

func FuzzDPTRedGreenBlueWhite(f *testing.F) {
	d := DPTRedGreenBlueWhite{Red: 10, Green: 20, Blue: 30, White: 40}
	data, _ := d.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTRedGreenBlueWhite{}
		_ = other.UnmarshalBinary(data)
	})
}
