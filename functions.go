package main

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/ilius/crock32"
)

func expand1(sep rune, in []rune) []rune {
	if len(in) < 1 {
		return nil
	}
	out := make([]rune, 2*len(in)-1)
	out[0] = in[0]
	for i, c := range in[1:] {
		out[2*i+1] = sep
		out[2*i+2] = c
	}
	return out
}

var encoderFunctions = map[string]func(in []rune) []rune{
	"base64": func(in []rune) []rune {
		return []rune(base64.StdEncoding.EncodeToString([]byte(string(in))))
	},
	"base64url": func(in []rune) []rune {
		return []rune(base64.URLEncoding.EncodeToString([]byte(string(in))))
	},

	// Crockford's Base32 encode functions (lowercase and uppercase)
	"base32": func(in []rune) []rune {
		return []rune(strings.ToLower(crock32.Encode([]byte(string(in)))))
	},
	"BASE32": func(in []rune) []rune {
		return []rune(crock32.Encode([]byte(string(in))))
	},

	// standard Base32 encode function (uppercase, with no padding)
	"base32std": func(in []rune) []rune {
		return []rune(
			base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(
				[]byte(string(in)),
			),
		)
	},

	// Hex encode functions (lowercase and uppercase)
	"hex": func(in []rune) []rune {
		return []rune(hex.EncodeToString([]byte(string(in))))
	},
	"HEX": func(in []rune) []rune {
		return []rune(strings.ToUpper(
			hex.EncodeToString([]byte(string(in))),
		))
	},

	"space": func(in []rune) []rune {
		return expand1(' ', in)
	},
	"expand": func(in []rune) []rune {
		if len(in) < 1 {
			return nil
		}
		return expand1(in[0], in[1:])
	},

	// Escape unicode characters, non-printable characters and double quote
	// The returned string uses Go escape sequences (\t, \n, \xFF, \u0100)
	// for non-ASCII characters and non-printable characters
	"escape": func(in []rune) []rune {
		q := strconv.QuoteToASCII(string(in))
		return []rune(q[1 : len(q)-1])
	},

	// BIP-39 encode function
	"bip39encode": bip39encode,

	// Japanese Kana to Latin
	"romaji": func(in []rune) []rune {
		return []rune(KanaToRomaji(string(in)))
	},
}

type encoderFunctionCallGenerator struct {
	entropy    *float64
	funcName   string
	argPattern []rune
}

func (g *encoderFunctionCallGenerator) Generate(s *State) error {
	funcName := g.funcName
	funcObj, ok := encoderFunctions[funcName]
	if !ok {
		return s.errorValue("invalid function '%v'", funcName)
	}
	err := baseFunctionCallGenerator(
		s,
		NewState(s.SharedState, []rune(g.argPattern)),
		funcName,
		funcObj,
	)
	if err != nil {
		return err
	}
	g.entropy = &s.patternEntropy
	return nil
}

func (g *encoderFunctionCallGenerator) Entropy(s *State) (float64, error) {
	if g.entropy != nil {
		return *g.entropy, nil
	}
	return 0, s.errorUnknown("entropy is not calculated")
}

func getFuncGenerator(s *State, funcName string, arg []rune) (generatorIface, error) {
	if _, ok := encoderFunctions[funcName]; ok {
		return &encoderFunctionCallGenerator{
			funcName:   funcName,
			argPattern: arg,
		}, nil
	}
	switch funcName {
	case "bip39word":
		return newBIP39WordGenerator(s, string(arg))
	case "shuffle":
		return newShuffleGenerator(arg)
	case "date":
		return newDateGenerator(s, string(arg))
	case "?":
		return newOnceOrNoneGenerator(arg)
	}
	return nil, s.errorValue("invalid function '%v'", funcName)
}
