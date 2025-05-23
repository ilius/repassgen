package passgen_test

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ilius/bip39-coder/bip39"
	"github.com/ilius/is/v2"
	passgen "github.com/ilius/repassgen/lib"
	"github.com/ilius/repassgen/lib/crock32"
)

const wordChars = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_`

var (
	verbose      = os.Getenv("TEST_VERBOSE") == "1"
	bip39WordMap = getBIP39WordMap()
)

func decodeHex(s string) []byte {
	pwBytes, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return pwBytes
}

func getBIP39WordMap() map[string]bool {
	count := bip39.WordCount()
	m := make(map[string]bool, count)
	for i := range count {
		word, ok := bip39.GetWord(i)
		if !ok {
			panic("NOT OK")
		}
		m[word] = true
	}
	return m
}

func checkErrorIsInList(is *is.Is, err error, expMsgs []any) {
	if !is.Err(err) {
		return
	}
	// tErr, okErr := err.(*Error)
	// if !okErr {
	is.OneOf(err.Error(), expMsgs...)
}

func testGen(t *testing.T, tc *genCase) {
	is := is.New(t).AddMsg("pattern=%#v", tc.Pattern)
	is = is.Lax()
	out, _, err := passgen.Generate(passgen.GenerateInput{Pattern: []rune(tc.Pattern)})
	if !is.NotErr(err) {
		tErr, okErr := err.(*passgen.Error)
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

func testGenErr(t *testing.T, tc *genErrCase) {
	is := is.New(t).AddMsg("pattern=%#v", tc.Pattern)
	is = is.Lax()
	out, _, err := passgen.Generate(passgen.GenerateInput{Pattern: []rune(tc.Pattern)})
	tErr, okErr := err.(*passgen.Error)
	switch expErr := tc.Error.(type) {
	case string:
		if okErr {
			expErrTyped := passgen.ParseSpacedError(expErr)
			if expErrTyped == nil {
				t.Errorf("bad spaced error %#v", expErr)
				is.Equal(tErr.SpacedError(), expErr)
				return
			}
			is.Equal(tErr.Type(), expErrTyped.Type())
			is.Equal(tErr.Message(), expErrTyped.Message())
			is = is.AddMsg(
				"msg=%#v", tErr.Message(),
			)
			is.AddMsg("mismatch pos").Equal(tErr.Pos(), expErrTyped.Pos())
			is.AddMsg("mismatch markLen").Equal(tErr.MarkLen(), expErrTyped.MarkLen())
		} else {
			is.ErrMsg(err, expErr)
		}
	case []any:
		checkErrorIsInList(is, err, expErr)
	case *passgen.Error:
		is.Equal(tErr.Type(), expErr.Type())
		is.Equal(tErr.Message(), expErr.Message())
		is = is.AddMsg(
			"msg=%#v", tErr.Message(),
		)
		is.AddMsg("mismatch pos").Equal(tErr.Pos(), expErr.Pos())
		is.AddMsg("mismatch markLen").Equal(tErr.MarkLen(), expErr.MarkLen())
	case error:
		is.ErrMsg(err, expErr.Error())
	default:
		panic(fmt.Errorf("invalid type %T for Error: %v", tc.Error, tc.Error))
	}
	if okErr && verbose {
		t.Log(tc.Pattern)
		t.Log(tErr.SpacedError())
	}
	if !is.Nil(out) {
		t.Log(string(out.Password))
	}
	if verbose {
		t.Log("------------------------------------")
	}
}

func TestGenerateCharClassShort(t *testing.T) {
	testGen(t, &genCase{
		Pattern: `\w{4}`,
		PassLen: [2]int{4, 4},
		Entropy: [2]float64{23.9, 23.91},
		Validate: func(p string) bool {
			for _, c := range p {
				if !strings.ContainsRune(wordChars, c) {
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern: `\d{4}`,
		PassLen: [2]int{4, 4},
		Entropy: [2]float64{13.2, 13.3},
		Validate: func(p string) bool {
			for _, c := range p {
				if c < '0' || c > '9' {
					return false
				}
			}
			return true
		},
	})
}

func TestGenerateCharClass(t *testing.T) {
	testGen(t, &genCase{
		Pattern: `[:alpha:]{4}`,
		PassLen: [2]int{4, 4},
		Entropy: [2]float64{22.8, 22.81},
		Validate: func(p string) bool {
			for _, c := range p {
				if !strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", c) {
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern: `[:alnum:]{4}`,
		PassLen: [2]int{4, 4},
		Entropy: [2]float64{23.81, 23.82},
		Validate: func(p string) bool {
			for _, c := range p {
				if !strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", c) {
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern: `[:word:]{4}`,
		PassLen: [2]int{4, 4},
		Entropy: [2]float64{23.9, 23.91},
		Validate: func(p string) bool {
			for _, c := range p {
				if !strings.ContainsRune(wordChars, c) {
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern: `[:w:]{4}`,
		PassLen: [2]int{4, 4},
		Entropy: [2]float64{23.9, 23.91},
		Validate: func(p string) bool {
			for _, c := range p {
				if !strings.ContainsRune(wordChars, c) {
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern: `[:lower:]{4}`,
		PassLen: [2]int{4, 4},
		Entropy: [2]float64{18.8, 18.81},
		Validate: func(p string) bool {
			for _, c := range p {
				if c < 'a' || c > 'z' {
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern: `[:upper:]{4}`,
		PassLen: [2]int{4, 4},
		Entropy: [2]float64{18.8, 18.81},
		Validate: func(p string) bool {
			for _, c := range p {
				if c < 'A' || c > 'Z' {
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern: `[:digit:]{4}`,
		PassLen: [2]int{4, 4},
		Entropy: [2]float64{13.28, 13.29},
		Validate: func(p string) bool {
			for _, c := range p {
				if c < '0' || c > '9' {
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern: `[:d:]{4}`,
		PassLen: [2]int{4, 4},
		Entropy: [2]float64{13.28, 13.29},
		Validate: func(p string) bool {
			for _, c := range p {
				if c < '0' || c > '9' {
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern: `[:xdigit:]{4}`,
		PassLen: [2]int{4, 4},
		Entropy: [2]float64{17.83, 17.84},
		Validate: func(p string) bool {
			for _, c := range p {
				if !strings.ContainsRune("0123456789abcdefABCDEF", c) {
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern: `[:punct:]{4}`,
		PassLen: [2]int{4, 4},
		Entropy: [2]float64{20, 20},
		Validate: func(p string) bool {
			for _, c := range p {
				if !strings.ContainsRune("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~", c) {
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern: `[:b32:]{4}`,
		PassLen: [2]int{4, 4},
		Entropy: [2]float64{20, 20},
		Validate: func(p string) bool {
			for _, c := range p {
				if !strings.ContainsRune("0123456789abcdefghjkmnpqrstvwxyz", c) {
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern: `[:B32:]{4}`,
		PassLen: [2]int{4, 4},
		Entropy: [2]float64{20, 20},
		Validate: func(p string) bool {
			for _, c := range p {
				if !strings.ContainsRune("0123456789ABCDEFGHJKMNPQRSTVWXYZ", c) {
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern: `[:B32STD:]{4}`,
		PassLen: [2]int{4, 4},
		Entropy: [2]float64{20, 20},
		Validate: func(p string) bool {
			for _, c := range p {
				if !strings.ContainsRune("ABCDEFGHIJKLMNOPQRSTUVWXYZ234567", c) {
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern: `[:b64:]{4}`,
		PassLen: [2]int{4, 4},
		Entropy: [2]float64{24, 24},
		Validate: func(p string) bool {
			for _, c := range p {
				if !strings.ContainsRune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/", c) {
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern: `[:b64url:]{4}`,
		PassLen: [2]int{4, 4},
		Entropy: [2]float64{24, 24},
		Validate: func(p string) bool {
			for _, c := range p {
				if !strings.ContainsRune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_", c) {
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern: `[:print:]{4}`,
		PassLen: [2]int{4, 4},
		Entropy: [2]float64{26.27, 26.28},
	})
	testGen(t, &genCase{
		Pattern: `[:graph:]{4}`,
		PassLen: [2]int{4, 4},
		Entropy: [2]float64{26.21, 26.22},
	})
	testGen(t, &genCase{
		Pattern: `[:space:]{4}`,
		PassLen: [2]int{4, 4},
		Entropy: [2]float64{10.33, 10.34},
		Validate: func(p string) bool {
			for _, c := range p {
				switch c {
				case ' ', '\t', '\r', '\n', '\v', '\f':
				default:
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern: `[:blank:]{4}`,
		PassLen: [2]int{4, 4},
		Entropy: [2]float64{4, 4},
		Validate: func(p string) bool {
			for _, c := range p {
				switch c {
				case ' ', '\t':
				default:
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern: `[:cntrl:]{4}`,
		PassLen: [2]int{4, 4},
		Entropy: [2]float64{20.17, 20.18},
		Validate: func(p string) bool {
			for _, c := range p {
				switch {
				case c >= '\x00' && c <= '\x1F':
				case c == '\x7F':
				default:
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern: `[:ascii:]{4}`,
		PassLen: [2]int{4, 4},
		Entropy: [2]float64{28, 28},
		Validate: func(p string) bool {
			for _, c := range p {
				switch {
				case c <= '\x7F':
				default:
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
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
	testGen(t, &genCase{
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
	testGen(t, &genCase{
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
	testGen(t, &genCase{
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
	testGen(t, &genCase{
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
	testGen(t, &genCase{
		Pattern:  `\U000103a0 \U000103c3`,
		PassLen:  [2]int{3, 3},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`𐎠 𐏃`),
	})
	testGen(t, &genCase{
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
	testGen(t, &genCase{
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
	testGen(t, &genCase{
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
	testGen(t, &genCase{
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
	testGen(t, &genCase{
		Pattern:  `[]`,
		PassLen:  [2]int{0, 0},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(``),
	})
	testGen(t, &genCase{
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
	testGen(t, &genCase{
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
	testGen(t, &genCase{
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
	testGen(t, &genCase{
		Pattern: `[^^]{10}`,
		PassLen: [2]int{10, 10},
		Entropy: [2]float64{65.5, 65.6},
		Validate: func(p string) bool {
			for _, c := range p {
				if c == '^' {
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
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
}

func TestGenerateAlter(t *testing.T) {
	testGen(t, &genCase{
		Pattern: `(ab|()){8}`,
		PassLen: [2]int{0, 16},
		Entropy: [2]float64{8, 8},
	})
	testGen(t, &genCase{
		Pattern: `(ab|c\\)`,
		PassLen: [2]int{2, 2},
		Entropy: [2]float64{1, 1},
	})
	testGen(t, &genCase{
		Pattern: `(ab|c\)`,
		PassLen: [2]int{2, 2},
		Entropy: [2]float64{1, 1},
	})
}

func TestGenerateGroups(t *testing.T) {
	testGen(t, &genCase{
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
	testGen(t, &genCase{
		Pattern: `([a-z]{5}[1-9]{2}){2}`,
		PassLen: [2]int{14, 14},
		Entropy: [2]float64{59.6, 59.7},
		Validate: func(p string) bool {
			for i := range 2 {
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
	testGen(t, &genCase{
		Pattern:  `(\)){2}`,
		PassLen:  [2]int{2, 2},
		Entropy:  [2]float64{0, 0},
		Password: strPtr("))"),
	})
	testGen(t, &genCase{
		Pattern:  `(\\){2}`,
		PassLen:  [2]int{2, 2},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`\\`),
	})
	testGen(t, &genCase{
		Pattern:  `(\\\)\(){2}`,
		PassLen:  [2]int{6, 6},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`\)(\)(`),
	})
	testGen(t, &genCase{
		Pattern: `([a-z]{5}[1-9]{2}-){2}`,
		PassLen: [2]int{16, 16},
		Entropy: [2]float64{59.6, 59.7},
		Validate: func(p string) bool {
			for i := range 2 {
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

	// alteration
	testGen(t, &genCase{
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
	testGen(t, &genCase{
		Pattern: `(ab|\\c){8}`,
		PassLen: [2]int{16, 16},
		Entropy: [2]float64{8, 8},
		Validate: func(p string) bool {
			for i := 0; i < len(p); i += 2 {
				switch p[i : i+2] {
				case "ab", `\c`:
				default:
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern: `(ab|()){8}`,
		PassLen: [2]int{16, 16},
		Entropy: [2]float64{8, 8},
		Validate: func(p string) bool {
			for i := 0; i < len(p); i += 2 {
				switch p[i : i+2] {
				case "ab", `()`:
				default:
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern: `(ab|c\)`,
		PassLen: [2]int{2, 2},
		Entropy: [2]float64{1, 1},
		Validate: func(p string) bool {
			for i := 0; i < len(p); i += 2 {
				switch p[i : i+2] {
				case "ab", `c\`:
				default:
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern: `(a|b|[cde]|f){8}`,
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{16, 16},
		Validate: func(p string) bool {
			for i := range len(p) {
				switch p[i : i+1] {
				case "a", "b", "c", "d", "e", "f":
				default:
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern: `([ab]|[cdef]|[gh]|[ij])`,
		PassLen: [2]int{1, 1},
		Entropy: [2]float64{3, 3},
		Validate: func(p string) bool {
			for i := range len(p) {
				switch p[i : i+1] {
				case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j":
				default:
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern: `([ab]|[cdef]|[gh]|[ij]){8}`,
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{24, 24},
		Validate: func(p string) bool {
			for i := range len(p) {
				switch p[i : i+1] {
				case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j":
				default:
					return false
				}
			}
			return true
		},
	})
}

func TestGenerateFuncJustify(t *testing.T) {
	testGen(t, &genCase{
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
	testGen(t, &genCase{
		Pattern:  `$rjust(abc,7)`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`    abc`),
	})
	testGen(t, &genCase{
		Pattern:  `$rjust(abc,2)`,
		PassLen:  [2]int{3, 3},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`abc`),
	})
	testGen(t, &genCase{
		Pattern:  `$rjust(abc,7,0)`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`0000abc`),
	})
	testGen(t, &genCase{
		Pattern: `$rjust([a-z]{5},7,0)`,
		PassLen: [2]int{7, 7},
		Entropy: [2]float64{23.5, 23.6},
	})
	testGen(t, &genCase{
		Pattern:  `$rjust((abc,),7,0)`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`000abc,`),
	})
	testGen(t, &genCase{
		Pattern:  `$rjust(abc\,,7,0)`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`000abc,`),
	})
	testGen(t, &genCase{
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
	testGen(t, &genCase{
		Pattern: `$rjust([)(}]{3},7,0)`,
		PassLen: [2]int{7, 7},
		Entropy: [2]float64{4.7, 4.8},
		Validate: func(p string) bool {
			for _, c := range p {
				switch {
				case c == ')' || c == '(' || c == '}':
				case c == '0':
				default:
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern:  `$rjust(abc\(,7,0)`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`000abc(`),
	})

	testGen(t, &genCase{
		Pattern:  `$ljust(abc,2)`,
		PassLen:  [2]int{3, 3},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`abc`),
	})
	testGen(t, &genCase{
		Pattern:  `$ljust((abc,),7,0)`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`abc,000`),
	})
	testGen(t, &genCase{
		Pattern:  `$ljust((abc,),7)`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`abc,   `),
	})
	testGen(t, &genCase{
		Pattern:  `$center(abc,2)`,
		PassLen:  [2]int{3, 3},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`abc`),
	})
	testGen(t, &genCase{
		Pattern:  `$center((abc,),7)`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(` abc,  `),
	})
	testGen(t, &genCase{
		Pattern:  `$center((abc,),8)`,
		PassLen:  [2]int{8, 8},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`  abc,  `),
	})
}

func TestGenerateFuncBase64(t *testing.T) {
	// base64 length: ((bytes + 2) / 3) * 4
	testGen(t, &genCase{
		Pattern: `$base64($hex([:alnum:]{10}))`,
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
	testGen(t, &genCase{
		Pattern: `$base64($hex([:alnum:]{9}))`,
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
	testGen(t, &genCase{
		Pattern: `$base64($hex([:alnum:]{5}))`,
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

	testGen(t, &genCase{
		Pattern: `$base64url($hex([:alnum:]{5}))`,
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
}

func TestGenerateFuncBase32(t *testing.T) {
	testGen(t, &genCase{
		Pattern: `$base32($hex([:alnum:]{5}))`,
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{29.7, 29.8},
		Validate: func(p string) bool {
			if strings.ToLower(p) != p {
				return false
			}
			pwBytes, err := crock32.DecodeStrings(p)
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
	testGen(t, &genCase{
		Pattern: `$BASE32($hex([:alnum:]{5}))`,
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{29.7, 29.8},
		Validate: func(p string) bool {
			if strings.ToUpper(p) != p {
				return false
			}
			pwBytes, err := crock32.DecodeStrings(p)
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
	testGen(t, &genCase{
		Pattern: `$base32std($hex([:alnum:]{5}))`,
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
}

func TestGenerateFuncHex(t *testing.T) {
	testGen(t, &genCase{
		Pattern: `$hex([:alnum:]{8})`,
		PassLen: [2]int{16, 16},
		Entropy: [2]float64{47.6, 47.7},
		Validate: func(p string) bool {
			if strings.ToLower(p) != p {
				return false
			}
			pwBytes := decodeHex(p)
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
	testGen(t, &genCase{
		Pattern: `$HEX([:alnum:]{8})`,
		PassLen: [2]int{16, 16},
		Entropy: [2]float64{47.6, 47.7},
		Validate: func(p string) bool {
			if strings.ToUpper(p) != p {
				return false
			}
			pwBytes := decodeHex(p)
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
	testGen(t, &genCase{
		Pattern: `$hex([a-c)(]{4})`,
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{9.28, 9.29},
		Validate: func(p string) bool {
			if strings.ToLower(p) != p {
				return false
			}
			pwBytes := decodeHex(p)
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
	testGen(t, &genCase{
		Pattern: `$hex(([a-e]{4}))`,
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{9.28, 9.29},
		Validate: func(p string) bool {
			if strings.ToLower(p) != p {
				return false
			}
			pwBytes := decodeHex(p)
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
	testGen(t, &genCase{
		Pattern:  `$hex2dec(616263)`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`6382179`),
	})
}

func TestGenerateFuncBIP39(t *testing.T) {
	// each bip39 word is at least 3 chars, and at most 8 chars
	testGen(t, &genCase{
		Pattern:   "$bip39word(10)",
		WordCount: 10,
		PassLen:   [2]int{39, 89}, // 10*4 - 1, 10*9 - 1
		Entropy:   [2]float64{110, 110},
		Validate: func(p string) bool {
			words := strings.Split(p, " ")
			if len(words) != 10 {
				return false
			}
			for _, word := range words {
				if !bip39WordMap[word] {
					return false
				}
			}
			return true
		},
	})
	testGen(t, &genCase{
		Pattern:   "$bip39word()",
		WordCount: 1,
		PassLen:   [2]int{3, 8},
		Entropy:   [2]float64{11, 11},
		Validate: func(p string) bool {
			return bip39WordMap[p]
		},
	})

	// 1 bip39 word   => 11 bits entropy
	// 8 bip39 words  => 11 bytes (88 bits) entropy
	// but there is also a checksum
	// https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki#generating-the-mnemonic
	// Criticism: https://github.com/bitcoin/bips/wiki/Comments:BIP-0039
	// The following table describes the relation between the initial
	// entropy length (ENT), the checksum length (CS), and the length of
	// the generated mnemonic sentence (MS) in words.
	// CS = ENT / 32
	// MS = (ENT + CS) / 11 = (ENT + ENT/32) / 11
	testGen(t, &genCase{
		Pattern: "$bip39encode($byte(){8})",
		PassLen: [2]int{23, 62}, // 6*4-1, 7*9-1
		Entropy: [2]float64{64, 64},
		Validate: func(p string) bool {
			words := strings.Split(p, " ")
			if len(words) < 6 {
				return false
			}
			if len(words) > 7 {
				return false
			}
			if len(words) == 7 && words[len(words)-1] != "abandon" {
				return false
			}
			for _, word := range words {
				if !bip39WordMap[word] {
					return false
				}
			}
			return true
		},
	})
}

func TestGenerateOK(t *testing.T) {
	testGen(t, &genCase{
		Pattern: ``,
		PassLen: [2]int{0, 0},
		Entropy: [2]float64{0, 0},
	})

	testGen(t, &genCase{
		Pattern:  `abc\(`,
		PassLen:  [2]int{4, 4},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`abc(`),
	})

	testGen(t, &genCase{
		Pattern:  `\`,
		PassLen:  [2]int{0, 0},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(""),
	})
	testGen(t, &genCase{
		Pattern:  `abc\`,
		PassLen:  [2]int{3, 3},
		Entropy:  [2]float64{0, 0},
		Password: strPtr("abc"),
	})

	testGen(t, &genCase{
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

	testGen(t, &genCase{
		Pattern:  `$escape(")`,
		PassLen:  [2]int{2, 2},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`\"`),
	})
	testGen(t, &genCase{
		Pattern:  `a[\t][\r][\n][\v][\f]b`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr("a\t\r\n\v\fb"),
	})
	testGen(t, &genCase{
		Pattern:  `a\t\r\n\v\fb\c`,
		PassLen:  [2]int{8, 8},
		Entropy:  [2]float64{0, 0},
		Password: strPtr("a\t\r\n\v\fbc"),
	})

	testGen(t, &genCase{
		Pattern: `$byte(){4}`,
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{32, 32},
		Validate: func(p string) bool {
			if strings.ToLower(p) != p {
				return false
			}
			return len(decodeHex(p)) == 4
		},
	})
	testGen(t, &genCase{
		Pattern: `$BYTE(){4}`,
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{32, 32},
		Validate: func(p string) bool {
			if strings.ToUpper(p) != p {
				return false
			}
			return len(decodeHex(p)) == 4
		},
	})

	testGen(t, &genCase{
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
	testGen(t, &genCase{
		Pattern:  `\u00e0-\u00e6`,
		PassLen:  [2]int{3, 3},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`à-æ`),
	})

	testGen(t, &genCase{
		Pattern:  `test1 \u00e1 test2 \u00e2 test3`,
		PassLen:  [2]int{21, 21},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`test1 á test2 â test3`),
	})
	testGen(t, &genCase{
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
	testGen(t, &genCase{
		Pattern: `$date(2000,2020,-)`,
		PassLen: [2]int{10, 10},
		Entropy: [2]float64{12.8, 12.9},
	})
	testGen(t, &genCase{
		Pattern: `$date(2000,2020)`,
		PassLen: [2]int{10, 10},
		Entropy: [2]float64{12.8, 12.9},
	})
	testGen(t, &genCase{
		Pattern: `$date(2000,2020,\,)`,
		PassLen: [2]int{10, 10},
		Entropy: [2]float64{12.8, 12.9},
	})

	testGen(t, &genCase{
		Pattern:  `$space()`,
		PassLen:  [2]int{0, 0},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(``),
	})
	testGen(t, &genCase{
		Pattern:  `$space(abcd)`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`a b c d`),
	})
	testGen(t, &genCase{
		Pattern:  `$expand()`,
		PassLen:  [2]int{0, 0},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(``),
	})
	testGen(t, &genCase{
		Pattern:  `$expand(|abcd)`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`a|b|c|d`),
	})
	testGen(t, &genCase{
		Pattern:  `$romaji()`,
		PassLen:  [2]int{0, 0},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(``),
	})
	testGen(t, &genCase{
		Pattern:  `$romaji(そうたい)`,
		PassLen:  [2]int{6, 6},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`soutai`),
	})
	testGen(t, &genCase{
		Pattern:  `$romaji(こうげきてき)`,
		PassLen:  [2]int{11, 11},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`kougekiteki`),
	})
	testGen(t, &genCase{
		Pattern:  `$romaji(レザーレット)`,
		PassLen:  [2]int{10, 10},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`reza-retto`),
	})
	testGen(t, &genCase{
		Pattern:  `$romaji(ーレット)`,
		PassLen:  [2]int{5, 5},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`retto`),
	})
	testGen(t, &genCase{
		Pattern:  `$romaji(あかんかった)`,
		PassLen:  [2]int{9, 9},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`akankatta`),
	})
	testGen(t, &genCase{
		Pattern:  `$romaji(あかんかっった)`,
		PassLen:  [2]int{9, 9},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`akankatta`),
	})
	testGen(t, &genCase{
		Pattern:  `$romaji(累減税)`,
		PassLen:  [2]int{3, 3},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`累減税`),
	})
	testGen(t, &genCase{
		Pattern:  `$romaji(test123)`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`test123`),
	})
	testGen(t, &genCase{
		Pattern:  `$json(test)`,
		PassLen:  [2]int{4, 4},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`test`),
	})
	testGen(t, &genCase{
		Pattern:  `$json("")`,
		PassLen:  [2]int{4, 4},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`\"\"`),
	})
	testGen(t, &genCase{
		Pattern:  `(abc) test1 \1 test2`,
		PassLen:  [2]int{19, 19},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`abc test1 abc test2`),
	})

	testGen(t, &genCase{
		Pattern:  `(a(b(c))) 1:'\1' 2:'\2' 3:'\3'`,
		PassLen:  [2]int{24, 24},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`abc 1:'abc' 2:'bc' 3:'c'`),
	})
	testGen(t, &genCase{
		Pattern:  `$hex((abc)) 1:'\1'`,
		PassLen:  [2]int{14, 14},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`616263 1:'abc'`),
	})

	testGen(t, &genCase{
		Pattern:  `$pyhex($hex(test))`,
		PassLen:  [2]int{19, 19}, // byteCount * 4 + 3
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`b'\x74\x65\x73\x74'`),
	})
	testGen(t, &genCase{
		Pattern:  `kana: (そうたい) romaji: $romaji(\1)`,
		PassLen:  [2]int{25, 25},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`kana: そうたい romaji: soutai`),
	})
}

func TestGenerateError(t *testing.T) {
	testGenErr(t, &genErrCase{
		Pattern: `abc(test\`,
		Error:   `         ^ syntax error: '(' not closed`,
	})
	testGenErr(t, &genErrCase{
		// FIXME: if one part of alteration has no error, test becomes flaky
		Pattern: `([:foobar1:]|[:foobar2:])`,
		Error: []any{
			`value error near index 9: invalid character class "foobar1"`,
			`value error near index 20: invalid character class "foobar2"`,
		},
	})
	testGenErr(t, &genErrCase{
		Pattern: `$base64(gh)`,
		Error:   `         ^ value error: invalid hex number "gh"`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$base64url(gh)`,
		Error:   `            ^ value error: invalid hex number "gh"`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$base32(gh)`,
		Error:   `         ^ value error: invalid hex number "gh"`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$BASE32(gh)`,
		Error:   `         ^ value error: invalid hex number "gh"`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$base32std(gh)`,
		Error:   `            ^ value error: invalid hex number "gh"`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `[\`,
		Error:   ` ^ syntax error: '[' not closed`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$\(\(\(`,
		Error:   ` ^ syntax error: expected a function call`,
	})

	testGenErr(t, &genErrCase{
		Pattern: `[0-\\`,
		Error:   `    ^ syntax error: '[' not closed`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `[::\(\(]`,
		Error:   ` ^^ value error: invalid character class ""`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `([::((]|(`,
		Error:   `  ^^ value error: invalid character class ""`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$\(\(\(|0`,
		Error:   ` ^ syntax error: expected a function call`,
	})

	testGenErr(t, &genErrCase{
		Pattern: `test (`,
		Error:   `      ^ syntax error: '(' not closed`,
	})
	// TODO: this should raise error
	// testGenErr(t, &genErrCase{
	// 	Pattern: `test ())`,
	// 	Error:   ``,
	// })
	testGenErr(t, &genErrCase{
		Pattern: `test [`,
		Error:   `      ^ syntax error: '[' not closed`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `test {`,
		Error:   `     ^ syntax error: '{' not closed`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `test {{`,
		Error:   `      ^ syntax error: nested '{'`,
	})

	testGenErr(t, &genErrCase{
		Pattern: `[a`,
		Error:   `  ^ syntax error: '[' not closed`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `[[]]`,
		Error:   ` ^ syntax error: nested '['`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `test [:x]`,
		Error:   `      ^^^ syntax error: ':' not closed`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `test [:abcde]`,
		Error:   `      ^^^^^^^ syntax error: ':' not closed`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `test [:x`,
		Error:   `      ^^^ syntax error: ':' not closed`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `test [a-`,
		Error:   `     ^^^^ syntax error: '[' not closed`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `test [hello-`,
		Error:   `     ^^^^^^^^ syntax error: '[' not closed`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `test [hello-]`,
		Error:   `            ^ syntax error: no character after '-'`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `[-a]`,
		Error:   ` ^ syntax error: no character before '-'`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `[a-z]{a}`,
		Error:   `      ^ syntax error: invalid natural number inside {...}`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `[a-z]{2.5}`,
		Error:   `      ^^ syntax error: invalid natural number inside {...}`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `[a-z]{2.0}`,
		Error:   `      ^^ syntax error: invalid natural number inside {...}`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `[a-z]{1-3}`,
		Error:   `       ^ syntax error: repetition range syntax is '{M,N}' not '{M-N}'`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `[a-z]{1,-3}`,
		Error:   `        ^ syntax error: invalid natural number`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `test([a-z]{1-3})`,
		Error:   `            ^ syntax error: repetition range syntax is '{M,N}' not '{M-N}'`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `test([a-z]{1,})`,
		Error:   `             ^ syntax error: no number after ','`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `test([a-z]{,3333})`,
		Error:   `           ^ syntax error: no number before ','`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `test([a-z]{1,2,3})`,
		Error:   `              ^ syntax error: multiple ',' inside {...}`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `[a-z]{{}`,
		Error:   `      ^ syntax error: nested '{'`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `[a-z]{[}`,
		Error:   `      ^ syntax error: '[' inside {...}`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `[a-z]{$}`,
		Error:   `      ^ syntax error: '$' inside {...}`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `test([a-z]{1a})`,
		Error:   `           ^^ syntax error: invalid natural number inside {...}`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `test([a-z]{})`,
		Error:   `           ^ syntax error: missing number inside {}`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `[a-z]{3,1}`,
		Error:   `      ^^^ value error: invalid numbers 3 > 1 inside {...}`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `{3}`,
		Error:   `^ syntax error: nothing to repeat`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `x{0}`,
		Error:   `  ^ syntax error: invalid natural number '0'`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `x{000}`,
		Error:   `  ^^^ syntax error: invalid natural number '000'`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `test{000000,00000000}`,
		Error:   `     ^^^^^^ syntax error: invalid natural number '000000'`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `test{100000,0000000a}`,
		Error:   `     ^^^^^^^^^^^^^^^ syntax error: invalid natural number inside {...}`,
	})
	// ^ FIXME: Error:   `            ^^^^^^^^ syntax error: invalid natural number inside {...}`,
	//  testGenErr(t, &genErrCase{
	//		Pattern: `(a|)`,
	//		Error:   `  ^ '|' at the end of group`,
	//  })
	testGenErr(t, &genErrCase{
		Pattern: `$hex2dec(abcdefg)`,
		Error:   `               ^ value error: invalid hex number "abcdefg"`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$byte(12)`,
		Error:   `      ^ value error: function does not accept any arguments`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$BYTE(12)`,
		Error:   `      ^ value error: function does not accept any arguments`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$bip39word(abcd)`,
		Error:   `           ^^^^ value error: invalid number 'abcd'`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$bip39encode(gh)`,
		Error:   `             ^^ value error: invalid hex number "gh"`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$()`,
		Error:   ` ^ syntax error: missing function name`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$hex([a-z]`,
		Error:   `          ^ syntax error: '(' not closed`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$hex(([a-z]`,
		Error:   `           ^ syntax error: '(' not closed`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$hex([:x:])`,
		Error:   `      ^^^ value error: invalid character class "x"`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `[:hello:]`,
		Error:   ` ^^^^^^^ value error: invalid character class "hello"`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$hex([:hello:])`,
		Error:   `      ^^^^^^^ value error: invalid character class "hello"`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `(`,
		Error:   ` ^ syntax error: '(' not closed`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$foo`,
		Error:   `    ^ syntax error: expected a function call`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `($foo`,
		Error:   `     ^ syntax error: '(' not closed`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$foo(`,
		Error:   `     ^ syntax error: '(' not closed`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `test $foo(123)`,
		Error:   `     ^^^^^ value error: invalid function 'foo'`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$foo\()`,
		Error:   `    ^ syntax error: expected a function call`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `test($foo(123))`,
		Error:   `     ^^^^^ value error: invalid function 'foo'`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `test $foo`,
		Error:   `         ^ syntax error: expected a function call`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `test($foo)`,
		Error:   `         ^ syntax error: expected a function call`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$rjust(a,10000)[`,
		Error:   `                ^ syntax error: '[' not closed`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$rjust(a,10000)[a-]`,
		Error:   `                  ^ syntax error: no character after '-'`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `(((a{10,20})))[`,
		Error:   `               ^ syntax error: '[' not closed`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `(((a{10,20})))[a-]`,
		Error:   `                 ^ syntax error: no character after '-'`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `\u00e0-\u00e`,
		Error:   `       ^^^^^ syntax error: invalid escape sequence`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `\u00e0-\U00e6`,
		Error:   `       ^^^^^^ syntax error: invalid escape sequence`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `\u00mn`,
		Error:   `^^^^^^ syntax error: invalid escape sequence`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `test1 \u00mn test2`,
		Error:   `      ^^^^^^ syntax error: invalid escape sequence`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `(test1 \u00mn test2){2}`,
		Error:   `       ^^^^^^ syntax error: invalid escape sequence`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `test[\u00mn-\u00e0]abc`,
		Error:   `     ^^^^^^ syntax error: invalid escape sequence`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `test[\u00e0-\u00mn]abc`,
		Error:   `            ^^^^^^ syntax error: invalid escape sequence`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$date()`,
		Error:   `      ^ argument error: date: too few characters as arguments`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$date(2000)`,
		Error:   `          ^ argument error: date: at least 2 arguments are required`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$date(2000a,2000b)`,
		Error:   `      ^^^^^ value error: invalid year 2000a`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$date(2000,2000b)`,
		Error:   `           ^^^^^ value error: invalid year 2000b`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$date(2000,{{2000}})`,
		Error:   `nested '{'`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$rjust(abc)`,
		Error:   `          ^ argument error: rjust: at least 2 arguments are required`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$rjust(abc,aaaa)`,
		Error:   `           ^^^^ value error: invalid width aaaa`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$rjust(abc,0)`,
		Error:   `           ^ value error: invalid width 0`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$rjust(abc,-100)`,
		Error:   `           ^^^^ value error: invalid width -100`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$rjust(abc,1,ab)`,
		Error:   `             ^^ value error: invalid fillChar="ab", must have length 1`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$rjust({{}},7)`,
		Error:   fmt.Errorf(`nested '{'`),
	})
	testGenErr(t, &genErrCase{
		Pattern: `(abc) test1 \2 test2`,
		Error:   `            ^^ value error: invalid group id '2'`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `(abc) test1 \20 test2`,
		Error:   `            ^^^ value error: invalid group id '20'`,
	})
	testGenErr(t, &genErrCase{
		Pattern: `$pyhex(gh)`,
		Error:   `        ^ value error: invalid hex number "gh"`,
	})
	testGenErr(t, &genErrCase{
		Pattern: strings.Repeat("a", 1001),
		Error:   `pattern is too long`,
	})
}
