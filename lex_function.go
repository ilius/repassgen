package main

func lexIdentFuncCall(s *State) (LexType, error) {
	if s.end() {
		if s.openParenth > 0 {
			return nil, s.errorSyntax("'(' not closed")
		}
		return nil, s.errorSyntax("expected a function call")
	}
	n := uint(len(s.patternBuff))
	// "$a()"  -->  c.patternBuffStart == 1
	c := s.pattern[s.patternPos]
	s.move(1)
	switch c {
	case '\\':
		return lexIdentFuncCallBackslash, nil
	case '(':
		if s.openBracket {
			break
		}
		s.openParenth++
	case '[':
		if s.openBracket {
			return nil, s.errorSyntax("nested '['")
		}
		s.openBracket = true
	case ']':
		s.openBracket = false
	case ')':
		if s.openBracket {
			break
		}
		s.openParenth--
		if s.openParenth > 0 {
			break
		}
		s2 := NewState(&SharedState{}, s.pattern)
		s2.output = s.output
		s2.absPos = s.absPos - (uint(len(s.patternBuff)) - s.patternBuffStart + 1)
		s2.patternEntropy = s.patternEntropy
		funcName := string(s.patternBuff[:s.patternBuffStart])
		if funcName == "" {
			return nil, s2.errorSyntax("missing function name")
		}
		arg := s.patternBuff[s.patternBuffStart:n]
		gen, err := getFuncGenerator(s2, funcName, arg)
		if err != nil {
			return nil, err
		}
		err = gen.Generate(s2)
		if err != nil {
			return nil, err
		}
		s.output = s2.output
		s.patternEntropy = s2.patternEntropy
		s.patternBuff = nil
		s.lastGen = gen
		return LexRoot, nil
	}
	s.patternBuff = append(s.patternBuff, c)
	return lexIdentFuncCall, nil
}

func lexIdentFuncCallBackslash(s *State) (LexType, error) {
	c := s.pattern[s.patternPos]
	s.move(1)
	s.patternBuff = append(s.patternBuff, '\\', c)
	return lexIdentFuncCall, nil
}
