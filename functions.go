package main

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/ilius/crock32"
)

var encoderFunctions = map[string]func(in []byte) []byte{
	"base64": func(in []byte) []byte {
		return []byte(base64.StdEncoding.EncodeToString(in))
		//out := make([]byte, base64.StdEncoding.EncodedLen(len(in)))
		//base64.StdEncoding.Encode(out, in)
		//return out
	},
	"base64url": func(in []byte) []byte {
		return []byte(base64.URLEncoding.EncodeToString(in))
		//out := make([]byte, base64.URLEncoding.EncodedLen(len(in)))
		//base64.URLEncoding.Encode(out, in)
		//return out
	},

	// Crockford's Base32 encode functions (lowercase and uppercase)
	"base32": func(in []byte) []byte {
		return bytes.ToLower([]byte(crock32.Encode(in)))
	},
	"BASE32": func(in []byte) []byte {
		return []byte(crock32.Encode(in))
	},

	// standard Base32 encode function (uppercase, with no padding)
	"base32std": func(in []byte) []byte {
		// FIXME: directly to []byte
		return []byte(
			base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(
				in,
			),
		)
	},

	// Hex encode functions (lowercase and uppercase)
	"hex": func(in []byte) []byte {
		// FIXME: directly to []byte
		return []byte(hex.EncodeToString(in))
	},
	"HEX": func(in []byte) []byte {
		// FIXME: directly to []byte
		return []byte(strings.ToUpper(
			hex.EncodeToString(in),
		))
	},

	// Escape unicode characters, non-printable characters and double quote
	// The returned string uses Go escape sequences (\t, \n, \xFF, \u0100)
	// for non-ASCII characters and non-printable characters
	"escape": func(in []byte) []byte {
		q := strconv.QuoteToASCII(string(in))
		return []byte(q[1 : len(q)-1])
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
		return newBIP39WordGenerator(arg)
	case "shuffle":
		return newShuffleGenerator(arg)
	}
	return nil, s.errorValue("invalid function '%v'", funcName)
}
