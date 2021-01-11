package main

import (
	"testing"

	"github.com/ilius/is/v2"
)

func newState(pattern string) *State {
	s := NewState(&SharedState{}, []rune(pattern))
	s.absPos = uint(len(pattern))
	return s
}

func TestEncoderFunctionCallGenerator(t *testing.T) {
	is := is.New(t)
	{
		s := newState("$foo()")
		gen := &encoderFunctionCallGenerator{
			funcName:   "foo",
			argPattern: []rune("$foo()"),
		}
		err := gen.Generate(s)
		is.ErrMsg(err, `value error near index 5: invalid function 'foo'`)
		entropy, entropyErr := gen.Entropy()
		is.Equal(entropy, 0)
		is.ErrMsg(entropyErr, `entropy is not calculated`)
	}
	{
		s := newState("$hex([:x:])")
		gen := &encoderFunctionCallGenerator{
			funcName:   "hex",
			argPattern: []rune("$hex([:x:])"),
		}
		err := gen.Generate(s)
		is.ErrMsg(err, `value error near index 19: invalid character class "x"`)
	}
	{
		s := newState("$hex([a-z]{2})")
		gen := &encoderFunctionCallGenerator{
			funcName:   "hex",
			argPattern: []rune("$hex([a-z]{2})"),
		}
		err := gen.Generate(s)
		is.NotErr(err)
		entropy, entropyErr := gen.Entropy()
		is.NotErr(entropyErr)
		isFloatBetween(is, entropy, 9.4, 9.5)
	}
}
