package baos

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"slices"
)

type DatapointValueInd struct {
	StartDatapoint DatapointID
	Values         []Datapoint
}

func (r *DatapointValueInd) Service() byte {
	return MainService
}

func (r *DatapointValueInd) SubService() byte {
	return SubServiceDatapointValueInd
}

func (r *DatapointValueInd) Valid() error {
	if len(r.Values) == 0 {
		return fmt.Errorf("%w indication: it should contain at least 1 datapoint", ErrInvalid)
	}
	return nil
}

func (r *DatapointValueInd) Equal(other *DatapointValueInd) bool {
	return r.StartDatapoint == other.StartDatapoint && slices.EqualFunc(r.Values, other.Values, func(v1, v2 Datapoint) bool {
		return v1.Equal(&v2)
	})
}

func (r *DatapointValueInd) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString(fmt.Sprintf("datapointValueInd[startDatapoint=%d,values=[", r.StartDatapoint))
	for i, val := range r.Values {
		buffer.WriteString(val.String())
		if i < len(r.Values)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("]]")
	return buffer.String()
}

func (r *DatapointValueInd) MarshalBinary() ([]byte, error) {
	result := make([]byte, 6+datapointsSizeInBytes(r.Values))
	result[0] = MainService
	result[1] = SubServiceDatapointValueInd
	binary.BigEndian.PutUint16(result[2:4], uint16(r.StartDatapoint))
	binary.BigEndian.PutUint16(result[4:6], uint16(len(r.Values)))
	if err := marshalDatapointsTo(result[6:], r.Values); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *DatapointValueInd) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) < 6:
		return fmt.Errorf("%w data length %d, expected at least %d", ErrInvalid, len(data), 6)
	case data[0] != byte(MainService):
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], MainService)
	case data[1] != byte(SubServiceDatapointValueInd):
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], SubServiceDatapointValueInd)
	}
	startDatapoint := DatapointID(binary.BigEndian.Uint16(data[2:4]))
	num := int(binary.BigEndian.Uint16(data[4:6]))
	values, err := unmarshalDatapointsFrom[Datapoint](data[6:], num)
	if err != nil {
		return err
	}
	*r = DatapointValueInd{startDatapoint, values}
	return nil
}

func IsDatapointValueInd(data []byte) bool {
	return len(data) > 6 && data[0] == MainService && data[1] == SubServiceDatapointValueInd
}
