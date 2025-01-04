package interpreter

import (
	"fmt"
	"strconv"
)

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

// Slice implements the Value interface
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
			var isFunctionWithValue bool = false
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
				isFunctionWithValue = true
				paramType := EvaluateType(function.parameters[0].Type, function.closure)
				if !paramType.Equals(s._type.elementType) {
					panic(fmt.Sprintf("Parameter type must be %s, but got %s", s._type.elementType.String(), paramType.String()))
				}
			} else if len(function.parameters) != 0 {
				panic("each method expects exactly two, one or zero arguments")
			}

			for index, element := range *s.elements {
				var callValue Value
				var err error
				if isFunctionWithIndex {
					callValue, err = function.Call(NewNumber(float64(index)), element)
				} else if isFunctionWithValue {
					callValue, err = function.Call(element)
				} else {
					callValue, err = function.Call()
				}

				if err != nil {
					panic(err)
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
			var isFunctionWithValue bool = false
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
				isFunctionWithValue = true

				paramType := EvaluateType(function.parameters[0].Type, function.closure)
				if !paramType.Equals(s._type.elementType) {
					panic(fmt.Sprintf("Parameter type must be %s, but got %s", s._type.elementType.String(), paramType.String()))
				}
			} else if len(function.parameters) != 0 {
				panic("filter method expects exactly two, one or zero arguments")
			}

			for index, element := range *s.elements {
				var callValue Value
				var err error
				if isFunctionWithIndex {
					callValue, err = function.Call(NewNumber(float64(index)), element)
				} else if isFunctionWithValue {
					callValue, err = function.Call(element)
				} else {
					callValue, err = function.Call()
				}

				if err != nil {
					panic(err)
				}

				if booleanValue, ok := callValue.(Boolean); ok && booleanValue.value {
					slice.Append(element.Clone())
				}
			}
			return slice
		},
		[]Type{PrimitiveType{AnyType}},
		NewSliceType(PrimitiveType{AnyType}),
	)
}
