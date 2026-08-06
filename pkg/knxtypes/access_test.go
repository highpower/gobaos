package knxtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDPTAccessMarshal(t *testing.T) {
	type testCase struct {
		name           string
		key            uint8
		id             uint8
		index          uint8
		additionalData uint8
	}

	testCases := []testCase{
		{name: "Simple", key: 10, id: 20, index: 5, additionalData: 1},
		{name: "Zero", key: 0, id: 0, index: 0, additionalData: 0},
		{name: "Max", key: 127, id: 127, index: 15, additionalData: 3},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			acc := DPTAccess{Key: tc.key, ID: tc.id, Index: tc.index, AdditionalData: tc.additionalData}
			data, err := acc.MarshalBinary()
			require.NoError(t, err)

			other := DPTAccess{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.Equal(t, acc, other)
		})
	}
}

func TestDPTAccessBadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTAccess](t)
}

func TestDPTAccessSet(t *testing.T) {
	type testCase struct {
		name    string
		value   string
		result  DPTAccess
		wantErr bool
	}

	testCases := []testCase{
		{name: "Valid", value: "10:20:5:1", result: DPTAccess{Key: 10, ID: 20, Index: 5, AdditionalData: 1}},
		{name: "Min", value: "0:0:0:0", result: DPTAccess{Key: 0, ID: 0, Index: 0, AdditionalData: 0}},
		{name: "Max", value: "127:127:15:3", result: DPTAccess{Key: 127, ID: 127, Index: 15, AdditionalData: 3}},
		{name: "KeyTooLarge", value: "128:0:0:0", wantErr: true},
		{name: "IDTooLarge", value: "0:128:0:0", wantErr: true},
		{name: "IndexTooLarge", value: "0:0:16:0", wantErr: true},
		{name: "AddDataTooLarge", value: "0:0:0:4", wantErr: true},
		{name: "TooFewParts", value: "1:2:3", wantErr: true},
		{name: "TooManyParts", value: "1:2:3:4:5", wantErr: true},
		{name: "InvalidKey", value: "abc:0:0:0", wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res := DPTAccess{}
			err := res.Set(tc.value)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.result, res)
		})
	}
}

func FuzzDPTAccess(f *testing.F) {
	a := DPTAccess{Key: 10, ID: 20, Index: 5, AdditionalData: 1}
	data, _ := a.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTAccess{}
		_ = other.UnmarshalBinary(data)
	})
}
