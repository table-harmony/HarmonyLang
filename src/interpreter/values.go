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
func (n Number) String() string  { return fmt.Sprintf("%g", n.value) }
func (n Number) Value() float64  { return n.value }

// String implementation
func (s String) Type() ValueType { return StringType }
func (s String) Clone() Value    { return String{s.value} }
func (s String) String() string  { return s.value }
func (s String) Value() string   { return s.value }

// Boolean implementation
func (b Boolean) Type() ValueType { return BooleanType }
func (b Boolean) Clone() Value    { return Boolean{b.value} }
func (b Boolean) String() string  { return fmt.Sprintf("%t", b.value) }
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

// Implement Value interface
func (s *VariableReference) Type() ValueType { return s.value.Type() }
func (s *VariableReference) Clone() Value    { return s.value.Clone() }
func (s *VariableReference) String() string  { return s.value.String() }

// Implement Reference interface
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

type ReferenceValue struct {
	value Value
}

func NewReference(value Value) *ReferenceValue {
	return &ReferenceValue{value}
}

// Implement Value interface
func (r *ReferenceValue) Type() ValueType { return ReferenceType }
func (r *ReferenceValue) Clone() Value    { return NewReference(r.value.Clone()) }
func (r *ReferenceValue) String() string  { return r.value.String() }

// Implement Reference interface
func (r *ReferenceValue) Load() Value { return r.value }
func (r *ReferenceValue) Store(v Value) error {
	r.value = v
	return nil
}
func (r *ReferenceValue) Address() Value {
	return NewPointer(r)
}

type Pointer struct {
	target Reference
}

func NewPointer(target Reference) *Pointer {
	return &Pointer{target}
}

// Implement Value interface
func (p *Pointer) Type() ValueType { return PointerType }
func (p *Pointer) Clone() Value    { return NewPointer(p.target) }
func (p *Pointer) String() string {
	if p.target == nil {
		return "nil"
	}
	return fmt.Sprintf("&%v", p.target.String())
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
