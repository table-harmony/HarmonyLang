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
	Name string
	From string
}

func (ImportStatement) statement() {}

type ForStatement struct {
	Initializer Statement
	Condition   Expression
	Post        []Expression
	Body        []Statement
}

func (ForStatement) statement() {}

type BreakStatement struct{}

func (BreakStatement) statement() {}

type ContinueStatement struct{}

func (ContinueStatement) statement() {}

type ReturnStatement struct {
	Value Expression
}

func (ReturnStatement) statement() {}
