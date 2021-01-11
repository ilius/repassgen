package main

import "fmt"

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
	if len(s.alterPos) > 0 {
		fmt.Printf("s.alterPos = %#v\n", s.alterPos)
		// fmt.Printf("s = %#v\n", *s)
		fmt.Printf("s.pattern = %#v\n", string(s.pattern))
		fmt.Printf("s.patternBuff = %#v\n", string(s.patternBuff))
		fmt.Printf("s.patternBuffStart = %d\n", s.patternBuffStart)
		fmt.Printf("s.patternPos = %d\n", s.patternPos)
		fmt.Printf("s.outputGroupPos = %v\n", s.outputGroupPos)
		if len(s.outputGroupPos) > 0 {
			lastPos := s.outputGroupPos[len(s.outputGroupPos)-1]
			s.output = s.output[:lastPos]
		} else {
			s.output = nil
		}
		var childPattern []rune
		if len(s.patternBuff) > 0 {
			childPattern = s.patternBuff
		} else {
			childPattern = s.pattern
		}
		gen := newAlterationGenerator(childPattern, s.alterPos)
		err := gen.Generate(s)
		if err != nil {
			return err
		}
		g.entropy = &s.patternEntropy
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
