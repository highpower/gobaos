package baos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerItemIndEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          serverItemInd
		other          serverItemInd
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal", value: serverItemInd{startItem: 1, items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, other: serverItemInd{startItem: 1, items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, resultExpected: true},
		{name: "StartUnequal", value: serverItemInd{startItem: 1, items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, other: serverItemInd{startItem: 2, items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, resultExpected: false},
		{name: "ItemsLenUnequal", value: serverItemInd{startItem: 1, items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, other: serverItemInd{startItem: 1, items: []ServerItem{}}, resultExpected: false},
		{name: "ItemDataUnequal", value: serverItemInd{startItem: 1, items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, other: serverItemInd{startItem: 1, items: []ServerItem{{ID: 1, Data: []byte{0x02}}}}, resultExpected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.resultExpected, tc.value.equal(&tc.other))
		})
	}
}

func TestServerItemIndMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value serverItemInd
	}
	testCases := []testCase{
		{name: "WithItems", value: serverItemInd{startItem: 2, items: []ServerItem{{ID: 10, Data: []byte{0x01, 0x02}}, {ID: 11, Data: []byte{0x03}}}}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			require.NoError(t, err)

			other := serverItemInd{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.True(t, tc.value.equal(&other))
		})
	}
}

func TestServerItemIndUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name          string
		data          []byte
		errorExpected bool
	}
	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0xC2, 0x00}, errorExpected: true},
		{name: "WrongService", data: []byte{0x00, 0xC2, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00}, errorExpected: true},
		{name: "TruncatedItems", data: []byte{0xF0, 0xC2, 0x00, 0x00, 0x00, 0x01}, errorExpected: true},
		{name: "ItemDataTooShort", data: []byte{0xF0, 0xC2, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x02}, errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := serverItemInd{}
			err := r.UnmarshalBinary(tc.data)
			if tc.errorExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func FuzzServerItemInd(f *testing.F) {
	ind := serverItemInd{startItem: 1, items: []ServerItem{{ID: 1, Data: []byte{0x01, 0x02}}}}
	data, _ := ind.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := serverItemInd{}
		_ = other.UnmarshalBinary(data)
	})
}
