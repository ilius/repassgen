package main

import (
	"testing"

	"github.com/ilius/is/v2"
)

func TestGroupGenerator(t *testing.T) {
	is := is.New(t)
	pattern := "([a-z]{5}[1-9]{2}){2}"
	g := newGroupGenerator(pattern)
	s := NewState(&SharedState{}, pattern)
	{
		entropy, err := g.Entropy()
		is.ErrMsg(err, "entropy is not calculated")
		is.Equal(0.0, entropy)
	}
	{
		err := g.Generate(s)
		is.NotErr(err)
	}
	{
		entropy, err := g.Entropy()
		is.NotErr(err)
		isFloatBetween(is, entropy, 59.6, 59.7)
	}
}
