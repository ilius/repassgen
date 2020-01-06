package main

import (
	rand "crypto/rand"
	"encoding/binary"
	math_rand "math/rand"
)

type CryptoRandSource struct{}

func NewRandSource() *math_rand.Rand {
	return math_rand.New(CryptoRandSource{})
}

func (_ CryptoRandSource) Int63() int64 {
	var b [8]byte
	rand.Read(b[:])
	// mask off sign bit to ensure positive number
	return int64(binary.LittleEndian.Uint64(b[:]) & (1<<63 - 1))
}

func (_ CryptoRandSource) Seed(_ int64) {}
