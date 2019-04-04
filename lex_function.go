package main

type functionCallGenerator struct {
	funcName   string
	argPattern string
}

func (g *functionCallGenerator) Generate(s *State) error {
	funcName := g.funcName
	funcObj, ok := functions[funcName]
	if !ok {
		return s.errorValue("invalid function '%v'", funcName)
	}
	argOut := Generate(GenerateInput{
		Pattern:            g.argPattern,
		CalcPatternEntropy: s.calcPatternEntropy,
	})
	result, err := funcObj(argOut.Password)
	if err != nil {
		lexErr, ok := err.(*LexError)
		if ok {
			lexErr.MovePos(int(s.patternBuffStart))
			lexErr.PrependMsg("function " + funcName)
			return lexErr
		}
		return s.errorUnknown("%v returned error: %v", funcName, err)
	}
	err = s.addOutputNonRepeatable(result)
	if err != nil {
		return err
	}
	if s.calcPatternEntropy {
		s.patternEntropy += argOut.PatternEntropy
	}
	return nil
}

func (g *functionCallGenerator) Level() int {
	return 0
}

func lexIdentFuncCall(s *State) (LexType, error) {
	if s.end() {
		return nil, s.errorSyntax("'(' not closed")
	}
	n := uint(len(s.patternBuff))
	// "$a()"  -->  c.patternBuffStart == 1
	c := s.pattern[s.patternPos]
	s.patternPos++
	switch c {
	case ')':
		funcName := string(s.patternBuff[:s.patternBuffStart])
		if funcName == "" {
			return nil, s.errorSyntax("missing function name")
		}
		argPattern := string(s.patternBuff[s.patternBuffStart:n])
		gen := &functionCallGenerator{
			funcName:   funcName,
			argPattern: argPattern,
		}
		err := gen.Generate(s)
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
