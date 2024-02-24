package passgen_test

import (
	"strings"
	"testing"

	"github.com/ilius/is/v2"
	passgen "github.com/ilius/repassgen/lib"
)

func TestByteGenerator(t *testing.T) {
	is := is.New(t)
	pattern := `$byte()`
	g := passgen.NewByteGenerator(false)
	s := passgen.NewState(passgen.NewSharedState(), []rune(pattern))
	{
		err := g.Generate(s)
		is.NotErr(err)
		password := s.Output()
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
	g := passgen.NewByteGenerator(true)
	s := passgen.NewState(passgen.NewSharedState(), []rune(pattern))
	{
		err := g.Generate(s)
		is.NotErr(err)
		password := s.Output()
		is.Equal(2, len(password))
		is.Equal(password, strings.ToUpper(password))
	}
	{
		entropy, err := g.Entropy(s)
		is.NotErr(err)
		is.Equal(8, entropy)
	}
}
