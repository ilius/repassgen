package main

func lexIdentFuncCall(s *State) (LexType, error) {
	if s.end() {
		return nil, s.errorSyntax("'(' not closed")
	}
	n := uint(len(s.patternBuff))
	// "$a()"  -->  c.patternBuffStart == 1
	c := s.pattern[s.patternPos]
	s.patternPos++
	switch c {
	case '(':
		s.openParenth++
	case ')':
		s.openParenth--
		if s.openParenth > 0 {
			break
		}
		funcName := string(s.patternBuff[:s.patternBuffStart])
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
