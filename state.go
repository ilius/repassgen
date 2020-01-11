package main

import (
	"fmt"
)

type SharedState struct {
	patternEntropy float64
	output         []rune
}

// State is lex inputs, output and temp state
type State struct {
	*SharedState

	pattern []rune

	patternPos       uint
	patternBuff      []rune
	patternBuffStart uint

	openParenth uint

	lastGen generatorIface
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
func NewState(pattern string) *State {
	s := &State{
		SharedState: &SharedState{},
		pattern:     []rune(pattern),
	}
	return s
}
