package baos

import (
	"testing"
)

func TestGetDatapointDescriptionReqValid(t *testing.T) {
	type testCase struct {
		name          string
		value         GetDatapointDescriptionReq
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Empty", value: GetDatapointDescriptionReq{0, 0}, errorExpected: false},
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

func TestGetDatapointDescriptionReqEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          GetDatapointDescriptionReq
		other          GetDatapointDescriptionReq
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal", value: GetDatapointDescriptionReq{0, 5},
			other: GetDatapointDescriptionReq{0, 5}, resultExpected: true},
		{name: "StartUnequal", value: GetDatapointDescriptionReq{1, 5},
			other: GetDatapointDescriptionReq{2, 5}, resultExpected: false},
		{name: "NumberUnequal", value: GetDatapointDescriptionReq{1, 5},
			other: GetDatapointDescriptionReq{1, 6}, resultExpected: false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if res := tc.value.Equal(&tc.other); res != tc.resultExpected {
				t.Errorf("incorrect result: %s=%s is %t, expected %t", tc.value.String(), tc.other.String(),
					res, tc.resultExpected)
			}
		})
	}
}

func TestGetDatapointDescriptionReqMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value GetDatapointDescriptionReq
	}
	testCases := []testCase{
		{name: "Empty", value: GetDatapointDescriptionReq{0, 0}},
		{name: "Nonempty", value: GetDatapointDescriptionReq{2, 4}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			if err != nil {
				t.Fatalf("%s", err.Error())
			}
			other := GetDatapointDescriptionReq{}
			if err := other.UnmarshalBinary(data); err != nil {
				t.Fatalf("%s", err.Error())
			}
			if !tc.value.Equal(&other) {
				t.Errorf("incorrect result: %s=%s expected to be equal", tc.value.String(), other.String())
			}
		})
	}
}

func TestGetDatapointDescriptionReqUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name        string
		data        []byte
		errExpected bool
	}
	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0x03, 0x00}, errExpected: true},
		{name: "WrongService", data: []byte{0x00, 0x03, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "Valid", data: []byte{0xF0, 0x03, 0x00, 0x02, 0x00, 0x04}, errExpected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := GetDatapointDescriptionReq{}
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

func TestGetDatapointDescriptionReqUnmarshalFields(t *testing.T) {
	r := GetDatapointDescriptionReq{}
	data := []byte{0xF0, 0x03, 0x00, 0x05, 0x00, 0x0A}
	if err := r.UnmarshalBinary(data); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if r.StartDatapoint != 5 {
		t.Errorf("StartDatapoint: got %d, want 5", r.StartDatapoint)
	}
	if r.NumberOfDatapoints != 10 {
		t.Errorf("NumberOfDatapoints: got %d, want 10", r.NumberOfDatapoints)
	}
}

func TestGetDatapointDescriptionResValid(t *testing.T) {
	type testCase struct {
		name          string
		value         GetDatapointDescriptionRes
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Empty", value: GetDatapointDescriptionRes{}, errorExpected: true},
		{name: "WithDescriptions", value: GetDatapointDescriptionRes{StartDatapoint: 1, Descriptions: []DatapointDescription{{ID: 1, ValueType: 7, ConfigFlags: 0x02, DPT: 1}}}, errorExpected: false},
		{name: "WithErrorCode", value: GetDatapointDescriptionRes{StartDatapoint: 1, ErrorCode: 5}, errorExpected: false},
		{name: "DescriptionsAndErrorCode",
			value: GetDatapointDescriptionRes{StartDatapoint: 1, Descriptions: []DatapointDescription{{ID: 1}}, ErrorCode: 5}, errorExpected: true},
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

func TestGetDatapointDescriptionResEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          GetDatapointDescriptionRes
		other          GetDatapointDescriptionRes
		resultExpected bool
	}
	testCases := []testCase{
		{name: "Equal", value: GetDatapointDescriptionRes{StartDatapoint: 1, Descriptions: []DatapointDescription{{ID: 1, ValueType: 7, ConfigFlags: 0x02, DPT: 1}}}, other: GetDatapointDescriptionRes{StartDatapoint: 1, Descriptions: []DatapointDescription{{ID: 1, ValueType: 7, ConfigFlags: 0x02, DPT: 1}}}, resultExpected: true},
		{name: "StartUnequal", value: GetDatapointDescriptionRes{StartDatapoint: 1}, other: GetDatapointDescriptionRes{StartDatapoint: 2}, resultExpected: false},
		{name: "ErrorCodeUnequal", value: GetDatapointDescriptionRes{ErrorCode: 1}, other: GetDatapointDescriptionRes{ErrorCode: 2}, resultExpected: false},
		{name: "DescriptionsLenUnequal", value: GetDatapointDescriptionRes{Descriptions: []DatapointDescription{{ID: 1}}}, other: GetDatapointDescriptionRes{Descriptions: []DatapointDescription{}}, resultExpected: false},
		{name: "DescriptionDPTUnequal", value: GetDatapointDescriptionRes{Descriptions: []DatapointDescription{{ID: 1, DPT: 1}}}, other: GetDatapointDescriptionRes{Descriptions: []DatapointDescription{{ID: 1, DPT: 2}}}, resultExpected: false},
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

func TestGetDatapointDescriptionResMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value GetDatapointDescriptionRes
	}
	testCases := []testCase{
		{name: "Empty", value: GetDatapointDescriptionRes{}},
		{name: "WithDescriptions", value: GetDatapointDescriptionRes{StartDatapoint: 2, Descriptions: []DatapointDescription{{ID: 10, ValueType: 7, ConfigFlags: 0x02, DPT: 1}, {ID: 11, ValueType: 8, ConfigFlags: 0x04, DPT: 2}}}},
		{name: "Error", value: GetDatapointDescriptionRes{StartDatapoint: 3, ErrorCode: 5}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			if err != nil {
				t.Fatalf("marshal: %s", err.Error())
			}
			other := GetDatapointDescriptionRes{}
			if err := other.UnmarshalBinary(data); err != nil {
				t.Fatalf("unmarshal: %s", err.Error())
			}
			if !tc.value.Equal(&other) {
				t.Errorf("round-trip failed: expected equal values")
			}
		})
	}
}

func TestGetDatapointDescriptionResUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name        string
		data        []byte
		errExpected bool
	}
	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0x83, 0x00}, errExpected: true},
		{name: "WrongService", data: []byte{0x00, 0x83, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "TruncatedDescriptions", data: []byte{0xF0, 0x83, 0x00, 0x00, 0x00, 0x01}, errExpected: true},
		{name: "DescriptionDataTooShort", data: []byte{0xF0, 0x83, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x07, 0x02}, errExpected: true},
		{name: "ExtraBytes", data: []byte{0xF0, 0x83, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x07, 0x02, 0x01, 0xAB}, errExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := GetDatapointDescriptionRes{}
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

func FuzzGetDatapointDescriptionReq(f *testing.F) {
	req := GetDatapointDescriptionReq{StartDatapoint: 1, NumberOfDatapoints: 2}
	data, _ := req.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := GetDatapointDescriptionReq{}
		_ = other.UnmarshalBinary(data)
	})
}

func FuzzGetDatapointDescriptionRes(f *testing.F) {
	res := GetDatapointDescriptionRes{StartDatapoint: 1, Descriptions: []DatapointDescription{{ID: 1, ValueType: 7, ConfigFlags: 0x02, DPT: 1}}}
	data, _ := res.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := GetDatapointDescriptionRes{}
		_ = other.UnmarshalBinary(data)
	})
}