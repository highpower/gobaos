package knxtypes

import "errors"

// Sentinel errors for the knxtypes package.
var (
	// ErrUnmarshal is returned when binary unmarshaling fails.
	ErrUnmarshal = errors.New("unmarshal error")
	// ErrUnknownType is returned when a type code string cannot be parsed.
	ErrUnknownType = errors.New("unknown type")
	// ErrBadValue is returned when a value is out of range or otherwise invalid.
	ErrBadValue = errors.New("bad value")
)
