package passgen

import (
	rand "crypto/rand"
	"math"
	"math/big"
)

type charClassGenerator struct {
	entropy     *float64
	charClasses [][]rune
}

func (g *charClassGenerator) Generate(s *State) error {
	for _, chars := range g.charClasses {
		if len(chars) == 0 {
			continue
		}
		ibig, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			panic(err) // not sure how to trigger this in test
		}
		i := int(ibig.Int64())
		s.output = append(s.output, chars[i])
	}
	entropy := g.getEntropy()
	s.patternEntropy += entropy
	return nil
}

func (g *charClassGenerator) getEntropy() float64 {
	if g.entropy != nil {
		return *g.entropy
	}
	entropy := 0.0
	for _, chars := range g.charClasses {
		if len(chars) == 0 {
			continue
		}
		entropy += math.Log2(float64(len(chars)))
	}
	g.entropy = &entropy
	return entropy
}

func (g *charClassGenerator) Entropy(_ *State) (float64, error) {
	return g.getEntropy(), nil
}
