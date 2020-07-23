package main

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ilius/crock32"
	"github.com/ilius/is/v2"
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
		pwStr := string(out.Password)
		is = is.AddMsg("password=%#v", pwStr)
		if tc.Validate != nil {
			is.True(tc.Validate(pwStr))
		}

		isFloatBetween(is, out.PatternEntropy, tc.Entropy[0], tc.Entropy[1])

		if tc.Password != nil {
			is.Equal(pwStr, *tc.Password)
		}
		if tc.WordCount != 0 {
			actual := len(strings.Split(pwStr, " "))
			is.Equal(actual, tc.WordCount)
		}
	}
	test(&genCase{
		Pattern: "",
		PassLen: [2]int{0, 0},
		Entropy: [2]float64{0, 0},
	})
	test(&genCase{
		Pattern: "[a-z]{a}",
		Error:   "syntax error near index 6: invalid natural number inside {...}",
	})
	test(&genCase{
		Pattern: "[a-z]{2.5}",
		Error:   "syntax error near index 7: invalid natural number inside {...}",
	})
	test(&genCase{
		Pattern: "[a-z]{2.0}",
		Error:   "syntax error near index 7: invalid natural number inside {...}",
	})
	test(&genCase{
		Pattern: "[a-z]{1-3}",
		Error:   "syntax error near index 7: repetition range syntax is '{M,N}' not '{M-N}'",
	})
	test(&genCase{
		Pattern: "test([a-z]{1-3})",
		Error:   "syntax error near index 12: repetition range syntax is '{M,N}' not '{M-N}'",
	})
	test(&genCase{
		Pattern: "test([a-z]{1,})",
		Error:   "syntax error near index 13: no number after ','",
	})
	test(&genCase{
		Pattern: "test([a-z]{,3})",
		Error:   "syntax error near index 13: no number before ','",
	})
	test(&genCase{
		Pattern: "test([a-z]{1,2,3})",
		Error:   "syntax error near index 16: multiple ',' inside {...}",
	})
	test(&genCase{
		Pattern: "[a-z]{{}",
		Error:   "syntax error near index 6: nested '{'",
	})
	test(&genCase{
		Pattern: "[a-z]{[}",
		Error:   "syntax error near index 6: '[' inside {...}",
	})
	test(&genCase{
		Pattern: "[a-z]{$}",
		Error:   "syntax error near index 6: '$' inside {...}",
	})
	test(&genCase{
		Pattern: "test([a-z]{1a})",
		Error:   "syntax error near index 12: invalid natural number inside {...}",
	})
	test(&genCase{
		Pattern: "test([a-z]{})",
		Error:   "syntax error near index 11: missing number inside {}",
	})
	test(&genCase{
		Pattern: "[a-z]{3,1}",
		Error:   "value error near index 9: invalid numbers 3 > 1 inside {...}",
	})
	test(&genCase{
		Pattern: "{3}",
		Error:   "syntax error near index 2: nothing to repeat",
	})
	test(&genCase{
		Pattern: "x{0}",
		Error:   "syntax error near index 3: invalid natural number inside {...}",
	})
	test(&genCase{
		Pattern: "[abcd]{8}",
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{16, 16},
		Validate: func(p string) bool {
			for _, c := range p {
				if c < 'a' || c > 'd' {
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: "[abcccdab]{8}",
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{16, 16},
		Validate: func(p string) bool {
			for _, c := range p {
				if c < 'a' || c > 'd' {
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: "[a-z]{8}",
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{37.6, 37.7},
		Validate: func(p string) bool {
			for _, c := range p {
				if c < 'a' || c > 'z' {
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: "[a-z]{8,10}",
		PassLen: [2]int{8, 10},
		Entropy: [2]float64{37.6, 47.01},
		Validate: func(p string) bool {
			for _, c := range p {
				if c < 'a' || c > 'z' {
					return false
				}
			}
			return true
		},
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
		Validate: func(p string) bool {
			for _, c := range p[:8] {
				if c < 'a' || c > 'z' {
					return false
				}
			}
			for _, c := range p[8:] {
				if c < '1' || c > '9' {
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: "([a-z]{8}[1-9]{3})",
		PassLen: [2]int{11, 11},
		Entropy: [2]float64{47.1, 47.2},
		Validate: func(p string) bool {
			for _, c := range p[:8] {
				if c < 'a' || c > 'z' {
					return false
				}
			}
			for _, c := range p[8:] {
				if c < '1' || c > '9' {
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: "([a-z]{5}[1-9]{2}){2}",
		PassLen: [2]int{14, 14},
		Entropy: [2]float64{59.6, 59.7},
		Validate: func(p string) bool {
			for i := 0; i < 2; i++ {
				k := 7 * i
				for _, c := range p[k : k+5] {
					if c < 'a' || c > 'z' {
						return false
					}
				}
				for _, c := range p[k+5 : k+7] {
					if c < '1' || c > '9' {
						return false
					}
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern:  `(\)){2}`,
		PassLen:  [2]int{2, 2},
		Password: strPtr("))"),
		Entropy:  [2]float64{0, 0},
	})
	test(&genCase{
		Pattern:  `(\\){2}`,
		PassLen:  [2]int{2, 2},
		Password: strPtr(`\\`),
		Entropy:  [2]float64{0, 0},
	})
	test(&genCase{
		Pattern: "([a-z]{5}[1-9]{2}-){2}",
		PassLen: [2]int{16, 16},
		Entropy: [2]float64{59.6, 59.7},
		Validate: func(p string) bool {
			for i := 0; i < 2; i++ {
				k := 8 * i
				for _, c := range p[k : k+5] {
					if c < 'a' || c > 'z' {
						return false
					}
				}
				for _, c := range p[k+5 : k+7] {
					if c < '1' || c > '9' {
						return false
					}
				}
				if p[k+7] != '-' {
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: `[^ :punct:]{128}`,
		PassLen: [2]int{128, 128},
		Entropy: [2]float64{762.1, 762.2},
	})
	test(&genCase{
		Pattern: `[^^]{10}`,
		PassLen: [2]int{10, 10},
		Entropy: [2]float64{65.5, 65.6},
	})
	test(&genCase{
		Pattern: `[!-~]{10}`,
		PassLen: [2]int{10, 10},
		Entropy: [2]float64{65.5, 65.6},
	})
	// base64 length: ((bytes + 2) / 3) * 4
	test(&genCase{
		Pattern: "$base64([:alnum:]{10})",
		PassLen: [2]int{16, 16},
		Entropy: [2]float64{59.5, 59.6},
		Validate: func(p string) bool {
			pwBytes, err := base64.StdEncoding.DecodeString(p)
			if err != nil {
				panic(err)
			}
			if len(pwBytes) != 10 {
				return false
			}
			for _, b := range bytes.ToLower(pwBytes) {
				c := rune(b)
				if c < '0' || (c > '9' && c < 'a') || c > 'z' {
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: "$base64([:alnum:]{9})",
		PassLen: [2]int{12, 12},
		Entropy: [2]float64{53.5, 53.6},
		Validate: func(p string) bool {
			pwBytes, err := base64.StdEncoding.DecodeString(p)
			if err != nil {
				panic(err)
			}
			if len(pwBytes) != 9 {
				return false
			}
			for _, b := range bytes.ToLower(pwBytes) {
				c := rune(b)
				if c < '0' || (c > '9' && c < 'a') || c > 'z' {
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: "$base64([:alnum:]{5})",
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{29.7, 29.8},
		Validate: func(p string) bool {
			pwBytes, err := base64.StdEncoding.DecodeString(p)
			if err != nil {
				panic(err)
			}
			if len(pwBytes) != 5 {
				return false
			}
			for _, b := range bytes.ToLower(pwBytes) {
				c := rune(b)
				if c < '0' || (c > '9' && c < 'a') || c > 'z' {
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: "$base64url([:alnum:]{5})",
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{29.7, 29.8},
		Validate: func(p string) bool {
			pwBytes, err := base64.URLEncoding.DecodeString(p)
			if err != nil {
				panic(err)
			}
			if len(pwBytes) != 5 {
				return false
			}
			for _, b := range bytes.ToLower(pwBytes) {
				c := rune(b)
				if c < '0' || (c > '9' && c < 'a') || c > 'z' {
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: "$base32([:alnum:]{5})",
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{29.7, 29.8},
		Validate: func(p string) bool {
			if strings.ToLower(p) != p {
				return false
			}
			pwBytes, err := crock32.Decode(p)
			if err != nil {
				panic(err)
			}
			if len(pwBytes) != 5 {
				return false
			}
			for _, b := range bytes.ToLower(pwBytes) {
				c := rune(b)
				if c < '0' || (c > '9' && c < 'a') || c > 'z' {
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: "$BASE32([:alnum:]{5})",
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{29.7, 29.8},
		Validate: func(p string) bool {
			if strings.ToUpper(p) != p {
				return false
			}
			pwBytes, err := crock32.Decode(p)
			if err != nil {
				panic(err)
			}
			if len(pwBytes) != 5 {
				return false
			}
			for _, b := range bytes.ToLower(pwBytes) {
				c := rune(b)
				if c < '0' || (c > '9' && c < 'a') || c > 'z' {
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: "$base32std([:alnum:]{5})",
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{29.7, 29.8},
		Validate: func(p string) bool {
			if strings.ToUpper(p) != p {
				return false
			}
			pwBytes, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(p)
			if err != nil {
				panic(err)
			}
			if len(pwBytes) != 5 {
				return false
			}
			for _, b := range bytes.ToLower(pwBytes) {
				c := rune(b)
				if c < '0' || (c > '9' && c < 'a') || c > 'z' {
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: "$hex([:alnum:]{8})",
		PassLen: [2]int{16, 16},
		Entropy: [2]float64{47.6, 47.7},
		Validate: func(p string) bool {
			if strings.ToLower(p) != p {
				return false
			}
			pwBytes, err := hex.DecodeString(p)
			if err != nil {
				panic(err)
			}
			if len(pwBytes) != 8 {
				return false
			}
			for _, b := range bytes.ToLower(pwBytes) {
				c := rune(b)
				if c < '0' || (c > '9' && c < 'a') || c > 'z' {
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: "$HEX([:alnum:]{8})",
		PassLen: [2]int{16, 16},
		Entropy: [2]float64{47.6, 47.7},
		Validate: func(p string) bool {
			if strings.ToUpper(p) != p {
				return false
			}
			pwBytes, err := hex.DecodeString(p)
			if err != nil {
				panic(err)
			}
			if len(pwBytes) != 8 {
				return false
			}
			for _, b := range bytes.ToLower(pwBytes) {
				c := rune(b)
				if c < '0' || (c > '9' && c < 'a') || c > 'z' {
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: "$hex([a-c)(]{4})",
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{9.28, 9.29},
		Validate: func(p string) bool {
			if strings.ToLower(p) != p {
				return false
			}
			pwBytes, err := hex.DecodeString(p)
			if err != nil {
				panic(err)
			}
			if len(pwBytes) != 4 {
				return false
			}
			for _, b := range bytes.ToLower(pwBytes) {
				switch rune(b) {
				case 'a', 'b', 'c', ')', '(':
				default:
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: "$hex(([a-e]{4}))",
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{9.28, 9.29},
		Validate: func(p string) bool {
			if strings.ToLower(p) != p {
				return false
			}
			pwBytes, err := hex.DecodeString(p)
			if err != nil {
				panic(err)
			}
			if len(pwBytes) != 4 {
				return false
			}
			for _, b := range bytes.ToLower(pwBytes) {
				c := rune(b)
				if c < 'a' || c > 'e' {
					return false
				}
			}
			return true
		},
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
	test(&genCase{
		Pattern:  `a\t\r\n\v\fb\c`,
		PassLen:  [2]int{8, 8},
		Entropy:  [2]float64{0, 0},
		Password: strPtr("a\t\r\n\v\fbc"),
	})
	// each bip39 is at least 3 chars and max 8 chars
	test(&genCase{
		Pattern:   "$bip39word(10)",
		WordCount: 10,
		PassLen:   [2]int{39, 89}, // 10*4-1, 10*9+-1
		Entropy:   [2]float64{110, 110},
		Validate: func(p string) bool {
			// TODO
			return true
		},
	})
	test(&genCase{
		Pattern:   "$bip39word()",
		WordCount: 1,
		PassLen:   [2]int{3, 8},
		Entropy:   [2]float64{11, 11},
		Validate: func(p string) bool {
			// TODO
			return true
		},
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
		PassLen:   [2]int{38, 98}, // 11*4-1, 11*9-1 // FIXME: why 38?
		Entropy:   [2]float64{62.7, 62.8},
		Validate: func(p string) bool {
			// TODO
			return true
		},
	})
	test(&genCase{
		Pattern: "$()",
		Error:   "syntax error near index 1: missing function name",
	})
	test(&genCase{
		Pattern: "$hex([a-z]",
		Error:   "syntax error near index 9: '(' not closed",
	})
	test(&genCase{
		Pattern: "$hex(([a-z]",
		Error:   "syntax error near index 10: '(' not closed",
	})
	test(&genCase{
		Pattern: "(",
		Error:   "syntax error near index 0: '(' not closed",
	})
	test(&genCase{
		Pattern: "$foo",
		Error:   "syntax error near index 3: expected a function call",
	})
	test(&genCase{
		Pattern: "$foo(123)",
		Error:   "value error near index 4: invalid function 'foo'",
	})
	test(&genCase{
		Pattern: `$foo\()`,
		Error:   "syntax error near index 4: expected a function call",
	})
	test(&genCase{
		Pattern: "test($foo(123))",
		Error:   "value error near index 9: invalid function 'foo'",
	})
	test(&genCase{
		Pattern: "test $foo",
		Error:   "syntax error near index 8: expected a function call",
	})
	test(&genCase{
		Pattern: "test($foo)",
		Error:   "syntax error near index 8: expected a function call",
	})
	test(&genCase{
		Pattern: "$shuffle([a-z]{5}[1-9]{2})",
		PassLen: [2]int{7, 7},
		Entropy: [2]float64{29.8, 29.9},
		Validate: func(p string) bool {
			alpha := 0
			num := 0
			for _, c := range p {
				if c >= 'a' && c <= 'z' {
					alpha++
				} else if c >= '0' && c <= '9' {
					num++
				}
			}
			return alpha == 5 && num == 2
		},
	})
}
