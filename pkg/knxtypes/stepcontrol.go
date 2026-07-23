package knxtypes

import (
	"fmt"
	"strconv"
	"strings"
)

type DPTStepControl struct {
	Control  bool
	StepCode uint8
}

func (s *DPTStepControl) Get() string {
	return fmt.Sprintf("%t:%d", s.Control, s.StepCode)
}

func (s *DPTStepControl) Set(value string) error {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return fmt.Errorf("%w: type %s, requires 2 parts, got %d", ErrBadValue, Control.String(), len(parts))
	}
	control, err := strconv.ParseBool(parts[0])
	if err != nil {
		return err
	}
	stepCode, err := strconv.ParseUint(parts[1], 10, 8)
	if err != nil {
		return err
	}
	if stepCode > 7 {
		return fmt.Errorf("%w: step code %d out of range [0,7]", ErrBadValue, stepCode)
	}
	s.Control = control
	s.StepCode = uint8(stepCode)
	return nil
}

func (s *DPTStepControl) TypeCode() TypeCode {
	return StepControl
}

func (s *DPTStepControl) BytesRequired() int {
	return 1
}

func (s *DPTStepControl) String() string {
	return fmt.Sprintf("stepControl[%t:%d]", s.Control, s.StepCode)
}

func (s *DPTStepControl) MarshalBinary() ([]byte, error) {
	if s.StepCode > 7 {
		return nil, fmt.Errorf("%w: step code %d out of range [0,7]", ErrBadValue, s.StepCode)
	}
	result := []byte{0}
	if s.Control {
		result[0] |= 0x08
	}
	result[0] |= s.StepCode & 0x07
	return result, nil
}

func (s *DPTStepControl) UnmarshalBinary(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("%w: type %s expected 1 byte, got %d", ErrUnmarshal, StepControl.String(), len(data))
	}
	s.Control = data[0]&0x08 != 0
	s.StepCode = data[0] & 0x07
	return nil
}
