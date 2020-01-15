package main

func baseFunctionCallGenerator(
	s *State,
	argState *State,
	funcName string,
	funcObj func(in []rune) ([]rune, error),
) error {
	g := NewRootGenerator()
	err := g.Generate(argState)
	if err != nil {
		return err
	}
	result, err := funcObj(argState.output)
	if err != nil {
		lexErr, ok := err.(*Error)
		if ok {
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
