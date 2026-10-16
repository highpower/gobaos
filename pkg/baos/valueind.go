package baos

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"slices"
)

// datapointValueInd is an unsolicited indication carrying datapoint values.
type datapointValueInd struct {
	StartDatapoint DatapointID
	Values         []Datapoint
}

// Service returns the main service identifier.
func (r *datapointValueInd) Service() byte {
	return MainService
}

// SubService returns the sub-service identifier for the indication.
func (r *datapointValueInd) SubService() byte {
	return SubServiceDatapointValueInd
}

// Equal reports whether r and other are identical indications.
func (r *datapointValueInd) Equal(other *datapointValueInd) bool {
	return r.StartDatapoint == other.StartDatapoint && slices.EqualFunc(r.Values, other.Values, func(v1, v2 Datapoint) bool {
		return v1.Equal(&v2)
	})
}

// String returns a human-readable representation of the indication.
func (r *datapointValueInd) String() string {
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

// MarshalBinary encodes the indication into its BAOS wire format.
func (r *datapointValueInd) MarshalBinary() ([]byte, error) {
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

// UnmarshalBinary decodes the BAOS wire format into the indication.
func (r *datapointValueInd) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) < 6:
		return fmt.Errorf("%w data length %d, expected at least %d", ErrInvalid, len(data), 6)
	case data[0] != MainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], MainService)
	case data[1] != SubServiceDatapointValueInd:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], SubServiceDatapointValueInd)
	}
	startDatapoint := DatapointID(binary.BigEndian.Uint16(data[2:4]))
	num := int(binary.BigEndian.Uint16(data[4:6]))
	values, err := unmarshalDatapointsFrom[Datapoint](data[6:], num)
	if err != nil {
		return err
	}
	*r = datapointValueInd{startDatapoint, values}
	return nil
}

// isDatapointValueInd reports whether data looks like a datapointValueInd message.
func isDatapointValueInd(data []byte) bool {
	return len(data) > 6 && data[0] == MainService && data[1] == SubServiceDatapointValueInd
}
