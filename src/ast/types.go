package ast

type SymbolType struct {
	Value string
}

func (SymbolType) _type() {}

type StringType struct {
}

func (StringType) _type() {}

type BooleanType struct {
}

func (BooleanType) _type() {}

type NumberType struct {
}

func (NumberType) _type() {}

type ArrayType struct {
	Underlying Type
}

func (ArrayType) _type() {}
