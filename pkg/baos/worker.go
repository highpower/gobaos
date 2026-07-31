package baos

import (
	"encoding"
	"fmt"
	"gobaos/pkg/ft12"
	"gobaos/pkg/knxtypes"
)

// Worker handles BAOS communication over an FT1.2 connection.
type Worker struct {
	conn    *ft12.Conn
	current *call
	stop    chan bool
	buffer  []byte
	calls   chan *call
}

type WorkerOption interface {
	apply(w *Worker)
}

type call struct {
	request  encoding.BinaryMarshaler
	response encoding.BinaryUnmarshaler
	err      error
	ready    chan bool
}

func (w *Worker) Run() error {
	for {
		w.buffer = w.buffer[:0]
		select {
		case <-w.stop:
			return nil
		case c := <-w.calls:
			if err := w.handleCall(c); err != nil {
				return err
			}
		default:
			if err := w.handleIndication(); err != nil {
				return err
			}
		}
	}
}

func (w *Worker) Stop() {
	w.stop <- true
}

func (w *Worker) Bool(id DatapointID) (bool, error) {
	dpt := knxtypes.DPTBool{}
	if err := w.GetValue(id, &dpt); err != nil {
		return false, err
	}
	return dpt.Value, nil
}

func (w *Worker) SetBool(id DatapointID, value bool) error {
	dpt := knxtypes.DPTBool{Value: value}
	return w.SetValue(id, &dpt)
}

func (w *Worker) handleCall(current *call) error {
	data, err := current.request.MarshalBinary()
	if err != nil {
		return err
	}
	if _, err := w.conn.Write(data); err != nil {
		return err
	}
	for w.conn.Available() == 0 {
		if err := w.conn.Poll(); err != nil {
			return err
		}
	}
	if cap(w.buffer) < w.conn.Available() {
		w.buffer = make([]byte, 0, w.conn.Available())
	}
	w.buffer = w.buffer[:w.conn.Available()]
	if _, err := w.conn.Read(w.buffer); err != nil {
		return err
	}
	if err := current.response.UnmarshalBinary(w.buffer); err != nil {
		current.err = err
	}
	current.ready <- true
	return nil
}

func (w *Worker) handleIndication() error {
	if err := w.conn.Poll(); err != nil {
		return err
	}
	if cap(w.buffer) < w.conn.Available() {
		w.buffer = make([]byte, 0, w.conn.Available())
	}
	w.buffer = w.buffer[:w.conn.Available()]
	if _, err := w.conn.Read(w.buffer); err != nil {
		return err
	}
	if !IsDatapointValueInd(w.buffer) {
		return nil
	}
	indication := DatapointValueInd{}
	if err := indication.UnmarshalBinary(w.buffer); err != nil {
		return err
	}
	return nil
}

func (w *Worker) Descriptions(start DatapointID, number uint16) ([]DatapointDescription, error) {
	req := GetDatapointDescriptionReq{StartDatapoint: start, NumberOfDatapoints: uint16(number)}
	res := GetDatapointDescriptionRes{}

	current := w.makeCall(&req, &res)
	w.calls <- current
	current.waitForReady()
	w.releaseCall(current)

	switch {
	case current.err != nil:
		return nil, current.err
	case res.ErrorCode != 0:
		return nil, Error{Code: res.ErrorCode}
	}
	return res.Descriptions, nil
}

func (w *Worker) GetValue(id DatapointID, dpt knxtypes.DPT) error {
	req := GetDatapointValueReq{StartDatapoint: id, NumberOfDatapoints: 1}
	res := GetDatapointValueRes{}

	current := w.makeCall(&req, &res)
	w.calls <- current
	current.waitForReady()
	w.releaseCall(current)

	switch {
	case current.err != nil:
		return current.err
	case res.ErrorCode != 0:
		return Error{res.ErrorCode}
	case len(res.Values) == 0:
		return fmt.Errorf("%w response: 0 values for datapoint %d", ErrInvalid, id)
	}
	if err := dpt.UnmarshalBinary(res.Values[0].Data); err != nil {
		return err
	}
	return nil
}

func (w *Worker) SetValue(id DatapointID, dpt knxtypes.DPT) error {
	data, err := dpt.MarshalBinary()
	if err != nil {
		return err
	}
	req := SetDatapointValueReq{Commands: []DatapointCommand{{ID: id, Command: SetAndSendOnBus, Data: data}}}
	res := SetDatapointValueRes{ErrorCode: ErrorCode(0xff)}

	current := w.makeCall(&req, &res)
	w.calls <- current
	current.waitForReady()
	w.releaseCall(current)

	switch {
	case current.err != nil:
		return current.err
	case res.ErrorCode != 0:
		return Error{res.ErrorCode}
	}
	return nil
}

func (w *Worker) makeCall(req encoding.BinaryMarshaler, resp encoding.BinaryUnmarshaler) *call {
	result := call{request: req, response: resp, ready: make(chan bool, 1)}
	return &result
}

func (w *Worker) releaseCall(call *call) {
	close(call.ready)
}

func (c *call) waitForReady() {
	<-c.ready
}

// NewWorker creates a Worker backed by the given FT1.2 connection.
func NewWorker(conn *ft12.Conn, options ...WorkerOption) *Worker {
	result := Worker{
		conn:    conn,
		current: nil,
		stop:    make(chan bool),
		buffer:  nil,
		calls:   make(chan *call),
	}
	for _, option := range options {
		option.apply(&result)
	}
	return &result
}
