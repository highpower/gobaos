package baos

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"slices"
)

// ItemID identifies a single BAOS server item.
type ItemID uint16

// GetServerItemReq is a request to read server items.
type GetServerItemReq struct {
	StartItem     ItemID
	NumberOfItems uint16
}

// GetServerItemRes is a response containing server items or an error.
type GetServerItemRes struct {
	StartItem ItemID
	Items     []ServerItem
	ErrorCode ErrorCode
}

// Service returns the main service identifier.
func (r *GetServerItemReq) Service() byte {
	return MainService
}

// SubService returns the sub-service identifier for the request.
func (r *GetServerItemReq) SubService() byte {
	return SubServiceGetServerItemReq
}

// Valid validates the request fields.
func (r *GetServerItemReq) Valid() error {
	return nil
}

// Equal reports whether r and other are identical requests.
func (r *GetServerItemReq) Equal(other *GetServerItemReq) bool {
	return r.StartItem == other.StartItem && r.NumberOfItems == other.NumberOfItems
}

// String returns a human-readable representation of the request.
func (r *GetServerItemReq) String() string {
	return fmt.Sprintf("getServerItemReq[startItem=%d,numberOfItems=%d]", r.StartItem, r.NumberOfItems)
}

// MarshalBinary encodes the request into its BAOS wire format.
func (r *GetServerItemReq) MarshalBinary() ([]byte, error) {
	b := make([]byte, 6)
	b[0] = MainService
	b[1] = SubServiceGetServerItemReq
	binary.BigEndian.PutUint16(b[2:4], uint16(r.StartItem))
	binary.BigEndian.PutUint16(b[4:6], r.NumberOfItems)
	return b, nil
}

// UnmarshalBinary decodes the BAOS wire format into the request.
func (r *GetServerItemReq) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) != 6:
		return fmt.Errorf("%w data length %d, expected %d", ErrInvalid, len(data), 6)
	case data[0] != MainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], MainService)
	case data[1] != SubServiceGetServerItemReq:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], SubServiceGetServerItemReq)
	}
	r.StartItem = ItemID(binary.BigEndian.Uint16(data[2:4]))
	r.NumberOfItems = binary.BigEndian.Uint16(data[4:6])
	return nil
}

// Service returns the main service identifier.
func (r *GetServerItemRes) Service() byte {
	return MainService
}

// SubService returns the sub-service identifier for the response.
func (r *GetServerItemRes) SubService() byte {
	return SubServiceGetServerItemRes
}

// Valid validates the response fields.
func (r *GetServerItemRes) Valid() error {
	switch {
	case len(r.Items) == 0 && r.ErrorCode == 0:
		return fmt.Errorf("%w response containes 0 items but error code is 0", ErrInvalid)
	case len(r.Items) != 0 && r.ErrorCode != 0:
		return fmt.Errorf("%w response contains %d items and error code is %d", ErrInvalid, len(r.Items), r.ErrorCode)
	}
	return nil
}

// Equal reports whether r and other are identical responses.
func (r *GetServerItemRes) Equal(other *GetServerItemRes) bool {
	if r.StartItem != other.StartItem || byte(r.ErrorCode) != byte(other.ErrorCode) {
		return false
	}
	return slices.EqualFunc(r.Items, other.Items, func(i1, i2 ServerItem) bool {
		return i1.Equal(&i2)
	})
}

// String returns a human-readable representation of the response.
func (r *GetServerItemRes) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString(fmt.Sprintf("getServerItemRes[startItem=%d,items=[", r.StartItem))
	for i, it := range r.Items {
		buffer.WriteString(it.String())
		if i != len(r.Items)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString(fmt.Sprintf("],errorCode=%d]", r.ErrorCode))
	return buffer.String()
}

// MarshalBinary encodes the response into its BAOS wire format.
func (r *GetServerItemRes) MarshalBinary() ([]byte, error) {
	if r.ErrorCode != 0 {
		result := make([]byte, 7)
		result[0] = MainService
		result[1] = SubServiceGetServerItemRes
		binary.BigEndian.PutUint16(result[2:4], uint16(r.StartItem))
		binary.BigEndian.PutUint16(result[4:6], uint16(0))
		result[6] = byte(r.ErrorCode)
		return result, nil
	}
	result := make([]byte, 6+sizeInBytes(r.Items))
	result[0] = MainService
	result[1] = SubServiceGetServerItemRes
	binary.BigEndian.PutUint16(result[2:4], uint16(r.StartItem))
	binary.BigEndian.PutUint16(result[4:6], uint16(len(r.Items)))
	if err := marshalItemsTo(result[6:], r.Items); err != nil {
		return nil, err
	}
	return result, nil
}

// UnmarshalBinary decodes the BAOS wire format into the response.
func (r *GetServerItemRes) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) < 6:
		return fmt.Errorf("%w data length %d, expected at least 6", ErrInvalid, len(data))
	case data[0] != MainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], MainService)
	case data[1] != SubServiceGetServerItemRes:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], SubServiceGetServerItemRes)
	}
	startItem := ItemID(binary.BigEndian.Uint16(data[2:4]))
	num := int(binary.BigEndian.Uint16(data[4:6]))
	if len(data) == 7 {
		*r = GetServerItemRes{StartItem: startItem, ErrorCode: ErrorCode(data[6])}
		return nil
	}
	items, err := unmarshalItemsFrom(data[6:], num)
	if err != nil {
		return err
	}
	*r = GetServerItemRes{StartItem: startItem, Items: items}
	return nil
}
