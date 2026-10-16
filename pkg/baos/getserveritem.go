package baos

import (
	"encoding/binary"
	"fmt"
	"slices"
)

// ItemID identifies a single BAOS server item.
type ItemID uint16

// getServerItemReq is a request to read server items.
type getServerItemReq struct {
	startItem     ItemID
	numberOfItems uint16
}

// getServerItemRes is a response containing server items or an error.
type getServerItemRes struct {
	startItem ItemID
	items     []ServerItem
	errorCode ErrorCode
}

// MarshalBinary encodes the request into its BAOS wire format.
func (r *getServerItemReq) MarshalBinary() ([]byte, error) {
	b := make([]byte, 6)
	b[0] = mainService
	b[1] = subServiceGetServerItemReq
	binary.BigEndian.PutUint16(b[2:4], uint16(r.startItem))
	binary.BigEndian.PutUint16(b[4:6], r.numberOfItems)
	return b, nil
}

// UnmarshalBinary decodes the BAOS wire format into the request.
func (r *getServerItemReq) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) != 6:
		return fmt.Errorf("%w data length %d, expected %d", ErrInvalid, len(data), 6)
	case data[0] != mainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], mainService)
	case data[1] != subServiceGetServerItemReq:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], subServiceGetServerItemReq)
	}
	r.startItem = ItemID(binary.BigEndian.Uint16(data[2:4]))
	r.numberOfItems = binary.BigEndian.Uint16(data[4:6])
	return nil
}

// equal reports whether r and other are identical requests.
func (r *getServerItemReq) equal(other *getServerItemReq) bool {
	return r.startItem == other.startItem && r.numberOfItems == other.numberOfItems
}

// MarshalBinary encodes the response into its BAOS wire format.
func (r *getServerItemRes) MarshalBinary() ([]byte, error) {
	if r.errorCode != 0 {
		result := make([]byte, 7)
		result[0] = mainService
		result[1] = subServiceGetServerItemRes
		binary.BigEndian.PutUint16(result[2:4], uint16(r.startItem))
		binary.BigEndian.PutUint16(result[4:6], uint16(0))
		result[6] = byte(r.errorCode)
		return result, nil
	}
	result := make([]byte, 6+sizeInBytes(r.items))
	result[0] = mainService
	result[1] = subServiceGetServerItemRes
	binary.BigEndian.PutUint16(result[2:4], uint16(r.startItem))
	binary.BigEndian.PutUint16(result[4:6], uint16(len(r.items)))
	if err := marshalItemsTo(result[6:], r.items); err != nil {
		return nil, err
	}
	return result, nil
}

// UnmarshalBinary decodes the BAOS wire format into the response.
func (r *getServerItemRes) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) < 6:
		return fmt.Errorf("%w data length %d, expected at least 6", ErrInvalid, len(data))
	case data[0] != mainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], mainService)
	case data[1] != subServiceGetServerItemRes:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], subServiceGetServerItemRes)
	}
	startItem := ItemID(binary.BigEndian.Uint16(data[2:4]))
	num := int(binary.BigEndian.Uint16(data[4:6]))
	if len(data) == 7 {
		*r = getServerItemRes{startItem: startItem, errorCode: ErrorCode(data[6])}
		return nil
	}
	items, err := unmarshalItemsFrom(data[6:], num)
	if err != nil {
		return err
	}
	*r = getServerItemRes{startItem: startItem, items: items}
	return nil
}

// Equal reports whether r and other are identical responses.
func (r *getServerItemRes) Equal(other *getServerItemRes) bool {
	if r.startItem != other.startItem || byte(r.errorCode) != byte(other.errorCode) {
		return false
	}
	return slices.EqualFunc(r.items, other.items, func(i1, i2 ServerItem) bool {
		return i1.Equal(&i2)
	})
}
