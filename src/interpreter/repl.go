package interpreter

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/lexer"
	"github.com/table-harmony/HarmonyLang/src/parser"
)

type REPL struct {
	environment *Environment
	reader      *bufio.Reader
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
		environment: create_environment(nil),
		reader:      bufio.NewReader(os.Stdin),
	}
}

func (repl *REPL) evaluate(input string) RuntimeValue {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error: %v\n", r)
		}
	}()

	tokens := lexer.Tokenize(input)
	ast := parser.Parse(tokens)

	var lastResult RuntimeValue
	for _, statement := range ast {
		lastResult = repl.evaluate_statement(statement)
	}

	return lastResult
}

func (repl *REPL) evaluate_statement(statement ast.Statement) RuntimeValue {
	switch statement := statement.(type) {
	case ast.ExpressionStatement:
		return evaluate_expression(statement.Expression, repl.environment)
	default:
		evaluate_statement(statement, repl.environment)
		return nil
	}
}

func print_value(value RuntimeValue) {
	switch v := value.(type) {
	case RuntimeNumber:
		if v.Value == float64(int(v.Value)) {
			fmt.Printf("%d\n", int(v.Value))
		} else {
			fmt.Printf("%g\n", v.Value)
		}
	case RuntimeString:
		fmt.Printf("%q\n", v.Value)
	case RuntimeBoolean:
		fmt.Printf("%t\n", v.Value)
	case RuntimeVariable:
		print_value(v.Value)
	default:
		fmt.Printf("%v\n", value)
	}
}
