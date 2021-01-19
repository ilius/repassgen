package main

import (
	"testing"

	"github.com/ilius/is/v2"
)

func TestEncoderFunctionCallGenerator(t *testing.T) {
	is := is.New(t)
	{
		s := newTestState("$foo()")
		gen := &encoderFunctionCallGenerator{
			funcName:   "foo",
			argPattern: []rune("$foo()"),
		}
		err := gen.Generate(s)
		is.ErrMsg(err, `value error near index 0: invalid function 'foo'`)
		entropy, entropyErr := gen.Entropy(s)
		is.Equal(entropy, 0)
		is.ErrMsg(entropyErr, `unknown error near index 0: entropy is not calculated`)
	}
	{
		s := newTestState("$hex([:x:])")
		gen := &encoderFunctionCallGenerator{
			funcName:   "hex",
			argPattern: []rune("$hex([:x:])"),
		}
		err := gen.Generate(s)
		is.ErrMsg(err, `value error near index 7: invalid character class "x"`)
	}
	{
		s := newTestState("$hex([a-z]{2})")
		gen := &encoderFunctionCallGenerator{
			funcName:   "hex",
			argPattern: []rune("$hex([a-z]{2})"),
		}
		err := gen.Generate(s)
		is.NotErr(err)
		entropy, entropyErr := gen.Entropy(s)
		is.NotErr(entropyErr)
		isFloatBetween(is, entropy, 9.4, 9.5)
	}
}
