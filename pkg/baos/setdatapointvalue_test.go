package baos

import (
	"testing"
)

func TestSetDatapointValueReqValid(t *testing.T) {
	type testCase struct {
		name          string
		value         SetDatapointValueReq
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Empty",
			value: SetDatapointValueReq{}, errorExpected: true},
		{name: "WithValues",
			value: SetDatapointValueReq{Commands: []DatapointCommand{{ID: 1, Command: 1, Data: []byte{0x01}}}}, errorExpected: false},
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

func TestSetDatapointValueReqEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          SetDatapointValueReq
		other          SetDatapointValueReq
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal",
			value: SetDatapointValueReq{Commands: []DatapointCommand{{ID: 1, Command: 1, Data: []byte{0x01}}}},
			other: SetDatapointValueReq{Commands: []DatapointCommand{{ID: 1, Command: 1, Data: []byte{0x01}}}}, resultExpected: true},
		{name: "IDUnequal",
			value: SetDatapointValueReq{Commands: []DatapointCommand{{ID: 1}}},
			other: SetDatapointValueReq{Commands: []DatapointCommand{{ID: 2}}}, resultExpected: false},
		{name: "CommandUnequal",
			value: SetDatapointValueReq{Commands: []DatapointCommand{{ID: 1, Command: 1}}},
			other: SetDatapointValueReq{Commands: []DatapointCommand{{ID: 1, Command: 2}}}, resultExpected: false},
		{name: "DataUnequal",
			value: SetDatapointValueReq{Commands: []DatapointCommand{{ID: 1, Data: []byte{0x01}}}},
			other: SetDatapointValueReq{Commands: []DatapointCommand{{ID: 1, Data: []byte{0x02}}}}, resultExpected: false},
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

func TestSetDatapointValueReqMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value SetDatapointValueReq
	}
	testCases := []testCase{
		{name: "SingleValue",
			value: SetDatapointValueReq{Commands: []DatapointCommand{{ID: 1, Command: 1, Data: []byte{0x01}}}}},
		{name: "MultipleValues",
			value: SetDatapointValueReq{Commands: []DatapointCommand{
				{ID: 10, Command: 2, Data: []byte{0x01, 0x02}},
				{ID: 11, Command: 3, Data: []byte{0x03}},
			}}},
		{name: "EmptyData",
			value: SetDatapointValueReq{Commands: []DatapointCommand{{ID: 5, Command: 5, Data: []byte{}}}}},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.value.MarshalBinary()
			if err != nil {
				t.Fatalf("%s", err.Error())
			}
			other := SetDatapointValueReq{}
			if err := other.UnmarshalBinary(data); err != nil {
				t.Fatalf("%s", err.Error())
			}
			if !tc.value.Equal(&other) {
				t.Errorf("incorrect result: %s=%s expected to be equal", tc.value.String(), other.String())
			}
		})
	}
}

func TestSetDatapointValueReqUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name        string
		data        []byte
		errExpected bool
	}
	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0x06, 0x00}, errExpected: true},
		{name: "WrongService", data: []byte{0x00, 0x06, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "TruncatedValues", data: []byte{0xF0, 0x06, 0x00, 0x00, 0x00, 0x01}, errExpected: true},
		{name: "ValueDataTooShort", data: []byte{0xF0, 0x06, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x02}, errExpected: true},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := SetDatapointValueReq{}
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

func TestSetDatapointValueResValid(t *testing.T) {
	type testCase struct {
		name          string
		value         SetDatapointValueRes
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Empty", value: SetDatapointValueRes{}, errorExpected: false},
		{name: "WithErrorCode", value: SetDatapointValueRes{StartDatapoint: 1, ErrorCode: 5}, errorExpected: false},
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

func TestSetDatapointValueResEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          SetDatapointValueRes
		other          SetDatapointValueRes
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal", value: SetDatapointValueRes{StartDatapoint: 1, NumberOfDatapoints: 0, ErrorCode: 2},
			other: SetDatapointValueRes{StartDatapoint: 1, NumberOfDatapoints: 0, ErrorCode: 2}, resultExpected: true},
		{name: "StartUnequal", value: SetDatapointValueRes{StartDatapoint: 1},
			other: SetDatapointValueRes{StartDatapoint: 2}, resultExpected: false},
		{name: "ErrorCodeUnequal", value: SetDatapointValueRes{ErrorCode: 1},
			other: SetDatapointValueRes{ErrorCode: 2}, resultExpected: false},
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

func TestSetDatapointValueResMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value SetDatapointValueRes
	}
	testCases := []testCase{
		{name: "Empty", value: SetDatapointValueRes{}},
		{name: "WithError", value: SetDatapointValueRes{StartDatapoint: 3, ErrorCode: 5}},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.value.MarshalBinary()
			if err != nil {
				t.Fatalf("marshal: %s", err.Error())
			}
			other := SetDatapointValueRes{}
			if err := other.UnmarshalBinary(data); err != nil {
				t.Fatalf("unmarshal: %s", err.Error())
			}
			if !tc.value.Equal(&other) {
				t.Errorf("round-trip failed: expected equal values")
			}
		})
	}
}

func TestSetDatapointValueResUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name        string
		data        []byte
		errExpected bool
	}
	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0x86, 0x00}, errExpected: true},
		{name: "WrongService", data: []byte{0x00, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := SetDatapointValueRes{}
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
