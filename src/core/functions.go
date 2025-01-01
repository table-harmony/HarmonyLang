package core

import (
	"fmt"
	"strings"

	"github.com/table-harmony/HarmonyLang/src/ast"
)

type FunctionType struct {
	parameters []ParameterType
	returnType Type
}

type ParameterType struct {
	identifier string
	valueType  Type
}

// FunctionType implements the Type interface
func (f FunctionType) String() string {
	params := make([]string, len(f.parameters))
	for i, param := range f.parameters {
		if param.identifier != "" {
			params[i] = fmt.Sprintf("%s: %s", param.identifier, param.valueType.String())
		} else {
			params[i] = param.valueType.String()
		}
	}
	return fmt.Sprintf("fn(%s) -> %s", strings.Join(params, ", "), f.returnType.String())
}
func (f FunctionType) Equals(other Type) bool {
	otherFn, ok := other.(FunctionType)
	if !ok {
		return false
	}
	if len(f.parameters) != len(otherFn.parameters) {
		return false
	}
	if !f.returnType.Equals(otherFn.returnType) {
		return false
	}
	for i := range f.parameters {
		if !f.parameters[i].valueType.Equals(otherFn.parameters[i].valueType) {
			return false
		}
	}
	return true
}

// TODO: implementation for default values for functions
func (f FunctionType) DefaultValue() Value {
	return Nil{}
}

type FunctionValue struct {
	parameters []ast.Parameter
	body       []ast.Statement
	returnType Type
	closure    *Scope
}

func NewFunctionValue(params []ast.Parameter, body []ast.Statement, returnType Type, closure *Scope) *FunctionValue {
	return &FunctionValue{
		parameters: params,
		body:       body,
		returnType: returnType,
		closure:    closure,
	}
}

// FunctionValue implements the Value interface
func (f FunctionValue) Type() Type {
	params := make([]ParameterType, len(f.parameters))
	for i, param := range f.parameters {
		params[i] = ParameterType{
			identifier: param.Name,
			valueType:  EvaluateType(param.Type),
		}
	}

	return FunctionType{
		parameters: params,
		returnType: f.returnType,
	}
}
func (f FunctionValue) Clone() Value {
	paramsCopy := make([]ast.Parameter, len(f.parameters))
	copy(paramsCopy, f.parameters)

	bodyCopy := make([]ast.Statement, len(f.body))
	copy(bodyCopy, f.body)

	return FunctionValue{
		parameters: paramsCopy,
		body:       bodyCopy,
		returnType: f.returnType,
		closure:    f.closure,
	}
}
func (f FunctionValue) String() string { return "function" }
func (f FunctionValue) CreateScope(params []Value) (*Scope, error) {
	functionScope := NewScope(f.closure)

	if len(params) != len(f.parameters) {
		return nil, fmt.Errorf("expected %d arguments but got %d",
			len(f.parameters), len(params))
	}

	for i, param := range f.parameters {
		paramType := EvaluateType(param.Type)
		paramValue := params[i]

		if paramValue.Type().Equals(paramType) && paramType.Equals(PrimitiveType{AnyType}) {
			return nil, fmt.Errorf("parameter '%s' expected type '%s' but got '%s'",
				param.Name, paramType.String(), paramValue.Type())
		}

		paramRef := &VariableReference{param.Name, false, paramValue, paramType}
		functionScope.Declare(paramRef)
	}

	return functionScope, nil
}
func (f FunctionValue) Body() []ast.Statement { return f.body }

type FunctionReference struct {
	identifier string
	value      FunctionValue
}

func NewFunctionReference(identifier string, value FunctionValue) *FunctionReference {
	return &FunctionReference{
		identifier: identifier,
		value:      value,
	}
}

// FunctionReference implements the Value interface
func (f *FunctionReference) Type() Type   { return f.value.Type() }
func (f *FunctionReference) Clone() Value { return f.value.Clone() }
func (f *FunctionReference) String() string {
	return fmt.Sprintf("type: function, identifier: %s", f.identifier)
}

// FunctionReference implements the Reference interface
func (f *FunctionReference) Load() Value { return f.value }
func (f *FunctionReference) Store(v Value) error {
	fn, ok := v.(FunctionValue)
	if !ok {
		return fmt.Errorf("cannot assign non-function value to function reference '%s'", f.identifier)
	}

	fnType := fn.Type().(FunctionType)
	expectedType := f.value.Type().(FunctionType)

	if !fnType.Equals(expectedType) {
		return fmt.Errorf("type mismatch: cannot assign function of type %v to function '%s' of type %v",
			fnType.String(), f.identifier, expectedType.String())
	}

	f.value = fn
	return nil
}
func (f *FunctionReference) Address() Value {
	return NewPointer(f)
}
