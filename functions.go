package main

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
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
		return []rune(
			base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(
				[]byte(string(in)),
			),
		), nil
	},

	// Hex encode functions (lowercase and uppercase)
	"hex": func(in []rune) ([]rune, error) {
		return []rune(hex.EncodeToString([]byte(string(in)))), nil
	},
	"HEX": func(in []rune) ([]rune, error) {
		return []rune(strings.ToUpper(
			hex.EncodeToString([]byte(string(in))),
		)), nil
	},

	// Escape unicode characters, non-printable characters and double quote
	// The returned string uses Go escape sequences (\t, \n, \xFF, \u0100)
	// for non-ASCII characters and non-printable characters
	"escape": func(in []rune) ([]rune, error) {
		q := strconv.QuoteToASCII(string(in))
		return []rune(q[1 : len(q)-1]), nil
	},

	// BIP-39 encode function
	"bip39encode": bip39encode,
}

type encoderFunctionCallGenerator struct {
	funcName   string
	argPattern string
	entropy    *float64
}

func (g *encoderFunctionCallGenerator) Generate(s *State) error {
	funcName := g.funcName
	funcObj, ok := encoderFunctions[funcName]
	if !ok {
		return s.errorValue("invalid function '%v'", funcName)
	}
	err := baseFunctionCallGenerator(
		s,
		NewState(s.SharedState, g.argPattern),
		funcName,
		funcObj,
	)
	if err != nil {
		return err
	}
	g.entropy = &s.patternEntropy
	return nil
}

func (g *encoderFunctionCallGenerator) Level() int {
	return 0
}

func (g *encoderFunctionCallGenerator) Entropy() (float64, error) {
	if g.entropy != nil {
		return *g.entropy, nil
	}
	return 0, fmt.Errorf("entropy is not calculated")
}

func getFuncGenerator(s *State, funcName string, arg string) (generatorIface, error) {
	if _, ok := encoderFunctions[funcName]; ok {
		return &encoderFunctionCallGenerator{
			funcName:   funcName,
			argPattern: arg,
		}, nil
	}
	switch funcName {
	case "bip39word":
		return newBIP99WordGenerator(arg)
	case "shuffle":
		return newShuffleGenerator(arg)
	}
	return nil, s.errorValue("invalid function '%v'", funcName)
}
