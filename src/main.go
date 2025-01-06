package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sanity-io/litter"
	"github.com/table-harmony/HarmonyLang/src/interpreter"
	"github.com/table-harmony/HarmonyLang/src/lexer"
	"github.com/table-harmony/HarmonyLang/src/parser"
)

func main() {
	start := time.Now()
	run("src/struct-test.ham")
	//run_repl()
	duration := time.Since(start)

	fmt.Printf("Duration: %v\n", duration)
}

func run(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	source := string(bytes)
	tokens := lexer.Tokenize(source)
	ast := parser.Parse(tokens)
	litter.Dump(ast)
	interpreter.Interpret(ast)
}

func run_repl() {
	interpreter.StartREPL()
}
