package main

import (
	math_rand "math/rand"
	"strconv"
	"strings"
)

const badRepeatCount = "invalid natural number inside {...}"

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
	case ',', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
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
		count, err := parseRepeatCount(s, string(s.patternBuff))
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
	return nil, s.errorSyntax(badRepeatCount)
}

func parseRepeatCount(s *State, countStr string) (int, error) {
	parts := strings.Split(countStr, ",")
	// we know that len(parts) >= 1
	if len(parts) > 2 {
		return 0, s.errorSyntax("multiple ',' inside {...}")
	}
	if len(parts) == 1 {
		countI64, err := strconv.ParseInt(countStr, 10, 64)
		if err != nil {
			return 0, s.errorSyntax(badRepeatCount)
		}
		if countI64 < 1 {
			return 0, s.errorSyntax(badRepeatCount)
		}
		return int(countI64), nil
	}
	// now we know len(parts) == 2
	if countStr[0] == ',' {
		return 0, s.errorSyntax("no number before ','")
	}
	if countStr[len(countStr)-1] == ',' {
		return 0, s.errorSyntax("no number after ','")
	}
	minStr := parts[0]
	maxStr := parts[1]
	minCount, err := strconv.ParseInt(minStr, 10, 64)
	if err != nil {
		return 0, s.errorSyntax(badRepeatCount)
	}
	if minCount < 1 {
		return 0, s.errorSyntax(badRepeatCount)
	}
	maxCount, err := strconv.ParseInt(maxStr, 10, 64)
	if err != nil {
		return 0, s.errorSyntax(badRepeatCount)
	}
	if maxCount < minCount {
		return 0, s.errorValue("invalid numbers %v > %v inside {...}", minCount, maxCount)
	}
	return int(minCount) + math_rand.Intn(int(maxCount-minCount+1)), nil
}
