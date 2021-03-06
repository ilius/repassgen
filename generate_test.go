package main

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"os"

	"github.com/ilius/crock32"
	"github.com/ilius/is/v2"
)

var verbose = os.Getenv("TEST_VERBOSE") == "1"

func TestGenerate(t *testing.T) {
	test := func(tc *genCase) {
		is := is.New(t).AddMsg("pattern=%#v", tc.Pattern)
		is = is.Lax()
		out, _, err := Generate(GenerateInput{Pattern: []rune(tc.Pattern)})
		if !is.NotErr(err) {
			tErr, okErr := err.(*Error)
			if okErr {
				t.Logf("Error: `" + tErr.SpacedError() + "`")
			}
			return
		}
		pwStr := string(out.Password)
		is = is.AddMsg("password=%#v", pwStr)
		if tc.Password != nil {
			if !is.Equal(pwStr, *tc.Password) {
				return
			}
		}
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
		if tc.Validate != nil {
			is.AddMsg("validation failed").True(tc.Validate(pwStr))
		}

		isFloatBetween(is, out.PatternEntropy, tc.Entropy[0], tc.Entropy[1])

		if tc.WordCount != 0 {
			actual := len(strings.Split(pwStr, " "))
			is.Equal(actual, tc.WordCount)
		}
	}
	testErr := func(tc *genErrCase) {
		is := is.New(t).AddMsg("pattern=%#v", tc.Pattern)
		is = is.Lax()
		out, s, err := Generate(GenerateInput{Pattern: []rune(tc.Pattern)})
		tErr, okErr := err.(*Error)
		switch expErr := tc.Error.(type) {
		case string:
			if okErr {
				expErrTyped := ParseSpacedError(expErr)
				if expErrTyped == nil {
					t.Errorf("bad spaced error %#v", expErr)
					is.Equal(tErr.SpacedError(), expErr)
					return
				}
				is.Equal(tErr.typ, expErrTyped.typ)
				is.Equal(tErr.Message(), expErrTyped.Message())
				is.AddMsg(
					"msg=%#v", tErr.Message(),
				).Equal(tErr.pos, expErrTyped.pos)
			} else {
				is.Equal(err.Error(), expErr)
			}
		case *Error:
			is.Equal(tErr.typ, expErr.typ)
			is.Equal(tErr.Message(), expErr.Message())
			is.AddMsg(
				"msg=%#v", tErr.Message(),
			).Equal(tErr.pos, expErr.pos)
		case error:
			is.ErrMsg(err, expErr.Error())
		}
		if okErr && verbose {
			t.Log(string(s.pattern))
			t.Log(tErr.SpacedError())
		}
		is.Nil(out)
		if verbose {
			t.Log("------------------------------------")
		}
	}
	test(&genCase{
		Pattern: ``,
		PassLen: [2]int{0, 0},
		Entropy: [2]float64{0, 0},
	})
	testErr(&genErrCase{
		Pattern: `[a`,
		Error:   `  ^ syntax error: '[' not closed`,
	})
	testErr(&genErrCase{
		Pattern: `[[]]`,
		Error:   ` ^ syntax error: nested '['`,
	})
	testErr(&genErrCase{
		Pattern: `[:x]`,
		Error:   `   ^ syntax error: ':' not closed`,
	})
	testErr(&genErrCase{
		Pattern: `[:x`,
		Error:   `   ^ syntax error: ':' not closed`,
	})
	testErr(&genErrCase{
		Pattern: `[a-`,
		Error:   `   ^ syntax error: '[' not closed`,
	})
	testErr(&genErrCase{
		Pattern: `[a-]`,
		Error:   `   ^ syntax error: no character after '-'`,
	})
	testErr(&genErrCase{
		Pattern: `[-a]`,
		Error:   ` ^ syntax error: no character before '-'`,
	})
	testErr(&genErrCase{
		Pattern: `[a-z]{a}`,
		Error:   `      ^ syntax error: invalid natural number inside {...}`,
	})
	testErr(&genErrCase{
		Pattern: `[a-z]{2.5}`,
		Error:   `       ^ syntax error: invalid natural number inside {...}`,
	})
	testErr(&genErrCase{
		Pattern: `[a-z]{2.0}`,
		Error:   `       ^ syntax error: invalid natural number inside {...}`,
	})
	testErr(&genErrCase{
		Pattern: `[a-z]{1-3}`,
		Error:   `       ^ syntax error: repetition range syntax is '{M,N}' not '{M-N}'`,
	})
	testErr(&genErrCase{
		Pattern: `test([a-z]{1-3})`,
		Error:   `            ^ syntax error: repetition range syntax is '{M,N}' not '{M-N}'`,
	})
	testErr(&genErrCase{
		Pattern: `test([a-z]{1,})`,
		Error:   `             ^ syntax error: no number after ','`,
	})
	testErr(&genErrCase{
		Pattern: `test([a-z]{,3333})`,
		Error:   `           ^ syntax error: no number before ','`,
	})
	testErr(&genErrCase{
		Pattern: `test([a-z]{1,2,3})`,
		Error:   `              ^ syntax error: multiple ',' inside {...}`,
	})
	testErr(&genErrCase{
		Pattern: `[a-z]{{}`,
		Error:   `      ^ syntax error: nested '{'`,
	})
	testErr(&genErrCase{
		Pattern: `[a-z]{[}`,
		Error:   `      ^ syntax error: '[' inside {...}`,
	})
	testErr(&genErrCase{
		Pattern: `[a-z]{$}`,
		Error:   `      ^ syntax error: '$' inside {...}`,
	})
	testErr(&genErrCase{
		Pattern: `test([a-z]{1a})`,
		Error:   `            ^ syntax error: invalid natural number inside {...}`,
	})
	testErr(&genErrCase{
		Pattern: `test([a-z]{})`,
		Error:   `           ^ syntax error: missing number inside {}`,
	})
	testErr(&genErrCase{
		Pattern: `[a-z]{3,1}`,
		Error:   `         ^ value error: invalid numbers 3 > 1 inside {...}`,
	})
	testErr(&genErrCase{
		Pattern: `{3}`,
		Error:   `^ syntax error: nothing to repeat`,
	})
	testErr(&genErrCase{
		Pattern: `x{0}`,
		Error:   `   ^ syntax error: invalid natural number '0'`,
	})
	test(&genCase{
		Pattern: `[abc$]{8}`,
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{16, 16},
		Validate: func(p string) bool {
			for _, c := range p {
				switch c {
				case 'a', 'b', 'c', '$':
				default:
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: `[ab}{]{8}`,
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{16, 16},
		Validate: func(p string) bool {
			for _, c := range p {
				switch c {
				case 'a', 'b', '}', '{':
				default:
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: `$rjust([ab}{]{8}, 10)`,
		PassLen: [2]int{10, 10},
		Entropy: [2]float64{16, 16},
		Validate: func(p string) bool {
			for _, c := range p {
				switch c {
				case 'a', 'b', '}', '{', ' ':
				default:
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: `[abcccdab]{8}`,
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
		Pattern: `[a-z]{8}`,
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
		Pattern: `[a-\u007a]{8}`,
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
		Pattern:  `\U000103a0 \U000103c3`,
		PassLen:  [2]int{3, 3},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`𐎠 𐏃`),
	})
	test(&genCase{
		Pattern: `[\U000103a0-\U000103c3]{8}`,
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{41.3, 41.4},
		Validate: func(p string) bool {
			for _, c := range p {
				if c < '𐎠' || c > '𐏃' {
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: `[\u0009-\u000a]{8}`,
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{8, 8},
		Validate: func(p string) bool {
			for _, c := range p {
				if c < '\t' || c > '\n' {
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: `[\t-\n]{8}`,
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{8, 8},
		Validate: func(p string) bool {
			for _, c := range p {
				if c < '\t' || c > '\n' {
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: `[a-z]{8,10}`,
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
		Pattern:  `[]`,
		PassLen:  [2]int{0, 0},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(``),
	})
	test(&genCase{
		Pattern: `[a-z]{8}[1-9]{3}`,
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
		Pattern: `([a-z]{8}[1-9]{3})`,
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
		Pattern: `([a-z]{5}[1-9]{2}){2}`,
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
	testErr(&genErrCase{
		Pattern: `abc(test\`,
		Error:   `         ^ syntax error: '(' not closed`,
	})
	test(&genCase{
		Pattern:  `(\)){2}`,
		PassLen:  [2]int{2, 2},
		Entropy:  [2]float64{0, 0},
		Password: strPtr("))"),
	})
	test(&genCase{
		Pattern:  `(\\){2}`,
		PassLen:  [2]int{2, 2},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`\\`),
	})
	test(&genCase{
		Pattern:  `(\\\)\(){2}`,
		PassLen:  [2]int{6, 6},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`\)(\)(`),
	})
	test(&genCase{
		Pattern: `([a-z]{5}[1-9]{2}-){2}`,
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
		Pattern: `(ab|cd|ef|gh){8}`,
		PassLen: [2]int{16, 16},
		Entropy: [2]float64{16, 16},
		Validate: func(p string) bool {
			for i := 0; i < len(p); i += 2 {
				switch p[i : i+2] {
				case "ab", "cd", "ef", "gh":
				default:
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: `$?(a)$?(b)$?(c)$?(d)`,
		PassLen: [2]int{0, 4},
		Entropy: [2]float64{4, 4},
		Validate: func(p string) bool {
			for _, c := range p {
				switch c {
				case 'a', 'b', 'c', 'd':
				default:
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
		Validate: func(p string) bool {
			for _, c := range p {
				switch c {
				case ' ', '!', '"', '#', '$', '%', '&', '\'', '(', ')', '*',
					'+', ',', '-', '.', '/', ':', ';', '<', '=', '>', '?', '@',
					'[', '\\', ']', '^', '_', '`', '{', '|', '}', '~':
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: `[\^abc]{8}`,
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{16, 16},
		Validate: func(p string) bool {
			for _, c := range p {
				switch c {
				case '^', 'a', 'b', 'c':
				default:
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: `[^^]{10}`,
		PassLen: [2]int{10, 10},
		Entropy: [2]float64{65.5, 65.6},
		Validate: func(p string) bool {
			for _, c := range p {
				switch c {
				case '^':
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: `[!-~]{10}`,
		PassLen: [2]int{10, 10},
		Entropy: [2]float64{65.5, 65.6},
		Validate: func(p string) bool {
			for _, c := range p {
				if c < '!' {
					return false
				}
				if c > '~' {
					return false
				}
			}
			return true
		},
	})
	// base64 length: ((bytes + 2) / 3) * 4
	test(&genCase{
		Pattern: `$base64([:alnum:]{10})`,
		PassLen: [2]int{16, 16},
		Entropy: [2]float64{59.5, 59.6},
		Validate: func(p string) bool {
			pwBytes, err := base64.StdEncoding.DecodeString(p)
			if err != nil {
				t.Logf("p=%#v", p)
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
		Pattern: `$base64([:alnum:]{9})`,
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
		Pattern: `$base64([:alnum:]{5})`,
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
		Pattern: `$base64url([:alnum:]{5})`,
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
		Pattern: `$base32([:alnum:]{5})`,
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
		Pattern: `$BASE32([:alnum:]{5})`,
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
		Pattern: `$base32std([:alnum:]{5})`,
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
		Pattern: `$hex([:alnum:]{8})`,
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
		Pattern: `$HEX([:alnum:]{8})`,
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
		Pattern: `$hex([a-c)(]{4})`,
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
		Pattern: `$hex(([a-e]{4}))`,
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
		Pattern:  `$hex2dec(616263)`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`6382179`),
	})
	testErr(&genErrCase{
		Pattern: `$hex2dec(abcdefg)`,
		Error:   `               ^ value error: invalid hex number "abcdefg"`,
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
	testErr(&genErrCase{
		Pattern: `$bip39word(abcd)`,
		Error:   `               ^ value error: invalid number 'abcd'`,
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
	testErr(&genErrCase{
		Pattern: `$()`,
		Error:   ` ^ syntax error: missing function name`,
	})
	testErr(&genErrCase{
		Pattern: `$hex([a-z]`,
		Error:   `          ^ syntax error: '(' not closed`,
	})
	testErr(&genErrCase{
		Pattern: `$hex(([a-z]`,
		Error:   `           ^ syntax error: '(' not closed`,
	})
	testErr(&genErrCase{
		Pattern: `$hex([:x:])`,
		Error:   `       ^ value error: invalid character class "x"`,
	})
	testErr(&genErrCase{
		Pattern: `(`,
		Error:   ` ^ syntax error: '(' not closed`,
	})
	testErr(&genErrCase{
		Pattern: `$foo`,
		Error:   `    ^ syntax error: expected a function call`,
	})
	testErr(&genErrCase{
		Pattern: `($foo`,
		Error:   `     ^ syntax error: '(' not closed`,
	})
	testErr(&genErrCase{
		Pattern: `$foo(`,
		Error:   `     ^ syntax error: '(' not closed`,
	})
	testErr(&genErrCase{
		Pattern: `$foo(123)`,
		Error:   `    ^ value error: invalid function 'foo'`,
	})
	testErr(&genErrCase{
		Pattern: `$foo\()`,
		Error:   `    ^ syntax error: expected a function call`,
	})
	testErr(&genErrCase{
		Pattern: `test($foo(123))`,
		Error:   `         ^ value error: invalid function 'foo'`,
	})
	testErr(&genErrCase{
		Pattern: `test $foo`,
		Error:   `         ^ syntax error: expected a function call`,
	})
	testErr(&genErrCase{
		Pattern: `test($foo)`,
		Error:   `         ^ syntax error: expected a function call`,
	})
	testErr(&genErrCase{
		Pattern: `$rjust(a,10000)[`,
		Error:   `                ^ syntax error: '[' not closed`,
	})
	testErr(&genErrCase{
		Pattern: `$rjust(a,10000)[a-]`,
		Error:   `                  ^ syntax error: no character after '-'`,
	})
	testErr(&genErrCase{
		Pattern: `(((a{10,20})))[`,
		Error:   `               ^ syntax error: '[' not closed`,
	})
	testErr(&genErrCase{
		Pattern: `(((a{10,20})))[a-]`,
		Error:   `                 ^ syntax error: no character after '-'`,
	})
	test(&genCase{
		Pattern: `$shuffle([a-z]{5}[1-9]{2})`,
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
	test(&genCase{
		Pattern:  `\u00e0-\u00e6`,
		PassLen:  [2]int{3, 3},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`à-æ`),
	})
	testErr(&genErrCase{
		Pattern: `\u00e0-\u00e`,
		Error:   `       ^ syntax error: invalid escape sequence`,
	})
	testErr(&genErrCase{
		Pattern: `\u00e0-\U00e6`,
		Error:   `       ^ syntax error: invalid escape sequence`,
	})
	test(&genCase{
		Pattern:  `test1 \u00e1 test2 \u00e2 test3`,
		PassLen:  [2]int{21, 21},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`test1 á test2 â test3`),
	})
	testErr(&genErrCase{
		Pattern: `\u00mn`,
		Error:   `^ syntax error: invalid escape sequence`,
	})
	testErr(&genErrCase{
		Pattern: `test1 \u00mn test2`,
		Error:   `      ^ syntax error: invalid escape sequence`,
	})
	testErr(&genErrCase{
		Pattern: `(test1 \u00mn test2){2}`,
		Error:   `       ^ syntax error: invalid escape sequence`,
	})
	testErr(&genErrCase{
		Pattern: `test[\u00mn-\u00e0]abc`,
		Error:   `     ^ syntax error: invalid escape sequence`,
	})
	testErr(&genErrCase{
		Pattern: `test[\u00e0-\u00mn]abc`,
		Error:   `            ^ syntax error: invalid escape sequence`,
	})
	test(&genCase{
		Pattern: `[\u00e0-\u00e6]{10}`,
		PassLen: [2]int{10, 10},
		Entropy: [2]float64{28, 29},
		Validate: func(p string) bool {
			for _, c := range p {
				if !(c >= 'à' && c <= 'æ') {
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: `$date(2000,2020,-)`,
		PassLen: [2]int{10, 10},
		Entropy: [2]float64{12.8, 12.9},
	})
	test(&genCase{
		Pattern: `$date(2000,2020)`,
		PassLen: [2]int{10, 10},
		Entropy: [2]float64{12.8, 12.9},
	})
	test(&genCase{
		Pattern: `$date(2000,2020,\,)`,
		PassLen: [2]int{10, 10},
		Entropy: [2]float64{12.8, 12.9},
	})
	testErr(&genErrCase{
		Pattern: `$date()`,
		Error:   `      ^ argument error: date: too few characters as arguments`,
	})
	testErr(&genErrCase{
		Pattern: `$date(2000)`,
		Error:   `          ^ argument error: date: at least 2 arguments are required`,
	})
	testErr(&genErrCase{
		Pattern: `$date(2000a,2000b)`,
		Error:   `          ^ value error: invalid year 2000a`,
	})
	testErr(&genErrCase{
		Pattern: `$date(2000,2000b)`,
		Error:   `               ^ value error: invalid year 2000b`,
	})
	test(&genCase{
		Pattern:  `$space()`,
		PassLen:  [2]int{0, 0},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(``),
	})
	test(&genCase{
		Pattern:  `$space(abcd)`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`a b c d`),
	})
	test(&genCase{
		Pattern:  `$expand()`,
		PassLen:  [2]int{0, 0},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(``),
	})
	test(&genCase{
		Pattern:  `$expand(|abcd)`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`a|b|c|d`),
	})
	test(&genCase{
		Pattern:  `$romaji()`,
		PassLen:  [2]int{0, 0},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(``),
	})
	test(&genCase{
		Pattern:  `$romaji(そうたい)`,
		PassLen:  [2]int{6, 6},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`soutai`),
	})
	test(&genCase{
		Pattern:  `$romaji(こうげきてき)`,
		PassLen:  [2]int{11, 11},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`kougekiteki`),
	})
	test(&genCase{
		Pattern:  `$romaji(レザーレット)`,
		PassLen:  [2]int{10, 10},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`reza-retto`),
	})
	test(&genCase{
		Pattern:  `$romaji(ーレット)`,
		PassLen:  [2]int{5, 5},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`retto`),
	})
	test(&genCase{
		Pattern:  `$romaji(あかんかった)`,
		PassLen:  [2]int{9, 9},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`akankatta`),
	})
	test(&genCase{
		Pattern:  `$romaji(あかんかっった)`,
		PassLen:  [2]int{9, 9},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`akankatta`),
	})
	test(&genCase{
		Pattern:  `$romaji(累減税)`,
		PassLen:  [2]int{3, 3},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`累減税`),
	})
	test(&genCase{
		Pattern:  `$romaji(test123)`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`test123`),
	})
	test(&genCase{
		Pattern:  `$rjust(abc,7)`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`    abc`),
	})
	test(&genCase{
		Pattern:  `$rjust(abc,2)`,
		PassLen:  [2]int{3, 3},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`abc`),
	})
	test(&genCase{
		Pattern:  `$rjust(abc,7,0)`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`0000abc`),
	})
	test(&genCase{
		Pattern: `$rjust([a-z]{5},7,0)`,
		PassLen: [2]int{7, 7},
		Entropy: [2]float64{23.5, 23.6},
	})
	test(&genCase{
		Pattern:  `$rjust((abc,),7,0)`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`000abc,`),
	})
	test(&genCase{
		Pattern:  `$rjust(abc\,,7,0)`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`000abc,`),
	})
	test(&genCase{
		Pattern: `$rjust([abc]{3},7,0)`,
		PassLen: [2]int{7, 7},
		Entropy: [2]float64{4.7, 4.8},
		Validate: func(p string) bool {
			for _, c := range p {
				switch c {
				case 'a', 'b', 'c', '0':
				default:
					return false
				}
			}
			return true
		},
	})
	test(&genCase{
		Pattern: `$rjust([)(}]{3},7,0)`,
		PassLen: [2]int{7, 7},
		Entropy: [2]float64{4.7, 4.8},
		Validate: func(p string) bool {
			for _, c := range p {
				if c == ')' || c == '(' || c == '}' {
					return true
				}
				if c == '0' {
					return true
				}
				return false
			}
			return true
		},
	})
	test(&genCase{
		Pattern:  `abc\(`,
		PassLen:  [2]int{4, 4},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`abc(`),
	})
	test(&genCase{
		Pattern:  `$rjust(abc\(,7,0)`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`000abc(`),
	})
	testErr(&genErrCase{
		Pattern: `$rjust(abc)`,
		Error:   `          ^ argument error: rjust: at least 2 arguments are required`,
	})
	testErr(&genErrCase{
		Pattern: `$rjust(abc,a)`,
		Error:   `           ^ value error: invalid width a`,
	})
	testErr(&genErrCase{
		Pattern: `$rjust(abc,0)`,
		Error:   `           ^ value error: invalid width 0`,
	})
	testErr(&genErrCase{
		Pattern: `$rjust(abc,1,ab)`,
		Error:   `              ^ value error: invalid fillChar="ab", must have length 1`,
	})
	testErr(&genErrCase{
		Pattern: `$rjust({{}},7)`,
		Error:   fmt.Errorf(`nested '{'`),
	})
	test(&genCase{
		Pattern:  `$ljust((abc,),7,0)`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`abc,000`),
	})
	test(&genCase{
		Pattern:  `$ljust((abc,),7)`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`abc,   `),
	})
	test(&genCase{
		Pattern:  `$center((abc,),7)`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(` abc,  `),
	})
	test(&genCase{
		Pattern:  `$center((abc,),8)`,
		PassLen:  [2]int{8, 8},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`  abc,  `),
	})
	test(&genCase{
		Pattern:  `(abc) test1 \1 test2`,
		PassLen:  [2]int{19, 19},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`abc test1 abc test2`),
	})
	testErr(&genErrCase{
		Pattern: `(abc) test1 \2 test2`,
		Error:   `             ^ value error: invalid group id '2'`,
	})
	testErr(&genErrCase{
		Pattern: `(abc) test1 \20 test2`,
		Error:   `              ^ value error: invalid group id '20'`,
	})
	test(&genCase{
		Pattern:  `(a(b(c))) 1:'\1' 2:'\2' 3:'\3'`,
		PassLen:  [2]int{24, 24},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`abc 1:'abc' 2:'bc' 3:'c'`),
	})
	test(&genCase{
		Pattern:  `$hex((abc)) 1:'\1'`,
		PassLen:  [2]int{14, 14},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`616263 1:'abc'`),
	})
	test(&genCase{
		Pattern:  `kana: (そうたい) romaji: $romaji(\1)`,
		PassLen:  [2]int{25, 25},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`kana: そうたい romaji: soutai`),
	})
}
