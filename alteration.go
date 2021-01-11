package main

import (
	rand "crypto/rand"
	"fmt"
	"math"
	"math/big"
)

func newAlterationGenerator(pattern []rune, alterPos []uint) *alterationGenerator {
	fmt.Printf("\nnewAlterationGenerator: pattern=%#v, alterPos=%#v\n", string(pattern), alterPos)
	return &alterationGenerator{
		childGen: NewRootGenerator(),
		pattern:  pattern,
		alterPos: alterPos,
	}
}

type alterationGenerator struct {
	childGen *RootGenerator
	entropy  *float64
	pattern  []rune
	alterPos []uint
}

func (g *alterationGenerator) Generate(s *State) error {
	ss := s.SharedState
	var output []rune
	count := len(g.alterPos) + 1
	randBig, err := rand.Int(rand.Reader, big.NewInt(int64(count)))
	if err != nil {
		panic(err)
	}
	alterIndex := int(randBig.Int64())
	var pattern []rune
	if alterIndex == 0 {
		pattern = g.pattern[:g.alterPos[0]-1]
	} else if alterIndex == len(g.alterPos) {
		pattern = g.pattern[g.alterPos[alterIndex-1]:]
	} else {
		pattern = g.pattern[g.alterPos[alterIndex-1] : g.alterPos[alterIndex]-1]
	}
	fmt.Printf("alterIndex = %v, pattern=%#v\n", alterIndex, string(pattern))
	{
		s := NewState(ss, pattern)
		err := g.childGen.Generate(s)
		if err != nil {
			return err
		}
		output = s.output
	}
	s.output = append(s.output, output...)
	s.lastGen = nil
	s.patternEntropy += math.Log2(float64(count))
	g.entropy = &s.patternEntropy
	return nil
}

func (g *alterationGenerator) Entropy(s *State) (float64, error) {
	if g.entropy != nil {
		return *g.entropy, nil
	}
	return 0, s.errorUnknown("entropy is not calculated")
}
