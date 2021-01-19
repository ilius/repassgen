package main

import (
	"fmt"
)

// SharedState is the shared part of State
type SharedState struct {
	absPos uint

	errorOffset int

	patternEntropy float64
}

// State is lex inputs, output and temp state
type State struct {
	*SharedState
	lastGen generatorIface

	pattern     []rune
	patternBuff []rune
	output      []rune

	patternPos uint

	patternBuffStart  uint
	patternBuffStart2 uint

	openParenth uint
	openBracket bool
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

func (s *State) getErrorPos() uint {
	pos := int(s.absPos) + s.errorOffset - 1
	if pos < 0 {
		fmt.Printf("Warning: getErrorPos: pos=%v\n", pos)
		pos = 0
	}
	return uint(pos)
}

func (s *State) errorSyntax(msg string, args ...interface{}) error {
	return NewError(
		ErrorSyntax,
		s.getErrorPos(),
		fmt.Sprintf(msg, args...),
	)
}

func (s *State) errorArg(msg string, args ...interface{}) error {
	return NewError(
		ErrorArg,
		s.getErrorPos(),
		fmt.Sprintf(msg, args...),
	)
}

func (s *State) errorValue(msg string, args ...interface{}) error {
	return NewError(
		ErrorValue,
		s.getErrorPos(),
		fmt.Sprintf(msg, args...),
	)
}

func (s *State) errorUnknown(msg string, args ...interface{}) error {
	return NewError(
		ErrorUnknown,
		s.getErrorPos(),
		fmt.Sprintf(msg, args...),
	)
}

// NewState is factory function for State
func NewState(ss *SharedState, pattern []rune) *State {
	s := &State{
		SharedState: ss,
		pattern:     pattern,
	}
	return s
}
