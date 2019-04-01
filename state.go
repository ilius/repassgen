package main

import "math/rand"

// State is lex inputs, output and temp state
type State struct {
	pattern          []rune
	patternPos       uint
	patternBuff      []rune
	patternBuffStart uint
	patternRange     [2]rune
	lastCharset      []rune
	output           []rune
}

func (s *State) addOutput(c rune) {
	s.lastCharset = []rune{c}
	s.output = append(s.output, c)
}

func (s *State) addRandomOutput(charset []rune) {
	s.lastCharset = charset
	if len(charset) == 0 {
		return
	}
	i := rand.Intn(len(charset))
	s.output = append(s.output, charset[i])
}

func (s *State) end() bool {
	return s.patternPos >= uint(len(s.pattern))
}

// NewState is factory function for State
func NewState(pattern string) *State {
	return &State{
		pattern: []rune(pattern),
	}
}
