package interpreter

import (
	"fmt"
	"strings"
)

type StructProperty struct {
	defaultValue Value
	_type        Type
	isStatic     bool
}

type StructMethod struct {
	isStatic bool
	value    FunctionValue
}

type StructType struct {
	properties    map[string]StructProperty
	methods       map[string]StructMethod
	staticMethods map[string]Value
	staticValues  map[string]Value
}

func NewStructType(properties map[string]StructProperty, methods map[string]StructMethod) StructType {
	staticMethods := make(map[string]Value)
	staticValues := make(map[string]Value)

	for name, method := range methods {
		if method.isStatic {
			staticMethods[name] = NewPointer(NewFunctionReference(name, method.value))
			delete(methods, name)
		}
	}

	for name, prop := range properties {
		if prop.isStatic {
			staticValues[name] = NewPointer(NewVariableReference(name, false, prop.defaultValue, prop._type))
			delete(properties, name)
		}
	}

	return StructType{
		properties:    properties,
		methods:       methods,
		staticMethods: staticMethods,
		staticValues:  staticValues,
	}
}

func (s StructType) String() string {
	var str strings.Builder

	for name, value := range s.staticValues {
		str.WriteString(fmt.Sprintf("static %s: %s = %s\n",
			name,
			value.Type().String(),
			value.String()))
	}

	for name, prop := range s.properties {
		str.WriteString(fmt.Sprintf("%s: %s",
			name,
			prop._type.String()))
		if prop.defaultValue != nil {
			str.WriteString(fmt.Sprintf(" = %s", prop.defaultValue.String()))
		}
		str.WriteString("\n")
	}

	for name, method := range s.staticMethods {
		str.WriteString(fmt.Sprintf("static %s: %s\n",
			name,
			method.Type().String()))
	}

	for name, method := range s.methods {
		str.WriteString(fmt.Sprintf("%s: %s\n",
			name,
			method.value.Type().String()))
	}

	return str.String()
}
func (s StructType) DefaultValue() Value    { return NewNil() }
func (s StructType) Equals(other Type) bool { return false }

type StructReference struct {
	identifier string
	_type      StructType
}

func NewStructReference(identifier string, _type StructType) Reference {
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

	_type := NewStructType(properties, methods)
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

type StructInstance struct {
	_type      StructType
	properties map[string]Value
	methods    map[string]Function
}

func NewStructInstance(structType StructType) *StructInstance {
	instance := &StructInstance{
		_type:      structType,
		properties: make(map[string]Value),
		methods:    make(map[string]Function),
	}

	for name, prop := range structType.properties {
		if !prop.isStatic {
			instance.properties[name] = prop.defaultValue.Clone()
		}
	}

	for name, method := range structType.methods {
		if !method.isStatic {
			instance.methods[name] = method.value.Clone().(Function)
		}
	}

	return instance
}

// StructInstance implements Value interface
func (si *StructInstance) Type() Type { return si._type }
func (si *StructInstance) Clone() Value {
	clone := NewStructInstance(si._type)
	for name, value := range si.properties {
		clone.properties[name] = value.Clone()
	}
	return clone
}
func (si *StructInstance) String() string {
	var str string
	str += "instance {\n"
	for name, value := range si.properties {
		str += fmt.Sprintf("  %s: %s\n", name, value.String())
	}
	str += "}"
	return str
}

// Add methods for property access
func (si *StructInstance) GetProperty(name string) Value {
	if value, exists := si.properties[name]; exists {
		return value
	}
	if prop, exists := si._type.properties[name]; exists && prop.isStatic {
		return prop.defaultValue
	}
	panic(fmt.Sprintf("property '%s' not found", name))
}

func (si *StructInstance) SetProperty(name string, value Value) error {
	if prop, exists := si._type.properties[name]; exists {
		if prop.isStatic {
			return fmt.Errorf("cannot modify static property '%s'", name)
		}
		if !prop._type.Equals(value.Type()) {
			return fmt.Errorf("type mismatch for property '%s': expected %s, got %s",
				name, prop._type.String(), value.Type().String())
		}
		si.properties[name] = value.Clone()
		return nil
	}
	return fmt.Errorf("property '%s' not found", name)
}

// Add method for method invocation
func (si *StructInstance) InvokeMethod(name string, args ...Value) (Value, error) {
	if method, exists := si.methods[name]; exists {
		return method.Call(args...)
	}
	if method, exists := si._type.methods[name]; exists && method.isStatic {
		return method.value.Call(args...)
	}
	return nil, fmt.Errorf("method '%s' not found", name)
}
