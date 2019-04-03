package main

import (
	"io"
)

type Generator interface {
	Generate(s *State) error
	Level() int
}

// LexType is the type for lex functions
type LexType func(*State) (LexType, error)

func lexNil(s *State) (LexType, error) {
	s.patternPos++
	return lexNil, nil
}

// LexRoot is the root lex implementation
func LexRoot(s *State) (LexType, error) {
	if s.patternBuff != nil {
		return lexNil, s.errorUnknown("incomplete buffer: %s", string(s.patternBuff))
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
		return lexRepeat, nil
	// case '(': // TODO: repeating groups, like `([a-z][1-9]){5}``
	case '$':
		return lexIdent, nil
	}
	err := s.addOutput(c)
	if err != nil {
		return lexNil, err
	}
	return LexRoot, nil
}

func backslashEscape(c rune) rune {
	switch c {
	case 't':
		return '\t'
	case 'r':
		return '\r'
	case 'n':
		return '\n'
	case 'v':
		return '\v'
	case 'f':
		return '\f'
	}
	return c
}

func lexBackslash(s *State) (LexType, error) {
	c := s.pattern[s.patternPos]
	s.patternPos++
	err := s.addOutput(backslashEscape(c))
	if err != nil {
		return lexNil, err
	}
	return LexRoot, nil
}

func lexIdent(s *State) (LexType, error) {
	if s.end() {
		return lexNil, s.errorSyntax("'(' not closed")
	}
	c := s.pattern[s.patternPos]
	s.patternPos++
	switch c {
	case '\\':
		return lexRangeBackslash, nil
	case '[', '{', '$':
		return lexRange, s.errorSyntax("expected a function call after $")
	case '(':
		s.patternBuffStart = uint(len(s.patternBuff))
		return lexIdentFuncCall, nil
	}
	s.patternBuff = append(s.patternBuff, c)
	return lexIdent, nil
}
