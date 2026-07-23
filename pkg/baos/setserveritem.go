package baos

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"slices"
)

type SetServerItemReq struct {
	Items []ServerItem
}

type SetServerItemRes struct {
	StartItem     ItemID
	NumberOfItems uint16
	ErrorCode     ErrorCode
}

func (r *SetServerItemReq) Service() byte {
	return MainService
}

func (r *SetServerItemReq) SubService() byte {
	return SubServiceSetServerItemReq
}

func (r *SetServerItemReq) Valid() error {
	if len(r.Items) == 0 {
		return fmt.Errorf("%w request must contain at least one item", ErrInvalid)
	}
	return nil
}

func (r *SetServerItemReq) Equal(other *SetServerItemReq) bool {
	return slices.EqualFunc(r.Items, other.Items, func(i1, i2 ServerItem) bool {
		return i1.Equal(&i2)
	})
}

func (r *SetServerItemReq) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString("setServerItemReq[items=[")
	for i, it := range r.Items {
		buffer.WriteString(it.String())
		if i != len(r.Items)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("]]")
	return buffer.String()
}

func (r *SetServerItemReq) MarshalBinary() ([]byte, error) {
	result := make([]byte, 6+sizeInBytes(r.Items))
	result[0] = MainService
	result[1] = SubServiceSetServerItemReq
	startItem := ItemID(0)
	if len(r.Items) > 0 {
		startItem = r.Items[0].ID
	}
	binary.BigEndian.PutUint16(result[2:4], uint16(startItem))
	binary.BigEndian.PutUint16(result[4:6], uint16(len(r.Items)))
	if err := marshalItemsTo(result[6:], r.Items); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *SetServerItemReq) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) < 6:
		return fmt.Errorf("%w data length %d, expected at least %d", ErrInvalid, len(data), 6)
	case data[0] != MainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], MainService)
	case data[1] != SubServiceSetServerItemReq:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], SubServiceSetServerItemReq)
	}
	num := int(binary.BigEndian.Uint16(data[4:6]))
	items, err := unmarshalItemsFrom(data[6:], num)
	if err != nil {
		return err
	}
	r.Items = items
	return nil
}

func (r *SetServerItemRes) Service() byte {
	return MainService
}

func (r *SetServerItemRes) SubService() byte {
	return SubServiceSetServerItemRes
}

func (r *SetServerItemRes) Valid() error {
	switch {
	case r.NumberOfItems == 0 && r.ErrorCode == 0:
		return fmt.Errorf("%w response must contain at least one item", ErrInvalid)
	case r.NumberOfItems != 0 && r.ErrorCode != 0:

		return fmt.Errorf("%w response contains %d items and error code is %d", ErrInvalid, r.NumberOfItems,
			r.ErrorCode)
	}
	return nil
}

func (r *SetServerItemRes) Equal(other *SetServerItemRes) bool {
	return r.StartItem == other.StartItem && r.NumberOfItems == other.NumberOfItems &&
		byte(r.ErrorCode) == byte(other.ErrorCode)
}

func (r *SetServerItemRes) String() string {
	return fmt.Sprintf("setServerItemRes[startItem=%d,numberOfItems=%d,errorCode=%d]", r.StartItem,
		r.NumberOfItems, r.ErrorCode)
}

func (r *SetServerItemRes) MarshalBinary() ([]byte, error) {
	result := make([]byte, 7)
	result[0] = MainService
	result[1] = SubServiceSetServerItemRes
	binary.BigEndian.PutUint16(result[2:4], uint16(r.StartItem))
	binary.BigEndian.PutUint16(result[4:6], r.NumberOfItems)
	result[6] = byte(r.ErrorCode)
	return result, nil
}

func (r *SetServerItemRes) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) != 7:
		return fmt.Errorf("%w data length %d, expected to be %d", ErrInvalid, len(data), 7)
	case data[0] != MainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], MainService)
	case data[1] != SubServiceSetServerItemRes:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], SubServiceSetServerItemRes)
	}
	r.StartItem = ItemID(binary.BigEndian.Uint16(data[2:4]))
	r.NumberOfItems = binary.BigEndian.Uint16(data[4:6])
	r.ErrorCode = ErrorCode(data[6])
	return nil
}
