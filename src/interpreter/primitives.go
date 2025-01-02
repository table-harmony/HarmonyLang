package interpreter

import "fmt"

type PrimitiveType struct {
	kind PrimitiveKind
}

type PrimitiveKind int

const (
	NumberType PrimitiveKind = iota
	StringType
	BooleanType
	NilType
	AnyType
)

// PrimitiveType implements Type interface
func (p PrimitiveType) DefaultValue() Value {
	switch p.kind {
	case NumberType:
		return NewNumber(0)
	case StringType:
		return NewString("")
	case BooleanType:
		return NewBoolean(false)
	case NilType:
		return NewNil()
	case AnyType:
		return NewNil()
	default:
		panic("unknown primitive type")
	}
}
func (p PrimitiveType) String() string {
	switch p.kind {
	case NumberType:
		return "number"
	case StringType:
		return "string"
	case BooleanType:
		return "boolean"
	case NilType:
		return "nil"
	case AnyType:
		return "any"
	default:
		return "unknown"
	}
}
func (p PrimitiveType) Equals(other Type) bool {
	if p.kind == AnyType {
		return true
	}
	if otherPrim, ok := other.(PrimitiveType); ok {
		return p.kind == otherPrim.kind
	}
	return false
}

// primitive values
type Number struct{ value float64 }
type String struct{ value string }
type Boolean struct{ value bool }
type Nil struct{}

// Number implements Value interface
func (n Number) Type() Type         { return PrimitiveType{NumberType} }
func (n Number) Clone() Value       { return NewNumber(n.value) }
func (n Number) String() string     { return fmt.Sprintf("%g", n.value) }
func (n Number) Value() float64     { return n.value }
func (n Number) IsInteger() bool    { return n.value == float64(int(n.value)) }
func NewNumber(value float64) Value { return Number{value} }

// Number implements Value interface
func (s String) Type() Type        { return PrimitiveType{StringType} }
func (s String) Clone() Value      { return NewString(s.value) }
func (s String) String() string    { return fmt.Sprintf("\"%s\"", s.value) }
func (s String) Value() string     { return s.value }
func NewString(value string) Value { return String{value} }

// Boolean implements Value interface
func (b Boolean) Type() Type      { return PrimitiveType{BooleanType} }
func (b Boolean) Clone() Value    { return NewBoolean(b.value) }
func (b Boolean) String() string  { return fmt.Sprintf("%t", b.value) }
func (b Boolean) Value() bool     { return b.value }
func NewBoolean(value bool) Value { return Boolean{value} }

// Nil implements Value interface
func (Nil) Type() Type     { return PrimitiveType{NilType} }
func (Nil) Clone() Value   { return Nil{} }
func (Nil) String() string { return "nil" }
func NewNil() Value        { return Nil{} }
