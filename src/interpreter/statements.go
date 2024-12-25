package interpreter

import (
	"fmt"
	"reflect"

	"github.com/sanity-io/litter"
	"github.com/table-harmony/HarmonyLang/src/ast"
)

func evalute_statement(interpreter *interpreter, enviorment *Environment) {
	statement := interpreter.current_statement()
	statement_type := reflect.TypeOf(statement)

	if handler, exists := statement_lookup[statement_type]; exists {
		var value = handler(statement, enviorment)
		litter.Dump(value)
	} else {
		panic(fmt.Sprintf("No handler registered for statement type: %v", statement_type))
	}
}

func evaluate_expression_statement(statement ast.Statement, env *Environment) RuntimeValue {
	expression_statement, err := ast.ExpectStatement[ast.ExpressionStatement](statement)

	if err != nil {
		panic(fmt.Sprintf("Expected expression statement, got %v", statement))
	}

	return evaluate_expression(expression_statement.Expression, env)
}

func evaluate_variable_declaration_statement(statement ast.Statement, env *Environment) RuntimeValue {
	expected_statement, err := ast.ExpectStatement[ast.VariableDeclarationStatement](statement)

	if err != nil {
		panic(err)
	}

	variable := RuntimeVariable{
		Value:        evaluate_expression(expected_statement.Value, env),
		IsConstant:   expected_statement.IsConstant,
		Identifier:   expected_statement.Identifier,
		ExplicitType: expected_statement.ExplicitType,
	}
	err = env.declare_variable(variable)

	if err != nil {
		panic(err)
	}

	return variable
}
