package main

import (
	rand "crypto/rand"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/ilius/bip39-coder/bip39"
)

func bip39encode(in []rune) []rune {
	return []rune(bip39.Encode([]byte(string(in))))
}

type bip39WordGenerator struct {
	wordCount int
}

func (g *bip39WordGenerator) Generate(s *State) error {
	count := g.wordCount
	words := make([]string, count)
	for ai := 0; ai < count; ai++ {
		ibig, err := rand.Int(rand.Reader, big.NewInt(int64(bip39.WordCount())))
		if err != nil {
			panic(err)
		}
		index := int(ibig.Int64())
		word, ok := bip39.GetWord(index)
		if !ok {
			return fmt.Errorf("internal error, index=%v > 2048", index)
		}
		words[ai] = word
	}
	result := []rune(strings.Join(words, " "))

	s.addOutputNonRepeatable(result)
	entropy, err := g.Entropy()
	if err != nil {
		return err
	}
	s.patternEntropy += entropy
	return nil
}

func (g *bip39WordGenerator) Entropy() (float64, error) {
	return float64(g.wordCount) * math.Log2(float64(bip39.WordCount())), nil
}

func newBIP39WordGenerator(arg string) (*bip39WordGenerator, error) {
	if arg == "" {
		return &bip39WordGenerator{
			wordCount: 1,
		}, nil
	}
	argInt64, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number '%v'", arg)
	}
	return &bip39WordGenerator{
		wordCount: int(argInt64),
	}, nil
}
