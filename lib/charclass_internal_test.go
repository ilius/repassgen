package passgen

import (
	"testing"

	"github.com/ilius/is/v2"
)

func Test_charClassGenerator_Entropy(t *testing.T) {
	is := is.New(t)
	test := func(pattern string, chars string, entropy float64) {
		s := newTestState(pattern)
		gen := &charClassGenerator{
			charClasses: [][]rune{[]rune(chars)},
		}
		entropyActual, err := gen.Entropy(s)
		is.NotErr(err)
		is.Equal(entropy, entropyActual)
	}
	test("[]", "", 0)
}
