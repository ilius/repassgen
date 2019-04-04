package main

// GenerateInput is struct given to Generate
type GenerateInput struct {
	Pattern            string
	CalcPatternEntropy bool
}

// GenerateOutput is struct returned by Generate
type GenerateOutput struct {
	Password       []rune
	PatternEntropy float64
}

// Generate generates random password based on given pattern
// see README.md for examples of pattern
func Generate(in GenerateInput) (*GenerateOutput, error) {
	s := NewState(in.Pattern, in.CalcPatternEntropy)
	lex := LexRoot
	var err error
	for {
		lex, err = lex(s)
		if err != nil {
			return nil, err
		}
		if lex == nil {
			break
		}
	}
	return &GenerateOutput{
		Password:       s.output,
		PatternEntropy: s.patternEntropy,
	}, nil
}
