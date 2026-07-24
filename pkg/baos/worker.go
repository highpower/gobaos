package baos

import (
	"gobaos/pkg/ft12"
	"gobaos/pkg/knxtypes"
)

type dataFunc func(id DatapointID, data []byte) error
type indicationHandlerMap map[DatapointID]dataFunc

// BoolValue pairs a datapoint ID with a boolean value.
type BoolValue struct {
	ID    DatapointID
	Value bool
}

// ControlValue pairs a datapoint ID with a control and value.
type ControlValue struct {
	ID      DatapointID
	Control bool
	Value   bool
}

// Worker handles BAOS communication over an FT1.2 connection.
type Worker struct {
	conn     *ft12.Conn
	handlers indicationHandlerMap
}

type dptWrapper[G any] interface {
	*G
	knxtypes.DPT
}

// SetBool sends boolean values to the configured datapoints.
func (w *Worker) SetBool(values []BoolValue) error {
	//TODO: implement
	return nil
}

// SetControl sends control values to the configured datapoints.
func (w *Worker) SetControl(values []ControlValue) error {
	//TODO: implement
	return nil
}

// OnBoolIndication registers a handler for boolean indications on the given datapoint IDs.
func (w *Worker) OnBoolIndication(ids []DatapointID, f func(id DatapointID, val knxtypes.DPTBool) error) {
	for _, id := range ids {
		w.handlers[id] = makeHandler[knxtypes.DPTBool](f)
	}
}

// OnControlIndication registers a handler for control indications on the given datapoint IDs.
func (w *Worker) OnControlIndication(ids []DatapointID, f func(id DatapointID, val knxtypes.DPTControl) error) {
	for _, id := range ids {
		w.handlers[id] = makeHandler[knxtypes.DPTControl](f)
	}
}

// OnStepControlIndication registers a handler for step-control indications on the given datapoint IDs.
func (w *Worker) OnStepControlIndication(ids []DatapointID, f func(id DatapointID, val knxtypes.DPTStepControl) error) {
	for _, id := range ids {
		w.handlers[id] = makeHandler[knxtypes.DPTStepControl](f)
	}
}

// NewWorker creates a Worker backed by the given FT1.2 connection.
func NewWorker(conn *ft12.Conn) Worker {
	result := Worker{conn: conn, handlers: make(map[DatapointID]dataFunc)}
	return result
}

func handle[G any, W dptWrapper[G]](id DatapointID, data []byte, f func(id DatapointID, dpt G) error) error {
	var dpt G
	if err := W(&dpt).UnmarshalBinary(data); err != nil {
		return err
	}
	return f(id, dpt)
}

func makeHandler[G any, W dptWrapper[G]](f func(id DatapointID, dpt G) error) dataFunc {
	return func(id DatapointID, data []byte) error {
		return handle[G, W](id, data, f)
	}
}
