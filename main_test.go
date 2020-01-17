package main

import (
	"os"
	"testing"
)

func TestMainFunc(t *testing.T) {
	args := os.Args
	defer func() {
		os.Args = args
	}()

	os.Setenv("REPASSGEN_FLOAT_ENTROPY", "")

	os.Args = []string{"repassgen", "[a-z]{6}"}
	main()

	os.Args = []string{"repassgen", "-entropy", "[a-z]{6}"}
	main()

	os.Setenv("REPASSGEN_FLOAT_ENTROPY", "true")

	os.Args = []string{"repassgen", "-entropy", "[a-z]{6}"}
	main()
}
