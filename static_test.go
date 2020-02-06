package main

import (
	"testing"

	"github.com/ilius/is"
)

func TestStaticStringGenerator(t *testing.T) {
	is := is.New(t)
	pattern := `abc\\d`
	g := &staticStringGenerator{str: []byte(pattern)}
	s := NewState(&SharedState{}, pattern)
	{
		err := g.Generate(s)
		is.NotErr(err)
		is.Equal(pattern, s.output.String())
	}
	{
		entropy, err := g.Entropy()
		is.NotErr(err)
		is.Equal(0, entropy)
	}
}
