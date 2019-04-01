package main

import (
	"fmt"
	"io"
	"math/rand"
	"strconv"
)

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

// LexType is the type for lex functions
type LexType func(*State) (LexType, error)

// LexRoot is the root lex implementation
func LexRoot(s *State) (LexType, error) {
	if s.patternBuff != nil {
		return lexNil, fmt.Errorf("incomplete buffer: %s", string(s.patternBuff))
	}
	if s.end() {
		return lexNil, io.EOF
	}
	c := s.pattern[s.patternPos]
	s.patternPos++
	switch c {
	case '\\':
		return lexBackslash, nil
	case '[':
		return lexRange, nil
	case '{':
		return lexCount, nil
	case '$':
		return lexIdent, nil
	}
	s.addOutput(c)
	return LexRoot, nil
}

func lexBackslash(s *State) (LexType, error) {
	c := s.pattern[s.patternPos]
	s.patternPos++
	s.addOutput(c)
	return LexRoot, nil
}

func lexRangeBackslash(s *State) (LexType, error) {
	c := s.pattern[s.patternPos]
	s.patternPos++
	s.patternBuff = append(s.patternBuff, c)
	return lexRange, nil
}

func lexRangeDash(s *State) (LexType, error) {
	if s.end() {
		return lexNil, fmt.Errorf("'[' not closed")
	}
	c1 := s.pattern[s.patternPos]
	s.patternPos++
	if s.end() {
		return lexNil, fmt.Errorf("no character after '-'")
	}
	n := len(s.patternBuff)
	if n < 1 {
		return lexNil, fmt.Errorf("no character before '-'")
	}
	c0 := s.patternBuff[n-1]
	for b := int(c0); b <= int(c1); b++ {
		s.patternBuff = append(s.patternBuff, rune(b))
	}
	return lexRange, nil
}

func lexRangeColon(s *State) (LexType, error) {
	if s.end() {
		return lexNil, fmt.Errorf("'[' not closed")
	}
	n := uint(len(s.patternBuff))
	// "[:digit:]"  -->  c.patternBuffStart == 0
	// "[abc:digit:]"  -->  c.patternBuffStart == 3
	c := s.pattern[s.patternPos]
	s.patternPos++
	switch c {
	case ':':
		name := string(s.patternBuff[s.patternBuffStart:n])
		charset, ok := charsets[name]
		if !ok {
			return lexNil, fmt.Errorf("invalid charset %#v", name)
		}
		s.patternBuff = append(s.patternBuff[:s.patternBuffStart], []rune(charset)...)
		return lexRange, nil
	}
	s.patternBuff = append(s.patternBuff, c)
	return lexRangeColon, nil
}

func lexRange(s *State) (LexType, error) {
	if s.end() {
		return lexNil, fmt.Errorf("'[' not closed")
	}
	c := s.pattern[s.patternPos]
	s.patternPos++
	switch c {
	case '[':
		return lexRange, fmt.Errorf("nested '['")
	case '{':
		return lexCount, fmt.Errorf("'{' inside [...]")
	case '$':
		return lexIdent, fmt.Errorf("'$' inside [...]")
	case '\\':
		return lexRangeBackslash, nil
	case ':':
		s.patternBuffStart = uint(len(s.patternBuff))
		return lexRangeColon, nil
	case '-':
		return lexRangeDash, nil
	case ']':
		s.addRandomOutput(s.patternBuff)
		s.patternBuff = nil
		return LexRoot, nil
	}
	s.patternBuff = append(s.patternBuff, c)
	return lexRange, nil
}

func lexCount(s *State) (LexType, error) {
	if s.end() {
		return lexNil, fmt.Errorf("'{' not closed")
	}
	c := s.pattern[s.patternPos]
	s.patternPos++
	switch c {
	case '[':
		return lexNil, fmt.Errorf("'[' inside {...}")
	case '{':
		return lexNil, fmt.Errorf("nested '{'")
	case '$':
		return lexNil, fmt.Errorf("'$' inside {...}")
	case '}':
		count, err := strconv.ParseInt(string(s.patternBuff), 10, 64)
		if err != nil {
			return lexNil, fmt.Errorf("non-numeric string inside {...}")
		}
		if s.lastCharset == nil {
			return lexNil, fmt.Errorf("nothing to repeat")
		}
		if count > 1 {
			for i := int64(0); i < count-1; i++ {
				s.addRandomOutput(s.lastCharset)
			}
		}
		s.patternBuff = nil
		return LexRoot, nil
	}
	s.patternBuff = append(s.patternBuff, c)
	return lexCount, nil
}

func lexIdent(s *State) (LexType, error) {
	return lexNil, fmt.Errorf("'$' not implemented yet")
}

func lexNil(s *State) (LexType, error) {
	s.patternPos++
	return lexNil, nil
}
