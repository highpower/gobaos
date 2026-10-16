package baos

import (
	"encoding/binary"
	"fmt"
	"gobaos/pkg/knxtypes"
	"slices"
)

// DatapointDescription holds the static description of a single BAOS datapoint.
type DatapointDescription struct {
	ID          DatapointID
	ValueType   ValueType
	ConfigFlags byte
	Code        knxtypes.TypeCode
}

// getDatapointDescriptionReq is a request to read datapoint descriptions.
type getDatapointDescriptionReq struct {
	startDatapoint     DatapointID
	numberOfDatapoints uint16
}

// getDatapointDescriptionRes is a response containing datapoint descriptions or an error.
type getDatapointDescriptionRes struct {
	startDatapoint DatapointID
	descriptions   []DatapointDescription
	errorCode      ErrorCode
}

const (
	datapointDescriptionSize    = 5
	getDescriptionResHeaderSize = 6
)

// String returns a human-readable representation of the datapoint description.
func (d *DatapointDescription) String() string {
	return fmt.Sprintf("datapointDescription[id=%d,valueType=%d,configFlags=%d,dpt=%d]",
		d.ID, d.ValueType, d.ConfigFlags, d.Code)
}

// Equal reports whether d and other represent the same datapoint description.
func (d *DatapointDescription) Equal(other *DatapointDescription) bool {
	return d.ID == other.ID && d.ValueType == other.ValueType &&
		d.ConfigFlags == other.ConfigFlags && d.Code == other.Code
}

// MarshalBinary encodes the request into its BAOS wire format.
func (r *getDatapointDescriptionReq) MarshalBinary() ([]byte, error) {
	b := make([]byte, 6)
	b[0] = mainService
	b[1] = subServiceGetDatapointDescriptionReq
	binary.BigEndian.PutUint16(b[2:4], uint16(r.startDatapoint))
	binary.BigEndian.PutUint16(b[4:6], r.numberOfDatapoints)
	return b, nil
}

// UnmarshalBinary decodes the request from its BAOS wire format.
func (r *getDatapointDescriptionReq) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) != 6:
		return fmt.Errorf("%w data length %d, expected %d", ErrInvalid, len(data), 6)
	case data[0] != mainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], mainService)
	case data[1] != subServiceGetDatapointDescriptionReq:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1],
			subServiceGetDatapointDescriptionReq)
	}
	r.startDatapoint = DatapointID(binary.BigEndian.Uint16(data[2:4]))
	r.numberOfDatapoints = binary.BigEndian.Uint16(data[4:6])
	return nil
}

// equal reports whether r and other represent the same request.
func (r *getDatapointDescriptionReq) equal(other *getDatapointDescriptionReq) bool {
	return r.startDatapoint == other.startDatapoint && r.numberOfDatapoints == other.numberOfDatapoints
}

// MarshalBinary encodes the response into its BAOS wire format.
func (r *getDatapointDescriptionRes) MarshalBinary() ([]byte, error) {
	if r.errorCode != 0 {
		result := make([]byte, getDescriptionResHeaderSize+1)
		result[0] = mainService
		result[1] = subServiceGetDatapointDescriptionRes
		binary.BigEndian.PutUint16(result[2:4], uint16(r.startDatapoint))
		binary.BigEndian.PutUint16(result[4:6], uint16(0))
		result[6] = byte(r.errorCode)
		return result, nil
	}
	result := make([]byte, getDescriptionResHeaderSize+len(r.descriptions)*datapointDescriptionSize)
	result[0] = mainService
	result[1] = subServiceGetDatapointDescriptionRes
	binary.BigEndian.PutUint16(result[2:4], uint16(r.startDatapoint))
	binary.BigEndian.PutUint16(result[4:6], uint16(len(r.descriptions)))
	if err := marshalDescriptionsTo(result[6:], r.descriptions); err != nil {
		return nil, err
	}
	return result, nil
}

// UnmarshalBinary decodes the response from its BAOS wire format.
func (r *getDatapointDescriptionRes) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) < getDescriptionResHeaderSize:
		return fmt.Errorf("%w data length %d, expected at least %d", ErrInvalid, len(data), 6)
	case data[0] != mainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], mainService)
	case data[1] != subServiceGetDatapointDescriptionRes:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], subServiceGetDatapointDescriptionRes)
	}
	startDatapoint := DatapointID(binary.BigEndian.Uint16(data[2:4]))
	if len(data) == 7 {
		*r = getDatapointDescriptionRes{startDatapoint: startDatapoint, errorCode: ErrorCode(data[6])}
		return nil
	}
	num := int(binary.BigEndian.Uint16(data[4:6]))
	descriptions, err := unmarshalDescriptionsFrom(data[6:], num)
	if err != nil {
		return err
	}
	*r = getDatapointDescriptionRes{startDatapoint: startDatapoint, descriptions: descriptions}
	return nil
}

// equal reports whether r and other represent the same response.
func (r *getDatapointDescriptionRes) equal(other *getDatapointDescriptionRes) bool {
	if r.startDatapoint != other.startDatapoint || byte(r.errorCode) != byte(other.errorCode) {
		return false
	}
	return slices.EqualFunc(r.descriptions, other.descriptions, func(v1, v2 DatapointDescription) bool {
		return v1.Equal(&v2)
	})
}

func marshalDescriptionsTo(data []byte, descriptions []DatapointDescription) error {
	for i := range descriptions {
		if len(data) < datapointDescriptionSize {
			return fmt.Errorf("%w buffer has no capacity: %d while %d needed", ErrInvalid, len(data), datapointDescriptionSize)
		}
		binary.BigEndian.PutUint16(data[0:2], uint16(descriptions[i].ID))
		data[2] = byte(descriptions[i].ValueType)
		data[3] = descriptions[i].ConfigFlags
		data[4] = byte(descriptions[i].Code)
		data = data[datapointDescriptionSize:]
	}
	if len(data) != 0 {
		return fmt.Errorf("%w buffer has %d extra bytes", ErrInvalid, len(data))
	}
	return nil
}

func unmarshalDescriptionsFrom(data []byte, n int) ([]DatapointDescription, error) {
	result := make([]DatapointDescription, 0, n)
	for range n {
		if len(data) < datapointDescriptionSize {
			return nil, fmt.Errorf("%w buffer has only %d bytes while %d needed", ErrInvalid,
				len(data), datapointDescriptionSize)
		}
		result = append(result, DatapointDescription{
			ID:          DatapointID(binary.BigEndian.Uint16(data[0:2])),
			ValueType:   ValueType(data[2]),
			ConfigFlags: data[3],
			Code:        knxtypes.TypeCode(data[4]),
		})
		data = data[datapointDescriptionSize:]
	}
	if len(data) != 0 {
		return nil, fmt.Errorf("%w buffer has %d extra bytes", ErrInvalid, len(data))
	}
	return result, nil
}
