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
		return RuntimeNumber{
			Value: expression.(ast.NumberExpression).Value,
		}
	case reflect.TypeOf(ast.StringExpression{}):
		return RuntimeString{
			Value: expression.(ast.StringExpression).Value,
		}
	case reflect.TypeOf(ast.BooleanExpression{}):
		return RuntimeBoolean{
			Value: expression.(ast.BooleanExpression).Value,
		}
	case reflect.TypeOf(ast.NilExpression{}):
		return RuntimeNil{}
	default:
		panic(fmt.Sprintf("Unknown statement type %s", expression_type))
	}
}

// TODO: write a prettier evaluate binary expression method
func evalute_binary_expression(expression ast.Expression, env *Environment) RuntimeValue {
	expected_expression, err := ast.ExpectExpression[ast.BinaryExpression](expression)
	if err != nil {
		panic(err)
	}

	left := evaluate_expression(expected_expression.Left, env)
	right := evaluate_expression(expected_expression.Right, env)

	leftValue := left.getValue()
	rightValue := right.getValue()

	switch expected_expression.Operator.Kind {
	case lexer.EQUALS:
		return RuntimeBoolean{Value: isEqual(leftValue, rightValue)}
	case lexer.NOT_EQUALS:
		return RuntimeBoolean{Value: !isEqual(leftValue, rightValue)}
	case lexer.OR:
		leftValue, err1 := ExpectRuntimeValue[RuntimeBoolean](leftValue)
		rightValue, err2 := ExpectRuntimeValue[RuntimeBoolean](rightValue)

		if err1 != nil {
			panic(err1)
		}

		if err2 != nil {
			panic(err2)
		}

		return RuntimeBoolean{Value: leftValue.Value || rightValue.Value}
	case lexer.AND:
		leftValue, err1 := ExpectRuntimeValue[RuntimeBoolean](leftValue)
		rightValue, err2 := ExpectRuntimeValue[RuntimeBoolean](rightValue)

		if err1 != nil {
			panic(err1)
		}

		if err2 != nil {
			panic(err2)
		}

		return RuntimeBoolean{Value: leftValue.Value && rightValue.Value}
	case lexer.PLUS:
		if leftValue.getType() == StringType || rightValue.getType() == StringType {
			return handle_string_concatenation(leftValue, rightValue)
		}

		if leftValue.getType() == NumberType && rightValue.getType() == NumberType {
			leftNum, _ := ExpectRuntimeValue[RuntimeNumber](leftValue)
			rightNum, _ := ExpectRuntimeValue[RuntimeNumber](rightValue)
			return RuntimeNumber{Value: leftNum.Value + rightNum.Value}
		}
		panic(fmt.Sprintf("Invalid addition between types %v and %v",
			leftValue.getType(), rightValue.getType()))
	}

	if leftValue.getType() != NumberType || rightValue.getType() != NumberType {
		panic(fmt.Sprintf("Invalid operation %v between types %v and %v",
			expected_expression.Operator.Kind.ToString(), leftValue.getType(), rightValue.getType()))
	}

	leftNum, _ := ExpectRuntimeValue[RuntimeNumber](leftValue)
	rightNum, _ := ExpectRuntimeValue[RuntimeNumber](rightValue)

	switch expected_expression.Operator.Kind {
	case lexer.DASH:
		return RuntimeNumber{Value: leftNum.Value - rightNum.Value}
	case lexer.STAR:
		return RuntimeNumber{Value: leftNum.Value * rightNum.Value}
	case lexer.SLASH:
		if rightNum.Value == 0 {
			panic("Division by zero")
		}
		return RuntimeNumber{Value: leftNum.Value / rightNum.Value}
	case lexer.PERCENT:
		if rightNum.Value == 0 {
			panic("Modulo by zero")
		}
		return RuntimeNumber{Value: float64(int64(leftNum.Value) % int64(rightNum.Value))}
	case lexer.LESS:
		return RuntimeBoolean{Value: leftNum.Value < rightNum.Value}
	case lexer.GREATER:
		return RuntimeBoolean{Value: leftNum.Value > rightNum.Value}
	case lexer.LESS_EQUALS:
		return RuntimeBoolean{Value: leftNum.Value <= rightNum.Value}
	case lexer.GREATER_EQUALS:
		return RuntimeBoolean{Value: leftNum.Value >= rightNum.Value}
	}

	panic(fmt.Sprintf("Unknown operator %v", expected_expression.Operator.Kind.ToString()))
}

func handle_string_concatenation(left RuntimeValue, right RuntimeValue) RuntimeValue {
	var leftStr string
	var rightStr string

	switch v := left.(type) {
	case RuntimeString:
		leftStr = v.Value
	case RuntimeNumber:
		leftStr = fmt.Sprintf("%g", v.Value)
	case RuntimeBoolean:
		leftStr = fmt.Sprintf("%t", v.Value)
	default:
		panic(fmt.Sprintf("Cannot convert type '%v' to string", left.getType().ToString()))
	}

	switch v := right.(type) {
	case RuntimeString:
		rightStr = v.Value
	case RuntimeNumber:
		rightStr = fmt.Sprintf("%g", v.Value)
	case RuntimeBoolean:
		rightStr = fmt.Sprintf("%t", v.Value)
	default:
		panic(fmt.Sprintf("Cannot convert type '%v' to string", right.getType().ToString()))
	}

	return RuntimeString{Value: leftStr + rightStr}
}

func evaluate_prefix_expression(expression ast.Expression, env *Environment) RuntimeValue {
	prefix_expression := expression.(ast.PrefixExpression)

	right := evaluate_expression(prefix_expression.Right, env)
	rightType := right.getType()

	switch prefix_expression.Operator.Kind {
	case lexer.NOT:
		right, err := ExpectRuntimeValue[RuntimeBoolean](right)

		if err != nil {
			panic(fmt.Sprintf("Invalid operation %v with type %v",
				lexer.NOT.ToString(), rightType))
		}

		return RuntimeBoolean{Value: !right.Value}
	case lexer.DASH:
		right, err := ExpectRuntimeValue[RuntimeNumber](right)

		if err != nil {
			panic(fmt.Sprintf("Invalid operation %v with type %v",
				lexer.DASH.ToString(), rightType))
		}

		return RuntimeNumber{Value: -right.Value}
	case lexer.PLUS:
		right, err := ExpectRuntimeValue[RuntimeNumber](right)

		if err != nil {
			panic(fmt.Sprintf("Invalid operation %v with type %v",
				lexer.PLUS.ToString(), rightType))
		}

		return RuntimeNumber{Value: right.Value}
	default:
		panic(fmt.Sprintf("Invalid operation %v with type %v",
			prefix_expression.Operator.Kind.ToString(), rightType.ToString()))
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

func evaluate_switch_expression(expression ast.Expression, env *Environment) RuntimeValue {
	expectedExpression, err := ast.ExpectExpression[ast.SwitchExpression](expression)

	if err != nil {
		panic(err)
	}

	value := evaluate_expression(expectedExpression.Value, env)
	var defaultCase *ast.DefaultSwitchCase

	for _, switchCase := range expectedExpression.Cases {
		if defaultSwitchCase, ok := switchCase.(ast.DefaultSwitchCase); ok {
			if defaultCase != nil {
				panic("duplicate default patterns in switch expression")
			}

			defaultCase = &defaultSwitchCase
		}
	}

	for _, switchCase := range expectedExpression.Cases {
		for _, pattern := range switchCase.GetPatterns() {
			casePatternValue := evaluate_expression(pattern, env)

			if isEqual(casePatternValue, value) {
				scope := create_enviorment(env)
				return evaluate_expression(switchCase.GetValue(), scope)
			}
		}
	}

	if defaultCase == nil {
		return nil
	}

	return evaluate_expression(defaultCase.Value, env)
}

func evaluate_ternary_expression(expression ast.Expression, env *Environment) RuntimeValue {
	expected_expression, err := ast.ExpectExpression[ast.TernaryExpression](expression)

	if err != nil {
		panic(err)
	}

	condition_value := evaluate_expression(expected_expression.Condition, env)
	expected_value, err := ExpectRuntimeValue[RuntimeBoolean](condition_value)

	if err != nil {
		panic(err)
	}

	if expected_value.Value {
		return evaluate_expression(expected_expression.Consequent, env)
	}

	return evaluate_expression(expected_expression.Alternate, env)
}
