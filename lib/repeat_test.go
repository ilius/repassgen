package passgen_test

import (
	"fmt"
	"testing"

	"github.com/ilius/is/v2"
	passgen "github.com/ilius/repassgen/lib"
)

func TestRepeatGeneratorByGroup(t *testing.T) {
	is := is.New(t)
	count := int64(4)
	pattern := []rune(fmt.Sprintf("[a-c]{%d}", count))
	g := passgen.NewRepeatGenerator(
		passgen.NewGroupGenerator([]rune("[a-c]")),
		count,
	)
	s := newTestState(pattern)
	{
		entropy, err := g.Entropy(s)
		is.ErrMsg(err, "unknown error near index 0: entropy is not calculated")
		is.Equal(0, entropy)
	}
	{
		err := g.Generate(s)
		is.NotErr(err)
	}
	{
		entropy, err := g.Entropy(s)
		is.NotErr(err)
		isFloatBetween(is, entropy, 25.3, 25.4)
	}
}

func TestRepeatGeneratorByStatic(t *testing.T) {
	is := is.New(t)
	str := "a"
	count := int64(4)
	pattern := []rune(fmt.Sprintf("%s{%d}", str, count))
	// if len(str) > 1, pattern would be wrong, but it won't effect this test
	g := passgen.NewRepeatGenerator(
		passgen.NewStaticStringGenerator([]rune(str)),
		count,
	)
	s := passgen.NewState(passgen.NewSharedState(), pattern)
	{
		err := g.Generate(s)
		is.NotErr(err)
	}
	{
		entropy, err := g.Entropy(s)
		is.NotErr(err)
		is.Equal(0, entropy)
	}
}

func TestRepeatGeneratorGroupFuncError(t *testing.T) {
	is := is.New(t)
	g := passgen.NewRepeatGenerator(
		passgen.NewGroupGenerator([]rune("$foo()")),
		2,
	)
	s := passgen.NewState(passgen.NewSharedState(), []rune(`($foo()){2}`))
	err := g.Generate(s)
	tErr := err.(*passgen.Error)
	is.Equal(
		`^^^^^ value error: invalid function 'foo'`,
		tErr.SpacedError(),
	)
}
