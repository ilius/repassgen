package passgen

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/ilius/repassgen/lib/crock32"
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

var encoderFunctions = map[string]func(s *State, in []rune) ([]rune, error){
	"base64": func(s *State, in []rune) ([]rune, error) {
		data, err := hex.DecodeString(string(in))
		if err != nil {
			return nil, s.errorValue(s_invalid_hex_num, string(in))
		}
		return []rune(
			base64.StdEncoding.EncodeToString(data),
		), nil
	},
	"base64url": func(s *State, in []rune) ([]rune, error) {
		data, err := hex.DecodeString(string(in))
		if err != nil {
			return nil, s.errorValue(s_invalid_hex_num, string(in))
		}
		return []rune(
			base64.URLEncoding.EncodeToString(data),
		), nil
	},

	// Crockford's Base32 encode functions (lowercase and uppercase)
	"base32": func(s *State, in []rune) ([]rune, error) {
		data, err := hex.DecodeString(string(in))
		if err != nil {
			return nil, s.errorValue(s_invalid_hex_num, string(in))
		}
		return []rune(
			strings.ToLower(crock32.EncodeToString(data)),
		), nil
	},
	"BASE32": func(s *State, in []rune) ([]rune, error) {
		data, err := hex.DecodeString(string(in))
		if err != nil {
			return nil, s.errorValue(s_invalid_hex_num, string(in))
		}
		return []rune(
			crock32.EncodeToString(data),
		), nil
	},

	// standard Base32 encode function (uppercase, with no padding)
	"base32std": func(s *State, in []rune) ([]rune, error) {
		data, err := hex.DecodeString(string(in))
		if err != nil {
			return nil, s.errorValue(s_invalid_hex_num, string(in))
		}
		return []rune(
			base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(
				data,
			),
		), nil
	},

	// Hex encode functions (lowercase and uppercase)
	"hex": func(s *State, in []rune) ([]rune, error) {
		return []rune(
			hex.EncodeToString([]byte(string(in))),
		), nil
	},
	"HEX": func(s *State, in []rune) ([]rune, error) {
		return []rune(
			strings.ToUpper(hex.EncodeToString([]byte(string(in)))),
		), nil
	},

	"hex2dec": func(s *State, in []rune) ([]rune, error) {
		// if len(in) > 2 && in[0] == '0' && in[1] == 'x' {
		//	in = in[2:]
		//}
		i64, err := strconv.ParseInt(string(in), 16, 64)
		if err != nil {
			return nil, s.errorValue(s_invalid_hex_num, string(in))
		}
		return []rune(strconv.FormatInt(i64, 10)), nil
	},

	// pyhex converts hex-encoded bytes into a python bytes consisting hex values
	"pyhex": func(s *State, in []rune) ([]rune, error) {
		data, err := hex.DecodeString(string(in))
		if err != nil {
			return nil, s.errorValue(s_invalid_hex_num, string(in))
		}
		out := ""
		for _, cbyte := range data {
			chex := make([]byte, 2)
			n := hex.Encode(chex, []byte{cbyte})
			if n != 2 {
				return nil, s.errorUnknown("failed converting byte %x to hex", cbyte) // TODO: cover in test
			}
			out += "\\x" + string(chex)
		}
		return []rune("b'" + out + "'"), nil
	},

	"space": func(s *State, in []rune) ([]rune, error) {
		return expand1(' ', in), nil
	},
	"expand": func(s *State, in []rune) ([]rune, error) {
		if len(in) < 1 {
			return nil, nil
		}
		return expand1(in[0], in[1:]), nil
	},

	// Escape unicode characters, non-printable characters and double quote
	// The returned string uses Go escape sequences (\t, \n, \xFF, \u0100)
	// for non-ASCII characters and non-printable characters
	"escape": func(s *State, in []rune) ([]rune, error) {
		q := strconv.QuoteToASCII(string(in))
		return []rune(q[1 : len(q)-1]), nil
	},

	// Escape characters according to json
	// Escaped: ampersand, double quote
	// Not escaped: Unicode characters, newline, tab, backslash
	"json": func(s *State, in []rune) ([]rune, error) {
		outB, err := json.Marshal(string(in))
		if err != nil {
			return nil, s.errorUnknown("error in json Marshal: %v", err) // how to cover in test?
		}
		outS := string(outB)
		outS = outS[1 : len(outS)-1]
		return []rune(outS), nil
	},

	// BIP-39 encode function
	"bip39encode": bip39encode,

	// Japanese Kana to Latin
	"romaji": func(s *State, in []rune) ([]rune, error) {
		return []rune(KanaToRomaji(string(in))), nil
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
		s.errorMarkLen = len(funcName) + 2
		return s.errorValue("invalid function '%v'", funcName)
	}
	err := baseFunctionCallGenerator(
		s,
		NewState(s.SharedState, g.argPattern),
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
	return 0, s.errorUnknown(s_entropy_not_calc)
}

func getFuncGenerator(s *State, funcName string, arg []rune) (GeneratorIface, error) {
	if _, ok := encoderFunctions[funcName]; ok {
		return &encoderFunctionCallGenerator{
			funcName:   funcName,
			argPattern: arg,
		}, nil
	}
	switch funcName {
	case "byte":
		return newByteGenerator(s, arg, false)
	case "BYTE":
		return newByteGenerator(s, arg, true)
	case "bip39word":
		return newBIP39WordGenerator(s, string(arg))
	case "shuffle":
		return newShuffleGenerator(arg)
	case "date":
		return newDateGenerator(s, arg)
	case "?":
		return newOnceOrNoneGenerator(arg)
	case "rjust":
		return newRjustGenerator(s, arg)
	case "ljust":
		return newLjustGenerator(s, arg)
	case "center":
		return newCenterGenerator(s, arg)
	}
	s.errorMarkLen = len(funcName) + 2
	return nil, s.errorValue("invalid function '%v'", funcName)
}
