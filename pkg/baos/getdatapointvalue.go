package baos

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"slices"
)

type Filter byte

type GetDatapointValueReq struct {
	StartDatapoint     DatapointID
	NumberOfDatapoints uint16
	Filter             Filter
}

type GetDatapointValueRes struct {
	StartDatapoint DatapointID
	Values         []Datapoint
	ErrorCode      ErrorCode
}

const (
	AllDatapoints     Filter = 0
	ValidDatapoints   Filter = 1
	UpdatedDatapoints Filter = 2
)

func (r *GetDatapointValueReq) Service() byte {
	return MainService
}

func (r *GetDatapointValueReq) SubService() byte {
	return SubServiceGetDatapointValueReq
}

func (r *GetDatapointValueReq) Valid() error {
	switch r.Filter {
	case AllDatapoints, ValidDatapoints, UpdatedDatapoints:
		return nil
	}
	return fmt.Errorf("%w filter: %d", ErrInvalid, r.Filter)
}

func (r *GetDatapointValueReq) Equal(other *GetDatapointValueReq) bool {
	return r.StartDatapoint == other.StartDatapoint && r.NumberOfDatapoints == other.NumberOfDatapoints &&
		r.Filter == other.Filter
}

func (r *GetDatapointValueReq) String() string {
	return fmt.Sprintf("getDatapointValueReq[startDatapoint=%d,numberOfDatapoints=%d,filter=%d]",
		r.StartDatapoint, r.NumberOfDatapoints, r.Filter)
}

func (r *GetDatapointValueReq) MarshalBinary() ([]byte, error) {
	b := make([]byte, 7)
	b[0] = MainService
	b[1] = SubServiceGetDatapointValueReq
	binary.BigEndian.PutUint16(b[2:4], uint16(r.StartDatapoint))
	binary.BigEndian.PutUint16(b[4:6], r.NumberOfDatapoints)
	b[6] = byte(r.Filter)
	return b, nil
}

func (r *GetDatapointValueReq) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) != 7:
		return fmt.Errorf("%w data length %d, expected %d", ErrInvalid, len(data), 7)
	case data[0] != MainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], MainService)
	case data[1] != SubServiceGetDatapointValueReq:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], SubServiceGetDatapointValueReq)
	}
	r.StartDatapoint = DatapointID(binary.BigEndian.Uint16(data[2:4]))
	r.NumberOfDatapoints = binary.BigEndian.Uint16(data[4:6])
	r.Filter = Filter(data[6])
	return nil
}

func (r *GetDatapointValueRes) Service() byte {
	return MainService
}

func (r *GetDatapointValueRes) SubService() byte {
	return SubServiceGetDatapointValueRes
}

func (r *GetDatapointValueRes) Valid() error {
	switch {
	case len(r.Values) == 0 && r.ErrorCode == 0:
		return fmt.Errorf("%w response contains 0 values but error code is 0", ErrInvalid)
	case len(r.Values) != 0 && r.ErrorCode != 0:
		return fmt.Errorf("%w response contains %d values but error code is %d", ErrInvalid, len(r.Values), r.ErrorCode)
	}
	return nil
}

func (r *GetDatapointValueRes) Equal(other *GetDatapointValueRes) bool {
	if r.StartDatapoint != other.StartDatapoint || byte(r.ErrorCode) != byte(other.ErrorCode) {
		return false
	}
	return slices.EqualFunc(r.Values, other.Values, func(v1, v2 Datapoint) bool {
		return v1.Equal(&v2)
	})
}

func (r *GetDatapointValueRes) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString(fmt.Sprintf("getDatapointValueRes[startDatapoint=%d,datapoints=[", r.StartDatapoint))
	for i, val := range r.Values {
		buffer.WriteString(val.String())
		if i < len(r.Values)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString(fmt.Sprintf("],errorCode=%d]", r.ErrorCode))
	return buffer.String()
}

func (r *GetDatapointValueRes) MarshalBinary() ([]byte, error) {
	if r.ErrorCode != 0 {
		result := make([]byte, 7)
		result[0] = MainService
		result[1] = SubServiceGetDatapointValueRes
		binary.BigEndian.PutUint16(result[2:4], uint16(r.StartDatapoint))
		binary.BigEndian.PutUint16(result[4:6], uint16(0))
		result[6] = byte(r.ErrorCode)
		return result, nil
	}
	result := make([]byte, 6+datapointsSizeInBytes(r.Values))
	result[0] = MainService
	result[1] = SubServiceGetDatapointValueRes
	binary.BigEndian.PutUint16(result[2:4], uint16(r.StartDatapoint))
	binary.BigEndian.PutUint16(result[4:6], uint16(len(r.Values)))
	if err := marshalDatapointsTo(result[6:], r.Values); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *GetDatapointValueRes) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) < 6:
		return fmt.Errorf("%w data length %d, expected at least %d", ErrInvalid, len(data), 6)
	case data[0] != MainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], MainService)
	case data[1] != SubServiceGetDatapointValueRes:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], SubServiceGetDatapointValueRes)
	}
	startDatapoint := DatapointID(binary.BigEndian.Uint16(data[2:4]))
	if len(data) == 7 {
		*r = GetDatapointValueRes{StartDatapoint: startDatapoint, ErrorCode: ErrorCode(data[6])}
		return nil
	}
	num := int(binary.BigEndian.Uint16(data[4:6]))
	values, err := unmarshalDatapointsFrom[Datapoint](data[6:], num)
	if err != nil {
		return err
	}
	*r = GetDatapointValueRes{StartDatapoint: startDatapoint, Values: values}
	return nil
}
