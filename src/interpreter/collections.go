package interpreter

import (
	"fmt"
	"strconv"
)

type ArrayType struct {
	size        int
	elementType Type
}

func NewArrayType(size Value, elementType Type) ArrayType {
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
	return ArrayType{length, elementType}
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

	return otherArray.size != a.size || !a.elementType.Equals(otherArray.elementType)
}

type Array struct {
	elements []Value
	_type    ArrayType
	methods  map[string]NativeFunctionValue
}

func NewEmptyArray(size Value, elementType Type) Array {
	_type := NewArrayType(size, elementType)
	defaultValue := elementType.DefaultValue()

	elements := make([]Value, 0, _type.size)
	for i := len(elements); i < _type.size; i++ {
		elements = append(elements, defaultValue)
	}

	arr := Array{
		elements: elements,
		_type:    _type,
		methods:  make(map[string]NativeFunctionValue),
	}
	arr.init_methods()
	return arr
}

func NewArray(elements []Value, size Value, elementType Type) Array {
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

	arr := Array{
		elements: elements,
		_type:    _type,
		methods:  make(map[string]NativeFunctionValue),
	}
	arr.init_methods()
	return arr
}

func (a *Array) init_methods() {
	lengthFunc := func(args ...Value) Value {
		if len(args) != 0 {
			panic("Length method expects no arguments")
		}
		return NewNumber(float64(len(a.elements)))
	}
	a.methods["len"] = *NewNativeFunction(
		lengthFunc,
		[]Type{},
		PrimitiveType{NumberType},
	)

	getFunc := func(args ...Value) Value {
		if len(args) != 1 {
			panic("Get method expects exactly one argument")
		}
		return a.Get(args[0])
	}
	a.methods["get"] = *NewNativeFunction(
		getFunc,
		[]Type{PrimitiveType{NumberType}},
		a._type.elementType,
	)

	setFunc := func(args ...Value) Value {
		if len(args) != 2 {
			panic("Set method expects exactly two arguments")
		}
		a.Set(args[0], args[1])
		return NewNil()
	}
	a.methods["set"] = *NewNativeFunction(
		setFunc,
		[]Type{PrimitiveType{NumberType}, a._type.elementType},
		PrimitiveType{NilType},
	)

	a.methods["slice"] = *NewNativeFunction(
		func(args ...Value) Value {
			if len(args) != 2 {
				panic("Slice method expects exactly two arguments")
			}
			return a.Slice(args[0], args[1])
		},
		[]Type{PrimitiveType{NumberType}, PrimitiveType{NumberType}},
		NewSliceType(a._type.elementType),
	)
}

// Array implements the Value interface
func (a Array) Type() Type { return a._type }
func (a Array) Clone() Value {
	newElements := make([]Value, len(a.elements))
	copy(newElements, a.elements)
	return NewArray(newElements, NewNumber(float64(a._type.size)), a._type.elementType)
}
func (a Array) String() string {
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

	if index < 0 || index >= len(a.elements) {
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

	if arrayIndex >= len(a.elements) || arrayIndex < 0 {
		panic(fmt.Sprintf("index out of range %v with length %v", arrayIndex, len(a.elements)))
	}

	if !a._type.elementType.Equals(newValue.Type()) {
		panic(fmt.Sprintf("cannot assign value of type %s to array of type %s",
			a.Type().String(), a._type.elementType.String()))
	}

	a.elements[arrayIndex] = newValue
}
func (a *Array) Slice(start Value, end Value) Slice {
	startValue, ok := start.(Number)
	if !ok {
		panic("Invalid start index must be a number")
	}

	endValue, ok := end.(Number)
	if !ok {
		panic("Invalid end index must be a number")
	}

	if !startValue.IsInteger() || !endValue.IsInteger() {
		panic("Invalid indices must be integers")
	}

	startIndex := int(startValue.value)
	endIndex := int(endValue.value)

	if startIndex < 0 {
		startIndex = a._type.size + startIndex
	}
	if endIndex < 0 {
		endIndex = a._type.size + endIndex
	}

	if startIndex < 0 || endIndex > a._type.size || startIndex > endIndex {
		panic(fmt.Sprintf("Invalid array indices [%d:%d] with length %d",
			startIndex, endIndex, a._type.size))
	}

	elements := a.elements[startIndex:endIndex]
	newSlice := Slice{
		elements: &elements,
		_type:    *NewSliceType(a._type.elementType),
		length:   endIndex - startIndex,
		capacity: a._type.size - startIndex,
		methods:  make(map[string]NativeFunctionValue),
	}
	newSlice.init_methods()
	return newSlice
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
		if !entry.key.Type().Equals(keyType) || !entry.value.Type().Equals(valueType) {
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
				panic("exists methods accepts only one argument for key")
			}
			return NewBoolean(m.IsExist(args[0]))
		},
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
}

type SliceType struct {
	elementType Type
}

func NewSliceType(elementType Type) *SliceType {
	return &SliceType{elementType}
}

// SliceType implements the Type interface
func (s SliceType) String() string      { return fmt.Sprintf("[]%s", s.elementType) }
func (s SliceType) DefaultValue() Value { return NewSlice(nil, s.elementType) }
func (s SliceType) Equals(other Type) bool {
	if other == nil {
		return true
	}
	if primitive, ok := other.(PrimitiveType); ok {
		return primitive.kind == NilType
	}

	otherSlice, ok := other.(SliceType)
	if !ok {
		return false
	}
	return s.elementType.Equals(otherSlice.elementType)
}

type Slice struct {
	elements *[]Value
	_type    SliceType
	length   int
	capacity int
	methods  map[string]NativeFunctionValue
}

const INITIAL_CAPACITY = 8

func NewSlice(elements []Value, elementType Type) Slice {
	var capacity int
	if elements != nil {
		capacity = max(INITIAL_CAPACITY, len(elements)*2)
	}

	sliceElements := make([]Value, 0, capacity)
	slice := Slice{
		elements: &sliceElements,
		_type:    *NewSliceType(elementType),
		length:   0,
		capacity: capacity,
		methods:  make(map[string]NativeFunctionValue),
	}

	for _, element := range elements {
		if !elementType.Equals(element.Type()) {
			panic(fmt.Sprintf("Slice element type mismatch: expected %s, got %s",
				elementType.String(), element.Type().String()))
		}

		*(slice.elements) = append(*(slice.elements), element)
		slice.length++
	}

	slice.init_methods()
	return slice
}

func (s Slice) Type() Type { return s._type }
func (s Slice) Clone() Value {
	copyElements := make([]Value, len(*s.elements))
	copy(copyElements, *s.elements)
	return NewSlice(copyElements, s._type.elementType)
}
func (s Slice) String() string {
	str := s._type.String() + "["
	for i, element := range *s.elements {
		str += element.String()

		if i < len(*s.elements)-1 {
			str += ", "
		}
	}
	str += "] { len: " + strconv.Itoa(s.length) + ", cap: " + strconv.Itoa(s.capacity) + " }"
	return str
}

// Slice specific methods
func (s *Slice) Append(value Value) {
	if !s._type.elementType.Equals(value.Type()) {
		panic(fmt.Sprintf("Cannot append %s to slice of %s",
			value.Type().String(), s._type.elementType.String()))
	}

	if s.length == s.capacity {
		newCap := s.capacity * 2
		newElements := make([]Value, s.length, newCap)
		copy(newElements, *s.elements)
		s.elements = &newElements
		s.capacity = newCap
	}

	*(s.elements) = append(*s.elements, value)
	s.length++
}
func (s *Slice) Get(property Value) Value {
	position, ok := property.(Number)
	if !ok {
		panic(fmt.Sprintf("expected index to be a number, but got %T", property))
	}

	if !position.IsInteger() {
		panic(fmt.Sprintf("expected index to be an integer, but got %v", property))
	}

	var index int = int(position.value)
	if index < 0 {
		index = len(*s.elements) + index
	}

	if index >= len(*s.elements) {
		panic(fmt.Sprintf("index out of range %v with length %v", index, len(*s.elements)))
	}

	return (*s.elements)[index]
}
func (s *Slice) Set(property Value, value Value) {
	position, ok := property.(Number)
	if !ok {
		panic(fmt.Sprintf("expected index to be a number, but got %T", property))
	}

	if !position.IsInteger() {
		panic(fmt.Sprintf("expected index to be an integer, but got %v", property))
	}

	var index int = int(position.value)
	if index < 0 {
		index = len(*s.elements) + index
	}

	if index >= s.length {
		panic(fmt.Sprintf("Index out of range [%d] with length %d", index, s.length))
	}
	if !s._type.elementType.Equals(value.Type()) {
		panic(fmt.Sprintf("Cannot set %s in slice of %s",
			value.Type().String(), s._type.elementType.String()))
	}
	(*s.elements)[index] = value
}
func (s *Slice) Slice(start Value, end Value) Slice {
	startValue, ok := start.(Number)
	if !ok {
		panic("Invalid start index must be a number")
	}

	endValue, ok := end.(Number)
	if !ok {
		panic("Invalid end index must be a number")
	}

	if !startValue.IsInteger() || !endValue.IsInteger() {
		panic("Invalid indices must be integers")
	}

	startIndex := int(startValue.value)
	endIndex := int(endValue.value)

	if startIndex < 0 {
		startIndex = s.length + startIndex
	}
	if endIndex < 0 {
		endIndex = s.length + endIndex
	}

	if startIndex < 0 || endIndex > s.length || startIndex > endIndex {
		panic(fmt.Sprintf("Invalid slice indices [%d:%d] with length %d",
			startIndex, endIndex, s.length))
	}

	elements := (*s.elements)[startIndex:endIndex]
	newSlice := Slice{
		elements: &elements,
		_type:    s._type,
		length:   endIndex - startIndex,
		capacity: s.capacity - startIndex,
		methods:  make(map[string]NativeFunctionValue),
	}
	newSlice.init_methods()
	return newSlice
}

func (s *Slice) init_methods() {
	s.methods["len"] = *NewNativeFunction(
		func(args ...Value) Value { return NewNumber(float64(s.length)) },
		[]Type{},
		PrimitiveType{NumberType},
	)

	s.methods["cap"] = *NewNativeFunction(
		func(args ...Value) Value { return NewNumber(float64(s.capacity)) },
		[]Type{},
		PrimitiveType{NumberType},
	)

	s.methods["append"] = *NewNativeFunction(
		func(args ...Value) Value {
			if len(args) != 1 {
				panic("Append method expects exactly one argument")
			}
			s.Append(args[0])
			return NewNil()
		},
		[]Type{s._type.elementType},
		PrimitiveType{NilType},
	)

	s.methods["get"] = *NewNativeFunction(
		func(args ...Value) Value {
			if len(args) != 1 {
				panic("Get method expects exactly one argument")
			}
			return s.Get(args[0])
		},
		[]Type{PrimitiveType{NumberType}},
		s._type.elementType,
	)

	s.methods["set"] = *NewNativeFunction(
		func(args ...Value) Value {
			if len(args) != 2 {
				panic("Set method expects exactly two arguments")
			}
			s.Set(args[0], args[1])
			return NewNil()
		},
		[]Type{PrimitiveType{NumberType}, PrimitiveType{NumberType}},
		PrimitiveType{NilType},
	)

	s.methods["slice"] = *NewNativeFunction(
		func(args ...Value) Value {
			if len(args) != 2 {
				panic("Slice method expects exactly two arguments")
			}
			return s.Slice(args[0], args[1])
		},
		[]Type{PrimitiveType{NumberType}, PrimitiveType{NumberType}},
		NewSliceType(s._type.elementType),
	)

	s.methods["each"] = *NewNativeFunction(
		func(args ...Value) Value {
			function, ok := args[0].(FunctionValue)
			if !ok {
				panic("only functions are allowed as parameters for the each function")
			}

			slice := NewSlice(make([]Value, 0), function.returnType)

			var isFunctionWithIndex bool = false
			if len(function.parameters) == 2 {
				isFunctionWithIndex = true

				paramType := EvaluateType(function.parameters[0].Type, function.closure)
				if !paramType.Equals(PrimitiveType{NumberType}) {
					panic("First parameter type must be a number")
				}

				paramType = EvaluateType(function.parameters[1].Type, function.closure)
				if !paramType.Equals(s._type.elementType) {
					panic(fmt.Sprintf("Second parameter type must be %s, but got %s", s._type.elementType.String(), paramType.String()))
				}
			} else if len(function.parameters) == 1 {
				paramType := EvaluateType(function.parameters[0].Type, function.closure)
				if !paramType.Equals(s._type.elementType) {
					panic(fmt.Sprintf("Parameter type must be %s, but got %s", s._type.elementType.String(), paramType.String()))
				}
			} else {
				panic("each method expects exactly two arguments or one argument")
			}

			for index, element := range *s.elements {
				var callValue Value
				var err error
				if isFunctionWithIndex {
					callValue, err = function.Call(NewNumber(float64(index)), element)
					if err != nil {
						panic(err)
					}
				} else {
					callValue, err = function.Call(element)
					if err != nil {
						panic(err)
					}
				}

				slice.Append(callValue)
			}

			return slice
		},
		[]Type{PrimitiveType{AnyType}},
		NewSliceType(PrimitiveType{AnyType}),
	)

	s.methods["filter"] = *NewNativeFunction(
		func(args ...Value) Value {
			function, ok := args[0].(FunctionValue)
			if !ok {
				panic("only functions are allowed as parameters for the filter function")
			}

			if !function.returnType.Equals(PrimitiveType{BooleanType}) {
				panic("function return type must be a boolean for the filter function")
			}

			slice := NewSlice(make([]Value, 0), function.returnType)

			var isFunctionWithIndex bool = false
			if len(function.parameters) == 2 {
				isFunctionWithIndex = true

				paramType := EvaluateType(function.parameters[0].Type, function.closure)
				if !paramType.Equals(PrimitiveType{NumberType}) {
					panic("First parameter type must be a number")
				}

				paramType = EvaluateType(function.parameters[1].Type, function.closure)
				if !paramType.Equals(s._type.elementType) {
					panic(fmt.Sprintf("Second parameter type must be %s, but got %s", s._type.elementType.String(), paramType.String()))
				}
			} else if len(function.parameters) == 1 {
				paramType := EvaluateType(function.parameters[0].Type, function.closure)
				if !paramType.Equals(s._type.elementType) {
					panic(fmt.Sprintf("Parameter type must be %s, but got %s", s._type.elementType.String(), paramType.String()))
				}
			} else {
				panic("filter method expects exactly two arguments or one argument")
			}

			for index, element := range *s.elements {
				var callValue Value
				var err error
				if isFunctionWithIndex {
					callValue, err = function.Call(NewNumber(float64(index)), element)
					if err != nil {
						panic(err)
					}
				} else {
					callValue, err = function.Call(element)
					if err != nil {
						panic(err)
					}
				}

				if booleanValue, ok := callValue.(Boolean); ok && booleanValue.value {
					slice.Append(callValue)
				}
			}
			return slice
		},
		[]Type{PrimitiveType{AnyType}},
		NewSliceType(PrimitiveType{AnyType}),
	)
}
