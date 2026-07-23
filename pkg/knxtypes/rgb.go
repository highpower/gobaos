package knxtypes

import (
	"fmt"
	"strconv"
	"strings"
)

type DPTRedGreenBlue struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

func (d *DPTRedGreenBlue) Get() string {
	return fmt.Sprintf("%d,%d,%d", d.Red, d.Green, d.Blue)
}

func (d *DPTRedGreenBlue) Set(value string) error {
	parts := strings.Split(value, ",")
	if len(parts) != 3 {
		return fmt.Errorf("%w redGreenBlue '%s'", ErrBadValue, value)
	}
	numbers := [3]uint8{0, 0, 0}
	for i, part := range parts {
		v, err := strconv.ParseUint(part, 10, 8)
		if err != nil {
			return err
		}
		numbers[i] = uint8(v)
	}
	d.Red, d.Green, d.Blue = numbers[0], numbers[1], numbers[2]
	return nil
}

func (d *DPTRedGreenBlue) TypeCode() TypeCode {
	return RedGreenBlue
}

func (d *DPTRedGreenBlue) String() string {
	return fmt.Sprintf("dptRedGreenBlue[%d,%d,%d]", d.Red, d.Green, d.Blue)
}

func (d *DPTRedGreenBlue) MarshalBinary() ([]byte, error) {
	return []byte{d.Red, d.Green, d.Blue}, nil
}

func (d *DPTRedGreenBlue) UnmarshalBinary(data []byte) error {
	if len(data) != 3 {
		return fmt.Errorf("%w: type %s requires 3 bytes, got %d", ErrUnmarshal, RedGreenBlue.String(), len(data))
	}
	d.Red = data[0]
	d.Green = data[1]
	d.Blue = data[2]
	return nil
}
