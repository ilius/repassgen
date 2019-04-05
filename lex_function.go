package main

import "fmt"

type functionCallGenerator struct {
	funcName   string
	argPattern string
	entropy    *float64
}

func (g *functionCallGenerator) Generate(s *State) error {
	funcName := g.funcName
	funcObj, ok := functions[funcName]
	if !ok {
		return s.errorValue("invalid function '%v'", funcName)
	}
	argOut, err := Generate(GenerateInput{
		Pattern: g.argPattern,
	})
	if err != nil {
		lexErr, ok := err.(*LexError)
		if ok {
			lexErr.MovePos(int(s.patternBuffStart + 1))
			return lexErr
		}
		return s.errorUnknown(err.Error())
	}
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
	s.patternEntropy += argOut.PatternEntropy
	g.entropy = &argOut.PatternEntropy
	return nil
}

func (g *functionCallGenerator) Level() int {
	return 0
}

func (g *functionCallGenerator) Entropy() (float64, error) {
	if g.entropy != nil {
		return *g.entropy, nil
	}
	return 0, fmt.Errorf("entropy is not calculated")
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
