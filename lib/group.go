package passgen

func newGroupGenerator(pattern []rune) *groupGenerator {
	return &groupGenerator{
		pattern: pattern,
	}
}

type groupGenerator struct {
	entropy *float64
	pattern []rune
}

func (g *groupGenerator) Generate(s *State) error {
	err := subGenerate(s, g.pattern)
	if err != nil {
		return err
	}
	g.entropy = &s.patternEntropy
	return nil
}

func (g *groupGenerator) Entropy(s *State) (float64, error) {
	if g.entropy != nil {
		return *g.entropy, nil
	}
	return 0, s.errorUnknown("entropy is not calculated")
}
