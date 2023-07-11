package main

import (
	"flag"
	"fmt"
	"os"

	passgen "github.com/ilius/repassgen/lib"
)

var entropyFlag = flag.Bool(
	"entropy",
	false,
	"repassgen [-entropy] PATTERN",
)

func printError(s *passgen.State, err error, pattern string) {
	myErr, ok := err.(*passgen.Error)
	if !ok {
		fmt.Println(err)
		return
	}
	fmt.Println(string(pattern))
	fmt.Println(myErr.SpacedError())
}

func main() {
	flag.Parse()

	calcEnropy := entropyFlag != nil && *entropyFlag

	pattern := flag.Arg(0)
	out, s, err := passgen.Generate(passgen.GenerateInput{
		Pattern: []rune(pattern),
	})
	if err != nil {
		printError(s, err, pattern)
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
