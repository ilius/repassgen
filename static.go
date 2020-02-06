package main

type staticStringGenerator struct {
	str []byte
}

func (g *staticStringGenerator) Generate(s *State) error {
	s.output.Write(g.str)
	return nil
}

func (g *staticStringGenerator) Entropy() (float64, error) {
	return 0, nil
}
