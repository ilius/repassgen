package main

import (
	"flag"
	"fmt"
	"os"
)

var entropyFlag = flag.Bool(
	"entropy",
	false,
	"repassgen [-entropy] PATTERN",
)

func main() {
	InitKana()

	flag.Parse()

	calcEnropy := entropyFlag != nil && *entropyFlag

	out, err := Generate(GenerateInput{
		Pattern: []rune(flag.Arg(0)),
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
