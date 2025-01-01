package interpreter

import (
	"fmt"
)

func evaluate_addition(left, right Value) Value {
	switch left := left.(type) {
	case Number:
		right, err := ExpectValue[Number](right)
		if err != nil {
			panic("right operand must be a number")
		}
		return NewNumber(left.Value() + right.Value())
	case String:
		right, err := ExpectValue[String](right)
		if err != nil {
			panic("right operand must be a string")
		}
		return NewString(left.Value() + right.Value())
	default:
		panic(fmt.Sprintf("cannot add values of type %v and %v", left.Type(), right.Type()))
	}
}

func evaluate_subtraction(left, right Value) Value {
	leftNum, err := ExpectValue[Number](left)
	if err != nil {
		panic("left operand must be a number")
	}
	rightNum, err := ExpectValue[Number](right)
	if err != nil {
		panic("right operand must be a number")
	}
	return NewNumber(leftNum.Value() - rightNum.Value())
}

func evaluate_multiplication(left, right Value) Value {
	leftNum, err := ExpectValue[Number](left)
	if err != nil {
		panic("left operand must be a number")
	}
	rightNum, err := ExpectValue[Number](right)
	if err != nil {
		panic("right operand must be a number")
	}
	return NewNumber(leftNum.Value() * rightNum.Value())
}

func evaluate_division(left, right Value) Value {
	leftNum, err := ExpectValue[Number](left)
	if err != nil {
		panic("left operand must be a number")
	}
	rightNum, err := ExpectValue[Number](right)
	if err != nil {
		panic("right operand must be a number")
	}
	if rightNum.Value() == 0 {
		panic("division by zero")
	}
	return NewNumber(leftNum.Value() / rightNum.Value())
}

func evaluate_modulo(left, right Value) Value {
	leftNum, err := ExpectValue[Number](left)
	if err != nil {
		panic("left operand must be a number")
	}
	rightNum, err := ExpectValue[Number](right)
	if err != nil {
		panic("right operand must be a number")
	}
	if rightNum.Value() == 0 {
		panic("modulo by zero")
	}
	return NewNumber(float64(int(leftNum.Value()) % int(rightNum.Value())))
}

func evaluate_less_than(left, right Value) Value {
	switch left := left.(type) {
	case Number:
		right, err := ExpectValue[Number](right)
		if err != nil {
			panic("right operand must be a number")
		}
		return NewBoolean(left.Value() < right.Value())
	case String:
		right, err := ExpectValue[String](right)
		if err != nil {
			panic("right operand must be a string")
		}
		return NewBoolean(left.Value() < right.Value())
	default:
		panic(fmt.Sprintf("cannot compare values of type %v", left.Type()))
	}
}

func evaluate_less_equals(left, right Value) Value {
	switch left := left.(type) {
	case Number:
		right, err := ExpectValue[Number](right)
		if err != nil {
			panic("right operand must be a number")
		}
		return NewBoolean(left.Value() <= right.Value())
	case String:
		right, err := ExpectValue[String](right)
		if err != nil {
			panic("right operand must be a string")
		}
		return NewBoolean(left.Value() <= right.Value())
	default:
		panic(fmt.Sprintf("cannot compare values of type %v", left.Type()))
	}
}

func evaluate_greater_than(left, right Value) Value {
	switch left := left.(type) {
	case Number:
		right, err := ExpectValue[Number](right)
		if err != nil {
			panic("right operand must be a number")
		}
		return NewBoolean(left.Value() > right.Value())
	case String:
		right, err := ExpectValue[String](right)
		if err != nil {
			panic("right operand must be a string")
		}
		return NewBoolean(left.Value() > right.Value())
	default:
		panic(fmt.Sprintf("cannot compare values of type %v", left.Type()))
	}
}

func evaluate_greater_equals(left, right Value) Value {
	switch left := left.(type) {
	case Number:
		right, err := ExpectValue[Number](right)
		if err != nil {
			panic("right operand must be a number")
		}
		return NewBoolean(left.Value() >= right.Value())
	case String:
		right, err := ExpectValue[String](right)
		if err != nil {
			panic("right operand must be a string")
		}
		return NewBoolean(left.Value() >= right.Value())
	default:
		panic(fmt.Sprintf("cannot compare values of type %v", left.Type()))
	}
}

func evaluate_equals(left, right Value) Value {
	if left.Type() != right.Type() {
		return NewBoolean(false)
	}

	switch left := left.(type) {
	case Number:
		right, _ := ExpectValue[Number](right)
		return NewBoolean(left.Value() == right.Value())
	case String:
		right, _ := ExpectValue[String](right)
		return NewBoolean(left.Value() == right.Value())
	case Boolean:
		right, _ := ExpectValue[Boolean](right)
		return NewBoolean(left.Value() == right.Value())
	case Nil:
		return NewBoolean(true) // nil equals nil
	default:
		panic(fmt.Sprintf("cannot compare values of type %v", left.Type()))
	}
}

func evaluate_not_equals(left, right Value) Value {
	equals := evaluate_equals(left, right)
	boolVal, _ := ExpectValue[Boolean](equals)
	return NewBoolean(!boolVal.Value())
}

func evaluate_logical_or(left, right Value) Value {
	leftBool, err := ExpectValue[Boolean](left)
	if err != nil {
		panic("left operand must be a boolean")
	}

	if leftBool.Value() {
		return NewBoolean(true)
	}

	rightBool, err := ExpectValue[Boolean](right)
	if err != nil {
		panic("right operand must be a boolean")
	}
	return NewBoolean(rightBool.Value())
}

func evaluate_logical_and(left, right Value) Value {
	leftBool, err := ExpectValue[Boolean](left)
	if err != nil {
		panic("left operand must be a boolean")
	}

	if !leftBool.Value() {
		return NewBoolean(false)
	}

	rightBool, err := ExpectValue[Boolean](right)
	if err != nil {
		panic("right operand must be a boolean")
	}
	return NewBoolean(rightBool.Value())
}
