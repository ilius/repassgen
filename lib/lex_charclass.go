package passgen

func processRange(s *State, charset []rune) (LexType, error) {
	reverse := s.rangeReverse
	s.openBracket = false
	s.rangeReverse = false
	charset = removeDuplicateRunes(charset)
	if reverse {
		charset = excludeCharsASCII(charset)
	}
	gen := &charClassGenerator{
		charClasses: [][]rune{charset},
	}
	err := gen.Generate(s)
	if err != nil {
		return nil, err
	}
	s.buffer = nil
	s.lastGen = gen
	return LexRoot, nil
}

func lexRange(s *State) (LexType, error) {
	if s.end() {
		s.errorOffset++
		return nil, s.errorSyntax("'[' not closed")
	}
	c := s.input[s.inputPos]
	s.move(1)
	switch c {
	case '\\':
		return lexRangeBackslash, nil
	case '[':
		return nil, s.errorSyntax("nested '['")
	case ':':
		return lexRangeColon, nil
	case '-':
		return lexRangeDashInit, nil
	case '^':
		if !s.rangeReverse && len(s.buffer) == 0 {
			s.rangeReverse = true
			return lexRange, nil
		}
	case ']':
		return processRange(s, s.buffer)
	}
	s.buffer = append(s.buffer, c)
	return lexRange, nil
}

func lexRangeColon(s *State) (LexType, error) {
	if s.end() {
		s.errorOffset++
		return nil, s.errorSyntax("':' not closed")
	}
	nameRunes := []rune{}
	for !s.end() {
		c := s.input[s.inputPos]
		s.move(1)
		switch c {
		case ':':
			name := string(nameRunes)
			charset, ok := charClasses[name]
			if !ok {
				s.errorMarkLen = len(name) + 2
				return nil, s.errorValue("invalid character class %#v", name)
			}
			s.buffer = append(s.buffer, charset...)
			return lexRange, nil
		case ']':
			s.errorMarkLen = len(nameRunes) + 2
			return nil, s.errorSyntax("':' not closed")
		}
		nameRunes = append(nameRunes, c)
	}
	s.errorOffset++
	s.errorMarkLen = len(nameRunes) + 2
	return nil, s.errorSyntax("':' not closed")
}

func lexRangeDashInit(s *State) (LexType, error) {
	if s.end() {
		s.errorOffset++
		s.errorMarkLen = len(s.buffer) + 3
		return nil, s.errorSyntax("'[' not closed")
	}
	s.buffer = append(s.buffer, s.input[s.inputPos-1], s.input[s.inputPos])
	s.move(1)
	if s.end() {
		return nil, s.errorSyntax("no character after '-'")
	}
	return lexRangeDash, nil
}

func lexRangeDash(s *State) (LexType, error) {
	n := len(s.buffer)
	if n < 3 {
		s.errorOffset--
		return nil, s.errorSyntax("no character before '-'")
	}
	c1 := s.buffer[n-1]
	if c1 == '\\' {
		s.buffer = s.buffer[:n-1]
		return lexRangeDashBackslash, nil
	}
	if s.buffer[n-2] != '-' {
		return nil, s.errorUnknown("expected '-', buffer=%#v", string(s.buffer))
	}
	c0 := s.buffer[n-3]
	s.buffer = s.buffer[:n-2]
	for b := int(c0) + 1; b <= int(c1); b++ {
		s.buffer = append(s.buffer, rune(b))
	}
	return lexRange, nil
}

func lexRangeBackslash(s *State) (LexType, error) {
	if s.end() {
		return nil, s.errorSyntax("'[' not closed")
	}
	c := s.input[s.inputPos]
	s.move(1)
	if c == 'u' {
		return makeLexUnicode(lexRange, 'u', 6, true), nil
	}
	if c == 'U' {
		return makeLexUnicode(lexRange, 'U', 10, true), nil
	}
	s.buffer = append(s.buffer, backslashEscape(c))
	return lexRange, nil
}

func lexRangeDashBackslash(s *State) (LexType, error) {
	if s.end() {
		return nil, s.errorSyntax("'[' not closed")
	}
	c := s.input[s.inputPos]
	s.move(1)
	if c == 'u' {
		return makeLexUnicode(lexRangeDash, 'u', 6, true), nil
	}
	if c == 'U' {
		return makeLexUnicode(lexRangeDash, 'U', 10, true), nil
	}
	s.buffer = append(s.buffer, backslashEscape(c))
	return lexRangeDash, nil
}
