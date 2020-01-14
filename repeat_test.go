package main

import (
	"fmt"
	"testing"

	"github.com/ilius/is"
)

func TestRepeatGenerator(t *testing.T) {
	is := is.New(t)
	str := "a"
	count := 4
	pattern := fmt.Sprintf("%s{%d}", str, count)
	// if len(str) > 1, pattern would be wrong, but it won't effect this test
	g := &repeatGenerator{
		child: &staticStringGenerator{str: []rune(str)},
		count: count,
	}
	s := NewState(&SharedState{}, pattern)
	{
		err := g.Generate(s)
		is.NotErr(err)
	}
	{
		entropy, err := g.Entropy()
		is.NotErr(err)
		is.Equal(0, entropy)
	}
}
