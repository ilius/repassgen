package main

import (
	rand "crypto/rand"
	"math"
	"math/big"
)

type charsetGroupGenerator struct {
	charsets [][]rune
	entropy  *float64
}

func (g *charsetGroupGenerator) Generate(s *State) error {
	for _, charset := range g.charsets {
		if len(charset) == 0 {
			continue
		}
		ibig, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}
		i := int(ibig.Int64())
		s.output = append(s.output, charset[i])
	}
	entropy, err := g.Entropy()
	if err != nil {
		return err
	}
	s.patternEntropy += entropy
	return nil
}

func (g *charsetGroupGenerator) Entropy() (float64, error) {
	if g.entropy != nil {
		return *g.entropy, nil
	}
	entropy := 0.0
	for _, charset := range g.charsets {
		if len(charset) == 0 {
			continue
		}
		entropy += math.Log2(float64(len(charset)))
	}
	g.entropy = &entropy
	return entropy, nil
}
