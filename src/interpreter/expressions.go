package interpreter

import (
	"fmt"
	"reflect"

	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/lexer"
)

func evaluate_expression(expression ast.Expression, env *Environment) RuntimeValue {
	if expression == nil {
		return nil
	}

	expression_type := reflect.TypeOf(expression)

	if handler, exists := expression_lookup[expression_type]; exists {
		return handler(expression, env)
	} else {
		panic(fmt.Sprintf("No handler registered for statement type: %v", expression_type))
	}
}

func evaluate_primary_statement(expression ast.Expression, env *Environment) RuntimeValue {
	expression_type := reflect.TypeOf(expression)

	switch expression_type {
	case reflect.TypeOf(ast.NumberExpression{}):
		return NumberRuntime{
			Value: expression.(ast.NumberExpression).Value,
		}
	case reflect.TypeOf(ast.StringExpression{}):
		return StringRuntime{
			Value: expression.(ast.StringExpression).Value,
		}
	case reflect.TypeOf(ast.BooleanExpression{}):
		return BooleanRuntime{
			Value: expression.(ast.BooleanExpression).Value,
		}
	default:
		panic("Unknown statement type")
	}
}

func evalute_binary_expression(expression ast.Expression, env *Environment) RuntimeValue {
	binary_expression := expression.(ast.BinaryExpression)

	left := evaluate_expression(binary_expression.Left, env)
	right := evaluate_expression(binary_expression.Right, env)

	left_type, right_type := left.Type(), right.Type()

	switch binary_expression.Operator.Kind {
	case lexer.PLUS:
		if left_type == StringType || right_type == StringType {
			left_value, err1 := left.AsString()
			right_value, err2 := right.AsString()

			if err1 == nil && err2 == nil {
				return StringRuntime{Value: left_value + right_value}
			}
		}

		if left_type == NumberType || right_type == NumberType {
			left_value, err1 := left.AsNumber()
			right_value, err2 := right.AsNumber()

			if err1 == nil && err2 == nil {
				return NumberRuntime{Value: left_value + right_value}
			}
		}

		panic("Operand + not supported for these types")
	case lexer.DASH:
		left_value, err1 := left.AsNumber()
		right_value, err2 := right.AsNumber()

		if err1 == nil && err2 == nil {
			return NumberRuntime{Value: left_value - right_value}
		}

		panic("Operand - not supported for these types")
	case lexer.STAR:
		left_value, err1 := left.AsNumber()
		right_value, err2 := right.AsNumber()

		if err1 == nil && err2 == nil {
			return NumberRuntime{Value: left_value * right_value}
		}

		panic("Operand * not supported for these types")
	case lexer.SLASH:
		left_value, err1 := left.AsNumber()
		right_value, err2 := right.AsNumber()

		if err1 == nil && err2 == nil {
			return NumberRuntime{Value: left_value / right_value}
		}

		panic("Operand / not supported for these types")
	case lexer.PERCENT:
		left_value, err1 := left.AsNumber()
		right_value, err2 := right.AsNumber()

		if err1 == nil && err2 == nil {
			return NumberRuntime{Value: float64(int(left_value) % int(right_value))}
		}

		panic("Operand % not supported for these types")
	case lexer.AND:
		left_value, err1 := left.AsBoolean()
		right_value, err2 := right.AsBoolean()

		if err1 == nil && err2 == nil {
			return BooleanRuntime{Value: left_value && right_value}
		}

		panic("Operand && not supported for these types")
	case lexer.OR:
		left_value, err1 := left.AsBoolean()
		right_value, err2 := right.AsBoolean()

		if err1 == nil && err2 == nil {
			return BooleanRuntime{Value: left_value || right_value}
		}

		panic("Operand || not supported for these types")
	case lexer.EQUALS:
		if left_type == StringType || right_type == StringType {
			left_value, err1 := left.AsString()
			right_value, err2 := right.AsString()

			if err1 == nil && err2 == nil {
				return BooleanRuntime{Value: left_value == right_value}
			}
		}

		if left_type == NumberType || right_type == NumberType {
			left_value, err1 := left.AsNumber()
			right_value, err2 := right.AsNumber()

			if err1 == nil && err2 == nil {
				return BooleanRuntime{Value: left_value == right_value}
			}
		}

		if left_type == BooleanType || right_type == BooleanType {
			left_value, err1 := left.AsBoolean()
			right_value, err2 := right.AsBoolean()

			if err1 == nil && err2 == nil {
				return BooleanRuntime{Value: left_value == right_value}
			}
		}

		panic("Operand == not supported for these types")
	case lexer.NOT_EQUALS:
		if left_type == StringType || right_type == StringType {
			left_value, err1 := left.AsString()
			right_value, err2 := right.AsString()

			if err1 == nil && err2 == nil {
				return BooleanRuntime{Value: left_value != right_value}
			}
		}

		if left_type == NumberType || right_type == NumberType {
			left_value, err1 := left.AsNumber()
			right_value, err2 := right.AsNumber()

			if err1 == nil && err2 == nil {
				return BooleanRuntime{Value: left_value != right_value}
			}
		}

		if left_type == BooleanType || right_type == BooleanType {
			left_value, err1 := left.AsBoolean()
			right_value, err2 := right.AsBoolean()

			if err1 == nil && err2 == nil {
				return BooleanRuntime{Value: left_value != right_value}
			}
		}

		panic("Operand != not supported for these types")
	case lexer.LESS:
		left_value, err1 := left.AsNumber()
		right_value, err2 := right.AsNumber()

		if err1 == nil && err2 == nil {
			return BooleanRuntime{Value: left_value < right_value}
		}

		left_str, err1 := left.AsString()
		right_str, err2 := right.AsString()

		if err1 == nil && err2 == nil {
			return BooleanRuntime{Value: left_str < right_str}
		}

		panic("Operand < not supported for these types")
	case lexer.LESS_EQUALS:
		left_value, err1 := left.AsNumber()
		right_value, err2 := right.AsNumber()

		if err1 == nil && err2 == nil {
			return BooleanRuntime{Value: left_value <= right_value}
		}

		left_str, err1 := left.AsString()
		right_str, err2 := right.AsString()

		if err1 == nil && err2 == nil {
			return BooleanRuntime{Value: left_str <= right_str}
		}

		panic("Operand <= not supported for these types")
	case lexer.GREATER:
		left_value, err1 := left.AsNumber()
		right_value, err2 := right.AsNumber()

		if err1 == nil && err2 == nil {
			return BooleanRuntime{Value: left_value > right_value}
		}

		left_str, err1 := left.AsString()
		right_str, err2 := right.AsString()

		if err1 == nil && err2 == nil {
			return BooleanRuntime{Value: left_str > right_str}
		}

		panic("Operand > not supported for these types")
	case lexer.GREATER_EQUALS:
		left_value, err1 := left.AsNumber()
		right_value, err2 := right.AsNumber()

		if err1 == nil && err2 == nil {
			return BooleanRuntime{Value: left_value >= right_value}
		}

		left_str, err1 := left.AsString()
		right_str, err2 := right.AsString()

		if err1 == nil && err2 == nil {
			return BooleanRuntime{Value: left_str >= right_str}
		}

		panic("Operand >= not supported for these types")
	default:
		panic("Unknown operator")
	}
	//TODO: relational tokens as well such as LESS, LESS EQUALS e.t.c
}

func evaluate_prefix_expression(expression ast.Expression, env *Environment) RuntimeValue {
	prefix_expression := expression.(ast.PrefixExpression)

	right := evaluate_expression(prefix_expression.Right, env)

	switch prefix_expression.Operator.Kind {
	case lexer.NOT:
		right_value, err := right.AsBoolean()

		if err != nil {
			panic(err)
		}

		return BooleanRuntime{Value: !right_value}
	case lexer.DASH:
		right_value, err := right.AsNumber()

		if err != nil {
			panic(err)
		}

		return NumberRuntime{Value: -right_value}
	default:
		panic("Unknown operator")
	}
}

func evaluate_symbol_expression(expression ast.Expression, env *Environment) RuntimeValue {
	expected_expression, err := ast.ExpectExpression[ast.SymbolExpression](expression)

	if err != nil {
		panic(err)
	}

	variable, err := env.get_variable(expected_expression.Value)

	if err != nil {
		panic(err)
	}

	return variable
}

func evaluate_assignment_expression(expression ast.Expression, env *Environment) RuntimeValue {
	expected_expression, err := ast.ExpectExpression[ast.AssignmentExpression](expression)

	if err != nil {
		panic(err)
	}

	//TODO: this is not expected it could be member or call
	expected_assigne_expression, _ := ast.ExpectExpression[ast.SymbolExpression](expected_expression.Assigne)
	declared_variable, err := env.get_variable(expected_assigne_expression.Value)

	if err != nil {
		panic(err)
	}

	err = env.assign_variable(expected_assigne_expression.Value,
		evaluate_expression(expected_expression.Value, env))

	if err != nil {
		panic(err)
	}

	return declared_variable
}

func evaluate_switch_expression(expression ast.Expression, env *Environment) RuntimeValue {
	expected_expression, err := ast.ExpectExpression[ast.SwitchExpression](expression)

	if err != nil {
		panic(err)
	}

	value := evaluate_expression(expected_expression.Value, env)

	underscore_pattern := ast.SymbolExpression{Value: "_"}

	for _, case_statement := range expected_expression.Cases {
		if case_statement.Pattern == underscore_pattern {
			return evaluate_expression(case_statement.Value, env)
		}

		case_value := evaluate_expression(case_statement.Pattern, env)

		//TODO: equality needs better support
		if case_value == value {
			return evaluate_expression(case_statement.Value, env)
		}
	}

	return nil
}
