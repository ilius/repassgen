package main

func baseFunctionCallGenerator(
	s *State,
	funcName string,
	funcObj func(in []rune) ([]rune, error),
	argPattern string,
) (*GenerateOutput, error) {
	ss := &SharedState{}
	err := generate(
		ss,
		GenerateInput{
			Pattern: argPattern,
		},
	)
	if err != nil {
		lexErr, ok := err.(*LexError)
		if ok {
			lexErr.MovePos(int(s.patternBuffStart + 1))
			return nil, lexErr
		}
		return nil, s.errorUnknown(err.Error())
	}
	result, err := funcObj(ss.output)
	if err != nil {
		lexErr, ok := err.(*LexError)
		if ok {
			lexErr.MovePos(int(s.patternBuffStart))
			lexErr.PrependMsg("function " + funcName)
			return nil, lexErr
		}
		return nil, s.errorUnknown("%v returned error: %v", funcName, err)
	}
	err = s.addOutputNonRepeatable(result)
	if err != nil {
		return nil, err
	}
	return &GenerateOutput{
		Password:       ss.output,
		PatternEntropy: ss.patternEntropy,
	}, nil
}
