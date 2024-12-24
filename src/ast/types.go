package ast

type SymbolType struct {
	Value string
}

func (t SymbolType) _type() {}

type ArrayType struct {
	Underlying Type
}

func (t ArrayType) _type() {}
