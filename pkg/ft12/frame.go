package ft12

import (
	"bytes"
	"fmt"
)

type frame []byte

type header [4]byte

func (h *header) len() int {
	return int((*h)[1])
}

func (h *header) valid() error {
	if (*h)[0] != 0x68 && (*h)[3] != 0x68 && (*h)[1] != (*h)[2] {
		return fmt.Errorf("invalid header: %x", (*h)[:])
	}
	return nil
}

func (h *header) isResetIndication() bool {
	return bytes.Equal((*h)[:], []byte{0x10, 0xC0, 0xC0, 0x16})
}

func marshal(data []byte) ([]frame, error) {
	return nil, nil
}

func unmarshal(data []byte) ([]byte, error) {
	return nil, nil
}
