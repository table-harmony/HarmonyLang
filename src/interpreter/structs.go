package interpreter

import "fmt"

type StructProperty struct {
	defaultValue Value
	_type        Type
	isStatic     bool
}

type StructMethod struct {
	isStatic bool
	value    FunctionValue
}

type Struct struct {
	properties map[string]StructProperty
	methods    map[string]StructMethod
}

func (s Struct) String() string {
	str := ""

	for identifier, property := range s.properties {
		if property.isStatic {
			str += "static "
		}
		str += identifier + ": " + property._type.String()
		if property.defaultValue != nil {
			str += " = " + property.defaultValue.String()
		}
		str += "\n"
	}

	for identifier, method := range s.methods {
		if method.isStatic {
			str += "static "
		}
		str += identifier + ": " + method.value.String()
		str += "\n"
	}

	return str
}
func (Struct) DefaultValue() Value { return NewNil() }
func (s Struct) Equals(other Type) bool {
	otherStruct, ok := other.(Struct)
	if !ok {
		return false
	}

	for identifier, property := range s.properties {
		otherProperty, exists := otherStruct.properties[identifier]
		if !exists {
			return false
		}

		if !otherProperty._type.Equals(property._type) {
			return false
		}
	}

	for identifier, method := range s.methods {
		match := false
		for otherIdentifier, otherMethod := range otherStruct.methods {
			if otherMethod.value.Type().Equals(method.value.Type()) && otherIdentifier == identifier {
				match = true
				break
			}
		}

		if !match {
			return false
		}
	}

	return true
}
func NewStruct(properties map[string]StructProperty, methods map[string]StructMethod) Struct {
	return Struct{
		properties, methods,
	}
}

type StructReference struct {
	identifier string
	_type      Struct
}

func NewStructReference(identifier string, _type Struct) *StructReference {
	return &StructReference{
		identifier, _type,
	}
}

// StructReference implements the Value interface
func (s *StructReference) Type() Type { return s._type }
func (s *StructReference) Clone() Value {
	properties := make(map[string]StructProperty, 0)
	for key, property := range s._type.properties {
		properties[key] = property
	}

	methods := make(map[string]StructMethod, 0)
	for identifier, method := range s._type.methods {
		ptr := method.value.Clone()
		methods[identifier] = StructMethod{isStatic: method.isStatic, value: ptr.(FunctionValue)}
	}

	_type := NewStruct(properties, methods)
	return NewStructReference(s.identifier, _type)
}
func (s *StructReference) String() string {
	str := s.identifier + "{\n"
	str += s._type.String()
	str += "}\n"
	return str
}

// StructReference implements the Reference interface
func (s *StructReference) Load() Value { return s }
func (s *StructReference) Store(v Value) error {
	return fmt.Errorf("cannot assign onto struct reference")
}
func (s *StructReference) Address() Value {
	return NewPointer(s)
}
