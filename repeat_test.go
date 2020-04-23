package main

import (
	"fmt"
	"testing"

	"github.com/ilius/is/v2"
)

func TestRepeatGeneratorByGroup(t *testing.T) {
	is := is.New(t)
	count := 4
	pattern := fmt.Sprintf("[a-c]{%d}", count)
	g := &repeatGenerator{
		child: newGroupGenerator("[a-c]"),
		count: count,
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
		isFloatBetween(is, entropy, 25.3, 25.4)
	}
}

func TestRepeatGeneratorByStatic(t *testing.T) {
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
