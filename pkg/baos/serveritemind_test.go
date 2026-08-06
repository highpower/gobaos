package baos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerItemIndValid(t *testing.T) {
	type testCase struct {
		name          string
		value         ServerItemInd
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Empty", value: ServerItemInd{}, errorExpected: true},
		{name: "WithItems", value: ServerItemInd{StartItem: 1, Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, errorExpected: false},
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

func TestServerItemIndEqual(t *testing.T) {
	type testCase struct {
		name           string
		value          ServerItemInd
		other          ServerItemInd
		resultExpected bool
	}

	testCases := []testCase{
		{name: "Equal", value: ServerItemInd{StartItem: 1, Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, other: ServerItemInd{StartItem: 1, Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, resultExpected: true},
		{name: "StartUnequal", value: ServerItemInd{StartItem: 1, Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, other: ServerItemInd{StartItem: 2, Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, resultExpected: false},
		{name: "ItemsLenUnequal", value: ServerItemInd{StartItem: 1, Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, other: ServerItemInd{StartItem: 1, Items: []ServerItem{}}, resultExpected: false},
		{name: "ItemDataUnequal", value: ServerItemInd{StartItem: 1, Items: []ServerItem{{ID: 1, Data: []byte{0x01}}}}, other: ServerItemInd{StartItem: 1, Items: []ServerItem{{ID: 1, Data: []byte{0x02}}}}, resultExpected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.resultExpected, tc.value.Equal(&tc.other))
		})
	}
}

func TestServerItemIndMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value ServerItemInd
	}
	testCases := []testCase{
		{name: "WithItems", value: ServerItemInd{StartItem: 2, Items: []ServerItem{{ID: 10, Data: []byte{0x01, 0x02}}, {ID: 11, Data: []byte{0x03}}}}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			require.NoError(t, err)

			other := ServerItemInd{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.True(t, tc.value.Equal(&other))
		})
	}
}

func TestServerItemIndUnmarshalErrors(t *testing.T) {
	type testCase struct {
		name        string
		data        []byte
		errExpected bool
	}
	testCases := []testCase{
		{name: "TooShort", data: []byte{0xF0, 0xC2, 0x00}, errExpected: true},
		{name: "WrongService", data: []byte{0x00, 0xC2, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "WrongSubService", data: []byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00}, errExpected: true},
		{name: "TruncatedItems", data: []byte{0xF0, 0xC2, 0x00, 0x00, 0x00, 0x01}, errExpected: true},
		{name: "ItemDataTooShort", data: []byte{0xF0, 0xC2, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x02}, errExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := ServerItemInd{}
			err := r.UnmarshalBinary(tc.data)
			if tc.errExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func FuzzServerItemInd(f *testing.F) {
	ind := ServerItemInd{StartItem: 1, Items: []ServerItem{{ID: 1, Data: []byte{0x01, 0x02}}}}
	data, _ := ind.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := ServerItemInd{}
		_ = other.UnmarshalBinary(data)
	})
}
