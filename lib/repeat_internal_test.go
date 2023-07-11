package passgen

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ilius/is/v2"
)

func newTestState(patternArg interface{}) *State {
	var pattern []rune
	switch patternTyped := patternArg.(type) {
	case string:
		pattern = []rune(patternTyped)
	case []rune:
		pattern = patternTyped
	default:
		panic("invalid patternArg")
	}
	s := NewState(NewSharedState(), pattern)
	return s
}

func Test_parseRepeatCount(t *testing.T) {
	is := is.New(t)
	pre := "test [a-z]"

	testError := func(count string) *Error {
		pattern := fmt.Sprintf("%s{%s}", pre, count)
		s := newTestState(pattern)
		s.move(uint64(len(pattern)))
		c, err := parseRepeatCount(s, []rune(count))
		is.Equal(0, c)
		return err.(*Error)
	}

	{
		err := testError("10a0")
		is.Equal(
			strings.Repeat(" ", len(pre)+1)+`^^^^ syntax error: invalid natural number '10a0'`,
			err.SpacedError(),
		)
	}
	{
		err := testError("10a0,abc")
		is.Equal(
			strings.Repeat(" ", len(pre)+1)+`^^^^ syntax error: invalid natural number '10a0'`,
			err.SpacedError(),
		)
	}
	{
		minStr := "1000"
		err := testError(minStr + ",abc")
		is.Equal(
			strings.Repeat(" ", len(pre+minStr)+1)+`^^^ syntax error: invalid natural number 'abc'`,
			err.SpacedError(),
		)
	}
	{
		err := testError("10,20,30")
		is.Equal(
			`                   ^ syntax error: multiple ',' inside {...}`,
			err.SpacedError(),
		)
	}
}
