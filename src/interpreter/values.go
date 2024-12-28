package interpreter

import (
	"fmt"

	"github.com/table-harmony/HarmonyLang/src/helpers"
)

type RuntimeValueType int

const (
	// Literals
	NumberType RuntimeValueType = iota
	StringType
	BooleanType
	NilType

	VariableType
	FunctionType
	AnyType
)

func (_type RuntimeValueType) ToString() string {
	switch _type {
	case NumberType:
		return "number"
	case StringType:
		return "string"
	case BooleanType:
		return "bool"
	case NilType:
		return "null"
	case VariableType:
		return "variable"
	case FunctionType:
		return "function"
	case AnyType:
		return "any"
	default:
		return fmt.Sprintf("unknown(%d)", _type)
	}
}

type RuntimeValue interface {
	getType() RuntimeValueType
	getValue() RuntimeValue
}

func ExpectRuntimeValue[T RuntimeValue](value RuntimeValue) (T, error) {
	return helpers.ExpectType[T](value)
}

func GetDefaultValue(valueType RuntimeValueType) RuntimeValue {
	switch valueType {
	case NumberType:
		return RuntimeNumber{Value: 0}
	case StringType:
		return RuntimeString{Value: ""}
	case BooleanType:
		return RuntimeBoolean{Value: false}
	case NilType:
		return RuntimeNil{}
	case FunctionType:
		return RuntimeNil{}
	default:
		return RuntimeNil{}
	}
}

func isEqual(variable1 RuntimeValue, variable2 RuntimeValue) bool {
	if variable1 == nil || variable2 == nil {
		return variable1 == variable2
	}

	value1 := variable1.getValue()
	value2 := variable2.getValue()

	if value1 == nil || value2 == nil {
		return value1 == value2
	}

	if value1.getType() != value2.getType() {
		return false
	}

	return value1 == value2
}

type RuntimeNumber struct {
	Value float64
}

func (RuntimeNumber) getType() RuntimeValueType { return NumberType }
func (n RuntimeNumber) getValue() RuntimeValue  { return n }

type RuntimeString struct {
	Value string
}

func (RuntimeString) getType() RuntimeValueType { return StringType }
func (s RuntimeString) getValue() RuntimeValue  { return s }

type RuntimeBoolean struct {
	Value bool
}

func (RuntimeBoolean) getType() RuntimeValueType { return BooleanType }
func (b RuntimeBoolean) getValue() RuntimeValue  { return b }

type RuntimeNil struct {
}

func (RuntimeNil) getType() RuntimeValueType { return NilType }
func (n RuntimeNil) getValue() RuntimeValue  { return n }

type RuntimeVariable struct {
	Identifier   string
	IsConstant   bool
	Value        RuntimeValue
	ExplicitType RuntimeValueType
}

func (RuntimeVariable) getType() RuntimeValueType { return VariableType }
func (v RuntimeVariable) getValue() RuntimeValue  { return v.Value }
