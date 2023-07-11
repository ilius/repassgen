package passgen_test

import (
	"github.com/ilius/is/v2"
	passgen "github.com/ilius/repassgen/lib"
)

func strPtr(s string) *string {
	s2 := s
	return &s2
}

func newTestState(patternArg interface{}) *passgen.State {
	var pattern []rune
	switch patternTyped := patternArg.(type) {
	case string:
		pattern = []rune(patternTyped)
	case []rune:
		pattern = patternTyped
	default:
		panic("invalid patternArg")
	}
	s := passgen.NewState(passgen.NewSharedState(), pattern)
	return s
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

	PassLen  [2]int // {min, max}
	Validate func(string) bool
	Entropy  [2]float64 // {min, max}

	Password *string

	WordCount int

	// TODO: CharClassCount map[string]int
}

type genErrCase struct {
	Pattern string

	Error interface{}
}
