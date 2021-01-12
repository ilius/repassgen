package main

import (
	"fmt"
	"testing"

	"github.com/ilius/is/v2"
)

func TestShuffleGenerator(t *testing.T) {
	is := is.New(t)
	argPattern := "[a-z]{5}[1-9]{2}"
	pattern := []rune(fmt.Sprintf("$shuffle(%s)", argPattern))
	g := &shuffleGenerator{
		argPattern: []rune(argPattern),
	}
	s := newTestState(pattern)
	{
		entropy, err := g.Entropy(s)
		is.ErrMsg(err, "unknown error near index 25: entropy is not calculated")
		is.Equal(0, entropy)
	}
	{
		err := g.Generate(s)
		is.NotErr(err)
	}
	{
		entropy, err := g.Entropy(s)
		is.NotErr(err)
		isFloatBetween(is, entropy, 29.8, 29.9)
	}
}
