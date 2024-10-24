package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	passgen "github.com/ilius/repassgen/lib"
	"github.com/ilius/repassgen/xflag"
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
	Main(os.Stdout, os.Args)
}

func Main(stdout io.Writer, args []string) {
	flagSet := &flag.FlagSet{}

	entropyFlag := flagSet.Bool(
		"entropy",
		false,
		"repassgen [-entropy] PATTERN",
	)

	err := xflag.ParseToEnd(flagSet, args)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	calcEnropy := entropyFlag != nil && *entropyFlag

	pattern := flagSet.Arg(1)
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
