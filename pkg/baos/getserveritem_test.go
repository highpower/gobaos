package baos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetServerItemReqEqual(t *testing.T) {

	type testCase struct {
		name           string
		value          getServerItemReq
		other          getServerItemReq
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal", value: getServerItemReq{0, 5},
			other: getServerItemReq{0, 5}, resultExpected: true},
		{name: "StartUnequal", value: getServerItemReq{1, 5},
			other: getServerItemReq{2, 5}, resultExpected: false},
		{name: "NumberUnequal", value: getServerItemReq{1, 5},
			other: getServerItemReq{1, 6}, resultExpected: false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.resultExpected, tc.value.Equal(&tc.other))
		})
	}

}

func TestGetServerItemReqMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value getServerItemReq
	}
	testCases := []testCase{
		{name: "Empty", value: getServerItemReq{0, 0}},
		{name: "Nonempty", value: getServerItemReq{2, 4}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			require.NoError(t, err)

			other := getServerItemReq{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.True(t, tc.value.Equal(&other))
		})
	}
}

func TestGetServerItemResEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          getServerItemRes
		other          getServerItemRes
		resultExpected bool
	}
	testCases := []testCase{
		{name: "Equal", value: getServerItemRes{StartItem: 1, Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, other: getServerItemRes{StartItem: 1, Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, resultExpected: true},
		{name: "StartUnequal", value: getServerItemRes{StartItem: 1}, other: getServerItemRes{StartItem: 2}, resultExpected: false},
		{name: "ErrorCodeUnequal", value: getServerItemRes{ErrorCode: 1}, other: getServerItemRes{ErrorCode: 2}, resultExpected: false},
		{name: "ItemsLenUnequal", value: getServerItemRes{Items: []ServerItem{{ID: 1}}}, other: getServerItemRes{Items: []ServerItem{}}, resultExpected: false},
		{name: "ItemDataUnequal", value: getServerItemRes{Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, other: getServerItemRes{Items: []ServerItem{{ID: 1, Data: []byte{0x02}}}}, resultExpected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.resultExpected, tc.value.Equal(&tc.other))
		})
	}
}

func TestGetServerItemResMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value getServerItemRes
	}
	testCases := []testCase{
		{name: "Empty", value: getServerItemRes{}},
		{name: "WithItems", value: getServerItemRes{StartItem: 2, Items: []ServerItem{{ID: 10, Data: []byte{0x01, 0x02}}, {ID: 11, Data: []byte{0x03}}}}},
		{name: "Error", value: getServerItemRes{StartItem: 3, ErrorCode: 5}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			require.NoError(t, err)

			other := getServerItemRes{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.True(t, tc.value.Equal(&other))
		})
	}
}

func TestGetServerItemResUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name          string
		data          []byte
		errorExpected bool
	}
	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0x81, 0x00}, errorExpected: true},
		{name: "WrongService", data: []byte{0x00, 0x81, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
		{name: "TruncatedItems", data: []byte{0xF0, 0x81, 0x00, 0x00, 0x00, 0x01}, errorExpected: true},
		{name: "ItemDataTooShort", data: []byte{0xF0, 0x81, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x02}, errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := getServerItemRes{}
			err := r.UnmarshalBinary(tc.data)
			if tc.errorExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func FuzzGetServerItemReq(f *testing.F) {
	req := getServerItemReq{StartItem: 1, NumberOfItems: 2}
	data, _ := req.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := getServerItemReq{}
		_ = other.UnmarshalBinary(data)
	})
}

func FuzzGetServerItemRes(f *testing.F) {
	res := getServerItemRes{StartItem: 1, Items: []ServerItem{{ID: 1, Data: []byte{0x01, 0x02}}}}
	data, _ := res.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := getServerItemRes{}
		_ = other.UnmarshalBinary(data)
	})
}
