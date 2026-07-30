package baos

import (
	"testing"
)

func TestGetDescriptionStringReqValid(t *testing.T) {
	type testCase struct {
		name          string
		value         GetDescriptionStringReq
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Empty", value: GetDescriptionStringReq{0, 0}, errorExpected: false},
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

func TestGetDescriptionStringReqEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          GetDescriptionStringReq
		other          GetDescriptionStringReq
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal", value: GetDescriptionStringReq{0, 5},
			other: GetDescriptionStringReq{0, 5}, resultExpected: true},
		{name: "StartUnequal", value: GetDescriptionStringReq{1, 5},
			other: GetDescriptionStringReq{2, 5}, resultExpected: false},
		{name: "NumberUnequal", value: GetDescriptionStringReq{1, 5},
			other: GetDescriptionStringReq{1, 6}, resultExpected: false},
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

func TestGetDescriptionStringReqMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value GetDescriptionStringReq
	}
	testCases := []testCase{
		{name: "Empty", value: GetDescriptionStringReq{0, 0}},
		{name: "Nonempty", value: GetDescriptionStringReq{2, 4}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			if err != nil {
				t.Fatalf("%s", err.Error())
			}
			other := GetDescriptionStringReq{}
			if err := other.UnmarshalBinary(data); err != nil {
				t.Fatalf("%s", err.Error())
			}
			if !tc.value.Equal(&other) {
				t.Errorf("incorrect result: %s=%s expected to be equal", tc.value.String(), other.String())
			}
		})
	}
}

func TestGetDescriptionStringReqUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name        string
		data        []byte
		errExpected bool
	}
	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0x04, 0x00}, errExpected: true},
		{name: "WrongService", data: []byte{0x00, 0x04, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "Valid", data: []byte{0xF0, 0x04, 0x00, 0x05, 0x00, 0x0A}, errExpected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := GetDescriptionStringReq{}
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

func TestGetDescriptionStringResValid(t *testing.T) {
	type testCase struct {
		name          string
		value         GetDescriptionStringRes
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Empty", value: GetDescriptionStringRes{}, errorExpected: true},
		{name: "WithStrings", value: GetDescriptionStringRes{StartDatapoint: 1, Strings: []string{"test"}}, errorExpected: false},
		{name: "WithErrorCode", value: GetDescriptionStringRes{StartDatapoint: 1, ErrorCode: 5}, errorExpected: false},
		{name: "StringsAndErrorCode",
			value: GetDescriptionStringRes{StartDatapoint: 1, Strings: []string{"test"}, ErrorCode: 5}, errorExpected: true},
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

func TestGetDescriptionStringResEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          GetDescriptionStringRes
		other          GetDescriptionStringRes
		resultExpected bool
	}
	testCases := []testCase{
		{name: "Equal", value: GetDescriptionStringRes{StartDatapoint: 1, Strings: []string{"a", "b"}}, other: GetDescriptionStringRes{StartDatapoint: 1, Strings: []string{"a", "b"}}, resultExpected: true},
		{name: "StartUnequal", value: GetDescriptionStringRes{StartDatapoint: 1}, other: GetDescriptionStringRes{StartDatapoint: 2}, resultExpected: false},
		{name: "ErrorCodeUnequal", value: GetDescriptionStringRes{ErrorCode: 1}, other: GetDescriptionStringRes{ErrorCode: 2}, resultExpected: false},
		{name: "StringsLenUnequal", value: GetDescriptionStringRes{Strings: []string{"a"}}, other: GetDescriptionStringRes{Strings: []string{}}, resultExpected: false},
		{name: "StringContentUnequal", value: GetDescriptionStringRes{Strings: []string{"a"}}, other: GetDescriptionStringRes{Strings: []string{"b"}}, resultExpected: false},
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

func TestGetDescriptionStringResMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value GetDescriptionStringRes
	}
	testCases := []testCase{
		{name: "Empty", value: GetDescriptionStringRes{}},
		{name: "WithStrings", value: GetDescriptionStringRes{StartDatapoint: 2, Strings: []string{"hello", "world"}}},
		{name: "WithEmptyString", value: GetDescriptionStringRes{StartDatapoint: 3, Strings: []string{""}}},
		{name: "Error", value: GetDescriptionStringRes{StartDatapoint: 3, ErrorCode: 5}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			if err != nil {
				t.Fatalf("marshal: %s", err.Error())
			}
			other := GetDescriptionStringRes{}
			if err := other.UnmarshalBinary(data); err != nil {
				t.Fatalf("unmarshal: %s", err.Error())
			}
			if !tc.value.Equal(&other) {
				t.Errorf("round-trip failed: expected equal values")
			}
		})
	}
}

func TestGetDescriptionStringResUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name        string
		data        []byte
		errExpected bool
	}
	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0x84, 0x00}, errExpected: true},
		{name: "WrongService", data: []byte{0x00, 0x84, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "TruncatedStrings", data: []byte{0xF0, 0x84, 0x00, 0x00, 0x00, 0x01}, errExpected: true},
		{name: "StringDataTooShort", data: []byte{0xF0, 0x84, 0x00, 0x00, 0x00, 0x01, 0x00, 0x05, 0x01}, errExpected: true},
		{name: "ExtraBytes", data: []byte{0xF0, 0x84, 0x00, 0x00, 0x00, 0x01, 0x00, 0x02, 0xAB, 0xCD, 0xEF}, errExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := GetDescriptionStringRes{}
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

func FuzzGetDescriptionStringReq(f *testing.F) {
	req := GetDescriptionStringReq{StartDatapoint: 1, NumberOfStrings: 2}
	data, _ := req.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := GetDescriptionStringReq{}
		_ = other.UnmarshalBinary(data)
	})
}

func FuzzGetDescriptionStringRes(f *testing.F) {
	res := GetDescriptionStringRes{StartDatapoint: 1, Strings: []string{"test", "data"}}
	data, _ := res.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := GetDescriptionStringRes{}
		_ = other.UnmarshalBinary(data)
	})
}
