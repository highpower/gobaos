package baos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatapointValueIndEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          datapointValueInd
		other          datapointValueInd
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal",
			value: datapointValueInd{StartDatapoint: 1, Values: []Datapoint{{ID: 1, State: 0x10, Data: []byte{0x01}}}},
			other: datapointValueInd{StartDatapoint: 1, Values: []Datapoint{{ID: 1, State: 0x10, Data: []byte{0x01}}}}, resultExpected: true},
		{name: "StartUnequal",
			value: datapointValueInd{StartDatapoint: 1}, other: datapointValueInd{StartDatapoint: 2}, resultExpected: false},
		{name: "ValuesLenUnequal",
			value: datapointValueInd{Values: []Datapoint{{ID: 1}}},
			other: datapointValueInd{Values: []Datapoint{{ID: 2}}}, resultExpected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.resultExpected, tc.value.Equal(&tc.other))
		})
	}
}

func TestDatapointValueIndMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value datapointValueInd
	}

	testCases := []testCase{
		{name: "Empty", value: datapointValueInd{}},
		{name: "WithValues", value: datapointValueInd{StartDatapoint: 2, Values: []Datapoint{
			{ID: 10, State: 0x10, Data: []byte{0x01}}, {ID: 11, State: 0x20, Data: []byte{0x03, 0x04}}}}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			require.NoError(t, err)

			other := datapointValueInd{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.True(t, tc.value.Equal(&other))
		})
	}
}

func TestDatapointValueIndUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name          string
		data          []byte
		errorExpected bool
	}

	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0xC1, 0x00}, errorExpected: true},
		{name: "WrongService", data: []byte{0x00, 0xC1, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
		{name: "TruncatedValues", data: []byte{0xF0, 0xC1, 0x00, 0x00, 0x00, 0x01}, errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := datapointValueInd{}
			err := r.UnmarshalBinary(tc.data)
			if tc.errorExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func FuzzDatapointValueInd(f *testing.F) {
	ind := datapointValueInd{StartDatapoint: 1, Values: []Datapoint{{ID: 1, State: 0x10, Data: []byte{0x01}}}}
	data, _ := ind.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := datapointValueInd{}
		_ = other.UnmarshalBinary(data)
	})
}
