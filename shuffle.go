package main

import (
	"fmt"
)

func shuffle(in []rune) []rune {
	r := NewRandSource()
	r.Shuffle(len(in), func(i int, j int) {
		in[i], in[j] = in[j], in[i]
	})
	return in
}

type shuffleGenerator struct {
	argPattern        string
	argPatternEntropy *float64
}

func (g *shuffleGenerator) Generate(s *State) error {
	argState := NewState(s.SharedState, g.argPattern)
	err := baseFunctionCallGenerator(
		s,
		argState,
		"shuffle",
		shuffle,
	)
	if err != nil {
		return err
	}
	g.argPatternEntropy = &s.patternEntropy
	return nil
}

func (g *shuffleGenerator) Entropy() (float64, error) {
	// FIXME: how to calculate entropy?
	if g.argPatternEntropy != nil {
		return *g.argPatternEntropy, nil
	}
	return 0, fmt.Errorf("entropy is not calculated")
}

func newShuffleGenerator(arg string) (*shuffleGenerator, error) {
	return &shuffleGenerator{
		argPattern: arg,
	}, nil
}
