package main

import (
	"fmt"
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

	lastGen generatorIface

	output []rune
}

func (s *State) addOutputOne(c rune) error {
	s.lastGen = &staticStringGenerator{str: []rune{c}}
	return s.lastGen.Generate(s)
}

func (s *State) addOutputNonRepeatable(data []rune) error {
	s.lastGen = nil
	s.output = append(s.output, data...)
	return nil
}

func (s *State) end() bool {
	return s.patternPos >= uint(len(s.pattern))
}

func (s *State) errorSyntax(msg string, args ...interface{}) error {
	return NewError(LexErrorSyntax, s.patternPos-1, fmt.Sprintf(msg, args...))
}

func (s *State) errorValue(msg string, args ...interface{}) error {
	return NewError(LexErrorValue, s.patternPos-1, fmt.Sprintf(msg, args...))
}

func (s *State) errorUnknown(msg string, args ...interface{}) error {
	return NewError(LexErrorUnknown, s.patternPos-1, fmt.Sprintf(msg, args...))
}

// NewState is factory function for State
func NewState(pattern string, calcPatternEntropy bool) *State {
	s := &State{
		pattern:            []rune(pattern),
		calcPatternEntropy: calcPatternEntropy,
	}
	return s
}
