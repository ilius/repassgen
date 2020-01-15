package main

import (
	"testing"

	"github.com/ilius/is"
)

func TestError(t *testing.T) {
	is := is.New(t)
	{
		err := NewError(LexErrorSyntax, 5, "dummy error")
		is.ErrMsg(err, "syntax error near index 5: dummy error")
	}
	{
		err := NewError(LexErrorSyntax, 5, "dummy error")
		err.AppendMsg("msg 2")
		is.ErrMsg(err, "syntax error near index 5: dummy error: msg 2")
	}
	{
		err := NewError(LexErrorSyntax, 5, "dummy error")
		err.AppendMsg("msg 2")
		err.PrependMsg("msg 3")
		is.ErrMsg(err, "syntax error near index 5: msg 3: dummy error: msg 2")
	}
}