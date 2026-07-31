package knxtypes

import (
	"testing"
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
			if err != nil {
				t.Fatalf("marshal: %s", err.Error())
			}
			other := DPTDateTime8Byte{}
			if err := other.UnmarshalBinary(data); err != nil {
				t.Fatalf("unmarshal: %s", err.Error())
			}
			if other != tc.value {
				t.Errorf("round-trip failed: got %v, want %v", other, tc.value)
			}
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
			result:        DPTDateTime8Byte{Year: 124, Month: 3, Day: 15, Hour: 14, Minute: 30, Second: 45},
			errorExpected: false},
		{name: "MinYear", value: "01/01/1900 00:00:00",
			result:        DPTDateTime8Byte{Year: 0, Month: 1, Day: 1, Hour: 0, Minute: 0, Second: 0},
			errorExpected: false},
		{name: "MaxYear", value: "31/12/2155 23:59:59",
			result:        DPTDateTime8Byte{Year: 255, Month: 12, Day: 31, Hour: 23, Minute: 59, Second: 59},
			errorExpected: false},
		{name: "Weekday", value: "31/12/2155,Wed 23:59:59",
			result:        DPTDateTime8Byte{Year: 255, Month: 12, Day: 31, Weekday: 3, Hour: 23, Minute: 59, Second: 59},
			errorExpected: false},
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
			switch {
			case err != nil && !tc.errorExpected:
				t.Errorf("error %s", err.Error())
			case err == nil && tc.errorExpected:
				t.Errorf("error expected")
			case err == nil && !tc.errorExpected:
				if res != tc.result {
					t.Errorf("incorrect result %v, expected %v", res, tc.result)
				}
			}
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
