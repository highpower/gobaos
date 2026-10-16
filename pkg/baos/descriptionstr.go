package baos

import (
	"encoding/binary"
	"fmt"
	"slices"
)

// getDescriptionStringReq is a request to read datapoint description strings.
type getDescriptionStringReq struct {
	startDatapoint  DatapointID
	numberOfStrings uint16
}

// getDescriptionStringRes is a response containing datapoint descriptions strings or an error.
type getDescriptionStringRes struct {
	startDatapoint DatapointID
	values         []string
	errorCode      ErrorCode
}

// MarshalBinary encodes the request into its BAOS wire format.
func (r *getDescriptionStringReq) MarshalBinary() ([]byte, error) {
	b := make([]byte, 6)
	b[0] = mainService
	b[1] = subServiceGetDescriptionStringReq
	binary.BigEndian.PutUint16(b[2:4], uint16(r.startDatapoint))
	binary.BigEndian.PutUint16(b[4:6], r.numberOfStrings)
	return b, nil
}

// UnmarshalBinary decodes the request from its BAOS wire format.
func (r *getDescriptionStringReq) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) != 6:
		return fmt.Errorf("%w data length %d, expected %d", ErrInvalid, len(data), 6)
	case data[0] != mainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], mainService)
	case data[1] != subServiceGetDescriptionStringReq:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1],
			subServiceGetDescriptionStringReq)
	}
	r.startDatapoint = DatapointID(binary.BigEndian.Uint16(data[2:4]))
	r.numberOfStrings = binary.BigEndian.Uint16(data[4:6])
	return nil
}

// equal reports whether r and other represent the same request.
func (r *getDescriptionStringReq) equal(other *getDescriptionStringReq) bool {
	return r.startDatapoint == other.startDatapoint && r.numberOfStrings == other.numberOfStrings
}

// MarshalBinary encodes the response into its BAOS wire format.
func (r *getDescriptionStringRes) MarshalBinary() ([]byte, error) {
	if r.errorCode != 0 {
		result := make([]byte, 7)
		result[0] = mainService
		result[1] = subServiceGetDescriptionStringRes
		binary.BigEndian.PutUint16(result[2:4], uint16(r.startDatapoint))
		binary.BigEndian.PutUint16(result[4:6], uint16(0))
		result[6] = byte(r.errorCode)
		return result, nil
	}
	result := make([]byte, 6+stringsSizeInBytes(r.values))
	result[0] = mainService
	result[1] = subServiceGetDescriptionStringRes
	binary.BigEndian.PutUint16(result[2:4], uint16(r.startDatapoint))
	binary.BigEndian.PutUint16(result[4:6], uint16(len(r.values)))
	if err := marshalStringsTo(result[6:], r.values); err != nil {
		return nil, err
	}
	return result, nil
}

// UnmarshalBinary decodes the response from its BAOS wire format.
func (r *getDescriptionStringRes) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) < 6:
		return fmt.Errorf("%w data length %d, expected at least %d", ErrInvalid, len(data), 6)
	case data[0] != mainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], mainService)
	case data[1] != subServiceGetDescriptionStringRes:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], subServiceGetDescriptionStringRes)
	}
	start := binary.BigEndian.Uint16(data[2:4])
	if len(data) == 7 {
		*r = getDescriptionStringRes{startDatapoint: DatapointID(start), errorCode: ErrorCode(data[6])}
		return nil
	}
	num := int(binary.BigEndian.Uint16(data[4:6]))
	strings, err := unmarshalStringsFrom(data[6:], num)
	if err != nil {
		return err
	}
	*r = getDescriptionStringRes{startDatapoint: DatapointID(start), values: strings}
	return nil
}

// equal reports whether r and other represent the same response.
func (r *getDescriptionStringRes) equal(other *getDescriptionStringRes) bool {
	if r.startDatapoint != other.startDatapoint || byte(r.errorCode) != byte(other.errorCode) {
		return false
	}
	return slices.Equal(r.values, other.values)
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
	for range n {
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
