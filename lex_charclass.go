package main

// FIXME: `[\^abc]` becomes `[^abc]` and excludes `abc`

func lexRange(s *State) (LexType, error) {
	if s.end() {
		return nil, s.errorSyntax("'[' not closed")
	}
	c := s.pattern[s.patternPos]
	s.move(1)
	switch c {
	case '\\':
		return lexRangeBackslash, nil
	case '[':
		return nil, s.errorSyntax("nested '['")
	case ':':
		s.patternBuffStart = uint(len(s.patternBuff))
		return lexRangeColon, nil
	case '-':
		return lexRangeDashInit, nil
	case ']':
		s.openBracket = false
		charset := s.patternBuff
		reverse := false
		if len(charset) > 0 && charset[0] == '^' {
			reverse = true
			charset = charset[1:]
		}
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
		s.patternBuff = nil
		s.lastGen = gen
		return LexRoot, nil
	}
	s.patternBuff = append(s.patternBuff, c)
	return lexRange, nil
}

func lexRangeColon(s *State) (LexType, error) {
	if s.end() {
		return nil, s.errorSyntax("':' not closed")
	}
	n := uint(len(s.patternBuff))
	// "[:digit:]"  -->  c.patternBuffStart == 0
	// "[abc:digit:]"  -->  c.patternBuffStart == 3
	c := s.pattern[s.patternPos]
	s.move(1)
	switch c {
	case ':':
		name := string(s.patternBuff[s.patternBuffStart:n])
		charset, ok := charClasses[name]
		if !ok {
			s.errorOffset -= len(name)
			return nil, s.errorValue("invalid character class %#v", name)
		}
		s.patternBuff = append(s.patternBuff[:s.patternBuffStart], charset...)
		return lexRange, nil
	case ']':
		return nil, s.errorSyntax("':' not closed")
	}
	s.patternBuff = append(s.patternBuff, c)
	return lexRangeColon, nil
}

func lexRangeDashInit(s *State) (LexType, error) {
	if s.end() {
		return nil, s.errorSyntax("'[' not closed")
	}
	s.patternBuff = append(s.patternBuff, s.pattern[s.patternPos-1], s.pattern[s.patternPos])
	s.move(1)
	if s.end() {
		return nil, s.errorSyntax("no character after '-'")
	}
	//s.patternBuff = append(s.patternBuff, s.pattern[s.patternPos])
	//s.move(1)
	return lexRangeDash, nil
}

func lexRangeDash(s *State) (LexType, error) {
	n := len(s.patternBuff)
	if n < 3 {
		return nil, s.errorSyntax("no character before '-'")
	}
	c1 := s.patternBuff[n-1]
	if c1 == '\\' {
		s.patternBuff = s.patternBuff[:n-1]
		return lexRangeDashBackslash, nil
	}
	if s.patternBuff[n-2] != '-' {
		return nil, s.errorUnknown("expected '-', buffer=%#v", string(s.patternBuff))
	}
	c0 := s.patternBuff[n-3]
	s.patternBuff = s.patternBuff[:n-2]
	for b := int(c0) + 1; b <= int(c1); b++ {
		s.patternBuff = append(s.patternBuff, rune(b))
	}
	return lexRange, nil
}

func lexRangeBackslash(s *State) (LexType, error) {
	c := s.pattern[s.patternPos]
	s.move(1)
	if c == 'u' {
		s.patternBuffStart2 = uint(len(s.patternBuff))
		s.patternBuff = append(s.patternBuff, '\\', 'u')
		return lexRangeUnicode, nil
	}
	s.patternBuff = append(s.patternBuff, backslashEscape(c))
	return lexRange, nil
}

func lexRangeDashBackslash(s *State) (LexType, error) {
	c := s.pattern[s.patternPos]
	s.move(1)
	if c == 'u' {
		s.patternBuffStart2 = uint(len(s.patternBuff))
		s.patternBuff = append(s.patternBuff, '\\', 'u')
		return lexRangeDashUnicode, nil
	}
	s.patternBuff = append(s.patternBuff, backslashEscape(c))
	return lexRangeDash, nil
}

func lexRangeUnicode(s *State) (LexType, error) {
	c := s.pattern[s.patternPos]
	s.move(1)
	s.patternBuff = append(s.patternBuff, c)
	start := int(s.patternBuffStart2)
	if len(s.patternBuff)-start == 6 {
		_, char, err := unescapeUnicodeSingle(s.patternBuff, start)
		if err != nil {
			s.errorOffset -= 5
			return nil, s.errorSyntax("invalid escape sequence")
		}
		s.patternBuff = append(s.patternBuff[:start], char)
		s.patternBuffStart2 = 0
		return lexRange, nil
	}
	return lexRangeUnicode, nil
}

func lexRangeDashUnicode(s *State) (LexType, error) {
	c := s.pattern[s.patternPos]
	s.move(1)
	s.patternBuff = append(s.patternBuff, c)
	start := int(s.patternBuffStart2)
	if len(s.patternBuff)-start == 6 {
		_, char, err := unescapeUnicodeSingle(s.patternBuff, start)
		if err != nil {
			s.errorOffset -= 5
			return nil, s.errorSyntax("invalid escape sequence")
		}
		s.patternBuff = append(s.patternBuff[:start], char)
		s.patternBuffStart2 = 0
		return lexRangeDash, nil
	}
	return lexRangeDashUnicode, nil
}
