package main

import (
	"testing"

	"github.com/ilius/is"
)

func TestGroupGenerator(t *testing.T) {
	is := is.New(t)
	pattern := "([a-z]{5}[1-9]{2}){2}"
	g := newGroupGenerator(pattern)
	s := NewState(&SharedState{}, pattern)
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
