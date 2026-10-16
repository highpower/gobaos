package baos

import (
	"encoding/binary"
	"fmt"
	"slices"
)

// datapointValueInd is an unsolicited indication carrying datapoint values.
type datapointValueInd struct {
	startDatapoint DatapointID
	values         []Datapoint
}

// MarshalBinary encodes the indication into its BAOS wire format.
func (r *datapointValueInd) MarshalBinary() ([]byte, error) {
	result := make([]byte, 6+datapointsSizeInBytes(r.values))
	result[0] = mainService
	result[1] = subServiceDatapointValueInd
	binary.BigEndian.PutUint16(result[2:4], uint16(r.startDatapoint))
	binary.BigEndian.PutUint16(result[4:6], uint16(len(r.values)))
	if err := marshalDatapointsTo(result[6:], r.values); err != nil {
		return nil, err
	}
	return result, nil
}

// UnmarshalBinary decodes the BAOS wire format into the indication.
func (r *datapointValueInd) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) < 6:
		return fmt.Errorf("%w data length %d, expected at least %d", ErrInvalid, len(data), 6)
	case data[0] != mainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], mainService)
	case data[1] != subServiceDatapointValueInd:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], subServiceDatapointValueInd)
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

// equal reports whether r and other are identical indications.
func (r *datapointValueInd) equal(other *datapointValueInd) bool {
	return r.startDatapoint == other.startDatapoint && slices.EqualFunc(r.values, other.values, func(v1, v2 Datapoint) bool {
		return v1.Equal(&v2)
	})
}

// isDatapointValueInd reports whether data looks like a datapointValueInd message.
func isDatapointValueInd(data []byte) bool {
	return len(data) > 6 && data[0] == mainService && data[1] == subServiceDatapointValueInd
}
