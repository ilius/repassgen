package main

import (
	"fmt"

	"github.com/spf13/cast"
)

// SharedState is the shared part of State
type SharedState struct {
	absPos uint64

	errorOffset int64

	patternEntropy float64

	lastGroupId  uint64
	groupsOutput map[uint64][]rune

	// character probability
	charProbEnable bool
	charProbIScale float64 // inversed scale of values in charProbMap
	charProbMap    map[rune]float64
}

func (s *SharedState) applyCharProbIScale() {
	if len(s.charProbMap) == 0 {
		return
	}
	iscale := s.charProbIScale
	if iscale == 1.0 {
		return
	}
	m := s.charProbMap
	for r, p := range m {
		m[r] = p / iscale
	}
	s.charProbIScale = 1.0
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
	s.absPos -= chars
}

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
		pos, err := cast.ToUintE(s.errorOffset)
		if err != nil {
			panic(err)
		}
		return pos
	}
	pos := int64(s.absPos) + s.errorOffset - 1
	if pos < 0 {
		fmt.Printf("Warning: getErrorPos: pos=%v, pattern=%#v\n", pos, string(s.pattern))
		pos = 0
	}
	pos2, err := cast.ToUintE(pos)
	if err != nil {
		panic(err)
	}
	return pos2
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

// NewSharedState is factory function for SharedState
func NewSharedState() *SharedState {
	return &SharedState{
		groupsOutput: map[uint64][]rune{},
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
