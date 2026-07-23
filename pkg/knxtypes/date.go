package knxtypes

import (
	"fmt"
	"strconv"
	"strings"
)

type DPTDate struct {
	Day   uint8
	Month uint8
	Year  uint8
}

func (d *DPTDate) Get() string {
	return fmt.Sprintf("%02d/%02d/%02d", d.Day, d.Month, d.Year)
}

func (d *DPTDate) Set(value string) error {
	parts := strings.Split(value, "/")
	if len(parts) != 3 {
		return fmt.Errorf("%w: time %s", ErrBadValue, value)
	}
	numbers := [3]uint8{0, 0, 0}
	for i, part := range parts {
		n, err := strconv.ParseUint(part, 10, 8)
		if err != nil {
			return err
		}
		numbers[i] = uint8(n)
	}
	if err := validDate(numbers[0], numbers[1], numbers[2]); err != nil {
		return err
	}
	d.Day, d.Month, d.Year = numbers[0], numbers[1], numbers[2]
	return nil
}

func (d *DPTDate) TypeCode() TypeCode {
	return Date
}

func (d *DPTDate) Valid() error {
	return validDate(d.Day, d.Month, d.Year)
}

func (d *DPTDate) String() string {
	return fmt.Sprintf("date[%02d/%02d/%02d]", d.Day, d.Month, d.Year)
}

func (d *DPTDate) MarshalBinary() ([]byte, error) {
	return []byte{d.Day, d.Month, d.Year}, nil
}

func (d *DPTDate) UnmarshalBinary(data []byte) error {
	if len(data) != 3 {
		return fmt.Errorf("%w: type %s requires 3 bytes, got %d", ErrUnmarshal, Date.String(), len(data))
	}
	d.Day = data[0]
	d.Month = data[1]
	d.Year = data[2]
	return nil
}

func validDate(day, month, year uint8) error {
	if day < 1 || day > 31 {
		return fmt.Errorf("%w: day %d out of range [1,31]", ErrBadValue, day)
	}
	if month < 1 || month > 12 {
		return fmt.Errorf("%w: month %d out of range [1,12]", ErrBadValue, month)
	}
	if year > 99 {
		return fmt.Errorf("%w: year %d out of range [0,99]", ErrBadValue, year)
	}
	return nil
}
