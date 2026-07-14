package knxtypes

import "fmt"

type DPTBool struct {
	Value bool
}

//goland:noinspection GoMixedReceiverTypes
func (d DPTBool) String() string {
	return fmt.Sprintf("dptBool[%t]", d.Value)
}

//goland:noinspection GoMixedReceiverTypes
func (d DPTBool) MarshalBinary() ([]byte, error) {
	result := []byte{0}
	if d.Value {
		result[0] = 1
	}
	return result, nil
}

//goland:noinspection GoMixedReceiverTypes
func (d *DPTBool) UnmarshalBinary(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("DPTBool %w: expected 1 byte, got %d", ErrUnmarshal, len(data))
	}
	d.Value = data[0]&0x01 != 0
	return nil
}

func NewDPTBool() DPTBool {
	return DPTBool{}
}
