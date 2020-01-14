package main

import (
	"strings"
	"testing"

	"github.com/ilius/is"
)

func TestGenerate(t *testing.T) {
	test := func(tc *genCase) {
		is := is.New(t).AddMsg("pattern=%#v", tc.Pattern)
		out, err := Generate(GenerateInput{Pattern: tc.Pattern})
		if tc.Error != "" {
			is.ErrMsg(err, tc.Error)
			is.Nil(out)
			return
		}
		is.NotErr(err)
		{
			length := len(out.Password)
			minLen := tc.PassLen[0]
			maxLen := tc.PassLen[1]
			is.AddMsg(
				"length=%v is not in range [%v, %v]",
				length,
				minLen,
				maxLen,
			).True(minLen <= length && length <= maxLen)
		}

		isFloatBetween(is, out.PatternEntropy, tc.Entropy[0], tc.Entropy[1])

		if tc.Password != nil {
			is.Equal(string(out.Password), *tc.Password)
		}
		if tc.WordCount != 0 {
			actual := len(strings.Split(string(out.Password), " "))
			is.Equal(actual, tc.WordCount)
		}
	}
	test(&genCase{
		Pattern: "",
		PassLen: [2]int{0, 0},
		Entropy: [2]float64{0, 0},
	})
	test(&genCase{
		Pattern: "[a-z]{8}",
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{37.6, 37.7},
	})
	test(&genCase{
		Pattern: "[a-z]{8,10}",
		PassLen: [2]int{8, 10},
		Entropy: [2]float64{37.6, 47.01},
	})
	test(&genCase{
		Pattern: "[]",
		PassLen: [2]int{0, 0},
		Entropy: [2]float64{0, 0},
	})
	test(&genCase{
		Pattern: "[a-z]{8}[1-9]{3}",
		PassLen: [2]int{11, 11},
		Entropy: [2]float64{47.1, 47.2},
	})
	test(&genCase{
		Pattern: "([a-z]{8}[1-9]{3})",
		PassLen: [2]int{11, 11},
		Entropy: [2]float64{47.1, 47.2},
	})
	test(&genCase{
		Pattern: "([a-z]{5}[1-9]{2}){2}",
		PassLen: [2]int{14, 14},
		Entropy: [2]float64{59.6, 59.7},
	})
	test(&genCase{
		Pattern: "([a-z]{5}[1-9]{2}-){2}",
		PassLen: [2]int{16, 16},
		Entropy: [2]float64{59.6, 59.7},
	})
	// base64 length: ((bytes + 2) / 3) * 4
	test(&genCase{
		Pattern: "$base64([:alnum:]{10})",
		PassLen: [2]int{16, 16},
		Entropy: [2]float64{59.5, 59.6},
	})
	test(&genCase{
		Pattern: "$base64([:alnum:]{9})",
		PassLen: [2]int{12, 12},
		Entropy: [2]float64{53.5, 53.6},
	})
	test(&genCase{
		Pattern: "$base64([:alnum:]{5})",
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{29.7, 29.8},
	})
	test(&genCase{
		Pattern: "$base64url([:alnum:]{5})",
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{29.7, 29.8},
	})
	test(&genCase{
		Pattern: "$base32([:alnum:]{5})",
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{29.7, 29.8},
	})
	test(&genCase{
		Pattern: "$BASE32([:alnum:]{5})",
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{29.7, 29.8},
	})
	test(&genCase{
		Pattern: "$base32std([:alnum:]{5})",
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{29.7, 29.8},
	})
	test(&genCase{
		Pattern: "$hex([:alnum:]{8})",
		PassLen: [2]int{16, 16},
		Entropy: [2]float64{47.6, 47.7},
	})
	test(&genCase{
		Pattern: "$HEX([:alnum:]{8})",
		PassLen: [2]int{16, 16},
		Entropy: [2]float64{47.6, 47.7},
	})
	test(&genCase{
		Pattern:  `$escape(")`,
		PassLen:  [2]int{2, 2},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`\"`),
	})
	test(&genCase{
		Pattern:  `a[\t][\r][\n][\v][\f]b`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr("a\t\r\n\v\fb"),
	})
	// each bip39 is at least 3 chars and max 8 chars
	test(&genCase{
		Pattern:   "$bip39word(10)",
		WordCount: 10,
		PassLen:   [2]int{39, 89}, // 10*4-1, 10*9+-1
		Entropy:   [2]float64{110, 110},
	})
	test(&genCase{
		Pattern:   "$bip39word()",
		WordCount: 1,
		PassLen:   [2]int{3, 8},
		Entropy:   [2]float64{11, 11},
	})
	test(&genCase{
		Pattern: "$bip39word(x)",
		Error:   "invalid number 'x'",
	})
	// 1 bip39 word   => 11 bits entropy
	// 8 bip39 words  => 11 bytes entropy
	test(&genCase{
		Pattern:   "$bip39encode([:alpha:]{11})",
		WordCount: 8,
		PassLen:   [2]int{42, 98}, // 11*4-1, 11*9-1 // FIXME: why 42?
		Entropy:   [2]float64{62.7, 62.8},
	})
	test(&genCase{
		Pattern: "$foo(123)",
		Error:   "value error near index 4: invalid function 'foo'",
	})
	test(&genCase{
		Pattern: "test($foo(123))",
		Error:   "value error near index 9: invalid function 'foo'",
	})
	test(&genCase{
		Pattern: "test $foo",
		Error:   "syntax error near index 8: '(' not closed",
	})
	test(&genCase{
		Pattern: "test($foo)",
		Error:   "syntax error near index 8: '(' not closed",
	})
	test(&genCase{
		Pattern: "[a-z]{1-3}",
		Error:   "value error near index 7: non-numeric character '-' inside {...}",
	})
	test(&genCase{
		Pattern: "test([a-z]{1-3})",
		Error:   "value error near index 12: non-numeric character '-' inside {...}",
	})
	test(&genCase{
		Pattern: "test([a-z]{1a})",
		Error:   "value error near index 12: non-numeric character 'a' inside {...}",
	})
	test(&genCase{
		Pattern: "test([a-z]{})",
		Error:   "syntax error near index 11: missing number inside {}",
	})
}
