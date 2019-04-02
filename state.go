package main

import (
	rand "crypto/rand"
	"math"
	"math/big"
)

// State is lex inputs, output and temp state
type State struct {
	pattern []rune

	calcPatternEntropy bool
	// patternEntropy is zero unless -entropy flag is given
	patternEntropy float64

	patternPos       uint
	patternBuff      []rune
	patternBuffStart uint
	lastCharset      []rune

	output []rune
}

func (s *State) addOutput(c rune) {
	s.lastCharset = []rune{c}
	s.output = append(s.output, c)
}

func (s *State) addOutputNonRepeatable(data []rune) {
	s.lastCharset = nil
	s.output = append(s.output, data...)
}

func (s *State) addRandomOutput(charset []rune) {
	s.lastCharset = charset
	if len(charset) == 0 {
		return
	}
	ibig, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
	if err != nil {
		panic(err)
	}
	i := int(ibig.Int64())
	s.output = append(s.output, charset[i])
	if s.calcPatternEntropy {
		s.patternEntropy += math.Log2(float64(len(charset)))
	}
}

func (s *State) end() bool {
	return s.patternPos >= uint(len(s.pattern))
}

// NewState is factory function for State
func NewState(pattern string, calcPatternEntropy bool) *State {
	s := &State{
		pattern:            []rune(pattern),
		calcPatternEntropy: calcPatternEntropy,
	}
	return s
}
