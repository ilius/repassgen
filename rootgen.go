package main

import (
	"fmt"
)

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
	lex := LexRoot
	var err error
	for {
		lex, err = lex(s)
		if err != nil {
			return err
		}
		if lex == nil {
			break
		}
	}
	g.entropy = &s.patternEntropy
	return nil
}

// Entropy returns the entropy after .Generate() is called
func (g *RootGenerator) Entropy() (float64, error) {
	if g.entropy != nil {
		return *g.entropy, nil
	}
	return 0, fmt.Errorf("entropy is not calculated")
}
