package baos

import (
	"testing"
)

func TestSetServerItemReqValid(t *testing.T) {
	type testCase struct {
		name          string
		value         SetServerItemReq
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Empty", value: SetServerItemReq{}, errorExpected: true},
		{name: "WithItems", value: SetServerItemReq{Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, errorExpected: false},
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

func TestSetServerItemReqEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          SetServerItemReq
		other          SetServerItemReq
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal", value: SetServerItemReq{Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, other: SetServerItemReq{Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, resultExpected: true},
		{name: "ItemIDUnequal", value: SetServerItemReq{Items: []ServerItem{{ID: 1}}}, other: SetServerItemReq{Items: []ServerItem{{ID: 2}}}, resultExpected: false},
		{name: "ItemDataUnequal", value: SetServerItemReq{Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, other: SetServerItemReq{Items: []ServerItem{{ID: 1, Data: []byte{0x02}}}}, resultExpected: false},
		{name: "ItemsLenUnequal", value: SetServerItemReq{Items: []ServerItem{{ID: 1}}}, other: SetServerItemReq{Items: []ServerItem{}}, resultExpected: false},
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

func TestSetServerItemReqMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value SetServerItemReq
	}

	testCases := []testCase{
		{name: "SingleItem", value: SetServerItemReq{Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}},
		{name: "MultipleItems", value: SetServerItemReq{Items: []ServerItem{{ID: 10, Data: []byte{0x01, 0x02}}, {ID: 11, Data: []byte{0x03}}}}},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.value.MarshalBinary()
			if err != nil {
				t.Fatalf("marshal: %s", err.Error())
			}
			other := SetServerItemReq{}
			if err := other.UnmarshalBinary(data); err != nil {
				t.Fatalf("unmarshal: %s", err.Error())
			}
			if !tc.value.Equal(&other) {
				t.Errorf("round-trip failed: expected equal values")
			}
		})
	}
}

func TestSetServerItemReqUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name        string
		data        []byte
		errExpected bool
	}

	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0x02, 0x00}, errExpected: true},
		{name: "WrongService", data: []byte{0x00, 0x02, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "TruncatedItems", data: []byte{0xF0, 0x02, 0x00, 0x00, 0x00, 0x01}, errExpected: true},
		{name: "ItemDataTooShort", data: []byte{0xF0, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x02}, errExpected: true},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := SetServerItemReq{}
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

func TestSetServerItemResValid(t *testing.T) {
	type testCase struct {
		name          string
		value         SetServerItemRes
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Empty", value: SetServerItemRes{}, errorExpected: true},
		{name: "WithNumberOfItems", value: SetServerItemRes{StartItem: 1, NumberOfItems: 5}, errorExpected: false},
		{name: "WithErrorCode", value: SetServerItemRes{StartItem: 1, ErrorCode: 5}, errorExpected: false},
		{name: "ItemsAndErrorCode", value: SetServerItemRes{StartItem: 1, NumberOfItems: 5, ErrorCode: 5}, errorExpected: true},
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

func TestSetServerItemResEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          SetServerItemRes
		other          SetServerItemRes
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal", value: SetServerItemRes{StartItem: 1, NumberOfItems: 5, ErrorCode: 0}, other: SetServerItemRes{StartItem: 1, NumberOfItems: 5, ErrorCode: 0}, resultExpected: true},
		{name: "StartUnequal", value: SetServerItemRes{StartItem: 1}, other: SetServerItemRes{StartItem: 2}, resultExpected: false},
		{name: "NumberUnequal", value: SetServerItemRes{NumberOfItems: 1}, other: SetServerItemRes{NumberOfItems: 2}, resultExpected: false},
		{name: "ErrorCodeUnequal", value: SetServerItemRes{ErrorCode: 1}, other: SetServerItemRes{ErrorCode: 2}, resultExpected: false},
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

func TestSetServerItemResMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value SetServerItemRes
	}

	testCases := []testCase{
		{name: "Success", value: SetServerItemRes{StartItem: 2, NumberOfItems: 4}},
		{name: "Error", value: SetServerItemRes{StartItem: 3, ErrorCode: 5}},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.value.MarshalBinary()
			if err != nil {
				t.Fatalf("marshal: %s", err.Error())
			}
			other := SetServerItemRes{}
			if err := other.UnmarshalBinary(data); err != nil {
				t.Fatalf("unmarshal: %s", err.Error())
			}
			if !tc.value.Equal(&other) {
				t.Errorf("round-trip failed: expected equal values")
			}
		})
	}
}

func TestSetServerItemResUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name        string
		data        []byte
		errExpected bool
	}

	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0x82, 0x00}, errExpected: true},
		{name: "WrongService", data: []byte{0x00, 0x82, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := SetServerItemRes{}
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
