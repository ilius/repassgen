package main

func lexIdentFuncCall(s *State) (LexType, error) {
	if s.end() {
		s.errorOffset++
		if s.openParenth > 0 {
			return nil, s.errorSyntax("'(' not closed")
		}
		return nil, s.errorSyntax("expected a function call")
	}
	buff := []rune{}
	for ; !s.end(); s.move(1) {
		c := s.pattern[s.patternPos]
		switch c {
		case '\\':
			buff = append(buff, '\\')
			s.move(1)
			if s.end() {
				continue
			}
			buff = append(buff, s.pattern[s.patternPos])
			continue
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
			s.move(1)
			s2 := NewState(NewSharedState(), s.pattern)
			s2.output = s.output
			s2.absPos = s.absPos - (uint(len(buff)) + 1)
			s2.patternEntropy = s.patternEntropy
			s2.lastGroupId = s.lastGroupId
			s2.groupsOutput = s.groupsOutput
			funcName := string(s.patternBuff)
			if funcName == "" {
				return nil, s2.errorSyntax("missing function name")
			}
			gen, err := getFuncGenerator(s2, funcName, buff)
			if err != nil {
				return nil, err
			}
			err = gen.Generate(s2)
			if err != nil {
				return nil, err
			}
			s.output = s2.output
			s.patternEntropy = s2.patternEntropy
			s.lastGroupId = s2.lastGroupId
			s.patternBuff = nil
			s.lastGen = gen
			return LexRoot, nil
		}
		buff = append(buff, c)
	}
	s.errorOffset++
	if s.openParenth > 0 {
		return nil, s.errorSyntax("'(' not closed")
	}
	return nil, s.errorSyntax("expected a function call")
}
