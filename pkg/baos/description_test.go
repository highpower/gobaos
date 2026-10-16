package baos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDatapointDescriptionReqEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          getDatapointDescriptionReq
		other          getDatapointDescriptionReq
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal", value: getDatapointDescriptionReq{0, 5},
			other: getDatapointDescriptionReq{0, 5}, resultExpected: true},
		{name: "StartUnequal", value: getDatapointDescriptionReq{1, 5},
			other: getDatapointDescriptionReq{2, 5}, resultExpected: false},
		{name: "NumberUnequal", value: getDatapointDescriptionReq{1, 5},
			other: getDatapointDescriptionReq{1, 6}, resultExpected: false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.resultExpected, tc.value.Equal(&tc.other))
		})
	}
}

func TestGetDatapointDescriptionReqMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value getDatapointDescriptionReq
	}
	testCases := []testCase{
		{name: "Empty", value: getDatapointDescriptionReq{0, 0}},
		{name: "Nonempty", value: getDatapointDescriptionReq{2, 4}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			require.NoError(t, err)

			other := getDatapointDescriptionReq{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.True(t, tc.value.Equal(&other))
		})
	}
}

func TestGetDatapointDescriptionReqUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name          string
		data          []byte
		errorExpected bool
	}
	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0x03, 0x00}, errorExpected: true},
		{name: "WrongService", data: []byte{0x00, 0x03, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
		{name: "Valid", data: []byte{0xF0, 0x03, 0x00, 0x02, 0x00, 0x04}, errorExpected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := getDatapointDescriptionReq{}
			err := r.UnmarshalBinary(tc.data)
			if tc.errorExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestGetDatapointDescriptionReqUnmarshalFields(t *testing.T) {
	r := getDatapointDescriptionReq{}
	data := []byte{0xF0, 0x03, 0x00, 0x05, 0x00, 0x0A}
	require.NoError(t, r.UnmarshalBinary(data))
	assert.Equal(t, DatapointID(5), r.StartDatapoint)
	assert.Equal(t, uint16(10), r.NumberOfDatapoints)
}

func TestGetDatapointDescriptionResEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          getDatapointDescriptionRes
		other          getDatapointDescriptionRes
		resultExpected bool
	}
	testCases := []testCase{
		{name: "Equal", value: getDatapointDescriptionRes{StartDatapoint: 1, Descriptions: []DatapointDescription{{ID: 1, ValueType: 7, ConfigFlags: 0x02, Code: 1}}}, other: getDatapointDescriptionRes{StartDatapoint: 1, Descriptions: []DatapointDescription{{ID: 1, ValueType: 7, ConfigFlags: 0x02, Code: 1}}}, resultExpected: true},
		{name: "StartUnequal", value: getDatapointDescriptionRes{StartDatapoint: 1}, other: getDatapointDescriptionRes{StartDatapoint: 2}, resultExpected: false},
		{name: "ErrorCodeUnequal", value: getDatapointDescriptionRes{ErrorCode: 1}, other: getDatapointDescriptionRes{ErrorCode: 2}, resultExpected: false},
		{name: "DescriptionsLenUnequal", value: getDatapointDescriptionRes{Descriptions: []DatapointDescription{{ID: 1}}}, other: getDatapointDescriptionRes{Descriptions: []DatapointDescription{}}, resultExpected: false},
		{name: "DescriptionDPTUnequal", value: getDatapointDescriptionRes{Descriptions: []DatapointDescription{{ID: 1, Code: 1}}}, other: getDatapointDescriptionRes{Descriptions: []DatapointDescription{{ID: 1, Code: 2}}}, resultExpected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.resultExpected, tc.value.Equal(&tc.other))
		})
	}
}

func TestGetDatapointDescriptionResMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value getDatapointDescriptionRes
	}
	testCases := []testCase{
		{name: "Empty", value: getDatapointDescriptionRes{}},
		{name: "WithDescriptions", value: getDatapointDescriptionRes{StartDatapoint: 2, Descriptions: []DatapointDescription{{ID: 10, ValueType: 7, ConfigFlags: 0x02, Code: 1}, {ID: 11, ValueType: 8, ConfigFlags: 0x04, Code: 2}}}},
		{name: "Error", value: getDatapointDescriptionRes{StartDatapoint: 3, ErrorCode: 5}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			require.NoError(t, err)

			other := getDatapointDescriptionRes{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.True(t, tc.value.Equal(&other))
		})
	}
}

func TestGetDatapointDescriptionResUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name          string
		data          []byte
		errorExpected bool
	}
	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0x83, 0x00}, errorExpected: true},
		{name: "WrongService", data: []byte{0x00, 0x83, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
		{name: "TruncatedDescriptions", data: []byte{0xF0, 0x83, 0x00, 0x00, 0x00, 0x01}, errorExpected: true},
		{name: "DescriptionDataTooShort", data: []byte{0xF0, 0x83, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x07, 0x02}, errorExpected: true},
		{name: "ExtraBytes", data: []byte{0xF0, 0x83, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x07, 0x02, 0x01, 0xAB}, errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := getDatapointDescriptionRes{}
			err := r.UnmarshalBinary(tc.data)
			if tc.errorExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func FuzzGetDatapointDescriptionReq(f *testing.F) {
	req := getDatapointDescriptionReq{StartDatapoint: 1, NumberOfDatapoints: 2}
	data, _ := req.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := getDatapointDescriptionReq{}
		_ = other.UnmarshalBinary(data)
	})
}

func FuzzGetDatapointDescriptionRes(f *testing.F) {
	res := getDatapointDescriptionRes{StartDatapoint: 1, Descriptions: []DatapointDescription{{ID: 1, ValueType: 7, ConfigFlags: 0x02, Code: 1}}}
	data, _ := res.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := getDatapointDescriptionRes{}
		_ = other.UnmarshalBinary(data)
	})
}
