package baos

import (
	"encoding/binary"
	"fmt"
	"slices"
)

type ServerItemInd struct {
	StartItem ItemID
	Items     []ServerItem
}

func (r *ServerItemInd) Service() byte {
	return MainService
}

func (r *ServerItemInd) SubService() byte {
	return SubServiceServerItemInd
}

func (r *ServerItemInd) Valid() error {
	if len(r.Items) == 0 {
		return fmt.Errorf("%w indication must contain at least one item", ErrInvalid)
	}
	return nil
}

func (r *ServerItemInd) Equal(other *ServerItemInd) bool {
	return r.StartItem == other.StartItem && slices.EqualFunc(r.Items, other.Items, func(i1, i2 ServerItem) bool {
		return i1.Equal(&i2)
	})
}

func (r *ServerItemInd) String() string {
	return fmt.Sprintf("serverItemInd[startItem=%d,items=%d]", r.StartItem, len(r.Items))
}

func (r *ServerItemInd) MarshalBinary() ([]byte, error) {
	result := make([]byte, 6+sizeInBytes(r.Items))
	result[0] = byte(MainService)
	result[1] = byte(SubServiceServerItemInd)
	binary.BigEndian.PutUint16(result[2:4], uint16(r.StartItem))
	binary.BigEndian.PutUint16(result[4:6], uint16(len(r.Items)))
	if err := marshalItemsTo(result[6:], r.Items); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ServerItemInd) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) < 6:
		return fmt.Errorf("%w data length %d, expected at least %d", ErrInvalid, len(data), 6)
	case data[0] != byte(MainService):
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], MainService)
	case data[1] != byte(SubServiceServerItemInd):
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], SubServiceServerItemInd)
	}
	startItem := ItemID(binary.BigEndian.Uint16(data[2:4]))
	num := int(binary.BigEndian.Uint16(data[4:6]))
	items, err := unmarshalItemsFrom(data[6:], num)
	if err != nil {
		return err
	}
	*r = ServerItemInd{StartItem: startItem, Items: items}
	return nil
}
