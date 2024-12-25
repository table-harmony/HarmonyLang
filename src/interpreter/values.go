package interpreter

type RuntimeValue interface {
	value()
}

type NumberValue struct {
	Value float64
}

func (NumberValue) value() {}

type StringValue struct {
	Value string
}

func (StringValue) value() {}

type BooleanValue struct {
	Value bool
}

func (BooleanValue) value() {}
