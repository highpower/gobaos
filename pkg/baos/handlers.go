package baos

import (
	"gobaos/pkg/knxtypes"
)

// Handlers maps datapoint IDs to type-safe indication handlers.
type Handlers map[DatapointID]ValueIndicationFunc

type dptWrapper[G any] interface {
	*G
	knxtypes.DPT
}

// Bool registers a handler for boolean DPT indications.
func (h Handlers) Bool(ids []DatapointID, f func(id DatapointID, value bool) error) {
	register[knxtypes.DPTBool](h, ids, func(id DatapointID, dpt knxtypes.DPTBool) error {
		return f(id, dpt.Value)
	})
}

// Control registers a handler for control DPT indications.
func (h Handlers) Control(ids []DatapointID, f func(id DatapointID, value Control) error) {
	register[knxtypes.DPTControl](h, ids, func(id DatapointID, dpt knxtypes.DPTControl) error {
		control := Control{Value: dpt.Value, Control: dpt.Control}
		return f(id, control)
	})
}

// StepControl registers a handler for step-control DPT indications.
func (h Handlers) StepControl(ids []DatapointID, f func(id DatapointID, value StepControl) error) {
	register[knxtypes.DPTStepControl](h, ids, func(id DatapointID, dpt knxtypes.DPTStepControl) error {
		stepControl := StepControl{Control: dpt.Control, StepCode: dpt.StepCode}
		return f(id, stepControl)
	})
}

// Char registers a handler for character DPT indications.
func (h Handlers) Char(ids []DatapointID, f func(id DatapointID, value rune) error) {
	register[knxtypes.DPTChar](h, ids, func(id DatapointID, dpt knxtypes.DPTChar) error {
		return f(id, dpt.Value)
	})
}

// Uint8 registers a handler for unsigned 8-bit integer DPT indications.
func (h Handlers) Uint8(ids []DatapointID, f func(id DatapointID, value uint8) error) {
	register[knxtypes.DPTUint8](h, ids, func(id DatapointID, dpt knxtypes.DPTUint8) error {
		return f(id, dpt.Value)
	})
}

// Int8 registers a handler for signed 8-bit integer DPT indications.
func (h Handlers) Int8(ids []DatapointID, f func(id DatapointID, value int8) error) {
	register[knxtypes.DPTInt8](h, ids, func(id DatapointID, dpt knxtypes.DPTInt8) error {
		return f(id, dpt.Value)
	})
}

// Uint16 registers a handler for unsigned 16-bit integer DPT indications.
func (h Handlers) Uint16(ids []DatapointID, f func(id DatapointID, value uint16) error) {
	register[knxtypes.DPTUint16](h, ids, func(id DatapointID, dpt knxtypes.DPTUint16) error {
		return f(id, dpt.Value)
	})
}

// Int16 registers a handler for signed 16-bit integer DPT indications.
func (h Handlers) Int16(ids []DatapointID, f func(id DatapointID, value int16) error) {
	register[knxtypes.DPTInt16](h, ids, func(id DatapointID, dpt knxtypes.DPTInt16) error {
		return f(id, dpt.Value)
	})
}

// Float registers a handler for 2-byte float DPT indications.
func (h Handlers) Float(ids []DatapointID, f func(id DatapointID, value float64) error) {
	register[knxtypes.DPTFloat](h, ids, func(id DatapointID, dpt knxtypes.DPTFloat) error {
		return f(id, dpt.Value)
	})
}

// TimeOfDay registers a handler for time-of-day DPT indications.
func (h Handlers) TimeOfDay(ids []DatapointID, f func(id DatapointID, value TimeOfDay) error) {
	register[knxtypes.DPTTimeOfDay](h, ids, func(id DatapointID, dpt knxtypes.DPTTimeOfDay) error {
		timeOfDay := TimeOfDay{Weekday: dpt.Weekday, Hour: dpt.Hour, Minute: dpt.Minute, Second: dpt.Second}
		return f(id, timeOfDay)
	})
}

// Date registers a handler for date DPT indications.
func (h Handlers) Date(ids []DatapointID, f func(id DatapointID, value Date) error) {
	register[knxtypes.DPTDate](h, ids, func(id DatapointID, dpt knxtypes.DPTDate) error {
		date := Date{Day: dpt.Day, Month: dpt.Month, Year: dpt.Year}
		return f(id, date)
	})
}

// Uint32 registers a handler for unsigned 32-bit integer DPT indications.
func (h Handlers) Uint32(ids []DatapointID, f func(id DatapointID, value uint32) error) {
	register[knxtypes.DPTUint32](h, ids, func(id DatapointID, dpt knxtypes.DPTUint32) error {
		return f(id, dpt.Value)
	})
}

// Int32 registers a handler for signed 32-bit integer DPT indications.
func (h Handlers) Int32(ids []DatapointID, f func(id DatapointID, value int32) error) {
	register[knxtypes.DPTInt32](h, ids, func(id DatapointID, dpt knxtypes.DPTInt32) error {
		return f(id, dpt.Value)
	})
}

// Float4Byte registers a handler for 4-byte float DPT indications.
func (h Handlers) Float4Byte(ids []DatapointID, f func(id DatapointID, value float64) error) {
	register[knxtypes.DPTFloat4Byte](h, ids, func(id DatapointID, dpt knxtypes.DPTFloat4Byte) error {
		return f(id, dpt.Value)
	})
}

// Access registers a handler for access-control DPT indications.
func (h Handlers) Access(ids []DatapointID, f func(value Access) error) {
	register[knxtypes.DPTAccess](h, ids, func(id DatapointID, dpt knxtypes.DPTAccess) error {
		access := Access{Key: dpt.Key, ID: dpt.ID, Index: dpt.Index, AdditionalData: dpt.AdditionalData}
		return f(access)
	})
}

// String registers a handler for string DPT indications.
func (h Handlers) String(ids []DatapointID, f func(id DatapointID, value string) error) {
	register[knxtypes.DPTString](h, ids, func(id DatapointID, dpt knxtypes.DPTString) error {
		return f(id, dpt.Value)
	})
}

// SceneNumber registers a handler for scene-number DPT indications.
func (h Handlers) SceneNumber(ids []DatapointID, f func(id DatapointID, value uint8) error) {
	register[knxtypes.DPTSceneNumber](h, ids, func(id DatapointID, dpt knxtypes.DPTSceneNumber) error {
		return f(id, dpt.Value)
	})
}

// SceneControl registers a handler for scene-control DPT indications.
func (h Handlers) SceneControl(ids []DatapointID, f func(id DatapointID, value SceneControl) error) {
	register[knxtypes.DPTSceneControl](h, ids, func(id DatapointID, dpt knxtypes.DPTSceneControl) error {
		sceneControl := SceneControl{Save: dpt.Save, Scene: dpt.Scene}
		return f(id, sceneControl)
	})
}

// DateTime registers a handler for 8-byte date-time DPT indications.
func (h Handlers) DateTime(ids []DatapointID, f func(id DatapointID, value DateTime) error) {
	register[knxtypes.DPTDateTime8Byte](h, ids, func(id DatapointID, dpt knxtypes.DPTDateTime8Byte) error {
		dateTime := DateTime{Day: dpt.Day, Month: dpt.Month, Year: dpt.Year, Weekday: dpt.Weekday,
			Hour: dpt.Hour, Minute: dpt.Minute, Second: dpt.Second, Flags: dpt.Flags}
		return f(id, dateTime)
	})
}

// Enum registers a handler for enumeration DPT indications.
func (h Handlers) Enum(ids []DatapointID, f func(id DatapointID, value uint8) error) {
	register[knxtypes.DPTEnum](h, ids, func(id DatapointID, dpt knxtypes.DPTEnum) error {
		return f(id, dpt.Value)
	})
}

// RGB registers a handler for RGB color DPT indications.
func (h Handlers) RGB(ids []DatapointID, f func(id DatapointID, value RGB) error) {
	register[knxtypes.DPTRedGreenBlue](h, ids, func(id DatapointID, dpt knxtypes.DPTRedGreenBlue) error {
		rgb := RGB{Red: dpt.Red, Green: dpt.Green, Blue: dpt.Blue}
		return f(id, rgb)
	})
}

// RGBW registers a handler for RGBW color DPT indications.
func (h Handlers) RGBW(ids []DatapointID, f func(id DatapointID, value RGBW) error) {
	register[knxtypes.DPTRedGreenBlueWhite](h, ids, func(id DatapointID, dpt knxtypes.DPTRedGreenBlueWhite) error {
		rgbw := RGBW{Red: dpt.Red, Green: dpt.Green, Blue: dpt.Blue, White: dpt.White, Flags: dpt.Flags}
		return f(id, rgbw)
	})
}

func (h Handlers) OnIndication(id DatapointID, data []byte) error {
	handler, ok := h[id]
	if ok {
		return handler(id, data)
	}
	return nil
}

func (h Handlers) apply(w *Worker) {
	w.handler = h
}

// WithHandlers returns a WorkerOption that installs the given handlers.
func WithHandlers(h Handlers) WorkerOption {
	return WorkerOption(h)
}

func register[G any, W dptWrapper[G]](h Handlers, ids []DatapointID, f func(id DatapointID, dpt G) error) {
	for _, datapointID := range ids {
		h[datapointID] = makeHandler[G, W](f)
	}
}

func handle[G any, W dptWrapper[G]](id DatapointID, data []byte, f func(id DatapointID, dpt G) error) error {
	var dpt G
	if err := W(&dpt).UnmarshalBinary(data); err != nil {
		return err
	}
	return f(id, dpt)
}

func makeHandler[G any, W dptWrapper[G]](f func(id DatapointID, dpt G) error) ValueIndicationFunc {
	return func(id DatapointID, data []byte) error {
		return handle[G, W](id, data, f)
	}
}
