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
		return lexRange, nil
	case '{':
		return lexRepeat, nil
	case '(':
		s.openParenth++
		return lexGroup, nil
	case '$':
		return lexIdent, nil
	}
	err := s.addOutputOne(c)
	if err != nil {
		return nil, err
	}
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
	err := s.addOutputOne(backslashEscape(c))
	if err != nil {
		return nil, err
	}
	return LexRoot, nil
}

func lexIdent(s *State) (LexType, error) {
	if s.end() {
		return nil, s.errorSyntax("'(' not closed")
	}
	c := s.pattern[s.patternPos]
	s.move(1)
	switch c {
	case '\\':
		return lexRangeBackslash, nil
	case '[', '{', '$':
		return lexRange, s.errorSyntax("expected a function call after $")
	case '(':
		s.patternBuffStart = uint(len(s.patternBuff))
		s.openParenth++
		return lexIdentFuncCall, nil
	}
	s.patternBuff = append(s.patternBuff, c)
	return lexIdent, nil
}
