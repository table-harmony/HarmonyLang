package interpreter

import (
	"fmt"
	"reflect"

	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/lexer"
)

func evaluate_expression(expression ast.Expression, env *Environment) RuntimeValue {
	expressionType := reflect.TypeOf(expression)

	if handler, exists := expression_lookup[expressionType]; exists {
		return handler(expression, env)
	} else {
		panic(fmt.Sprintf("No handler registered for statement type: %v", expressionType))
	}
}

func evaluate_primary_statement(expression ast.Expression, env *Environment) RuntimeValue {
	expressionType := reflect.TypeOf(expression)

	switch expressionType {
	case reflect.TypeOf(ast.NumberExpression{}):
		return RuntimeNumber{expression.(ast.NumberExpression).Value}
	case reflect.TypeOf(ast.StringExpression{}):
		return RuntimeString{expression.(ast.StringExpression).Value}
	case reflect.TypeOf(ast.BooleanExpression{}):
		return RuntimeBoolean{expression.(ast.BooleanExpression).Value}
	case reflect.TypeOf(ast.NilExpression{}):
		return RuntimeNil{}
	default:
		panic(fmt.Sprintf("Unknown statement type %s", expressionType))
	}
}

func evalute_binary_expression(expression ast.Expression, env *Environment) RuntimeValue {
	expectedExpression, err := ast.ExpectExpression[ast.BinaryExpression](expression)
	if err != nil {
		panic(err)
	}

	left := evaluate_expression(expectedExpression.Left, env)
	right := evaluate_expression(expectedExpression.Right, env)

	switch expectedExpression.Operator.Kind {
	case lexer.EQUALS, lexer.NOT_EQUALS:
		return handle_equality_operations(expectedExpression.Operator.Kind, left, right)
	case lexer.OR, lexer.AND:
		return handle_logical_operations(expectedExpression.Operator.Kind, left, right)
	case lexer.PLUS, lexer.DASH, lexer.STAR, lexer.SLASH, lexer.PERCENT:
		return handle_arithmetic_operations(expectedExpression.Operator.Kind, left, right)
	case lexer.LESS, lexer.GREATER, lexer.LESS_EQUALS, lexer.GREATER_EQUALS:
		return handle_relational_operations(expectedExpression.Operator.Kind, left, right)
	default:
		panic(fmt.Sprintf("Unknown operator %v", expectedExpression.Operator.Kind.ToString()))
	}
}

func handle_equality_operations(operator lexer.TokenKind, left, right RuntimeValue) RuntimeValue {
	leftValue := left.getValue()
	rightValue := right.getValue()

	switch operator {
	case lexer.EQUALS:
		return RuntimeBoolean{isEqual(leftValue, rightValue)}
	case lexer.NOT_EQUALS:
		return RuntimeBoolean{!isEqual(leftValue, rightValue)}
	default:
		panic(fmt.Sprintf("Unsupported equality operator: %v", operator.ToString()))
	}
}

func handle_logical_operations(operator lexer.TokenKind, left, right RuntimeValue) RuntimeValue {
	leftValue, err1 := ExpectRuntimeValue[RuntimeBoolean](left.getValue())
	rightValue, err2 := ExpectRuntimeValue[RuntimeBoolean](right.getValue())

	if err1 != nil || err2 != nil {
		panic(fmt.Sprintf("Logical operations require boolean values: %v, %v", err1, err2))
	}

	switch operator {
	case lexer.OR:
		return RuntimeBoolean{leftValue.Value || rightValue.Value}
	case lexer.AND:
		return RuntimeBoolean{leftValue.Value && rightValue.Value}
	default:
		panic(fmt.Sprintf("Unsupported logical operator: %v", operator.ToString()))
	}
}

func handle_arithmetic_operations(operator lexer.TokenKind, left, right RuntimeValue) RuntimeValue {
	if left == nil || right == nil {
		panic("Cannot perform arithmetic operation on nil value")
	}

	leftValue := left.getValue()
	rightValue := right.getValue()

	if leftValue.getType() == StringType || rightValue.getType() == StringType {
		if operator == lexer.PLUS {
			return RuntimeString{handle_string_concatenation(leftValue, rightValue).Value}
		}
		panic(fmt.Sprintf("Invalid operation %v between string types", operator.ToString()))
	}

	leftNum, leftErr := ExpectRuntimeValue[RuntimeNumber](leftValue)
	rightNum, rightErr := ExpectRuntimeValue[RuntimeNumber](rightValue)

	if leftErr != nil || rightErr != nil {
		panic(fmt.Sprintf("Invalid operation '%v' between %v and %v",
			operator.ToString(),
			left.getType().ToString(),
			right.getType().ToString()))
	}

	switch operator {
	case lexer.PLUS:
		return RuntimeNumber{leftNum.Value + rightNum.Value}
	case lexer.DASH:
		return RuntimeNumber{leftNum.Value - rightNum.Value}
	case lexer.STAR:
		return RuntimeNumber{leftNum.Value * rightNum.Value}
	case lexer.SLASH:
		if rightNum.Value == 0 {
			panic("Division by zero")
		}
		return RuntimeNumber{leftNum.Value / rightNum.Value}
	case lexer.PERCENT:
		if rightNum.Value == 0 {
			panic("Modulo by zero")
		}
		return RuntimeNumber{float64(int64(leftNum.Value) % int64(rightNum.Value))}
	default:
		panic(fmt.Sprintf("Unsupported arithmetic operator: %v", operator.ToString()))
	}
}

func handle_relational_operations(operator lexer.TokenKind, left, right RuntimeValue) RuntimeValue {
	leftNum, _ := ExpectRuntimeValue[RuntimeNumber](left.getValue())
	rightNum, _ := ExpectRuntimeValue[RuntimeNumber](right.getValue())

	switch operator {
	case lexer.LESS:
		return RuntimeBoolean{leftNum.Value < rightNum.Value}
	case lexer.GREATER:
		return RuntimeBoolean{leftNum.Value > rightNum.Value}
	case lexer.LESS_EQUALS:
		return RuntimeBoolean{leftNum.Value <= rightNum.Value}
	case lexer.GREATER_EQUALS:
		return RuntimeBoolean{leftNum.Value >= rightNum.Value}
	default:
		panic(fmt.Sprintf("Unsupported relational operator: %v", operator.ToString()))
	}
}

func handle_string_concatenation(left RuntimeValue, right RuntimeValue) RuntimeString {
	var leftStr string
	var rightStr string

	switch v := left.(type) {
	case RuntimeString:
		leftStr = v.Value
	case RuntimeNumber:
		leftStr = fmt.Sprintf("%g", v.Value)
	case RuntimeBoolean:
		leftStr = fmt.Sprintf("%v", v.Value)
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

	return RuntimeString{leftStr + rightStr}
}

func evaluate_prefix_expression(expression ast.Expression, env *Environment) RuntimeValue {
	prefixExpression := expression.(ast.PrefixExpression)

	right := evaluate_expression(prefixExpression.Right, env)
	rightType := right.getType()

	switch prefixExpression.Operator.Kind {
	case lexer.NOT:
		right, err := ExpectRuntimeValue[RuntimeBoolean](right)

		if err != nil {
			panic(fmt.Sprintf("Invalid operation %v with type %v",
				lexer.NOT.ToString(), rightType))
		}

		return RuntimeBoolean{!right.Value}
	case lexer.DASH:
		right, err := ExpectRuntimeValue[RuntimeNumber](right)

		if err != nil {
			panic(fmt.Sprintf("Invalid operation %v with type %v",
				lexer.DASH.ToString(), rightType))
		}

		return RuntimeNumber{-right.Value}
	case lexer.PLUS:
		right, err := ExpectRuntimeValue[RuntimeNumber](right)

		if err != nil {
			panic(fmt.Sprintf("Invalid operation %v with type %v",
				lexer.PLUS.ToString(), rightType))
		}

		return RuntimeNumber{right.Value}
	default:
		panic(fmt.Sprintf("Invalid operation %v with type %v",
			prefixExpression.Operator.Kind.ToString(), rightType.ToString()))
	}
}

func evaluate_symbol_expression(expression ast.Expression, env *Environment) RuntimeValue {
	expectedExpression, err := ast.ExpectExpression[ast.SymbolExpression](expression)
	if err != nil {
		panic(err)
	}

	variable, err := env.get_variable(expectedExpression.Value)
	if err == nil {
		return variable
	}

	function, err := env.get_function(expectedExpression.Value, make([]ast.Parameter, 0))
	if err == nil {
		return function
	}

	panic(err)
}

func evaluate_ternary_expression(expression ast.Expression, env *Environment) RuntimeValue {
	expectedExpression, err := ast.ExpectExpression[ast.TernaryExpression](expression)
	if err != nil {
		panic(err)
	}

	conditionValue := evaluate_expression(expectedExpression.Condition, env)

	expectedValue, err := ExpectRuntimeValue[RuntimeBoolean](conditionValue)
	if err != nil {
		panic(err)
	}

	if expectedValue.Value {
		return evaluate_expression(expectedExpression.Consequent, env)
	}

	return evaluate_expression(expectedExpression.Alternate, env)
}

func evaluate_block_expression(expression ast.Expression, env *Environment) RuntimeValue {
	expectedExpression, err := ast.ExpectExpression[ast.BlockExpression](expression)
	if err != nil {
		panic(err)
	}

	scope := create_environment(env)

	statements := expectedExpression.Statements
	for _, statement := range statements[:len(statements)-1] {
		evaluate_statement(statement, scope)
	}

	if len(statements) == 0 {
		return RuntimeNil{}
	}

	lastStatement := statements[len(statements)-1]
	if expressionStatement, ok := lastStatement.(ast.ExpressionStatement); ok {
		return evaluate_expression(expressionStatement.Expression, scope)
	}

	evaluate_statement(lastStatement, env)

	return RuntimeNil{}
}

func evaluate_if_expression(expression ast.Expression, env *Environment) RuntimeValue {
	expectedExpression, err := ast.ExpectExpression[ast.IfExpression](expression)

	if err != nil {
		panic(err)
	}

	conditionExpression := evaluate_expression(expectedExpression.Condition, env)
	expectedCondition, err := ExpectRuntimeValue[RuntimeBoolean](conditionExpression)

	if err != nil {
		panic(err)
	}

	if expectedCondition.Value {
		return evaluate_block_expression(expectedExpression.Consequent, env)
	} else if expectedExpression.Alternate != nil {
		alternateExpression, err := ast.ExpectExpression[ast.IfExpression](expectedExpression.Alternate)

		if err == nil {
			return evaluate_if_expression(alternateExpression, env)
		} else {
			return evaluate_block_expression(expectedExpression.Alternate, env)
		}
	}

	return RuntimeNil{}
}

func evaluate_switch_expression(expression ast.Expression, env *Environment) RuntimeValue {
	expectedExpression, err := ast.ExpectExpression[ast.SwitchExpression](expression)

	if err != nil {
		panic(err)
	}

	value := evaluate_expression(expectedExpression.Value, env)
	var defaultCase *ast.SwitchCaseStatement

	for _, switchCase := range expectedExpression.Cases {
		if switchCase.IsDefault {
			if defaultCase != nil {
				panic("duplicate default patterns in switch expression")
			}

			defaultCase = &switchCase
		}
	}

	for _, switchCase := range expectedExpression.Cases {
		for _, pattern := range switchCase.Patterns {
			casePatternValue := evaluate_expression(pattern, env)

			if isEqual(casePatternValue, value) {
				return evaluate_expression(switchCase.Body, env)
			}
		}
	}

	if defaultCase == nil {
		return nil
	}

	return evaluate_expression(defaultCase.Body, env)
}

func evaluate_call_expression(expression ast.Expression, env *Environment) (result RuntimeValue) {
	defer func() {
		if r := recover(); r != nil {
			switch err := r.(type) {
			case ReturnError:
				result = err.Value
			case error:
				// Re-panic other errors
				panic(err)
			default:
				// Re-panic unknown panic types
				panic(r)
			}
		}
	}()

	expectedExpression, err := ast.ExpectExpression[ast.CallExpression](expression)
	if err != nil {
		panic(err)
	}

	callerValue := evaluate_expression(expectedExpression.Caller, env)
	var function RuntimeFunction

	switch caller := callerValue.(type) {
	case RuntimeFunction:
		function = caller
	case RuntimeAnonymousFunction:
		function = RuntimeFunction{
			Parameters: caller.Parameters,
			Body:       caller.Body,
			ReturnType: caller.ReturnType,
			Closure:    caller.Closure,
		}
	default:
		panic(fmt.Sprintf("Cannot call value of type %v", caller.getType().ToString()))
	}

	functionEnv := create_environment(function.Closure)

	if len(expectedExpression.Params) != len(function.Parameters) {
		panic(fmt.Errorf("Expected %d arguments but got %d",
			len(function.Parameters), len(expectedExpression.Params)))
	}

	for i, param := range function.Parameters {
		argValue := evaluate_expression(expectedExpression.Params[i], env)
		if argValue == nil {
			argValue = RuntimeNil{}
		}

		err := functionEnv.declare_variable(RuntimeVariable{
			Identifier:   param.Name,
			Value:        argValue,
			ExplicitType: evaluate_type(param.Type),
		})
		if err != nil {
			panic(err)
		}
	}

	var lastValue RuntimeValue = RuntimeNil{}
	for _, statement := range function.Body {
		evaluate_statement(statement, functionEnv)
	}

	return lastValue
}

func evaluate_function_declaration_expression(expression ast.Expression, env *Environment) RuntimeValue {
	expectedExpression, err := ast.ExpectExpression[ast.FunctionDeclarationExpression](expression)
	if err != nil {
		panic(err)
	}

	return RuntimeAnonymousFunction{
		Parameters: expectedExpression.Parameters,
		Body:       expectedExpression.Body,
		ReturnType: evaluate_type(expectedExpression.ReturnType),
	}
}

func evaluate_try_catch_expression(expression ast.Expression, env *Environment) (result RuntimeValue) {
	expectedExpression, err := ast.ExpectExpression[ast.TryCatchExpression](expression)
	if err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			result = evaluate_expression(expectedExpression.CatchBlock, env)
		}
	}()

	result = evaluate_expression(expectedExpression.TryBlock, env)
	return result
}
