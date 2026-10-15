package baos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			if tc.errorExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
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
			assert.Equal(t, tc.resultExpected, tc.value.Equal(&tc.other))
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
			require.NoError(t, err)

			other := GetServerItemReq{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.True(t, tc.value.Equal(&other))
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
			if tc.errorExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
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
			assert.Equal(t, tc.resultExpected, tc.value.Equal(&tc.other))
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
			require.NoError(t, err)

			other := GetServerItemRes{}
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
			r := GetServerItemRes{}
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
	req := GetServerItemReq{StartItem: 1, NumberOfItems: 2}
	data, _ := req.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := GetServerItemReq{}
		_ = other.UnmarshalBinary(data)
	})
}

func FuzzGetServerItemRes(f *testing.F) {
	res := GetServerItemRes{StartItem: 1, Items: []ServerItem{{ID: 1, Data: []byte{0x01, 0x02}}}}
	data, _ := res.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := GetServerItemRes{}
		_ = other.UnmarshalBinary(data)
	})
}
