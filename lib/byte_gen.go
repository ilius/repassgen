package passgen

import (
	"crypto/rand"
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
		panic(err) // not sure how to trigger this in test
	}
	byteStr := hex.EncodeToString([]byte{uint8(randBig.Uint64())})
	if g.uppercase {
		byteStr = strings.ToUpper(byteStr)
	}
	s.addOutputNonRepeatable([]rune(byteStr))
	s.patternEntropy += 8
	return nil
}

func (g *byteGenerator) Entropy(_ *State) (float64, error) {
	return 8, nil
}

func newByteGenerator(s *State, argsStr []rune, uppercase bool) (*byteGenerator, error) {
	if len(argsStr) > 0 {
		s.errorOffset += 1
		return nil, s.errorValue("function does not accept any arguments")
	}
	return &byteGenerator{uppercase: uppercase}, nil
}
