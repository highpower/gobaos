package mask

import (
	"bytes"
)

const defaultMaskLength = 4

func Apply(s string) string {
	runes := []rune(s)
	buffer := bytes.Buffer{}
	switch {
	case len(runes) == 0:
		return ""
	case len(runes) <= defaultMaskLength:
		buffer.WriteRune(runes[0])
	default:
		for i := range defaultMaskLength {
			buffer.WriteRune(runes[i])
		}
	}
	buffer.WriteString("...")
	return buffer.String()
}
