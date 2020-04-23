package main

import (
	"fmt"
	"testing"

	"github.com/ilius/is/v2"
)

func TestShuffleGenerator(t *testing.T) {
	is := is.New(t)
	argPattern := "[a-z]{5}[1-9]{2}"
	pattern := fmt.Sprintf("$shuffle(%s)", argPattern)
	g := &shuffleGenerator{
		argPattern: argPattern,
	}
	s := NewState(&SharedState{}, pattern)
	{
		entropy, err := g.Entropy()
		is.ErrMsg(err, "entropy is not calculated")
		is.Equal(0, entropy)
	}
	{
		err := g.Generate(s)
		is.NotErr(err)
	}
	{
		entropy, err := g.Entropy()
		is.NotErr(err)
		isFloatBetween(is, entropy, 29.8, 29.9)
	}
}
