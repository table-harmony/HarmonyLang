package interpreter

import "fmt"

type PointerType struct {
	valueType Type
}

func NewPointerType(valueType Type) *PointerType {
	return &PointerType{valueType}
}

// PointerType implements the Type interface
func (p PointerType) String() string {
	return fmt.Sprintf("*%s", p.valueType.String())
}
func (p PointerType) Equals(other Type) bool {
	if other == nil {
		return true
	}
	if primitive, ok := other.(PrimitiveType); ok {
		return primitive.kind == NilType
	}

	otherPtr, ok := other.(PointerType)
	if !ok {
		return false
	}
	return p.valueType.Equals(otherPtr.valueType)
}
func (p PointerType) DefaultValue() Value {
	return Nil{}
}

type Pointer struct {
	target Reference
}

func NewPointer(target Reference) *Pointer {
	return &Pointer{target}
}

// Pointer implements the Value interface
func (p *Pointer) Type() Type   { return PointerType{valueType: p.target.Type()} }
func (p *Pointer) Clone() Value { return NewPointer(p.target) }
func (p *Pointer) String() string {
	if p.target == nil {
		return "nil"
	}
	return fmt.Sprintf("&%v", p.target.String())
}

// Deref returns the reference pointed to by the pointer
func (p *Pointer) Deref() Reference {
	return p.target
}

// Deref returns the value pointed to by the pointer
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
