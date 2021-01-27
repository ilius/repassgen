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

func printError(s *State, err error) {
	myErr, ok := err.(*Error)
	if !ok {
		fmt.Println(err)
		return
	}
	fmt.Println(string(s.pattern))
	fmt.Println(myErr.SpacedError())
}

func main() {
	flag.Parse()

	calcEnropy := entropyFlag != nil && *entropyFlag

	out, s, err := Generate(GenerateInput{
		Pattern: []rune(flag.Arg(0)),
	})
	if err != nil {
		printError(s, err)
		os.Exit(1)
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
