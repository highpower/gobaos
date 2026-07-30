package knxtypes

import (
	"fmt"
	"strconv"
	"strings"
)

// DPTAccess represents KNX DPT 15.xxx (access control).
// It encodes a 4-byte payload: Key (7 bit), ID (7 bit), Index (4 bit), AdditionalData (2 bit).
type DPTAccess struct {
	Key            uint8
	ID             uint8
	Index          uint8
	AdditionalData uint8
}

// Get returns the access fields in the order "key:id:index:additionalData"
// as colon-separated decimal values.
func (a *DPTAccess) Get() string {
	return fmt.Sprintf("%d:%d:%d:%d", a.Key, a.ID, a.Index, a.AdditionalData)
}

// Set parses a "key:id:index:additionalData" string, validates each field
// against its allowed range, and stores them. It returns ErrBadValue if
// the format is wrong or any value is out of range.
func (a *DPTAccess) Set(value string) error {
	parts := strings.Split(value, ":")
	if len(parts) != 4 {
		return fmt.Errorf("%w: type %s requires 4 parts, got %d", ErrBadValue, AccessCode.String(),
			len(parts))
	}
	parsed := make([]uint8, 4)
	for i, part := range parts {
		v, err := strconv.ParseUint(part, 10, 8)
		if err != nil {
			return err
		}
		parsed[i] = uint8(v)
	}
	if err := validAccess(parsed[0], parsed[1], parsed[2], parsed[3]); err != nil {
		return err
	}
	a.Key, a.ID, a.Index, a.AdditionalData = parsed[0], parsed[1], parsed[2], parsed[3]
	return nil
}

// TypeCode returns the KNX type code for DPTAccess.
func (a *DPTAccess) TypeCode() TypeCode {
	return AccessCode
}

// Valid reports whether the access fields are within their valid ranges.
func (a *DPTAccess) Valid() error {
	return validAccess(a.Key, a.ID, a.Index, a.AdditionalData)
}

// String returns a human-readable representation of DPTAccess.
func (a *DPTAccess) String() string {
	return fmt.Sprintf("dptAccess[key=%d,id=%d,index=%d,additionalData=%d]", a.Key, a.ID, a.Index,
		a.AdditionalData)
}

// MarshalBinary encodes DPTAccess into a 4-byte payload, masking any
// reserved upper bits to their allowed widths.
func (a *DPTAccess) MarshalBinary() ([]byte, error) {
	return []byte{a.Key & 0x7F, a.ID & 0x7F, a.Index & 0x0F, a.AdditionalData & 0x03}, nil
}

// UnmarshalBinary decodes a 4-byte payload into DPTAccess. It returns
// ErrUnmarshal if data is not exactly 4 bytes. Reserved upper bits are
// masked off during decoding.
func (a *DPTAccess) UnmarshalBinary(data []byte) error {
	if len(data) != 4 {
		return fmt.Errorf("%w: type %s requires 4 bytes, got %d", ErrUnmarshal, AccessCode.String(), len(data))
	}
	a.Key = data[0] & 0x7F
	a.ID = data[1] & 0x7F
	a.Index = data[2] & 0x0F
	a.AdditionalData = data[3] & 0x03
	return nil
}

func validAccess(key, id, index, additionalData uint8) error {
	if key > 127 {
		return fmt.Errorf("%w: key %d out of range [0,127]", ErrBadValue, key)
	}
	if id > 127 {
		return fmt.Errorf("%w: id %d out of range [0,127]", ErrBadValue, id)
	}
	if index > 15 {
		return fmt.Errorf("%w: index %d out of range [0,15]", ErrBadValue, index)
	}
	if additionalData > 3 {
		return fmt.Errorf("%w: additional data %d out of range [0,3]", ErrBadValue, additionalData)
	}
	return nil
}
