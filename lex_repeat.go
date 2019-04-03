package main

import (
	math_rand "math/rand"
	"strconv"
	"strings"
)

type repeatGenerator struct {
	child Generator
	count int
	level int
}

func (g *repeatGenerator) Generate(s *State) error {
	child := g.child
	count := g.count
	for i := 0; i < count; i++ {
		err := child.Generate(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *repeatGenerator) Level() int {
	return g.level
}

func lexRepeat(s *State) (LexType, error) {
	if s.end() {
		return lexNil, s.errorSyntax("'{' not closed")
	}
	c := s.pattern[s.patternPos]
	s.patternPos++
	switch c {
	case '[':
		return lexNil, s.errorSyntax("'[' inside {...}")
	case '{':
		return lexNil, s.errorSyntax("nested '{'")
	case '$':
		return lexNil, s.errorSyntax("'$' inside {...}")
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.patternBuff = append(s.patternBuff, c)
		return lexRepeat, nil
	case '}':
		if len(s.patternBuff) == 0 {
			return lexNil, s.errorSyntax("missing number inside {}")
		}
		countStr := string(s.patternBuff)
		count := 0
		if strings.Contains(countStr, "-") {
			if countStr[0] == '-' {
				return lexNil, s.errorSyntax("no number before '-'")
			}
			if countStr[len(countStr)-1] == '-' {
				return lexNil, s.errorSyntax("no number after '-'")
			}
			parts := strings.Split(countStr, "-")
			if len(parts) > 2 {
				return lexNil, s.errorSyntax("multiple '-' inside {...}")
			} else if len(parts) < 2 {
				return lexNil, s.errorUnknown("unexpected error near '-' inside {...}")
			}
			minStr := parts[0]
			maxStr := parts[1]
			minCount, err := strconv.ParseInt(minStr, 10, 64)
			if err != nil {
				return lexNil, s.errorValue("invalid number %v inside {...}", minCount)
			}
			if minCount < 1 {
				return lexNil, s.errorValue("invalid number %v inside {...}", minCount)
			}
			maxCount, err := strconv.ParseInt(maxStr, 10, 64)
			if err != nil {
				return lexNil, s.errorValue("invalid number %v inside {...}", maxCount)
			}
			if maxCount < minCount {
				return lexNil, s.errorValue("invalid numbers %v > %v inside {...}", minCount, maxCount)
			}
			count = int(minCount) + math_rand.Intn(int(maxCount-minCount+1))
		} else {
			countI64, err := strconv.ParseInt(countStr, 10, 64)
			if err != nil {
				return lexNil, s.errorValue("invalid number '%v' inside {...}", countStr)
			}
			count = int(countI64)
			if count < 1 {
				return lexNil, s.errorValue("invalid number '%v' inside {...}", countStr)
			}
		}
		if s.lastGen == nil {
			return lexNil, s.errorSyntax("nothing to repeat")
		}
		gen := &repeatGenerator{
			child: s.lastGen,
			count: count - 1,
			level: s.lastGen.Level() + 1,
		}
		err := gen.Generate(s)
		if err != nil {
			return lexNil, err
		}
		gen.count = count
		// we set the gen.count to count-1 initially, because we don't want to
		// undo adding the characters we already have added to output
		// but we need to re-set g.count after gen.Generate(s), because the whole thing
		// might be repeated again
		s.lastGen = gen
		s.patternBuff = nil
		return LexRoot, nil
	}
	return lexNil, s.errorValue("non-numeric character '%v' inside {...}", string(c))
}
