package knxtypes

import (
	"fmt"
	"strconv"
	"strings"
)

// DPTSceneControl represents KNX DPT 18.xxx (scene control).
type DPTSceneControl struct {
	Save  bool  // true = learn/save, false = execute/recall
	Scene uint8 // 1..64
}

// Get returns the scene control as "scene,save".
func (c *DPTSceneControl) Get() string {
	return fmt.Sprintf("%d,%t", c.Scene, c.Save)
}

// Set parses a "scene,save" string and updates the control values after validating the scene number.
func (c *DPTSceneControl) Set(value string) error {
	parts := strings.Split(value, ",")
	if len(parts) != 2 {
		return fmt.Errorf("%w: type %s expected format <value>,<value>", ErrBadValue, SceneControlCode.String())
	}
	scene, err := strconv.ParseUint(parts[0], 10, 8)
	if err != nil {
		return err
	}
	save, err := strconv.ParseBool(parts[1])
	if err != nil {
		return err
	}
	if err := validSceneNumber(uint8(scene)); err != nil {
		return err
	}
	c.Scene, c.Save = uint8(scene), save
	return nil
}

// TypeCode returns the KNX type code for DPTSceneControl.
func (c *DPTSceneControl) TypeCode() TypeCode {
	return SceneControlCode
}

// String returns a human-readable representation of DPTSceneControl.
func (c *DPTSceneControl) String() string {
	return fmt.Sprintf("dptSceneControl[%d,%t]", c.Scene, c.Save)
}

// Valid checks whether the scene number is within the valid range.
func (c *DPTSceneControl) Valid() error {
	return validSceneNumber(c.Scene)
}

// MarshalBinary encodes DPTSceneControl into a 1-byte payload.
func (c *DPTSceneControl) MarshalBinary() ([]byte, error) {
	result := []byte{0}
	if c.Save {
		result[0] |= 0x80
	}
	result[0] |= (c.Scene - 1) & 0x3F
	return result, nil
}

// UnmarshalBinary decodes a 1-byte payload into DPTSceneControl.
func (c *DPTSceneControl) UnmarshalBinary(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("%w: type %s requires 1 byte, got %d", ErrUnmarshal, SceneControlCode.String(), len(data))
	}
	c.Save = data[0]&0x80 != 0
	c.Scene = (data[0] & 0x3F) + 1
	return nil
}
