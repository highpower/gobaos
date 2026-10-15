package knxtypes

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// DPTDateTime8Byte represents KNX DPT 19.xxx (date and time, 8 bytes).
// The payload layout is: Year Month Day {Weekday|Hour} Minute Second Flags.
type DPTDateTime8Byte struct {
	Day     uint8 // 1..31
	Month   uint8 // 1..12
	Year    uint8 // 0 = 1900
	Weekday uint8 // 0..7 (0 = no day)
	Hour    uint8 // 0..23
	Minute  uint8 // 0..59
	Second  uint8 // 0..59
	Flags   uint8 // status flags
}

type dateStruct struct {
	day     uint8
	month   uint8
	year    uint8
	weekday uint8
}

const (
	dayBitSize   = 8
	monthBitSize = 8
	yearBitSize  = 16
)

// Get returns the date and time as "DD/MM/YYYY HH:MM:SS".
func (d *DPTDateTime8Byte) Get() string {
	year := 1900 + int(d.Year)
	if d.Weekday == 0 {
		return fmt.Sprintf("%02d/%02d/%04d %02d:%02d:%02d", d.Day, d.Month, year, d.Hour, d.Minute, d.Second)
	}
	wd, _ := weekdayName(d.Weekday)
	return fmt.Sprintf("%02d/%02d/%04d,%s %02d:%02d:%02d", d.Day, d.Month, year, wd,
		d.Hour, d.Minute, d.Second)
}

// Set parses a "DD/MM/YYYY HH:MM:SS" or "DD/MM/YYYY,WWW HH:MM:SS" string, validates the fields, and stores them.
// Flags are reset to zero.
func (d *DPTDateTime8Byte) Set(value string) error {
	parts := strings.Split(value, " ")
	if len(parts) != 2 {
		return fmt.Errorf("%w: type %s expected format DD/MM/YYYY[,WWW] HH:MM:SS", ErrBadValue,
			DateTime8ByteCode.String())
	}
	dt, err := parseDateWeekday(parts[0])
	if err != nil {
		return err
	}
	tm, err := parseTime(parts[1])
	if err != nil {
		return err
	}
	d.Hour, d.Minute, d.Second = tm.hour, tm.minute, tm.second
	d.Day, d.Month, d.Year, d.Weekday = dt.day, dt.month, dt.year, dt.weekday
	d.Flags = 0
	return nil
}

// TypeCode returns the KNX type code for DPTDateTime8Byte.
func (d *DPTDateTime8Byte) TypeCode() TypeCode {
	return DateTime8ByteCode
}

// Valid reports whether the date and time fields are within valid ranges.
func (d *DPTDateTime8Byte) Valid() error {
	return errors.Join(validDateFullYear(d.Day, d.Month, d.Year),
		validTime(d.Hour, d.Minute, d.Second))
}

// String returns a human-readable representation of DPTDateTime8Byte.
func (d *DPTDateTime8Byte) String() string {
	return fmt.Sprintf("dptDateTime8Byte[%s]", d.Get())
}

// MarshalBinary encodes DPTDateTime8Byte into an 8-byte payload.
func (d *DPTDateTime8Byte) MarshalBinary() ([]byte, error) {
	return []byte{d.Year, d.Month, d.Day, (d.Weekday&0x07)<<5 | (d.Hour & 0x1F),
		d.Minute & 0x3F, d.Second & 0x3F, d.Flags, 0}, nil
}

// UnmarshalBinary decodes an 8-byte payload into DPTDateTime8Byte.
func (d *DPTDateTime8Byte) UnmarshalBinary(data []byte) error {
	if len(data) != 8 {
		return fmt.Errorf("%w: type %s requires 8 bytes, got %d", ErrUnmarshal, DateTime8ByteCode.String(), len(data))
	}
	d.Hour, d.Minute, d.Second = data[3]&0x1F, data[4]&0x3F, data[5]&0x3F
	d.Day, d.Month, d.Year, d.Weekday = data[2], data[1], data[0], (data[3]>>5)&0x07
	d.Flags = data[6]
	return nil
}

func validDateFullYear(day, month, _ uint8) error {
	if month < 1 || month > 12 {
		return fmt.Errorf("%w: month %d out of range [1,12]", ErrBadValue, month)
	}
	if day < 1 || day > 31 {
		return fmt.Errorf("%w: day %d out of range [1,31]", ErrBadValue, day)
	}
	return nil
}

func parseDateWeekday(s string) (dateStruct, error) {
	parts := strings.Split(s, ",")
	switch len(parts) {
	case 1:
		return parseDateFullYear(parts[0])
	case 2:
		date, err := parseDateWeekday(parts[0])
		if err != nil {
			return dateStruct{}, err
		}
		weekday, err := parseWeekday(parts[1])
		if err != nil {
			return dateStruct{}, err
		}
		date.weekday = weekday
		return date, nil
	default:
		return dateStruct{}, fmt.Errorf("%w: invalid format %q", ErrBadValue, s)
	}
}

func parseDateFullYear(s string) (dateStruct, error) {
	parts := strings.Split(s, "/")
	if len(parts) != 3 {
		return dateStruct{}, fmt.Errorf("%w: invalid date %s", ErrBadValue, parts[0])
	}
	bits := []int{dayBitSize, monthBitSize, yearBitSize}
	numbers := [3]uint64{0, 0, 0}
	for i, part := range parts {
		n, err := strconv.ParseUint(part, 10, bits[i])
		if err != nil {
			return dateStruct{}, err
		}
		numbers[i] = n
	}
	day, month, year := uint8(numbers[0]), uint8(numbers[1]), numbers[2]
	if year < 1900 || year > 2155 {
		return dateStruct{}, fmt.Errorf("%w: year %d out of range [1900,2155]", ErrBadValue, year)
	}
	yearOffset := uint8(year - 1900)
	if err := validDateFullYear(day, month, yearOffset); err != nil {
		return dateStruct{}, err
	}
	return dateStruct{day: day, month: month, year: yearOffset}, nil
}
