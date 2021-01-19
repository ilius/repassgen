package main

import (
	"testing"

	"github.com/ilius/is/v2"
)

func TestRootGenerator(t *testing.T) {
	is := is.New(t)
	pattern := []rune("[a-z]{5}[1-9]{2}")
	g := NewRootGenerator()
	s := newTestState(pattern)
	{
		entropy, err := g.Entropy(s)
		is.ErrMsg(err, "unknown error near index 0: entropy is not calculated")
		is.Equal(0.0, entropy)
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
