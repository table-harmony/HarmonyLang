package interpreter

import (
	"fmt"

	"github.com/table-harmony/HarmonyLang/src/ast"
)

type Scope struct {
	parent       *Scope
	storage      map[string]Reference
	declarations map[string]Declaration
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		parent:       parent,
		storage:      make(map[string]Reference),
		declarations: make(map[string]Declaration),
	}
}

func NewRootScope() *Scope {
	scope := NewScope(nil)

	// Declare native printing functions
	scope.Declare(NewFunctionReference("print", native_print))
	scope.Declare(NewFunctionReference("println", native_println))
	scope.Declare(NewFunctionReference("printf", native_printf))

	// Declare native type conversion functions
	scope.Declare(NewFunctionReference("string", native_string))
	scope.Declare(NewFunctionReference("bool", native_bool))
	scope.Declare(NewFunctionReference("number", native_number))
	scope.Declare(NewFunctionReference("error", native_error))

	return scope
}

func (scope *Scope) Declare(ref Reference) error {
	var identifier string

	switch ref := ref.(type) {
	case *VariableReference:
		identifier = ref.identifier
	case *FunctionReference:
		identifier = ref.identifier
	case *StructReference:
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

func (s *Scope) DeclareForward(decl Declaration) error {
	name := decl.Identifier()
	if existing, exists := s.declarations[name]; exists {
		if existing.IsComplete() {
			return fmt.Errorf("%d '%s' already declared", existing.Kind(), name)
		}
		return nil // Allow multiple incomplete declarations
	}
	s.declarations[name] = decl
	return nil
}

func (s *Scope) CompleteDeclaration(name string, stmt ast.Statement) error {
	decl, exists := s.declarations[name]
	if !exists {
		return fmt.Errorf("no forward declaration found for '%s'", name)
	}

	if decl.IsComplete() {
		return fmt.Errorf("%d '%s' already implemented", decl.Kind(), name)
	}

	return decl.Complete(stmt, s)
}

func (s *Scope) GetDeclaration(name string) (Declaration, bool) {
	if decl, exists := s.declarations[name]; exists {
		return decl, true
	}
	if s.parent != nil {
		return s.parent.GetDeclaration(name)
	}
	return nil, false
}

func (scope *Scope) String() string {
	str := ""
	for identifier, ref := range scope.storage {
		str += fmt.Sprintf("%s: %s \n", identifier, ref.String())
	}
	return str
}
