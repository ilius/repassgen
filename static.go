package main

type staticStringGenerator struct {
	str []rune
}

func (g *staticStringGenerator) Generate(s *State) error {
	s.output = append(s.output, g.str...)
	return nil
}

func (g *staticStringGenerator) CharProb() map[rune]float64 {
	m := make(map[rune]float64, len(g.str))
	for _, r := range g.str {
		m[r] = 1.0
	}
	return m
}

func (g *staticStringGenerator) Level() int {
	return 0
}

func (g *staticStringGenerator) Entropy(s *State) (float64, error) {
	return 0, nil
}
