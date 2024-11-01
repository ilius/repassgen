package passgen

func _lexIdentFuncCallParanClose(s *State, buffer []rune) (LexType, error) {
	if s.openBracket {
		return nil, nil
	}
	s.openParenth--
	if s.openParenth > 0 {
		return nil, nil
	}
	s.move(1)
	s2 := NewState(s.SharedState.Copy(), s.input)
	s2.errorOffset -= int64(len(buffer) + 1)
	funcName := string(s.buffer)
	if funcName == "" {
		return nil, s2.errorSyntax("missing function name")
	}
	gen, err := getFuncGenerator(s2, funcName, buffer)
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
	s.buffer = nil
	s.lastGen = gen
	return LexRoot, nil
}

func _lexIdentFuncCallEndError(s *State) error {
	s.errorOffset++
	if s.openParenth > 0 {
		return s.errorSyntax(err_paranthNotClosed)
	}
	return s.errorSyntax(s_func_call_expected)
}

func _lexIdentFuncCallBackslash(s *State, buffer []rune) []rune {
	buffer = append(buffer, '\\')
	s.move(1)
	if !s.end() {
		buffer = append(buffer, s.input[s.inputPos])
	}
	return buffer
}

func lexIdentFuncCall(s *State) (LexType, error) {
	if s.end() {
		return nil, _lexIdentFuncCallEndError(s)
	}
	buffer := []rune{}
	for ; !s.end(); s.move(1) {
		c := s.input[s.inputPos]
		switch c {
		case '\\':
			buffer = _lexIdentFuncCallBackslash(s, buffer)
			continue
		case '(':
			if s.openBracket {
				break
			}
			s.openParenth++
		case '[':
			if s.openBracket {
				return nil, s.errorSyntax(err_nestedBracket)
			}
			s.openBracket = true
		case ']':
			s.openBracket = false
		case ')':
			lex, err := _lexIdentFuncCallParanClose(s, buffer)
			if err != nil {
				return nil, err
			}
			if lex != nil {
				return lex, nil
			}
		}
		buffer = append(buffer, c)
	}
	s.errorOffset++
	if s.openParenth > 0 {
		return nil, s.errorSyntax(err_paranthNotClosed)
	}
	return nil, s.errorSyntax(s_func_call_expected)
}
