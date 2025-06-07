package signer

import (
	"encoding/hex"
	"fmt"
	"strings"
)

func toHex(data []byte) string {
	return fmt.Sprintf("%#x", data)
}

func hexDecodeString(s string) ([]byte, error) {
	s = strings.TrimPrefix(s, "0x")

	if len(s)%2 != 0 {
		s = "0" + s
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func u8aToU8a(value interface{}) []byte {
	if value == nil {
		return make([]byte, 0)
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return bytes
}

func u8aConcat(list [][]byte) []byte {
	length := 0
	offset := 0
	u8as := make([][]byte, len(list))
	for i, u8a := range list {
		u8as[i] = u8aToU8a(u8a)
		length += len(u8as[i])
	}

	result := make([]byte, length)

	for i, u8a := range list {
		copy(result[i:], u8a)
		offset += len(u8as[i])
	}
	return result
}
