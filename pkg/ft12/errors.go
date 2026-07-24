package ft12

import "fmt"

// Sentinel errors for the FT1.2 package.
var (
	// ErrInvalid is returned when a frame or header is malformed.
	ErrInvalid = fmt.Errorf("invalid ft1.2")
	// ErrNoSpace is returned when the read buffer is too small.
	ErrNoSpace = fmt.Errorf("no space to read")
)
