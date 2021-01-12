package main

// NewRootGenerator creates a new RootGenerator
func NewRootGenerator() *RootGenerator {
	return &RootGenerator{}
}

// RootGenerator is the root Generator implementation
type RootGenerator struct {
	entropy *float64
}

// Generate generates a password
func (g *RootGenerator) Generate(s *State) error {
	err := g.lexLoop(s)
	if err != nil {
		return err
	}
	g.entropy = &s.patternEntropy
	return nil
}

func (g *RootGenerator) lexLoop(s *State) error {
	lex := LexRoot
	var err error
	for {
		lex, err = lex(s)
		if err != nil {
			return err
		}
		if lex == nil {
			return nil
		}
	}
}

// Entropy returns the entropy after .Generate() is called
func (g *RootGenerator) Entropy(s *State) (float64, error) {
	if g.entropy != nil {
		return *g.entropy, nil
	}
	return 0, s.errorUnknown("entropy is not calculated")
}
