package main

func lexGroup(s *State) (LexType, error) {
	if s.end() {
		s.errorOffset++
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
		s2 := NewState(&SharedState{}, s.pattern)
		s2.output = s.output
		s2.absPos = s.absPos - uint(len(s.patternBuff)) - 1
		s2.patternEntropy = s.patternEntropy
		childPattern := s.patternBuff[s.patternBuffStart:len(s.patternBuff)]
		gen := newGroupGenerator(childPattern)
		err := gen.Generate(s2)
		if err != nil {
			return nil, err
		}
		s.output = s2.output
		s.patternEntropy = s2.patternEntropy
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
	s.patternBuff = append(s.patternBuff, '\\', c)
	return lexGroup, nil
}
