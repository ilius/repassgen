package passgen

import (
	"crypto/rand"
	"encoding/binary"
	math_rand "math/rand/v2"
)

// CryptoRandSource is a source for math/rand that uses more secure crypto/rand
type CryptoRandSource struct{}

// NewRandSource creates a new source for math/rand that uses more secure crypto/rand
func NewRandSource() *math_rand.Rand {
	return math_rand.New(CryptoRandSource{})
}

func (CryptoRandSource) Uint64() uint64 {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err)
	}
	return binary.LittleEndian.Uint64(b[:])
}
