// MIT License

// Copyright (c) 2016 Leonid Bugaev
// Copyright (c) 2020 Saeed Rasooli

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"strings"
	"testing"

	"github.com/ilius/is/v2"
)

// unescape unescapes the string contained in 'in' and returns it as a slice.
func unescapeUnicode(in []rune) ([]rune, error) {
	firstBackslash := strings.IndexRune(string(in), '\\')
	if firstBackslash == -1 {
		return in, nil
	}

	out := make([]rune, 0, len(in))

	// Copy the first sequence of unescaped bytes to the output and obtain a buffer pointer (subslice)
	out = in[:firstBackslash]
	start := firstBackslash

	for start < len(in) {
		// Unescape the next escaped character
		inLen, r, err := unescapeUnicodeSingle(in, start)
		if err != nil {
			return nil, err
		}
		if inLen == -1 {
			out = append(out, in[start])
			start += 1
			continue
		}
		out = append(out, r)

		start += inLen

		// Copy everything up until the next backslash
		nextBackslash := strings.IndexRune(string(in[start:]), '\\')
		if nextBackslash == -1 {
			out = append(out, in[start:]...)
			break
		} else {
			out = append(out, in[start:start+nextBackslash]...)
			start += nextBackslash
		}
	}

	return out, nil
}

func TestH2I(t *testing.T) {
	hexChars := []rune{'0', '9', 'A', 'F', 'a', 'f', 'x', '\000'}
	hexValues := []int{0, 9, 10, 15, 10, 15, -1, -1}

	for i, c := range hexChars {
		if v := h2I(c); v != hexValues[i] {
			t.Errorf("h2I('%c') returned wrong value (obtained %d, expected %d)", c, v, hexValues[i])
		}
	}
}

type escapedUnicodeRuneTest struct {
	in    string
	isErr bool
	out   rune
	len   int
}

var commonUnicodeEscapeTests = []escapedUnicodeRuneTest{
	{in: `\u0041`, out: 'A', len: 6},
	{in: `\u0000`, out: 0, len: 6},
	{in: `\u00b0`, out: '°', len: 6},
	{in: `\u00B0`, out: '°', len: 6},

	{in: `\x1234`, out: 0x1234, len: 6}, // These functions do not check the \u prefix

	{in: ``, isErr: true},
	{in: `\`, isErr: true},
	{in: `\u`, isErr: true},
	{in: `\u1`, isErr: true},
	{in: `\u11`, isErr: true},
	{in: `\u111`, isErr: true},
	{in: `\u123X`, isErr: true},
}

var singleUnicodeEscapeTests = append([]escapedUnicodeRuneTest{
	{in: `\uD83D`, out: 0xD83D, len: 6},
	{in: `\uDE03`, out: 0xDE03, len: 6},
	{in: `\uFFFF`, out: 0xFFFF, len: 6},
	{in: `\uFF11`, out: '１', len: 6},
}, commonUnicodeEscapeTests...)

var multiUnicodeEscapeTests = append([]escapedUnicodeRuneTest{
	{in: `\uD83D`, isErr: true},
	{in: `\uDE03`, isErr: true},
	{in: `\uFFFF`, out: '\uFFFF', len: 6},
	{in: `\uFF11`, out: '１', len: 6},

	{in: `\uD83D\uDE03`, out: '\U0001F603', len: 12},
	{in: `\uD800\uDC00`, out: '\U00010000', len: 12},

	{in: `\uD800\`, isErr: true},
	{in: `\uD800\u`, isErr: true},
	{in: `\uD800\uD`, isErr: true},
	{in: `\uD800\uDC`, isErr: true},
	{in: `\uD800\uDC0`, isErr: true},
	{in: `\uD800\uDBFF`, isErr: true}, // invalid low surrogate
}, commonUnicodeEscapeTests...)

func TestDecodeSingleUnicodeEscape(t *testing.T) {
	for _, test := range singleUnicodeEscapeTests {
		r, ok := decodeSingleUnicodeEscape([]rune(test.in))
		isErr := !ok

		if isErr != test.isErr {
			t.Errorf("decodeSingleUnicodeEscape(%s) returned isErr mismatch: expected %t, obtained %t", test.in, test.isErr, isErr)
		} else if isErr {
			continue
		} else if r != test.out {
			t.Errorf("decodeSingleUnicodeEscape(%s) returned rune mismatch: expected %x (%c), obtained %x (%c)", test.in, test.out, test.out, r, r)
		}
	}
}

func TestDecodeUnicodeEscape(t *testing.T) {
	for _, test := range multiUnicodeEscapeTests {
		r, len := decodeUnicodeEscape([]rune(test.in))
		isErr := (len == -1)

		if isErr != test.isErr {
			t.Errorf("decodeUnicodeEscape(%s) returned isErr mismatch: expected %t, obtained %t", test.in, test.isErr, isErr)
		} else if isErr {
			continue
		} else if len != test.len {
			t.Errorf("decodeUnicodeEscape(%s) returned length mismatch: expected %d, obtained %d", test.in, test.len, len)
		} else if r != test.out {
			t.Errorf("decodeUnicodeEscape(%s) returned rune mismatch: expected %x (%c), obtained %x (%c)", test.in, test.out, test.out, r, r)
		}
	}
}

type unescapeTest struct {
	in    string // escaped string
	out   string // expected unescaped string
	isErr bool   // should this operation result in an error
}

var unescapeTests = []unescapeTest{
	{in: ``, out: ``},
	{in: `a`, out: `a`},
	{in: `abcde`, out: `abcde`},

	{in: `ab\\de`, out: `ab\\de`},
	{in: `ab\"de`, out: `ab\"de`},
	{in: `ab \u00B0 de`, out: `ab ° de`},
	{in: `ab \uFF11 de`, out: `ab １ de`},
	{in: `\uFFFF`, out: "\uFFFF"},
	{in: `ab \uD83D\uDE03 de`, out: "ab \U0001F603 de"},
	{in: `\u0000\u0000\u0000\u0000\u0000`, out: "\u0000\u0000\u0000\u0000\u0000"},
	{in: `\u0000 \u0000 \u0000 \u0000 \u0000`, out: "\u0000 \u0000 \u0000 \u0000 \u0000"},
	{in: ` \u0000 \u0000 \u0000 \u0000 \u0000 `, out: " \u0000 \u0000 \u0000 \u0000 \u0000 "},

	{in: `\uD800`, isErr: true},
	{in: `abcde\`, out: `abcde\`},
	{in: `abcde\x`, out: `abcde\x`},
	{in: `abcde\u`, isErr: true},
	{in: `abcde\u1`, isErr: true},
	{in: `abcde\u12`, isErr: true},
	{in: `abcde\u123`, isErr: true},
	{in: `abcde\uD800`, isErr: true},
	{in: `ab\uD800de`, isErr: true},
	{in: `\uD800abcde`, isErr: true},
}

// isSameMemory checks if two slices contain the same memory pointer (meaning one is a
// subslice of the other, with possibly differing lengths/capacities).
func isSameMemory(a, b []rune) bool {
	if cap(a) == 0 || cap(b) == 0 {
		return cap(a) == cap(b)
	} else if a, b = a[:1], b[:1]; a[0] != b[0] {
		return false
	} else {
		a[0]++
		same := (a[0] == b[0])
		a[0]--
		return same
	}

}

func TestUnescapeUnicode(t *testing.T) {
	for _, test := range unescapeTests {
		is := is.New(t).AddMsg("input=%#v", test.in)
		in := []rune(test.in)

		out, err := unescapeUnicode(in)
		if test.isErr {
			is.Err(err)
			continue
		}
		if !is.NotErr(err) {
			continue
		}
		is.Equal(string(out), string(test.out))
	}
}
