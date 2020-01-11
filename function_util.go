package main

func baseFunctionCallGenerator(
	s *State,
	funcName string,
	funcObj func(in []rune) ([]rune, error),
	argPattern string,
) error {
	err := generate(
		s.SharedState,
		GenerateInput{
			Pattern: argPattern,
		},
	)
	if err != nil {
		lexErr, ok := err.(*LexError)
		if ok {
			lexErr.MovePos(int(s.patternBuffStart + 1))
			return lexErr
		}
		return s.errorUnknown(err.Error())
	}
	result, err := funcObj(s.output)
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
	return nil
}
