package baos

import (
	"encoding/binary"
	"fmt"
	"slices"
)

// ServerItemInd is an unsolicited indication carrying server items.
type ServerItemInd struct {
	StartItem ItemID
	Items     []ServerItem
}

// Service returns the main service identifier.
func (r *ServerItemInd) Service() byte {
	return MainService
}

// SubService returns the sub-service identifier for the indication.
func (r *ServerItemInd) SubService() byte {
	return SubServiceServerItemInd
}

// Equal reports whether r and other are identical indications.
func (r *ServerItemInd) Equal(other *ServerItemInd) bool {
	return r.StartItem == other.StartItem && slices.EqualFunc(r.Items, other.Items, func(i1, i2 ServerItem) bool {
		return i1.Equal(&i2)
	})
}

// String returns a human-readable representation of the indication.
func (r *ServerItemInd) String() string {
	return fmt.Sprintf("serverItemInd[startItem=%d,items=%d]", r.StartItem, len(r.Items))
}

// MarshalBinary encodes the indication into its BAOS wire format.
func (r *ServerItemInd) MarshalBinary() ([]byte, error) {
	result := make([]byte, 6+sizeInBytes(r.Items))
	result[0] = MainService
	result[1] = SubServiceServerItemInd
	binary.BigEndian.PutUint16(result[2:4], uint16(r.StartItem))
	binary.BigEndian.PutUint16(result[4:6], uint16(len(r.Items)))
	if err := marshalItemsTo(result[6:], r.Items); err != nil {
		return nil, err
	}
	return result, nil
}

// UnmarshalBinary decodes the BAOS wire format into the indication.
func (r *ServerItemInd) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) < 6:
		return fmt.Errorf("%w data length %d, expected at least %d", ErrInvalid, len(data), 6)
	case data[0] != MainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], MainService)
	case data[1] != SubServiceServerItemInd:
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
