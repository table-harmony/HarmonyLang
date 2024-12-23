package ast

type PrimitiveType struct {
	Name string
}

func (t PrimitiveType) _type() {}

type ArrayType struct {
	Underlying Type
}

func (t ArrayType) _type() {}
