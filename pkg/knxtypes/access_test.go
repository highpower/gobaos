package knxtypes

import (
	"testing"
)

func TestDPTAccessMarshal(t *testing.T) {
	type testCase struct {
		name           string
		key            uint8
		id             uint8
		index          uint8
		additionalData uint8
	}

	testCases := []testCase{
		{name: "Simple", key: 10, id: 20, index: 5, additionalData: 1},
		{name: "Zero", key: 0, id: 0, index: 0, additionalData: 0},
		{name: "Max", key: 127, id: 127, index: 15, additionalData: 3},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			acc := DPTAccess{Key: tc.key, ID: tc.id, Index: tc.index, AdditionalData: tc.additionalData}
			data, err := acc.MarshalBinary()
			if err != nil {
				t.Fatalf("marshal: %s", err.Error())
			}
			other := DPTAccess{}
			if err := other.UnmarshalBinary(data); err != nil {
				t.Fatalf("unmarshal: %s", err.Error())
			}
			if other != acc {
				t.Errorf("got value=%s, want value=%s", other.String(), acc.String())
			}
		})
	}
}

func TestDPTAccessBadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTAccess](t)
}

func TestDPTAccessSet(t *testing.T) {
	type testCase struct {
		name          string
		value         string
		result        DPTAccess
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Valid", value: "10:20:5:1", result: DPTAccess{Key: 10, ID: 20, Index: 5, AdditionalData: 1}, errorExpected: false},
		{name: "Min", value: "0:0:0:0", result: DPTAccess{Key: 0, ID: 0, Index: 0, AdditionalData: 0}, errorExpected: false},
		{name: "Max", value: "127:127:15:3", result: DPTAccess{Key: 127, ID: 127, Index: 15, AdditionalData: 3}, errorExpected: false},
		{name: "KeyTooLarge", value: "128:0:0:0", result: DPTAccess{}, errorExpected: true},
		{name: "IDTooLarge", value: "0:128:0:0", result: DPTAccess{}, errorExpected: true},
		{name: "IndexTooLarge", value: "0:0:16:0", result: DPTAccess{}, errorExpected: true},
		{name: "AddDataTooLarge", value: "0:0:0:4", result: DPTAccess{}, errorExpected: true},
		{name: "TooFewParts", value: "1:2:3", result: DPTAccess{}, errorExpected: true},
		{name: "TooManyParts", value: "1:2:3:4:5", result: DPTAccess{}, errorExpected: true},
		{name: "InvalidKey", value: "abc:0:0:0", result: DPTAccess{}, errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res := DPTAccess{}
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

func FuzzDPTAccess(f *testing.F) {
	a := DPTAccess{Key: 10, ID: 20, Index: 5, AdditionalData: 1}
	data, _ := a.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTAccess{}
		_ = other.UnmarshalBinary(data)
	})
}
