package main

import (
	"testing"

	"github.com/ilius/is/v2"
)

func Test_splitArgsStr_error(t *testing.T) {
	is := is.New(t)
	sep := ','

	testError := func(pattern string, errMsg string) {
		res, indexList, err := splitArgsStr([]rune(pattern), sep)
		is.Nil(res)
		is.Nil(indexList)
		is.ErrMsg(err, errMsg)
	}

	testError(`test (`, "unclosed '('")
	testError(`test ())`, "too many ')'")
	testError(`test [`, "unclosed '['")
	testError(`test {`, "unclosed '{'")
	testError(`test {{`, "nested '{'")
}
