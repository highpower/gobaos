package pattern

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	type testCase struct {
		name          string
		str           string
		result        Parsed
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Empty", str: "", result: Parsed{}, errorExpected: true},
		{name: "OneToken", str: "+", result: Parsed{index: 0, parts: []string{"+"}}, errorExpected: false},
		{name: "Common", str: "thermo/+/set", result: Parsed{index: 1, parts: []string{"thermo", "+", "set"}}, errorExpected: false},
		{name: "NoWildcard", str: "thermo/all", result: Parsed{}, errorExpected: true},
		{name: "MultipleWildcard", str: "thermo/+/set/+", result: Parsed{}, errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := Parse(tc.str)
			if tc.errorExpected {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tc.result, p)
		})
	}
}

func TestApply(t *testing.T) {

	type testCase struct {
		name           string
		pattern        Parsed
		input          string
		resultExpected string
	}

	testCases := []testCase{
		{name: "OnlyWildcard", pattern: Parsed{index: 0, parts: []string{"+"}}, input: "foo",
			resultExpected: "foo"},
		{name: "WildcardInMiddle", pattern: Parsed{index: 1, parts: []string{"thermo", "+", "set"}}, input: "kitchen",
			resultExpected: "thermo/kitchen/set"},
		{name: "WildcardAtStart", pattern: Parsed{index: 0, parts: []string{"+", "status"}}, input: "device1",
			resultExpected: "device1/status"},
		{name: "WildcardAtEnd", pattern: Parsed{index: 2, parts: []string{"home", "room", "+"}}, input: "temp",
			resultExpected: "home/room/temp"},
		{name: "EmptyInput", pattern: Parsed{index: 1, parts: []string{"a", "+", "b"}}, input: "",
			resultExpected: "a//b"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res := tc.pattern.Apply(tc.input)
			assert.Equal(t, tc.resultExpected, res)
		})
	}
}
