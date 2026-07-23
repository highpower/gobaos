package baos

import (
	"gobaos/pkg/ft12"
	"gobaos/pkg/knxtypes"
)

type dataFunc func(id DatapointID, data []byte) error
type indicationHandlerMap map[DatapointID]dataFunc

type BoolValue struct {
	ID    DatapointID
	Value bool
}

type Worker struct {
	// conn *ft12.Conn
	ihm indicationHandlerMap
}

type dptWrapper[G any] interface {
	*G
	knxtypes.DPT
}

func (w *Worker) SetBool(values []BoolValue) {

}

func (w *Worker) OnBoolIndication(ids []DatapointID, f func(id DatapointID, val knxtypes.DPTBool) error) {
	for _, id := range ids {
		w.ihm[id] = makeHandler[knxtypes.DPTBool](f)
	}
}

func (w *Worker) OnControlIndication(ids []DatapointID, f func(id DatapointID, val knxtypes.DPTControl) error) {
	for _, id := range ids {
		w.ihm[id] = makeHandler[knxtypes.DPTControl](f)
	}
}

func (w *Worker) OnStepControlIndication(ids []DatapointID, f func(id DatapointID, val knxtypes.DPTStepControl) error) {
	for _, id := range ids {
		w.ihm[id] = makeHandler[knxtypes.DPTStepControl](f)
	}
}

func NewWorker(conn *ft12.Conn) Worker {
	return Worker{}
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
