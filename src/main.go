package main

import (
	"os"

	"github.com/table-harmony/HarmonyLang/src/lexer"
)

func main() {
	bytes, err := os.ReadFile("examples/01.harmony")

	if err != nil {
		panic(err)
	}

	source := string(bytes)

	var tokens = lexer.Tokenize(source)

	for _, token := range tokens {
		token.Print()
	}
}
