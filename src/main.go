package main

import (
	"os"

	"github.com/table-harmony/HarmonyLang/src/interpreter"
	"github.com/table-harmony/HarmonyLang/src/lexer"
	"github.com/table-harmony/HarmonyLang/src/parser"
)

func main() {
	source := read_file("examples/01.ham")

	tokens := lexer.Tokenize(source)
	ast := parser.Parse(tokens)
	interpreter.Interpret(ast)
}

func read_file(path string) string {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	source := string(bytes)
	return source
}
