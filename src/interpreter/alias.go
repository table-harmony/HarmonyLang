package interpreter

type ValueType struct {
	_type Type
}

func NewValueType(_type Type) ValueType {
	return ValueType{_type}
}

// ValueType implements the Value interface
func (t ValueType) Type() Type     { return t._type }
func (t ValueType) Clone() Value   { return NewValueType(t._type) }
func (t ValueType) String() string { return t._type.String() }

// ValueType implements the Type interface
func (t ValueType) Equals(other Type) bool { return t._type.Equals(other) }
func (t ValueType) DefaultValue() Value    { return t._type.DefaultValue() }
