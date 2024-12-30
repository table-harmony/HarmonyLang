package interpreter

import (
	"fmt"
	"reflect"

	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/lexer"
)

func (interpreter *interpreter) evalute_current_statement(enviorment *Environment) {
	statement := interpreter.current_statement()
	evaluate_statement(statement, enviorment)
}

func evaluate_statement(statement ast.Statement, enviornment *Environment) {
	statementType := reflect.TypeOf(statement)

	if handler, exists := statement_lookup[statementType]; exists {
		handler(statement, enviornment)
	} else {
		panic(fmt.Sprintf("No handler registered for statement type: %v", statementType))
	}
}

func evaluate_expression_statement(statement ast.Statement, env *Environment) {
	expression_statement, err := ast.ExpectStatement[ast.ExpressionStatement](statement)

	if err != nil {
		panic(fmt.Sprintf("Expected expression statement, got %T", statement))
	}

	evaluate_expression(expression_statement.Expression, env)
}

func evaluate_variable_declaration_statement(statement ast.Statement, env *Environment) {
	expectedStatement, err := ast.ExpectStatement[ast.VariableDeclarationStatement](statement)

	if err != nil {
		panic(err)
	}

	err = env.declare_variable(RuntimeVariable{
		expectedStatement.Identifier,
		expectedStatement.IsConstant,
		evaluate_expression(expectedStatement.Value, env),
		evaluate_type(expectedStatement.ExplicitType),
	})

	if err != nil {
		panic(err)
	}
}

func evaluate_multi_variable_declaration_statement(statement ast.Statement, env *Environment) {
	expected_statement, err := ast.ExpectStatement[ast.MultiVariableDeclarationStatement](statement)

	if err != nil {
		panic(err)
	}

	for _, declaration := range expected_statement.Declarations {
		evaluate_variable_declaration_statement(declaration, env)
	}
}

func evaluate_continue_statement(statement ast.Statement, env *Environment) {
	panic(ContinueError{})
}

func evaluate_break_statement(statement ast.Statement, env *Environment) {
	panic(BreakError{})
}

func evaluate_return_statement(statement ast.Statement, env *Environment) {
	expectedStatement, err := ast.ExpectStatement[ast.ReturnStatement](statement)

	if err != nil {
		panic(err)
	}

	panic(ReturnError{evaluate_expression(expectedStatement.Value, env)})
}

// TODO: i dont iterate each statement cause i changed it from ast.BlockStatement to []ast.Statement
func evaluate_for_statement(statement ast.Statement, env *Environment) {
	expected_statement, err := ast.ExpectStatement[ast.ForStatement](statement)
	if err != nil {
		panic(err)
	}

	loop_env := create_environment(env)

	evaluate_statement(expected_statement.Initializer, loop_env)

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(BreakError); ok {
				return
			}
			panic(r)
		}
	}()

	for {
		condition_value := evaluate_expression(expected_statement.Condition, loop_env)
		expected_value, err := ExpectRuntimeValue[RuntimeBoolean](condition_value)

		if err != nil {
			panic(err)
		}

		if !expected_value.Value {
			break
		}

		func() {
			defer func() {
				if r := recover(); r != nil {
					if _, ok := r.(ContinueError); ok {
						return
					}
					panic(r)
				}
			}()

			//evaluate_block_statement(expected_statement, loop_env)
		}()

		for _, post := range expected_statement.Post {
			evaluate_expression(post, loop_env)
		}
	}
}

func evaluate_function_declaration_statement(statement ast.Statement, env *Environment) {
	expectedStatement, err := ast.ExpectStatement[ast.FunctionDeclarationStatment](statement)
	if err != nil {
		panic(err)
	}

	function := RuntimeFunction{
		Identifier: expectedStatement.Identifier,
		Parameters: expectedStatement.Parameters,
		Body:       expectedStatement.Body,
		ReturnType: evaluate_type(expectedStatement.ReturnType),
		Closure:    env,
	}

	err = env.declare_function(function)
	if err != nil {
		panic(err)
	}
}

func evaluate_assignment_statement(statement ast.Statement, env *Environment) {
	expectedStatement, err := ast.ExpectStatement[ast.AssignmentStatement](statement)
	if err != nil {
		panic(err)
	}

	assignable, err := evaluate_assignable(expectedStatement.Assigne, env)
	if err != nil {
		panic(err)
	}

	if expectedStatement.Operator.Kind == lexer.NULLISH_ASSIGNMENT {
		if assignable.getValue().getType() != NilType {
			return
		}
	}

	value := evaluate_expression(expectedStatement.Value, env)

	err = assignable.assign(value)
	if err != nil {
		panic(err)
	}
}

func evaluate_assignable(expression ast.Expression, env *Environment) (AssignableValue, error) {
	switch expression := expression.(type) {
	case ast.SymbolExpression:
		variable, err := env.get_variable(expression.Value)
		if err != nil {
			return nil, err
		}
		return &variable, nil
	case ast.MemberExpression:
		return nil, fmt.Errorf("not implemented yet")
	case ast.ComputedMemberExpression:
		return nil, fmt.Errorf("not implemented yet")
	case ast.PrefixExpression:
		//TODO: for pointers e.t.c
		return nil, fmt.Errorf("not implemented yet")
	default:
		return nil, fmt.Errorf("invalid assignable expression type: %T", expression)
	}
}
