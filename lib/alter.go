package passgen

import (
	rand "crypto/rand"
	"math"
	"math/big"
)

type alterGenerator struct {
	entropy   *float64
	parts     [][]rune
	indexList []uint64
	length    uint64
}

func (g *alterGenerator) calcMinEntropy(s *State) (float64, error) {
	// TODO: optimize
	minEntropy := 0.0
	for partI, part := range g.parts {
		s2 := NewState(s.SharedState.Copy(), part)
		s2.errorOffset += int64(g.indexList[partI] - g.length)
		s2.patternEntropy = 0
		_, err := subGenerate(s2, part)
		if err != nil {
			return 0, err
		}
		entropy := s2.patternEntropy
		if entropy == 0 {
			return 0, nil
		}
		if partI == 0 || entropy < minEntropy {
			minEntropy = entropy
		}
	}
	return minEntropy, nil
}

func (g *alterGenerator) Generate(s *State) error {
	parts := g.parts
	ibig, err := rand.Int(rand.Reader, big.NewInt(int64(len(parts))))
	if err != nil {
		panic(err)
	}

	i := ibig.Int64()
	groupId := s.lastGroupId
	s2 := NewState(s.SharedState.Copy(), parts[i])
	s2.errorOffset += int64(g.indexList[i] - g.length)
	output, err := subGenerate(s2, parts[i])
	if err != nil {
		return err
	}
	s.output = append(s.output, output...)

	minEntropy, err := g.calcMinEntropy(s)
	if err != nil {
		return err
	}

	s.patternEntropy += math.Log2(float64(len(parts)))
	s.patternEntropy += minEntropy
	g.entropy = &s.patternEntropy

	s.lastGroupId = s2.lastGroupId
	s.groupsOutput[groupId] = output
	return nil
}

func (g *alterGenerator) Entropy(s *State) (float64, error) {
	if g.entropy != nil {
		return *g.entropy, nil
	}
	return 0, s.errorUnknown("entropy is not calculated")
}
