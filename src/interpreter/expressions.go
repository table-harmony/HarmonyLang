package interpreter

import (
	"fmt"
	"reflect"

	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/lexer"
)

func evaluate_expression(expression ast.Expression) RuntimeValue {
	expression_type := reflect.TypeOf(expression)

	if handler, exists := expression_lookup[expression_type]; exists {
		return handler(expression)
	} else {
		panic(fmt.Sprintf("No handler registered for statement type: %v", expression_type))
	}
}

func evaluate_primary_statement(expression ast.Expression) RuntimeValue {
	expression_type := reflect.TypeOf(expression)

	switch expression_type {
	case reflect.TypeOf(ast.NumberExpression{}):
		return NumberValue{
			Value: expression.(ast.NumberExpression).Value,
		}
	case reflect.TypeOf(ast.StringExpression{}):
		return StringValue{
			Value: expression.(ast.StringExpression).Value,
		}
	case reflect.TypeOf(ast.BooleanExpression{}):
		return BooleanValue{
			Value: expression.(ast.BooleanExpression).Value,
		}
	default:
		panic("Unknown statement type")
	}
}

func evalute_binary_expression(expression ast.Expression) RuntimeValue {
	binary_expression := expression.(ast.BinaryExpression)

	left := evaluate_expression(binary_expression.Left)
	right := evaluate_expression(binary_expression.Right)

	switch binary_expression.Operator.Kind {
	case lexer.PLUS:
		return NumberValue{Value: left.(NumberValue).Value + right.(NumberValue).Value}
	case lexer.DASH:
		return NumberValue{Value: left.(NumberValue).Value - right.(NumberValue).Value}
	case lexer.STAR:
		return NumberValue{Value: left.(NumberValue).Value * right.(NumberValue).Value}
	case lexer.SLASH:
		return NumberValue{Value: left.(NumberValue).Value / right.(NumberValue).Value}
	case lexer.PERCENT:
		return NumberValue{Value: float64(int(left.(NumberValue).Value) % int(right.(NumberValue).Value))}
	default:
		panic("Unknown operator")
	}
}

func evaluate_prefix_expression(expression ast.Expression) RuntimeValue {
	prefix_expression := expression.(ast.PrefixExpression)

	right := evaluate_expression(prefix_expression.Right)

	switch prefix_expression.Operator.Kind {
	case lexer.NOT:
		return BooleanValue{Value: !right.(BooleanValue).Value}
	case lexer.PLUS:
		return right
	case lexer.DASH:
		return NumberValue{Value: -right.(NumberValue).Value}
	default:
		panic("Unknown operator")
	}
}
