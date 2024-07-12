package passgen

import (
	math_rand "math/rand/v2"
	"slices"
	"strconv"
	"strings"
)

const maxRepeatCount = 1 << 28

func lexRepeat(s *State) (LexType, error) {
	if s.end() {
		return nil, s.errorSyntax("'{' not closed")
	}
	c := s.input[s.inputPos]
	s.move(1)
	switch c {
	case '[':
		return nil, s.errorSyntax("'[' inside {...}")
	case '{':
		return nil, s.errorSyntax("nested '{'")
	case '$':
		return nil, s.errorSyntax("'$' inside {...}")
	case ',':
		if slices.Contains(s.buffer, ',') {
			return nil, s.errorSyntax("multiple ',' inside {...}")
		}
		s.buffer = append(s.buffer, c)
		return lexRepeat, nil
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.buffer = append(s.buffer, c)
		return lexRepeat, nil
	case '-':
		if slices.Contains(s.buffer, ',') {
			return nil, s.errorSyntax("invalid natural number")
		}
		return nil, s.errorSyntax("repetition range syntax is '{M,N}' not '{M-N}'")
	case '}':
		return closeLexRepeat(s)
	}
	s.errorMarkLen = len(s.buffer) + 1
	return nil, s.errorSyntax("invalid natural number inside {...}")
}

func closeLexRepeat(s *State) (LexType, error) {
	if len(s.buffer) == 0 {
		return nil, s.errorSyntax("missing number inside {}")
	}
	if s.lastGen == nil {
		// I don't know how to test this without calling lexRepeat directly
		return nil, s.errorSyntax("nothing to repeat")
	}
	child := s.lastGen
	// FIXME: lastGen may have used another state
	count, err := parseRepeatCount(s, s.buffer)
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
			// I don't know how to test this without calling lexRepeat directly
			return nil, err
		}
	}
	gen.count = count
	// we set the gen.count to count-1 initially, because we don't want to
	// undo adding the characters we already have added to output
	// but we need to re-set g.count after gen.Generate(s), because the whole thing
	// might be repeated again
	s.lastGen = gen
	s.buffer = nil
	return LexRoot, nil
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
			return 0, s.errorSyntax(s_invalid_natural_num, countStr)
		}
		if countI64 < 1 {
			s.errorOffset--
			s.errorMarkLen = len(countStr)
			return 0, s.errorSyntax(s_invalid_natural_num, countStr)
		}
		if countI64 > maxRepeatCount {
			return 0, s.errorSyntax("count value is too large")
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
		// I don't know how to produce this by high-level Generate test
		s.errorOffset -= int64(len(maxStr)) + 2
		s.errorMarkLen = len(minStr)
		return 0, s.errorSyntax(s_invalid_natural_num, minStr)
	}
	if minCount < 1 {
		s.errorOffset -= int64(len(maxStr)) + 2
		s.errorMarkLen = len(minStr)
		return 0, s.errorSyntax(s_invalid_natural_num, minStr)
	}
	maxCount, err := strconv.ParseInt(maxStr, 10, 64)
	if err != nil {
		// I don't know how to produce this by high-level Generate test
		s.errorOffset -= 2
		s.errorMarkLen = len(maxStr)
		return 0, s.errorSyntax(s_invalid_natural_num, maxStr)
	}
	if maxCount < minCount {
		s.errorOffset--
		s.errorMarkLen = len(countRunes)
		return 0, s.errorValue("invalid numbers %v > %v inside {...}", minCount, maxCount)
	}
	if maxCount > maxRepeatCount {
		return 0, s.errorSyntax("count value is too large")
	}
	return minCount + math_rand.Int64N(maxCount-minCount+1), nil
}
