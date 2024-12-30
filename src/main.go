package main

import (
	"os"

	"github.com/table-harmony/HarmonyLang/src/interpreter"
	"github.com/table-harmony/HarmonyLang/src/lexer"
	"github.com/table-harmony/HarmonyLang/src/parser"
)

func main() {
	bytes, err := os.ReadFile("examples/01.ham")

	if err != nil {
		panic(err)
	}

	source := string(bytes)

	tokens := lexer.Tokenize(source)
	ast := parser.Parse(tokens)
	interpreter.Interpret(ast)
}
