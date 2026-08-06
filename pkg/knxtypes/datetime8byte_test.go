package knxtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDPTDateTime8ByteMarshal(t *testing.T) {
	type testCase struct {
		name  string
		value DPTDateTime8Byte
	}

	testCases := []testCase{
		{name: "Simple",
			value: DPTDateTime8Byte{Year: 124, Month: 3, Day: 15, Weekday: 5, Hour: 14, Minute: 30, Second: 45, Flags: 0x02}},
		{name: "Zero",
			value: DPTDateTime8Byte{Year: 0, Month: 1, Day: 1, Hour: 0, Minute: 0, Second: 0}},
		{name: "Max",
			value: DPTDateTime8Byte{Year: 255, Month: 12, Day: 31, Weekday: 7, Hour: 23, Minute: 59, Second: 59, Flags: 0xFF}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.value.MarshalBinary()
			require.NoError(t, err)

			other := DPTDateTime8Byte{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.Equal(t, tc.value, other)
		})
	}
}

func TestDPTDateTime8BadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTDateTime8Byte](t)
}

func TestDPTDateTime8ByteSet(t *testing.T) {
	type testCase struct {
		name          string
		value         string
		result        DPTDateTime8Byte
		errorExpected bool
	}

	testCases := []testCase{
		{name: "Valid", value: "15/03/2024 14:30:45",
			result: DPTDateTime8Byte{Year: 124, Month: 3, Day: 15, Hour: 14, Minute: 30, Second: 45}},
		{name: "MinYear", value: "01/01/1900 00:00:00",
			result: DPTDateTime8Byte{Year: 0, Month: 1, Day: 1, Hour: 0, Minute: 0, Second: 0}},
		{name: "MaxYear", value: "31/12/2155 23:59:59",
			result: DPTDateTime8Byte{Year: 255, Month: 12, Day: 31, Hour: 23, Minute: 59, Second: 59}},
		{name: "Weekday", value: "31/12/2155,Wed 23:59:59",
			result: DPTDateTime8Byte{Year: 255, Month: 12, Day: 31, Weekday: 3, Hour: 23, Minute: 59, Second: 59}},
		{name: "DayTooLarge", value: "32/01/2024 12:00:00", errorExpected: true},
		{name: "MonthTooLarge", value: "01/13/2024 12:00:00", errorExpected: true},
		{name: "HourTooLarge", value: "01/01/2024 24:00:00", errorExpected: true},
		{name: "MinuteTooLarge", value: "01/01/2024 12:60:00", errorExpected: true},
		{name: "SecondTooLarge", value: "01/01/2024 12:00:60", errorExpected: true},
		{name: "YearTooSmall", value: "01/01/1899 12:00:00", errorExpected: true},
		{name: "YearTooLarge", value: "01/01/2156 12:00:00", errorExpected: true},
		{name: "BadFormat", value: "01/01/2024", errorExpected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res := DPTDateTime8Byte{}
			err := res.Set(tc.value)
			if tc.errorExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.result, res)
		})
	}
}

func FuzzDPTDateTime8Byte(f *testing.F) {
	d := DPTDateTime8Byte{Year: 124, Month: 3, Day: 15, Weekday: 5, Hour: 14, Minute: 30, Second: 45, Flags: 0x02}
	data, _ := d.MarshalBinary()
	f.Add(data)
	f.Fuzz(func(t *testing.T, data []byte) {
		other := DPTDateTime8Byte{}
		_ = other.UnmarshalBinary(data)
	})
}
