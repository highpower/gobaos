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

var weekdays = map[uint8]string{
	1: "mon", 2: "tue", 3: "wed", 4: "thu", 5: "fri", 6: "sat", 7: "sun",
}

// Get returns the time as a string, optionally prefixed by the weekday.
func (t *DPTTimeOfDay) Get() string {
	if s, ok := weekdays[t.Weekday]; ok {
		return fmt.Sprintf("%s,%02d:%02d:%02d", s, t.Hour, t.Minute, t.Second)
	}
	return fmt.Sprintf("%02d:%02d:%02d", t.Hour, t.Minute, t.Second)
}

// Set parses a time string and stores the fields.
func (t *DPTTimeOfDay) Set(value string) error {
	parts := strings.Split(value, ",")
	switch len(parts) {
	case 1:
		tm, err := t.parseTime(parts[0])
		if err != nil {
			return err
		}
		t.Hour, t.Minute, t.Second = tm.hour, tm.minute, tm.second
		return nil
	case 2:
		tm, err := t.parseTime(parts[1])
		if err != nil {
			return err
		}
		wd, err := t.parseWeekday(parts[0])
		if err != nil {
			return err
		}
		t.Weekday = wd
		t.Hour, t.Minute, t.Second = tm.hour, tm.minute, tm.second
		return nil
	}
	return fmt.Errorf("%w: timeOfDay %s", ErrBadValue, value)
}

// TypeCode returns the type code for DPTTimeOfDay.
func (t *DPTTimeOfDay) TypeCode() TypeCode {
	return TimeOfDayCode
}

// Valid checks whether the time and weekday are within valid ranges.
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
	wd := weekdays[t.Weekday]
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

func (t *DPTTimeOfDay) parseTime(s string) (timeStruct, error) {
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

func (t *DPTTimeOfDay) parseWeekday(s string) (uint8, error) {
	for k, v := range weekdays {
		if strings.EqualFold(s, v) {
			return k, nil
		}
	}
	return 0, fmt.Errorf("%w: timeOfDay %s", ErrBadValue, s)
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
