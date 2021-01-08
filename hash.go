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

func Hash(c interface{}) Id { return HashAs(DefaultHash, c) }

// TODO: Want to add add/switch to xxhash
func HashAs(hashType HashType, c interface{}) Id {
	switch hashType {
	case SHA3:
		return sha3Hash(c)
	default: // XXH32
		return xxh32Hash(c)
	}
}

func xxh32Hash(c interface{}) Id {
	xxhash32 := NewXXHash32()
	h := xxhash32.Sum([]byte(fmt.Sprintf("%v", c)))
	return Id{
		stringValue:    fmt.Sprintf("%x", h),
		byteSliceValue: h,
		uint32Value:    u32(h),
	}

}

func sha3Hash(c interface{}) Id {
	h := make([]byte, 64)
	sha3.ShakeSum256(h, []byte(fmt.Sprintf("%v", c)))
	return Id{
		stringValue:    fmt.Sprintf("%x", h),
		byteSliceValue: h,
		uint32Value:    u32(h),
	}

}
