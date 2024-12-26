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
	interpreter := create_interpreter(ast)
	global_env := create_enviorment(nil)

	for !interpreter.is_empty() {
		interpreter.evalute_current_statement(global_env)
		interpreter.advance(1)
	}

	litter.Dump(global_env.variables)
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
