package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	pattern := os.Args[1]
	s := NewState(pattern)
	lex := LexRoot
	var err error
	for {
		lex, err = lex(s)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(fmt.Errorf("invalid syntax near character %d: %v", s.patternPos, err))
		}
	}
	fmt.Println(string(s.output))
}
