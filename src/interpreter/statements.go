package interpreter

import (
	"fmt"
	"reflect"

	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/core"
	"github.com/table-harmony/HarmonyLang/src/lexer"
)

func (interpreter *interpreter) evalute_current_statement(scope *core.Scope) {
	statement := interpreter.current_statement()
	evaluate_statement(statement, scope)
}

func evaluate_statement(statement ast.Statement, scope *core.Scope) {
	statementType := reflect.TypeOf(statement)

	if handler, exists := statement_lookup[statementType]; exists {
		handler(statement, scope)
	} else {
		panic(fmt.Sprintf("No handler registered for statement type: %v", statementType))
	}
}

func evaluate_expression_statement(statement ast.Statement, scope *core.Scope) {
	expectedStatement, err := ast.ExpectStatement[ast.ExpressionStatement](statement)
	if err != nil {
		panic(err)
	}

	evaluate_expression(expectedStatement.Expression, scope)
}

func evaluate_variable_declaration_statement(statement ast.Statement, scope *core.Scope) {
	expectedStatement, err := ast.ExpectStatement[ast.VariableDeclarationStatement](statement)

	if err != nil {
		panic(err)
	}

	variable := core.NewVariableReference(
		expectedStatement.Identifier,
		expectedStatement.IsConstant,
		evaluate_expression(expectedStatement.Value, scope),
		core.EvaluateType(expectedStatement.ExplicitType),
	)

	err = scope.Declare(variable)

	if err != nil {
		panic(err)
	}
}

func evaluate_multi_variable_declaration_statement(statement ast.Statement, scope *core.Scope) {
	expectedStatement, err := ast.ExpectStatement[ast.MultiVariableDeclarationStatement](statement)
	if err != nil {
		panic(err)
	}

	for _, declaration := range expectedStatement.Declarations {
		evaluate_variable_declaration_statement(declaration, scope)
	}
}

func evaluate_continue_statement(statement ast.Statement, scope *core.Scope) {
	panic(core.ContinueError{})
}

func evaluate_break_statement(statement ast.Statement, scope *core.Scope) {
	panic(core.BreakError{})
}

func evaluate_return_statement(statement ast.Statement, scope *core.Scope) {
	expectedStatement, err := ast.ExpectStatement[ast.ReturnStatement](statement)
	if err != nil {
		panic(err)
	}

	panic(core.NewReturnError(evaluate_expression(expectedStatement.Value, scope)))
}

// TODO: implement
func evaluate_for_statement(statement ast.Statement, env *core.Scope) {
}

func evaluate_function_declaration_statement(statement ast.Statement, scope *core.Scope) {
	expectedStatement, err := ast.ExpectStatement[ast.FunctionDeclarationStatment](statement)
	if err != nil {
		panic(err)
	}

	valuePtr := core.NewFunctionValue(
		expectedStatement.Parameters,
		expectedStatement.Body,
		core.EvaluateType(expectedStatement.ReturnType),
		scope,
	)

	ref := core.NewFunctionReference(
		expectedStatement.Identifier,
		*valuePtr,
	)

	err = scope.Declare(ref)

	if err != nil {
		panic(err)
	}
}

func evaluate_assignment_statement(statement ast.Statement, scope *core.Scope) {
	expectedStatement, err := ast.ExpectStatement[ast.AssignmentStatement](statement)
	if err != nil {
		panic(err)
	}

	var ref core.Reference
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
		ptr, err := core.ExpectValue[*core.Pointer](value)
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

func evaluate_throw_statement(statement ast.Statement, scope *core.Scope) {
	expectedStatement, err := ast.ExpectStatement[ast.ThrowStatement](statement)
	if err != nil {
		panic(err)
	}

	value := evaluate_expression(expectedStatement.Value, scope)
	panic(value)
}

func evaluate_type_declaration_statement(statement ast.Statement, scope *core.Scope) {
	panic("not implemented evaluate type declaration statement")
}
