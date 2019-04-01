package main

import (
	"fmt"
	"os"
)

func main() {
	pattern := os.Args[1]
	fmt.Println(string(Generate(pattern)))
}
