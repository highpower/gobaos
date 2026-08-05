package baos

import (
	"encoding/binary"
	"fmt"
	"slices"
)

// GetDescriptionStringReq is a request to read datapoint description strings.
type GetDescriptionStringReq struct {
	StartDatapoint  DatapointID
	NumberOfStrings uint16
}

// GetDescriptionStringRes is a response containing datapoint descriptions strings or an error.
type GetDescriptionStringRes struct {
	StartDatapoint DatapointID
	Strings        []string
	ErrorCode      ErrorCode
}

// Service returns the main service identifier.
func (r *GetDescriptionStringReq) Service() byte {
	return MainService
}

// SubService returns the sub-service identifier.
func (r *GetDescriptionStringReq) SubService() byte {
	return SubServiceGetDescriptionStringReq
}

// Valid checks whether the request fields are valid.
func (r *GetDescriptionStringReq) Valid() error {
	return nil
}

// Equal reports whether r and other represent the same request.
func (r *GetDescriptionStringReq) Equal(other *GetDescriptionStringReq) bool {
	return r.StartDatapoint == other.StartDatapoint && r.NumberOfStrings == other.NumberOfStrings
}

// String returns a human-readable representation of the request.
func (r *GetDescriptionStringReq) String() string {
	return fmt.Sprintf("getDescriptionStringReq[startDatapoint=%d,numberOfStrings=%d]",
		r.StartDatapoint, r.NumberOfStrings)
}

// MarshalBinary encodes the request into its BAOS wire format.
func (r *GetDescriptionStringReq) MarshalBinary() ([]byte, error) {
	b := make([]byte, 6)
	b[0] = MainService
	b[1] = SubServiceGetDescriptionStringReq
	binary.BigEndian.PutUint16(b[2:4], uint16(r.StartDatapoint))
	binary.BigEndian.PutUint16(b[4:6], r.NumberOfStrings)
	return b, nil
}

// UnmarshalBinary decodes the request from its BAOS wire format.
func (r *GetDescriptionStringReq) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) != 6:
		return fmt.Errorf("%w data length %d, expected %d", ErrInvalid, len(data), 6)
	case data[0] != MainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], MainService)
	case data[1] != SubServiceGetDescriptionStringReq:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1],
			SubServiceGetDescriptionStringReq)
	}
	r.StartDatapoint = DatapointID(binary.BigEndian.Uint16(data[2:4]))
	r.NumberOfStrings = binary.BigEndian.Uint16(data[4:6])
	return nil
}

// Service returns the main service identifier.
func (r *GetDescriptionStringRes) Service() byte {
	return MainService
}

// SubService returns the sub-service identifier.
func (r *GetDescriptionStringRes) SubService() byte {
	return SubServiceGetDescriptionStringRes
}

// Valid checks whether the response fields are consistent.
func (r *GetDescriptionStringRes) Valid() error {
	switch {
	case len(r.Strings) == 0 && r.ErrorCode == 0:
		return fmt.Errorf("%w response contains 0 strings but error code is 0", ErrInvalid)
	case len(r.Strings) != 0 && r.ErrorCode != 0:
		return fmt.Errorf("%w response contains %d strings but error code is %d", ErrInvalid, len(r.Strings), r.ErrorCode)
	}
	return nil
}

// Equal reports whether r and other represent the same response.
func (r *GetDescriptionStringRes) Equal(other *GetDescriptionStringRes) bool {
	if r.StartDatapoint != other.StartDatapoint || byte(r.ErrorCode) != byte(other.ErrorCode) {
		return false
	}
	return slices.Equal(r.Strings, other.Strings)
}

// String returns a human-readable representation of the response.
func (r *GetDescriptionStringRes) String() string {
	return fmt.Sprintf("getDescriptionStringRes[startString=%d,strings=%v,errorCode=%d]",
		r.StartDatapoint, r.Strings, r.ErrorCode)
}

// MarshalBinary encodes the response into its BAOS wire format.
func (r *GetDescriptionStringRes) MarshalBinary() ([]byte, error) {
	if r.ErrorCode != 0 {
		result := make([]byte, 7)
		result[0] = MainService
		result[1] = SubServiceGetDescriptionStringRes
		binary.BigEndian.PutUint16(result[2:4], uint16(r.StartDatapoint))
		binary.BigEndian.PutUint16(result[4:6], uint16(0))
		result[6] = byte(r.ErrorCode)
		return result, nil
	}
	result := make([]byte, 6+stringsSizeInBytes(r.Strings))
	result[0] = MainService
	result[1] = SubServiceGetDescriptionStringRes
	binary.BigEndian.PutUint16(result[2:4], uint16(r.StartDatapoint))
	binary.BigEndian.PutUint16(result[4:6], uint16(len(r.Strings)))
	if err := marshalStringsTo(result[6:], r.Strings); err != nil {
		return nil, err
	}
	return result, nil
}

// UnmarshalBinary decodes the response from its BAOS wire format.
func (r *GetDescriptionStringRes) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) < 6:
		return fmt.Errorf("%w data length %d, expected at least %d", ErrInvalid, len(data), 6)
	case data[0] != MainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], MainService)
	case data[1] != SubServiceGetDescriptionStringRes:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], SubServiceGetDescriptionStringRes)
	}
	start := binary.BigEndian.Uint16(data[2:4])
	if len(data) == 7 {
		*r = GetDescriptionStringRes{StartDatapoint: DatapointID(start), ErrorCode: ErrorCode(data[6])}
		return nil
	}
	num := int(binary.BigEndian.Uint16(data[4:6]))
	strings, err := unmarshalStringsFrom(data[6:], num)
	if err != nil {
		return err
	}
	*r = GetDescriptionStringRes{StartDatapoint: DatapointID(start), Strings: strings}
	return nil
}

func stringsSizeInBytes(strings []string) int {
	totalSize := 0
	for _, s := range strings {
		totalSize += 2 + len(s)
	}
	return totalSize
}

func marshalStringsTo(data []byte, strings []string) error {
	for _, s := range strings {
		if len(data) < 2+len(s) {
			return fmt.Errorf("%w buffer has no capacity: %d while %d needed", ErrInvalid, len(data), 2+len(s))
		}
		binary.BigEndian.PutUint16(data[0:2], uint16(len(s)))
		copy(data[2:], s)
		data = data[2+len(s):]
	}
	if len(data) != 0 {
		return fmt.Errorf("%w buffer has %d extra bytes", ErrInvalid, len(data))
	}
	return nil
}

func unmarshalStringsFrom(data []byte, n int) ([]string, error) {
	result := make([]string, 0, n)
	for i := 0; i < n; i++ {
		if len(data) < 2 {
			return nil, fmt.Errorf("%w buffer has only %d bytes while 2 needed for length", ErrInvalid, len(data))
		}
		length := int(binary.BigEndian.Uint16(data[0:2]))
		if len(data) < 2+length {
			return nil, fmt.Errorf("%w buffer has only %d bytes while %d needed", ErrInvalid, len(data), 2+length)
		}
		result = append(result, string(data[2:2+length]))
		data = data[2+length:]
	}
	if len(data) != 0 {
		return nil, fmt.Errorf("%w buffer has %d extra bytes", ErrInvalid, len(data))
	}
	return result, nil
}
