package passgen_test

import (
	"testing"

	"github.com/ilius/is/v2"
	passgen "github.com/ilius/repassgen/lib"
)

func TestEncoderFunctionCallGenerator(t *testing.T) {
	is := is.New(t)
	{
		s := newTestState("$foo()")
		gen := passgen.NewEncoderFunctionCallGenerator("foo", []rune("$foo()"))
		err := gen.Generate(s)
		is.ErrMsg(err, `value error near index 0: invalid function 'foo'`)
		entropy, entropyErr := gen.Entropy(s)
		is.Equal(entropy, 0)
		is.ErrMsg(entropyErr, `unknown error near index 0: entropy is not calculated`)
	}
	{
		s := newTestState("$hex([a-z]{2})")
		gen := passgen.NewEncoderFunctionCallGenerator("hex", []rune("$hex([a-z]{2})"))
		err := gen.Generate(s)
		is.NotErr(err)
		entropy, entropyErr := gen.Entropy(s)
		is.NotErr(entropyErr)
		isFloatBetween(is, entropy, 9.4, 9.5)
	}
}
