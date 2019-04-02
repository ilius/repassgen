package main

import (
	"flag"
	"fmt"
)

func main() {
	entropyFlag := flag.Bool(
		"entropy",
		false,
		"repassgen [-entropy] PATTERN",
	)

	flag.Parse()

	calcEnropy := entropyFlag != nil && *entropyFlag

	pattern := flag.Arg(0)
	out := Generate(GenerateInput{
		Pattern:            pattern,
		CalcPatternEntropy: calcEnropy,
	})

	fmt.Println(string(out.Password))
	if calcEnropy {
		fmt.Printf(
			"Entropy of pattern: %d bits\n",
			int(out.PatternEntropy),
		)
	}
}
