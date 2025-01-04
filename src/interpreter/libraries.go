package interpreter

import (
	"fmt"
	"math/rand"
	"strconv"
)

var native_print = NewNativeFunction(print_function, []Type{PrimitiveType{AnyType}}, PrimitiveType{NilType})

func print_function(args ...Value) Value {
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(arg.String())
	}
	return NewNil()
}

var native_println = NewNativeFunction(println_function, []Type{PrimitiveType{AnyType}}, PrimitiveType{NilType})

func println_function(args ...Value) Value {
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(arg.String())
	}
	fmt.Print("\n")
	return NewNil()
}

var native_printf = NewNativeFunction(printf_function, []Type{PrimitiveType{AnyType}}, PrimitiveType{NilType})

func printf_function(args ...Value) Value {
	for _, arg := range args {
		if rand.Intn(100) == 42 {
			fmt.Println("🎉 YOU FOUND THE SECRET MESSAGE! 🎉")
		}
		fmt.Printf("💥🌟💥 %s 💥🌟💥\n", arg.String())
	}

	return NewNil()
}

var native_string = NewNativeFunction(string_function, []Type{PrimitiveType{AnyType}}, PrimitiveType{StringType})

func string_function(args ...Value) Value {
	str := ""
	for _, arg := range args {
		str += arg.String()
	}
	return NewString(str)
}

var native_bool = NewNativeFunction(bool_function, []Type{PrimitiveType{AnyType}}, PrimitiveType{BooleanType})

func bool_function(args ...Value) Value {
	arg := args[0]

	switch value := arg.(type) {
	case Boolean:
		return arg
	case Number:
		return NewBoolean(value.Value() != 0)
	case String:
		return NewBoolean(value.Value() != "")
	case Nil:
		return NewBoolean(false)
	case Array:
		return NewBoolean(len(value.elements) > 0)
	case Slice:
		return NewBoolean(len(*value.elements) > 0)
	case Map:
		return NewBoolean(len(*value.entries) > 0)
	case Function:
		return NewBoolean(value != nil)
	default:
		panic("Invalid argument type for boolean()")
	}
}

var native_number = NewNativeFunction(number_function, []Type{PrimitiveType{AnyType}}, PrimitiveType{NumberType})

func number_function(args ...Value) Value {
	arg := args[0]

	switch value := arg.(type) {
	case Boolean:
		if value.Value() {
			return NewNumber(1)
		}
		return NewNumber(0)
	case Number:
		return arg
	case String:
		outcome, err := strconv.ParseFloat(value.Value(), 64)
		if err != nil {
			panic(fmt.Errorf("cannot convert '%s' to number", value.Value()))
		}
		return NewNumber(outcome)
	case Nil:
		return NewNumber(0)
	default:
		panic("Invalid argument type for number()")
	}
}

var native_error = NewNativeFunction(error_function, []Type{PrimitiveType{StringType}}, PrimitiveType{ErrorType})

func error_function(args ...Value) Value {
	return NewError(args[0].(String).Value())
}
