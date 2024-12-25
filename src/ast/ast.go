package ast

import "github.com/table-harmony/HarmonyLang/src/helpers"

type Statement interface {
	statement()
}

type Expression interface {
	expression()
}

type Type interface {
	_type()
}

func ExpectExpression[T Expression](exprssion Expression) (T, error) {
	return helpers.ExpectType[T](exprssion)
}

func ExpectStatement[T Statement](statement Statement) (T, error) {
	return helpers.ExpectType[T](statement)
}
