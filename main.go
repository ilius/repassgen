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

func printError(err error, pattern string) {
	myErr, ok := err.(*passgen.Error)
	if !ok {
		fmt.Println(err)
		return
	}
	fmt.Println(pattern)
	fmt.Println(myErr.SpacedError())
}

func main() {
	Main(os.Stdout)
}

func Main(stdout io.Writer) {
	flag.Parse()

	calcEnropy := entropyFlag != nil && *entropyFlag

	pattern := flag.Arg(0)
	out, _, err := passgen.Generate(passgen.GenerateInput{
		Pattern: []rune(pattern),
	})
	if err != nil {
		printError(err, pattern)
		os.Exit(1)
	}

	_, err = fmt.Fprintln(stdout, string(out.Password))
	if err != nil {
		panic(err)
	}
	if calcEnropy {
		if os.Getenv("REPASSGEN_FLOAT_ENTROPY") == "true" {
			_, err := fmt.Fprintf(
				stdout,
				"Entropy of pattern: %.2f bits\n",
				out.PatternEntropy,
			)
			if err != nil {
				panic(err)
			}
		} else {
			_, err := fmt.Fprintf(
				stdout,
				"Entropy of pattern: %d bits\n",
				int(out.PatternEntropy),
			)
			if err != nil {
				panic(err)
			}
		}
	}
}
