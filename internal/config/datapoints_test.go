package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseIDs(t *testing.T) {

	type testCase struct {
		name           string
		s              string
		resultExpected []int
		errorExpected  bool
	}

	testCases := []testCase{
		{name: "Empty", s: "", resultExpected: nil, errorExpected: false},
		{name: "List", s: "1,2,3,6", resultExpected: []int{1, 2, 3, 6}, errorExpected: false},
		{name: "Ranges", s: "1-3,6-9", resultExpected: []int{1, 2, 3, 6, 7, 8, 9}, errorExpected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := parseIDs(tc.s)
			if tc.errorExpected {
				require.Error(t, err)
			}
			assert.Equal(t, tc.resultExpected, res)
		})
	}

}
