package main

import (
	"fmt"
)

// SharedState is the shared part of State
type SharedState struct {
	absPos uint

	patternEntropy float64
}

// State is lex inputs, output and temp state
type State struct {
	*SharedState

	pattern []rune

	patternPos       uint
	patternBuff      []rune
	patternBuffStart uint

	openParenth uint
	openBracket uint

	lastGen generatorIface

	output []rune
}

func (s *State) move(chars uint) {
	s.patternPos += chars
	s.absPos += chars
}

func (s *State) addOutputOne(c rune) {
	s.lastGen = &staticStringGenerator{str: []rune{c}}
	s.lastGen.Generate(s)
}

func (s *State) addOutputNonRepeatable(data []rune) {
	s.lastGen = nil
	s.output = append(s.output, data...)
}

func (s *State) end() bool {
	return s.patternPos >= uint(len(s.pattern))
}

func (s *State) errorSyntax(msg string, args ...interface{}) error {
	return NewError(ErrorSyntax, s.absPos-1, fmt.Sprintf(msg, args...))
}

func (s *State) errorValue(msg string, args ...interface{}) error {
	return NewError(ErrorValue, s.absPos-1, fmt.Sprintf(msg, args...))
}

func (s *State) errorUnknown(msg string, args ...interface{}) error {
	return NewError(ErrorUnknown, s.absPos-1, fmt.Sprintf(msg, args...))
}

// NewState is factory function for State
func NewState(ss *SharedState, pattern string) *State {
	s := &State{
		SharedState: ss,
		pattern:     []rune(pattern),
	}
	return s
}
