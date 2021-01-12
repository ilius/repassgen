package main

import (
	"fmt"
)

func newGroupGenerator(pattern []rune) *groupGenerator {
	return &groupGenerator{
		pattern:  pattern,
		childGen: NewRootGenerator(),
	}
}

type groupGenerator struct {
	childGen *RootGenerator
	entropy  *float64
	pattern  []rune
}

func (g *groupGenerator) Generate(s *State) error {
	ss := s.SharedState
	var output []rune
	{
		s := NewState(ss, g.pattern)
		err := g.childGen.Generate(s)
		if err != nil {
			return err
		}
		output = s.output
	}
	s.output = append(s.output, output...)
	s.lastGen = nil
	g.entropy = &s.patternEntropy
	return nil
}

func (g *groupGenerator) Entropy(s *State) (float64, error) {
	if g.entropy != nil {
		return *g.entropy, nil
	}
	return 0, fmt.Errorf("entropy is not calculated")
}
