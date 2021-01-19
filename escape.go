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
	"fmt"
	"unicode/utf8"
)

// JSON Unicode stuff: see https://tools.ietf.org/html/rfc7159#section-7

const supplementalPlanesOffset = 0x10000
const highSurrogateOffset = 0xD800
const lowSurrogateOffset = 0xDC00

const basicMultilingualPlaneReservedOffset = 0xDFFF
const basicMultilingualPlaneOffset = 0xFFFF

func combineUTF16Surrogates(high, low rune) rune {
	return supplementalPlanesOffset + (high-highSurrogateOffset)<<10 + (low - lowSurrogateOffset)
}

const badHex = -1

func h2I(c rune) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'A' && c <= 'F':
		return int(c - 'A' + 10)
	case c >= 'a' && c <= 'f':
		return int(c - 'a' + 10)
	}
	return badHex
}

// decodeSingleUnicodeEscape decodes a single \uXXXX escape sequence. The prefix \u is assumed to be present and
// is not checked.
// In JSON, these escapes can either come alone or as part of "UTF16 surrogate pairs" that must be handled together.
// This function only handles one; decodeUnicodeEscape handles this more complex case.
func decodeSingleUnicodeEscape(in []rune) (rune, bool) {
	// We need at least 6 characters total
	if len(in) < 6 {
		return utf8.RuneError, false
	}

	// Convert hex to decimal
	h1, h2, h3, h4 := h2I(in[2]), h2I(in[3]), h2I(in[4]), h2I(in[5])
	if h1 == badHex || h2 == badHex || h3 == badHex || h4 == badHex {
		return utf8.RuneError, false
	}

	// Compose the hex digits
	return rune(h1<<12 + h2<<8 + h3<<4 + h4), true
}

// isUTF16EncodedRune checks if a rune is in the range for non-BMP characters,
// which is used to describe UTF16 chars.
// Source: https://en.wikipedia.org/wiki/Plane_(Unicode)#Basic_Multilingual_Plane
func isUTF16EncodedRune(r rune) bool {
	return highSurrogateOffset <= r && r <= basicMultilingualPlaneReservedOffset
}

func decodeUnicodeEscape(in []rune) (rune, int) {
	if r, ok := decodeSingleUnicodeEscape(in); !ok {
		// Invalid Unicode escape
		return utf8.RuneError, -1
	} else if r <= basicMultilingualPlaneOffset && !isUTF16EncodedRune(r) {
		// Valid Unicode escape in Basic Multilingual Plane
		return r, 6
	} else if r2, ok := decodeSingleUnicodeEscape(in[6:]); !ok { // Note: previous decodeSingleUnicodeEscape success guarantees at least 6 bytes remain
		// UTF16 "high surrogate" without manditory valid following Unicode escape for the "low surrogate"
		return utf8.RuneError, -1
	} else if r2 < lowSurrogateOffset {
		// Invalid UTF16 "low surrogate"
		return utf8.RuneError, -1
	} else {
		// Valid UTF16 surrogate pair
		return combineUTF16Surrogates(r, r2), 12
	}
}

// backslashCharEscapeTable: when '\X' is found for some byte X, it is to be replaced with backslashCharEscapeTable[X]
var backslashCharEscapeTable = [...]rune{
	'"':  '"',
	'\\': '\\',
	'/':  '/',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
}

// unescapeUnicodeSingle unescapes the single escape sequence starting at 'in' into 'out' and returns
// how many characters were consumed from 'in' and emitted into 'out'.
// If a valid escape sequence does not appear as a prefix of 'in', (-1, -1) to signal the error.
func unescapeUnicodeSingle(in []rune, start int) (int, rune, error) {
	if len(in)-start < 2 || in[start] != '\\' {
		// Invalid escape due to insufficient characters for any escape or no initial backslash
		return -1, 0, nil
	}

	// https://tools.ietf.org/html/rfc7159#section-7
	switch e := in[start+1]; e {
	case 'u':
		// Unicode escape
		if r, inLen := decodeUnicodeEscape(in[start:]); inLen == -1 {
			// Invalid Unicode escape
			return -1, 0, fmt.Errorf("invalid escape sequence near index %d", start)
		} else {
			return inLen, r, nil
		}
	}

	return -1, 0, nil
}
