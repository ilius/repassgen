package passgen

type staticStringGenerator struct {
	str []rune
}

func (g *staticStringGenerator) Generate(s *State) error {
	s.output = append(s.output, g.str...)
	return nil
}

func (g *staticStringGenerator) Entropy(_ *State) (float64, error) {
	return 0, nil
}
