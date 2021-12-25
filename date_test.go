package main

import (
	"testing"

	"github.com/ilius/is/v2"
)

func TestDateGenerator(t *testing.T) {
	is := is.New(t)
	pattern := `$date(2022,2023)`
	g := &dateGenerator{
		sep:     "-",
		startJd: 2459581,
		endJd:   2459946,
	}
	s := NewState(NewSharedState(), []rune(pattern))
	{
		err := g.Generate(s)
		is.NotErr(err)
		is.Equal(10, len(string(s.output)))
	}
	{
		entropy, err := g.Entropy(s)
		is.NotErr(err)
		isFloatBetween(is, entropy, 8.51, 8.52)
	}
}
