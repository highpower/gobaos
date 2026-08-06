package knxtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			require.NoError(t, err)

			other := DPTDate{}
			require.NoError(t, other.UnmarshalBinary(data))
			assert.Equal(t, dpt.Day, other.Day)
			assert.Equal(t, dpt.Month, other.Month)
			assert.Equal(t, dpt.Year, other.Year)
		})
	}
}

func TestDPTDateBadUnmarshal(t *testing.T) {
	testBadUnmarshal[DPTDate](t)
}

func TestDPTDateSet(t *testing.T) {
	type testCase struct {
		name     string
		input    string
		expected DPTDate
		wantErr  bool
	}

	testCases := []testCase{
		{name: "Valid", input: "15/06/23", expected: DPTDate{Day: 15, Month: 6, Year: 23}},
		{name: "WrongParts", input: "15/06/23/extra", wantErr: true},
		{name: "NotADay", input: "bad/06/23", wantErr: true},
		{name: "NotAMonth", input: "15/bad/23", wantErr: true},
		{name: "NotAYear", input: "15/06/bad", wantErr: true},
		{name: "BadDay", input: "32/06/23", wantErr: true},
		{name: "BadMonth", input: "15/13/23", wantErr: true},
		{name: "BadYear", input: "15/06/106", wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dpt := DPTDate{}
			err := dpt.Set(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, dpt)
		})
	}
}

func TestDPTDateValid(t *testing.T) {
	type testCase struct {
		name    string
		dpt     DPTDate
		wantErr bool
	}

	testCases := []testCase{
		{name: "Min", dpt: DPTDate{Day: 1, Month: 1, Year: 0}},
		{name: "Max", dpt: DPTDate{Day: 31, Month: 12, Year: 99}},
		{name: "DayZero", dpt: DPTDate{Day: 0, Month: 1, Year: 0}, wantErr: true},
		{name: "BadDay", dpt: DPTDate{Day: 32, Month: 1, Year: 0}, wantErr: true},
		{name: "BadMonth", dpt: DPTDate{Day: 1, Month: 0, Year: 0}, wantErr: true},
		{name: "BadMonthMax", dpt: DPTDate{Day: 1, Month: 13, Year: 0}, wantErr: true},
		{name: "BadYear", dpt: DPTDate{Day: 1, Month: 1, Year: 100}, wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.dpt.Valid()
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
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
