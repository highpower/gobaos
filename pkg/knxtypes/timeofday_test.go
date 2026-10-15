package knxtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			require.NoError(t, err)

			other := DPTTimeOfDay{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.Equal(t, dpt.Weekday, other.Weekday)
			assert.Equal(t, dpt.Hour, other.Hour)
			assert.Equal(t, dpt.Minute, other.Minute)
			assert.Equal(t, dpt.Second, other.Second)
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
		{name: "Zero", dpt: DPTTimeOfDay{Weekday: 0, Hour: 0, Minute: 0, Second: 0}},
		{name: "MaxValues", dpt: DPTTimeOfDay{Weekday: 7, Hour: 23, Minute: 59, Second: 59}},
		{name: "BadHour", dpt: DPTTimeOfDay{Weekday: 0, Hour: 24, Minute: 0, Second: 0}, errorExpected: true},
		{name: "BadMinute", dpt: DPTTimeOfDay{Weekday: 0, Hour: 0, Minute: 60, Second: 0}, errorExpected: true},
		{name: "BadSecond", dpt: DPTTimeOfDay{Weekday: 0, Hour: 0, Minute: 0, Second: 60}, errorExpected: true},
		{name: "BadWeekday", dpt: DPTTimeOfDay{Weekday: 8, Hour: 0, Minute: 0, Second: 0}, errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.dpt.Valid()
			if tc.errorExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
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
		{name: "ValidTimeOnly", input: "12:30:45", expected: DPTTimeOfDay{Weekday: 0, Hour: 12, Minute: 30, Second: 45}},
		{name: "ValidWithWeekday", input: "mon,12:30:45", expected: DPTTimeOfDay{Weekday: 1, Hour: 12, Minute: 30, Second: 45}},
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
			if tc.errorExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, dpt)
		})
	}
}

func FuzzDPTTimeOfDay(f *testing.F) {
	dpt := DPTTimeOfDay{Weekday: 3, Hour: 12, Minute: 30, Second: 45}
	data, _ := dpt.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTTimeOfDay{}
		_ = other.UnmarshalBinary(data)
	})
}
