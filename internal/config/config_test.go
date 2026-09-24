package config

import (
	"slices"

	"gobaos/pkg/baos"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	c, err := Parse("testdata/config.yaml")
	require.NoError(t, err)
	assert.Equal(t, c.Device, "/dev/ttyAMA0")
	assert.Len(t, c.Sensors, 1)
	assert.Len(t, c.Controllers, 1)
}

func TestVisitIDStrings(t *testing.T) {

	type testCase struct {
		name           string
		ids            []string
		resultExpected []baos.DatapointID
		errorExpected  bool
	}

	testCases := []testCase{
		{name: "Empty",
			ids:            []string{},
			resultExpected: []baos.DatapointID{},
			errorExpected:  false},
		{name: "NonEmpty",
			ids:            []string{"1", "2", "3", "45"},
			resultExpected: []baos.DatapointID{1, 2, 3, 45},
			errorExpected:  false},
		{name: "Ranges",
			ids:            []string{"1", "2", "3", "5-9"},
			resultExpected: []baos.DatapointID{1, 2, 3, 5, 6, 7, 8, 9},
			errorExpected:  false},
		{name: "BadID",
			ids:            []string{"1", "abc", "3", "5-9"},
			resultExpected: []baos.DatapointID{},
			errorExpected:  true},
		{name: "BadRange",
			ids:            []string{"1", "2", "3", "06-08-2026"},
			resultExpected: []baos.DatapointID{},
			errorExpected:  true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := parseIDStrings(tc.ids)
			if tc.errorExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			slices.Sort(res)
			slices.Sort(tc.resultExpected)
			assert.Equal(t, tc.resultExpected, res)
		})
	}
}
