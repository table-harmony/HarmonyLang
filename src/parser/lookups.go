package parser

import (
	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/lexer"
)

type BindingPower int

const (
	deflet BindingPower = iota
	comma
	assignment
	logical
	relational
	additive
	multiplicative
	unary
	call
	member
	primary
)

type statementHandler func(parser *parser) ast.Statement
type nudHandler func(parser *parser) ast.Expression
type ledHandler func(parser *parser, left ast.Expression, bp BindingPower) ast.Expression

type statementLookup map[lexer.TokenKind]statementHandler
type nudLookup map[lexer.TokenKind]nudHandler
type ledLookup map[lexer.TokenKind]ledHandler
type bindingPowerLookup map[lexer.TokenKind]BindingPower

var statementHandlers = statementLookup{}
var nudHandlers = nudLookup{}
var ledHandlers = ledLookup{}
var bindingPowers = bindingPowerLookup{}

func led(kind lexer.TokenKind, power BindingPower, handler ledHandler) {
	bindingPowers[kind] = power
	ledHandlers[kind] = handler
}

func nud(kind lexer.TokenKind, power BindingPower, handler nudHandler) {
	bindingPowers[kind] = power
	nudHandlers[kind] = handler
}

func stmt(kind lexer.TokenKind, handler statementHandler) {
	bindingPowers[kind] = deflet
	statementHandlers[kind] = handler
}

func createTokenLookups() {
}
