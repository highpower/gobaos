package knxtypes

import (
	"fmt"
	"strconv"
	"strings"
)

type DPTControl struct {
	Control bool
	Value   bool
}

func (d *DPTControl) Get() string {
	return fmt.Sprintf("%t:%t", d.Value, d.Control)
}

func (d *DPTControl) Set(s string) error {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return fmt.Errorf("%w: type %s expected format <value>:<value>", ErrBadValue, Control.String())
	}
	value, err := strconv.ParseBool(parts[0])
	if err != nil {
		return err
	}
	control, err := strconv.ParseBool(parts[1])
	if err != nil {
		return err
	}
	d.Value = value
	d.Control = control
	return nil
}

func (d *DPTControl) TypeCode() TypeCode {
	return Control
}

func (d *DPTControl) String() string {
	return fmt.Sprintf("dptControl[value=%t,control=%t]", d.Value, d.Control)
}

func (d *DPTControl) MarshalBinary() ([]byte, error) {
	result := []byte{0}
	if d.Control {
		result[0] |= 0x02
	}
	if d.Value {
		result[0] |= 0x01
	}
	return result, nil
}

func (d *DPTControl) UnmarshalBinary(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("%w: type %s, requires 1 byte, got %d", ErrUnmarshal, Control.String(), len(data))
	}
	d.Value = data[0]&0x01 != 0
	d.Control = data[0]&0x02 != 0
	return nil
}
