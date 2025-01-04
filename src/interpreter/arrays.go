package interpreter

import (
	"fmt"
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

	if otherArray.size != a.size || !a.elementType.Equals(otherArray.elementType) {
		return false
	}

	return true
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

	a.methods["each"] = *NewNativeFunction(
		func(args ...Value) Value {
			function, ok := args[0].(FunctionValue)
			if !ok {
				panic("only functions are allowed as parameters for the each function")
			}

			arr := NewArray(make([]Value, 0), NewNumber(float64(a._type.size)), PrimitiveType{AnyType})

			var isFunctionWithIndex bool = false
			var isFunctionWithValue bool = false
			if len(function.parameters) == 2 {
				isFunctionWithIndex = true

				paramType := EvaluateType(function.parameters[0].Type, function.closure)
				if !paramType.Equals(PrimitiveType{NumberType}) {
					panic("First parameter type must be a number")
				}

				paramType = EvaluateType(function.parameters[1].Type, function.closure)
				if !paramType.Equals(a._type.elementType) {
					panic(fmt.Sprintf("Second parameter type must be %s, but got %s", a._type.elementType.String(), paramType.String()))
				}
			} else if len(function.parameters) == 1 {
				isFunctionWithValue = true
				paramType := EvaluateType(function.parameters[0].Type, function.closure)
				if !paramType.Equals(a._type.elementType) {
					panic(fmt.Sprintf("Parameter type must be %s, but got %s", a._type.elementType.String(), paramType.String()))
				}
			} else if len(function.parameters) != 0 {
				panic("each method expects exactly two, one or zero arguments")
			}

			for index, element := range a.elements {
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

				arr.elements[index] = callValue
			}

			return arr
		},
		[]Type{PrimitiveType{AnyType}},
		NewArrayType(NewNumber(float64(a._type.size)), PrimitiveType{AnyType}),
	)

	a.methods["filter"] = *NewNativeFunction(
		func(args ...Value) Value {
			function, ok := args[0].(FunctionValue)
			if !ok {
				panic("only functions are allowed as parameters for the filter function")
			}

			if !function.returnType.Equals(PrimitiveType{BooleanType}) {
				panic("function return type must be a boolean for the filter function")
			}

			arr := NewArray(make([]Value, 0), NewNumber(float64(a._type.size)), function.returnType)

			var isFunctionWithIndex bool = false
			var isFunctionWithValue bool = false
			if len(function.parameters) == 2 {
				isFunctionWithIndex = true

				paramType := EvaluateType(function.parameters[0].Type, function.closure)
				if !paramType.Equals(PrimitiveType{NumberType}) {
					panic("First parameter type must be a number")
				}

				paramType = EvaluateType(function.parameters[1].Type, function.closure)
				if !paramType.Equals(a._type.elementType) {
					panic(fmt.Sprintf("Second parameter type must be %s, but got %s", a._type.elementType.String(), paramType.String()))
				}
			} else if len(function.parameters) == 1 {
				isFunctionWithValue = true
				paramType := EvaluateType(function.parameters[0].Type, function.closure)
				if !paramType.Equals(a._type.elementType) {
					panic(fmt.Sprintf("Parameter type must be %s, but got %s", a._type.elementType.String(), paramType.String()))
				}
			} else if len(function.parameters) != 0 {
				panic("filter method expects exactly two, one or zero arguments")
			}

			for index, element := range a.elements {
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
					arr.elements[index] = element.Clone()
				} else {
					arr.elements[index] = arr._type.DefaultValue()
				}
			}

			return arr
		},
		[]Type{PrimitiveType{AnyType}},
		NewArrayType(NewNumber(float64(a._type.size)), PrimitiveType{AnyType}),
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
