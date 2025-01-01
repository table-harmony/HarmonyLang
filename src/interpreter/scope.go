package interpreter

import "fmt"

type Scope struct {
	parent  *Scope
	storage map[string]Reference
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		parent:  parent,
		storage: make(map[string]Reference),
	}
}

func (scope *Scope) Declare(ref Reference) error {
	var identifier string

	switch ref := ref.(type) {
	case *VariableReference:
		identifier = ref.identifier
	case *FunctionReference:
		identifier = ref.identifier
	}

	if _, exists := scope.storage[identifier]; exists {
		return fmt.Errorf("redeclaration of '%s'", identifier)
	}

	scope.storage[identifier] = ref
	return nil
}

func (scope *Scope) Resolve(identifier string) (Reference, error) {
	if ref, exists := scope.storage[identifier]; exists {
		return ref, nil
	}
	if scope.parent != nil {
		return scope.parent.Resolve(identifier)
	}
	return nil, fmt.Errorf("undefined: %s", identifier)
}

func (scope *Scope) String() string {
	str := ""
	for identifier, ref := range scope.storage {
		str += fmt.Sprintf("%s: { %s } \n", identifier, ref.String())
	}
	return str
}
