package main

import (
	rand "crypto/rand"
	"encoding/hex"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/tyler-smith/go-bip39"
)

var bip39WordCount = len(bip39.GetWordList())

func bip39encode(s *State, in []rune) ([]rune, error) {
	data, err := hex.DecodeString(string(in))
	if err != nil {
		s.errorMarkLen = len(in)
		return nil, s.errorValue("invalid hex number %#v", string(in))
	}
	return []rune(bip39.Encode(data)), nil
}

type bip39WordGenerator struct {
	wordCount int64
}

func (g *bip39WordGenerator) Generate(s *State) error {
	count := g.wordCount
	words := make([]string, count)
	for ai := int64(0); ai < count; ai++ {
		ibig, err := rand.Int(rand.Reader, big.NewInt(int64(bip39WordCount)))
		if err != nil {
			panic(err)
		}
		index := int(ibig.Int64())
		word, ok := bip39.GetWordIndex(index)
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

func (g *bip39WordGenerator) Entropy(s *State) (float64, error) {
	return float64(g.wordCount) * math.Log2(float64(bip39WordCount)), nil
}

func newBIP39WordGenerator(s *State, arg string) (*bip39WordGenerator, error) {
	if arg == "" {
		return &bip39WordGenerator{
			wordCount: 1,
		}, nil
	}
	argInt64, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		s.errorOffset += int64(len(arg))
		s.errorMarkLen = len(arg)
		return nil, s.errorValue("invalid number '%v'", arg)
	}
	return &bip39WordGenerator{
		wordCount: argInt64,
	}, nil
}
