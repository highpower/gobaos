package baos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDatapointValueReqEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          getDatapointValueReq
		other          getDatapointValueReq
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal",
			value: getDatapointValueReq{1, 5, 0},
			other: getDatapointValueReq{1, 5, 0}, resultExpected: true},
		{name: "StartUnequal",
			value: getDatapointValueReq{1, 5, 0},
			other: getDatapointValueReq{2, 5, 0}, resultExpected: false},
		{name: "NumberUnequal",
			value: getDatapointValueReq{1, 5, 0},
			other: getDatapointValueReq{1, 6, 0}, resultExpected: false},
		{name: "FilterUnequal",
			value: getDatapointValueReq{1, 5, 0},
			other: getDatapointValueReq{1, 5, 1}, resultExpected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.resultExpected, tc.value.equal(&tc.other))
		})
	}
}

func TestGetDatapointValueReqMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value getDatapointValueReq
	}

	testCases := []testCase{
		{name: "Empty", value: getDatapointValueReq{0, 0, 0}},
		{name: "Nonempty", value: getDatapointValueReq{2, 4, ValidDatapoints}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			require.NoError(t, err)

			other := getDatapointValueReq{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.True(t, tc.value.equal(&other))
		})
	}
}

func TestGetDatapointValueReqUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name          string
		data          []byte
		errorExpected bool
	}

	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0x05, 0x00}, errorExpected: true},
		{name: "WrongService", data: []byte{0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := getDatapointValueReq{}
			err := r.UnmarshalBinary(tc.data)
			if tc.errorExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestGetDatapointValueResEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          getDatapointValueRes
		other          getDatapointValueRes
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal",
			value: getDatapointValueRes{startDatapoint: 1, values: []Datapoint{{ID: 1, State: 0x10, Data: []byte{0x01}}}},
			other: getDatapointValueRes{startDatapoint: 1, values: []Datapoint{{ID: 1, State: 0x10, Data: []byte{0x01}}}}, resultExpected: true},
		{name: "StartUnequal",
			value: getDatapointValueRes{startDatapoint: 1},
			other: getDatapointValueRes{startDatapoint: 2}, resultExpected: false},
		{name: "ErrorCodeUnequal",
			value: getDatapointValueRes{errorCode: 1},
			other: getDatapointValueRes{errorCode: 2}, resultExpected: false},
		{name: "ValuesLenUnequal",
			value: getDatapointValueRes{values: []Datapoint{{ID: 1}}},
			other: getDatapointValueRes{values: []Datapoint{}}, resultExpected: false},
		{name: "ValueDataUnequal",
			value: getDatapointValueRes{values: []Datapoint{{ID: 1, Data: []byte{0x01}}}},
			other: getDatapointValueRes{values: []Datapoint{{ID: 1, Data: []byte{0x02}}}}, resultExpected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.resultExpected, tc.value.equal(&tc.other))
		})
	}
}

func TestGetDatapointValueResMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value getDatapointValueRes
	}

	testCases := []testCase{
		{name: "Empty", value: getDatapointValueRes{}},
		{name: "WithValues", value: getDatapointValueRes{startDatapoint: 2, values: []Datapoint{
			{ID: 10, State: 0x10, Data: []byte{0x01}}, {ID: 11, State: 0x20, Data: []byte{0x03, 0x04}}}}},
		{name: "WithSubByteValue", value: getDatapointValueRes{startDatapoint: 1, values: []Datapoint{
			{ID: 1, State: 0x10, Data: []byte{0x0A}}}}},
		{name: "Error", value: getDatapointValueRes{startDatapoint: 3, errorCode: 5}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			require.NoError(t, err)

			other := getDatapointValueRes{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.True(t, tc.value.equal(&other))
		})
	}
}

func TestGetDatapointValueResUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name          string
		data          []byte
		errorExpected bool
	}

	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0x85, 0x00}, errorExpected: true},
		{name: "WrongService", data: []byte{0x00, 0x85, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
		{name: "TruncatedValues", data: []byte{0xF0, 0x85, 0x00, 0x00, 0x00, 0x01}, errorExpected: true},
		{name: "InvalidLength", data: []byte{0xF0, 0x85, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x10, 0xFF}, errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := getDatapointValueRes{}
			err := r.UnmarshalBinary(tc.data)
			if tc.errorExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func FuzzGetDatapointValueReq(f *testing.F) {
	req := getDatapointValueReq{startDatapoint: 1, numberOfDatapoints: 2, filter: AllDatapoints}
	data, _ := req.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := getDatapointValueReq{}
		_ = other.UnmarshalBinary(data)
	})
}

func FuzzGetDatapointValueRes(f *testing.F) {
	res := getDatapointValueRes{startDatapoint: 1, values: []Datapoint{{ID: 1, State: 0x10, Data: []byte{0x01}}}}
	data, _ := res.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := getDatapointValueRes{}
		_ = other.UnmarshalBinary(data)
	})
}
