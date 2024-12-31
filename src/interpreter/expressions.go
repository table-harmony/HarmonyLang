package interpreter

import (
	"fmt"
	"reflect"

	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/lexer"
)

func evaluate_expression(expression ast.Expression, scope *Scope) Value {
	expressionType := reflect.TypeOf(expression)

	if handler, exists := expression_lookup[expressionType]; exists {
		return handler(expression, scope)
	} else {
		panic(fmt.Sprintf("No handler registered for statement type: %v", expressionType))
	}
}

func evaluate_primary_expression(expression ast.Expression, scope *Scope) Value {
	switch expression := expression.(type) {
	case ast.NumberExpression:
		return Number{expression.Value}
	case ast.StringExpression:
		return String{expression.Value}
	case ast.BooleanExpression:
		return Boolean{expression.Value}
	case ast.NilExpression:
		return Nil{}
	default:
		panic(fmt.Sprintf("Unknown expression type %s"))
	}
}

func evalute_binary_expression(expression ast.Expression, env *Environment) Value {
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

func handle_equality_operations(operator lexer.TokenKind, left, right RuntimeValue) Value {
	leftValue := left.getValue()
	rightValue := right.getValue()

	switch operator {
	case lexer.EQUALS:
		return Boolean{isEqual(leftValue, rightValue)}
	case lexer.NOT_EQUALS:
		return Boolean{!isEqual(leftValue, rightValue)}
	default:
		panic(fmt.Sprintf("Unsupported equality operator: %v", operator.ToString()))
	}
}

func handle_logical_operations(operator lexer.TokenKind, left, right RuntimeValue) Value {
	leftValue, err1 := ExpectRuntimeValue[Boolean](left.getValue())
	rightValue, err2 := ExpectRuntimeValue[Boolean](right.getValue())

	if err1 != nil || err2 != nil {
		panic(fmt.Sprintf("Logical operations require boolean values: %v, %v", err1, err2))
	}

	switch operator {
	case lexer.OR:
		return Boolean{leftValue.Value || rightValue.Value}
	case lexer.AND:
		return Boolean{leftValue.Value && rightValue.Value}
	default:
		panic(fmt.Sprintf("Unsupported logical operator: %v", operator.ToString()))
	}
}

func handle_arithmetic_operations(operator lexer.TokenKind, left, right RuntimeValue) Value {
	if left == nil || right == nil {
		panic("Cannot perform arithmetic operation on nil value")
	}

	leftValue := left.getValue()
	rightValue := right.getValue()

	if leftValue.getType() == StringType || rightValue.getType() == StringType {
		if operator == lexer.PLUS {
			return String{handle_string_concatenation(leftValue, rightValue).Value}
		}
		panic(fmt.Sprintf("Invalid operation %v between string types", operator.ToString()))
	}

	leftNum, leftErr := ExpectRuntimeValue[Number](leftValue)
	rightNum, rightErr := ExpectRuntimeValue[Number](rightValue)

	if leftErr != nil || rightErr != nil {
		panic(fmt.Sprintf("Invalid operation '%v' between %v and %v",
			operator.ToString(),
			left.getType().ToString(),
			right.getType().ToString()))
	}

	switch operator {
	case lexer.PLUS:
		return Number{leftNum.Value + rightNum.Value}
	case lexer.DASH:
		return Number{leftNum.Value - rightNum.Value}
	case lexer.STAR:
		return Number{leftNum.Value * rightNum.Value}
	case lexer.SLASH:
		if rightNum.Value == 0 {
			panic("Division by zero")
		}
		return Number{leftNum.Value / rightNum.Value}
	case lexer.PERCENT:
		if rightNum.Value == 0 {
			panic("Modulo by zero")
		}
		return Number{float64(int64(leftNum.Value) % int64(rightNum.Value))}
	default:
		panic(fmt.Sprintf("Unsupported arithmetic operator: %v", operator.ToString()))
	}
}

func handle_relational_operations(operator lexer.TokenKind, left, right RuntimeValue) Value {
	leftNum, _ := ExpectRuntimeValue[Number](left.getValue())
	rightNum, _ := ExpectRuntimeValue[Number](right.getValue())

	switch operator {
	case lexer.LESS:
		return Boolean{leftNum.Value < rightNum.Value}
	case lexer.GREATER:
		return Boolean{leftNum.Value > rightNum.Value}
	case lexer.LESS_EQUALS:
		return Boolean{leftNum.Value <= rightNum.Value}
	case lexer.GREATER_EQUALS:
		return Boolean{leftNum.Value >= rightNum.Value}
	default:
		panic(fmt.Sprintf("Unsupported relational operator: %v", operator.ToString()))
	}
}

func handle_string_concatenation(left RuntimeValue, right RuntimeValue) String {
	var leftStr string
	var rightStr string

	switch v := left.(type) {
	case String:
		leftStr = v.Value
	case Number:
		leftStr = fmt.Sprintf("%g", v.Value)
	case Boolean:
		leftStr = fmt.Sprintf("%v", v.Value)
	default:
		panic(fmt.Sprintf("Cannot convert type '%v' to string", left.getType().ToString()))
	}

	switch v := right.(type) {
	case String:
		rightStr = v.Value
	case Number:
		rightStr = fmt.Sprintf("%g", v.Value)
	case Boolean:
		rightStr = fmt.Sprintf("%t", v.Value)
	default:
		panic(fmt.Sprintf("Cannot convert type '%v' to string", right.getType().ToString()))
	}

	return String{leftStr + rightStr}
}

func evaluate_prefix_expression(expression ast.Expression, scope *Scope) Value {
	expectedExpression, err := ast.ExpectExpression[ast.PrefixExpression](expression)
	if err != nil {
		panic(err)
	}

	right := evaluate_expression(expectedExpression.Right, scope)

	switch expectedExpression.Operator.Kind {
	case lexer.NOT:
		right, err := ExpectValue[Boolean](right)
		if err != nil {
			panic(fmt.Sprintf("Invalid operation %v with type %v",
				lexer.NOT.ToString(), right.Type().String()))
		}

		return Boolean{!right.Value()}
	case lexer.DASH:
		right, err := ExpectValue[Number](right)
		if err != nil {
			panic(fmt.Sprintf("Invalid operation %v with type %v",
				lexer.DASH.ToString(), right.Type().String()))
		}

		return Number{-right.Value()}
	case lexer.PLUS:
		right, err := ExpectValue[Number](right)
		if err != nil {
			panic(fmt.Sprintf("Invalid operation %v with type %v",
				lexer.PLUS.ToString(), right.Type().String()))
		}

		return Number{right.Value()}
	case lexer.AMPERSAND:
		//\right, err := ExpectReference[*VariableReference](right)
		//\if err != nil {
		//\	panic(err)
		//\}

		return right
	case lexer.STAR:

		//return NewPointer(right)
	default:
		panic(fmt.Sprintf("Invalid operation %v with type %v",
			expectedExpression.Operator.Kind.ToString(), right.Type().String()))
	}
}

// TODO: this implementation is made so there cannot exist a duplicate identifier in a scope even for a function and a variable
func evaluate_symbol_expression(expression ast.Expression, env *Environment) Value {
	expectedExpression, err := ast.ExpectExpression[ast.SymbolExpression](expression)
	if err != nil {
		panic(err)
	}

	variable, err := env.get_variable(expectedExpression.Value)
	if err == nil {
		return variable
	}

	function, err := env.get_function(expectedExpression.Value)
	if err == nil {
		return function
	}

	panic(fmt.Errorf("the name '%v' does not exist in the current context", expectedExpression.Value))
}

func evaluate_ternary_expression(expression ast.Expression, env *Environment) Value {
	expectedExpression, err := ast.ExpectExpression[ast.TernaryExpression](expression)
	if err != nil {
		panic(err)
	}

	conditionValue := evaluate_expression(expectedExpression.Condition, env)

	expectedValue, err := ExpectRuntimeValue[Boolean](conditionValue)
	if err != nil {
		panic(err)
	}

	if expectedValue.Value {
		return evaluate_expression(expectedExpression.Consequent, env)
	}

	return evaluate_expression(expectedExpression.Alternate, env)
}

func evaluate_block_expression(expression ast.Expression, env *Environment) Value {
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
		return Nil{}
	}

	lastStatement := statements[len(statements)-1]
	if expressionStatement, ok := lastStatement.(ast.ExpressionStatement); ok {
		return evaluate_expression(expressionStatement.Expression, scope)
	}

	evaluate_statement(lastStatement, env)

	return Nil{}
}

func evaluate_if_expression(expression ast.Expression, env *Environment) Value {
	expectedExpression, err := ast.ExpectExpression[ast.IfExpression](expression)

	if err != nil {
		panic(err)
	}

	conditionExpression := evaluate_expression(expectedExpression.Condition, env)
	expectedCondition, err := ExpectRuntimeValue[Boolean](conditionExpression)

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

	return Nil{}
}

func evaluate_switch_expression(expression ast.Expression, env *Environment) Value {
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

func evaluate_call_expression(expression ast.Expression, env *Environment) (result Value) {
	defer func() {
		if r := recover(); r != nil {
			switch err := r.(type) {
			case ReturnError:
				result = err.Value
			default:
				panic(r)
			}
		}
	}()

	expectedExpression, err := ast.ExpectExpression[ast.CallExpression](expression)
	if err != nil {
		panic(err)
	}

	caller := evaluate_expression(expectedExpression.Caller, env)
	callerValue := caller.getValue()

	if callerValue.getType() != AnonymousFunctionType && callerValue.getType() != FunctionType {
		panic(fmt.Errorf("cannot call non-function value of type %v", callerValue.getType().ToString()))
	}

	var function RuntimeAnonymousFunction
	if callerValue.getType() == AnonymousFunctionType {
		function, _ = ExpectRuntimeValue[RuntimeAnonymousFunction](callerValue)
	} else {
		functionValue, err := ExpectRuntimeValue[RuntimeFunction](callerValue)
		if err != nil {
			panic(err)
		}

		function = RuntimeAnonymousFunction{
			Parameters: functionValue.Parameters,
			Body:       functionValue.Body,
			ReturnType: functionValue.ReturnType,
			Closure:    functionValue.Closure,
		}
	}

	if len(expectedExpression.Params) != len(function.Parameters) {
		panic(fmt.Errorf("expected %d arguments but got %d",
			len(function.Parameters), len(expectedExpression.Params)))
	}

	functionScope := create_environment(function.Closure)

	for index, param := range expectedExpression.Params {
		expectedParam := function.Parameters[index]
		expectedParamType := evaluate_type(expectedParam.Type)

		paramValue := evaluate_expression(param, env).getValue()

		if paramValue.getType() != expectedParamType && expectedParamType != AnyType {
			panic(fmt.Sprintf("parameter '%s' expected type '%s' but got '%s'",
				expectedParam.Name,
				expectedParamType.ToString(),
				paramValue.getType().ToString()))
		}

		if expectedParamType == FunctionType {
			passedFunction, err := ExpectRuntimeValue[RuntimeFunction](paramValue)

			if err != nil {
				panic(err)
			}

			err = functionScope.declare_function(RuntimeFunction{
				Identifier: expectedParam.Name,
				Parameters: passedFunction.Parameters,
				Body:       passedFunction.Body,
				ReturnType: passedFunction.ReturnType,
				Closure:    passedFunction.Closure,
			})

			if err != nil {
				panic(err)
			}
		} else {
			err = functionScope.declare_variable(RuntimeVariable{
				Identifier:   expectedParam.Name,
				IsConstant:   false,
				Value:        paramValue,
				ExplicitType: expectedParamType,
			})

			if err != nil {
				panic(err)
			}
		}
	}

	for _, statement := range function.Body {
		evaluate_statement(statement, functionScope)
	}

	return Nil{}
}

func evaluate_function_declaration_expression(expression ast.Expression, env *Environment) Value {
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

func evaluate_try_catch_expression(expression ast.Expression, env *Environment) (result Value) {
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
	return
}
