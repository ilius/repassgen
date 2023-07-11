package passgen_test

import (
	"testing"

	"github.com/ilius/is/v2"
	passgen "github.com/ilius/repassgen/lib"
)

func TestError(t *testing.T) {
	is := is.New(t)
	{
		err := passgen.NewError(passgen.ErrorSyntax, 5, "dummy error")
		is.ErrMsg(err, "syntax error near index 5: dummy error")
		is.Equal(err.SpacedError(), `     ^ syntax error: dummy error`)
	}
	{
		err := passgen.NewError(passgen.ErrorSyntax, 5, "dummy error")
		err.AppendMsg("msg 2")
		is.ErrMsg(err, "syntax error near index 5: dummy error: msg 2")
		is.Equal(err.SpacedError(), `     ^ syntax error: dummy error: msg 2`)
	}
	{
		err := passgen.NewError(passgen.ErrorSyntax, 5, "dummy error")
		err.AppendMsg("msg 2")
		err.PrependMsg("msg 3")
		is.ErrMsg(err, "syntax error near index 5: msg 3: dummy error: msg 2")
		is.Equal(err.SpacedError(), `     ^ syntax error: msg 3: dummy error: msg 2`)
	}
	{
		err := passgen.NewError(passgen.ErrorUnknown, 5, "dummy error")
		err.AppendMsg("msg 2")
		err.PrependMsg("msg 3")
		is.ErrMsg(err, "unknown error near index 5: msg 3: dummy error: msg 2")
		is.Equal(err.SpacedError(), `     ^ unknown error: msg 3: dummy error: msg 2`)
	}
	{
		err := passgen.NewError(passgen.ErrorSyntax, 5, "dummy error").WithMarkLen(1)
		is.ErrMsg(err, "syntax error near index 5: dummy error")
		is.Equal(err.SpacedError(), `     ^ syntax error: dummy error`)
	}
	{
		err := passgen.NewError(passgen.ErrorSyntax, 5, "dummy error").WithMarkLen(4)
		is.ErrMsg(err, "syntax error near index 5: dummy error")
		is.Equal(err.SpacedError(), `  ^^^^ syntax error: dummy error`)
	}
	{
		err := passgen.NewError(passgen.ErrorSyntax, 5, "dummy error").WithMarkLen(0)
		is.ErrMsg(err, "syntax error near index 5: dummy error")
		is.Equal(err.SpacedError(), `     ^ syntax error: dummy error`)
	}
}

func TestParseSpacedError(t *testing.T) {
	is := is.New(t)
	{
		spaced := `  ^^^^ syntax error: dummy error`
		err := passgen.ParseSpacedError(spaced)
		is.Equal(spaced, err.SpacedError())
		is.Equal(err.Pos(), 5)
		is.Equal(err.MarkLen(), 4)
		is.Equal(err.Type(), passgen.ErrorSyntax)
		// is.Equal(err.msgs, []string{"dummy error"})
	}
	{
		spaced := `     ^ syntax error: dummy error`
		err := passgen.ParseSpacedError(spaced)
		is.Equal(spaced, err.SpacedError())
		is.Equal(err.Pos(), 5)
		is.Equal(err.MarkLen(), 1)
		is.Equal(err.Type(), passgen.ErrorSyntax)
		// is.Equal(err.msgs, []string{"dummy error"})
	}
}
