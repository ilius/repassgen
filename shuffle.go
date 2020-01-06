package main

import "fmt"

func shuffle(in []rune) ([]rune, error) {
	r := NewRandSource()
	r.Shuffle(len(in), func(i int, j int) {
		in[i], in[j] = in[j], in[i]
	})
	return in, nil
}

type shuffleGenerator struct {
	argPattern        string
	argPatternEntropy *float64
}

func (g *shuffleGenerator) Generate(s *State) error {
	argOut, err := baseFunctionCallGenerator(s, "shuffle", shuffle, g.argPattern)
	if err != nil {
		return err
	}
	g.argPatternEntropy = &argOut.PatternEntropy
	return nil
}

func (g *shuffleGenerator) Level() int {
	return 0
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
