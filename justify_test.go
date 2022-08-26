package main

import (
	"testing"

	"github.com/ilius/is/v2"
)

func TestRJustGenerator(t *testing.T) {
	is := is.New(t)
	pattern := `$rjust(abc,8)`
	s := NewState(NewSharedState(), []rune(pattern))
	g, err := newRjustGenerator(s, []rune("abc,8"))
	is.NotErr(err)
	is.NotNil(g)
	{
		err := g.Generate(s)
		is.NotErr(err)
		is.Equal(`     abc`, string(s.output))
	}
	{
		entropy, err := g.Entropy(s)
		is.NotErr(err)
		is.Equal(0, entropy)
	}
}

func TestLJustGenerator(t *testing.T) {
	is := is.New(t)
	pattern := `$ljust(abc,8)`
	s := NewState(NewSharedState(), []rune(pattern))
	g, err := newLjustGenerator(s, []rune("abc,8"))
	is.NotErr(err)
	is.NotNil(g)
	{
		err := g.Generate(s)
		is.NotErr(err)
		is.Equal(`abc     `, string(s.output))
	}
	{
		entropy, err := g.Entropy(s)
		is.NotErr(err)
		is.Equal(0, entropy)
	}
}

func TestRJustGeneratorErr1(t *testing.T) {
	is := is.New(t)
	s := NewState(NewSharedState(), []rune(`$rjust()`))
	g, err := newRjustGenerator(s, []rune(""))
	is.ErrMsg(err, "argument error near index 0: rjust: too few characters as arguments")
	is.Nil(g)
}

func TestLJustGeneratorErr1(t *testing.T) {
	is := is.New(t)
	s := NewState(NewSharedState(), []rune(`$ljust()`))
	g, err := newLjustGenerator(s, []rune(""))
	is.ErrMsg(err, "argument error near index 0: ljust: too few characters as arguments")
	is.Nil(g)
}

func TestCenterGeneratorErr1(t *testing.T) {
	is := is.New(t)
	s := NewState(NewSharedState(), []rune(`$center()`))
	g, err := newCenterGenerator(s, []rune(""))
	is.ErrMsg(err, "argument error near index 0: center: too few characters as arguments")
	is.Nil(g)
}
