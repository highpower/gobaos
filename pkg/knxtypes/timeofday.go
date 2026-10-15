package knxtypes

import (
	"fmt"
	"strconv"
	"strings"
)

// DPTTimeOfDay represents KNX DPT 10.xxx (time of day).
type DPTTimeOfDay struct {
	Weekday uint8
	Hour    uint8
	Minute  uint8
	Second  uint8
}

type timeStruct struct {
	hour   uint8
	minute uint8
	second uint8
}

// Get returns the time as a string, optionally prefixed by the weekday.
func (t *DPTTimeOfDay) Get() string {
	if s, ok := weekdayName(t.Weekday); ok {
		return fmt.Sprintf("%s,%02d:%02d:%02d", s, t.Hour, t.Minute, t.Second)
	}
	return fmt.Sprintf("%02d:%02d:%02d", t.Hour, t.Minute, t.Second)
}

// Set parses a time string and stores the fields.
func (t *DPTTimeOfDay) Set(value string) error {
	parts := strings.Split(value, ",")
	switch len(parts) {
	case 1:
		tm, err := parseTime(parts[0])
		if err != nil {
			return err
		}
		t.Hour, t.Minute, t.Second = tm.hour, tm.minute, tm.second
		return nil
	case 2:
		tm, err := parseTime(parts[1])
		if err != nil {
			return err
		}
		wd, err := parseWeekday(parts[0])
		if err != nil {
			return err
		}
		t.Weekday = wd
		t.Hour, t.Minute, t.Second = tm.hour, tm.minute, tm.second
		return nil
	}
	return fmt.Errorf("%w: timeOfDay %s", ErrBadValue, value)
}

// TypeCode returns the KNX type code for DPTTimeOfDay.
func (t *DPTTimeOfDay) TypeCode() TypeCode {
	return TimeOfDayCode
}

// Valid reports whether the time and weekday are within valid ranges.
func (t *DPTTimeOfDay) Valid() error {
	if err := validTime(t.Hour, t.Minute, t.Second); err != nil {
		return err
	}
	if err := validWeekday(t.Weekday); err != nil {
		return err
	}
	return nil
}

// String returns a human-readable representation of DPTTimeOfDay.
func (t *DPTTimeOfDay) String() string {
	wd, _ := weekdayName(t.Weekday)
	return fmt.Sprintf("dptTimeOfDay['%s',%d:%d:%d]", wd, t.Hour, t.Minute, t.Second)
}

// MarshalBinary encodes DPTTimeOfDay into a 3-byte payload.
func (t *DPTTimeOfDay) MarshalBinary() ([]byte, error) {
	return []byte{(t.Weekday&0x07)<<5 | (t.Hour & 0x1F), t.Minute & 0x3F, t.Second & 0x3F}, nil
}

// UnmarshalBinary decodes a 3-byte payload into DPTTimeOfDay.
func (t *DPTTimeOfDay) UnmarshalBinary(data []byte) error {
	if len(data) != 3 {
		return fmt.Errorf("%w: type %s requires 3 bytes, got %d", ErrUnmarshal, TimeOfDayCode.String(), len(data))
	}
	t.Weekday = (data[0] >> 5) & 0x07
	t.Hour = data[0] & 0x1F
	t.Minute = data[1] & 0x3F
	t.Second = data[2] & 0x3F
	return nil
}

func parseTime(s string) (timeStruct, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return timeStruct{}, fmt.Errorf("%w: time %s", ErrBadValue, s)
	}
	numbers := [3]uint8{0, 0, 0}
	for i, part := range parts {
		n, err := strconv.ParseUint(part, 10, 8)
		if err != nil {
			return timeStruct{}, err
		}
		numbers[i] = uint8(n)
	}
	if err := validTime(numbers[0], numbers[1], numbers[2]); err != nil {
		return timeStruct{}, err
	}
	return timeStruct{hour: numbers[0], minute: numbers[1], second: numbers[2]}, nil
}

func parseWeekday(s string) (uint8, error) {
	switch strings.ToLower(s) {
	case "mon":
		return 1, nil
	case "tue":
		return 2, nil
	case "wed":
		return 3, nil
	case "thu":
		return 4, nil
	case "fri":
		return 5, nil
	case "sat":
		return 6, nil
	case "sun":
		return 7, nil
	default:
		return 0, fmt.Errorf("%w: timeOfDay %s", ErrBadValue, s)
	}
}

// weekdayName returns the English abbreviation for a KNX weekday number (1..7).
func weekdayName(n uint8) (string, bool) {
	switch n {
	case 1:
		return "mon", true
	case 2:
		return "tue", true
	case 3:
		return "wed", true
	case 4:
		return "thu", true
	case 5:
		return "fri", true
	case 6:
		return "sat", true
	case 7:
		return "sun", true
	default:
		return "", false
	}
}

func validTime(hour, minute, second uint8) error {
	if hour > 23 {
		return fmt.Errorf("%w: hour %d out of range [0,23]", ErrBadValue, hour)
	}
	if minute > 59 {
		return fmt.Errorf("%w: minute %d out of range [0,59]", ErrUnmarshal, minute)
	}
	if second > 59 {
		return fmt.Errorf("%w: second %d out of range [0,59]", ErrUnmarshal, second)
	}
	return nil
}

func validWeekday(weekday uint8) error {
	if weekday > 7 {
		return fmt.Errorf("%w: weekday %d out of range [0,7]", ErrUnmarshal, weekday)
	}
	return nil
}
