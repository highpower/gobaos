package baos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDescriptionStringReqEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          getDescriptionStringReq
		other          getDescriptionStringReq
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal", value: getDescriptionStringReq{0, 5},
			other: getDescriptionStringReq{0, 5}, resultExpected: true},
		{name: "StartUnequal", value: getDescriptionStringReq{1, 5},
			other: getDescriptionStringReq{2, 5}, resultExpected: false},
		{name: "NumberUnequal", value: getDescriptionStringReq{1, 5},
			other: getDescriptionStringReq{1, 6}, resultExpected: false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.resultExpected, tc.value.equal(&tc.other))
		})
	}
}

func TestGetDescriptionStringReqMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value getDescriptionStringReq
	}
	testCases := []testCase{
		{name: "Empty", value: getDescriptionStringReq{0, 0}},
		{name: "Nonempty", value: getDescriptionStringReq{2, 4}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			require.NoError(t, err)

			other := getDescriptionStringReq{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.True(t, tc.value.equal(&other))
		})
	}
}

func TestGetDescriptionStringReqUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name          string
		data          []byte
		errorExpected bool
	}
	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0x04, 0x00}, errorExpected: true},
		{name: "WrongService", data: []byte{0x00, 0x04, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
		{name: "Valid", data: []byte{0xF0, 0x04, 0x00, 0x05, 0x00, 0x0A}, errorExpected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := getDescriptionStringReq{}
			err := r.UnmarshalBinary(tc.data)
			if tc.errorExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestGetDescriptionStringResEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          getDescriptionStringRes
		other          getDescriptionStringRes
		resultExpected bool
	}
	testCases := []testCase{
		{name: "Equal", value: getDescriptionStringRes{startDatapoint: 1, values: []string{"a", "b"}}, other: getDescriptionStringRes{startDatapoint: 1, values: []string{"a", "b"}}, resultExpected: true},
		{name: "StartUnequal", value: getDescriptionStringRes{startDatapoint: 1}, other: getDescriptionStringRes{startDatapoint: 2}, resultExpected: false},
		{name: "ErrorCodeUnequal", value: getDescriptionStringRes{errorCode: 1}, other: getDescriptionStringRes{errorCode: 2}, resultExpected: false},
		{name: "StringsLenUnequal", value: getDescriptionStringRes{values: []string{"a"}}, other: getDescriptionStringRes{values: []string{}}, resultExpected: false},
		{name: "StringContentUnequal", value: getDescriptionStringRes{values: []string{"a"}}, other: getDescriptionStringRes{values: []string{"b"}}, resultExpected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.resultExpected, tc.value.equal(&tc.other))
		})
	}
}

func TestGetDescriptionStringResMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value getDescriptionStringRes
	}
	testCases := []testCase{
		{name: "Empty", value: getDescriptionStringRes{}},
		{name: "WithStrings", value: getDescriptionStringRes{startDatapoint: 2, values: []string{"hello", "world"}}},
		{name: "WithEmptyString", value: getDescriptionStringRes{startDatapoint: 3, values: []string{""}}},
		{name: "Error", value: getDescriptionStringRes{startDatapoint: 3, errorCode: 5}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			require.NoError(t, err)

			other := getDescriptionStringRes{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.True(t, tc.value.equal(&other))
		})
	}
}

func TestGetDescriptionStringResUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name          string
		data          []byte
		errorExpected bool
	}
	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0x84, 0x00}, errorExpected: true},
		{name: "WrongService", data: []byte{0x00, 0x84, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
		{name: "TruncatedStrings", data: []byte{0xF0, 0x84, 0x00, 0x00, 0x00, 0x01}, errorExpected: true},
		{name: "StringDataTooShort", data: []byte{0xF0, 0x84, 0x00, 0x00, 0x00, 0x01, 0x00, 0x05, 0x01}, errorExpected: true},
		{name: "ExtraBytes", data: []byte{0xF0, 0x84, 0x00, 0x00, 0x00, 0x01, 0x00, 0x02, 0xAB, 0xCD, 0xEF}, errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := getDescriptionStringRes{}
			err := r.UnmarshalBinary(tc.data)
			if tc.errorExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func FuzzGetDescriptionStringReq(f *testing.F) {
	req := getDescriptionStringReq{startDatapoint: 1, numberOfStrings: 2}
	data, _ := req.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := getDescriptionStringReq{}
		_ = other.UnmarshalBinary(data)
	})
}

func FuzzGetDescriptionStringRes(f *testing.F) {
	res := getDescriptionStringRes{startDatapoint: 1, values: []string{"test", "data"}}
	data, _ := res.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := getDescriptionStringRes{}
		_ = other.UnmarshalBinary(data)
	})
}
