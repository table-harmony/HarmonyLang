package interpreter

import (
	"github.com/sanity-io/litter"
	"github.com/table-harmony/HarmonyLang/src/ast"
)

type interpreter struct {
	ast []ast.Statement
	pos int
}

func Interpret(ast []ast.Statement) {
	interpreter := createInterpreter(ast)

	for !interpreter.isEmpty() {
		litter.Dump(interpreter.currentStatement())
		evalute_statement(interpreter)
		interpreter.advance(1)
	}
}

func createInterpreter(ast []ast.Statement) *interpreter {
	create_lookups()

	return &interpreter{
		ast: ast,
		pos: 0,
	}
}

func (interpreter *interpreter) advance(n int) {
	interpreter.pos += n
}

func (interpreter *interpreter) currentStatement() ast.Statement {
	return interpreter.ast[interpreter.pos]
}

func (interpreter *interpreter) isEmpty() bool {
	return interpreter.pos >= len(interpreter.ast)
}
