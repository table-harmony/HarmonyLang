package interpreter

import (
	"fmt"
	"reflect"

	"github.com/table-harmony/HarmonyLang/src/ast"
)

func (interpreter *interpreter) evalute_current_statement(scope *Scope) {
	statement := interpreter.current_statement()
	evaluate_statement(statement, scope)
}

func evaluate_statement(statement ast.Statement, scope *Scope) {
	statementType := reflect.TypeOf(statement)

	if handler, exists := statement_lookup[statementType]; exists {
		handler(statement, scope)
	} else {
		panic(fmt.Sprintf("No handler registered for statement type: %v", statementType))
	}
}

func evaluate_expression_statement(statement ast.Statement, scope *Scope) {
	expectedStatement, err := ast.ExpectStatement[ast.ExpressionStatement](statement)
	if err != nil {
		panic(err)
	}

	evaluate_expression(expectedStatement.Expression, scope)
}

func evaluate_variable_declaration_statement(statement ast.Statement, scope *Scope) {
	expectedStatement, err := ast.ExpectStatement[ast.VariableDeclarationStatement](statement)

	if err != nil {
		panic(err)
	}

	variable := VariableReference{
		identifier:   expectedStatement.Identifier,
		isConstant:   expectedStatement.IsConstant,
		value:        evaluate_expression(expectedStatement.Value, scope),
		explicitType: evaluate_type(expectedStatement.ExplicitType),
	}

	scope.DeclareVariable(&variable)

	if err != nil {
		panic(err)
	}
}

func evaluate_multi_variable_declaration_statement(statement ast.Statement, scope *Scope) {
	expectedStatement, err := ast.ExpectStatement[ast.MultiVariableDeclarationStatement](statement)
	if err != nil {
		panic(err)
	}

	for _, declaration := range expectedStatement.Declarations {
		evaluate_variable_declaration_statement(declaration, scope)
	}
}

func evaluate_continue_statement(statement ast.Statement, scope *Scope) {
	panic(ContinueError{})
}

func evaluate_break_statement(statement ast.Statement, scope *Scope) {
	panic(BreakError{})
}

func evaluate_return_statement(statement ast.Statement, scope *Scope) {
	expectedStatement, err := ast.ExpectStatement[ast.ReturnStatement](statement)
	if err != nil {
		panic(err)
	}

	panic(ReturnError{evaluate_expression(expectedStatement.Value, scope)})
}

// TODO: implement
func evaluate_for_statement(statement ast.Statement, env *Scope) {
}

func evaluate_function_declaration_statement(statement ast.Statement, scope *Scope) {
	expectedStatement, err := ast.ExpectStatement[ast.FunctionDeclarationStatment](statement)
	if err != nil {
		panic(err)
	}

	//TODO: create a function reference
	//TODO: declare the reference in the scope
}

func evaluate_assignment_statement(statement ast.Statement, scope *Scope) {
	expectedStatement, err := ast.ExpectStatement[ast.AssignmentStatement](statement)
	if err != nil {
		panic(err)
	}

	//TODO: get the reference of the assigne
	//TODO: store the new value in the reference
}
