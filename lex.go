package main

import "strconv"

// LexType is the type for lex functions
type LexType func(*State) (LexType, error)

var (
	lexRootUnicode     LexType
	lexRootUnicodeWide LexType
)

func init() {
	lexRootUnicode = makeLexUnicode(LexRoot, 'u', 6, false)
	lexRootUnicodeWide = makeLexUnicode(LexRoot, 'U', 10, false)
}

// LexRoot is the root lex implementation
func LexRoot(s *State) (LexType, error) {
	if s.patternBuff != nil {
		return nil, s.errorUnknown("incomplete buffer: %s", string(s.patternBuff))
	}
	if s.end() {
		if s.openParenth > 0 {
			return nil, s.errorSyntax("unexpected: unclosed '('")
		}
		return nil, nil
	}
	c := s.pattern[s.patternPos]
	s.move(1)
	switch c {
	case '\\':
		return lexBackslash, nil
	case '[':
		if s.openBracket {
			return nil, s.errorSyntax("nested '['")
		}
		s.openBracket = true
		return lexRange, nil
	case '{':
		if s.lastGen == nil {
			return nil, s.errorSyntax("nothing to repeat")
		}
		return lexRepeat, nil
	case '(':
		s.openParenth++
		s.lastGroupId++
		return lexGroup, nil
	case '$':
		return lexIdent, nil
	}
	s.addOutputOne(c)
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
	if c >= '1' && c <= '9' {
		return processGroupRef(s, LexRoot)
	}
	s.move(1)
	if c == 'u' {
		if s.patternBuff != nil {
			return nil, s.errorUnknown("incomplete buffer: %s", string(s.patternBuff))
		}
		return lexRootUnicode, nil
	}
	if c == 'U' {
		if s.patternBuff != nil {
			return nil, s.errorUnknown("incomplete buffer: %s", string(s.patternBuff))
		}
		return lexRootUnicodeWide, nil
	}
	s.addOutputOne(backslashEscape(c))
	return LexRoot, nil
}

func lexIdent(s *State) (LexType, error) {
	if s.end() {
		s.errorOffset++
		return nil, s.errorSyntax("expected a function call")
	}
	c := s.pattern[s.patternPos]
	s.move(1)
	switch c {
	case '\\', '[', '{', '$':
		return nil, s.errorSyntax("expected a function call")
	case '(':
		s.openParenth++
		return lexIdentFuncCall, nil
	}
	s.patternBuff = append(s.patternBuff, c)
	return lexIdent, nil
}

func makeLexUnicode(parentLex LexType, symbol rune, width int, toBuff bool) LexType {
	return func(s *State) (LexType, error) {
		buff := make([]rune, 0, width)
		buff = append(buff, '\\', symbol)
		for ; len(buff) < width && !s.end(); s.move(1) {
			buff = append(buff, s.pattern[s.patternPos])
		}
		if len(buff) != width {
			s.errorOffset -= len(buff) - 1
			return nil, s.errorSyntax("invalid escape sequence")
		}
		char, _, _, err := strconv.UnquoteChar(string(buff), '"')
		if err != nil {
			s.errorOffset -= width - 1
			return nil, s.errorSyntax("invalid escape sequence")
		}
		if toBuff {
			s.patternBuff = append(s.patternBuff, char)
		} else {
			s.addOutputOne(char)
		}
		return parentLex, nil
	}
}
