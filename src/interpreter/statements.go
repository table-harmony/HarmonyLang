package interpreter

import (
	"fmt"
	"reflect"

	"github.com/sanity-io/litter"
	"github.com/table-harmony/HarmonyLang/src/ast"
)

func evalute_statement(interpreter *interpreter) {
	statement := interpreter.currentStatement()
	statement_type := reflect.TypeOf(statement)

	if handler, exists := statement_lookup[statement_type]; exists {
		handler(interpreter)
	} else {
		panic(fmt.Sprintf("No handler registered for statement type: %v", statement_type))
	}
}

func evaluate_expression_statement(interpreter *interpreter) RuntimeValue {
	expression_statement := interpreter.currentStatement().(ast.ExpressionStatement)

	var result = evaluate_expression(expression_statement.Expression)
	litter.Dump(result)
	return result
}
