package baos

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"
)

type ServerItem struct {
	ID   ItemID
	Data []byte
}

const serverItemHeaderSize = int(unsafe.Sizeof(ItemID(0))) + 1

func (i *ServerItem) String() string {
	return fmt.Sprintf("serverItem[id=%d,data=%x]", i.ID, i.Data)
}

func (i *ServerItem) Equal(other *ServerItem) bool {
	return i.ID == other.ID && bytes.Equal(i.Data, other.Data)
}

func sizeInBytes(items []ServerItem) int {
	totalSize := 0
	for _, item := range items {
		totalSize += serverItemHeaderSize + len(item.Data)
	}
	return totalSize
}

func marshalItemsTo(data []byte, items []ServerItem) error {
	for _, it := range items {
		if len(data) < serverItemHeaderSize+len(it.Data) {
			return fmt.Errorf("%w buffer has no capacity: %d while %d needed", ErrInvalid, len(it.Data),
				serverItemHeaderSize+len(it.Data))
		}
		binary.BigEndian.PutUint16(data[0:2], uint16(it.ID))
		data[2] = byte(len(it.Data))
		copy(data[serverItemHeaderSize:], it.Data)
		data = data[serverItemHeaderSize+len(it.Data):]
	}
	if len(data) != 0 {
		return fmt.Errorf("%w buffer has %d extra bytes", ErrInvalid, len(data))
	}
	return nil
}

func unmarshalItemsFrom(data []byte, n int) ([]ServerItem, error) {
	result := make([]ServerItem, 0, n)
	for i := 0; i < n; i++ {
		if len(data) < serverItemHeaderSize {
			return nil, fmt.Errorf("%w buffer has %d bytes while %d needed", ErrInvalid,
				len(data), serverItemHeaderSize)
		}
		id := ItemID(binary.BigEndian.Uint16(data[0:2]))
		length := int(data[2])
		if len(data) < serverItemHeaderSize+length {
			return nil, fmt.Errorf("%w buffer has only %d bytes while %d needed", ErrInvalid,
				len(data), serverItemHeaderSize+length)
		}
		result = append(result, ServerItem{ID: id, Data: data[serverItemHeaderSize : serverItemHeaderSize+length]})
		data = data[serverItemHeaderSize+length:]
	}
	if len(data) != 0 {
		return nil, fmt.Errorf("%w buffer has %d extra bytes", ErrInvalid, len(data))
	}
	return result, nil
}
