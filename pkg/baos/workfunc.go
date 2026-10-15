package baos

import (
	"context"

	"gobaos/pkg/knxtypes"
)

// Control represents a boolean control value with a separate control bit.
type Control struct {
	Control bool
	Value   bool
}

// StepControl represents a step-control value.
type StepControl struct {
	Control  bool
	StepCode uint8
}

// TimeOfDay represents a time-of-day value.
type TimeOfDay struct {
	Weekday uint8
	Hour    uint8
	Minute  uint8
	Second  uint8
}

// Date represents a date value.
type Date struct {
	Day   uint8
	Month uint8
	Year  uint8
}

// Access represents an access-control value.
type Access struct {
	Key            uint8
	ID             uint8
	Index          uint8
	AdditionalData uint8
}

// SceneControl represents a scene-control value.
type SceneControl struct {
	Save  bool
	Scene uint8
}

// DateTime represents an 8-byte date-and-time value.
type DateTime struct {
	Day     uint8
	Month   uint8
	Year    uint8
	Weekday uint8
	Hour    uint8
	Minute  uint8
	Second  uint8
	Flags   uint8
}

// RGB represents a red-green-blue color value.
type RGB struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

// RGBW represents a red-green-blue-white color value.
type RGBW struct {
	Red   uint8
	Green uint8
	Blue  uint8
	White uint8
	Flags uint16
}

// GetBool reads the boolean value of the given datapoint.
func GetBool(ctx context.Context, w *Worker, id DatapointID) (bool, error) {
	dpt := knxtypes.DPTBool{}
	if err := w.Get(ctx, id, &dpt); err != nil {
		return false, err
	}
	return dpt.Value, nil
}

// SetBool writes a boolean value to the given datapoint.
func SetBool(ctx context.Context, w *Worker, id DatapointID, value bool) error {
	dpt := knxtypes.DPTBool{Value: value}
	return w.Set(ctx, id, &dpt)
}

// GetControl reads the control value of the given datapoint.
func GetControl(ctx context.Context, w *Worker, id DatapointID) (Control, error) {
	dpt := knxtypes.DPTControl{}
	if err := w.Get(ctx, id, &dpt); err != nil {
		return Control{}, err
	}
	return Control{Control: dpt.Control, Value: dpt.Value}, nil
}

// SetControl writes a control value to the given datapoint.
func SetControl(ctx context.Context, w *Worker, id DatapointID, control Control) error {
	dpt := knxtypes.DPTControl{Control: control.Control, Value: control.Value}
	return w.Set(ctx, id, &dpt)
}

// GetStepControl reads the step-control value of the given datapoint.
func GetStepControl(ctx context.Context, w *Worker, id DatapointID) (StepControl, error) {
	dpt := knxtypes.DPTStepControl{}
	if err := w.Get(ctx, id, &dpt); err != nil {
		return StepControl{}, err
	}
	return StepControl{Control: dpt.Control, StepCode: dpt.StepCode}, nil
}

// SetStepControl writes a step-control value to the given datapoint.
func SetStepControl(ctx context.Context, w *Worker, id DatapointID, value StepControl) error {
	dpt := knxtypes.DPTStepControl{Control: value.Control, StepCode: value.StepCode}
	return w.Set(ctx, id, &dpt)
}

// GetChar reads the character value of the given datapoint.
func GetChar(ctx context.Context, w *Worker, id DatapointID) (rune, error) {
	dpt := knxtypes.DPTChar{}
	if err := w.Get(ctx, id, &dpt); err != nil {
		return 0, err
	}
	return dpt.Value, nil
}

// SetChar writes a character value to the given datapoint.
func SetChar(ctx context.Context, w *Worker, id DatapointID, value rune) error {
	dpt := knxtypes.DPTChar{Value: value}
	return w.Set(ctx, id, &dpt)
}

// GetUint8 reads the unsigned 8-bit integer value of the given datapoint.
func GetUint8(ctx context.Context, w *Worker, id DatapointID) (uint8, error) {
	dpt := knxtypes.DPTUint8{}
	if err := w.Get(ctx, id, &dpt); err != nil {
		return 0, err
	}
	return dpt.Value, nil
}

// SetUint8 writes an unsigned 8-bit integer value to the given datapoint.
func SetUint8(ctx context.Context, w *Worker, id DatapointID, value uint8) error {
	dpt := knxtypes.DPTUint8{Value: value}
	return w.Set(ctx, id, &dpt)
}

// GetInt8 reads the signed 8-bit integer value of the given datapoint.
func GetInt8(ctx context.Context, w *Worker, id DatapointID) (int8, error) {
	dpt := knxtypes.DPTInt8{}
	if err := w.Get(ctx, id, &dpt); err != nil {
		return 0, err
	}
	return dpt.Value, nil
}

// SetInt8 writes a signed 8-bit integer value to the given datapoint.
func SetInt8(ctx context.Context, w *Worker, id DatapointID, value int8) error {
	dpt := knxtypes.DPTInt8{Value: value}
	return w.Set(ctx, id, &dpt)
}

// GetUint16 reads the unsigned 16-bit integer value of the given datapoint.
func GetUint16(ctx context.Context, w *Worker, id DatapointID) (uint16, error) {
	dpt := knxtypes.DPTUint16{}
	if err := w.Get(ctx, id, &dpt); err != nil {
		return 0, err
	}
	return dpt.Value, nil
}

// SetUint16 writes an unsigned 16-bit integer value to the given datapoint.
func SetUint16(ctx context.Context, w *Worker, id DatapointID, value uint16) error {
	dpt := knxtypes.DPTUint16{Value: value}
	return w.Set(ctx, id, &dpt)
}

// GetInt16 reads the signed 16-bit integer value of the given datapoint.
func GetInt16(ctx context.Context, w *Worker, id DatapointID) (int16, error) {
	dpt := knxtypes.DPTInt16{}
	if err := w.Get(ctx, id, &dpt); err != nil {
		return 0, err
	}
	return dpt.Value, nil
}

// SetInt16 writes a signed 16-bit integer value to the given datapoint.
func SetInt16(ctx context.Context, w *Worker, id DatapointID, value int16) error {
	dpt := knxtypes.DPTInt16{Value: value}
	return w.Set(ctx, id, &dpt)
}

// GetFloat reads the 2-byte float value of the given datapoint.
func GetFloat(ctx context.Context, w *Worker, id DatapointID) (float64, error) {
	dpt := knxtypes.DPTFloat{}
	if err := w.Get(ctx, id, &dpt); err != nil {
		return 0, err
	}
	return dpt.Value, nil
}

// SetFloat writes a 2-byte float value to the given datapoint.
func SetFloat(ctx context.Context, w *Worker, id DatapointID, value float64) error {
	dpt := knxtypes.DPTFloat{Value: value}
	return w.Set(ctx, id, &dpt)
}

// GetTimeOfDay reads the time-of-day value of the given datapoint.
func GetTimeOfDay(ctx context.Context, w *Worker, id DatapointID) (TimeOfDay, error) {
	dpt := knxtypes.DPTTimeOfDay{}
	if err := w.Get(ctx, id, &dpt); err != nil {
		return TimeOfDay{}, err
	}
	return TimeOfDay{
		Weekday: dpt.Weekday,
		Hour:    dpt.Hour,
		Minute:  dpt.Minute,
		Second:  dpt.Second,
	}, nil
}

// SetTimeOfDay writes a time-of-day value to the given datapoint.
func SetTimeOfDay(ctx context.Context, w *Worker, id DatapointID, value TimeOfDay) error {
	dpt := knxtypes.DPTTimeOfDay{
		Weekday: value.Weekday,
		Hour:    value.Hour,
		Minute:  value.Minute,
		Second:  value.Second,
	}
	return w.Set(ctx, id, &dpt)
}

// GetDate reads the date value of the given datapoint.
func GetDate(ctx context.Context, w *Worker, id DatapointID) (Date, error) {
	dpt := knxtypes.DPTDate{}
	if err := w.Get(ctx, id, &dpt); err != nil {
		return Date{}, err
	}
	return Date{Day: dpt.Day, Month: dpt.Month, Year: dpt.Year}, nil
}

// SetDate writes a date value to the given datapoint.
func SetDate(ctx context.Context, w *Worker, id DatapointID, value Date) error {
	dpt := knxtypes.DPTDate{Day: value.Day, Month: value.Month, Year: value.Year}
	return w.Set(ctx, id, &dpt)
}

// GetUint32 reads the unsigned 32-bit integer value of the given datapoint.
func GetUint32(ctx context.Context, w *Worker, id DatapointID) (uint32, error) {
	dpt := knxtypes.DPTUint32{}
	if err := w.Get(ctx, id, &dpt); err != nil {
		return 0, err
	}
	return dpt.Value, nil
}

// SetUint32 writes an unsigned 32-bit integer value to the given datapoint.
func SetUint32(ctx context.Context, w *Worker, id DatapointID, value uint32) error {
	dpt := knxtypes.DPTUint32{Value: value}
	return w.Set(ctx, id, &dpt)
}

// GetInt32 reads the signed 32-bit integer value of the given datapoint.
func GetInt32(ctx context.Context, w *Worker, id DatapointID) (int32, error) {
	dpt := knxtypes.DPTInt32{}
	if err := w.Get(ctx, id, &dpt); err != nil {
		return 0, err
	}
	return dpt.Value, nil
}

// SetInt32 writes a signed 32-bit integer value to the given datapoint.
func SetInt32(ctx context.Context, w *Worker, id DatapointID, value int32) error {
	dpt := knxtypes.DPTInt32{Value: value}
	return w.Set(ctx, id, &dpt)
}

// GetFloat4Byte reads the 4-byte float value of the given datapoint.
func GetFloat4Byte(ctx context.Context, w *Worker, id DatapointID) (float64, error) {
	dpt := knxtypes.DPTFloat4Byte{}
	if err := w.Get(ctx, id, &dpt); err != nil {
		return 0, err
	}
	return dpt.Value, nil
}

// SetFloat4Byte writes a 4-byte float value to the given datapoint.
func SetFloat4Byte(ctx context.Context, w *Worker, id DatapointID, value float64) error {
	dpt := knxtypes.DPTFloat4Byte{Value: value}
	return w.Set(ctx, id, &dpt)
}

// GetAccess reads the access-control value of the given datapoint.
func GetAccess(ctx context.Context, w *Worker, id DatapointID) (Access, error) {
	dpt := knxtypes.DPTAccess{}
	if err := w.Get(ctx, id, &dpt); err != nil {
		return Access{}, err
	}
	return Access{
		Key:            dpt.Key,
		ID:             dpt.ID,
		Index:          dpt.Index,
		AdditionalData: dpt.AdditionalData,
	}, nil
}

// SetAccess writes an access-control value to the given datapoint.
func SetAccess(ctx context.Context, w *Worker, id DatapointID, value Access) error {
	dpt := knxtypes.DPTAccess{
		Key:            value.Key,
		ID:             value.ID,
		Index:          value.Index,
		AdditionalData: value.AdditionalData,
	}
	return w.Set(ctx, id, &dpt)
}

// GetString reads the string value of the given datapoint.
func GetString(ctx context.Context, w *Worker, id DatapointID) (string, error) {
	dpt := knxtypes.DPTString{}
	if err := w.Get(ctx, id, &dpt); err != nil {
		return "", err
	}
	return dpt.Value, nil
}

// SetString writes a string value to the given datapoint.
func SetString(ctx context.Context, w *Worker, id DatapointID, value string) error {
	dpt := knxtypes.DPTString{Value: value}
	return w.Set(ctx, id, &dpt)
}

// GetSceneNumber reads the scene number value of the given datapoint.
func GetSceneNumber(ctx context.Context, w *Worker, id DatapointID) (uint8, error) {
	dpt := knxtypes.DPTSceneNumber{}
	if err := w.Get(ctx, id, &dpt); err != nil {
		return 0, err
	}
	return dpt.Value, nil
}

// SetSceneNumber writes a scene number value to the given datapoint.
func SetSceneNumber(ctx context.Context, w *Worker, id DatapointID, value uint8) error {
	dpt := knxtypes.DPTSceneNumber{Value: value}
	return w.Set(ctx, id, &dpt)
}

// GetSceneControl reads the scene-control value of the given datapoint.
func GetSceneControl(ctx context.Context, w *Worker, id DatapointID) (SceneControl, error) {
	dpt := knxtypes.DPTSceneControl{}
	if err := w.Get(ctx, id, &dpt); err != nil {
		return SceneControl{}, err
	}
	return SceneControl{Save: dpt.Save, Scene: dpt.Scene}, nil
}

// SetSceneControl writes a scene-control value to the given datapoint.
func SetSceneControl(ctx context.Context, w *Worker, id DatapointID, value SceneControl) error {
	dpt := knxtypes.DPTSceneControl{Save: value.Save, Scene: value.Scene}
	return w.Set(ctx, id, &dpt)
}

// GetDateTime reads the date-time value of the given datapoint.
func GetDateTime(ctx context.Context, w *Worker, id DatapointID) (DateTime, error) {
	dpt := knxtypes.DPTDateTime8Byte{}
	if err := w.Get(ctx, id, &dpt); err != nil {
		return DateTime{}, err
	}
	return DateTime{
		Day:     dpt.Day,
		Month:   dpt.Month,
		Year:    dpt.Year,
		Weekday: dpt.Weekday,
		Hour:    dpt.Hour,
		Minute:  dpt.Minute,
		Second:  dpt.Second,
		Flags:   dpt.Flags,
	}, nil
}

// SetDateTime writes a date-time value to the given datapoint.
func SetDateTime(ctx context.Context, w *Worker, id DatapointID, value DateTime) error {
	dpt := knxtypes.DPTDateTime8Byte{
		Day:     value.Day,
		Month:   value.Month,
		Year:    value.Year,
		Weekday: value.Weekday,
		Hour:    value.Hour,
		Minute:  value.Minute,
		Second:  value.Second,
		Flags:   value.Flags,
	}
	return w.Set(ctx, id, &dpt)
}

// GetEnum reads the enumeration value of the given datapoint.
func GetEnum(ctx context.Context, w *Worker, id DatapointID) (uint8, error) {
	dpt := knxtypes.DPTEnum{}
	if err := w.Get(ctx, id, &dpt); err != nil {
		return 0, err
	}
	return dpt.Value, nil
}

// SetEnum writes an enumeration value to the given datapoint.
func SetEnum(ctx context.Context, w *Worker, id DatapointID, value uint8) error {
	dpt := knxtypes.DPTEnum{Value: value}
	return w.Set(ctx, id, &dpt)
}

// GetRGB reads the RGB color value of the given datapoint.
func GetRGB(ctx context.Context, w *Worker, id DatapointID) (RGB, error) {
	dpt := knxtypes.DPTRedGreenBlue{}
	if err := w.Get(ctx, id, &dpt); err != nil {
		return RGB{}, err
	}
	return RGB{Red: dpt.Red, Green: dpt.Green, Blue: dpt.Blue}, nil
}

// SetRGB writes an RGB color value to the given datapoint.
func SetRGB(ctx context.Context, w *Worker, id DatapointID, value RGB) error {
	dpt := knxtypes.DPTRedGreenBlue{Red: value.Red, Green: value.Green, Blue: value.Blue}
	return w.Set(ctx, id, &dpt)
}

// GetRGBW reads the RGBW color value of the given datapoint.
func GetRGBW(ctx context.Context, w *Worker, id DatapointID) (RGBW, error) {
	dpt := knxtypes.DPTRedGreenBlueWhite{}
	if err := w.Get(ctx, id, &dpt); err != nil {
		return RGBW{}, err
	}
	return RGBW{
		Red:   dpt.Red,
		Green: dpt.Green,
		Blue:  dpt.Blue,
		White: dpt.White,
		Flags: dpt.Flags,
	}, nil
}

// SetRGBW writes an RGBW color value to the given datapoint.
func SetRGBW(ctx context.Context, w *Worker, id DatapointID, value RGBW) error {
	dpt := knxtypes.DPTRedGreenBlueWhite{
		Red:   value.Red,
		Green: value.Green,
		Blue:  value.Blue,
		White: value.White,
		Flags: value.Flags,
	}
	return w.Set(ctx, id, &dpt)
}
