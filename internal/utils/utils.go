package utils

import (
	"github.com/mr-tron/base58"
)

// We use Base 58 encoding to make our addresses more human-readable
// This is because base 58 encoding removes characters that could be confused with each other
// For example, 0 and O, or 1 and l
func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)
	return []byte(encode)
}

func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	if err != nil {
		panic(err)
	}

	return decode
}
