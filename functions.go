package main

import "encoding/base64"

var functions = map[string]func(in []rune) ([]rune, error){
	"base64": func(in []rune) ([]rune, error) {
		return []rune(base64.StdEncoding.EncodeToString([]byte(string(in)))), nil
	},
	"base64url": func(in []rune) ([]rune, error) {
		return []rune(base64.URLEncoding.EncodeToString([]byte(string(in)))), nil
	},
}
