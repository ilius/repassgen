package main

import (
	"testing"

	"github.com/ilius/is/v2"
)

func TestError(t *testing.T) {
	is := is.New(t)
	{
		err := NewError(ErrorSyntax, 5, "dummy error")
		is.ErrMsg(err, "syntax error near index 5: dummy error")
		is.Equal(err.SpacedError(), `     ^ syntax error: dummy error`)
	}
	{
		err := NewError(ErrorSyntax, 5, "dummy error")
		err.AppendMsg("msg 2")
		is.ErrMsg(err, "syntax error near index 5: dummy error: msg 2")
		is.Equal(err.SpacedError(), `     ^ syntax error: dummy error: msg 2`)
	}
	{
		err := NewError(ErrorSyntax, 5, "dummy error")
		err.AppendMsg("msg 2")
		err.PrependMsg("msg 3")
		is.ErrMsg(err, "syntax error near index 5: msg 3: dummy error: msg 2")
		is.Equal(err.SpacedError(), `     ^ syntax error: msg 3: dummy error: msg 2`)
	}
	{
		err := NewError(ErrorUnknown, 5, "dummy error")
		err.AppendMsg("msg 2")
		err.PrependMsg("msg 3")
		is.ErrMsg(err, "unknown error near index 5: msg 3: dummy error: msg 2")
		is.Equal(err.SpacedError(), `     ^ unknown error: msg 3: dummy error: msg 2`)
	}
}
