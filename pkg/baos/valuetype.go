package baos

import "fmt"

// ValueType describes the size and format of a datapoint value.
type ValueType byte

// SizeString returns the size of the value type in bits, or -1 for unknown types.
func SizeString(t ValueType) string {
	switch t {
	case 0, 1, 2, 3, 4, 5, 6:
		return fmt.Sprintf("%d bits", int(t)+1)
	case 7, 8, 9, 10:
		return fmt.Sprintf("%d bytes", int(t)-6)
	case 11:
		return fmt.Sprintf("%d bytes", 6)
	case 12:
		return fmt.Sprintf("%d bytes", 8)
	case 13:
		return fmt.Sprintf("%d bytes", 10)
	case 14:
		return fmt.Sprintf("%d bytes", 14)
	default:
		return "unknown size"
	}
}
