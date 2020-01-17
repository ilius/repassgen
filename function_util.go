package main

func baseFunctionCallGenerator(
	s *State,
	argState *State,
	funcName string,
	funcObj func(in []rune) []rune,
) error {
	g := NewRootGenerator()
	err := g.Generate(argState)
	if err != nil {
		return err
	}
	result := funcObj(argState.output)
	err = s.addOutputNonRepeatable(result)
	if err != nil {
		return err
	}
	return nil
}
