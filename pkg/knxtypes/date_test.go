package knxtypes

import "testing"

func TestDPTDateMarshal(t *testing.T) {
	type testCase struct {
		name  string
		day   uint8
		month uint8
		year  uint8
	}

	testCases := []testCase{
		{name: "Min", day: 1, month: 1, year: 0},
		{name: "Normal", day: 15, month: 6, year: 23},
		{name: "Max", day: 31, month: 12, year: 99},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTDate{Day: tc.day, Month: tc.month, Year: tc.year}
			data, err := dpt.MarshalBinary()
			if err != nil {
				t.Fatalf("%s", err.Error())
			}
			other := DPTDate{}
			if err = other.UnmarshalBinary(data); err != nil {
				t.Fatalf("%s", err.Error())
			}
			if other.Day != dpt.Day || other.Month != dpt.Month || other.Year != dpt.Year {
				t.Errorf("got day=%d month=%d year=%d, want day=%d month=%d year=%d",
					other.Day, other.Month, other.Year,
					dpt.Day, dpt.Month, dpt.Year)
			}
		})
	}
}

func TestDPTDateUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTDate](t)
}

func TestDPTDateSet(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTDate
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Valid", input: "15/06/23", expected: DPTDate{Day: 15, Month: 6, Year: 23}, errorExpected: false},
		{name: "WrongParts", input: "15/06/23/extra", errorExpected: true},
		{name: "NotADay", input: "bad/06/23", errorExpected: true},
		{name: "NotAMonth", input: "15/bad/23", errorExpected: true},
		{name: "NotAYear", input: "15/06/bad", errorExpected: true},
		{name: "BadDay", input: "32/06/23", errorExpected: true},
		{name: "BadMonth", input: "15/13/23", errorExpected: true},
		{name: "BadYear", input: "15/06/106", errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTDate{}
			err := dpt.Set(tc.input)
			switch {
			case err == nil && tc.errorExpected:
				t.Fatalf("error expected but did not happen")
			case err != nil && !tc.errorExpected:
				t.Fatalf("%s", err.Error())
			case err != nil:
				return
			}
			if dpt != tc.expected {
				t.Errorf("got dpt=%s, expected %s", dpt.String(), tc.expected.String())
			}
		})
	}
}

func TestDPTDateValid(t *testing.T) {
	type testCase struct {
		name          string
		dpt           DPTDate
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Min", dpt: DPTDate{Day: 1, Month: 1, Year: 0}, errorExpected: false},
		{name: "Max", dpt: DPTDate{Day: 31, Month: 12, Year: 99}, errorExpected: false},
		{name: "DayZero", dpt: DPTDate{Day: 0, Month: 1, Year: 0}, errorExpected: true},
		{name: "BadDay", dpt: DPTDate{Day: 32, Month: 1, Year: 0}, errorExpected: true},
		{name: "BadMonth", dpt: DPTDate{Day: 1, Month: 0, Year: 0}, errorExpected: true},
		{name: "BadMonthMax", dpt: DPTDate{Day: 1, Month: 13, Year: 0}, errorExpected: true},
		{name: "BadYear", dpt: DPTDate{Day: 1, Month: 1, Year: 100}, errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.dpt.Valid()
			if tc.errorExpected && err == nil {
				t.Fatalf("error expected but did not happen")
			}
			if !tc.errorExpected && err != nil {
				t.Fatalf("%s", err.Error())
			}
		})
	}
}

func FuzzDPTDate(f *testing.F) {
	dpt := DPTDate{Day: 15, Month: 6, Year: 23}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTDate{}
		_ = other.UnmarshalBinary(data)
	})
}
