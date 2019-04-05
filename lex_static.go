package main

type staticStringGenerator struct {
	str []rune
}

func (g *staticStringGenerator) Generate(s *State) error {
	s.output = append(s.output, g.str...)
	return nil
}

func (g *staticStringGenerator) Level() int {
	return 0
}

func (g *staticStringGenerator) Entropy() (float64, error) {
	return 0, nil
}
