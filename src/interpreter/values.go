package interpreter

import (
	"fmt"

	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/helpers"
)

type ValueType int

const (
	NumberType ValueType = iota
	StringType
	BooleanType
	NilType
	ArrayType
	SliceType
	MapType
	StructType
	InterfaceType
	FunctionType
	PointerType
	ReferenceType
	AnyType
)

// TODO: evaluate type from ast to a value type
func evaluate_type(t ast.Type) ValueType {
	return AnyType
}

func ExpectValue[T Value](value Value) (T, error) {
	return helpers.ExpectType[T](value)
}

func ExpectReference[T Reference](value Reference) (T, error) {
	return helpers.ExpectType[T](value)
}

func (t ValueType) String() string {
	switch t {
	case NumberType:
		return "number"
	case StringType:
		return "string"
	case BooleanType:
		return "boolean"
	case NilType:
		return "nil"
	case ArrayType:
		return "array"
	case SliceType:
		return "slice"
	case MapType:
		return "map"
	case StructType:
		return "struct"
	case InterfaceType:
		return "interface"
	case FunctionType:
		return "function"
	case PointerType:
		return "pointer"
	case ReferenceType:
		return "reference"
	default:
		return "unknown"
	}
}

type Value interface {
	Type() ValueType
	Clone() Value
	String() string
}

type Reference interface {
	Value
	Load() Value
	Store(Value) error
	Address() Value
}

// Primitive type implementations
type Number struct{ value float64 }
type String struct{ value string }
type Boolean struct{ value bool }
type Nil struct{}

// Number implementation
func (n Number) Type() ValueType { return NumberType }
func (n Number) Clone() Value    { return Number{n.value} }
func (n Number) String() string  { return fmt.Sprintf("type: number, value: %g", n.value) }
func (n Number) Value() float64  { return n.value }

// String implementation
func (s String) Type() ValueType { return StringType }
func (s String) Clone() Value    { return String{s.value} }
func (s String) String() string  { return fmt.Sprintf("type: string, value: %s", s.value) }
func (s String) Value() string   { return s.value }

// Boolean implementation
func (b Boolean) Type() ValueType { return BooleanType }
func (b Boolean) Clone() Value    { return Boolean{b.value} }
func (b Boolean) String() string  { return fmt.Sprintf("type: boolean, value: %t", b.value) }
func (b Boolean) Value() bool     { return b.value }

// Nil implementation
func (Nil) Type() ValueType { return NilType }
func (Nil) Clone() Value    { return Nil{} }
func (Nil) String() string  { return "nil" }

type VariableReference struct {
	identifier   string
	isConstant   bool
	value        Value
	explicitType ValueType
}

func (s *VariableReference) Type() ValueType { return s.value.Type() }
func (s *VariableReference) Clone() Value    { return s.value.Clone() }
func (s *VariableReference) String() string  { return s.value.String() }

func (s *VariableReference) Load() Value { return s.value }
func (s *VariableReference) Store(v Value) error {
	if s.isConstant {
		return fmt.Errorf("cannot assign to constant '%s'", s.identifier)
	}
	if s.explicitType != v.Type() && s.explicitType != AnyType {
		return fmt.Errorf("type mismatch: cannot assign %v to %s of type %v",
			v.Type(), s.identifier, s.explicitType)
	}
	s.value = v
	return nil
}
func (s *VariableReference) Address() Value {
	return NewPointer(s)
}

type Pointer struct {
	target Reference
}

func NewPointer(target Reference) *Pointer {
	return &Pointer{target}
}

func (p *Pointer) Type() ValueType { return PointerType }
func (p *Pointer) Clone() Value    { return NewPointer(p.target) }
func (p *Pointer) String() string {
	if p.target == nil {
		return "nil"
	}
	return fmt.Sprintf("&{ %v }", p.target.String())
}

func (p *Pointer) Deref() Reference {
	return p.target
}

func Deref(v Value) (Value, error) {
	switch ptr := v.(type) {
	case *Pointer:
		if ptr.target == nil {
			return nil, fmt.Errorf("null pointer dereference")
		}
		return ptr.target.Load(), nil
	default:
		return nil, fmt.Errorf("cannot dereference non-pointer type %v", v.Type())
	}
}

type FunctionValue struct {
	parameters []ast.Parameter
	body       []ast.Statement
	returnType ValueType
	closure    *Scope
}

func (f FunctionValue) Type() ValueType { return FunctionType }
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

func (f FunctionValue) Call(params []Value, scope *Scope) (result Value) {
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

	functionScope := NewScope(f.closure)

	if len(params) != len(f.parameters) {
		panic(fmt.Errorf("expected %d arguments but got %d",
			len(f.parameters), len(params)))
	}

	for i, param := range f.parameters {
		paramType := evaluate_type(param.Type)
		paramValue := params[i]

		if paramValue.Type() != paramType && paramType != AnyType {
			panic(fmt.Sprintf("parameter '%s' expected type '%s' but got '%s'",
				param.Name, paramType.String(), paramValue.Type().String()))
		}

		paramRef := &VariableReference{param.Name, false, paramValue, paramType}
		functionScope.Declare(paramRef)
	}

	for _, statement := range f.body {
		evaluate_statement(statement, functionScope)
	}

	return Nil{}
}

type FunctionReference struct {
	identifier string
	value      FunctionValue
}

func (f *FunctionReference) Type() ValueType { return FunctionType }
func (f *FunctionReference) Clone() Value    { return f.value.Clone() }
func (f *FunctionReference) String() string {
	return fmt.Sprintf("type: function, identifier: %s", f.identifier)
}

func (f *FunctionReference) Load() Value { return f.value }

func (f *FunctionReference) Store(v Value) error {
	fn, ok := v.(FunctionValue)
	if !ok {
		return fmt.Errorf("cannot assign non-function value to function reference '%s'", f.identifier)
	}

	if fn.returnType != f.value.returnType {
		return fmt.Errorf("type mismatch: cannot assign function with return type %v to function '%s' expecting return type %v",
			fn.returnType, f.identifier, f.value.returnType)
	}

	if len(fn.parameters) != len(f.value.parameters) {
		return fmt.Errorf("parameter count mismatch: function '%s' expects %d parameters but got %d",
			f.identifier, len(f.value.parameters), len(fn.parameters))
	}

	f.value = fn
	return nil
}

func (f *FunctionReference) Address() Value {
	return NewPointer(f)
}
