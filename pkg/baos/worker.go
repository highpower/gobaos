package baos

import (
	"context"
	"encoding"
	"encoding/binary"
	"fmt"
	"sync/atomic"

	"gobaos/pkg/ft12"
	"gobaos/pkg/knxtypes"
)

// Worker handles BAOS communication over an FT1.2 connection.
type Worker struct {
	conn        *ft12.Conn
	buffer      []byte
	errorOrStop chan errorOrStop
	indications chan datapointValueInd
	calls       chan *call
	handler     ValueIndicationHandler
	stopper     atomic.Value
	goroutines  int
}

// ValueIndicationFunc is called when an unsolicited datapoint value indication arrives.
type ValueIndicationFunc func(id DatapointID, data []byte) error

type ValueIndicationHandler interface {
	OnIndication(id DatapointID, data []byte) error
}

// WorkerOption configures a Worker during construction.
type WorkerOption interface {
	apply(w *Worker)
}

type call struct {
	ctx      context.Context
	request  encoding.BinaryMarshaler
	response encoding.BinaryUnmarshaler
	err      error
	ready    chan bool
}

type errorOrStop struct {
	stop bool
	err  error
}

type defaultHandler struct{}

type goroutineCountOption int

type valueIndicationOption struct {
	f ValueIndicationFunc
}

type stopperFunc func()

const defaultIndicationHandlerCount = 5

// Run starts the worker's event loop. It blocks until Stop is called or an error occurs.
func (w *Worker) Run() error {
	defer close(w.errorOrStop)
	defer close(w.indications)
	defer close(w.calls)
	defer w.stopper.Store(stopperFunc(w.alreadyStopped))
	for {
		w.buffer = w.buffer[:0]
		select {
		case eos := <-w.errorOrStop:
			return eos.err
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
	stopper := w.stopper.Load().(stopperFunc)
	stopper()
}

// Get reads the current value of the given datapoint into dpt.
func (w *Worker) Get(ctx context.Context, id DatapointID, dpt knxtypes.DPT) error {
	return w.GetWithFilter(ctx, id, AllDatapoints, dpt)
}

func (w *Worker) GetWithFilter(ctx context.Context, id DatapointID, f Filter, dpt knxtypes.DPT) error {
	req := getDatapointValueReq{startDatapoint: id, numberOfDatapoints: 1, filter: f}
	res := getDatapointValueRes{}

	current := w.makeCall(ctx, &req, &res)
	w.calls <- current
	current.waitForReady()
	w.releaseCall(current)

	switch {
	case ctx.Err() != nil:
		return fmt.Errorf("worker.get: %w", ctx.Err())
	case current.err != nil:
		return fmt.Errorf("worker.get: %w", current.err)
	case res.errorCode != 0:
		return fmt.Errorf("worker.get: %w", Error{res.errorCode})
	case len(res.values) == 0:
		return fmt.Errorf("worker.get: %w response [0 values for datapoint %d]", ErrInvalid, id)
	}
	if err := dpt.UnmarshalBinary(res.values[0].Data); err != nil {
		return err
	}
	return nil
}

// Set writes dpt to the given datapoint and sends it on the bus.
func (w *Worker) Set(ctx context.Context, id DatapointID, dpt knxtypes.DPT) error {
	if err := knxtypes.IsValidDPT(dpt); err != nil {
		return fmt.Errorf("can not set dpt: %w", err)
	}
	data, err := dpt.MarshalBinary()
	if err != nil {
		return err
	}
	req := setDatapointValueReq{commands: []DatapointCommand{{ID: id, Command: SetAndSendOnBus, Data: data}}}
	res := setDatapointValueRes{errorCode: ErrorCode(0xff)}

	current := w.makeCall(ctx, &req, &res)
	w.calls <- current
	current.waitForReady()
	w.releaseCall(current)

	switch {
	case ctx.Err() != nil:
		return fmt.Errorf("worker.set: %w", ctx.Err())
	case current.err != nil:
		return fmt.Errorf("worker.set: %w", current.err)
	case res.errorCode != 0:
		return fmt.Errorf("worker.set: %w", Error{res.errorCode})
	}
	return nil
}

// Description reads datapoint description of  the given datapoint
func (w *Worker) Description(ctx context.Context, id DatapointID) (DatapointDescription, error) {
	result, err := w.descriptions(ctx, id, 1)
	switch {
	case err != nil:
		return DatapointDescription{}, err
	case len(result) == 0:
		return DatapointDescription{},
			fmt.Errorf("worker.description: %w response [no description for datapoint %d]", ErrInvalid, id)
	}
	return result[0], nil
}

// AllDescriptions reads descriptions for all configured datapoints.
func (w *Worker) AllDescriptions(ctx context.Context) ([]DatapointDescription, error) {
	item, err := w.ServerItem(ctx, ItemIDConfiguredDatapoints)
	if err != nil {
		return nil, err
	}
	// TODO(o.obolenskiy): add fetching functions instead of direct read
	maxID := binary.BigEndian.Uint16(item.Data)
	return w.descriptions(ctx, 1, maxID)
}

// DescriptionString reads the description string of the given datapoint.
func (w *Worker) DescriptionString(ctx context.Context, id DatapointID) (string, error) {
	req := getDescriptionStringReq{startDatapoint: id, numberOfStrings: 1}
	res := getDescriptionStringRes{}

	current := w.makeCall(ctx, &req, &res)
	w.calls <- current
	current.waitForReady()
	w.releaseCall(current)

	switch {
	case current.err != nil:
		return "", fmt.Errorf("error in descriptionString request: %w", current.err)
	case res.errorCode != 0:
		return "", Error{res.errorCode}
	case len(res.values) == 0:
		return "", fmt.Errorf("%w response: 0 values for datapoint %d", ErrInvalid, id)
	}
	return res.values[0], nil
}

// ServerItem reads a single server item by its ID.
func (w *Worker) ServerItem(ctx context.Context, id ItemID) (ServerItem, error) {
	req := getServerItemReq{startItem: id, numberOfItems: 1}
	res := getServerItemRes{}

	current := w.makeCall(ctx, &req, &res)
	w.calls <- current
	current.waitForReady()
	w.releaseCall(current)

	switch {
	case current.err != nil:
		return ServerItem{}, fmt.Errorf("error in getServerItem request: %w", current.err)
	case res.errorCode != 0:
		return ServerItem{}, Error{res.errorCode}
	case len(res.items) == 0:
		return ServerItem{}, fmt.Errorf("%w response: 0 values for item %d", ErrInvalid, id)
	}
	return res.items[0], nil
}

// SetServerItem writes data to a single server item by its ID.
func (w *Worker) SetServerItem(ctx context.Context, id ItemID, data []byte) error {
	req := setServerItemReq{items: []ServerItem{{ID: id, Data: data}}}
	res := setServerItemRes{}

	current := w.makeCall(ctx, &req, &res)
	w.calls <- current
	current.waitForReady()
	w.releaseCall(current)

	switch {
	case current.err != nil:
		return fmt.Errorf("error in setServerItem request: %w", current.err)
	case res.errorCode != 0:
		return Error{res.errorCode}
	}
	return nil
}

func (w *Worker) descriptions(ctx context.Context, start DatapointID, number uint16) ([]DatapointDescription, error) {
	result := []DatapointDescription(nil)
	descInResponse := (ft12.MaxMessageLength - getDescriptionResHeaderSize) / datapointDescriptionSize
	i := int(number)
	for i > 0 {
		req := getDatapointDescriptionReq{startDatapoint: start, numberOfDatapoints: uint16(min(descInResponse, i))}
		res := getDatapointDescriptionRes{}

		current := w.makeCall(ctx, &req, &res)
		w.calls <- current
		current.waitForReady()
		w.releaseCall(current)

		switch {
		case current.err != nil:
			return nil, fmt.Errorf("error in getDatapointDescription request: %w", current.err)
		case res.errorCode == NoElementFound:
			return result, nil
		case res.errorCode != 0:
			return nil, Error{Code: res.errorCode}
		}
		result = append(result, res.descriptions...)
		i -= len(res.descriptions)
		start = nextDatapointID(res.descriptions)
	}
	return result, nil
}

func (w *Worker) handleCall(current *call) error {
	data, err := current.request.MarshalBinary()
	if err != nil {
		return fmt.Errorf("marshal request error: %w", err)
	}
	if _, err := w.conn.Write(data); err != nil {
		return fmt.Errorf("write request error: %w", err)
	}
	for w.conn.Available() == 0 {
		if err := w.conn.Poll(); err != nil {
			return fmt.Errorf("handleCall poll error: %w", err)
		}
	}
	if cap(w.buffer) < w.conn.Available() {
		w.buffer = make([]byte, 0, w.conn.Available())
	}
	w.buffer = w.buffer[:w.conn.Available()]
	if _, err := w.conn.Read(w.buffer); err != nil {
		return fmt.Errorf("read response error: %w", err)
	}
	if err := current.response.UnmarshalBinary(w.buffer); err != nil {
		current.err = err
	}
	current.ready <- true
	return nil
}

func (w *Worker) handleIndication() error {
	if err := w.conn.Poll(); err != nil {
		return fmt.Errorf("handleIndication poll error: %w", err)
	}
	if cap(w.buffer) < w.conn.Available() {
		w.buffer = make([]byte, 0, w.conn.Available())
	}
	w.buffer = w.buffer[:w.conn.Available()]
	if _, err := w.conn.Read(w.buffer); err != nil {
		return fmt.Errorf("read indication error: %w", err)
	}
	if !isDatapointValueInd(w.buffer) {
		return nil
	}
	indication := datapointValueInd{}
	if err := indication.UnmarshalBinary(w.buffer); err != nil {
		return fmt.Errorf("unmarshal indication error: %w", err)
	}
	w.indications <- indication
	return nil
}

func (w *Worker) runIndication(indications <-chan datapointValueInd) {
	for indication := range indications {
		for _, value := range indication.values {
			if err := w.handler.OnIndication(value.ID, value.Data); err != nil {
				w.errorOrStop <- errorOrStop{err: err}
			}
		}
	}
}

func (w *Worker) makeCall(ctx context.Context, req encoding.BinaryMarshaler, resp encoding.BinaryUnmarshaler) *call {
	result := call{ctx: ctx, request: req, response: resp, ready: make(chan bool, 1)}
	return &result
}

func (w *Worker) releaseCall(call *call) {
	close(call.ready)
}

func (w *Worker) stop() {
	w.errorOrStop <- errorOrStop{stop: true}
}

func (w *Worker) alreadyStopped() {
}

func (c *call) waitForReady() {
	select {
	case <-c.ready:
	case <-c.ctx.Done():
	}
}

func (d defaultHandler) OnIndication(_ DatapointID, _ []byte) error {
	return nil
}

func (g goroutineCountOption) apply(w *Worker) {
	w.goroutines = int(g)
}

func (o valueIndicationOption) apply(w *Worker) {
	w.handler = ValueIndicationHandler(o)
}

func (o valueIndicationOption) OnIndication(id DatapointID, data []byte) error {
	return o.f(id, data)
}

// NewWorker creates a Worker backed by the given FT1.2 connection.
func NewWorker(conn *ft12.Conn, options ...WorkerOption) *Worker {

	result := Worker{
		conn:        conn,
		buffer:      nil,
		errorOrStop: make(chan errorOrStop),
		indications: make(chan datapointValueInd),
		calls:       make(chan *call),
		handler:     defaultHandler{},
		goroutines:  defaultIndicationHandlerCount,
	}
	result.stopper.Store(stopperFunc(result.stop))
	for _, option := range options {
		option.apply(&result)
	}
	for i := 0; i < result.goroutines; i++ {
		go result.runIndication(result.indications)
	}
	return &result
}

// WithHandler returns a WorkerOption that registers a raw value indication handler.
func WithHandler(f ValueIndicationFunc) WorkerOption {
	return valueIndicationOption{f: f}
}

func WithGoroutineCount(n int) WorkerOption {
	return goroutineCountOption(n)
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
