package id

import (
	"fmt"

	"golang.org/x/crypto/sha3"
)

type HashType int

const (
	XXH32 HashType = iota
	SHA3
)

var DefaultHash = XXH32

func Hash(c interface{}) Output { return HashAs(DefaultHash, c) }

// TODO: Want to add add/switch to xxhash
func HashAs(hashType HashType, c interface{}) Output {
	switch hashType {
	case SHA3:
		return sha3Hash(c)
	default: // XXH32
		return xxh32Hash(c)
	}
}

func xxh32Hash(c interface{}) Output {
	xxhash32 := NewXXHash32()
	h := xxhash32.Sum([]byte(fmt.Sprintf("%v", c)))
	return Output{
		stringValue:    fmt.Sprintf("%x", h),
		byteSliceValue: h,
	}

}

func sha3Hash(c interface{}) Output {
	h := make([]byte, 64)
	sha3.ShakeSum256(h, []byte(fmt.Sprintf("%v", c)))
	return Output{
		stringValue:    fmt.Sprintf("%x", h),
		byteSliceValue: h,
	}

}