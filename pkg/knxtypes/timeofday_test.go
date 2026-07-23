package knxtypes

import "testing"

func TestDPTTimeOfDayMarshal(t *testing.T) {
	type testCase struct {
		name    string
		weekday uint8
		hour    uint8
		minute  uint8
		second  uint8
	}

	testCases := []testCase{
		{name: "Zero", weekday: 0, hour: 0, minute: 0, second: 0},
		{name: "Weekday", weekday: 1, hour: 12, minute: 30, second: 45},
		{name: "MaxValues", weekday: 7, hour: 23, minute: 59, second: 59},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTTimeOfDay{Weekday: tc.weekday, Hour: tc.hour, Minute: tc.minute, Second: tc.second}
			data, err := dpt.MarshalBinary()
			if err != nil {
				t.Fatalf("%s", err.Error())
			}
			other := DPTTimeOfDay{}
			if err = other.UnmarshalBinary(data); err != nil {
				t.Fatalf("%s", err.Error())
			}
			if other.Weekday != dpt.Weekday || other.Hour != dpt.Hour || other.Minute != dpt.Minute || other.Second != dpt.Second {
				t.Errorf("got weekday=%d hour=%d minute=%d second=%d, want weekday=%d hour=%d minute=%d second=%d",
					other.Weekday, other.Hour, other.Minute, other.Second,
					dpt.Weekday, dpt.Hour, dpt.Minute, dpt.Second)
			}
		})
	}
}

func TestDPTTimeOfDayValid(t *testing.T) {
	type testCase struct {
		name          string
		dpt           DPTTimeOfDay
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Zero", dpt: DPTTimeOfDay{Weekday: 0, Hour: 0, Minute: 0, Second: 0}, errorExpected: false},
		{name: "MaxValues", dpt: DPTTimeOfDay{Weekday: 7, Hour: 23, Minute: 59, Second: 59}, errorExpected: false},
		{name: "BadHour", dpt: DPTTimeOfDay{Weekday: 0, Hour: 24, Minute: 0, Second: 0}, errorExpected: true},
		{name: "BadMinute", dpt: DPTTimeOfDay{Weekday: 0, Hour: 0, Minute: 60, Second: 0}, errorExpected: true},
		{name: "BadSecond", dpt: DPTTimeOfDay{Weekday: 0, Hour: 0, Minute: 0, Second: 60}, errorExpected: true},
		{name: "BadWeekday", dpt: DPTTimeOfDay{Weekday: 8, Hour: 0, Minute: 0, Second: 0}, errorExpected: true},
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

func TestDPTTimeOfDayUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTTimeOfDay](t)
}

func TestDPTTimeOfDaySet(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		expected      DPTTimeOfDay
		errorExpected bool
	}

	testCases := []testCase{
		{name: "ValidTimeOnly", input: "12:30:45", expected: DPTTimeOfDay{Weekday: 0, Hour: 12, Minute: 30, Second: 45}, errorExpected: false},
		{name: "ValidWithWeekday", input: "mon,12:30:45", expected: DPTTimeOfDay{Weekday: 1, Hour: 12, Minute: 30, Second: 45}, errorExpected: false},
		{name: "WrongParts", input: "mon,12:30:45,extra", errorExpected: true},
		{name: "BadTimeFormat", input: "12:30", errorExpected: true},
		{name: "NotAnHour", input: "mon,bad:00:00", errorExpected: true},
		{name: "NotAMinute", input: "00:bad:00", errorExpected: true},
		{name: "NotASecond", input: "tue,00:00:bad", errorExpected: true},
		{name: "BadHour", input: "25:00:00", errorExpected: true},
		{name: "BadMinute", input: "12:61:00", errorExpected: true},
		{name: "BadSecond", input: "12:00:61", errorExpected: true},
		{name: "BadWeekday", input: "foo,12:00:00", errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTTimeOfDay{}
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
