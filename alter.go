package main

import (
	rand "crypto/rand"
	"math"
	"math/big"
)

type alterGenerator struct {
	parts     [][]rune
	indexList []int
	entropy   *float64
	absPos    uint
}

func (g *alterGenerator) Generate(s *State) error {
	parts := g.parts
	indexList := g.indexList
	ibig, err := rand.Int(rand.Reader, big.NewInt(int64(len(parts))))
	if err != nil {
		panic(err)
	}
	s.patternEntropy += math.Log2(float64(len(parts)))
	i := int(ibig.Int64())
	groupId := s.lastGroupId
	s2 := NewState(&SharedState{}, parts[i])
	s2.absPos = g.absPos + uint(indexList[i])
	s2.lastGroupId = groupId
	s2.groupsOutput = s.groupsOutput
	output, err := subGenerate(s2, parts[i])
	if err != nil {
		return err
	}
	s.output = append(s.output, output...)
	s.patternEntropy += s2.patternEntropy
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
