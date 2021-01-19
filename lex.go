package main

// LexType is the type for lex functions
type LexType func(*State) (LexType, error)

// LexRoot is the root lex implementation
func LexRoot(s *State) (LexType, error) {
	if s.patternBuff != nil {
		return nil, s.errorUnknown("incomplete buffer: %s", string(s.patternBuff))
	}
	if s.end() {
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
		return lexRepeat, nil
	case '(':
		s.openParenth++
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
	s.move(1)
	if c == 'u' {
		if s.patternBuff != nil {
			return nil, s.errorUnknown("incomplete buffer: %s", string(s.patternBuff))
		}
		s.patternBuff = []rune(`\u`)
		return lexUnicode, nil
	}
	s.addOutputOne(backslashEscape(c))
	return LexRoot, nil
}

func lexUnicode(s *State) (LexType, error) {
	c := s.pattern[s.patternPos]
	s.move(1)
	s.patternBuff = append(s.patternBuff, c)
	if len(s.patternBuff) == 6 {
		_, char, err := unescapeUnicodeSingle(s.patternBuff, 0)
		if err != nil {
			s.errorOffset = -5
			return nil, s.errorSyntax("invalid escape sequence")
		}
		s.addOutputOne(char)
		s.patternBuff = nil
		return LexRoot, nil
	}
	return lexUnicode, nil
}

func lexIdent(s *State) (LexType, error) {
	if s.end() {
		return nil, s.errorSyntax("expected a function call")
	}
	c := s.pattern[s.patternPos]
	s.move(1)
	switch c {
	case '\\', '[', '{', '$':
		return nil, s.errorSyntax("expected a function call")
	case '(':
		s.patternBuffStart = uint(len(s.patternBuff))
		s.openParenth++
		return lexIdentFuncCall, nil
	}
	s.patternBuff = append(s.patternBuff, c)
	return lexIdent, nil
}
