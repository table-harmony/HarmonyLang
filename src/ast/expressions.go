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

// Complex Expressions

type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator lexer.Token
}

func (node BinaryExpression) expression() {}

type AssignmentExpression struct {
	Assigne       Expression
	AssignedValue Expression
}

func (node AssignmentExpression) expression() {}
