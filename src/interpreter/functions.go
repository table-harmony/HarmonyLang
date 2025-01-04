package interpreter

import (
	"fmt"

	"github.com/table-harmony/HarmonyLang/src/ast"
)

type Function interface {
	Value
	Call(args ...Value) (Value, error)
}

type FunctionType struct {
	parameters []ParameterType
	returnType Type
}

func NewFunctionType(params []ParameterType, returnType Type) FunctionType {
	return FunctionType{
		parameters: params,
		returnType: returnType,
	}
}

type ParameterType struct {
	identifier string
	valueType  Type
}

// FunctionType implements the Type interface
func (f FunctionType) String() string {
	str := "fn("

	for i, param := range f.parameters {
		if i > 0 {
			str += ", "
		}
		str += param.identifier + ": " + param.valueType.String()
	}

	str += ") -> " + f.returnType.String()
	return str
}
func (f FunctionType) Equals(other Type) bool {
	if other == nil {
		return true
	}
	if primitive, ok := other.(PrimitiveType); ok {
		return primitive.kind == NilType
	}

	otherFn, ok := other.(FunctionType)
	if !ok {
		return false
	}

	if len(f.parameters) != len(otherFn.parameters) {
		return false
	}

	for i := range f.parameters {
		if !f.parameters[i].valueType.Equals(otherFn.parameters[i].valueType) {
			return false
		}
	}

	if _, ok := f.returnType.(PrimitiveType); ok && f.returnType.(PrimitiveType).kind == AnyType {
		return true
	}

	return f.returnType.Equals(otherFn.returnType)
}
func (f FunctionType) DefaultValue() Value {
	return NewNil()
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
			valueType:  EvaluateType(param.Type, f.closure),
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
func (f FunctionValue) String() string {
	str := "fn("

	for i, param := range f.parameters {
		if i > 0 {
			str += ", "
		}
		str += param.Name + ": " + EvaluateType(param.Type, f.closure).String()
	}

	str += ") -> " + f.returnType.String()
	return str
}
func (f FunctionValue) Call(args ...Value) (result Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch e := r.(type) {
			case ReturnError:
				result = e.Value()
				err = nil
			case error:
				err = e
			default:
				panic(e)
			}
		}
	}()

	functionScope := NewScope(f.closure)
	if len(args) > len(f.parameters) {
		return nil, fmt.Errorf("expected at most %d arguments but got %d",
			len(f.parameters), len(args))
	}

	for i, param := range f.parameters {
		var paramValue Value
		if i < len(args) {
			paramValue = args[i]
		} else if param.DefaultValue != nil {
			paramValue = evaluate_expression(param.DefaultValue, functionScope)
		} else {
			return nil, fmt.Errorf("missing value for parameter '%s'", param.Name)
		}

		paramType := EvaluateType(param.Type, f.closure)
		if !paramValue.Type().Equals(paramType) && !paramType.Equals(PrimitiveType{AnyType}) {
			return nil, fmt.Errorf("parameter '%s' expected type '%s' but got '%s'",
				param.Name, paramType.String(), paramValue.Type())
		}

		paramRef := NewVariableReference(param.Name, false, paramValue, paramType)
		functionScope.Declare(paramRef)
	}

	for _, statement := range f.body {
		evaluate_statement(statement, functionScope)
	}

	return NewNil(), nil
}

type FunctionReference struct {
	identifier string
	value      Function
}

func NewFunctionReference(identifier string, value Function) *FunctionReference {
	return &FunctionReference{
		identifier: identifier,
		value:      value,
	}
}

// FunctionReference implements the Value interface
func (f *FunctionReference) Type() Type     { return f.value.Type() }
func (f *FunctionReference) Clone() Value   { return f.value.Clone() }
func (f *FunctionReference) String() string { return f.value.String() }

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

// NativeFunction represents the type signature for native Go functions that can be called
type NativeFunction func(args ...Value) Value

type NativeFunctionType struct {
	paramTypes []Type
	returnType Type
}

// NativeFunctionType implements the Type interface
func (f NativeFunctionType) String() string {
	return "native_fn"
}
func (f NativeFunctionType) DefaultValue() Value {
	return NewNil()
}
func (f NativeFunctionType) Equals(other Type) bool {
	if other == nil {
		return true
	}
	if primitive, ok := other.(PrimitiveType); ok {
		return primitive.kind == NilType || primitive.kind == AnyType
	}

	otherFn, ok := other.(NativeFunctionType)
	if !ok {
		return false
	}

	if len(f.paramTypes) != len(otherFn.paramTypes) {
		return false
	}

	for i := range f.paramTypes {
		if !f.paramTypes[i].Equals(otherFn.paramTypes[i]) {
			return false
		}
	}

	return f.returnType.Equals(otherFn.returnType)
}

type NativeFunctionValue struct {
	value      NativeFunction
	paramTypes []Type
	returnType Type
}

func NewNativeFunction(fn NativeFunction, paramTypes []Type, returnType Type) *NativeFunctionValue {
	return &NativeFunctionValue{
		value:      fn,
		paramTypes: paramTypes,
		returnType: returnType,
	}
}

// NativeFunction implements the Value interface
func (n NativeFunctionValue) Type() Type {
	return NativeFunctionType{
		paramTypes: n.paramTypes,
		returnType: n.returnType,
	}
}
func (n NativeFunctionValue) Clone() Value {
	return NewNativeFunction(n.value, n.paramTypes, n.returnType)
}
func (n NativeFunctionValue) String() string {
	str := "native_fn("

	for i, param := range n.paramTypes {
		if i > 0 {
			str += ", "
		}
		str += param.String()
	}

	str += ") -> " + n.returnType.String()
	return str
}

func (n NativeFunctionValue) Call(args ...Value) (Value, error) {
	if len(args) != len(n.paramTypes) {
		return NewNil(), fmt.Errorf("expected %d arguments but got %d", len(n.paramTypes), len(args))
	}
	for i, arg := range args {
		if !n.paramTypes[i].Equals(arg.Type()) {
			return NewNil(), fmt.Errorf("argument %d: expected %v but got %v", i, n.paramTypes[i], arg.Type())
		}
	}

	result := n.value(args...)

	if !n.returnType.Equals(result.Type()) {
		return NewNil(), fmt.Errorf("return value: expected %v but got %v", n.returnType, result.Type())
	}

	return result, nil
}
