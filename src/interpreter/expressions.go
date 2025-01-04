package interpreter

import (
	"fmt"
	"reflect"

	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/lexer"
)

func evaluate_expression(expression ast.Expression, scope *Scope) Value {
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

func evaluate_primary_expression(expression ast.Expression, scope *Scope) Value {
	switch expression := expression.(type) {
	case ast.NumberExpression:
		return NewNumber(expression.Value)
	case ast.StringExpression:
		return NewString(expression.Value)
	case ast.BooleanExpression:
		return NewBoolean(expression.Value)
	case ast.NilExpression:
		return NewNil()
	default:
		panic("Unknown expression type")
	}
}

func evaluate_prefix_expression(expression ast.Expression, scope *Scope) Value {
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
			return NewPointer(ref)

		case ast.ComputedMemberExpression:
			panic("TODO: computed member expression in prefix expression")

		case ast.MemberExpression:
			panic("TODO: member expression in prefix expression")

		case ast.CallExpression:
			result := evaluate_call_expression(right, scope)
			if ref, ok := result.(Reference); ok {
				return NewPointer(ref)
			}
			panic("cannot take address of function call result")

		default:
			panic("cannot take address of non-addressable expression")
		}
	case lexer.STAR:
		right := evaluate_expression(expectedExpression.Right, scope)
		derefed, err := Deref(right)
		if err != nil {
			panic(err)
		}
		return derefed
	}

	right := evaluate_expression(expectedExpression.Right, scope)

	switch expectedExpression.Operator.Kind {
	case lexer.NOT:
		rightResult, err := ExpectValue[Boolean](right)
		if err != nil {
			panic(fmt.Sprintf("Invalid operation %v with type %v",
				lexer.NOT.String(), right.Type().String()))
		}
		return NewBoolean(!rightResult.Value())

	case lexer.DASH:
		rightResult, err := ExpectValue[Number](right)
		if err != nil {
			panic(fmt.Sprintf("Invalid operation %v with type %v",
				lexer.DASH.String(), right.Type().String()))
		}
		return NewNumber(-rightResult.Value())

	case lexer.PLUS:
		rightResult, err := ExpectValue[Number](right)
		if err != nil {
			panic(fmt.Sprintf("Invalid operation %v with type %v",
				lexer.PLUS.String(), right.Type().String()))
		}
		return NewNumber(rightResult.Value())

	case lexer.TYPEOF:
		return NewValueType(right.Type())

	default:
		panic(fmt.Sprintf("Invalid operation %v with type %v",
			expectedExpression.Operator.Kind.String(), right.Type().String()))
	}
}

func evaluate_symbol_expression(expression ast.Expression, scope *Scope) Value {
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

func evaluate_binary_expression(expression ast.Expression, scope *Scope) Value {
	expectedExpression, err := ast.ExpectExpression[ast.BinaryExpression](expression)
	if err != nil {
		panic(err)
	}

	right := evaluate_expression(expectedExpression.Right, scope)
	left := evaluate_expression(expectedExpression.Left, scope)

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

func evaluate_ternary_expression(expression ast.Expression, scope *Scope) Value {
	expectedExpression, err := ast.ExpectExpression[ast.TernaryExpression](expression)
	if err != nil {
		panic(err)
	}

	conditionValue := evaluate_expression(expectedExpression.Condition, scope)

	expectedValue, err := ExpectValue[Boolean](conditionValue)
	if err != nil {
		panic(err)
	}

	if expectedValue.Value() {
		return evaluate_expression(expectedExpression.Consequent, scope)
	}

	return evaluate_expression(expectedExpression.Alternate, scope)
}

func evaluate_block_expression(expression ast.Expression, scope *Scope) Value {
	expectedExpression, err := ast.ExpectExpression[ast.BlockExpression](expression)
	if err != nil {
		panic(err)
	}

	blockScope := NewScope(scope)

	statements := expectedExpression.Statements
	for _, statement := range statements[:len(statements)-1] {
		evaluate_statement(statement, blockScope)
	}

	if len(statements) == 0 {
		return NewNil()
	}

	lastStatement := statements[len(statements)-1]
	if expressionStatement, ok := lastStatement.(ast.ExpressionStatement); ok {
		return evaluate_expression(expressionStatement.Expression, blockScope)
	}

	evaluate_statement(lastStatement, blockScope)

	return NewNil()
}

func evaluate_if_expression(expression ast.Expression, scope *Scope) Value {
	expectedExpression, err := ast.ExpectExpression[ast.IfExpression](expression)
	if err != nil {
		panic(err)
	}

	conditionExpression := evaluate_expression(expectedExpression.Condition, scope)
	expectedCondition, err := ExpectValue[Boolean](conditionExpression)

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

	return NewNil()
}

func evaluate_switch_expression(expression ast.Expression, scope *Scope) Value {
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
		return NewNil()
	}

	return evaluate_expression(defaultCase.Body, scope)
}

func evaluate_call_expression(expression ast.Expression, scope *Scope) (result Value) {
	expectedExpression, err := ast.ExpectExpression[ast.CallExpression](expression)
	if err != nil {
		panic(err)
	}

	var function Function
	switch caller := expectedExpression.Caller.(type) {
	case ast.SymbolExpression:
		ref, err := scope.Resolve(caller.Value)
		if err != nil {
			panic(fmt.Sprintf("cannot call undefined variable %s", caller.Value))
		}

		function, err = ExpectValue[FunctionValue](ref.Load())
		if err != nil {
			panic("cannot call non-function values")
		}

	case ast.PrefixExpression:
		if caller.Operator.Kind != lexer.STAR {
			panic("invalid call target")
		}

		value := evaluate_expression(caller.Right, scope)
		ptr, err := ExpectValue[*Pointer](value)
		if err != nil {
			panic("cannot dereference non-pointer type")
		}
		ref := ptr.Deref()

		function, err = ExpectValue[FunctionValue](ref.Load())
		if err != nil {
			panic("cannot call non-function values")
		}

	default:
		function = evaluate_expression(caller, scope).(Function)
		if false {
			panic("cannot call non-function values")
		}
	}

	params := make([]Value, 0)
	for _, param := range expectedExpression.Params {
		params = append(params, evaluate_expression(param, scope).Clone())
	}

	result, err = function.Call(params...)
	if err != nil {
		panic(err)
	}

	return result
}

func evaluate_function_declaration_expression(expression ast.Expression, scope *Scope) Value {
	expectedExpression, err := ast.ExpectExpression[ast.FunctionDeclarationExpression](expression)
	if err != nil {
		panic(err)
	}

	ptr := NewFunctionValue(
		expectedExpression.Parameters,
		expectedExpression.Body,
		EvaluateType(expectedExpression.ReturnType, scope),
		scope,
	)

	return *ptr
}

func evaluate_try_catch_expression(expression ast.Expression, scope *Scope) (result Value) {
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

func evaluate_array_instantiation_expression(expression ast.Expression, scope *Scope) Value {
	expectedExpression, err := ast.ExpectExpression[ast.ArrayInstantiationExpression](expression)
	if err != nil {
		panic(err)
	}

	size := evaluate_expression(expectedExpression.Size, scope)
	elementType := EvaluateType(expectedExpression.ElementType, scope)

	values := make([]Value, 0)
	for _, value := range expectedExpression.Elements {
		values = append(values, evaluate_expression(value, scope))
	}

	return NewArray(values, size, elementType)
}

func evaluate_map_instantiation_expression(expression ast.Expression, scope *Scope) Value {
	expectedExpression, err := ast.ExpectExpression[ast.MapInstantiationExpression](expression)
	if err != nil {
		panic(err)
	}

	keyType := EvaluateType(expectedExpression.KeyType, scope)
	valueType := EvaluateType(expectedExpression.ValueType, scope)

	entries := make([]MapEntry, 0)
	for _, entry := range expectedExpression.Entries {
		entries = append(entries, MapEntry{
			key:   evaluate_expression(entry.Key, scope),
			value: evaluate_expression(entry.Value, scope),
		})
	}

	return NewMap(entries, keyType, valueType)
}

func evaluate_slice_instantiation_expression(expression ast.Expression, scope *Scope) Value {
	expectedExpression, err := ast.ExpectExpression[ast.SliceInstantiationExpression](expression)
	if err != nil {
		panic(err)
	}

	var elementType Type
	if expectedExpression.ElementType != nil {
		elementType = EvaluateType(expectedExpression.ElementType, scope)
	} else {
		element := evaluate_expression(expectedExpression.Elements[0], scope)
		elementType = element.Type()
	}

	elements := make([]Value, 0)
	for _, value := range expectedExpression.Elements {
		elements = append(elements, evaluate_expression(value, scope))
	}

	return NewSlice(elements, elementType)
}

func evaluate_computed_member_expression(expression ast.Expression, scope *Scope) Value {
	expectedExpression, err := ast.ExpectExpression[ast.ComputedMemberExpression](expression)
	if err != nil {
		panic(err)
	}

	ownerValue := evaluate_expression(expectedExpression.Owner, scope)
	if ref, ok := ownerValue.(Reference); ok {
		ownerValue = ref.Load()
	}

	property := evaluate_expression(expectedExpression.Property, scope)
	if ref, ok := property.(Reference); ok {
		property = ref.Load()
	}

	switch owner := ownerValue.(type) {
	case Array:
		return owner.Get(property)
	case Map:
		return owner.Get(property)
	case Slice:
		return owner.Get(property)
	default:
		panic(fmt.Sprintf("Computed member expression not supported for type: %T", ownerValue))
	}
}

func evaluate_member_expression(expression ast.Expression, scope *Scope) Value {
	expectedExpression, err := ast.ExpectExpression[ast.MemberExpression](expression)
	if err != nil {
		panic(err)
	}

	ownerValue := evaluate_expression(expectedExpression.Owner, scope)
	if ref, ok := ownerValue.(Reference); ok {
		ownerValue = ref.Load()
	}

	switch owner := ownerValue.(type) {
	case Array:
		property, ok := expectedExpression.Property.(ast.SymbolExpression)
		if !ok {
			panic("Array method must be a symbol")
		}

		if method, exists := owner.methods[property.Value]; exists {
			return method
		}
		panic(fmt.Sprintf("Unknown array method: %s", property.Value))

	case Slice:
		property, ok := expectedExpression.Property.(ast.SymbolExpression)
		if !ok {
			panic("Slice method must be a symbol")
		}

		if method, exists := owner.methods[property.Value]; exists {
			return method
		}
		panic(fmt.Sprintf("Unknown slice method: %s", property.Value))

	case Map:
		property, ok := expectedExpression.Property.(ast.SymbolExpression)
		if !ok {
			panic("Map method must be a symbol")
		}

		if method, exists := owner.methods[property.Value]; exists {
			return method
		}

		panic(fmt.Sprintf("Unknown map method: %s", property.Value))
	}

	panic(fmt.Sprintf("Member expression not supported for type: %T", ownerValue))
}

func evaluate_range_expression(expression ast.Expression, scope *Scope) Value {
	expectedExpression, err := ast.ExpectExpression[ast.RangeExpression](expression)
	if err != nil {
		panic(err)
	}

	lower := evaluate_expression(expectedExpression.Lower, scope)
	upper := evaluate_expression(expectedExpression.Upper, scope)
	step := evaluate_expression(expectedExpression.Step, scope)

	lowerValue, err := ExpectValue[Number](lower)
	if err != nil {
		panic("Lower bound must be a number")
	}

	upperValue, err := ExpectValue[Number](upper)
	if err != nil {
		panic("Upper bound must be a number")
	}

	stepValue, err := ExpectValue[Number](step)
	if err != nil {
		panic("Step must be a number")
	}

	if stepValue.Value() == 0 {
		panic("Step cannot be zero")
	}

	values := make([]Value, 0)
	for i := lowerValue.Value(); (stepValue.Value() > 0 && i <= upperValue.Value()) || (stepValue.Value() < 0 && i >= upperValue.Value()); i += stepValue.Value() {
		values = append(values, NewNumber(i))
	}

	return NewArray(values, NewNumber(float64(len(values))), PrimitiveType{NumberType})
}
