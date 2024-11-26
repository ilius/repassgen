package passgen

import (
	"crypto/rand"
	"math/big"
)

type onceOrNoneGenerator struct {
	pattern []rune
	entropy float64
}

func randBool() bool {
	randBig, err := rand.Int(rand.Reader, big.NewInt(10))
	if err != nil {
		panic(err)
	}
	return randBig.Int64()%2 == 1
}

func (g *onceOrNoneGenerator) Generate(s *State) error {
	if randBool() {
		output, err := subGenerate(s, g.pattern)
		if err != nil {
			return err
		}
		s.output = append(s.output, output...)
	}
	s.patternEntropy += 1.0
	g.entropy = s.patternEntropy
	return nil
}

func (g *onceOrNoneGenerator) Entropy(_ *State) (float64, error) {
	return g.entropy, nil
}

func newOnceOrNoneGenerator(pattern []rune) (*onceOrNoneGenerator, error) {
	return &onceOrNoneGenerator{
		pattern: pattern,
	}, nil
}
