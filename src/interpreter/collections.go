package interpreter

import (
	"fmt"
)

type ArrayType struct {
	size        int
	elementType Type
}

func NewArrayType(size Value, elementType Type) *ArrayType {
	validSize, ok := size.(Number)
	if !ok {
		panic("Array size must be a number")
	}

	if !validSize.IsInteger() {
		panic("Array size must be an integer")
	}

	if validSize.value < 0 {
		panic("Array size must be greater than or equal to 0")
	}

	length := int(validSize.value)
	return &ArrayType{length, elementType}
}

// ArrayType implements the Type interface
func (a ArrayType) String() string      { return fmt.Sprintf("[%d]%s", a.size, a.elementType) }
func (a ArrayType) DefaultValue() Value { return NewNil() }
func (a ArrayType) Equals(other Type) bool {
	if other == nil {
		return true
	}
	if primitive, ok := other.(PrimitiveType); ok {
		return primitive.kind == NilType
	}

	otherArray, ok := other.(ArrayType)
	if !ok {
		return false
	}

	if otherArray.size != a.size || !a.elementType.Equals(otherArray.elementType) {
		return false
	}

	return true
}

type Array struct {
	elements []Value
	_type    ArrayType
}

func NewEmptyArray(size Value, elementType Type) *Array {
	_type := NewArrayType(size, elementType)
	defaultValue := elementType.DefaultValue()

	elements := make([]Value, 0, _type.size)
	for i := len(elements); i < _type.size; i++ {
		elements = append(elements, defaultValue)
	}

	return &Array{elements, *_type}
}

func NewArray(elements []Value, size Value, elementType Type) *Array {
	_type := NewArrayType(size, elementType)
	if _type.size < len(elements) {
		panic(fmt.Sprintf("Array size is less than the number of elements, expected %d", _type.size))
	}

	defaultValue := elementType.DefaultValue()
	for i := len(elements); i < _type.size; i++ {
		elements = append(elements, defaultValue)
	}

	for _, element := range elements {
		if !elementType.Equals(element.Type()) {
			panic(fmt.Sprintf("Array type is not compatible with element type %s, expected %s", element.Type().String(), elementType.String()))
		}
	}

	return &Array{elements, *_type}
}

// Array implements the Value interface
func (a *Array) Type() Type { return a._type }
func (a *Array) Clone() Value {
	return NewArray(a.elements, NewNumber(float64(a._type.size)), a._type.elementType)
}
func (a *Array) String() string {
	str := a._type.String() + "["
	for i, element := range a.elements {
		str += element.String()

		if i < len(a.elements)-1 {
			str += ", "
		}
	}
	str += "]"
	return str
}

// Array specific methods
func (a *Array) Get(property Value) Value {
	position, ok := property.(Number)
	if !ok {
		panic(fmt.Sprintf("expected index to be a number, but got %T", property))
	}

	if !position.IsInteger() {
		panic(fmt.Sprintf("expected index to be an integer, but got %v", property))
	}

	var index int = int(position.value)
	if index < 0 {
		index = len(a.elements) + index
	}

	if index >= len(a.elements) {
		panic(fmt.Sprintf("index out of range %v with length %v", index, len(a.elements)))
	}

	return a.elements[index]
}
func (a *Array) Set(property Value, newValue Value) {
	index, ok := property.(Number)
	if !ok {
		panic(fmt.Sprintf("expected index to be a number, but got %T", property))
	}
	if !index.IsInteger() {
		panic(fmt.Sprintf("expected index to be an integer, but got %v", index))
	}

	arrayIndex := int(index.value)
	if arrayIndex < 0 {
		arrayIndex = len(a.elements) + arrayIndex
	}

	if arrayIndex >= len(a.elements) {
		panic(fmt.Sprintf("index out of range %v with length %v", arrayIndex, len(a.elements)))
	}

	if !a._type.elementType.Equals(newValue.Type()) {
		panic(fmt.Sprintf("cannot assign value of type %s to array of type %s",
			a.Type().String(), a._type.elementType.String()))
	}

	a.elements[arrayIndex] = newValue
}

type MapType struct {
	keyType   Type
	valueType Type
}

func NewMapType(keyType Type, valueType Type) *MapType {
	if keyType == nil {
		keyType = NewNil().Type()
	}
	if valueType == nil {
		valueType = NewNil().Type()
	}

	return &MapType{keyType, valueType}
}

// MapType implements the Type interface
func (m MapType) String() string      { return fmt.Sprintf("map[%s -> %s]", m.keyType, m.valueType) }
func (m MapType) DefaultValue() Value { return NewNil() }
func (m MapType) Equals(other Type) bool {
	if other == nil {
		return true
	}
	if primitive, ok := other.(PrimitiveType); ok {
		return primitive.kind == NilType
	}

	otherMap, ok := other.(MapType)
	if !ok {
		return false
	}

	if !m.keyType.Equals(otherMap.keyType) || !m.valueType.Equals(otherMap.keyType) {
		return false
	}

	return true
}

type MapEntry struct {
	key   Value
	value Value
}

type Map struct {
	entries []MapEntry
	_type   MapType
	methods map[string]NativeFunction
}

func NewMap(entries []MapEntry, keyType Type, valueType Type) *Map {
	_type := NewMapType(keyType, valueType)

	for _, entry := range entries {
		if !entry.key.Type().Equals(keyType) || !entry.value.Type().Equals(valueType) {
			panic(fmt.Sprintf("Map entry type is not compatible with key type %s or value type %s, expected key type %s and value type %s", entry.key.Type().String(), entry.value.Type().String(), keyType.String(), valueType.String()))
		}
	}

	m := &Map{
		entries: entries,
		_type:   *_type,
		methods: make(map[string]NativeFunction),
	}
	m.init_methods()

	return m
}

// Map implements the Value interface
func (m *Map) Type() Type { return m._type }
func (m *Map) Clone() Value {
	return NewMap(m.entries, m._type.keyType, m._type.valueType)
}
func (m *Map) String() string {
	str := m._type.String() + "{"
	for i, entry := range m.entries {
		str += entry.key.String() + " -> " + entry.value.String()

		if i < len(m.entries)-1 {
			str += ", \n"
		}
	}
	str += "}"
	return str
}

// Map specific methods
func (m *Map) Get(key Value) Value {
	for _, entry := range m.entries {
		if entry.key == key {
			return entry.value
		}
	}
	panic(fmt.Sprintf("key %s not found in map", key.String()))
}
func (m *Map) Set(key Value, newValue Value) {
	if !m._type.keyType.Equals(key.Type()) {
		panic(fmt.Sprintf("cannot use key of type %s for map with key type %s",
			m.Type().String(), m._type.keyType.String()))
	}
	if !m._type.valueType.Equals(newValue.Type()) {
		panic(fmt.Sprintf("cannot assign value of type %s to map with value type %s",
			newValue.Type().String(), m._type.valueType.String()))
	}

	for i, entry := range m.entries {
		if entry.key == key {
			m.entries[i].value = newValue
			return
		}
	}

	m.entries = append(m.entries, MapEntry{key, newValue})
}
func (m *Map) Keys() []Value {
	keys := make([]Value, 0)

	for _, entry := range m.entries {
		keys = append(keys, entry.key)
	}

	return keys
}
func (m *Map) Values() []Value {
	values := make([]Value, 0)
	for _, entry := range m.entries {
		values = append(values, entry.value)
	}
	return values
}
func (m *Map) IsExist(key Value) bool {
	for _, entry := range m.entries {
		if entry.key == key {
			return true
		}
	}
	return false
}
func (m *Map) Pop(key Value) Value {
	index := -1

	for i, entry := range m.entries {
		if entry.key == key {
			index = i
			break
		}
	}

	if index != -1 {
		m.entries = append(m.entries[:index], m.entries[index+1:]...)
		return NewBoolean(true)
	}

	return NewBoolean(false)
}

func (m *Map) init_methods() {
	getFunc := func(args ...Value) Value {
		if len(args) != 1 {
			panic("Get method expects exactly one argument")
		}
		return m.Get(args[0])
	}

	setFunc := func(args ...Value) Value {
		if len(args) != 2 {
			panic("Set method expects exactly 2 argument")
		}

		m.Set(args[0], args[1])
		return NewNil()
	}

	popFunc := func(args ...Value) Value {
		if len(args) != 1 {
			panic("Pop method expects exactly 1 argument")
		}

		return m.Pop(args[0])
	}

	m.methods["get"] = *NewNativeFunction(
		getFunc,
		[]Type{m._type.keyType},
		m._type.valueType,
	)

	m.methods["set"] = *NewNativeFunction(
		setFunc,
		[]Type{m._type.keyType, m._type.valueType},
		PrimitiveType{NilType},
	)

	m.methods["pop"] = *NewNativeFunction(
		popFunc,
		[]Type{m._type.keyType},
		PrimitiveType{BooleanType},
	)
}
