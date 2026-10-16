package baos

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"slices"
)

// setDatapointValueReq is a request to write datapoint values.
type setDatapointValueReq struct {
	Commands []DatapointCommand
}

// setDatapointValueRes is a response confirming the set operation or reporting an error.
type setDatapointValueRes struct {
	StartDatapoint     DatapointID
	NumberOfDatapoints uint16
	ErrorCode          ErrorCode
}

// Service returns the main service identifier.
func (r *setDatapointValueReq) Service() byte {
	return MainService
}

// SubService returns the sub-service identifier for the request.
func (r *setDatapointValueReq) SubService() byte {
	return SubServiceSetDatapointValueReq
}

// Equal reports whether r and other are identical requests.
func (r *setDatapointValueReq) Equal(other *setDatapointValueReq) bool {
	return slices.EqualFunc(r.Commands, other.Commands, func(v1, v2 DatapointCommand) bool {
		return v1.Equal(&v2)
	})
}

// String returns a human-readable representation of the request.
func (r *setDatapointValueReq) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString("setDatapointValueReq[commands=[")
	for i, v := range r.Commands {
		buffer.WriteString(v.String())
		if i != len(r.Commands)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("]]")
	return buffer.String()
}

// MarshalBinary encodes the request into its BAOS wire format.
func (r *setDatapointValueReq) MarshalBinary() ([]byte, error) {
	result := make([]byte, 6+datapointsSizeInBytes(r.Commands))
	result[0] = MainService
	result[1] = SubServiceSetDatapointValueReq
	startDatapoint := DatapointID(0)
	if len(r.Commands) > 0 {
		startDatapoint = r.Commands[0].ID
	}
	binary.BigEndian.PutUint16(result[2:4], uint16(startDatapoint))
	binary.BigEndian.PutUint16(result[4:6], uint16(len(r.Commands)))
	if err := marshalDatapointsTo(result[6:], r.Commands); err != nil {
		return nil, err
	}
	return result, nil
}

// UnmarshalBinary decodes the BAOS wire format into the request.
func (r *setDatapointValueReq) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) < 6:
		return fmt.Errorf("%w data length %d, expected at least %d", ErrInvalid, len(data), 6)
	case data[0] != MainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], MainService)
	case data[1] != SubServiceSetDatapointValueReq:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], SubServiceSetDatapointValueReq)
	}
	num := int(binary.BigEndian.Uint16(data[4:6]))
	values, err := unmarshalDatapointsFrom[DatapointCommand](data[6:], num)
	if err != nil {
		return err
	}
	r.Commands = values
	return nil
}

// Service returns the main service identifier.
func (r *setDatapointValueRes) Service() byte {
	return MainService
}

// SubService returns the sub-service identifier for the response.
func (r *setDatapointValueRes) SubService() byte {
	return SubServiceSetDatapointValueRes
}

// Equal reports whether r and other are identical responses.
func (r *setDatapointValueRes) Equal(other *setDatapointValueRes) bool {
	return r.StartDatapoint == other.StartDatapoint &&
		r.NumberOfDatapoints == other.NumberOfDatapoints &&
		byte(r.ErrorCode) == byte(other.ErrorCode)
}

// String returns a human-readable representation of the response.
func (r *setDatapointValueRes) String() string {
	return fmt.Sprintf("setDatapointValueRes[startDatapoint=%d,numberOfDatapoints=%d,errorCode=%d]",
		r.StartDatapoint, r.NumberOfDatapoints, r.ErrorCode)
}

// MarshalBinary encodes the response into its BAOS wire format.
func (r *setDatapointValueRes) MarshalBinary() ([]byte, error) {
	result := make([]byte, 7)
	result[0] = MainService
	result[1] = SubServiceSetDatapointValueRes
	binary.BigEndian.PutUint16(result[2:4], uint16(r.StartDatapoint))
	binary.BigEndian.PutUint16(result[4:6], r.NumberOfDatapoints)
	result[6] = byte(r.ErrorCode)
	return result, nil
}

// UnmarshalBinary decodes the BAOS wire format into the response.
func (r *setDatapointValueRes) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) != 7:
		return fmt.Errorf("%w data length %d, expected to be %d", ErrInvalid, len(data), 7)
	case data[0] != MainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], MainService)
	case data[1] != SubServiceSetDatapointValueRes:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], SubServiceSetDatapointValueRes)
	}
	r.StartDatapoint = DatapointID(binary.BigEndian.Uint16(data[2:4]))
	r.NumberOfDatapoints = binary.BigEndian.Uint16(data[4:6])
	r.ErrorCode = ErrorCode(data[6])
	return nil
}
