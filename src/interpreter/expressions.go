package interpreter

import (
	"fmt"
	"reflect"

	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/helpers"
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

	if varRef, ok := right.(*VariableReference); ok {
		right = varRef.value
	}

	if varRef, ok := left.(*VariableReference); ok {
		left = varRef.value
	}

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

			if reflect.DeepEqual(casePatternValue, value) {
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

		var ok bool
		if function, ok = ref.Load().(Function); !ok {
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
		var ok bool
		value := evaluate_expression(caller, scope)
		function, ok = value.(Function)
		if ok {
			break
		}

		if ref, ok := value.(*FunctionReference); ok {
			function = ref.value
			break
		}

		if ref, ok := value.(*VariableReference); ok {
			function = ref.value.(Function)
			break
		}

		panic("cannot call non-function values")
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

	errorIdentifier := expectedExpression.ErrorIdentifier

	defer func() {
		if r := recover(); r != nil {
			catchScope := NewScope(scope)
			if errorIdentifier != "" {
				var value Value
				switch e := r.(type) {
				case error:
					value = NewError(e.Error())
				case string:
					value = NewError(e)
				default:
					value = NewError(fmt.Sprintf("%v", e))
				}
				ref := NewVariableReference(
					errorIdentifier,
					true,
					value,
					PrimitiveType{ErrorType},
				)
				catchScope.Declare(ref)
			}

			result = evaluate_expression(expectedExpression.CatchBlock, catchScope)
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
	case String:
		numberProperty, ok := property.(Number)
		if ok {
			if len(owner.value) <= int(numberProperty.value) {
				panic(fmt.Sprintf("Index out of range '%b'", numberProperty.value))
			}
			return NewString(string(owner.value[int(numberProperty.value)]))
		}

		propertyName, ok := property.(String)
		if !ok {
			panic("Computed member access must use string expression for property")
		}

		if method, exists := owner.methods[propertyName.value]; exists {
			return method
		}

		panic(fmt.Sprintf("Unknown string method: %s", propertyName.value))
	case Server:
		propertyName, ok := property.(String)
		if !ok {
			panic("Computed member access must use string expression for property")
		}

		if method, exists := owner.methods[propertyName.value]; exists {
			return method
		}

		panic(fmt.Sprintf("Unknown server method: %s", propertyName.value))
	case Response:
		propertyName, ok := property.(String)
		if !ok {
			panic("Computed member access must use string expression for property")
		}

		if method, exists := owner.Methods[propertyName.value]; exists {
			return method
		}
		value := reflect.ValueOf(owner)
		field := value.FieldByName(helpers.Capitalize(propertyName.value))
		if field.IsValid() {
			return convert_to_value(field.Interface())
		}

		panic(fmt.Sprintf("Unknown response method or property: %s", propertyName.value))
	case Request:
		propertyName, ok := property.(String)
		if !ok {
			panic("Computed member access must use string expression for property")
		}

		value := reflect.ValueOf(owner)
		field := value.FieldByName(helpers.Capitalize(propertyName.value))
		if field.IsValid() {
			return convert_to_value(field.Interface())
		}

		panic(fmt.Sprintf("Unknown request method or property: %s", propertyName.value))
	case *Struct:
		propertyName, ok := property.(String)
		if !ok {
			panic("Computed member access must use string expression for property")
		}
		attr, exists := owner._type.storage[propertyName.Value()]
		if !exists {
			panic(fmt.Sprintf("Unknown struct member: %s", propertyName.Value()))
		}

		if !attr.isStatic {
			panic(fmt.Sprintf("Cannot access non-static member '%s' on struct type", propertyName.Value()))
		}

		return attr.Reference
	case StructInstantiation:
		propertyName, ok := property.(String)
		if !ok {
			panic("Computed member access must use string expression for property")
		}
		attr, exists := owner.constructor._type.storage[propertyName.Value()]
		if !exists {
			panic(fmt.Sprintf("Unknown struct member: %s", propertyName.Value()))
		}

		if attr.isStatic {
			panic(fmt.Sprintf("Cannot access static member '%s' on struct instantiation type", propertyName.Value()))
		}

		ref, exists := owner.storage[propertyName.Value()]
		if !exists {
			if ref, ok := attr.Reference.(*FunctionReference); ok {
				if fn, ok := ref.value.(*FunctionValue); ok {
					// Create a new function value with the closure containing self
					newFn := *fn
					newClosure := NewScope(fn.closure)
					newClosure.Declare(NewVariableReference("self", true, owner, owner.constructor.Type()))
					newFn.closure = newClosure
					return NewFunctionReference(ref.identifier, &newFn)
				}
			}
			panic(fmt.Sprintf("Member '%s' not initialized", propertyName.Value()))
		}

		// If the referenced value is a struct and has methods, we need to handle method calls on it
		if loaded := ref.Load(); loaded != nil {
			if structInst, ok := loaded.(StructInstantiation); ok {
				if methodRef, ok := structInst.storage[propertyName.Value()]; ok {
					if fn, ok := methodRef.Load().(*FunctionValue); ok {
						// Create a new function value with the closure containing self
						newFn := *fn
						newClosure := NewScope(fn.closure)
						newClosure.Declare(NewVariableReference("self", true, structInst, structInst.constructor.Type()))
						newFn.closure = newClosure
						return NewFunctionReference(propertyName.Value(), &newFn)
					}
				}
			}
		}

		return ref

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

	property, ok := expectedExpression.Property.(ast.SymbolExpression)
	if !ok {
		panic("Member access must use symbol expression")
	}

	switch owner := ownerValue.(type) {
	case Array:
		if method, exists := owner.methods[property.Value]; exists {
			return method
		}
		panic(fmt.Sprintf("Unknown array method: %s", property.Value))

	case Slice:
		if method, exists := owner.methods[property.Value]; exists {
			return method
		}
		panic(fmt.Sprintf("Unknown slice method: %s", property.Value))

	case Map:
		if method, exists := owner.methods[property.Value]; exists {
			return method
		}

		panic(fmt.Sprintf("Unknown map method: %s", property.Value))

	case String:
		if method, exists := owner.methods[property.Value]; exists {
			return method
		}

		panic(fmt.Sprintf("Unknown server method: %s", property.Value))

	case Server:
		if method, exists := owner.methods[property.Value]; exists {
			return method
		}

		panic(fmt.Sprintf("Unknown server method: %s", property.Value))

	case Request:
		value := reflect.ValueOf(owner)
		field := value.FieldByName(helpers.Capitalize(property.Value))
		if field.IsValid() {
			return convert_to_value(field.Interface())
		}

		panic(fmt.Sprintf("Unknown request method or property: %s", property.Value))

	case Response:
		if method, exists := owner.Methods[property.Value]; exists {
			return method
		}

		value := reflect.ValueOf(owner)
		field := value.FieldByName(helpers.Capitalize(property.Value))
		if field.IsValid() {
			return convert_to_value(field.Interface())
		}

		panic(fmt.Sprintf("Unknown response method or property: %s", property.Value))

	case *Error:
		if method, exists := owner.methods[property.Value]; exists {
			return method
		}

		panic(fmt.Sprintf("Unknown error method: %s", property.Value))

	case *Module:
		if method, exists := owner.exports[property.Value]; exists {
			return method
		}

		panic("Unknown module member")

	case *Struct:
		attr, exists := owner._type.storage[property.Value]
		if !exists {
			panic(fmt.Sprintf("Unknown struct member: %s", property.Value))
		}

		if !attr.isStatic {
			panic(fmt.Sprintf("Cannot access non-static member '%s' on struct type", property.Value))
		}

		return attr.Reference

	case StructInstantiation:
		attr, exists := owner.constructor._type.storage[property.Value]
		if !exists {
			panic(fmt.Sprintf("Unknown struct member: %s", property.Value))
		}

		if attr.isStatic {
			panic(fmt.Sprintf("Cannot access static member '%s' on struct instantiation type", property.Value))
		}

		ref, exists := owner.storage[property.Value]
		if !exists {
			if ref, ok := attr.Reference.(*FunctionReference); ok {
				if fn, ok := ref.value.(*FunctionValue); ok {
					newFn := *fn
					newClosure := NewScope(fn.closure)
					newClosure.Declare(NewVariableReference("self", true, owner, owner.constructor.Type()))
					newFn.closure = newClosure
					return NewFunctionReference(ref.identifier, &newFn)
				}
			}
			panic(fmt.Sprintf("Member '%s' not initialized", property.Value))
		}

		if loaded := ref.Load(); loaded != nil {
			if structInst, ok := loaded.(StructInstantiation); ok {
				if methodRef, ok := structInst.storage[property.Value]; ok {
					if fn, ok := methodRef.Load().(*FunctionValue); ok {
						newFn := *fn
						newClosure := NewScope(fn.closure)
						newClosure.Declare(NewVariableReference("self", true, structInst, structInst.constructor.Type()))
						newFn.closure = newClosure
						return NewFunctionReference(property.Value, &newFn)
					}
				}
			}
		}

		return ref

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

func evaluate_struct_instantiation_expression(expression ast.Expression, scope *Scope) Value {
	expectedExpression, err := ast.ExpectExpression[ast.StructLiteralExpression](expression)
	if err != nil {
		panic(err)
	}

	constructor := evaluate_expression(expectedExpression.Constructor, scope)
	constructorStruct, err := ExpectValue[*Struct](constructor)
	if err != nil {
		panic(err)
	}

	storage := make(map[string]Reference)

	for identifier, property := range constructorStruct._type.storage {
		if !property.isStatic {
			if ref, ok := property.Reference.(*VariableReference); ok {
				var defaultValue Value
				if ref.value != nil {
					defaultValue = ref.value.Clone()
				} else {
					defaultValue = ref.explicitType.DefaultValue()
				}

				storage[identifier] = NewVariableReference(
					identifier,
					ref.isConstant,
					defaultValue,
					ref.explicitType,
				)
			}
		}
	}

	for _, propertyExpression := range expectedExpression.Properties {
		var propertyName string

		if propertyExpression.Identifier != nil {
			identifier, err := ast.ExpectExpression[ast.SymbolExpression](propertyExpression.Identifier)
			if err != nil {
				panic(err)
			}
			propertyName = identifier.Value
		} else {
			for name, attr := range constructorStruct._type.storage {
				if _, ok := attr.Reference.(*FunctionReference); !ok && !attr.isStatic && storage[name] == nil {
					propertyName = name
					break
				}
			}
		}

		structAttr, exists := constructorStruct._type.storage[propertyName]
		if !exists {
			panic(fmt.Sprintf("Property '%s' does not exist in struct", propertyName))
		}

		if structAttr.isStatic {
			panic(fmt.Sprintf("Cannot assign to static property '%s'", propertyName))
		}

		propertyValue := evaluate_expression(propertyExpression.Value, scope)

		if ref, ok := structAttr.Reference.(*VariableReference); ok {
			if !ref.explicitType.Equals(propertyValue.Type()) {
				panic(fmt.Sprintf("Type mismatch: cannot assign value of type %v to property '%s' of type %v",
					propertyValue.Type(), propertyName, ref.explicitType))
			}

			storage[propertyName] = NewVariableReference(
				propertyName,
				ref.isConstant,
				propertyValue,
				ref.explicitType,
			)
		}
	}

	return NewStructInstaniation(*constructorStruct, storage)
}
