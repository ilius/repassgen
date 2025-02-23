package passgen

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

// SharedState is the shared part of State
type SharedState struct {
	groupsOutput   map[uint64][]rune
	absPos         uint64
	errorOffset    int64
	errorMarkLen   int
	patternEntropy float64
	lastGroupId    uint64

	maxOutputLength int
}

func (ss *SharedState) Copy() *SharedState {
	new_ss := *ss
	return &new_ss
}

// State is lex inputs, output and temp state
type State struct {
	*SharedState
	lastGen generatorIface

	input  []rune
	buffer []rune
	output []rune

	inputPos uint64

	openParenth uint64
	openBracket bool

	rangeReverse bool
}

func (s *State) move(chars uint64) {
	s.inputPos += chars
	s.absPos += chars
}

func (s *State) moveBack(chars uint64) {
	s.inputPos -= chars
	if s.absPos < chars {
		log.Printf("moveBack(%v) with absPos=%v", chars, s.absPos)
		return
	}
	s.absPos -= chars
}

func (s *State) addOutputOne(c rune) error {
	if s.tooLong() {
		s.lastGen = nil
		return nil
	}
	s.lastGen = &staticStringGenerator{str: []rune{c}}
	return s.lastGen.Generate(s)
}

func (s *State) addOutput(str []rune) error {
	s.lastGen = &staticStringGenerator{str: str}
	return s.lastGen.Generate(s)
}

func (s *State) addOutputNonRepeatable(data []rune) {
	s.lastGen = nil
	s.output = append(s.output, data...)
}

func (s *State) tooLong() bool {
	return s.maxOutputLength > 0 && len(s.output) > s.maxOutputLength
}

func (s *State) end() bool {
	if s.tooLong() {
		return true
	}
	return s.inputPos >= uint64(len(s.input))
}

func (s *State) getErrorPos() uint {
	if s.absPos == 0 {
		if s.errorOffset < 0 {
			log.Printf("errorOffset = %v < 0, pattern: %v", s.errorOffset, string(s.input))
			return 0
		}
		return uint(s.errorOffset)
	}
	pos := int64(s.absPos) + s.errorOffset - 1
	if pos < 0 {
		log.Printf("errorOffset=%v, absPos=%v, pos=%v, pattern: %v", s.errorOffset, s.absPos, pos, string(s.input))
		return 0
	}
	return uint(pos)
}

func (s *State) errorSyntax(msg string, args ...any) error {
	return NewError(
		ErrorSyntax,
		s.getErrorPos(),
		fmt.Sprintf(msg, args...),
	).WithMarkLen(s.errorMarkLen)
}

func (s *State) errorArg(msg string, args ...any) error {
	return NewError(
		ErrorArg,
		s.getErrorPos(),
		fmt.Sprintf(msg, args...),
	).WithMarkLen(s.errorMarkLen)
}

func (s *State) errorValue(msg string, args ...any) error {
	return NewError(
		ErrorValue,
		s.getErrorPos(),
		fmt.Sprintf(msg, args...),
	).WithMarkLen(s.errorMarkLen)
}

func (s *State) errorUnknown(msg string, args ...any) error {
	return NewError(
		ErrorUnknown,
		s.getErrorPos(),
		fmt.Sprintf(msg, args...),
	).WithMarkLen(s.errorMarkLen)
}

// NewSharedState is factory function for SharedState
func NewSharedState() *SharedState {
	ss := &SharedState{
		groupsOutput: map[uint64][]rune{},
		errorMarkLen: 1,
	}
	maxLengthStr := os.Getenv("REPASSGEN_MAX_LENGTH")
	if maxLengthStr != "" {
		maxLength, err := strconv.ParseInt(maxLengthStr, 10, 64)
		if err != nil {
			panic("invalid REPASSGEN_MAX_LENGTH: must be an integer")
		}
		if maxLength > math.MaxInt {
			panic("REPASSGEN_MAX_LENGTH is larger than MaxInt")
		}
		ss.maxOutputLength = int(maxLength)
	}
	return ss
}

// NewState is factory function for State
func NewState(ss *SharedState, pattern []rune) *State {
	s := &State{
		SharedState: ss,
		input:       pattern,
	}
	return s
}
