package main

import (
	rand "crypto/rand"
	"encoding/binary"
	math_rand "math/rand"
)

// CryptoRandSource is a source for math/rand that uses more secure crypto/rand
type CryptoRandSource struct{}

// NewRandSource creates a new source for math/rand that uses more secure crypto/rand
func NewRandSource() *math_rand.Rand {
	return math_rand.New(CryptoRandSource{})
}

// Int63 ...
func (CryptoRandSource) Int63() int64 {
	var b [8]byte
	rand.Read(b[:])
	// mask off sign bit to ensure positive number
	return int64(binary.LittleEndian.Uint64(b[:]) & (1<<63 - 1))
}

// Seed ...
func (CryptoRandSource) Seed(_ int64) {}
