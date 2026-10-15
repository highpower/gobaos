package baos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDatapointValueReqValid(t *testing.T) {
	type testCase struct {
		name          string
		value         GetDatapointValueReq
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Empty",
			value: GetDatapointValueReq{0, 0, 0}, errorExpected: false},
		{name: "Nonempty",
			value: GetDatapointValueReq{1, 5, AllDatapoints}, errorExpected: false},
		{name: "BadFilter",
			value: GetDatapointValueReq{1, 5, 6}, errorExpected: true},
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

func TestGetDatapointValueReqEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          GetDatapointValueReq
		other          GetDatapointValueReq
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal",
			value: GetDatapointValueReq{1, 5, 0},
			other: GetDatapointValueReq{1, 5, 0}, resultExpected: true},
		{name: "StartUnequal",
			value: GetDatapointValueReq{1, 5, 0},
			other: GetDatapointValueReq{2, 5, 0}, resultExpected: false},
		{name: "NumberUnequal",
			value: GetDatapointValueReq{1, 5, 0},
			other: GetDatapointValueReq{1, 6, 0}, resultExpected: false},
		{name: "FilterUnequal",
			value: GetDatapointValueReq{1, 5, 0},
			other: GetDatapointValueReq{1, 5, 1}, resultExpected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.resultExpected, tc.value.Equal(&tc.other))
		})
	}
}

func TestGetDatapointValueReqMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value GetDatapointValueReq
	}

	testCases := []testCase{
		{name: "Empty", value: GetDatapointValueReq{0, 0, 0}},
		{name: "Nonempty", value: GetDatapointValueReq{2, 4, ValidDatapoints}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			require.NoError(t, err)

			other := GetDatapointValueReq{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.True(t, tc.value.Equal(&other))
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
			r := GetDatapointValueReq{}
			err := r.UnmarshalBinary(tc.data)
			if tc.errorExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestGetDatapointValueResValid(t *testing.T) {
	type testCase struct {
		name          string
		value         GetDatapointValueRes
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Empty", value: GetDatapointValueRes{}, errorExpected: true},
		{name: "WithValues", value: GetDatapointValueRes{
			StartDatapoint: 1, Values: []Datapoint{{ID: 1, State: 0x10, Data: []byte{0x01}}}}, errorExpected: false},
		{name: "WithErrorCode", value: GetDatapointValueRes{
			StartDatapoint: 1, ErrorCode: 5}, errorExpected: false},
		{name: "ValuesAndErrorCode", value: GetDatapointValueRes{
			StartDatapoint: 1, Values: []Datapoint{{ID: 1, Data: []byte{0x01}}}, ErrorCode: 5}, errorExpected: true},
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

func TestGetDatapointValueResEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          GetDatapointValueRes
		other          GetDatapointValueRes
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal",
			value: GetDatapointValueRes{StartDatapoint: 1, Values: []Datapoint{{ID: 1, State: 0x10, Data: []byte{0x01}}}},
			other: GetDatapointValueRes{StartDatapoint: 1, Values: []Datapoint{{ID: 1, State: 0x10, Data: []byte{0x01}}}}, resultExpected: true},
		{name: "StartUnequal",
			value: GetDatapointValueRes{StartDatapoint: 1},
			other: GetDatapointValueRes{StartDatapoint: 2}, resultExpected: false},
		{name: "ErrorCodeUnequal",
			value: GetDatapointValueRes{ErrorCode: 1},
			other: GetDatapointValueRes{ErrorCode: 2}, resultExpected: false},
		{name: "ValuesLenUnequal",
			value: GetDatapointValueRes{Values: []Datapoint{{ID: 1}}},
			other: GetDatapointValueRes{Values: []Datapoint{}}, resultExpected: false},
		{name: "ValueDataUnequal",
			value: GetDatapointValueRes{Values: []Datapoint{{ID: 1, Data: []byte{0x01}}}},
			other: GetDatapointValueRes{Values: []Datapoint{{ID: 1, Data: []byte{0x02}}}}, resultExpected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.resultExpected, tc.value.Equal(&tc.other))
		})
	}
}

func TestGetDatapointValueResMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value GetDatapointValueRes
	}

	testCases := []testCase{
		{name: "Empty", value: GetDatapointValueRes{}},
		{name: "WithValues", value: GetDatapointValueRes{StartDatapoint: 2, Values: []Datapoint{
			{ID: 10, State: 0x10, Data: []byte{0x01}}, {ID: 11, State: 0x20, Data: []byte{0x03, 0x04}}}}},
		{name: "WithSubByteValue", value: GetDatapointValueRes{StartDatapoint: 1, Values: []Datapoint{
			{ID: 1, State: 0x10, Data: []byte{0x0A}}}}},
		{name: "Error", value: GetDatapointValueRes{StartDatapoint: 3, ErrorCode: 5}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			require.NoError(t, err)

			other := GetDatapointValueRes{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.True(t, tc.value.Equal(&other))
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
			r := GetDatapointValueRes{}
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
	req := GetDatapointValueReq{StartDatapoint: 1, NumberOfDatapoints: 2, Filter: AllDatapoints}
	data, _ := req.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := GetDatapointValueReq{}
		_ = other.UnmarshalBinary(data)
	})
}

func FuzzGetDatapointValueRes(f *testing.F) {
	res := GetDatapointValueRes{StartDatapoint: 1, Values: []Datapoint{{ID: 1, State: 0x10, Data: []byte{0x01}}}}
	data, _ := res.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := GetDatapointValueRes{}
		_ = other.UnmarshalBinary(data)
	})
}
