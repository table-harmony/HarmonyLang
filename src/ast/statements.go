package ast

type BlockStatement struct {
	Body []Statement
}

func (BlockStatement) statement() {}

type ExpressionStatement struct {
	Expression Expression
}

func (ExpressionStatement) statement() {}

type VariableDeclarationStatement struct {
	Identifier   string
	IsConstant   bool
	Value        Expression
	ExplicitType Type
}

func (VariableDeclarationStatement) statement() {}

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

func (FunctionDeclarationStatment) statement() {}

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

func (StructDeclarationStatement) statement() {}

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

func (IfStatement) statement() {}

type SwitchCaseStatement struct {
	Pattern Expression
	Body    []Statement
}

type SwitchStatement struct {
	Value Expression
	Cases []SwitchCaseStatement
}

func (SwitchStatement) statement() {}

type ForStatement struct {
	Initializer Statement
	Condition   Expression
	Post        []Expression
	Body        []Statement
}

func (ForStatement) statement() {}
