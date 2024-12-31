package interpreter

import "fmt"

type Scope struct {
	parent    *Scope
	variables map[string]*VariableReference
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		parent:    parent,
		variables: make(map[string]*VariableReference),
	}
}

func (scope *Scope) DeclareVariable(variable *VariableReference) error {
	if _, exists := scope.variables[variable.identifier]; exists {
		return fmt.Errorf("redeclaration of '%s'", variable.identifier)
	}

	scope.variables[variable.identifier] = variable
	return nil
}

func (scope *Scope) Resolve(identifier string) (Reference, error) {
	if ref, exists := scope.variables[identifier]; exists {
		return ref, nil
	}
	if scope.parent != nil {
		return scope.parent.Resolve(identifier)
	}
	return nil, fmt.Errorf("undefined: %s", identifier)
}
