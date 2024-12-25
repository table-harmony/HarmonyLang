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

type Parameter struct {
	Name string
	Type Type
}

type FunctionDeclarationStatment struct {
	Modifiers  []Statement
	Parameters []Parameter
	Name       string
	Body       []Statement
	ReturnType Type
}

func (node FunctionDeclarationStatment) statement() {}

type StructProperty struct {
	IsStatic bool
	Type     Type
}

type StructMethod struct {
	IsStatic bool
}

type StructDeclarationStatement struct {
	Identifier string
	Properties map[string]StructProperty
	Methods    map[string]StructMethod
}

func (node StructDeclarationStatement) statement() {}

type ImportStatement struct {
	Name string
	From string
}

func (node ImportStatement) statement() {}

type IfStatement struct {
	Condition  Expression
	Consequent Statement
	Alternate  Statement
}

func (node IfStatement) statement() {}

type SwitchCaseStatement struct {
	Pattern Expression
	Body    Statement
}

type SwitchStatement struct {
	Value Expression
	Cases []SwitchCaseStatement
}

func (node SwitchStatement) statement() {}

type ForStatement struct {
	Initializer Statement
	Condition   Expression
	Post        []Expression
	Body        Statement
}

func (node ForStatement) statement() {}
