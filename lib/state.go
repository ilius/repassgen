package passgen

import (
	"fmt"
	"log"
)

// SharedState is the shared part of State
type SharedState struct {
	groupsOutput   map[uint64][]rune
	absPos         uint64
	errorOffset    int64
	errorMarkLen   int
	patternEntropy float64
	lastGroupId    uint64
}

// State is lex inputs, output and temp state
type State struct {
	*SharedState
	lastGen generatorIface

	pattern     []rune
	patternBuff []rune
	output      []rune

	patternPos uint64

	openParenth uint64
	openBracket bool

	rangeReverse bool
}

func (s *State) move(chars uint64) {
	s.patternPos += chars
	s.absPos += chars
}

func (s *State) moveBack(chars uint64) {
	s.patternPos -= chars
	if s.absPos < chars {
		log.Printf("moveBack(%v) with absPos=%v", chars, s.absPos)
		return
	}
	s.absPos -= chars
}

// func (s *State) moveBackAbs(chars uint64) {
// 	if s.absPos < chars {
// 		log.Printf("moveBack(%v) with absPos=%v", chars, s.absPos)
// 		return
// 	}
// 	s.absPos -= chars
// }

func (s *State) addOutputOne(c rune) {
	s.lastGen = &staticStringGenerator{str: []rune{c}}
	s.lastGen.Generate(s)
}

func (s *State) addOutput(str []rune) {
	s.lastGen = &staticStringGenerator{str: str}
	s.lastGen.Generate(s)
}

func (s *State) addOutputNonRepeatable(data []rune) {
	s.lastGen = nil
	s.output = append(s.output, data...)
}

func (s *State) end() bool {
	return s.patternPos >= uint64(len(s.pattern))
}

func (s *State) getErrorPos() uint {
	if s.absPos == 0 {
		if s.errorOffset < 0 {
			log.Printf("errorOffset = %v < 0, pattern: %v", s.errorOffset, string(s.pattern))
			return 0
		}
		return uint(s.errorOffset)
	}
	pos := int64(s.absPos) + s.errorOffset - 1
	if pos < 0 {
		// errorOffset=0, absPos=18446744073709551615, pos=-2, pattern: $\(\(\(
		log.Printf("errorOffset=%v, absPos=%v, pos=%v, pattern: %v", s.errorOffset, s.absPos, pos, string(s.pattern))
		return 0
	}
	return uint(pos)
}

func (s *State) errorSyntax(msg string, args ...interface{}) error {
	return NewError(
		ErrorSyntax,
		s.getErrorPos(),
		fmt.Sprintf(msg, args...),
	).WithMarkLen(s.errorMarkLen)
}

func (s *State) errorArg(msg string, args ...interface{}) error {
	return NewError(
		ErrorArg,
		s.getErrorPos(),
		fmt.Sprintf(msg, args...),
	).WithMarkLen(s.errorMarkLen)
}

func (s *State) errorValue(msg string, args ...interface{}) error {
	return NewError(
		ErrorValue,
		s.getErrorPos(),
		fmt.Sprintf(msg, args...),
	).WithMarkLen(s.errorMarkLen)
}

func (s *State) errorUnknown(msg string, args ...interface{}) error {
	return NewError(
		ErrorUnknown,
		s.getErrorPos(),
		fmt.Sprintf(msg, args...),
	).WithMarkLen(s.errorMarkLen)
}

// NewSharedState is factory function for SharedState
func NewSharedState() *SharedState {
	return &SharedState{
		groupsOutput: map[uint64][]rune{},
		errorMarkLen: 1,
	}
}

// NewState is factory function for State
func NewState(ss *SharedState, pattern []rune) *State {
	s := &State{
		SharedState: ss,
		pattern:     pattern,
	}
	return s
}
