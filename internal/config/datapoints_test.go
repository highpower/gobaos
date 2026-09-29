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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := parseIDs(tc.s)
			if tc.errorExpected {
				require.Error(t, err)
			}
			assert.Equal(t, res, tc.resultExpected)
		})
	}

}
