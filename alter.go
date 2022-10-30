package main

import (
	rand "crypto/rand"
	"math"
	"math/big"
)

type alterGenerator struct {
	entropy   *float64
	parts     [][]rune
	indexList []uint64
	absPos    uint64
}

func (g *alterGenerator) calcMinEntropy(s *State) (float64, error) {
	// TODO: optimize
	minEntropy := 0.0
	groupId := s.lastGroupId
	for i, part := range g.parts {
		s2 := NewState(NewSharedState(), part)
		s2.absPos = g.absPos + g.indexList[i]
		s2.lastGroupId = groupId
		s2.groupsOutput = s.groupsOutput
		_, err := subGenerate(s2, part)
		if err != nil {
			return 0, err
		}
		entropy := s2.patternEntropy
		if i == 0 || entropy < minEntropy {
			minEntropy = entropy
			if minEntropy == 0 {
				break
			}
		}
	}
	return minEntropy, nil
}

func (g *alterGenerator) Generate(s *State) error {
	parts := g.parts
	indexList := g.indexList
	ibig, err := rand.Int(rand.Reader, big.NewInt(int64(len(parts))))
	if err != nil {
		panic(err)
	}

	i := ibig.Int64()
	groupId := s.lastGroupId
	s2 := NewState(NewSharedState(), parts[i])
	s2.absPos = g.absPos + indexList[i]
	s2.lastGroupId = groupId
	s2.groupsOutput = s.groupsOutput
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
