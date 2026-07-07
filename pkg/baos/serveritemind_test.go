package baos

import (
	"testing"
)

func TestServerItemIndValid(t *testing.T) {
	type testCase struct {
		name          string
		value         ServerItemInd
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Empty", value: ServerItemInd{}, errorExpected: true},
		{name: "WithItems", value: ServerItemInd{StartItem: 1, Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, errorExpected: false},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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

func TestServerItemIndEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          ServerItemInd
		other          ServerItemInd
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal", value: ServerItemInd{StartItem: 1, Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, other: ServerItemInd{StartItem: 1, Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, resultExpected: true},
		{name: "StartUnequal", value: ServerItemInd{StartItem: 1, Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, other: ServerItemInd{StartItem: 2, Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, resultExpected: false},
		{name: "ItemsLenUnequal", value: ServerItemInd{StartItem: 1, Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, other: ServerItemInd{StartItem: 1, Items: []ServerItem{}}, resultExpected: false},
		{name: "ItemDataUnequal", value: ServerItemInd{StartItem: 1, Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, other: ServerItemInd{StartItem: 1, Items: []ServerItem{{ID: 1, Data: []byte{0x02}}}}, resultExpected: false},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if res := tc.value.Equal(&tc.other); res != tc.resultExpected {
				t.Errorf("incorrect result: expected %t, got %t", tc.resultExpected, res)
			}
		})
	}
}

func TestServerItemIndMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value ServerItemInd
	}
	testCases := []testCase{
		{name: "WithItems", value: ServerItemInd{StartItem: 2, Items: []ServerItem{{ID: 10, Data: []byte{0x01, 0x02}}, {ID: 11, Data: []byte{0x03}}}}},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.value.MarshalBinary()
			if err != nil {
				t.Fatalf("marshal: %s", err.Error())
			}
			other := ServerItemInd{}
			if err := other.UnmarshalBinary(data); err != nil {
				t.Fatalf("unmarshal: %s", err.Error())
			}
			if !tc.value.Equal(&other) {
				t.Errorf("round-trip failed: expected equal values")
			}
		})
	}
}

func TestServerItemIndUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name        string
		data        []byte
		errExpected bool
	}
	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0xC2, 0x00}, errExpected: true},
		{name: "WrongService", data: []byte{0x00, 0xC2, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "TruncatedItems", data: []byte{0xF0, 0xC2, 0x00, 0x00, 0x00, 0x01}, errExpected: true},
		{name: "ItemDataTooShort", data: []byte{0xF0, 0xC2, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x02}, errExpected: true},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := ServerItemInd{}
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
