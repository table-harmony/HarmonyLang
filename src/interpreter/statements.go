package interpreter

import (
	"fmt"
	"reflect"

	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/lexer"
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

	err = scope.Declare(&variable)

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

	function := FunctionReference{
		identifier: expectedStatement.Identifier,
		value: FunctionValue{
			parameters: expectedStatement.Parameters,
			body:       expectedStatement.Body,
			returnType: evaluate_type(expectedStatement.ReturnType),
			closure:    scope,
		},
	}

	err = scope.Declare(&function)

	if err != nil {
		panic(err)
	}
}

func evaluate_assignment_statement(statement ast.Statement, scope *Scope) {
	expectedStatement, err := ast.ExpectStatement[ast.AssignmentStatement](statement)
	if err != nil {
		panic(err)
	}

	var ref Reference
	switch assigne := expectedStatement.Assigne.(type) {
	case ast.SymbolExpression:
		ref, err = scope.Resolve(assigne.Value)
		if err != nil {
			panic(fmt.Sprintf("cannot assign to undefined variable %s", assigne.Value))
		}

	case ast.PrefixExpression:
		if assigne.Operator.Kind != lexer.STAR {
			panic("invalid assignment target")
		}

		value := evaluate_expression(assigne.Right, scope)
		ptr, err := ExpectValue[*Pointer](value)
		if err != nil {
			panic("cannot dereference non-pointer type")
		}
		ref = ptr.Deref()

	case ast.ComputedMemberExpression:
		panic("TODO: computed member expression in assignment statement evaluation")

	case ast.MemberExpression:
		panic("TODO: member expression in assignment statement evaluation")

	default:
		panic("invalid assignment target")
	}

	value := evaluate_expression(expectedStatement.Value, scope)

	if err := ref.Store(value); err != nil {
		panic(err)
	}
}
