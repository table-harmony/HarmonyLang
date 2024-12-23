package ast

type BlockStatement struct {
	Body []Statement
}

func (node BlockStatement) statement() {}

type ExpressionStatement struct {
	Expression Expression
}

func (node ExpressionStatement) statement() {}

type VariableDeclarationStatement struct {
	Identifier   string
	IsConstant   bool
	Value        Expression
	ExplicitType Type
}

func (node VariableDeclarationStatement) statement() {}
