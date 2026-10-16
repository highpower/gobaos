package baos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetServerItemReqEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          setServerItemReq
		other          setServerItemReq
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal", value: setServerItemReq{items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, other: setServerItemReq{items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, resultExpected: true},
		{name: "ItemIDUnequal", value: setServerItemReq{items: []ServerItem{{ID: 1}}}, other: setServerItemReq{items: []ServerItem{{ID: 2}}}, resultExpected: false},
		{name: "ItemDataUnequal", value: setServerItemReq{items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, other: setServerItemReq{items: []ServerItem{{ID: 1, Data: []byte{0x02}}}}, resultExpected: false},
		{name: "ItemsLenUnequal", value: setServerItemReq{items: []ServerItem{{ID: 1}}}, other: setServerItemReq{items: []ServerItem{}}, resultExpected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.resultExpected, tc.value.equal(&tc.other))
		})
	}
}

func TestSetServerItemReqMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value setServerItemReq
	}

	testCases := []testCase{
		{name: "SingleItem", value: setServerItemReq{items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}},
		{name: "MultipleItems", value: setServerItemReq{items: []ServerItem{{ID: 10, Data: []byte{0x01, 0x02}}, {ID: 11, Data: []byte{0x03}}}}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			require.NoError(t, err)

			other := setServerItemReq{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.True(t, tc.value.equal(&other))
		})
	}
}

func TestSetServerItemReqUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name          string
		data          []byte
		errorExpected bool
	}

	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0x02, 0x00}, errorExpected: true},
		{name: "WrongService", data: []byte{0x00, 0x02, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
		{name: "TruncatedItems", data: []byte{0xF0, 0x02, 0x00, 0x00, 0x00, 0x01}, errorExpected: true},
		{name: "ItemDataTooShort", data: []byte{0xF0, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x02}, errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := setServerItemReq{}
			err := r.UnmarshalBinary(tc.data)
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
		value          setServerItemRes
		other          setServerItemRes
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal", value: setServerItemRes{startItem: 1, numberOfItems: 5, errorCode: 0}, other: setServerItemRes{startItem: 1, numberOfItems: 5, errorCode: 0}, resultExpected: true},
		{name: "StartUnequal", value: setServerItemRes{startItem: 1}, other: setServerItemRes{startItem: 2}, resultExpected: false},
		{name: "NumberUnequal", value: setServerItemRes{numberOfItems: 1}, other: setServerItemRes{numberOfItems: 2}, resultExpected: false},
		{name: "ErrorCodeUnequal", value: setServerItemRes{errorCode: 1}, other: setServerItemRes{errorCode: 2}, resultExpected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.resultExpected, tc.value.equal(&tc.other))
		})
	}
}

func TestSetServerItemResMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value setServerItemRes
	}

	testCases := []testCase{
		{name: "Success", value: setServerItemRes{startItem: 2, numberOfItems: 4}},
		{name: "Error", value: setServerItemRes{startItem: 3, errorCode: 5}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			require.NoError(t, err)

			other := setServerItemRes{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.True(t, tc.value.equal(&other))
		})
	}
}

func TestSetServerItemResUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name          string
		data          []byte
		errorExpected bool
	}

	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0x82, 0x00}, errorExpected: true},
		{name: "WrongService", data: []byte{0x00, 0x82, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := setServerItemRes{}
			err := r.UnmarshalBinary(tc.data)
			if tc.errorExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func FuzzSetServerItemReq(f *testing.F) {
	req := setServerItemReq{items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}
	data, _ := req.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := setServerItemReq{}
		_ = other.UnmarshalBinary(data)
	})
}

func FuzzSetServerItemRes(f *testing.F) {
	res := setServerItemRes{startItem: 1, numberOfItems: 0, errorCode: 0}
	data, _ := res.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := setServerItemRes{}
		_ = other.UnmarshalBinary(data)
	})
}
