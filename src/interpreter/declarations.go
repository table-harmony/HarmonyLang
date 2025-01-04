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
	Name() string
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
func (f *FunctionDecl) Name() string                   { return f.identifier }
func (f *FunctionDecl) IsComplete() bool               { return f.isComplete }
func (f *FunctionDecl) Implementation() *FunctionValue { return f.impl }

func (f *FunctionDecl) Complete(stmt ast.Statement, scope *Scope) error {
	funcStmt, ok := stmt.(ast.FunctionDeclarationStatment)
	if !ok {
		return fmt.Errorf("expected function declaration, got %T", stmt)
	}

	f.impl = NewFunctionValue(
		funcStmt.Parameters,
		funcStmt.Body,
		EvaluateType(funcStmt.ReturnType, scope),
		scope,
	)
	f.isComplete = true

	return nil
}
