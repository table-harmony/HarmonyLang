package interpreter

import (
	"fmt"

	"github.com/table-harmony/HarmonyLang/src/ast"
)

type interpreter struct {
	ast []ast.Statement
	pos int
}

func Interpret(ast []ast.Statement) {
	interpreter := create_interpreter(ast)
	scope := NewScope(nil)

	for !interpreter.is_empty() {
		interpreter.evalute_current_statement(scope)
		interpreter.advance(1)
	}

	fmt.Println(scope.String())
}

func create_interpreter(ast []ast.Statement) *interpreter {
	create_lookups()

	return &interpreter{
		ast: ast,
		pos: 0,
	}
}

func (interpreter *interpreter) advance(n int) {
	interpreter.pos += n
}

func (interpreter *interpreter) current_statement() ast.Statement {
	return interpreter.ast[interpreter.pos]
}

func (interpreter *interpreter) is_empty() bool {
	return interpreter.pos >= len(interpreter.ast)
}
