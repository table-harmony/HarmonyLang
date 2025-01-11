package interpreter

import (
	"fmt"
)

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

	if !m.keyType.Equals(otherMap.keyType) || !m.valueType.Equals(otherMap.valueType) {
		return false
	}

	return true
}

type MapEntry struct {
	key   Value
	value Value
}

type Map struct {
	entries *[]MapEntry
	_type   MapType
	methods map[string]NativeFunctionValue
}

func NewMap(entries []MapEntry, keyType Type, valueType Type) Map {
	_type := NewMapType(keyType, valueType)

	for _, entry := range entries {
		if !keyType.Equals(entry.key.Type()) || !valueType.Equals(entry.value.Type()) {
			panic(fmt.Sprintf("Map entry type is not compatible with key type %s or value type %s, expected key type %s and value type %s", entry.key.Type().String(), entry.value.Type().String(), keyType.String(), valueType.String()))
		}
	}

	m := Map{
		entries: &entries,
		_type:   *_type,
		methods: make(map[string]NativeFunctionValue),
	}
	m.init_methods()

	return m
}

// Map implements the Value interface
func (m Map) Type() Type { return m._type }
func (m Map) Clone() Value {
	copyEntries := make([]MapEntry, len(*m.entries))
	copy(copyEntries, *m.entries)
	return NewMap(copyEntries, m._type.keyType, m._type.valueType)
}
func (m Map) String() string {
	str := m._type.String() + "{\n"
	for i, entry := range *m.entries {
		str += "  " + entry.key.String() + " -> " + entry.value.String()

		if i < len(*m.entries)-1 {
			str += ", \n"
		}
	}
	str += "\n}"
	return str
}

// Map specific methods
func (m *Map) Get(key Value) Value {
	for _, entry := range *(m.entries) {
		if entry.key.String() == key.String() {
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

	for i, entry := range *m.entries {
		if entry.key == key {
			(*m.entries)[i].value = newValue
			return
		}
	}
	*(m.entries) = append(*m.entries, MapEntry{key, newValue})
}
func (m *Map) Keys() []Value {
	keys := make([]Value, 0)

	for _, entry := range *(m.entries) {
		keys = append(keys, entry.key)
	}

	return keys
}
func (m *Map) Values() []Value {
	values := make([]Value, 0)
	for _, entry := range *(m.entries) {
		values = append(values, entry.value)
	}
	return values
}
func (m *Map) IsExist(key Value) bool {
	for _, entry := range *(m.entries) {
		if entry.key == key {
			return true
		}
	}
	return false
}
func (m *Map) Pop(key Value) Value {
	index := -1

	for i, entry := range *(m.entries) {
		if entry.key == key {
			index = i
			break
		}
	}

	if index != -1 {
		*m.entries = append((*m.entries)[:index], (*m.entries)[index+1:]...)
		return NewBoolean(true)
	}

	return NewBoolean(false)
}
func (m *Map) Intersect(other Value) Value {
	otherMap, ok := other.(Map)
	if !ok {
		panic("Intersect method expects a map as the argument")
	}

	if !m._type.Equals(otherMap._type) {
		panic(fmt.Sprintf("Map types mismatch: expected map of type %s, but got %s", m._type.String(), otherMap._type.String()))
	}

	newMap := NewMap(make([]MapEntry, 0), m._type.keyType, m._type.valueType)
	for _, entry := range *m.entries {
		if otherMap.IsExist(entry.key) {
			*(newMap.entries) = append(*newMap.entries, MapEntry{entry.key.Clone(), entry.value.Clone()})
		}
	}

	return newMap
}
func (m *Map) Union(other Value) Value {
	otherMap, ok := other.(Map)
	if !ok {
		panic("Union method expects a map as the argument")
	}

	if !m._type.Equals(otherMap._type) {
		panic(fmt.Sprintf("Map types mismatch: expected map of type %s, but got %s", m._type.String(), otherMap._type.String()))
	}

	newMap := NewMap(make([]MapEntry, 0), m._type.keyType, m._type.valueType)
	for _, entry := range *m.entries {
		*(newMap.entries) = append(*newMap.entries, MapEntry{entry.key.Clone(), entry.value.Clone()})
	}

	for _, entry := range *otherMap.entries {
		if !newMap.IsExist(entry.key) {
			*(newMap.entries) = append(*newMap.entries, MapEntry{entry.key.Clone(), entry.value.Clone()})
		}
	}

	return newMap
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

	m.methods["exists"] = *NewNativeFunction(
		func(args ...Value) Value {
			if len(args) != 1 {
				panic("Exists method expects exactly 1 argument")
			}
			return NewBoolean(m.IsExist(args[0]))
		},
		[]Type{m._type.keyType},
		PrimitiveType{BooleanType},
	)

	m.methods["intersect"] = *NewNativeFunction(
		func(args ...Value) Value {
			if len(args) != 1 {
				panic("Intersect method expects exactly 1 argument")
			}
			return m.Intersect(args[0])
		},
		[]Type{m._type},
		m._type,
	)

	m.methods["union"] = *NewNativeFunction(
		func(args ...Value) Value {
			if len(args) != 1 {
				panic("Union method expects exactly 1 argument")
			}
			return m.Union(args[0])
		},
		[]Type{m._type},
		m._type,
	)

	m.methods["keys"] = *NewNativeFunction(
		func(args ...Value) Value {
			if len(args) != 1 {
				panic("Union method expects exactly 1 argument")
			}
			return m.Union(args[0])
		},
		[]Type{m._type},
		m._type,
	)

	m.methods["values"] = *NewNativeFunction(
		func(args ...Value) Value {
			if len(args) != 0 {
				panic("Values method expects exactly 0 arguments")
			}
			values := m.Values()
			return NewArray(values, NewNumber(float64(len(*(m.entries)))), m._type.valueType)
		},
		[]Type{},
		NewArrayType(NewNumber(float64(len(*(m.entries)))), m._type.valueType),
	)

	m.methods["keys"] = *NewNativeFunction(
		func(args ...Value) Value {
			if len(args) != 0 {
				panic("Keys method expects exactly 0 arguments")
			}
			values := m.Keys()
			return NewArray(values, NewNumber(float64(len(*(m.entries)))), m._type.keyType)
		},
		[]Type{},
		NewArrayType(NewNumber(float64(len(*(m.entries)))), m._type.keyType),
	)
}
