package interpreter

import (
	"fmt"
	"reflect"

	"github.com/table-harmony/HarmonyLang/src/ast"
)

func (interpreter *interpreter) evalute_current_statement(enviorment *Environment) {
	statement := interpreter.current_statement()
	evaluate_statement(statement, enviorment)
}

func evaluate_statement(statement ast.Statement, enviornment *Environment) {
	statement_type := reflect.TypeOf(statement)

	if handler, exists := statement_lookup[statement_type]; exists {
		handler(statement, enviornment)
	} else {
		panic(fmt.Sprintf("No handler registered for statement type: %v", statement_type))
	}
}

func evaluate_expression_statement(statement ast.Statement, env *Environment) {
	expression_statement, err := ast.ExpectStatement[ast.ExpressionStatement](statement)

	if err != nil {
		panic(fmt.Sprintf("Expected expression statement, got %v", statement))
	}

	evaluate_expression(expression_statement.Expression, env)
}

func evaluate_variable_declaration_statement(statement ast.Statement, env *Environment) {
	expected_statement, err := ast.ExpectStatement[ast.VariableDeclarationStatement](statement)

	if err != nil {
		panic(err)
	}

	err = env.declare_variable(RuntimeVariable{
		Value:        evaluate_expression(expected_statement.Value, env),
		IsConstant:   expected_statement.IsConstant,
		Identifier:   expected_statement.Identifier,
		ExplicitType: expected_statement.ExplicitType,
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

func evaluate_block_statement(statement ast.Statement, env *Environment) {
	expected_statement, err := ast.ExpectStatement[ast.BlockStatement](statement)

	if err != nil {
		panic(err)
	}

	scope := create_enviorment(env)
	for _, underlying_statement := range expected_statement.Body {
		evaluate_statement(underlying_statement, scope)
	}
}

func evaluate_if_statement(statement ast.Statement, env *Environment) {
	expected_statement, err := ast.ExpectStatement[ast.IfStatement](statement)

	if err != nil {
		panic(err)
	}

	condition_value := evaluate_expression(expected_statement.Condition, env)
	expected_value, err := ExpectRuntimeValue[RuntimeBoolean](condition_value)

	if err != nil {
		panic(err)
	}

	if expected_value.Value {
		evaluate_block_statement(expected_statement.Consequent, env)
	} else if expected_statement.Alternate != nil {
		if alternate_if_statement, ok := expected_statement.Alternate.(ast.IfStatement); ok {
			evaluate_if_statement(alternate_if_statement, env)
		} else {
			evaluate_block_statement(expected_statement.Alternate, env)
		}
	}
}

func evaluate_continue_statement(statement ast.Statement, env *Environment) {
	panic(ContinueError{})
}

func evaluate_break_statement(statement ast.Statement, env *Environment) {
	panic(BreakError{})
}

func evaluate_for_statement(statement ast.Statement, env *Environment) {
	expected_statement, err := ast.ExpectStatement[ast.ForStatement](statement)
	if err != nil {
		panic(err)
	}

	loop_env := create_enviorment(env)

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

			evaluate_block_statement(expected_statement, loop_env)
		}()

		for _, post := range expected_statement.Post {
			evaluate_expression(post, loop_env)
		}
	}
}

func evaluate_switch_statement(statement ast.Statement, env *Environment) {
	expected_statement, err := ast.ExpectStatement[ast.SwitchStatement](statement)

	if err != nil {
		panic(err)
	}

	value := evaluate_expression(expected_statement.Value, env)
	var default_case ast.BlockStatement

	for _, case_statement := range expected_statement.Cases {
		if case_statement.Pattern == nil {
			default_case = case_statement.Body
			continue
		}

		case_value := evaluate_expression(case_statement.Pattern, env)

		if isEqual(case_value, value) {
			sub_environment := create_enviorment(env)
			evaluate_block_statement(case_statement.Body, sub_environment)
			return
		}
	}

	if len(default_case.Body) > 0 {
		sub_environment := create_enviorment(env)
		evaluate_block_statement(default_case, sub_environment)
	}
}

func evaluate_function_declaration_statement(statement ast.Statement, env *Environment) {
	panic("Not implemented yet")
}
