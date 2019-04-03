package main

import (
	rand "crypto/rand"
	"math"
	"math/big"
)

type charsetGroupGenerator struct {
	charsets [][]rune
}

func (g *charsetGroupGenerator) Generate(s *State) error {
	for _, charset := range g.charsets {
		if len(charset) == 0 {
			continue
		}
		ibig, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}
		i := int(ibig.Int64())
		s.output = append(s.output, charset[i])
		if s.calcPatternEntropy {
			s.patternEntropy += math.Log2(float64(len(charset)))
		}
	}
	return nil
}

func (g *charsetGroupGenerator) Level() int {
	return 0
}

func lexRange(s *State) (LexType, error) {
	if s.end() {
		return lexNil, s.errorSyntax("'[' not closed")
	}
	c := s.pattern[s.patternPos]
	s.patternPos++
	switch c {
	case '\\':
		return lexRangeBackslash, nil
	case '[':
		return lexNil, s.errorSyntax("nested '['")
	case '{':
		return lexNil, s.errorSyntax("'{' inside [...]")
	case '$':
		return lexNil, s.errorSyntax("'$' inside [...]")
	case ':':
		s.patternBuffStart = uint(len(s.patternBuff))
		return lexRangeColon, nil
	case '-':
		return lexRangeDash, nil
	case ']':
		s.lastGen = &charsetGroupGenerator{
			charsets: [][]rune{s.patternBuff},
		}
		err := s.lastGen.Generate(s)
		if err != nil {
			return lexNil, err
		}
		s.patternBuff = nil
		return LexRoot, nil
	}
	s.patternBuff = append(s.patternBuff, c)
	return lexRange, nil
}

func lexRangeColon(s *State) (LexType, error) {
	if s.end() {
		return lexNil, s.errorSyntax("':' not closed")
	}
	n := uint(len(s.patternBuff))
	// "[:digit:]"  -->  c.patternBuffStart == 0
	// "[abc:digit:]"  -->  c.patternBuffStart == 3
	c := s.pattern[s.patternPos]
	s.patternPos++
	switch c {
	case ':':
		name := string(s.patternBuff[s.patternBuffStart:n])
		charset, ok := charClasses[name]
		if !ok {
			return lexNil, s.errorValue("invalid character class %#v", name)
		}
		s.patternBuff = append(s.patternBuff[:s.patternBuffStart], charset...)
		return lexRange, nil
	case ']':
		return lexNil, s.errorSyntax("':' not closed")
	}
	s.patternBuff = append(s.patternBuff, c)
	return lexRangeColon, nil
}

func lexRangeDash(s *State) (LexType, error) {
	if s.end() {
		return lexNil, s.errorSyntax("'[' not closed")
	}
	c1 := s.pattern[s.patternPos]
	s.patternPos++
	if s.end() {
		return lexNil, s.errorSyntax("no character after '-'")
	}
	n := len(s.patternBuff)
	if n < 1 {
		return lexNil, s.errorSyntax("no character before '-'")
	}
	c0 := s.patternBuff[n-1]
	for b := int(c0) + 1; b <= int(c1); b++ {
		s.patternBuff = append(s.patternBuff, rune(b))
	}
	return lexRange, nil
}

func lexRangeBackslash(s *State) (LexType, error) {
	c := s.pattern[s.patternPos]
	s.patternPos++
	s.patternBuff = append(s.patternBuff, backslashEscape(c))
	return lexRange, nil
}
