package interpreter

import "fmt"

var print = NewNativeFunction(native_print, []Type{PrimitiveType{AnyType}}, PrimitiveType{NilType})

func native_print(args ...Value) Value {
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(arg.String())
	}
	return NewNil()
}

var println = NewNativeFunction(native_println, []Type{PrimitiveType{AnyType}}, PrimitiveType{NilType})

func native_println(args ...Value) Value {
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(arg.String())
	}
	fmt.Print("\n")
	return NewNil()
}
