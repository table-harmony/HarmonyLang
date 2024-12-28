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

type NullExpression struct {
}

func (NullExpression) expression() {}

// Complex Expressions

type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator lexer.Token
}

func (BinaryExpression) expression() {}

type AssignmentExpression struct {
	Assigne  Expression
	Value    Expression
	Operator lexer.Token
}

func (AssignmentExpression) expression() {}

type PrefixExpression struct {
	Operator lexer.Token
	Right    Expression
}

func (PrefixExpression) expression() {}

type SwitchExpression struct {
	Value Expression
	Cases []SwitchCase
}

func (SwitchExpression) expression() {}

type SwitchCase interface {
	GetValue() Expression
	GetPatterns() []Expression
}

type NormalSwitchCase struct {
	Patterns []Expression
	Value    Expression
}

func (n NormalSwitchCase) GetValue() Expression      { return n.Value }
func (n NormalSwitchCase) GetPatterns() []Expression { return n.Patterns }

type DefaultSwitchCase struct {
	Value Expression
}

func (n DefaultSwitchCase) GetValue() Expression      { return n.Value }
func (n DefaultSwitchCase) GetPatterns() []Expression { return make([]Expression, 0) }

type TernaryExpression struct {
	Condition  Expression
	Alternate  Expression
	Consequent Expression
}

func (TernaryExpression) expression() {}

type FunctionDeclarationExpression struct {
	Parameters []Parameter
	Body       []Statement
	ReturnType Type
}
