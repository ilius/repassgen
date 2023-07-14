package passgen_test

import (
	"testing"
	"time"

	passgen "github.com/ilius/repassgen/lib"
)

func FuzzGenerate(f *testing.F) {
	testcases := []string{
		"Hello, world",
		" ",
		"!1234567890",
		`Hello, world 1234567890!@#$%^&*()_+{}[];\,./`,
		time.Now().Format(time.RFC3339Nano),
	}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, pattern string) {
		// TODO: define a max output length to prevent out of memory and crash
		out, _, err := passgen.Generate(passgen.GenerateInput{
			Pattern: []rune(pattern),
		})
		if err != nil {
			return
		}
		if out == nil {
			panic("out == nil")
		}
	})
}
