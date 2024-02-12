package main

import (
	"bytes"
	"os"
	"testing"
)

func TestMainFunc(t *testing.T) {
	stdout := bytes.NewBuffer(nil)
	args := os.Args
	defer func() {
		os.Args = args
	}()

	err := os.Setenv("REPASSGEN_FLOAT_ENTROPY", "")
	if err != nil {
		panic(err)
	}

	os.Args = []string{"repassgen", "[a-z]{6}"}
	Main(stdout)

	os.Args = []string{"repassgen", "-entropy", "[a-z]{6}"}
	Main(stdout)

	err = os.Setenv("REPASSGEN_FLOAT_ENTROPY", "true")
	if err != nil {
		panic(err)
	}

	os.Args = []string{"repassgen", "-entropy", "[a-z]{6}"}
	Main(stdout)
}
