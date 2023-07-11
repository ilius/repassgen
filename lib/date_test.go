package passgen_test

import (
	"testing"

	"github.com/ilius/is/v2"
	passgen "github.com/ilius/repassgen/lib"
)

func TestDateGenerator(t *testing.T) {
	is := is.New(t)
	pattern := `$date(2022,2023)`
	g := passgen.NewDateGenerator("-", 2459581, 2459946)
	s := passgen.NewState(passgen.NewSharedState(), []rune(pattern))
	{
		err := g.Generate(s)
		is.NotErr(err)
		is.Equal(10, len(string(s.Output())))
	}
	{
		entropy, err := g.Entropy(s)
		is.NotErr(err)
		isFloatBetween(is, entropy, 8.51, 8.52)
	}
}
