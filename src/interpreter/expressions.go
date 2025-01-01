package interpreter

import (
	"fmt"
	"reflect"

	"github.com/sanity-io/litter"
	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/core"
	"github.com/table-harmony/HarmonyLang/src/lexer"
)

func evaluate_expression(expression ast.Expression, scope *core.Scope) core.Value {
	if expression == nil {
		return nil
	}

	expressionType := reflect.TypeOf(expression)

	if handler, exists := expression_lookup[expressionType]; exists {
		return handler(expression, scope)
	} else {
		panic(fmt.Sprintf("No handler registered for statement type: %v", expressionType))
	}
}

func evaluate_primary_expression(expression ast.Expression, scope *core.Scope) core.Value {
	switch expression := expression.(type) {
	case ast.NumberExpression:
		return core.NewNumber(expression.Value)
	case ast.StringExpression:
		return core.NewString(expression.Value)
	case ast.BooleanExpression:
		return core.NewBoolean(expression.Value)
	case ast.NilExpression:
		return core.NewNil()
	default:
		panic("Unknown expression type")
	}
}

func evaluate_prefix_expression(expression ast.Expression, scope *core.Scope) core.Value {
	expectedExpression, err := ast.ExpectExpression[ast.PrefixExpression](expression)
	if err != nil {
		panic(err)
	}

	switch expectedExpression.Operator.Kind {
	case lexer.AMPERSAND:
		switch right := expectedExpression.Right.(type) {
		case ast.SymbolExpression:
			ref, err := scope.Resolve(right.Value)
			if err != nil {
				panic(fmt.Sprintf("cannot take address of undefined variable %s", right.Value))
			}
			return core.NewPointer(ref)

		case ast.ComputedMemberExpression:
			panic("TODO: computed member expression in prefix expression")

		case ast.MemberExpression:
			panic("TODO: member expression in prefix expression")

		case ast.CallExpression:
			result := evaluate_call_expression(right, scope)
			if ref, ok := result.(core.Reference); ok {
				return core.NewPointer(ref)
			}
			panic("cannot take address of function call result")

		default:
			panic("cannot take address of non-addressable expression")
		}
	case lexer.STAR:
		right := evaluate_expression(expectedExpression.Right, scope)
		derefed, err := core.Deref(right)
		if err != nil {
			panic(err)
		}
		return derefed
	}

	right := evaluate_expression(expectedExpression.Right, scope)

	switch expectedExpression.Operator.Kind {
	case lexer.NOT:
		right, err := core.ExpectValue[core.Boolean](right)
		if err != nil {
			panic(fmt.Sprintf("Invalid operation %v with type %v",
				lexer.NOT.String(), right.Type().String()))
		}
		return core.NewBoolean(!right.Value())

	case lexer.DASH:
		right, err := core.ExpectValue[core.Number](right)
		if err != nil {
			panic(fmt.Sprintf("Invalid operation %v with type %v",
				lexer.DASH.String(), right.Type().String()))
		}
		return core.NewNumber(-right.Value())

	case lexer.PLUS:
		right, err := core.ExpectValue[core.Number](right)
		if err != nil {
			panic(fmt.Sprintf("Invalid operation %v with type %v",
				lexer.PLUS.String(), right.Type().String()))
		}
		return core.NewNumber(right.Value())

	default:
		panic(fmt.Sprintf("Invalid operation %v with type %v",
			expectedExpression.Operator.Kind.String(), right.Type().String()))
	}
}

func evaluate_symbol_expression(expression ast.Expression, scope *core.Scope) core.Value {
	expectedExpression, err := ast.ExpectExpression[ast.SymbolExpression](expression)
	if err != nil {
		panic(err)
	}

	ref, err := scope.Resolve(expectedExpression.Value)
	if err == nil {
		return ref.Load()
	}

	panic(fmt.Errorf("the name '%v' does not exist in the current scope", expectedExpression.Value))
}

func evaluate_binary_expression(expression ast.Expression, scope *core.Scope) core.Value {
	expectedExpression, err := ast.ExpectExpression[ast.BinaryExpression](expression)
	if err != nil {
		panic(err)
	}

	left := evaluate_expression(expectedExpression.Left, scope)
	right := evaluate_expression(expectedExpression.Right, scope)

	switch expectedExpression.Operator.Kind {
	case lexer.PLUS:
		return evaluate_addition(left, right)
	case lexer.DASH:
		return evaluate_subtraction(left, right)
	case lexer.STAR:
		return evaluate_multiplication(left, right)
	case lexer.SLASH:
		return evaluate_division(left, right)
	case lexer.PERCENT:
		return evaluate_modulo(left, right)
	case lexer.LESS:
		return evaluate_less_than(left, right)
	case lexer.LESS_EQUALS:
		return evaluate_less_equals(left, right)
	case lexer.GREATER:
		return evaluate_greater_than(left, right)
	case lexer.GREATER_EQUALS:
		return evaluate_greater_equals(left, right)
	case lexer.EQUALS:
		return evaluate_equals(left, right)
	case lexer.NOT_EQUALS:
		return evaluate_not_equals(left, right)
	case lexer.OR:
		return evaluate_logical_or(left, right)
	case lexer.AND:
		return evaluate_logical_and(left, right)
	default:
		panic(fmt.Sprintf("unknown binary operator: %v", expectedExpression.Operator.Kind))
	}
}

func evaluate_ternary_expression(expression ast.Expression, scope *core.Scope) core.Value {
	expectedExpression, err := ast.ExpectExpression[ast.TernaryExpression](expression)
	if err != nil {
		panic(err)
	}

	conditionValue := evaluate_expression(expectedExpression.Condition, scope)

	expectedValue, err := core.ExpectValue[core.Boolean](conditionValue)
	if err != nil {
		panic(err)
	}

	if expectedValue.Value() {
		return evaluate_expression(expectedExpression.Consequent, scope)
	}

	return evaluate_expression(expectedExpression.Alternate, scope)
}

func evaluate_block_expression(expression ast.Expression, scope *core.Scope) core.Value {
	expectedExpression, err := ast.ExpectExpression[ast.BlockExpression](expression)
	if err != nil {
		panic(err)
	}

	blockScope := core.NewScope(scope)

	statements := expectedExpression.Statements
	for _, statement := range statements[:len(statements)-1] {
		evaluate_statement(statement, blockScope)
	}

	if len(statements) == 0 {
		return core.NewNil()
	}

	lastStatement := statements[len(statements)-1]
	if expressionStatement, ok := lastStatement.(ast.ExpressionStatement); ok {
		return evaluate_expression(expressionStatement.Expression, blockScope)
	}

	evaluate_statement(lastStatement, blockScope)

	return core.NewNil()
}

func evaluate_if_expression(expression ast.Expression, scope *core.Scope) core.Value {
	expectedExpression, err := ast.ExpectExpression[ast.IfExpression](expression)
	if err != nil {
		panic(err)
	}

	conditionExpression := evaluate_expression(expectedExpression.Condition, scope)
	expectedCondition, err := core.ExpectValue[core.Boolean](conditionExpression)

	if err != nil {
		panic(err)
	}

	if expectedCondition.Value() {
		return evaluate_block_expression(expectedExpression.Consequent, scope)
	} else if expectedExpression.Alternate != nil {
		alternateExpression, err := ast.ExpectExpression[ast.IfExpression](expectedExpression.Alternate)

		if err == nil {
			return evaluate_if_expression(alternateExpression, scope)
		} else {
			return evaluate_block_expression(expectedExpression.Alternate, scope)
		}
	}

	return core.NewNil()
}

func evaluate_switch_expression(expression ast.Expression, scope *core.Scope) core.Value {
	expectedExpression, err := ast.ExpectExpression[ast.SwitchExpression](expression)
	if err != nil {
		panic(err)
	}

	value := evaluate_expression(expectedExpression.Value, scope)
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
			casePatternValue := evaluate_expression(pattern, scope)

			//TODO: IS THIS EQUALITY OK ? sense we are comparing Value
			if casePatternValue == value {
				return evaluate_expression(switchCase.Body, scope)
			}
		}
	}

	if defaultCase == nil {
		return core.NewNil()
	}

	return evaluate_expression(defaultCase.Body, scope)
}

func evaluate_call_expression(expression ast.Expression, scope *core.Scope) (result core.Value) {
	defer func() {
		if r := recover(); r != nil {
			switch err := r.(type) {
			case core.ReturnError:
				result = err.Value()
			default:
				panic(r)
			}
		}
	}()

	expectedExpression, err := ast.ExpectExpression[ast.CallExpression](expression)
	if err != nil {
		panic(err)
	}

	var function core.FunctionValue
	switch caller := expectedExpression.Caller.(type) {
	case ast.SymbolExpression:
		ref, err := scope.Resolve(caller.Value)
		if err != nil {
			panic(fmt.Sprintf("cannot assign to undefined variable %s", caller.Value))
		}

		function, err = core.ExpectValue[core.FunctionValue](ref.Load())
		if err != nil {
			panic("cannot call non-function values")
		}

	case ast.PrefixExpression:
		if caller.Operator.Kind != lexer.STAR {
			panic("invalid call target")
		}

		value := evaluate_expression(caller.Right, scope)
		ptr, err := core.ExpectValue[*core.Pointer](value)
		if err != nil {
			panic("cannot dereference non-pointer type")
		}
		ref := ptr.Deref()

		function, err = core.ExpectValue[core.FunctionValue](ref.Load())
		if err != nil {
			panic("cannot call non-function values")
		}

	case ast.FunctionDeclarationExpression:
		function, err = core.ExpectValue[core.FunctionValue](evaluate_expression(caller, scope))
		if err != nil {
			panic("cannot call non-function values")
		}

	case ast.ComputedMemberExpression:
		panic("TODO: computed member expression in call expression evaluation")

	case ast.MemberExpression:
		panic("TODO: member expression in call expression evaluation")

	default:
		panic("invalid call target")
	}
	litter.Dump(function)

	params := make([]core.Value, 0)
	for _, param := range expectedExpression.Params {
		params = append(params, evaluate_expression(param, scope))
	}
	litter.Dump("s123123ss")

	functionScope, err := function.CreateScope(params)
	if err != nil {
		panic(err)
	}
	litter.Dump("sss")
	for _, statement := range function.Body() {
		evaluate_statement(statement, functionScope)
	}

	return core.NewNil()
}

func evaluate_function_declaration_expression(expression ast.Expression, scope *core.Scope) core.Value {
	expectedExpression, err := ast.ExpectExpression[ast.FunctionDeclarationExpression](expression)
	if err != nil {
		panic(err)
	}

	ptr := core.NewFunctionValue(
		expectedExpression.Parameters,
		expectedExpression.Body,
		core.EvaluateType(expectedExpression.ReturnType),
		scope,
	)

	return *ptr
}

func evaluate_try_catch_expression(expression ast.Expression, scope *core.Scope) (result core.Value) {
	expectedExpression, err := ast.ExpectExpression[ast.TryCatchExpression](expression)
	if err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			result = evaluate_expression(expectedExpression.CatchBlock, scope)
		}
	}()

	result = evaluate_expression(expectedExpression.TryBlock, scope)

	return
}
