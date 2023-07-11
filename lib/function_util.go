package passgen

func baseFunctionCallGenerator(
	s *State,
	argState *State,
	funcName string,
	funcObj func(s *State, in []rune) ([]rune, error),
) error {
	g := NewRootGenerator()
	err := g.Generate(argState)
	if err != nil {
		return err
	}
	result, err := funcObj(s, argState.output)
	if err != nil {
		return err
	}
	s.addOutputNonRepeatable(result)
	return nil
}
