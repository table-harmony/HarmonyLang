package interpreter

import "errors"

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
	_, err := env.resolve(variable.Identifier)

	if err == nil {
		return errors.New("variable already declared")
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
		return err
	}

	variable := env.variables[identifier]
	variable.Value = value
	env.variables[identifier] = variable

	return nil
}

func (env *Environment) resolve(identifier string) (*Environment, error) {
	if _, exists := env.variables[identifier]; exists {
		return env, nil
	}

	if env.parent == nil {
		return nil, errors.New("variable not declared")
	}

	return env.parent.resolve(identifier)
}
