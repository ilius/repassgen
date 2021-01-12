package main

import (
	rand "crypto/rand"
	"math/big"
)

type onceOrNoneGenerator struct {
	entropy    float64
	argPattern []rune
}

func randBool() bool {
	randBig, err := rand.Int(rand.Reader, big.NewInt(10))
	if err != nil {
		panic(err)
	}
	return randBig.Int64()%2 == 1
}

func (g *onceOrNoneGenerator) Generate(s *State) error {
	childGen := NewRootGenerator()
	ss := s.SharedState
	var output []rune
	{
		s := NewState(ss, g.argPattern)
		err := childGen.Generate(s)
		if err != nil {
			return err
		}
		output = s.output
	}
	if randBool() {
		s.output = append(s.output, output...)
	}
	s.lastGen = nil
	s.patternEntropy += 1.0
	g.entropy = s.patternEntropy
	return nil
}

func (g *onceOrNoneGenerator) Entropy(s *State) (float64, error) {
	return g.entropy, nil
}

func newOnceOrNoneGenerator(argPattern []rune) (*onceOrNoneGenerator, error) {
	return &onceOrNoneGenerator{
		argPattern: argPattern,
	}, nil
}
