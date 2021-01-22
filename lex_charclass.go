package main

// FIXME: `[\^abc]` becomes `[^abc]` and excludes `abc`

var (
	lexRangeUnicode         LexType
	lexRangeUnicodeWide     LexType
	lexRangeDashUnicode     LexType
	lexRangeDashUnicodeWide LexType
)

func init() {
	lexRangeUnicode = makeLexUnicode(lexRange, 'u', 6, true)
	lexRangeUnicodeWide = makeLexUnicode(lexRange, 'U', 10, true)
	lexRangeDashUnicode = makeLexUnicode(lexRangeDash, 'u', 6, true)
	lexRangeDashUnicodeWide = makeLexUnicode(lexRangeDash, 'U', 10, true)
}

func processRange(s *State, charset []rune, reverse bool) (LexType, error) {
	s.openBracket = false
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

func lexRange(s *State) (LexType, error) {
	if s.end() {
		s.errorOffset++
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
		charset := s.patternBuff
		reverse := false
		if len(charset) > 0 && charset[0] == '^' {
			reverse = true
			charset = charset[1:]
		}
		return processRange(s, charset, reverse)
	}
	s.patternBuff = append(s.patternBuff, c)
	return lexRange, nil
}

func lexRangeColon(s *State) (LexType, error) {
	if s.end() {
		s.errorOffset++
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
		s.errorOffset++
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
		s.errorOffset--
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
		return makeLexUnicode(lexRange, 'u', 6, true), nil
	}
	if c == 'U' {
		return makeLexUnicode(lexRange, 'U', 10, true), nil
	}
	s.patternBuff = append(s.patternBuff, backslashEscape(c))
	return lexRange, nil
}

func lexRangeDashBackslash(s *State) (LexType, error) {
	c := s.pattern[s.patternPos]
	s.move(1)
	if c == 'u' {
		return makeLexUnicode(lexRangeDash, 'u', 6, true), nil
	}
	if c == 'U' {
		return makeLexUnicode(lexRangeDash, 'U', 10, true), nil
	}
	s.patternBuff = append(s.patternBuff, backslashEscape(c))
	return lexRangeDash, nil
}
