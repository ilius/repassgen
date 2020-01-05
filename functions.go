package main

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/ilius/crock32"
)

var encoderFunctions = map[string]func(in []rune) ([]rune, error){
	"base64": func(in []rune) ([]rune, error) {
		return []rune(base64.StdEncoding.EncodeToString([]byte(string(in)))), nil
	},
	"base64url": func(in []rune) ([]rune, error) {
		return []rune(base64.URLEncoding.EncodeToString([]byte(string(in)))), nil
	},

	// Crockford's Base32 encode functions (lowercase and uppercase)
	"base32": func(in []rune) ([]rune, error) {
		return []rune(strings.ToLower(crock32.Encode([]byte(string(in))))), nil
	},
	"BASE32": func(in []rune) ([]rune, error) {
		return []rune(crock32.Encode([]byte(string(in)))), nil
	},

	// standard Base32 encode function (uppercase, with no padding)
	"base32std": func(in []rune) ([]rune, error) {
		return []rune(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString([]byte(string(in)))), nil
	},

	// Hex encode functions (lowercase and uppercase)
	"hex": func(in []rune) ([]rune, error) {
		return []rune(hex.EncodeToString([]byte(string(in)))), nil
	},
	"HEX": func(in []rune) ([]rune, error) {
		return []rune(strings.ToUpper(hex.EncodeToString([]byte(string(in))))), nil
	},

	// Escape unicode characters, non-printable characters and double quote
	// The returned string uses Go escape sequences (\t, \n, \xFF, \u0100) for non-ASCII
	// characters and non-printable characters
	"escape": func(in []rune) ([]rune, error) {
		q := strconv.QuoteToASCII(string(in))
		return []rune(q[1 : len(q)-1]), nil
	},

	// BIP-39 encode function
	"bip39encode": bip39encode,
}
