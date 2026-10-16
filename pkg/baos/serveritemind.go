package baos

import (
	"encoding/binary"
	"fmt"
	"slices"
)

// serverItemInd is an unsolicited indication carrying server items.
type serverItemInd struct {
	startItem ItemID
	items     []ServerItem
}

// MarshalBinary encodes the indication into its BAOS wire format.
func (r *serverItemInd) MarshalBinary() ([]byte, error) {
	result := make([]byte, 6+sizeInBytes(r.items))
	result[0] = mainService
	result[1] = subServiceServerItemInd
	binary.BigEndian.PutUint16(result[2:4], uint16(r.startItem))
	binary.BigEndian.PutUint16(result[4:6], uint16(len(r.items)))
	if err := marshalItemsTo(result[6:], r.items); err != nil {
		return nil, err
	}
	return result, nil
}

// UnmarshalBinary decodes the BAOS wire format into the indication.
func (r *serverItemInd) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) < 6:
		return fmt.Errorf("%w data length %d, expected at least %d", ErrInvalid, len(data), 6)
	case data[0] != mainService:
		return fmt.Errorf("%w data service %d, expected %d", ErrInvalid, data[0], mainService)
	case data[1] != subServiceServerItemInd:
		return fmt.Errorf("%w data subservice %d, expected %d", ErrInvalid, data[1], subServiceServerItemInd)
	}
	startItem := ItemID(binary.BigEndian.Uint16(data[2:4]))
	num := int(binary.BigEndian.Uint16(data[4:6]))
	items, err := unmarshalItemsFrom(data[6:], num)
	if err != nil {
		return err
	}
	*r = serverItemInd{startItem: startItem, items: items}
	return nil
}

// equal reports whether r and other are identical indications.
func (r *serverItemInd) equal(other *serverItemInd) bool {
	return r.startItem == other.startItem && slices.EqualFunc(r.items, other.items, func(i1, i2 ServerItem) bool {
		return i1.Equal(&i2)
	})
}
