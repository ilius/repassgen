package main

import (
	"io"
	math_rand "math/rand"
	"strconv"
	"strings"
)

// LexType is the type for lex functions
type LexType func(*State) (LexType, error)

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
		return lexCount, nil
	case '$':
		return lexIdent, nil
	}
	s.addOutput(c)
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
	s.addOutput(backslashEscape(c))
	return LexRoot, nil
}

func lexRangeBackslash(s *State) (LexType, error) {
	c := s.pattern[s.patternPos]
	s.patternPos++
	s.patternBuff = append(s.patternBuff, backslashEscape(c))
	return lexRange, nil
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

func lexRange(s *State) (LexType, error) {
	if s.end() {
		return lexNil, s.errorSyntax("'[' not closed")
	}
	c := s.pattern[s.patternPos]
	s.patternPos++
	switch c {
	case '[':
		return lexRange, s.errorSyntax("nested '['")
	case '{':
		return lexCount, s.errorSyntax("'{' inside [...]")
	case '$':
		return lexIdent, s.errorSyntax("'$' inside [...]")
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
		return lexCount, nil
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
		if s.lastCharset == nil {
			return lexNil, s.errorSyntax("nothing to repeat")
		}
		for i := 0; i < count-1; i++ {
			s.addRandomOutput(s.lastCharset)
		}
		s.patternBuff = nil
		return LexRoot, nil
	}
	return lexNil, s.errorValue("non-numeric character '%v' inside {...}", string(c))
}

func lexIdentParen(s *State) (LexType, error) {
	if s.end() {
		return lexNil, s.errorSyntax("'(' not closed")
	}
	n := uint(len(s.patternBuff))
	// "$a()"  -->  c.patternBuffStart == 1
	c := s.pattern[s.patternPos]
	s.patternPos++
	switch c {
	case ')':
		funcName := string(s.patternBuff[:s.patternBuffStart])
		if funcName == "" {
			return lexNil, s.errorSyntax("missing function name")
		}
		funcObj, ok := functions[funcName]
		if !ok {
			return lexNil, s.errorValue("invalid function '%v'", funcName)
		}
		argPattern := string(s.patternBuff[s.patternBuffStart:n])
		argOut := Generate(GenerateInput{
			Pattern:            argPattern,
			CalcPatternEntropy: s.calcPatternEntropy,
		})
		result, err := funcObj(argOut.Password)
		if err != nil {
			lexErr, ok := err.(*LexError)
			if ok {
				lexErr.MovePos(int(s.patternBuffStart))
				lexErr.PrependMsg("function " + funcName)
				return lexNil, lexErr
			}
			return lexNil, s.errorUnknown("%v returned error: %v", funcName, err)
		}
		s.addOutputNonRepeatable(result)
		if s.calcPatternEntropy {
			s.patternEntropy += argOut.PatternEntropy
		}
		s.patternBuff = nil
		return LexRoot, nil
	}
	s.patternBuff = append(s.patternBuff, c)
	return lexIdentParen, nil
}

func lexIdent(s *State) (LexType, error) {
	if s.end() {
		return lexNil, s.errorSyntax("'(' not closed")
	}
	c := s.pattern[s.patternPos]
	s.patternPos++
	switch c {
	case '[', '{', '$':
		return lexRange, s.errorSyntax("expected a function call after $")
	case '\\':
		return lexRangeBackslash, nil
	case '(':
		s.patternBuffStart = uint(len(s.patternBuff))
		return lexIdentParen, nil
	}
	s.patternBuff = append(s.patternBuff, c)
	return lexIdent, nil
}

func lexNil(s *State) (LexType, error) {
	s.patternPos++
	return lexNil, nil
}
