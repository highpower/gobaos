package baos

import (
	"encoding/binary"
	"fmt"
	"slices"

	"gobaos/pkg/knxtypes"
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
	StartDatapoint     DatapointID
	NumberOfDatapoints uint16
}

// getDatapointDescriptionRes is a response containing datapoint descriptions or an error.
type getDatapointDescriptionRes struct {
	StartDatapoint DatapointID
	Descriptions   []DatapointDescription
	ErrorCode      ErrorCode
}

const (
	datapointDescriptionSize    = 5
	getDescriptionResHeaderSize = 6
)

// Equal reports whether d and other represent the same datapoint description.
func (d *DatapointDescription) Equal(other *DatapointDescription) bool {
	return d.ID == other.ID && d.ValueType == other.ValueType &&
		d.ConfigFlags == other.ConfigFlags && d.Code == other.Code
}

// String returns a human-readable representation of the datapoint description.
func (d *DatapointDescription) String() string {
	return fmt.Sprintf("datapointDescription[id=%d,valueType=%d,configFlags=%d,dpt=%d]",
		d.ID, d.ValueType, d.ConfigFlags, d.Code)
}

// Service returns the main service identifier.
func (r *getDatapointDescriptionReq) Service() byte {
	return MainService
}

// SubService returns the sub-service identifier.
func (r *getDatapointDescriptionReq) SubService() byte {
	return SubServiceGetDatapointDescriptionReq
}

// Equal reports whether r and other represent the same request.
func (r *getDatapointDescriptionReq) Equal(other *getDatapointDescriptionReq) bool {
	return r.StartDatapoint == other.StartDatapoint && r.NumberOfDatapoints == other.NumberOfDatapoints
}

// String returns a human-readable representation of the request.
func (r *getDatapointDescriptionReq) String() string {
	return fmt.Sprintf("getDatapointDescriptionReq[startDatapoint=%d,numberOfDatapoints=%d]",
		r.StartDatapoint, r.NumberOfDatapoints)
}

// MarshalBinary encodes the request into its BAOS wire format.
func (r *getDatapointDescriptionReq) MarshalBinary() ([]byte, error) {
	b := make([]byte, 6)
	b[0] = MainService
	b[1] = SubServiceGetDatapointDescriptionReq
	binary.BigEndian.PutUint16(b[2:4], uint16(r.StartDatapoint))
	binary.BigEndian.PutUint16(b[4:6], r.NumberOfDatapoints)
	return b, nil
}

// UnmarshalBinary decodes the request from its BAOS wire format.
func (r *getDatapointDescriptionReq) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) != 6:
		return fmt.Errorf("%w data length %d, expected %d", ErrInvalid, len(data), 6)
	case data[0] != MainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], MainService)
	case data[1] != SubServiceGetDatapointDescriptionReq:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1],
			SubServiceGetDatapointDescriptionReq)
	}
	r.StartDatapoint = DatapointID(binary.BigEndian.Uint16(data[2:4]))
	r.NumberOfDatapoints = binary.BigEndian.Uint16(data[4:6])
	return nil
}

// Service returns the main service identifier.
func (r *getDatapointDescriptionRes) Service() byte {
	return MainService
}

// SubService returns the sub-service identifier.
func (r *getDatapointDescriptionRes) SubService() byte {
	return SubServiceGetDatapointDescriptionRes
}

// Equal reports whether r and other represent the same response.
func (r *getDatapointDescriptionRes) Equal(other *getDatapointDescriptionRes) bool {
	if r.StartDatapoint != other.StartDatapoint || byte(r.ErrorCode) != byte(other.ErrorCode) {
		return false
	}
	return slices.EqualFunc(r.Descriptions, other.Descriptions, func(v1, v2 DatapointDescription) bool {
		return v1.Equal(&v2)
	})
}

// MarshalBinary encodes the response into its BAOS wire format.
func (r *getDatapointDescriptionRes) MarshalBinary() ([]byte, error) {
	if r.ErrorCode != 0 {
		result := make([]byte, getDescriptionResHeaderSize+1)
		result[0] = MainService
		result[1] = SubServiceGetDatapointDescriptionRes
		binary.BigEndian.PutUint16(result[2:4], uint16(r.StartDatapoint))
		binary.BigEndian.PutUint16(result[4:6], uint16(0))
		result[6] = byte(r.ErrorCode)
		return result, nil
	}
	result := make([]byte, getDescriptionResHeaderSize+len(r.Descriptions)*datapointDescriptionSize)
	result[0] = MainService
	result[1] = SubServiceGetDatapointDescriptionRes
	binary.BigEndian.PutUint16(result[2:4], uint16(r.StartDatapoint))
	binary.BigEndian.PutUint16(result[4:6], uint16(len(r.Descriptions)))
	if err := marshalDescriptionsTo(result[6:], r.Descriptions); err != nil {
		return nil, err
	}
	return result, nil
}

// UnmarshalBinary decodes the response from its BAOS wire format.
func (r *getDatapointDescriptionRes) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) < getDescriptionResHeaderSize:
		return fmt.Errorf("%w data length %d, expected at least %d", ErrInvalid, len(data), 6)
	case data[0] != MainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], MainService)
	case data[1] != SubServiceGetDatapointDescriptionRes:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], SubServiceGetDatapointDescriptionRes)
	}
	startDatapoint := DatapointID(binary.BigEndian.Uint16(data[2:4]))
	if len(data) == 7 {
		*r = getDatapointDescriptionRes{StartDatapoint: startDatapoint, ErrorCode: ErrorCode(data[6])}
		return nil
	}
	num := int(binary.BigEndian.Uint16(data[4:6]))
	descriptions, err := unmarshalDescriptionsFrom(data[6:], num)
	if err != nil {
		return err
	}
	*r = getDatapointDescriptionRes{StartDatapoint: startDatapoint, Descriptions: descriptions}
	return nil
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
