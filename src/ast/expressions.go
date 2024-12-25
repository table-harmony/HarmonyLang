package ast

import "github.com/table-harmony/HarmonyLang/src/lexer"

// Literal Expressions

type NumberExpression struct {
	Value float64
}

func (node NumberExpression) expression() {}

type StringExpression struct {
	Value string
}

func (node StringExpression) expression() {}

type SymbolExpression struct {
	Value string
}

func (node SymbolExpression) expression() {}

type BooleanExpression struct {
	Value bool
}

func (node BooleanExpression) expression() {}

// Complex Expressions

type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator lexer.Token
}

func (node BinaryExpression) expression() {}

type AssignmentExpression struct {
	Assigne  Expression
	Value    Expression
	Operator lexer.Token
}

func (node AssignmentExpression) expression() {}

type PrefixExpression struct {
	Operator lexer.Token
	Right    Expression
}

func (node PrefixExpression) expression() {}

type SwitchExpression struct {
	Value Expression
	Cases []SwitchCase
}

func (node SwitchExpression) expression() {}

type SwitchCase struct {
	Pattern Expression
	Value   Expression
}
