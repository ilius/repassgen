package main

import (
	"strings"
	"testing"

	"github.com/ilius/is"
)

func strPtr(s string) *string {
	s2 := s
	return &s2
}

type genCase struct {
	Pattern string

	Error string

	PassLen [2]int     // {min, max}
	Entropy [2]float64 // {min, max}

	Password *string

	WordCount int

	// TODO: CharClassCount map[string]int
}

func TestGenerate(t *testing.T) {
	test := func(tc *genCase) {
		is := is.New(t).AddMsg("pattern=%#v", tc.Pattern)
		out, err := Generate(GenerateInput{Pattern: tc.Pattern})
		if tc.Error != "" {
			is.ErrMsg(err, tc.Error)
			is.Nil(out)
			return
		}
		is.NotErr(err)
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
		{
			entropy := out.PatternEntropy
			minEnt := tc.Entropy[0]
			maxEnt := tc.Entropy[1]
			is.AddMsg(
				"entropy=%v is not in range [%v, %v]",
				entropy,
				minEnt,
				maxEnt,
			).True(minEnt <= entropy && entropy <= maxEnt)
		}
		if tc.Password != nil {
			is.Equal(string(out.Password), *tc.Password)
		}
		if tc.WordCount != 0 {
			actual := len(strings.Split(string(out.Password), " "))
			is.Equal(actual, tc.WordCount)
		}
	}
	test(&genCase{
		Pattern: "",
		PassLen: [2]int{0, 0},
		Entropy: [2]float64{0, 0},
	})
	test(&genCase{
		Pattern: "[a-z]{8}",
		PassLen: [2]int{8, 8},
		Entropy: [2]float64{37.6, 37.7},
	})
	test(&genCase{
		Pattern: "[a-z]{8,10}",
		PassLen: [2]int{8, 10},
		Entropy: [2]float64{37.6, 47.01},
	})
	test(&genCase{
		Pattern: "[a-z]{8}[1-9]{3}",
		PassLen: [2]int{11, 11},
		Entropy: [2]float64{47.1, 47.2},
	})
	test(&genCase{
		Pattern: "([a-z]{8}[1-9]{3})",
		PassLen: [2]int{11, 11},
		Entropy: [2]float64{47.1, 47.2},
	})
	test(&genCase{
		Pattern: "([a-z]{5}[1-9]{2}){2}",
		PassLen: [2]int{14, 14},
		Entropy: [2]float64{59.6, 59.7},
	})
	test(&genCase{
		Pattern: "([a-z]{5}[1-9]{2}-){2}",
		PassLen: [2]int{16, 16},
		Entropy: [2]float64{59.6, 59.7},
	})
	// base64 length: ((bytes + 2) / 3) * 4
	test(&genCase{
		Pattern: "$base64([:byte:]{10})",
		PassLen: [2]int{16, 28},
		Entropy: [2]float64{80, 80},
	})
	test(&genCase{
		Pattern: "$base64([:byte:]{9})",
		PassLen: [2]int{12, 24},
		Entropy: [2]float64{72, 72},
	})
	test(&genCase{
		Pattern: "$base64([:byte:]{5})",
		PassLen: [2]int{8, 16},
		Entropy: [2]float64{40, 40},
	})
	test(&genCase{
		Pattern: "$base64url([:byte:]{5})",
		PassLen: [2]int{8, 16},
		Entropy: [2]float64{40, 40},
	})
	test(&genCase{
		Pattern: "$base32([:byte:]{5})",
		PassLen: [2]int{8, 16},
		Entropy: [2]float64{40, 40},
	})
	test(&genCase{
		Pattern: "$BASE32([:byte:]{5})",
		PassLen: [2]int{8, 16},
		Entropy: [2]float64{40, 40},
	})
	test(&genCase{
		Pattern: "$base32std([:byte:]{5})",
		PassLen: [2]int{8, 16},
		Entropy: [2]float64{40, 40},
	})
	test(&genCase{
		Pattern: "$hex([:byte:]{8})",
		PassLen: [2]int{16, 32},
		Entropy: [2]float64{64, 64},
	})
	test(&genCase{
		Pattern: "$HEX([:byte:]{8})",
		PassLen: [2]int{16, 30},
		Entropy: [2]float64{64, 64},
	})
	test(&genCase{
		Pattern:  `$escape(")`,
		PassLen:  [2]int{2, 2},
		Entropy:  [2]float64{0, 0},
		Password: strPtr(`\"`),
	})
	test(&genCase{
		Pattern:  `a[\t][\r][\n][\v][\f]b`,
		PassLen:  [2]int{7, 7},
		Entropy:  [2]float64{0, 0},
		Password: strPtr("a\t\r\n\v\fb"),
	})
	// each bip39 is at least 3 chars and max 8 chars
	test(&genCase{
		Pattern:   "$bip39word(10)",
		WordCount: 10,
		PassLen:   [2]int{39, 89}, // 10*4-1, 10*9+-1
		Entropy:   [2]float64{110, 110},
	})
	// 1 bip39 word   => 11 bits entropy
	// 8 bip39 words  => 11 bytes entropy
	test(&genCase{
		Pattern:   "$bip39encode([:alpha:]{11})",
		WordCount: 8,
		PassLen:   [2]int{43, 98}, // 11*4-1, 11*9+-1
		Entropy:   [2]float64{62.7, 62.8},
	})
}