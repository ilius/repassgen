package main

import (
	math_rand "math/rand"
	"strconv"
	"strings"
)

func lexRepeat(s *State) (LexType, error) {
	if s.end() {
		return nil, s.errorSyntax("'{' not closed")
	}
	c := s.pattern[s.patternPos]
	s.move(1)
	switch c {
	case '[':
		return nil, s.errorSyntax("'[' inside {...}")
	case '{':
		return nil, s.errorSyntax("nested '{'")
	case '$':
		return nil, s.errorSyntax("'$' inside {...}")
	case ',':
		if hasRune(s.patternBuff, ',') {
			return nil, s.errorSyntax("multiple ',' inside {...}")
		}
		s.patternBuff = append(s.patternBuff, c)
		return lexRepeat, nil
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.patternBuff = append(s.patternBuff, c)
		return lexRepeat, nil
	case '-':
		return nil, s.errorSyntax("repetition range syntax is '{M,N}' not '{M-N}'")
	case '}':
		if len(s.patternBuff) == 0 {
			return nil, s.errorSyntax("missing number inside {}")
		}
		if s.lastGen == nil {
			return nil, s.errorSyntax("nothing to repeat")
		}
		child := s.lastGen
		// FIXME: lastGen may have used another state
		count, err := parseRepeatCount(s, s.patternBuff)
		if err != nil {
			return nil, err
		}
		gen := &repeatGenerator{
			child: child,
			count: count - 1,
		}
		{
			err = gen.Generate(s)
			if err != nil {
				return nil, err
			}
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
	s.errorMarkLen = len(s.patternBuff) + 1
	return nil, s.errorSyntax("invalid natural number inside {...}")
}

func parseRepeatCount(s *State, countRunes []rune) (int64, error) {
	countStr := string(countRunes)
	parts := strings.Split(countStr, ",")
	// we know that len(parts) >= 1
	if len(parts) > 2 {
		return 0, s.errorSyntax("multiple ',' inside {...}")
	}
	if len(parts) == 1 {
		countI64, err := strconv.ParseInt(countStr, 10, 64)
		if err != nil {
			s.errorOffset--
			s.errorMarkLen = len(countStr)
			return 0, s.errorSyntax("invalid natural number '%v'", countStr)
		}
		if countI64 < 1 {
			s.errorOffset--
			s.errorMarkLen = len(countStr)
			return 0, s.errorSyntax("invalid natural number '%v'", countStr)
		}
		return countI64, nil
	}
	// now we know len(parts) == 2
	if countStr[0] == ',' {
		s.errorOffset -= int64(len(countRunes))
		return 0, s.errorSyntax("no number before ','")
	}
	if countStr[len(countStr)-1] == ',' {
		return 0, s.errorSyntax("no number after ','")
	}
	minStr := parts[0]
	maxStr := parts[1]
	minCount, err := strconv.ParseInt(minStr, 10, 64)
	if err != nil {
		// s.errorMarkLen = len(minStr)
		return 0, s.errorSyntax("invalid natural number '%v'", minStr)
	}
	if minCount < 1 {
		// s.errorMarkLen = len(minStr)
		return 0, s.errorSyntax("invalid natural number '%v'", minStr)
	}
	maxCount, err := strconv.ParseInt(maxStr, 10, 64)
	if err != nil {
		// s.errorMarkLen = len(maxStr)
		return 0, s.errorSyntax("invalid natural number '%v'", maxCount)
	}
	if maxCount < minCount {
		s.errorOffset--
		s.errorMarkLen = len(countRunes)
		return 0, s.errorValue("invalid numbers %v > %v inside {...}", minCount, maxCount)
	}
	return minCount + math_rand.Int63n(maxCount-minCount+1), nil
}
