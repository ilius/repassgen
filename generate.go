package main

// GenerateInput is struct given to Generate
type GenerateInput struct {
	Pattern []rune
}

// GenerateOutput is struct returned by Generate
type GenerateOutput struct {
	Password       []rune
	PatternEntropy float64
}

// Generate generates random password based on given pattern
// see README.md for examples of pattern
func Generate(in GenerateInput) (*GenerateOutput, error) {
	ss := &SharedState{}
	escapedPattern, err := unescapeUnicode(in.Pattern)
	if err != nil {
		return nil, err
	}
	s := NewState(ss, escapedPattern)
	g := NewRootGenerator()
	{
		err := g.Generate(s)
		if err != nil {
			return nil, err
		}
	}
	return &GenerateOutput{
		Password:       s.output,
		PatternEntropy: s.patternEntropy,
	}, nil
}

func subGenerate(s *State, pattern []rune) ([]rune, error) {
	childGen := NewRootGenerator()
	ss := s.SharedState
	var output []rune
	{
		s := NewState(ss, pattern)
		err := childGen.Generate(s)
		if err != nil {
			return nil, err
		}
		output = s.output
	}
	s.lastGen = nil
	return output, nil
}
