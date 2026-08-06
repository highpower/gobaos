package baos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.resultExpected, tc.value.Equal(&tc.other))
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			require.NoError(t, err)

			other := SetServerItemReq{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.True(t, tc.value.Equal(&other))
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := SetServerItemReq{}
			err := r.UnmarshalBinary(tc.data)
			if tc.errExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.resultExpected, tc.value.Equal(&tc.other))
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			require.NoError(t, err)

			other := SetServerItemRes{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.True(t, tc.value.Equal(&other))
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := SetServerItemRes{}
			err := r.UnmarshalBinary(tc.data)
			if tc.errExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func FuzzSetServerItemReq(f *testing.F) {
	req := SetServerItemReq{Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}
	data, _ := req.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := SetServerItemReq{}
		_ = other.UnmarshalBinary(data)
	})
}

func FuzzSetServerItemRes(f *testing.F) {
	res := SetServerItemRes{StartItem: 1, NumberOfItems: 0, ErrorCode: 0}
	data, _ := res.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := SetServerItemRes{}
		_ = other.UnmarshalBinary(data)
	})
}
