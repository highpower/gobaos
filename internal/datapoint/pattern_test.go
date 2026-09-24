package datapoint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePattern(t *testing.T) {
	type testCase struct {
		name          string
		str           string
		result        pattern
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Empty", str: "", result: pattern{}, errorExpected: true},
		{name: "OneToken", str: "+", result: pattern{index: 0, parts: []string{"+"}}, errorExpected: false},
		{name: "Common", str: "thermo/+/set", result: pattern{index: 1, parts: []string{"thermo", "+", "set"}}, errorExpected: false},
		{name: "NoWildcard", str: "thermo/all", result: pattern{}, errorExpected: true},
		{name: "MultipleWildcard", str: "thermo/+/set/+", result: pattern{}, errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := parsePattern(tc.str)
			if tc.errorExpected {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tc.result, p)
		})
	}
}

func TestPatternApply(t *testing.T) {

	type testCase struct {
		name           string
		pattern        pattern
		input          string
		resultExpected string
	}

	testCases := []testCase{
		{name: "OnlyWildcard", pattern: pattern{index: 0, parts: []string{"+"}}, input: "foo",
			resultExpected: "foo"},
		{name: "WildcardInMiddle", pattern: pattern{index: 1, parts: []string{"thermo", "+", "set"}}, input: "kitchen",
			resultExpected: "thermo/kitchen/set"},
		{name: "WildcardAtStart", pattern: pattern{index: 0, parts: []string{"+", "status"}}, input: "device1",
			resultExpected: "device1/status"},
		{name: "WildcardAtEnd", pattern: pattern{index: 2, parts: []string{"home", "room", "+"}}, input: "temp",
			resultExpected: "home/room/temp"},
		{name: "EmptyInput", pattern: pattern{index: 1, parts: []string{"a", "+", "b"}}, input: "",
			resultExpected: "a//b"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res := tc.pattern.apply(tc.input)
			assert.Equal(t, tc.resultExpected, res)
		})
	}
}
