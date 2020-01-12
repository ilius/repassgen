package main

import (
	"testing"

	"github.com/ilius/is"
)

type genCase struct {
	Pattern string

	Error string

	PassLen [2]int     // {min, max}
	Entropy [2]float64 // {min, max}

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
}
