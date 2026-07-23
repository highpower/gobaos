package baos

import "fmt"

type ErrorCode byte

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

var ErrInvalid = fmt.Errorf("invalid")

func (e ErrorCode) Byte() byte {
	return byte(e)
}

func (e ErrorCode) Error() string {
	switch byte(e) {
	case byte(NoError):
		return "no error"
	case byte(InternalError):
		return "internal error"
	case byte(NoElementFound):
		return "no element found"
	case byte(BufferIsTooSmall):
		return "buffer is too small"
	case byte(ItemIsNotWriteable):
		return "item is not writeable"
	case byte(ServiceIsNotSupported):
		return "service is not supported"
	case byte(BadServiceParameter):
		return "bad service parameter"
	case byte(BadID):
		return "bad id"
	case byte(BadCommandValue):
		return "bad command/value"
	case byte(BadLength):
		return "bad length"
	case byte(MessageInconsistent):
		return "message inconsistent"
	case byte(ObjectServerIsBusy):
		return "object server is busy"
	}
	return "unknown error"
}
