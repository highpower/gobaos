package baos

type ValueType byte

func (t ValueType) SizeInBits() int {
	switch t {
	case 0, 1, 2, 3, 4, 5, 6:
		return int(t) + 1
	case 7, 8, 9, 10:
		return (int(t) - 6) * 8
	case 11:
		return 48
	case 12:
		return 64
	case 13:
		return 80
	case 14:
		return 112
	default:
		return -1
	}
}
