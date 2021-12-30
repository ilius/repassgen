package main

import (
	rand "crypto/rand"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/ilius/bip39-coder/bip39"
	"github.com/spf13/cast"
)

func bip39encode(s *State, in []rune) ([]rune, error) {
	return []rune(bip39.Encode([]byte(string(in)))), nil
}

type bip39WordGenerator struct {
	wordCount int64
}

func (g *bip39WordGenerator) Generate(s *State) error {
	count := g.wordCount
	words := make([]string, count)
	for ai := int64(0); ai < count; ai++ {
		ibig, err := rand.Int(rand.Reader, big.NewInt(int64(bip39.WordCount())))
		if err != nil {
			panic(err)
		}
		index, err := cast.ToIntE(ibig.Int64())
		if err != nil {
			return err
		}
		word, ok := bip39.GetWord(index)
		if !ok {
			return s.errorUnknown("internal error, index=%v > 2048", index)
		}
		words[ai] = word
	}
	result := []rune(strings.Join(words, " "))

	s.addOutputNonRepeatable(result)
	entropy, err := g.Entropy(s)
	if err != nil {
		return err
	}
	s.patternEntropy += entropy
	return nil
}

func (g *bip39WordGenerator) CharProb() map[rune]float64 {
	// TODO: charProb
	return nil
}

func (g *bip39WordGenerator) Level() int {
	return 0
}

func (g *bip39WordGenerator) Entropy(s *State) (float64, error) {
	return float64(g.wordCount) * math.Log2(float64(bip39.WordCount())), nil
}

func newBIP39WordGenerator(s *State, arg string) (*bip39WordGenerator, error) {
	if arg == "" {
		return &bip39WordGenerator{
			wordCount: 1,
		}, nil
	}
	argInt64, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		s.errorOffset += int64(len(arg) + 1)
		return nil, s.errorValue("invalid number '%v'", arg)
	}
	return &bip39WordGenerator{
		wordCount: argInt64,
	}, nil
}
