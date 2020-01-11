package main

import "fmt"

type groupGenerator struct {
	pattern string
	entropy *float64
}

func (g *groupGenerator) Generate(s *State) error {
	err := generate(
		s.SharedState,
		GenerateInput{
			Pattern: g.pattern,
		},
	)
	if err != nil {
		lexErr, ok := err.(*LexError)
		if ok {
			lexErr.MovePos(int(s.patternBuffStart))
			return lexErr
		}
		return s.errorUnknown(err.Error())
	}
	s.lastGen = nil
	if err != nil {
		return err
	}
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
