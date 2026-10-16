package baos

import (
	"encoding/binary"
	"fmt"
	"slices"
)

// setDatapointValueReq is a request to write datapoint values.
type setDatapointValueReq struct {
	commands []DatapointCommand
}

// setDatapointValueRes is a response confirming the set operation or reporting an error.
type setDatapointValueRes struct {
	startDatapoint     DatapointID
	numberOfDatapoints uint16
	errorCode          ErrorCode
}

// MarshalBinary encodes the request into its BAOS wire format.
func (r *setDatapointValueReq) MarshalBinary() ([]byte, error) {
	result := make([]byte, 6+datapointsSizeInBytes(r.commands))
	result[0] = mainService
	result[1] = subServiceSetDatapointValueReq
	startDatapoint := DatapointID(0)
	if len(r.commands) > 0 {
		startDatapoint = r.commands[0].ID
	}
	binary.BigEndian.PutUint16(result[2:4], uint16(startDatapoint))
	binary.BigEndian.PutUint16(result[4:6], uint16(len(r.commands)))
	if err := marshalDatapointsTo(result[6:], r.commands); err != nil {
		return nil, err
	}
	return result, nil
}

// UnmarshalBinary decodes the BAOS wire format into the request.
func (r *setDatapointValueReq) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) < 6:
		return fmt.Errorf("%w data length %d, expected at least %d", ErrInvalid, len(data), 6)
	case data[0] != mainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], mainService)
	case data[1] != subServiceSetDatapointValueReq:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], subServiceSetDatapointValueReq)
	}
	num := int(binary.BigEndian.Uint16(data[4:6]))
	values, err := unmarshalDatapointsFrom[DatapointCommand](data[6:], num)
	if err != nil {
		return err
	}
	r.commands = values
	return nil
}

// equal reports whether r and other are identical requests.
func (r *setDatapointValueReq) equal(other *setDatapointValueReq) bool {
	return slices.EqualFunc(r.commands, other.commands, func(v1, v2 DatapointCommand) bool {
		return v1.Equal(&v2)
	})
}

// MarshalBinary encodes the response into its BAOS wire format.
func (r *setDatapointValueRes) MarshalBinary() ([]byte, error) {
	result := make([]byte, 7)
	result[0] = mainService
	result[1] = subServiceSetDatapointValueRes
	binary.BigEndian.PutUint16(result[2:4], uint16(r.startDatapoint))
	binary.BigEndian.PutUint16(result[4:6], r.numberOfDatapoints)
	result[6] = byte(r.errorCode)
	return result, nil
}

// UnmarshalBinary decodes the BAOS wire format into the response.
func (r *setDatapointValueRes) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) != 7:
		return fmt.Errorf("%w data length %d, expected to be %d", ErrInvalid, len(data), 7)
	case data[0] != mainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], mainService)
	case data[1] != subServiceSetDatapointValueRes:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], subServiceSetDatapointValueRes)
	}
	r.startDatapoint = DatapointID(binary.BigEndian.Uint16(data[2:4]))
	r.numberOfDatapoints = binary.BigEndian.Uint16(data[4:6])
	r.errorCode = ErrorCode(data[6])
	return nil
}

// equal reports whether r and other are identical responses.
func (r *setDatapointValueRes) equal(other *setDatapointValueRes) bool {
	return r.startDatapoint == other.startDatapoint &&
		r.numberOfDatapoints == other.numberOfDatapoints &&
		byte(r.errorCode) == byte(other.errorCode)
}
