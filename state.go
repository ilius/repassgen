package main

import (
	"unicode/utf8"
	"bytes"
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
	patternBuff      []byte
	patternBuffStart uint

	openParenth uint
	openBracket uint

	lastGen generatorIface

	output *bytes.Buffer
}

func (s *State) move(chars uint) {
	s.patternPos += chars
	s.absPos += chars
}

func (s *State) addPatternBuffRune(c rune) {
	// FIXME
	n := utf8.RuneLen(c)
	cb := make([]byte, n)
	utf8.EncodeRune(cb, c)
	s.patternBuff = append(s.patternBuff, cb...)
}

func (s *State) addOutputOne(c rune) {
	// FIXME
	s.lastGen = &staticStringGenerator{str: []byte(string(c))}
	s.lastGen.Generate(s)
}

func (s *State) addOutputNonRepeatable(data []byte) {
	s.lastGen = nil
	s.output.Write(data)
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
		output:      bytes.NewBuffer(nil),
	}
	return s
}
