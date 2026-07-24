package baos

import "fmt"

// ErrorCode represents a BAOS protocol error code.
type ErrorCode byte

// Supported BAOS error codes.
const (
	NoError               ErrorCode = 0
	InternalError         ErrorCode = 1
	NoElementFound        ErrorCode = 2
	BufferIsTooSmall      ErrorCode = 3
	ItemIsNotWriteable    ErrorCode = 4
	ServiceIsNotSupported ErrorCode = 5
	BadServiceParameter   ErrorCode = 6
	BadID                 ErrorCode = 7
	BadCommandValue       ErrorCode = 8
	BadLength             ErrorCode = 9
	MessageInconsistent   ErrorCode = 10
	ObjectServerIsBusy    ErrorCode = 11
)

// ErrInvalid is returned when a message or value fails validation.
var ErrInvalid = fmt.Errorf("invalid")

// Byte returns the error code as a byte.
func (e ErrorCode) Byte() byte {
	return byte(e)
}

// Error returns a human-readable description of the error code.
func (e ErrorCode) Error() string {
	switch e {
	case NoError:
		return "no error"
	case InternalError:
		return "internal error"
	case NoElementFound:
		return "no element found"
	case BufferIsTooSmall:
		return "buffer is too small"
	case ItemIsNotWriteable:
		return "item is not writeable"
	case ServiceIsNotSupported:
		return "service is not supported"
	case BadServiceParameter:
		return "bad service parameter"
	case BadID:
		return "bad id"
	case BadCommandValue:
		return "bad command/value"
	case BadLength:
		return "bad length"
	case MessageInconsistent:
		return "message inconsistent"
	case ObjectServerIsBusy:
		return "object server is busy"
	}
	return "unknown error"
}
