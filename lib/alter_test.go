package passgen_test

import (
	"testing"

	"github.com/ilius/is/v2"
	passgen "github.com/ilius/repassgen/lib"
)

func TestAlterGenerator(t *testing.T) {
	is := is.New(t)
	pattern := []rune("(a|bc)")
	g := passgen.NewAlterGenerator(
		[][]rune{[]rune("a"), []rune("bc")},
		[]uint64{0, 1},
	)
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
		isFloatBetween(is, entropy, 1, 1)
	}
}
