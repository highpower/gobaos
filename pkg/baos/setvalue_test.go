package baos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetDatapointValueReqEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          setDatapointValueReq
		other          setDatapointValueReq
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal",
			value: setDatapointValueReq{commands: []DatapointCommand{{ID: 1, Command: 1, Data: []byte{0x01}}}},
			other: setDatapointValueReq{commands: []DatapointCommand{{ID: 1, Command: 1, Data: []byte{0x01}}}}, resultExpected: true},
		{name: "IDUnequal",
			value: setDatapointValueReq{commands: []DatapointCommand{{ID: 1}}},
			other: setDatapointValueReq{commands: []DatapointCommand{{ID: 2}}}, resultExpected: false},
		{name: "CommandUnequal",
			value: setDatapointValueReq{commands: []DatapointCommand{{ID: 1, Command: 1}}},
			other: setDatapointValueReq{commands: []DatapointCommand{{ID: 1, Command: 2}}}, resultExpected: false},
		{name: "DataUnequal",
			value: setDatapointValueReq{commands: []DatapointCommand{{ID: 1, Data: []byte{0x01}}}},
			other: setDatapointValueReq{commands: []DatapointCommand{{ID: 1, Data: []byte{0x02}}}}, resultExpected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.resultExpected, tc.value.equal(&tc.other))
		})
	}
}

func TestSetDatapointValueReqMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value setDatapointValueReq
	}
	testCases := []testCase{
		{name: "SingleValue",
			value: setDatapointValueReq{commands: []DatapointCommand{{ID: 1, Command: 1, Data: []byte{0x01}}}}},
		{name: "MultipleValues",
			value: setDatapointValueReq{commands: []DatapointCommand{
				{ID: 10, Command: 2, Data: []byte{0x01, 0x02}},
				{ID: 11, Command: 3, Data: []byte{0x03}},
			}}},
		{name: "EmptyData",
			value: setDatapointValueReq{commands: []DatapointCommand{{ID: 5, Command: 5, Data: []byte{}}}}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			require.NoError(t, err)

			other := setDatapointValueReq{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.True(t, tc.value.equal(&other))
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
			r := setDatapointValueReq{}
			err := r.UnmarshalBinary(tc.data)
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
		value          setDatapointValueRes
		other          setDatapointValueRes
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal", value: setDatapointValueRes{startDatapoint: 1, numberOfDatapoints: 0, errorCode: 2},
			other: setDatapointValueRes{startDatapoint: 1, numberOfDatapoints: 0, errorCode: 2}, resultExpected: true},
		{name: "StartUnequal", value: setDatapointValueRes{startDatapoint: 1},
			other: setDatapointValueRes{startDatapoint: 2}, resultExpected: false},
		{name: "ErrorCodeUnequal", value: setDatapointValueRes{errorCode: 1},
			other: setDatapointValueRes{errorCode: 2}, resultExpected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.resultExpected, tc.value.equal(&tc.other))
		})
	}
}

func TestSetDatapointValueResMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value setDatapointValueRes
	}
	testCases := []testCase{
		{name: "Empty", value: setDatapointValueRes{}},
		{name: "WithError", value: setDatapointValueRes{startDatapoint: 3, errorCode: 5}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			require.NoError(t, err)

			other := setDatapointValueRes{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.True(t, tc.value.equal(&other))
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
			r := setDatapointValueRes{}
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
	req := setDatapointValueReq{commands: []DatapointCommand{{ID: 1, Command: SetNewValue, Data: []byte{0x01}}}}
	data, _ := req.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := setDatapointValueReq{}
		_ = other.UnmarshalBinary(data)
	})
}

func FuzzSetDatapointValueRes(f *testing.F) {
	res := setDatapointValueRes{startDatapoint: 1, numberOfDatapoints: 0, errorCode: 0}
	data, _ := res.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := setDatapointValueRes{}
		_ = other.UnmarshalBinary(data)
	})
}
