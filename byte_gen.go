package main

import (
	rand "crypto/rand"
	"encoding/hex"
	"math/big"
	"strings"
)

type byteGenerator struct {
	uppercase bool
}

func (g *byteGenerator) Generate(s *State) error {
	randBig, err := rand.Int(rand.Reader, big.NewInt(0xff))
	if err != nil {
		panic(err)
	}
	byteStr := hex.EncodeToString([]byte{uint8(randBig.Uint64())})
	if g.uppercase {
		byteStr = strings.ToUpper(byteStr)
	}
	s.addOutputNonRepeatable([]rune(byteStr))
	s.patternEntropy += g.entropy()
	return nil
}

func (g *byteGenerator) entropy() float64 {
	return 8
}

func (g *byteGenerator) Entropy(s *State) (float64, error) {
	return g.entropy(), nil
}

func newByteGenerator(s *State, argsStr []rune, uppercase bool) (*byteGenerator, error) {
	if len(argsStr) > 0 {
		s.errorOffset += 1
		return nil, s.errorValue("function does not accept any arguments")
	}
	return &byteGenerator{uppercase: uppercase}, nil
}
