package baos

import "testing"

func TestGetDatapointValueReqValid(t *testing.T) {
	type testCase struct {
		name          string
		value         GetDatapointValueReq
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Empty",
			value: GetDatapointValueReq{0, 0, 0}, errorExpected: false},
		{name: "Nonempty",
			value: GetDatapointValueReq{1, 5, AllDatapoints}, errorExpected: false},
		{name: "BadFilter",
			value: GetDatapointValueReq{1, 5, 6}, errorExpected: true},
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

func TestGetDatapointValueReqEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          GetDatapointValueReq
		other          GetDatapointValueReq
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal",
			value: GetDatapointValueReq{1, 5, 0},
			other: GetDatapointValueReq{1, 5, 0}, resultExpected: true},
		{name: "StartUnequal",
			value: GetDatapointValueReq{1, 5, 0},
			other: GetDatapointValueReq{2, 5, 0}, resultExpected: false},
		{name: "NumberUnequal",
			value: GetDatapointValueReq{1, 5, 0},
			other: GetDatapointValueReq{1, 6, 0}, resultExpected: false},
		{name: "FilterUnequal",
			value: GetDatapointValueReq{1, 5, 0},
			other: GetDatapointValueReq{1, 5, 1}, resultExpected: false},
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

func TestGetDatapointValueReqMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value GetDatapointValueReq
	}

	testCases := []testCase{
		{name: "Empty", value: GetDatapointValueReq{0, 0, 0}},
		{name: "Nonempty", value: GetDatapointValueReq{2, 4, ValidDatapoints}},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.value.MarshalBinary()
			if err != nil {
				t.Fatalf("%s", err.Error())
			}
			other := GetDatapointValueReq{}
			if err := other.UnmarshalBinary(data); err != nil {
				t.Fatalf("%s", err.Error())
			}
			if !tc.value.Equal(&other) {
				t.Errorf("incorrect result: %s=%s expected to be equal", tc.value.String(), other.String())
			}
		})
	}
}

func TestGetDatapointValueReqUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name        string
		data        []byte
		errExpected bool
	}

	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0x05, 0x00}, errExpected: true},
		{name: "WrongService", data: []byte{0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := GetDatapointValueReq{}
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

func TestGetDatapointValueResValid(t *testing.T) {
	type testCase struct {
		name          string
		value         GetDatapointValueRes
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Empty", value: GetDatapointValueRes{}, errorExpected: true},
		{name: "WithValues", value: GetDatapointValueRes{
			StartDatapoint: 1, Values: []Datapoint{{ID: 1, State: 0x10, Data: []byte{0x01}}}}, errorExpected: false},
		{name: "WithErrorCode", value: GetDatapointValueRes{
			StartDatapoint: 1, ErrorCode: 5}, errorExpected: false},
		{name: "ValuesAndErrorCode", value: GetDatapointValueRes{
			StartDatapoint: 1, Values: []Datapoint{{ID: 1, Data: []byte{0x01}}}, ErrorCode: 5}, errorExpected: true},
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

func TestGetDatapointValueResEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          GetDatapointValueRes
		other          GetDatapointValueRes
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal",
			value: GetDatapointValueRes{StartDatapoint: 1, Values: []Datapoint{{ID: 1, State: 0x10, Data: []byte{0x01}}}},
			other: GetDatapointValueRes{StartDatapoint: 1, Values: []Datapoint{{ID: 1, State: 0x10, Data: []byte{0x01}}}}, resultExpected: true},
		{name: "StartUnequal",
			value: GetDatapointValueRes{StartDatapoint: 1},
			other: GetDatapointValueRes{StartDatapoint: 2}, resultExpected: false},
		{name: "ErrorCodeUnequal",
			value: GetDatapointValueRes{ErrorCode: 1},
			other: GetDatapointValueRes{ErrorCode: 2}, resultExpected: false},
		{name: "ValuesLenUnequal",
			value: GetDatapointValueRes{Values: []Datapoint{{ID: 1}}},
			other: GetDatapointValueRes{Values: []Datapoint{}}, resultExpected: false},
		{name: "ValueDataUnequal",
			value: GetDatapointValueRes{Values: []Datapoint{{ID: 1, Data: []byte{0x01}}}},
			other: GetDatapointValueRes{Values: []Datapoint{{ID: 1, Data: []byte{0x02}}}}, resultExpected: false},
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

func TestGetDatapointValueResMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value GetDatapointValueRes
	}

	testCases := []testCase{
		{name: "Empty", value: GetDatapointValueRes{}},
		{name: "WithValues", value: GetDatapointValueRes{StartDatapoint: 2, Values: []Datapoint{
			{ID: 10, State: 0x10, Data: []byte{0x01}}, {ID: 11, State: 0x20, Data: []byte{0x03, 0x04}}}}},
		{name: "WithSubByteValue", value: GetDatapointValueRes{StartDatapoint: 1, Values: []Datapoint{
			{ID: 1, State: 0x10, Data: []byte{0x0A}}}}},
		{name: "Error", value: GetDatapointValueRes{StartDatapoint: 3, ErrorCode: 5}},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.value.MarshalBinary()
			if err != nil {
				t.Fatalf("marshal: %s", err.Error())
			}
			other := GetDatapointValueRes{}
			if err := other.UnmarshalBinary(data); err != nil {
				t.Fatalf("unmarshal: %s", err.Error())
			}
			if !tc.value.Equal(&other) {
				t.Errorf("round-trip failed: expected equal values")
			}
		})
	}
}

func TestGetDatapointValueResUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name        string
		data        []byte
		errExpected bool
	}

	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0x85, 0x00}, errExpected: true},
		{name: "WrongService", data: []byte{0x00, 0x85, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "TruncatedValues", data: []byte{0xF0, 0x85, 0x00, 0x00, 0x00, 0x01}, errExpected: true},
		{name: "InvalidLength", data: []byte{0xF0, 0x85, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x10, 0xFF}, errExpected: true},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := GetDatapointValueRes{}
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
