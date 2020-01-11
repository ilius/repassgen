package main

import "fmt"

func newGroupGenerator(pattern string) *groupGenerator {
	return &groupGenerator{
		pattern:  pattern,
		childGen: NewRootGenerator(),
	}
}

type groupGenerator struct {
	pattern  string
	entropy  *float64
	childGen *RootGenerator
}

func (g *groupGenerator) Generate(s *State) error {
	ss := s.SharedState
	{
		s := NewState(ss, g.pattern)
		err := g.childGen.Generate(s)
		if err != nil {
			lexErr, ok := err.(*LexError)
			if ok {
				lexErr.MovePos(1)
				return lexErr
			}
			return s.errorUnknown(err.Error())
		}
	}
	s.lastGen = nil
	g.entropy = &s.patternEntropy
	return nil
}

func (g *groupGenerator) Level() int {
	return 0
}

func (g *groupGenerator) Entropy() (float64, error) {
	if g.entropy != nil {
		return *g.entropy, nil
	}
	return 0, fmt.Errorf("entropy is not calculated")
}
