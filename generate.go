package main

import (
	"fmt"
	"io"
)

// Generate generates random password based on given pattern
// see README.md for examples of pattern
func Generate(pattern string) []rune {
	s := NewState(pattern)
	lex := LexRoot
	var err error
	for {
		lex, err = lex(s)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(fmt.Errorf("invalid syntax near index %d: %v", s.patternPos-1, err))
		}
	}
	return s.output
}
