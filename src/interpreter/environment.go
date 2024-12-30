package interpreter

import (
	"fmt"

	"github.com/table-harmony/HarmonyLang/src/ast"
)

type Environment struct {
	parent    *Environment
	variables map[string]RuntimeVariable
	functions []RuntimeFunction
}

func create_environment(parent *Environment) *Environment {
	return &Environment{
		parent:    parent,
		variables: make(map[string]RuntimeVariable),
	}
}

func (env *Environment) declare_variable(variable RuntimeVariable) error {
	_, exists := env.variables[variable.Identifier]

	if exists {
		return fmt.Errorf("variable '%s' already declared", variable.Identifier)
	}

	value := variable.getValue()

	if value == nil {
		value = GetDefaultValue(variable.ExplicitType)
	}

	if value.getType() != variable.ExplicitType && variable.ExplicitType != AnyType {
		return fmt.Errorf("type mismatch: variable '%s' declared as %v but got %v",
			variable.Identifier, variable.ExplicitType.ToString(),
			value.getType().ToString(),
		)
	}

	env.variables[variable.Identifier] = variable
	return nil
}

func (env *Environment) declare_function(function RuntimeFunction) error {
	_, err := env.get_function(function.Identifier, function.Parameters)
	if err == nil {
		return fmt.Errorf("function '%s' already declared", function.Identifier)
	}

	env.functions = append(env.functions, function)

	return nil
}

func (env *Environment) get_variable(identifier string) (RuntimeVariable, error) {
	env, err := env.resolve_variable(identifier)

	if err != nil {
		return RuntimeVariable{}, err
	}

	return env.variables[identifier], nil
}

func (env *Environment) get_function(identifier string, params []ast.Parameter) (RuntimeFunction, error) {
	function, err := env.resolve_function(identifier, params)
	if err != nil {
		return RuntimeFunction{}, err
	}

	return function, nil
}

func compare_parameters(params1, params2 []ast.Parameter) bool {
	if len(params1) != len(params2) {
		return false
	}

	for i := range params1 {
		if params1[i] != params2[i] {
			return false
		}
	}

	return true
}

func (env *Environment) resolve_variable(identifier string) (*Environment, error) {
	if _, exists := env.variables[identifier]; exists {
		return env, nil
	}

	if env.parent == nil {
		return nil, fmt.Errorf("variable '%s' not declared", identifier)
	}

	return env.parent.resolve_variable(identifier)
}

func (env *Environment) resolve_function(identifier string, params []ast.Parameter) (RuntimeFunction, error) {
	for _, function := range env.functions {
		if function.Identifier == identifier && compare_parameters(function.Parameters, params) {
			return function, nil
		}
	}

	if env.parent == nil {
		return RuntimeFunction{}, fmt.Errorf("function '%s' not declared", identifier)
	}

	return env.parent.resolve_function(identifier, params)
}
