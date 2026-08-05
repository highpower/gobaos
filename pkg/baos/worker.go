package baos

import (
	"encoding"
	"encoding/binary"
	"fmt"
	"gobaos/pkg/ft12"
	"gobaos/pkg/knxtypes"
)

// Worker handles BAOS communication over an FT1.2 connection.
type Worker struct {
	conn    *ft12.Conn
	current *call
	buffer  []byte
	stop    chan bool
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

// Run starts the worker's event loop. It blocks until Stop is called or an error occurs.
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

// Stop signals the worker to shut down.
func (w *Worker) Stop() {
	w.stop <- true
}

func (w *Worker) ServerItem(id ItemID) (ServerItem, error) {
	req := GetServerItemReq{StartItem: id, NumberOfItems: 1}
	res := GetServerItemRes{}

	current := w.makeCall(&req, &res)
	w.calls <- current
	current.waitForReady()
	w.releaseCall(current)

	switch {
	case current.err != nil:
		return ServerItem{}, current.err
	case res.ErrorCode != 0:
		return ServerItem{}, Error{res.ErrorCode}
	case len(res.Items) == 0:
		return ServerItem{}, fmt.Errorf("%w response: 0 values for item %d", ErrInvalid, id)
	}
	return res.Items[0], nil
}

// Get reads the current value of the given datapoint into dpt.
func (w *Worker) Get(id DatapointID, dpt knxtypes.DPT) error {
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

// Set writes dpt to the given datapoint and sends it on the bus.
func (w *Worker) Set(id DatapointID, dpt knxtypes.DPT) error {
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

// Description reads datapoint description of  the given datapoint
func (w *Worker) Description(id DatapointID) (DatapointDescription, error) {
	result, err := w.descriptions(id, 1)
	switch {
	case err != nil:
		return DatapointDescription{}, err
	case len(result) == 0:
		return DatapointDescription{}, fmt.Errorf("%w response: 0 values for datapoint %d", ErrInvalid, id)
	}
	return result[0], nil
}

func (w *Worker) AllDescriptions() ([]DatapointDescription, error) {
	item, err := w.ServerItem(ItemIDConfiguredDatapoints)
	if err != nil {
		return nil, err
	}
	//TODO: add fetching functions instead of direct read
	maxId := binary.BigEndian.Uint16(item.Data)
	return w.descriptions(1, maxId)
}

func (w *Worker) DescriptionString(id DatapointID) (string, error) {
	req := GetDescriptionStringReq{StartDatapoint: id, NumberOfStrings: 1}
	res := GetDescriptionStringRes{}

	current := w.makeCall(&req, &res)
	w.calls <- current
	current.waitForReady()
	w.releaseCall(current)

	switch {
	case current.err != nil:
		return "", current.err
	case res.ErrorCode != 0:
		return "", Error{res.ErrorCode}
	case len(res.Strings) == 0:
		return "", fmt.Errorf("%w response: 0 values for datapoint %d", ErrInvalid, id)
	}
	return res.Strings[0], nil
}

func (w *Worker) descriptions(start DatapointID, number uint16) ([]DatapointDescription, error) {
	result := []DatapointDescription(nil)
	descInResponse := (ft12.MaxMessageLength - getDescriptionResHeaderSize) / datapointDescriptionSize
	i := int(number)
	for i > 0 {
		req := GetDatapointDescriptionReq{StartDatapoint: start, NumberOfDatapoints: uint16(min(descInResponse, i))}
		res := GetDatapointDescriptionRes{}

		current := w.makeCall(&req, &res)
		w.calls <- current
		current.waitForReady()
		w.releaseCall(current)

		switch {
		case current.err != nil:
			return nil, current.err
		case res.ErrorCode == NoElementFound:
			return result, nil
		case res.ErrorCode != 0:
			return nil, Error{Code: res.ErrorCode}
		}
		result = append(result, res.Descriptions...)
		i -= len(res.Descriptions)
		start = nextDatapointID(res.Descriptions)
	}
	return result, nil
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
		buffer:  nil,
		stop:    make(chan bool),
		calls:   make(chan *call),
	}
	for _, option := range options {
		option.apply(&result)
	}
	return &result
}

func nextDatapointID(desc []DatapointDescription) DatapointID {
	result := DatapointID(0)
	for _, d := range desc {
		if d.ID > result {
			result = d.ID
		}
	}
	return result + 1
}
