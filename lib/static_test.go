package passgen_test

import (
	"testing"

	"github.com/ilius/is/v2"
	passgen "github.com/ilius/repassgen/lib"
)

func TestStaticStringGenerator(t *testing.T) {
	is := is.New(t)
	pattern := `abc\\d`
	g := passgen.NewStaticStringGenerator([]rune(pattern))
	s := passgen.NewState(passgen.NewSharedState(), []rune(pattern))
	{
		err := g.Generate(s)
		is.NotErr(err)
		is.Equal(pattern, s.Output())
	}
	{
		entropy, err := g.Entropy(s)
		is.NotErr(err)
		is.Equal(0, entropy)
	}
}
