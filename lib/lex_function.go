package passgen

func _lexIdentFuncCallParanClose(s *State, buff []rune) (LexType, error) {
	if s.openBracket {
		return nil, nil
	}
	s.openParenth--
	if s.openParenth > 0 {
		return nil, nil
	}
	s.move(1)
	s2 := NewState(s.SharedState.Copy(), s.input)
	s2.errorOffset -= int64(len(buff) + 1)
	funcName := string(s.buff)
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
	s.output = append(s.output, s2.output...)
	s.patternEntropy = s2.patternEntropy
	s.lastGroupId = s2.lastGroupId
	s.buff = nil
	s.lastGen = gen
	return LexRoot, nil
}

func _lexIdentFuncCallEndError(s *State) error {
	s.errorOffset++
	if s.openParenth > 0 {
		return s.errorSyntax("'(' not closed")
	}
	return s.errorSyntax("expected a function call")
}

func _lexIdentFuncCallBackslash(s *State, buff []rune) []rune {
	buff = append(buff, '\\')
	s.move(1)
	if !s.end() {
		buff = append(buff, s.input[s.inputPos])
	}
	return buff
}

func lexIdentFuncCall(s *State) (LexType, error) {
	if s.end() {
		return nil, _lexIdentFuncCallEndError(s)
	}
	buff := []rune{}
	for ; !s.end(); s.move(1) {
		c := s.input[s.inputPos]
		switch c {
		case '\\':
			buff = _lexIdentFuncCallBackslash(s, buff)
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
			lex, err := _lexIdentFuncCallParanClose(s, buff)
			if err != nil {
				return nil, err
			}
			if lex != nil {
				return lex, nil
			}
		}
		buff = append(buff, c)
	}
	s.errorOffset++
	if s.openParenth > 0 {
		return nil, s.errorSyntax("'(' not closed")
	}
	return nil, s.errorSyntax("expected a function call")
}
