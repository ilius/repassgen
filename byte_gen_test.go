package main

import (
	"strings"
	"testing"

	"github.com/ilius/is/v2"
)

func TestByteGenerator(t *testing.T) {
	is := is.New(t)
	pattern := `$byte()`
	g := &byteGenerator{}
	s := NewState(NewSharedState(), []rune(pattern))
	{
		err := g.Generate(s)
		is.NotErr(err)
		password := string(s.output)
		is.Equal(2, len(password))
		is.Equal(password, strings.ToLower(password))
	}
	{
		entropy, err := g.Entropy(s)
		is.NotErr(err)
		is.Equal(8, entropy)
	}
}

func TestByteGenerator2(t *testing.T) {
	is := is.New(t)
	pattern := `$BYTE()`
	g := &byteGenerator{uppercase: true}
	s := NewState(NewSharedState(), []rune(pattern))
	{
		err := g.Generate(s)
		is.NotErr(err)
		password := string(s.output)
		is.Equal(2, len(password))
		is.Equal(password, strings.ToUpper(password))
	}
	{
		entropy, err := g.Entropy(s)
		is.NotErr(err)
		is.Equal(8, entropy)
	}
}
