package passgen

import "bytes"

func baseFunctionCallGenerator(
	s *State,
	argState *State,
	funcName string,
	funcObj func(s *State, in []rune) ([]rune, error),
) error {
	g := NewRootGenerator()
	outBuf := bytes.NewBuffer(nil)
	argState.output = outBuf
	err := g.Generate(argState)
	if err != nil {
		return err
	}
	result, err := funcObj(s, []rune(outBuf.String()))
	if err != nil {
		return err
	}
	s.addOutputNonRepeatable(result)
	return nil
}
