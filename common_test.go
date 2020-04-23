package main

import (
	"github.com/ilius/is/v2"
)

func strPtr(s string) *string {
	s2 := s
	return &s2
}

func isFloatBetween(is *is.Is, actual float64, min float64, max float64) {
	is.AddMsg(
		"%v is not in range [%v, %v]",
		actual,
		min,
		max,
	).True(min <= actual && actual <= max)
}

type genCase struct {
	Pattern string

	Error string

	PassLen  [2]int // {min, max}
	Validate func(string) bool
	Entropy  [2]float64 // {min, max}

	Password *string

	WordCount int

	// TODO: CharClassCount map[string]int
}
