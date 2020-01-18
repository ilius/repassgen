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
		s.openBracket++
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
	s.addOutputOne(backslashEscape(c))
	return LexRoot, nil
}

func lexIdent(s *State) (LexType, error) {
	if s.end() {
		if s.openParenth > 0 {
			return nil, s.errorSyntax("'(' not closed")
		}
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
