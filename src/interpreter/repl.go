package interpreter

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/core"
	"github.com/table-harmony/HarmonyLang/src/lexer"
	"github.com/table-harmony/HarmonyLang/src/parser"
)

type REPL struct {
	scope  *core.Scope
	reader *bufio.Reader
}

func StartREPL() {
	repl := create_repl()

	fmt.Println("Harmony Lang REPL")
	fmt.Println("Type 'exit' to quit")

	for {
		fmt.Print(">> ")
		input, err := repl.reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "exit" {
			break
		}

		if input == "" {
			continue
		}

		result := repl.evaluate(input)
		if result != nil {
			print_value(result)
		}
	}
}

func create_repl() REPL {
	create_lookups()

	return REPL{
		scope:  core.NewScope(nil),
		reader: bufio.NewReader(os.Stdin),
	}
}

func (repl *REPL) evaluate(input string) core.Value {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error: %v\n", r)
		}
	}()

	tokens := lexer.Tokenize(input)
	ast := parser.Parse(tokens)

	var lastResult core.Value
	for _, statement := range ast {
		lastResult = repl.evaluate_statement(statement)
	}

	return lastResult
}

func (repl *REPL) evaluate_statement(statement ast.Statement) core.Value {
	switch statement := statement.(type) {
	case ast.ExpressionStatement:
		return evaluate_expression(statement.Expression, repl.scope)
	default:
		evaluate_statement(statement, repl.scope)
		return nil
	}
}

func print_value(value core.Value) {
	switch v := value.(type) {
	case core.Number:
		if v.Value() == float64(int(v.Value())) {
			fmt.Printf("%d\n", int(v.Value()))
		} else {
			fmt.Printf("%g\n", v.Value())
		}
	case core.String:
		fmt.Printf("%q\n", v.Value())
	case core.Boolean:
		fmt.Printf("%t\n", v.Value())
	case core.Reference:
		print_value(v.Load())
	default:
		fmt.Printf("%v\n", value)
	}
}
