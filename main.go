package main

import (
	"flag"
	"fmt"
	"io"
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
	Main(os.Stdout)
}

func Main(stdout io.Writer) {
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

	fmt.Fprintln(stdout, string(out.Password))
	if calcEnropy {
		if os.Getenv("REPASSGEN_FLOAT_ENTROPY") == "true" {
			fmt.Fprintf(
				stdout,
				"Entropy of pattern: %.2f bits\n",
				out.PatternEntropy,
			)
		} else {
			fmt.Fprintf(
				stdout,
				"Entropy of pattern: %d bits\n",
				int(out.PatternEntropy),
			)
		}
	}
}
