package passgen

import (
	"testing"

	"github.com/ilius/is/v2"
)

func Test_splitArgsStr_parseError(t *testing.T) {
	is := is.New(t)
	test := func(input string, sep rune, errMsg string) {
		is := is.AddMsg("input=%#v", input)
		res, indexList, err := splitArgsStr([]rune(input), sep)
		is.ErrMsg(err, errMsg)
		is.Nil(res)
		is.Nil(indexList)
	}
	test("(", ',', `unclosed '('`)
	test(")", ',', `too many ')'`)
	test("[", ',', `unclosed '['`)
	test("{", ',', `unclosed '{'`)
}
