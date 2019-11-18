package events

import (
	"encoding/binary"
)

func SimpleSerialization(data interface{}) (value []byte) {
	switch data.(type) {
	case uint32:
		value = make([]byte, 4)
		binary.LittleEndian.PutUint32(value, data.(uint32))
	case uint64:
		value = make([]byte, 8)
		binary.LittleEndian.PutUint64(value, data.(uint64))
	case int64:
		value = make([]byte, 8)
		binary.LittleEndian.PutUint64(value, uint64(data.(int64)))
	case string:
		value = []byte(data.(string))
	}

	return
}

func ReadUint32(data []byte) uint32 {
	return binary.LittleEndian.Uint32(data)
}

func ReadUint64(data []byte) uint64 {
	return binary.LittleEndian.Uint64(data)
}

func ReadInt64(data []byte) int64 {
	return int64(binary.LittleEndian.Uint64(data))
}
