package ast

import "github.com/table-harmony/HarmonyLang/src/lexer"

// Literal Expressions

type NumberExpression struct {
	Value float64
}

func (NumberExpression) expression() {}

type StringExpression struct {
	Value string
}

func (StringExpression) expression() {}

type SymbolExpression struct {
	Value string
}

func (SymbolExpression) expression() {}

type BooleanExpression struct {
	Value bool
}

func (BooleanExpression) expression() {}

type NilExpression struct {
}

func (NilExpression) expression() {}

// Complex Expressions

type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator lexer.Token
}

func (BinaryExpression) expression() {}

type PrefixExpression struct {
	Operator lexer.Token
	Right    Expression
}

func (PrefixExpression) expression() {}

type TernaryExpression struct {
	Condition  Expression
	Alternate  Expression
	Consequent Expression
}

func (TernaryExpression) expression() {}

type CallExpression struct {
	Caller Expression
	Params []Expression
}

func (CallExpression) expression() {}

type MemberExpression struct {
	Owner    Expression
	Property Expression
}

func (MemberExpression) expression() {}

type ComputedMemberExpression struct {
	Owner    Expression
	Property Expression
}

func (ComputedMemberExpression) expression() {}

type BlockExpression struct {
	Statements []Statement
}

func (BlockExpression) expression() {}

type IfExpression struct {
	Condition  Expression
	Consequent BlockExpression
	Alternate  BlockExpression
}

func (IfExpression) expression() {}

type SwitchExpression struct {
	Value Expression
	Cases []SwitchCaseStatement
}

func (SwitchExpression) expression() {}

type SwitchCaseStatement struct {
	Patterns  []Expression
	Body      BlockExpression
	IsDefault bool
}

func (SwitchCaseStatement) statement() {}
