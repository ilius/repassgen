package main

import (
	rand "crypto/rand"
	"math"
	"math/big"
)

type charsetGroupGenerator struct {
	charsets [][]rune
	entropy  *float64
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
	}
	entropy, err := g.Entropy()
	if err != nil {
		return err
	}
	s.patternEntropy += entropy
	return nil
}

func (g *charsetGroupGenerator) Level() int {
	return 0
}

func (g *charsetGroupGenerator) Entropy() (float64, error) {
	if g.entropy != nil {
		return *g.entropy, nil
	}
	entropy := 0.0
	for _, charset := range g.charsets {
		entropy += math.Log2(float64(len(charset)))
	}
	g.entropy = &entropy
	return entropy, nil
}

func lexRange(s *State) (LexType, error) {
	if s.end() {
		return nil, s.errorSyntax("'[' not closed")
	}
	c := s.pattern[s.patternPos]
	s.patternPos++
	switch c {
	case '\\':
		return lexRangeBackslash, nil
	case '[':
		return nil, s.errorSyntax("nested '['")
	case '{':
		return nil, s.errorSyntax("'{' inside [...]")
	case '$':
		return nil, s.errorSyntax("'$' inside [...]")
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
			return nil, err
		}
		s.patternBuff = nil
		return LexRoot, nil
	}
	s.patternBuff = append(s.patternBuff, c)
	return lexRange, nil
}

func lexRangeColon(s *State) (LexType, error) {
	if s.end() {
		return nil, s.errorSyntax("':' not closed")
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
			return nil, s.errorValue("invalid character class %#v", name)
		}
		s.patternBuff = append(s.patternBuff[:s.patternBuffStart], charset...)
		return lexRange, nil
	case ']':
		return nil, s.errorSyntax("':' not closed")
	}
	s.patternBuff = append(s.patternBuff, c)
	return lexRangeColon, nil
}

func lexRangeDash(s *State) (LexType, error) {
	if s.end() {
		return nil, s.errorSyntax("'[' not closed")
	}
	c1 := s.pattern[s.patternPos]
	s.patternPos++
	if s.end() {
		return nil, s.errorSyntax("no character after '-'")
	}
	n := len(s.patternBuff)
	if n < 1 {
		return nil, s.errorSyntax("no character before '-'")
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
