package main

import (
	"os"

	"github.com/sanity-io/litter"

	"github.com/table-harmony/HarmonyLang/src/lexer"
	"github.com/table-harmony/HarmonyLang/src/parser"
)

func main() {
	bytes, err := os.ReadFile("examples/01.harmony")

	if err != nil {
		panic(err)
	}

	source := string(bytes)

	tokens := lexer.Tokenize(source)
	ast := parser.Parse(tokens)

	litter.Dump(ast)
}
