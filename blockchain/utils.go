package blockchain

import (
	"bytes"
	"encoding/binary"
)

func IntToBytes(n uint64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, n)
	if err != nil {
		panic(err)
	}
	return buff.Bytes()
}

func IsEqualToZeroHash(hash []byte) bool {
	if len(hash) != 32 {
		return false
	}
	for _, n := range hash {
		if n != 0 {
			return false
		}
	}
	return true
}
