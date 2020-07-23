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
	case '{':
		return nil, s.errorSyntax("'{' inside [...]")
	case '$':
		return nil, s.errorSyntax("'$' inside [...]")
	case ':':
		s.patternBuffStart = uint(len(s.patternBuff))
		return lexRangeColon, nil
	case '-':
		return lexRangeDash, nil
	case ']':
		s.openBracket--
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
		s.lastGen = &charClassGenerator{
			charClasses: [][]rune{charset},
		}
		err := s.lastGen.Generate(s)
		if err != nil {
			return nil, err
		}
		s.patternBuff = nil
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

func lexRangeDash(s *State) (LexType, error) {
	if s.end() {
		return nil, s.errorSyntax("'[' not closed")
	}
	c1 := s.pattern[s.patternPos]
	s.move(1)
	if s.end() {
		return nil, s.errorSyntax("no character after '-'")
	}
	n := len(s.patternBuff)
	if n < 1 {
		return nil, s.errorSyntax("no character before '-'")
	}
	c0 := s.patternBuff[n-1]
	for b := int(c0) + 1; b <= int(c1); b++ {
		s.patternBuff = append(s.patternBuff, rune(b))
	}
	return lexRange, nil
}

func lexRangeBackslash(s *State) (LexType, error) {
	c := s.pattern[s.patternPos]
	s.move(1)
	s.patternBuff = append(s.patternBuff, backslashEscape(c))
	return lexRange, nil
}
