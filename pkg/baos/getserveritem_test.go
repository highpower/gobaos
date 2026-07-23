package baos

import (
	"testing"
)

func TestGetServerItemReqValid(t *testing.T) {
	type testCase struct {
		name          string
		value         GetServerItemReq
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Empty", value: GetServerItemReq{0, 0}, errorExpected: false},
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

func TestGetServerItemReqEqual(t *testing.T) {

	type testCase struct {
		name           string
		value          GetServerItemReq
		other          GetServerItemReq
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal", value: GetServerItemReq{0, 5},
			other: GetServerItemReq{0, 5}, resultExpected: true},
		{name: "StartUnequal", value: GetServerItemReq{1, 5},
			other: GetServerItemReq{2, 5}, resultExpected: false},
		{name: "NumberUnequal", value: GetServerItemReq{1, 5},
			other: GetServerItemReq{1, 6}, resultExpected: false},
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

func TestGetServerItemReqMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value GetServerItemReq
	}
	testCases := []testCase{
		{name: "Empty", value: GetServerItemReq{0, 0}},
		{name: "Nonempty", value: GetServerItemReq{2, 4}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			if err != nil {
				t.Fatalf("%s", err.Error())
			}
			other := GetServerItemReq{}
			if err := other.UnmarshalBinary(data); err != nil {
				t.Fatalf("%s", err.Error())
			}
			if !tc.value.Equal(&other) {
				t.Errorf("incorrect result: %s=%s expected to be equal", tc.value.String(), other.String())
			}
		})
	}
}

func TestGetServerItemResValid(t *testing.T) {
	type testCase struct {
		name          string
		value         GetServerItemRes
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Empty", value: GetServerItemRes{}, errorExpected: true},
		{name: "WithItems", value: GetServerItemRes{StartItem: 1, Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, errorExpected: false},
		{name: "WithErrorCode", value: GetServerItemRes{StartItem: 1, ErrorCode: 5}, errorExpected: false},
		{name: "ItemsAndErrorCode",
			value: GetServerItemRes{StartItem: 1, Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}, ErrorCode: 5}, errorExpected: true},
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

func TestGetServerItemResEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          GetServerItemRes
		other          GetServerItemRes
		resultExpected bool
	}
	testCases := []testCase{
		{name: "Equal", value: GetServerItemRes{StartItem: 1, Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, other: GetServerItemRes{StartItem: 1, Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, resultExpected: true},
		{name: "StartUnequal", value: GetServerItemRes{StartItem: 1}, other: GetServerItemRes{StartItem: 2}, resultExpected: false},
		{name: "ErrorCodeUnequal", value: GetServerItemRes{ErrorCode: 1}, other: GetServerItemRes{ErrorCode: 2}, resultExpected: false},
		{name: "ItemsLenUnequal", value: GetServerItemRes{Items: []ServerItem{{ID: 1}}}, other: GetServerItemRes{Items: []ServerItem{}}, resultExpected: false},
		{name: "ItemDataUnequal", value: GetServerItemRes{Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, other: GetServerItemRes{Items: []ServerItem{{ID: 1, Data: []byte{0x02}}}}, resultExpected: false},
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

func TestGetServerItemResMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value GetServerItemRes
	}
	testCases := []testCase{
		{name: "Empty", value: GetServerItemRes{}},
		{name: "WithItems", value: GetServerItemRes{StartItem: 2, Items: []ServerItem{{ID: 10, Data: []byte{0x01, 0x02}}, {ID: 11, Data: []byte{0x03}}}}},
		{name: "Error", value: GetServerItemRes{StartItem: 3, ErrorCode: 5}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			if err != nil {
				t.Fatalf("marshal: %s", err.Error())
			}
			other := GetServerItemRes{}
			if err := other.UnmarshalBinary(data); err != nil {
				t.Fatalf("unmarshal: %s", err.Error())
			}
			if !tc.value.Equal(&other) {
				t.Errorf("round-trip failed: expected equal values")
			}
		})
	}
}

func TestGetServerItemResUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name        string
		data        []byte
		errExpected bool
	}
	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0x81, 0x00}, errExpected: true},
		{name: "WrongService", data: []byte{0x00, 0x81, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "TruncatedItems", data: []byte{0xF0, 0x81, 0x00, 0x00, 0x00, 0x01}, errExpected: true},
		{name: "ItemDataTooShort", data: []byte{0xF0, 0x81, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x02}, errExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := GetServerItemRes{}
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
