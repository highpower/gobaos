package knxtypes

import (
	"fmt"
	"strconv"
	"strings"
)

// DPTControl represents KNX DPT 2.xxx (control with value).
type DPTControl struct {
	Control bool
	Value   bool
}

// Get returns the value and control as a colon-separated string.
func (d *DPTControl) Get() string {
	return fmt.Sprintf("%t:%t", d.Value, d.Control)
}

// Set parses a "value:control" string and stores the fields.
func (d *DPTControl) Set(s string) error {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return fmt.Errorf("%w: type %s expected format <value>:<value>", ErrBadValue, ControlCode.String())
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

// TypeCode returns the type code for DPTControl.
func (d *DPTControl) TypeCode() TypeCode {
	return ControlCode
}

// String returns a human-readable representation of DPTControl.
func (d *DPTControl) String() string {
	return fmt.Sprintf("dptControl[%t,%t]", d.Value, d.Control)
}

// MarshalBinary encodes DPTControl into a 1-byte payload.
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

// UnmarshalBinary decodes a 1-byte payload into DPTControl.
func (d *DPTControl) UnmarshalBinary(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("%w: type %s requires 1 byte, got %d", ErrUnmarshal, ControlCode.String(), len(data))
	}
	d.Value = data[0]&0x01 != 0
	d.Control = data[0]&0x02 != 0
	return nil
}
