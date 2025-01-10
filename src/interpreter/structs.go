package interpreter

import (
	"fmt"
)

type StructAttribute struct {
	Reference
	isStatic bool
}

type StructType struct {
	storage map[string]StructAttribute
}

func NewStructType(storage map[string]StructAttribute) StructType {
	for name, attr := range storage {
		if attr.Reference == nil {
			panic(fmt.Sprintf("struct attribute '%s' has nil reference", name))
		}

		attrType := attr.Reference.Type()
		if attrType == nil {
			panic(fmt.Sprintf("struct attribute '%s' has nil type", name))
		}

		if _, ok := attrType.(FunctionType); ok {
			if attr.Reference.Load() == nil {
				panic(fmt.Sprintf("struct method '%s' has nil value but type %v", name, attrType))
			}
		}
	}

	return StructType{storage}
}

func (s StructType) String() string {
	str := ""
	for identifier, item := range s.storage {
		if item.isStatic {
			str += "static "
		}
		str += identifier + ": " + item.Load().String() + "\n"
	}
	return str
}
func (s StructType) DefaultValue() Value    { return NewNil() }
func (s StructType) Equals(other Type) bool { return true } //TODO:....

type Struct struct {
	identifier string
	_type      StructType
}

func NewStruct(identifier string, _type StructType) *Struct {
	return &Struct{
		identifier: identifier,
		_type:      _type,
	}
}

// Struct implements the Value interface
func (s *Struct) Type() Type { return s._type }
func (s *Struct) Clone() Value {
	return NewStruct(s.identifier, s._type)
}
func (s *Struct) String() string {
	str := fmt.Sprintf("struct %s {\n", s.identifier)
	for name, attr := range s._type.storage {
		str += "  "
		if attr.isStatic {
			str += "static "
		}

		str += fmt.Sprintf("%s: %s\n", name, attr.Reference.Type().String())
	}
	str += "}"
	return str
}

// Struct implements the Reference interface
func (s *Struct) Load() Value { return s }
func (s *Struct) Store(v Value) error {
	return fmt.Errorf("cannot assign to struct type %s", s.identifier)
}
func (s *Struct) Address() Value { return NewPointer(s) }

type StructInstantiation struct {
	constructor Struct
	storage     map[string]Reference
}

func NewStructInstaniation(constructor Struct, storage map[string]Reference) StructInstantiation {
	return StructInstantiation{
		constructor: constructor,
		storage:     storage,
	}
}

// StructInstantiation implements the Value interface
func (s StructInstantiation) Type() Type { return s.constructor._type }
func (s StructInstantiation) Clone() Value {
	newStorage := make(map[string]Reference)
	for name, ref := range s.storage {
		newStorage[name] = NewVariableReference(
			name,
			false,
			ref.Load().Clone(),
			ref.Type(),
		)
	}
	return NewStructInstaniation(s.constructor, newStorage)
}
func (s StructInstantiation) String() string {
	str := fmt.Sprintf("%s {\n", s.constructor.identifier)
	for name, ref := range s.storage {
		str += fmt.Sprintf("  %s: %v\n", name, ref.Load())
	}
	str += "}"
	return str
}
