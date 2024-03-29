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
			panic(err)
		}
		i := int(ibig.Int64())
		s.output = append(s.output, chars[i])
	}
	entropy, err := g.Entropy(s)
	if err != nil {
		return err
	}
	s.patternEntropy += entropy
	return nil
}

func (g *charClassGenerator) Entropy(s *State) (float64, error) {
	if g.entropy != nil {
		return *g.entropy, nil
	}
	entropy := 0.0
	for _, chars := range g.charClasses {
		if len(chars) == 0 {
			continue
		}
		entropy += math.Log2(float64(len(chars)))
	}
	g.entropy = &entropy
	return entropy, nil
}
