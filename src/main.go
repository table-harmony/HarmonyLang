package main

import (
	"fmt"
	"os"

	"github.com/table-harmony/HarmonyLang/src/lexer"
)

func main() {
	bytes, err := os.ReadFile("examples/00.harmony")

	if err != nil {
		panic(err)
	}

	source := string(bytes)

	var tokens = lexer.Tokenize(source)

	for _, token := range tokens {
		token.Print()
	}

	fmt.Printf("Code: %s\n", source)
}
