package interpreter

import (
	"fmt"
	"reflect"

	"github.com/sanity-io/litter"
	"github.com/table-harmony/HarmonyLang/src/ast"
)

func evalute_current_statement(interpreter *interpreter, enviorment *Environment) {
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

// TODO: fix break to end loop not the block
func evaluate_block_statement(statement ast.Statement, env *Environment) {
	expected_statement, err := ast.ExpectStatement[ast.BlockStatement](statement)

	if err != nil {
		panic(err)
	}

	sub_environment := create_enviorment(env)
	for _, underlying_statement := range expected_statement.Body {
		if _, ok := underlying_statement.(ast.BreakStatement); ok {
			break
		}

		if _, ok := underlying_statement.(ast.ContinueStatement); ok {
			break
		}

		evaluate_statement(underlying_statement, sub_environment)
	}
}

func evaluate_if_statement(statement ast.Statement, env *Environment) {
	expected_statement, err := ast.ExpectStatement[ast.IfStatement](statement)

	if err != nil {
		panic(err)
	}

	condition_value := evaluate_expression(expected_statement.Condition, env)
	condition_met, err := condition_value.AsBoolean()

	if err != nil {
		panic(err)
	}

	if condition_met {
		evaluate_block_statement(expected_statement.Consequent, env)
	} else if expected_statement.Alternate != nil {
		if alternate_if_statement, ok := expected_statement.Alternate.(ast.IfStatement); ok {
			evaluate_if_statement(alternate_if_statement, env)
		} else {
			evaluate_block_statement(expected_statement.Alternate, env)
		}
	}
}

func evaluate_for_statement(statement ast.Statement, env *Environment) {
	expected_statement, err := ast.ExpectStatement[ast.ForStatement](statement)

	if err != nil {
		panic(err)
	}

	sub_environment := create_enviorment(env)

	evaluate_statement(expected_statement.Initializer, sub_environment)

	condition_value := evaluate_expression(expected_statement.Condition, sub_environment)
	condition_met, err := condition_value.AsBoolean()

	if err != nil {
		panic(err)
	}

	for condition_met {
		evaluate_block_statement(ast.BlockStatement{Body: expected_statement.Body}, sub_environment)

		for _, post := range expected_statement.Post {
			evaluate_expression(post, sub_environment)
		}

		condition_value = evaluate_expression(expected_statement.Condition, sub_environment)
		condition_met, err = condition_value.AsBoolean()
		litter.Dump(env.variables)

		if err != nil {
			panic(err)
		}
	}
}

func evaluate_switch_statement(statement ast.Statement, env *Environment) {
	expected_statement, err := ast.ExpectStatement[ast.SwitchStatement](statement)

	litter.Dump(expected_statement)
	if err != nil {
		panic(err)
	}

	value := evaluate_expression(expected_statement.Value, env)
	var default_case ast.BlockStatement

	for _, case_statement := range expected_statement.Cases {
		if case_statement.Pattern == nil {
			default_case = ast.BlockStatement{Body: case_statement.Body}
			continue
		}

		case_value := evaluate_expression(case_statement.Pattern, env)

		//TODO: equality needs better support
		if case_value == value {
			sub_environment := create_enviorment(env)
			evaluate_block_statement(ast.BlockStatement{Body: case_statement.Body}, sub_environment)
			return
		}
	}

	if len(default_case.Body) > 0 {
		sub_environment := create_enviorment(env)
		evaluate_block_statement(default_case, sub_environment)
	}
}
