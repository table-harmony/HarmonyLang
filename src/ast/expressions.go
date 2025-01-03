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
	Alternate  Expression
}

func (IfExpression) expression() {}

type SwitchExpression struct {
	Value Expression
	Cases []SwitchCaseStatement
}

func (SwitchExpression) expression() {}

type ArrayInstantiationExpression struct {
	Size        Expression
	ElementType Type
	Elements    []Expression
}

func (ArrayInstantiationExpression) expression() {}

type SliceInstantiationExpression struct {
	ElementType Type
	Elements    []Expression
}

func (SliceInstantiationExpression) expression() {}

type MapInstantiationExpression struct {
	KeyType   Type
	ValueType Type
	Entries   []MapEntry
}

func (MapInstantiationExpression) expression() {}

type MapEntry struct {
	Key   Expression
	Value Expression
}

type FunctionDeclarationExpression struct {
	Parameters []Parameter
	Body       []Statement
	ReturnType Type
}

func (FunctionDeclarationExpression) expression() {}

type TryCatchExpression struct {
	TryBlock   Expression
	CatchBlock Expression
}

func (TryCatchExpression) expression() {}

type RangeExpression struct {
	Lower Expression
	Upper Expression
	Step  Expression
}

func (RangeExpression) expression() {}
