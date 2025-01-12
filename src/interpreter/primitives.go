package interpreter

import (
	"fmt"
	"net/url"
	"strings"
)

type PrimitiveType struct {
	kind PrimitiveKind
}

type PrimitiveKind int

const (
	NumberType PrimitiveKind = iota
	StringType
	BooleanType
	ErrorType
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
type String struct {
	value   string
	methods map[string]Function
}
type Boolean struct{ value bool }
type Error struct {
	value   string
	methods map[string]Function
}
type Nil struct{}

// Number implements Value interface
func (n Number) Type() Type         { return PrimitiveType{NumberType} }
func (n Number) Clone() Value       { return NewNumber(n.value) }
func (n Number) String() string     { return fmt.Sprintf("%g", n.value) }
func (n Number) Value() float64     { return n.value }
func (n Number) IsInteger() bool    { return n.value == float64(int(n.value)) }
func NewNumber(value float64) Value { return Number{value} }

// Number implements Value interface
func (s String) Type() Type     { return PrimitiveType{StringType} }
func (s String) Clone() Value   { return NewString(s.value) }
func (s String) String() string { return s.value }
func (s String) Value() string  { return s.value }
func NewString(value string) Value {
	methods := make(map[string]Function)

	s := String{value, methods}

	methods["len"] = NewNativeFunction(
		func(args ...Value) Value {
			return NewNumber(float64(len(s.value)))
		},
		[]Type{},
		PrimitiveType{NumberType},
	)

	methods["index"] = NewNativeFunction(
		func(args ...Value) Value {
			substr := args[0].(String)
			return NewNumber(float64(strings.Index(s.value, substr.value)))
		},
		[]Type{PrimitiveType{StringType}},
		PrimitiveType{NumberType},
	)

	methods["contains"] = NewNativeFunction(
		func(args ...Value) Value {
			substr := args[0].(String)
			return NewBoolean(strings.Contains(s.value, substr.value))
		},
		[]Type{PrimitiveType{StringType}},
		PrimitiveType{BooleanType},
	)

	methods["substr"] = NewNativeFunction(
		func(args ...Value) Value {
			start := args[0].(Number).value
			end := args[1].(Number).value
			if start < 0 || int(end) > len(s.value) || start > end {
				panic("substring indices out of range")
			}
			return NewString(s.value[int(start):int(end)])
		},
		[]Type{PrimitiveType{NumberType}, PrimitiveType{NumberType}},
		PrimitiveType{StringType},
	)

	methods["upper"] = NewNativeFunction(
		func(args ...Value) Value {
			return NewString(strings.ToUpper(s.value))
		},
		[]Type{},
		PrimitiveType{StringType},
	)

	methods["lower"] = NewNativeFunction(
		func(args ...Value) Value {
			return NewString(strings.ToLower(s.value))
		},
		[]Type{},
		PrimitiveType{StringType},
	)

	methods["trim"] = NewNativeFunction(
		func(args ...Value) Value {
			return NewString(strings.TrimSpace(s.value))
		},
		[]Type{},
		PrimitiveType{StringType},
	)

	methods["split"] = NewNativeFunction(
		func(args ...Value) Value {
			sep := args[0].(String).value
			parts := strings.Split(s.value, sep)
			result := make([]Value, len(parts))
			for i, part := range parts {
				result[i] = NewString(part)
			}
			return NewSlice(result, PrimitiveType{StringType})
		},
		[]Type{PrimitiveType{StringType}},
		NewSliceType(PrimitiveType{StringType}),
	)

	methods["url_decode"] = NewNativeFunction(
		func(args ...Value) Value {
			decoded, _ := url.QueryUnescape(s.value)
			return NewString(decoded)
		},
		[]Type{},
		PrimitiveType{StringType},
	)

	methods["url_encode"] = NewNativeFunction(
		func(args ...Value) Value {
			decoded := url.QueryEscape(s.value)
			return NewString(decoded)
		},
		[]Type{},
		PrimitiveType{StringType},
	)

	return s
}

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

// Error implements Value interface
func (e Error) Type() Type     { return PrimitiveType{ErrorType} }
func (e Error) Clone() Value   { return NewError(e.value) }
func (e Error) String() string { return fmt.Sprintf("Error: %s", e.value) }
func NewError(value string) Value {
	err := &Error{value, map[string]Function{}}
	err.init_methods()

	return err
}

func (e *Error) init_methods() {
	e.methods["message"] = NewNativeFunction(func(args ...Value) Value {
		return NewString(e.value)
	}, []Type{}, PrimitiveType{StringType})
}
