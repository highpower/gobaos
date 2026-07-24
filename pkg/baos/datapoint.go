package baos

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// DatapointID identifies a single BAOS datapoint.
type DatapointID uint16

// State represents the state bits of a datapoint value.
type State byte

// Datapoint holds a single datapoint value with its state.
type Datapoint struct {
	ID    DatapointID
	State State
	Data  []byte
}

// DatapointCommand holds a command to be sent to a datapoint.
type DatapointCommand struct {
	ID      DatapointID
	Command Command
	Data    []byte
}

type datapointData interface {
	data() (id DatapointID, extra byte, data []byte)
	attachData(id DatapointID, extra byte, data []byte)
}

type datapointWrapper[G any] interface {
	*G
	datapointData
}

const datapointHeaderSize = 4

// Equal reports whether d and other represent the same datapoint value.
func (d *Datapoint) Equal(other *Datapoint) bool {
	return d.ID == other.ID && d.State == other.State && bytes.Equal(d.Data, other.Data)
}

// String returns a human-readable representation of the datapoint.
func (d *Datapoint) String() string {
	return fmt.Sprintf("datapointValue[id=%d,state=%d,data=%x]", d.ID, d.State, d.Data)
}

func (d *Datapoint) data() (id DatapointID, extra byte, data []byte) {
	return d.ID, byte(d.State), d.Data
}

func (d *Datapoint) attachData(id DatapointID, extra byte, data []byte) {
	d.ID = id
	d.State = State(extra)
	d.Data = data
}

// Equal reports whether c and other represent the same datapoint command.
func (c *DatapointCommand) Equal(other *DatapointCommand) bool {
	return c.ID == other.ID && c.Command == other.Command && bytes.Equal(c.Data, other.Data)
}

// String returns a human-readable representation of the datapoint command.
func (c *DatapointCommand) String() string {
	return fmt.Sprintf("datapointCommand[id=%d,command=%d,data=%x]", c.ID, c.Command, c.Data)
}

func (c *DatapointCommand) data() (id DatapointID, extra byte, data []byte) {
	return c.ID, byte(c.Command), c.Data
}

func (c *DatapointCommand) attachData(id DatapointID, extra byte, data []byte) {
	c.ID = id
	c.Command = Command(extra)
	c.Data = data
}

func datapointsSizeInBytes[G any, W datapointWrapper[G]](values []G) int {
	totalSize := 0
	for i := range values {
		_, _, data := W(&values[i]).data()
		totalSize += datapointHeaderSize + len(data)
	}
	return totalSize
}

func marshalDatapointsTo[G any, W datapointWrapper[G]](data []byte, datapoints []G) error {
	for i := range datapoints {
		id, extra, d := W(&datapoints[i]).data()
		if len(data) < len(d)+4 {
			return fmt.Errorf("%w buffer has no capacity: %d while %d needed", ErrInvalid, len(data), len(d)+4)
		}
		binary.BigEndian.PutUint16(data[0:2], uint16(id))
		data[2] = extra
		data[3] = byte(len(d))
		copy(data[datapointHeaderSize:datapointHeaderSize+len(d)], d)
		data = data[datapointHeaderSize+len(d):]
	}
	if len(data) != 0 {
		return fmt.Errorf("%w buffer has %d extra bytes", ErrInvalid, len(data))
	}
	return nil
}

func unmarshalDatapointsFrom[G any, W datapointWrapper[G]](data []byte, n int) ([]G, error) {
	result := make([]G, 0, n)
	for i := 0; i < n; i++ {
		if len(data) < datapointHeaderSize {
			return nil, fmt.Errorf("%w buffer has only %d bytes while %d needed", ErrInvalid,
				len(data), datapointHeaderSize)
		}
		id := DatapointID(binary.BigEndian.Uint16(data[0:2]))
		state := data[2]
		length := int(data[3])
		if len(data) < datapointHeaderSize+length {
			return nil, fmt.Errorf("%w buffer has only %d bytes while %d needed", ErrInvalid,
				len(data), datapointHeaderSize+length)
		}
		var item G
		W(&item).attachData(id, state, data[datapointHeaderSize:datapointHeaderSize+length])
		result = append(result, item)
		data = data[datapointHeaderSize+length:]
	}
	if len(data) != 0 {
		return nil, fmt.Errorf("%w buffer has %d extra bytes", ErrInvalid, len(data))
	}
	return result, nil
}
