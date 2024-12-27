package interpreter

import (
	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/helpers"
)

type RuntimeValueType int

const (
	NumberType RuntimeValueType = iota
	StringType
	BooleanType
	VariableType
	FunctionType
)

type RuntimeValue interface {
	getType() RuntimeValueType
	getValue() RuntimeValue
}

func ExpectRuntimeValue[T RuntimeValue](value RuntimeValue) (T, error) {
	return helpers.ExpectType[T](value)
}

func isEqual(variable1 RuntimeValue, variable2 RuntimeValue) bool {
	if variable1.getType() != variable2.getType() {
		return false
	}

	//TODO: equality incorrect it checks refrences not values
	return variable1 == variable2
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

type RuntimeVariable struct {
	Identifier   string
	IsConstant   bool
	Value        RuntimeValue
	ExplicitType ast.Type
}

func (RuntimeVariable) getType() RuntimeValueType { return VariableType }
func (v RuntimeVariable) getValue() RuntimeValue  { return v.Value }
