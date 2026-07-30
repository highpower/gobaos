package baos

import (
	"gobaos/pkg/ft12"
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

// NewWorker creates a Worker backed by the given FT1.2 connection.
func NewWorker(conn *ft12.Conn) Worker {
	result := Worker{conn: conn, handlers: make(map[DatapointID]dataFunc)}
	return result
}
