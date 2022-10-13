package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ilius/is/v2"
)

func TestRepeatGeneratorByGroup(t *testing.T) {
	is := is.New(t)
	count := int64(4)
	pattern := []rune(fmt.Sprintf("[a-c]{%d}", count))
	g := &repeatGenerator{
		child: newGroupGenerator([]rune("[a-c]")),
		count: count,
	}
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
	g := &repeatGenerator{
		child: &staticStringGenerator{str: []rune(str)},
		count: count,
	}
	s := NewState(NewSharedState(), pattern)
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

func Test_parseRepeatCount(t *testing.T) {
	is := is.New(t)
	pre := "test [a-z]"
	{
		count := "10a0,abc"
		pattern := fmt.Sprintf("%s{%s}", pre, count)
		s := newTestState(pattern)
		s.move(uint64(len(pattern)))
		c, err := parseRepeatCount(s, []rune(count))
		is.Equal(0, c)
		is.Equal(
			strings.Repeat(" ", len(pre)+1)+`^^^^ syntax error: invalid natural number '10a0'`,
			err.(*Error).SpacedError(),
		)
	}
	{
		minStr := "1000"
		count := minStr + ",abc"
		pattern := fmt.Sprintf("%s{%s}", pre, count)
		s := newTestState(pattern)
		s.move(uint64(len(pattern)))
		c, err := parseRepeatCount(s, []rune(count))
		is.Equal(0, c)
		is.Equal(
			err.(*Error).SpacedError(),
			strings.Repeat(" ", len(pre+minStr)+1)+`^^^ syntax error: invalid natural number 'abc'`,
		)
	}
}
