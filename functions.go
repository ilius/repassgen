package main

import (
	"encoding/base64"
	"strings"

	"github.com/ilius/crock32"
)

var functions = map[string]func(in []rune) ([]rune, error){
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
}
