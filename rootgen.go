package main

import (
	"fmt"
)

func NewRootGenerator() *RootGenerator {
	return &RootGenerator{}
}

type RootGenerator struct {
	entropy *float64
}

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

func (g *RootGenerator) Entropy() (float64, error) {
	if g.entropy != nil {
		return *g.entropy, nil
	}
	return 0, fmt.Errorf("entropy is not calculated")
}
