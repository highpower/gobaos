package knxtypes

import "errors"

var (
	ErrUnmarshal   = errors.New("unmarshal error")
	ErrUnknownType = errors.New("unknown type")
	ErrBadValue    = errors.New("bad value")
)
