package interpreter

import (
	"fmt"
)

type Environment struct {
	parent    *Environment
	variables map[string]RuntimeVariable
}

func create_enviorment(parent *Environment) *Environment {
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
	//variable.Value = value

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

func (env *Environment) get_variable(identifier string) (RuntimeVariable, error) {
	env, err := env.resolve(identifier)

	if err != nil {
		return RuntimeVariable{}, err
	}

	return env.variables[identifier], nil
}

func (env *Environment) assign_variable(identifier string, value RuntimeValue) error {
	env, err := env.resolve(identifier)

	if err != nil {
		panic(err)
	}

	panic("not implemented yet use assign function")
}

func (env *Environment) resolve(identifier string) (*Environment, error) {
	if _, exists := env.variables[identifier]; exists {
		return env, nil
	}

	if env.parent == nil {
		return nil, fmt.Errorf("variable '%s' not declared", identifier)
	}

	return env.parent.resolve(identifier)
}
