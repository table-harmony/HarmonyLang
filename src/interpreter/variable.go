package interpreter

import (
	"fmt"
)

type VariableReference struct {
	identifier   string
	isConstant   bool
	value        Value
	explicitType Type
}

func NewVariableReference(identifier string, isConstant bool, value Value, explicitType Type) *VariableReference {
	if explicitType == nil {
		explicitType = PrimitiveType{AnyType}
	}
	if value == nil {
		value = explicitType.DefaultValue()
	}

	value = value.Clone()
	variable := VariableReference{
		identifier,
		isConstant,
		value,
		explicitType,
	}

	if !explicitType.Equals(value.Type()) && !explicitType.Equals(PrimitiveType{AnyType}) {
		panic(fmt.Sprintf("variable '%s' expected type '%s' but got '%s'",
			variable.identifier, explicitType.String(), value.Type().String()))
	}

	return &variable
}

// VariableReference implements the Value interface
func (s *VariableReference) Type() Type     { return s.explicitType }
func (s *VariableReference) Clone() Value   { return s.value.Clone() }
func (s *VariableReference) String() string { return s.value.String() }

// VariableReference implements the Reference interface
func (s *VariableReference) Load() Value { return s.value }
func (s *VariableReference) Store(v Value) error {
	if s.isConstant {
		return fmt.Errorf("cannot assign to constant variable '%s'", s.identifier)
	}

	if !s.explicitType.Equals(v.Type()) && !s.explicitType.Equals(PrimitiveType{AnyType}) {
		return fmt.Errorf("type mismatch: cannot assign %v to %s of type %v",
			v.Type(), s.identifier, s.explicitType)
	}

	s.value = v
	return nil
}
func (s *VariableReference) Address() Value { return NewPointer(s) }
