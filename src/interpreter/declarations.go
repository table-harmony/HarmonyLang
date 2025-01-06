package interpreter

import (
	"fmt"

	"github.com/table-harmony/HarmonyLang/src/ast"
)

type DeclarationKind int

const (
	FunctionDeclaration DeclarationKind = iota
	StructDeclaration
	InterfaceDeclaration
)

type Declaration interface {
	Kind() DeclarationKind
	Identifier() string
	IsComplete() bool
	Complete(ast.Statement, *Scope) error
}

// FunctionDecl represents a function declaration
type FunctionDecl struct {
	identifier string
	isComplete bool
	signature  FunctionType
	impl       *FunctionValue
}

func NewFunctionDecl(identifier string, signature FunctionType) *FunctionDecl {
	return &FunctionDecl{
		identifier: identifier,
		signature:  signature,
		isComplete: false,
	}
}

func (f *FunctionDecl) Kind() DeclarationKind          { return FunctionDeclaration }
func (f *FunctionDecl) Identifier() string             { return f.identifier }
func (f *FunctionDecl) IsComplete() bool               { return f.isComplete }
func (f *FunctionDecl) Implementation() *FunctionValue { return f.impl }

func (f *FunctionDecl) Complete(statement ast.Statement, scope *Scope) error {
	expectedStatement, ok := statement.(ast.FunctionDeclarationStatment)
	if !ok {
		return fmt.Errorf("expected function declaration, got %T", expectedStatement)
	}

	f.impl = NewFunctionValue(
		expectedStatement.Parameters,
		expectedStatement.Body,
		EvaluateType(expectedStatement.ReturnType, scope),
		scope,
	)
	f.isComplete = true

	return nil
}
