package baos

import (
	"encoding/binary"
	"fmt"
	"slices"
)

// setServerItemReq is a request to write server items.
type setServerItemReq struct {
	items []ServerItem
}

// setServerItemRes is a response confirming the set operation or reporting an error.
type setServerItemRes struct {
	startItem     ItemID
	numberOfItems uint16
	errorCode     ErrorCode
}

// MarshalBinary encodes the request into its BAOS wire format.
func (r *setServerItemReq) MarshalBinary() ([]byte, error) {
	result := make([]byte, 6+sizeInBytes(r.items))
	result[0] = mainService
	result[1] = subServiceSetServerItemReq
	startItem := ItemID(0)
	if len(r.items) > 0 {
		startItem = r.items[0].ID
	}
	binary.BigEndian.PutUint16(result[2:4], uint16(startItem))
	binary.BigEndian.PutUint16(result[4:6], uint16(len(r.items)))
	if err := marshalItemsTo(result[6:], r.items); err != nil {
		return nil, err
	}
	return result, nil
}

// UnmarshalBinary decodes the BAOS wire format into the request.
func (r *setServerItemReq) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) < 6:
		return fmt.Errorf("%w data length %d, expected at least %d", ErrInvalid, len(data), 6)
	case data[0] != mainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], mainService)
	case data[1] != subServiceSetServerItemReq:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], subServiceSetServerItemReq)
	}
	num := int(binary.BigEndian.Uint16(data[4:6]))
	items, err := unmarshalItemsFrom(data[6:], num)
	if err != nil {
		return err
	}
	r.items = items
	return nil
}

// equal reports whether r and other are identical requests.
func (r *setServerItemReq) equal(other *setServerItemReq) bool {
	return slices.EqualFunc(r.items, other.items, func(i1, i2 ServerItem) bool {
		return i1.Equal(&i2)
	})
}

// MarshalBinary encodes the response into its BAOS wire format.
func (r *setServerItemRes) MarshalBinary() ([]byte, error) {
	result := make([]byte, 7)
	result[0] = mainService
	result[1] = subServiceSetServerItemRes
	binary.BigEndian.PutUint16(result[2:4], uint16(r.startItem))
	binary.BigEndian.PutUint16(result[4:6], r.numberOfItems)
	result[6] = byte(r.errorCode)
	return result, nil
}

// UnmarshalBinary decodes the BAOS wire format into the response.
func (r *setServerItemRes) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) != 7:
		return fmt.Errorf("%w data length %d, expected to be %d", ErrInvalid, len(data), 7)
	case data[0] != mainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], mainService)
	case data[1] != subServiceSetServerItemRes:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], subServiceSetServerItemRes)
	}
	r.startItem = ItemID(binary.BigEndian.Uint16(data[2:4]))
	r.numberOfItems = binary.BigEndian.Uint16(data[4:6])
	r.errorCode = ErrorCode(data[6])
	return nil
}

// equal reports whether r and other are identical responses.
func (r *setServerItemRes) equal(other *setServerItemRes) bool {
	return r.startItem == other.startItem && r.numberOfItems == other.numberOfItems &&
		byte(r.errorCode) == byte(other.errorCode)
}
