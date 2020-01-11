package main

import (
	"flag"
	"fmt"
	"os"
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
	out, err := Generate(GenerateInput{
		Pattern: pattern,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out.Password))
	if calcEnropy {
		if os.Getenv("REPASSGEN_FLOAT_ENTROPY") == "true" {
			fmt.Printf(
				"Entropy of pattern: %.2f bits\n",
				out.PatternEntropy,
			)
		} else {
			fmt.Printf(
				"Entropy of pattern: %d bits\n",
				int(out.PatternEntropy),
			)
		}
	}
}
