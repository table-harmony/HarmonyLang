package interpreter

import (
	"fmt"
	"reflect"

	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/lexer"
)

func (interpreter *interpreter) evalute_current_statement(scope *Scope) {
	statement := interpreter.current_statement()
	evaluate_statement(statement, scope)
}

func evaluate_statement(statement ast.Statement, scope *Scope) {
	if statement == nil {
		return
	}

	statementType := reflect.TypeOf(statement)
	if handler, exists := statement_lookup[statementType]; exists {
		handler(statement, scope)
	} else {
		panic(fmt.Sprintf("No handler registered for statement type: %v", statementType))
	}
}

func evaluate_expression_statement(statement ast.Statement, scope *Scope) {
	expectedStatement, err := ast.ExpectStatement[ast.ExpressionStatement](statement)
	if err != nil {
		panic(err)
	}

	evaluate_expression(expectedStatement.Expression, scope)
}

func evaluate_variable_declaration_statement(statement ast.Statement, scope *Scope) {
	expectedStatement, err := ast.ExpectStatement[ast.VariableDeclarationStatement](statement)

	if err != nil {
		panic(err)
	}

	value := evaluate_expression(expectedStatement.Value, scope)

	var _type Type = value.Type()
	if expectedStatement.ExplicitType != nil {
		_type = EvaluateType(expectedStatement.ExplicitType, scope)
	}

	variable := NewVariableReference(
		expectedStatement.Identifier,
		expectedStatement.IsConstant,
		value,
		_type,
	)

	err = scope.Declare(variable)
	if err != nil {
		panic(err)
	}
}

func evaluate_multi_variable_declaration_statement(statement ast.Statement, scope *Scope) {
	expectedStatement, err := ast.ExpectStatement[ast.MultiVariableDeclarationStatement](statement)
	if err != nil {
		panic(err)
	}

	for _, declaration := range expectedStatement.Declarations {
		evaluate_variable_declaration_statement(declaration, scope)
	}
}

func evaluate_continue_statement(statement ast.Statement, scope *Scope) {
	panic(ContinueError{})
}

func evaluate_break_statement(statement ast.Statement, scope *Scope) {
	panic(BreakError{})
}

func evaluate_return_statement(statement ast.Statement, scope *Scope) {
	expectedStatement, err := ast.ExpectStatement[ast.ReturnStatement](statement)
	if err != nil {
		panic(err)
	}

	panic(NewReturnError(evaluate_expression(expectedStatement.Value, scope)))
}

func evaluate_traditional_for_statement(statement ast.Statement, scope *Scope) {
	expectedStatement, err := ast.ExpectStatement[ast.TraditionalForStatement](statement)
	if err != nil {
		panic(err)
	}

	loopScope := NewScope(scope)
	evaluate_statement(expectedStatement.Initializer, loopScope)

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(BreakError); ok {
				return
			}
			panic(r)
		}
	}()

	for {
		condition := evaluate_expression(expectedStatement.Condition, loopScope)
		if condition == nil {
			condition = NewBoolean(true)
		}
		conditionValue, err := ExpectValue[Boolean](condition)
		if err != nil {
			panic(err)
		}

		if !conditionValue.Value() {
			return
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

			for _, statement := range expectedStatement.Body {
				evaluate_statement(statement, loopScope)
			}
		}()

		for _, statement := range expectedStatement.Post {
			evaluate_statement(statement, loopScope)
		}
	}
}

func evaluate_iterator_for_statement(statement ast.Statement, scope *Scope) {
	expectedStatement, err := ast.ExpectStatement[ast.IteratorForStatement](statement)
	if err != nil {
		panic(err)
	}

	loopScope := NewScope(scope)
	iteratorValue := evaluate_expression(expectedStatement.Iterator, loopScope)

	var keyType, valueType Type
	iterations := 0
	switch iterator := iteratorValue.(type) {
	case Array:
		keyType = PrimitiveType{NumberType}
		valueType = iterator._type.elementType
		iterations = len(iterator.elements)
	case Slice:
		keyType = PrimitiveType{NumberType}
		valueType = iterator._type.elementType
		iterations = iterator.length
	case Map:
		keyType = iterator._type.keyType
		valueType = iterator._type.valueType
		iterations = len(*iterator.entries)
	}

	key := NewVariableReference(
		expectedStatement.KeyIdentifier,
		false,
		keyType.DefaultValue(),
		keyType,
	)

	err = loopScope.Declare(key)
	if err != nil {
		panic(err)
	}

	var value *VariableReference
	if expectedStatement.ValueIdentifier != "" {
		value = NewVariableReference(
			expectedStatement.ValueIdentifier,
			false,
			valueType.DefaultValue(),
			valueType,
		)

		err = loopScope.Declare(value)
		if err != nil {
			panic(err)
		}
	}

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(BreakError); ok {
				return
			}
			panic(r)
		}
	}()

	for i := 0; i < iterations; i++ {
		switch iterator := iteratorValue.(type) {
		case Array:
			keyValue := NewNumber(float64(i))
			key.Store(keyValue)

			if value != nil {
				value.Store(iterator.elements[i])
			}
		case Slice:
			keyValue := NewNumber(float64(i))
			key.Store(keyValue)

			if value != nil {
				value.Store((*iterator.elements)[i])
			}
		case Map:
			current := (*iterator.entries)[i]
			keyValue := current.key
			key.Store(keyValue)
			if value != nil {
				value.Store(current.value)
			}
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

			for _, statement := range expectedStatement.Body {
				evaluate_statement(statement, loopScope)
			}
		}()
	}
}

func evaluate_function_declaration_statement(statement ast.Statement, scope *Scope) {
	expectedStatement, err := ast.ExpectStatement[ast.FunctionDeclarationStatment](statement)
	if err != nil {
		panic(err)
	}

	valuePtr := NewFunctionValue(
		expectedStatement.Parameters,
		expectedStatement.Body,
		EvaluateType(expectedStatement.ReturnType, scope),
		scope,
	)

	ref := NewFunctionReference(
		expectedStatement.Identifier,
		*valuePtr,
	)

	err = scope.Declare(ref)

	if err != nil {
		panic(err)
	}
}

func evaluate_assignment_statement(statement ast.Statement, scope *Scope) {
	expectedStatement, err := ast.ExpectStatement[ast.AssignmentStatement](statement)
	if err != nil {
		panic(err)
	}

	value := evaluate_expression(expectedStatement.Value, scope)
	if arrayValue, ok := value.(Array); ok {
		value = arrayValue.Clone()
	}
	switch assigne := expectedStatement.Assigne.(type) {
	case ast.SymbolExpression:
		ref, err := scope.Resolve(assigne.Value)
		if err != nil {
			panic(fmt.Sprintf("cannot assign to undefined variable %s", assigne.Value))
		}
		if err := ref.Store(value); err != nil {
			panic(err)
		}

	case ast.PrefixExpression:
		if assigne.Operator.Kind != lexer.STAR {
			panic("invalid assignment target")
		}
		ptrValue := evaluate_expression(assigne.Right, scope)
		ptr, err := ExpectValue[*Pointer](ptrValue)
		if err != nil {
			panic("cannot dereference non-pointer type")
		}
		if err := ptr.Deref().Store(value); err != nil {
			panic(err)
		}

	case ast.ComputedMemberExpression:
		ownerValue := evaluate_expression(assigne.Owner, scope)
		property := evaluate_expression(assigne.Property, scope)
		if ref, ok := ownerValue.(Reference); ok {
			ownerValue = ref.Load()
		}
		if ref, ok := property.(Reference); ok {
			property = ref.Load()
		}

		switch owner := ownerValue.(type) {
		case Array:
			owner.Set(property, value)
		case Map:
			owner.Set(property, value)
		case Slice:
			owner.Set(property, value)
		default:
			panic(fmt.Sprintf("cannot index into value of type %T", owner))
		}

	case ast.MemberExpression:
		panic("TODO: member expression in assignment statement evaluation")

	default:
		panic("invalid assignment target")
	}
}

func evaluate_throw_statement(statement ast.Statement, scope *Scope) {
	expectedStatement, err := ast.ExpectStatement[ast.ThrowStatement](statement)
	if err != nil {
		panic(err)
	}

	value := evaluate_expression(expectedStatement.Value, scope)
	panic(value)
}

func evaluate_type_declaration_statement(statement ast.Statement, scope *Scope) {
	expectedStatement, err := ast.ExpectStatement[ast.TypeDeclarationStatement](statement)
	if err != nil {
		panic(err)
	}

	valueType := NewValueType(EvaluateType(expectedStatement.Type, scope))
	variable := NewVariableReference(
		expectedStatement.Identifier,
		true,
		valueType,
		valueType,
	)

	err = scope.Declare(variable)
	if err != nil {
		panic(err)
	}
}
