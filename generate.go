package main

// GenerateInput is struct given to Generate
type GenerateInput struct {
	Pattern string
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
	s := NewState(ss, in.Pattern)
	g := NewRootGenerator()
	err := g.Generate(s)
	if err != nil {
		return nil, err
	}
	return &GenerateOutput{
		Password:       s.output,
		PatternEntropy: s.patternEntropy,
	}, nil
}
