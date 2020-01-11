package main

func lexGroup(s *State) (LexType, error) {
	if s.end() {
		return nil, s.errorSyntax("'(' not closed")
	}
	c := s.pattern[s.patternPos]
	s.move(1)
	switch c {
	case '\\':
		return lexGroupBackslash, nil
	case '(':
		s.openParenth++
	case ')':
		s.openParenth--
		if s.openParenth > 0 {
			break
		}
		s.absPos -= uint(len(s.patternBuff)) + 1
		childPattern := string(s.patternBuff[s.patternBuffStart:len(s.patternBuff)])
		gen := newGroupGenerator(childPattern)
		err := gen.Generate(s)
		if err != nil {
			return nil, err
		}
		s.lastGen = gen
		s.patternBuff = nil
		return LexRoot, nil
	}
	s.patternBuff = append(s.patternBuff, c)
	return lexGroup, nil
}

func lexGroupBackslash(s *State) (LexType, error) {
	c := s.pattern[s.patternPos]
	s.move(1)
	if c == ')' {
		s.patternBuff = append(s.patternBuff, c)
	} else {
		s.patternBuff = append(s.patternBuff, '\\', c)
	}
	return lexGroup, nil
}
