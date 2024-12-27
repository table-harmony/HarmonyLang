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
	NullType

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
	case NullType:
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

func isEqual(variable1 RuntimeValue, variable2 RuntimeValue) bool {
	if variable1.getValue().getType() != variable2.getValue().getType() {
		return false
	}

	return variable1.getValue() == variable2.getValue()
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

type RuntimeNull struct {
}

func (RuntimeNull) getType() RuntimeValueType { return NullType }
func (n RuntimeNull) getValue() RuntimeValue  { return n }

type RuntimeVariable struct {
	Identifier   string
	IsConstant   bool
	Value        RuntimeValue
	ExplicitType RuntimeValueType
}

func (RuntimeVariable) getType() RuntimeValueType { return VariableType }
func (v RuntimeVariable) getValue() RuntimeValue  { return v.Value }
