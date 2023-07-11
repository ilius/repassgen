package passgen

func shuffle(s *State, in []rune) ([]rune, error) {
	r := NewRandSource()
	r.Shuffle(len(in), func(i int, j int) {
		in[i], in[j] = in[j], in[i]
	})
	return in, nil
}

type shuffleGenerator struct {
	argPatternEntropy *float64
	argPattern        []rune
}

func (g *shuffleGenerator) Generate(s *State) error {
	argState := NewState(s.SharedState, g.argPattern)
	err := baseFunctionCallGenerator(
		s,
		argState,
		"shuffle",
		shuffle,
	)
	if err != nil {
		return err
	}
	g.argPatternEntropy = &s.patternEntropy
	return nil
}

func (g *shuffleGenerator) Entropy(s *State) (float64, error) {
	// FIXME: how to calculate entropy?
	if g.argPatternEntropy != nil {
		return *g.argPatternEntropy, nil
	}
	return 0, s.errorUnknown("entropy is not calculated")
}

func newShuffleGenerator(arg []rune) (*shuffleGenerator, error) {
	return &shuffleGenerator{
		argPattern: arg,
	}, nil
}
