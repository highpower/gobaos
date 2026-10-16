package baos

import (
	"encoding/binary"
	"fmt"
	"slices"
)

// Filter controls which datapoints are returned in a get-value request.
type Filter byte

// getDatapointValueReq is a request to read datapoint values.
type getDatapointValueReq struct {
	startDatapoint     DatapointID
	numberOfDatapoints uint16
	filter             Filter
}

// getDatapointValueRes is a response containing datapoint values or an error.
type getDatapointValueRes struct {
	startDatapoint DatapointID
	values         []Datapoint
	errorCode      ErrorCode
}

// Supported datapoint filters.
const (
	AllDatapoints     Filter = 0
	ValidDatapoints   Filter = 1
	UpdatedDatapoints Filter = 2
)

// MarshalBinary encodes the request into its BAOS wire format.
func (r *getDatapointValueReq) MarshalBinary() ([]byte, error) {
	b := make([]byte, 7)
	b[0] = mainService
	b[1] = subServiceGetDatapointValueReq
	binary.BigEndian.PutUint16(b[2:4], uint16(r.startDatapoint))
	binary.BigEndian.PutUint16(b[4:6], r.numberOfDatapoints)
	b[6] = byte(r.filter)
	return b, nil
}

// UnmarshalBinary decodes the BAOS wire format into the request.
func (r *getDatapointValueReq) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) != 7:
		return fmt.Errorf("%w data length %d, expected %d", ErrInvalid, len(data), 7)
	case data[0] != mainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], mainService)
	case data[1] != subServiceGetDatapointValueReq:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], subServiceGetDatapointValueReq)
	}
	r.startDatapoint = DatapointID(binary.BigEndian.Uint16(data[2:4]))
	r.numberOfDatapoints = binary.BigEndian.Uint16(data[4:6])
	r.filter = Filter(data[6])
	return nil
}

// equal reports whether r and other are identical requests.
func (r *getDatapointValueReq) equal(other *getDatapointValueReq) bool {
	return r.startDatapoint == other.startDatapoint && r.numberOfDatapoints == other.numberOfDatapoints &&
		r.filter == other.filter
}

// MarshalBinary encodes the response into its BAOS wire format.
func (r *getDatapointValueRes) MarshalBinary() ([]byte, error) {
	if r.errorCode != 0 {
		result := make([]byte, 7)
		result[0] = mainService
		result[1] = subServiceGetDatapointValueRes
		binary.BigEndian.PutUint16(result[2:4], uint16(r.startDatapoint))
		binary.BigEndian.PutUint16(result[4:6], uint16(0))
		result[6] = byte(r.errorCode)
		return result, nil
	}
	result := make([]byte, 6+datapointsSizeInBytes(r.values))
	result[0] = mainService
	result[1] = subServiceGetDatapointValueRes
	binary.BigEndian.PutUint16(result[2:4], uint16(r.startDatapoint))
	binary.BigEndian.PutUint16(result[4:6], uint16(len(r.values)))
	if err := marshalDatapointsTo(result[6:], r.values); err != nil {
		return nil, err
	}
	return result, nil
}

// UnmarshalBinary decodes the BAOS wire format into the response.
func (r *getDatapointValueRes) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) < 6:
		return fmt.Errorf("%w data length %d, expected at least %d", ErrInvalid, len(data), 6)
	case data[0] != mainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], mainService)
	case data[1] != subServiceGetDatapointValueRes:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], subServiceGetDatapointValueRes)
	}
	startDatapoint := DatapointID(binary.BigEndian.Uint16(data[2:4]))
	if len(data) == 7 {
		*r = getDatapointValueRes{startDatapoint: startDatapoint, errorCode: ErrorCode(data[6])}
		return nil
	}
	num := int(binary.BigEndian.Uint16(data[4:6]))
	values, err := unmarshalDatapointsFrom[Datapoint](data[6:], num)
	if err != nil {
		return err
	}
	*r = getDatapointValueRes{startDatapoint: startDatapoint, values: values}
	return nil
}

// equal reports whether r and other are identical responses.
func (r *getDatapointValueRes) equal(other *getDatapointValueRes) bool {
	if r.startDatapoint != other.startDatapoint || byte(r.errorCode) != byte(other.errorCode) {
		return false
	}
	return slices.EqualFunc(r.values, other.values, func(v1, v2 Datapoint) bool {
		return v1.Equal(&v2)
	})
}
