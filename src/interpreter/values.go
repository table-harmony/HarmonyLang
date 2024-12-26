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
	get_type() ValueType
	as_number() (float64, error)
	as_string() (string, error)
	as_boolean() (bool, error)
}

func type_mismatch_error(got, want ValueType) error {
	return fmt.Errorf("type mismatch: got %v, want %v", got, want)
}

type NumberRuntime struct {
	Value float64
}

func (n NumberRuntime) get_type() ValueType         { return NumberType }
func (n NumberRuntime) as_number() (float64, error) { return n.Value, nil }
func (n NumberRuntime) as_string() (string, error) {
	return fmt.Sprintf("%.g", n.Value), nil
}
func (n NumberRuntime) as_boolean() (bool, error) {
	return n.Value != 0, nil
}

type StringRuntime struct {
	Value string
}

func (s StringRuntime) get_type() ValueType { return StringType }
func (s StringRuntime) as_number() (float64, error) {
	return 0, type_mismatch_error(s.get_type(), NumberType)
}
func (s StringRuntime) as_string() (string, error) { return s.Value, nil }
func (s StringRuntime) as_boolean() (bool, error) {
	return s.Value != "", nil
}

type BooleanRuntime struct {
	Value bool
}

func (b BooleanRuntime) get_type() ValueType { return BooleanType }
func (b BooleanRuntime) as_number() (float64, error) {
	if b.Value {
		return 1, nil
	}

	return 0, nil
}
func (b BooleanRuntime) as_string() (string, error) {
	return fmt.Sprintf("%t", b.Value), nil
}
func (b BooleanRuntime) as_boolean() (bool, error) { return b.Value, nil }

type RuntimeVariable struct {
	Identifier   string
	IsConstant   bool
	Value        RuntimeValue
	ExplicitType ast.Type
}

func (variable RuntimeVariable) get_type() ValueType         { return variable.Value.get_type() }
func (variable RuntimeVariable) as_number() (float64, error) { return variable.Value.as_number() }
func (variable RuntimeVariable) as_string() (string, error)  { return variable.Value.as_string() }
func (variable RuntimeVariable) as_boolean() (bool, error)   { return variable.Value.as_boolean() }

// TODO: decide whether there is equality between other types
func is_equal(variable1 RuntimeValue, variable2 RuntimeValue) bool {
	variable1_type, variable2_type := variable1.get_type(), variable2.get_type()

	// Handle string comparison
	if variable1_type == StringType && variable2_type == StringType {
		variable1_value, err1 := variable1.as_string()
		variable2_value, err2 := variable2.as_string()

		if err1 == nil && err2 == nil {
			return variable1_value == variable2_value
		}
	}

	// Handle number comparison
	if variable1_type == NumberType && variable2_type == NumberType {
		variable1_value, err1 := variable1.as_number()
		variable2_value, err2 := variable2.as_number()

		if err1 == nil && err2 == nil {
			return variable1_value == variable2_value
		}
	}

	// Handle boolean comparison
	if variable1_type == BooleanType && variable2_type == BooleanType {
		variable1_value, err1 := variable1.as_boolean()
		variable2_value, err2 := variable2.as_boolean()

		if err1 == nil && err2 == nil {
			return variable1_value == variable2_value
		}
	}

	//TODO: complex variables equality

	return false
}
