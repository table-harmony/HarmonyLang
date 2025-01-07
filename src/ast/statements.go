package ast

import "github.com/table-harmony/HarmonyLang/src/lexer"

type ExpressionStatement struct {
	Expression Expression
}

func (ExpressionStatement) statement() {}

type AssignmentStatement struct {
	Assigne  Expression
	Value    Expression
	Operator lexer.Token
}

func (AssignmentStatement) statement() {}

type VariableDeclarationStatement struct {
	Identifier   string
	IsConstant   bool
	Value        Expression
	ExplicitType Type
}

func (VariableDeclarationStatement) statement() {}

type MultiVariableDeclarationStatement struct {
	Declarations []VariableDeclarationStatement
}

func (MultiVariableDeclarationStatement) statement() {}

type Parameter struct {
	Name         string
	Type         Type
	DefaultValue Expression
}

type FunctionDeclarationStatment struct {
	Identifier string
	Parameters []Parameter
	Body       []Statement
	ReturnType Type
}

func (FunctionDeclarationStatment) statement() {}

type ImportStatement struct {
	Module       string
	Alias        string
	NamedImports map[string]string
}

func (ImportStatement) statement() {}

type TraditionalForStatement struct {
	Initializer Statement
	Condition   Expression
	Post        []Statement
	Body        []Statement
}

func (TraditionalForStatement) statement() {}

type IteratorForStatement struct {
	KeyIdentifier   string
	ValueIdentifier string
	Iterator        Expression
	Body            []Statement
}

func (IteratorForStatement) statement() {}

type BreakStatement struct{}

func (BreakStatement) statement() {}

type ContinueStatement struct{}

func (ContinueStatement) statement() {}

type ReturnStatement struct {
	Value Expression
}

func (ReturnStatement) statement() {}

type SwitchCaseStatement struct {
	Patterns  []Expression
	Body      BlockExpression
	IsDefault bool
}

func (SwitchCaseStatement) statement() {}

type ThrowStatement struct {
	Value Expression
}

func (ThrowStatement) statement() {}

type TypeDeclarationStatement struct {
	Identifier string
	Type       Type
}

func (TypeDeclarationStatement) statement() {}

type StructProperty struct {
	Identifier   string
	Type         Type
	DefaultValue Expression
	IsStatic     bool
}

type StructMethod struct {
	IsStatic    bool
	Declaration FunctionDeclarationStatment
}

type StructDeclarationStatement struct {
	Identifier string
	Properties []StructProperty
	Methods    []StructMethod
}

func (StructDeclarationStatement) statement() {}
