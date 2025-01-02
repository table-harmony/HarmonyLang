package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sanity-io/litter"
	"github.com/table-harmony/HarmonyLang/src/lexer"
	"github.com/table-harmony/HarmonyLang/src/parser"
)

func main() {
	start := time.Now()
	run("examples/01.ham")
	duration := time.Since(start)

	fmt.Printf("Duration: %v\n", duration)
}

func run(path string) {
	source := read_file(path)

	tokens := lexer.Tokenize(source)
	ast := parser.Parse(tokens)
	litter.Dump(ast)
	//interpreter.Interpret(ast)
}

func read_file(path string) string {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	source := string(bytes)
	return source
}
