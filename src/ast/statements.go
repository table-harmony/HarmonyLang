package ast

type BlockStatement struct {
	Body []Statement
}

func (node BlockStatement) statement() {

}

type ExpressionStatement struct {
	Expression Expression
}

func (node ExpressionStatement) statement() {

}
