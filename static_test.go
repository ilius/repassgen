package main

import (
	"testing"

	"github.com/ilius/is/v2"
)

func TestStaticStringGenerator(t *testing.T) {
	is := is.New(t)
	pattern := `abc\\d`
	g := &staticStringGenerator{str: []rune(pattern)}
	s := NewState(NewSharedState(), []rune(pattern))
	{
		err := g.Generate(s)
		is.NotErr(err)
		is.Equal(pattern, string(s.output))
	}
	{
		entropy, err := g.Entropy(s)
		is.NotErr(err)
		is.Equal(0, entropy)
	}
}
