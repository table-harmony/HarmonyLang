package interpreter

import (
	"fmt"

	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/helpers"
)

// Tthe core interface that all runtime values must implement
type Value interface {
	Type() Type     // Type returns the type of the value
	Clone() Value   // Clone creates a deep copy of the value
	String() string // String provides a string representation for debugging/REPL
}

// Reference represents a reference to a storage location
type Reference interface {
	Value
	Load() Value       // Load retrieves the current value from the storage location
	Store(Value) error // Store updates the value at the storage location
	Address() Value    // Address returns a pointer to this reference
}

// Type represents a runtime type
type Type interface {
	String() string      // String provides a string representation for debugging/REPL
	Equals(Type) bool    // Equals checks if two types are equal
	DefaultValue() Value // DefaultValue returns the default value for the type
}

func ExpectValue[T Value](value Value) (T, error) {
	return helpers.ExpectType[T](value)
}

func ExpectReference[T Reference](value Reference) (T, error) {
	return helpers.ExpectType[T](value)
}

func ExpectType[T Type](value Type) (T, error) {
	return helpers.ExpectType[T](value)
}

// EvaluateType evaluates an AST type into a runtime type
func EvaluateType(astType ast.Type, scope *Scope) Type {
	if astType == nil {
		return PrimitiveType{AnyType}
	}

	switch t := astType.(type) {
	case ast.StringType:
		return PrimitiveType{StringType}
	case ast.NumberType:
		return PrimitiveType{NumberType}
	case ast.BooleanType:
		return PrimitiveType{BooleanType}
	case ast.NilType:
		return PrimitiveType{NilType}
	case ast.FunctionType:
		params := make([]ParameterType, len(t.Parameters))
		for i, param := range t.Parameters {
			params[i] = ParameterType{
				identifier: param.Name,
				valueType:  EvaluateType(param.Type, scope),
			}
		}
		return FunctionType{
			parameters: params,
			returnType: EvaluateType(t.Return, scope),
		}
	case ast.ArrayType:
		size := evaluate_expression(t.Size, scope)
		return NewArrayType(size, EvaluateType(t.Underlying, scope))
	case ast.SliceType:
		return NewSliceType(EvaluateType(t.Underlying, scope))
	case ast.PointerType:
		value := EvaluateType(t.Target, scope)
		return NewPointerType(value)
	case ast.SymbolType:
		ref, err := scope.Resolve(t.Value)
		if err != nil {
			panic(err)
		}
		return ref.Load().Type()
	case ast.AnyType:
		return PrimitiveType{AnyType}
	default:
		panic(fmt.Sprintf("invalid type: %T", t))
	}
}
