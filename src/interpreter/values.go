package interpreter

import (
	"fmt"

	"github.com/table-harmony/HarmonyLang/src/ast"
)

type ValueType int

const (
	NumberType ValueType = iota
	StringType
	BooleanType
)

type RuntimeValue interface {
	Type() ValueType
	AsNumber() (float64, error)
	AsString() (string, error)
	AsBoolean() (bool, error)
}

func type_mismatch_error(got, want ValueType) error {
	return fmt.Errorf("type mismatch: got %v, want %v", got, want)
}

type NumberRuntime struct {
	Value float64
}

func (n NumberRuntime) Type() ValueType            { return NumberType }
func (n NumberRuntime) AsNumber() (float64, error) { return n.Value, nil }
func (n NumberRuntime) AsString() (string, error) {
	return fmt.Sprintf("%f", n.Value), nil
}
func (n NumberRuntime) AsBoolean() (bool, error) {
	return n.Value != 0, nil
}

type StringRuntime struct {
	Value string
}

func (s StringRuntime) Type() ValueType { return StringType }
func (s StringRuntime) AsNumber() (float64, error) {
	return 0, type_mismatch_error(s.Type(), NumberType)
}
func (s StringRuntime) AsString() (string, error) { return s.Value, nil }
func (s StringRuntime) AsBoolean() (bool, error) {
	return s.Value != "", nil
}

type BooleanRuntime struct {
	Value bool
}

func (b BooleanRuntime) Type() ValueType { return BooleanType }
func (b BooleanRuntime) AsNumber() (float64, error) {
	if b.Value {
		return 1, nil
	}

	return 0, nil
}
func (b BooleanRuntime) AsString() (string, error) {
	return fmt.Sprintf("%b", b.Value), nil
}
func (b BooleanRuntime) AsBoolean() (bool, error) { return b.Value, nil }

type RuntimeVariable struct {
	Identifier   string
	IsConstant   bool
	Value        RuntimeValue
	ExplicitType ast.Type
}

func (variable RuntimeVariable) Type() ValueType            { return variable.Value.Type() }
func (variable RuntimeVariable) AsNumber() (float64, error) { return variable.Value.AsNumber() }
func (variable RuntimeVariable) AsString() (string, error)  { return variable.Value.AsString() }
func (variable RuntimeVariable) AsBoolean() (bool, error)   { return variable.Value.AsBoolean() }
