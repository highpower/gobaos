package baos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetDatapointValueReqValid(t *testing.T) {
	type testCase struct {
		name          string
		value         SetDatapointValueReq
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Empty",
			value: SetDatapointValueReq{}, errorExpected: true},
		{name: "WithValues",
			value: SetDatapointValueReq{Commands: []DatapointCommand{{ID: 1, Command: 1, Data: []byte{0x01}}}}, errorExpected: false},
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

func TestSetDatapointValueReqEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          SetDatapointValueReq
		other          SetDatapointValueReq
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal",
			value: SetDatapointValueReq{Commands: []DatapointCommand{{ID: 1, Command: 1, Data: []byte{0x01}}}},
			other: SetDatapointValueReq{Commands: []DatapointCommand{{ID: 1, Command: 1, Data: []byte{0x01}}}}, resultExpected: true},
		{name: "IDUnequal",
			value: SetDatapointValueReq{Commands: []DatapointCommand{{ID: 1}}},
			other: SetDatapointValueReq{Commands: []DatapointCommand{{ID: 2}}}, resultExpected: false},
		{name: "CommandUnequal",
			value: SetDatapointValueReq{Commands: []DatapointCommand{{ID: 1, Command: 1}}},
			other: SetDatapointValueReq{Commands: []DatapointCommand{{ID: 1, Command: 2}}}, resultExpected: false},
		{name: "DataUnequal",
			value: SetDatapointValueReq{Commands: []DatapointCommand{{ID: 1, Data: []byte{0x01}}}},
			other: SetDatapointValueReq{Commands: []DatapointCommand{{ID: 1, Data: []byte{0x02}}}}, resultExpected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.resultExpected, tc.value.Equal(&tc.other))
		})
	}
}

func TestSetDatapointValueReqMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value SetDatapointValueReq
	}
	testCases := []testCase{
		{name: "SingleValue",
			value: SetDatapointValueReq{Commands: []DatapointCommand{{ID: 1, Command: 1, Data: []byte{0x01}}}}},
		{name: "MultipleValues",
			value: SetDatapointValueReq{Commands: []DatapointCommand{
				{ID: 10, Command: 2, Data: []byte{0x01, 0x02}},
				{ID: 11, Command: 3, Data: []byte{0x03}},
			}}},
		{name: "EmptyData",
			value: SetDatapointValueReq{Commands: []DatapointCommand{{ID: 5, Command: 5, Data: []byte{}}}}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			require.NoError(t, err)

			other := SetDatapointValueReq{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.True(t, tc.value.Equal(&other))
		})
	}
}

func TestSetDatapointValueReqUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name          string
		data          []byte
		errorExpected bool
	}
	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0x06, 0x00}, errorExpected: true},
		{name: "WrongService", data: []byte{0x00, 0x06, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
		{name: "TruncatedValues", data: []byte{0xF0, 0x06, 0x00, 0x00, 0x00, 0x01}, errorExpected: true},
		{name: "ValueDataTooShort", data: []byte{0xF0, 0x06, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x02}, errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := SetDatapointValueReq{}
			err := r.UnmarshalBinary(tc.data)
			if tc.errorExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestSetDatapointValueResValid(t *testing.T) {
	type testCase struct {
		name          string
		value         SetDatapointValueRes
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Empty", value: SetDatapointValueRes{}, errorExpected: false},
		{name: "WithErrorCode", value: SetDatapointValueRes{StartDatapoint: 1, ErrorCode: 5}, errorExpected: false},
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

func TestSetDatapointValueResEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          SetDatapointValueRes
		other          SetDatapointValueRes
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal", value: SetDatapointValueRes{StartDatapoint: 1, NumberOfDatapoints: 0, ErrorCode: 2},
			other: SetDatapointValueRes{StartDatapoint: 1, NumberOfDatapoints: 0, ErrorCode: 2}, resultExpected: true},
		{name: "StartUnequal", value: SetDatapointValueRes{StartDatapoint: 1},
			other: SetDatapointValueRes{StartDatapoint: 2}, resultExpected: false},
		{name: "ErrorCodeUnequal", value: SetDatapointValueRes{ErrorCode: 1},
			other: SetDatapointValueRes{ErrorCode: 2}, resultExpected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.resultExpected, tc.value.Equal(&tc.other))
		})
	}
}

func TestSetDatapointValueResMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value SetDatapointValueRes
	}
	testCases := []testCase{
		{name: "Empty", value: SetDatapointValueRes{}},
		{name: "WithError", value: SetDatapointValueRes{StartDatapoint: 3, ErrorCode: 5}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			require.NoError(t, err)

			other := SetDatapointValueRes{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.True(t, tc.value.Equal(&other))
		})
	}
}

func TestSetDatapointValueResUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name          string
		data          []byte
		errorExpected bool
	}
	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0x86, 0x00}, errorExpected: true},
		{name: "WrongService", data: []byte{0x00, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := SetDatapointValueRes{}
			err := r.UnmarshalBinary(tc.data)
			if tc.errorExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func FuzzSetDatapointValueReq(f *testing.F) {
	req := SetDatapointValueReq{Commands: []DatapointCommand{{ID: 1, Command: SetNewValue, Data: []byte{0x01}}}}
	data, _ := req.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := SetDatapointValueReq{}
		_ = other.UnmarshalBinary(data)
	})
}

func FuzzSetDatapointValueRes(f *testing.F) {
	res := SetDatapointValueRes{StartDatapoint: 1, NumberOfDatapoints: 0, ErrorCode: 0}
	data, _ := res.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := SetDatapointValueRes{}
		_ = other.UnmarshalBinary(data)
	})
}
