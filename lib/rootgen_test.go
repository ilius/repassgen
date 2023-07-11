package passgen_test

import (
	"testing"

	"github.com/ilius/is/v2"
	passgen "github.com/ilius/repassgen/lib"
)

func TestRootGenerator(t *testing.T) {
	is := is.New(t)
	pattern := []rune("[a-z]{5}[1-9]{2}")
	g := passgen.NewRootGenerator()
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
