package mask

import "testing"

func TestApply(t *testing.T) {
	type testCase struct {
		name           string
		value          string
		resultExpected string
	}

	//goland:noinspection SpellCheckingInspection
	testCases := []testCase{
		{name: "Empty", value: "", resultExpected: ""},
		{name: "Short", value: "abcd", resultExpected: "a..."},
		{name: "Long", value: "abcdefgh", resultExpected: "abcd..."},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if res := Apply(tc.value); res != tc.resultExpected {
				t.Errorf("incorrect result: expected %q, got %q", tc.resultExpected, res)
			}
		})
	}

}

func FuzzApply(f *testing.F) {
	f.Add("test")
	f.Fuzz(func(t *testing.T, value string) {
		Apply(value)
	})
}
