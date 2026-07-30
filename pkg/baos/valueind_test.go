package baos

import (
	"testing"
)

func TestDatapointValueIndValid(t *testing.T) {
	type testCase struct {
		name          string
		value         DatapointValueInd
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Empty", value: DatapointValueInd{}, errorExpected: true},
		{name: "Valid", value: DatapointValueInd{StartDatapoint: 1,
			Values: []Datapoint{{ID: 1, State: 0x10, Data: []byte{0x01}}}}, errorExpected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.value.Valid()
			switch {
			case err == nil && tc.errorExpected:
				t.Errorf("error expected but did not happen")
			case err != nil && !tc.errorExpected:
				t.Errorf("%s", err.Error())
			}
		})
	}
}

func TestDatapointValueIndEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          DatapointValueInd
		other          DatapointValueInd
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal",
			value: DatapointValueInd{StartDatapoint: 1, Values: []Datapoint{{ID: 1, State: 0x10, Data: []byte{0x01}}}},
			other: DatapointValueInd{StartDatapoint: 1, Values: []Datapoint{{ID: 1, State: 0x10, Data: []byte{0x01}}}}, resultExpected: true},
		{name: "StartUnequal",
			value: DatapointValueInd{StartDatapoint: 1}, other: DatapointValueInd{StartDatapoint: 2}, resultExpected: false},
		{name: "ValuesLenUnequal",
			value: DatapointValueInd{Values: []Datapoint{{ID: 1}}},
			other: DatapointValueInd{Values: []Datapoint{{ID: 2}}}, resultExpected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if res := tc.value.Equal(&tc.other); res != tc.resultExpected {
				t.Errorf("incorrect result: expected %t, got %t", tc.resultExpected, res)
			}
		})
	}
}

func TestDatapointValueIndMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value DatapointValueInd
	}

	testCases := []testCase{
		{name: "Empty", value: DatapointValueInd{}},
		{name: "WithValues", value: DatapointValueInd{StartDatapoint: 2, Values: []Datapoint{
			{ID: 10, State: 0x10, Data: []byte{0x01}}, {ID: 11, State: 0x20, Data: []byte{0x03, 0x04}}}}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			if err != nil {
				t.Fatalf("marshal: %s", err.Error())
			}
			other := DatapointValueInd{}
			if err := other.UnmarshalBinary(data); err != nil {
				t.Fatalf("unmarshal: %s", err.Error())
			}
			if !tc.value.Equal(&other) {
				t.Errorf("round-trip failed: expected equal values")
			}
		})
	}
}

func TestDatapointValueIndUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name        string
		data        []byte
		errExpected bool
	}

	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0xC1, 0x00}, errExpected: true},
		{name: "WrongService", data: []byte{0x00, 0xC1, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "TruncatedValues", data: []byte{0xF0, 0xC1, 0x00, 0x00, 0x00, 0x01}, errExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := DatapointValueInd{}
			err := r.UnmarshalBinary(tc.data)
			switch {
			case err == nil && tc.errExpected:
				t.Errorf("error expected but did not happen")
			case err != nil && !tc.errExpected:
				t.Errorf("unexpected error: %s", err.Error())
			}
		})
	}
}

func FuzzDatapointValueInd(f *testing.F) {

	ind := DatapointValueInd{StartDatapoint: 1, Values: []Datapoint{{ID: 1, State: 0x10, Data: []byte{0x01}}}}
	data, _ := ind.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DatapointValueInd{}
		_ = other.UnmarshalBinary(data)
	})
}
