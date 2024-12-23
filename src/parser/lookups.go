package parser

import (
	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/lexer"
)

type binding_power int

const (
	defalt_bp binding_power = iota
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

type stmt_handler func(p *parser) ast.Statement
type nud_handler func(p *parser) ast.Expression
type led_handler func(p *parser, left ast.Expression, bp binding_power) ast.Expression

var binding_power_lookup = map[lexer.TokenKind]binding_power{}
var nud_lookup = map[lexer.TokenKind]nud_handler{}
var led_lookup = map[lexer.TokenKind]led_handler{}
var statement_lookup = map[lexer.TokenKind]stmt_handler{}

func led(kind lexer.TokenKind, bp binding_power, handler led_handler) {
	binding_power_lookup[kind] = bp
	led_lookup[kind] = handler
}

func nud(kind lexer.TokenKind, bp binding_power, handler nud_handler) {
	binding_power_lookup[kind] = bp
	nud_lookup[kind] = handler
}

func statement(kind lexer.TokenKind, handler stmt_handler) {
	binding_power_lookup[kind] = defalt_bp
	statement_lookup[kind] = handler
}

func createTokenLookups() {
	// Assignment

	// Logical
	led(lexer.AND, logical, parse_binary_expression)
	led(lexer.OR, logical, parse_binary_expression)

	// Relational
	led(lexer.LESS, relational, parse_binary_expression)
	led(lexer.LESS_EQUALS, relational, parse_binary_expression)
	led(lexer.GREATER, relational, parse_binary_expression)
	led(lexer.GREATER_EQUALS, relational, parse_binary_expression)
	led(lexer.EQUALS, relational, parse_binary_expression)
	led(lexer.NOT_EQUALS, relational, parse_binary_expression)

	// Additive & Multiplicitave
	led(lexer.PLUS, additive, parse_binary_expression)
	led(lexer.DASH, additive, parse_binary_expression)
	led(lexer.SLASH, multiplicative, parse_binary_expression)
	led(lexer.STAR, multiplicative, parse_binary_expression)
	led(lexer.PERCENT, multiplicative, parse_binary_expression)

	// Literals & Symbols
	nud(lexer.NUMBER, primary, parse_primary_expression)
	nud(lexer.STRING, primary, parse_primary_expression)
	nud(lexer.IDENTIFIER, primary, parse_primary_expression)

	// Unary/Prefix

	// Member / Computed // Call

	// Grouping Expression
}
