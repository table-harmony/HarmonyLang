package interpreter

import (
	"fmt"

	"github.com/table-harmony/HarmonyLang/src/core"
)

func evaluate_addition(left, right core.Value) core.Value {
	switch left := left.(type) {
	case core.Number:
		right, err := core.ExpectValue[core.Number](right)
		if err != nil {
			panic("right operand must be a number")
		}
		return core.NewNumber(left.Value() + right.Value())
	case core.String:
		right, err := core.ExpectValue[core.String](right)
		if err != nil {
			panic("right operand must be a string")
		}
		return core.NewString(left.Value() + right.Value())
	default:
		panic(fmt.Sprintf("cannot add values of type %v and %v", left.Type(), right.Type()))
	}
}

func evaluate_subtraction(left, right core.Value) core.Value {
	leftNum, err := core.ExpectValue[core.Number](left)
	if err != nil {
		panic("left operand must be a number")
	}
	rightNum, err := core.ExpectValue[core.Number](right)
	if err != nil {
		panic("right operand must be a number")
	}
	return core.NewNumber(leftNum.Value() - rightNum.Value())
}

func evaluate_multiplication(left, right core.Value) core.Value {
	leftNum, err := core.ExpectValue[core.Number](left)
	if err != nil {
		panic("left operand must be a number")
	}
	rightNum, err := core.ExpectValue[core.Number](right)
	if err != nil {
		panic("right operand must be a number")
	}
	return core.NewNumber(leftNum.Value() * rightNum.Value())
}

func evaluate_division(left, right core.Value) core.Value {
	leftNum, err := core.ExpectValue[core.Number](left)
	if err != nil {
		panic("left operand must be a number")
	}
	rightNum, err := core.ExpectValue[core.Number](right)
	if err != nil {
		panic("right operand must be a number")
	}
	if rightNum.Value() == 0 {
		panic("division by zero")
	}
	return core.NewNumber(leftNum.Value() / rightNum.Value())
}

func evaluate_modulo(left, right core.Value) core.Value {
	leftNum, err := core.ExpectValue[core.Number](left)
	if err != nil {
		panic("left operand must be a number")
	}
	rightNum, err := core.ExpectValue[core.Number](right)
	if err != nil {
		panic("right operand must be a number")
	}
	if rightNum.Value() == 0 {
		panic("modulo by zero")
	}
	return core.NewNumber(float64(int(leftNum.Value()) % int(rightNum.Value())))
}

func evaluate_less_than(left, right core.Value) core.Value {
	switch left := left.(type) {
	case core.Number:
		right, err := core.ExpectValue[core.Number](right)
		if err != nil {
			panic("right operand must be a number")
		}
		return core.NewBoolean(left.Value() < right.Value())
	case core.String:
		right, err := core.ExpectValue[core.String](right)
		if err != nil {
			panic("right operand must be a string")
		}
		return core.NewBoolean(left.Value() < right.Value())
	default:
		panic(fmt.Sprintf("cannot compare values of type %v", left.Type()))
	}
}

func evaluate_less_equals(left, right core.Value) core.Value {
	switch left := left.(type) {
	case core.Number:
		right, err := core.ExpectValue[core.Number](right)
		if err != nil {
			panic("right operand must be a number")
		}
		return core.NewBoolean(left.Value() <= right.Value())
	case core.String:
		right, err := core.ExpectValue[core.String](right)
		if err != nil {
			panic("right operand must be a string")
		}
		return core.NewBoolean(left.Value() <= right.Value())
	default:
		panic(fmt.Sprintf("cannot compare values of type %v", left.Type()))
	}
}

func evaluate_greater_than(left, right core.Value) core.Value {
	switch left := left.(type) {
	case core.Number:
		right, err := core.ExpectValue[core.Number](right)
		if err != nil {
			panic("right operand must be a number")
		}
		return core.NewBoolean(left.Value() > right.Value())
	case core.String:
		right, err := core.ExpectValue[core.String](right)
		if err != nil {
			panic("right operand must be a string")
		}
		return core.NewBoolean(left.Value() > right.Value())
	default:
		panic(fmt.Sprintf("cannot compare values of type %v", left.Type()))
	}
}

func evaluate_greater_equals(left, right core.Value) core.Value {
	switch left := left.(type) {
	case core.Number:
		right, err := core.ExpectValue[core.Number](right)
		if err != nil {
			panic("right operand must be a number")
		}
		return core.NewBoolean(left.Value() >= right.Value())
	case core.String:
		right, err := core.ExpectValue[core.String](right)
		if err != nil {
			panic("right operand must be a string")
		}
		return core.NewBoolean(left.Value() >= right.Value())
	default:
		panic(fmt.Sprintf("cannot compare values of type %v", left.Type()))
	}
}

func evaluate_equals(left, right core.Value) core.Value {
	if left.Type() != right.Type() {
		return core.NewBoolean(false)
	}

	switch left := left.(type) {
	case core.Number:
		right, _ := core.ExpectValue[core.Number](right)
		return core.NewBoolean(left.Value() == right.Value())
	case core.String:
		right, _ := core.ExpectValue[core.String](right)
		return core.NewBoolean(left.Value() == right.Value())
	case core.Boolean:
		right, _ := core.ExpectValue[core.Boolean](right)
		return core.NewBoolean(left.Value() == right.Value())
	case core.Nil:
		return core.NewBoolean(true) // nil equals nil
	default:
		panic(fmt.Sprintf("cannot compare values of type %v", left.Type()))
	}
}

func evaluate_not_equals(left, right core.Value) core.Value {
	equals := evaluate_equals(left, right)
	boolVal, _ := core.ExpectValue[core.Boolean](equals)
	return core.NewBoolean(!boolVal.Value())
}

func evaluate_logical_or(left, right core.Value) core.Value {
	leftBool, err := core.ExpectValue[core.Boolean](left)
	if err != nil {
		panic("left operand must be a boolean")
	}

	if leftBool.Value() {
		return core.NewBoolean(true)
	}

	rightBool, err := core.ExpectValue[core.Boolean](right)
	if err != nil {
		panic("right operand must be a boolean")
	}
	return core.NewBoolean(rightBool.Value())
}

func evaluate_logical_and(left, right core.Value) core.Value {
	leftBool, err := core.ExpectValue[core.Boolean](left)
	if err != nil {
		panic("left operand must be a boolean")
	}

	if !leftBool.Value() {
		return core.NewBoolean(false)
	}

	rightBool, err := core.ExpectValue[core.Boolean](right)
	if err != nil {
		panic("right operand must be a boolean")
	}
	return core.NewBoolean(rightBool.Value())
}
