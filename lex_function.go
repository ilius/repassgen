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
	case '(':
		if s.openBracket > 0 {
			break
		}
		s.openParenth++
	case '[':
		s.openBracket++
	case ']':
		s.openBracket--
	case ')':
		if s.openBracket > 0 {
			break
		}
		s.openParenth--
		if s.openParenth > 0 {
			break
		}
		funcName := string(s.patternBuff[:s.patternBuffStart])
		s.absPos -= uint(len(s.patternBuff)) - s.patternBuffStart + 1
		if funcName == "" {
			return nil, s.errorSyntax("missing function name")
		}
		arg := string(s.patternBuff[s.patternBuffStart:n])
		gen, err := getFuncGenerator(s, funcName, arg)
		if err != nil {
			return nil, err
		}
		err = gen.Generate(s)
		if err != nil {
			return nil, err
		}
		s.patternBuff = nil
		s.lastGen = gen
		return LexRoot, nil
	}
	s.patternBuff = append(s.patternBuff, c)
	return lexIdentFuncCall, nil
}
