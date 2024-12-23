package parser

import (
	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/lexer"
)

type binding_power int

const (
	default_bp binding_power = iota
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

type statement_handler func(parser *parser) ast.Statement
type nud_handler func(parser *parser) ast.Expression
type led_handler func(parser *parser, left ast.Expression, bp binding_power) ast.Expression

var binding_power_lookup = map[lexer.TokenKind]binding_power{}
var nud_lookup = map[lexer.TokenKind]nud_handler{}
var led_lookup = map[lexer.TokenKind]led_handler{}
var statement_lookup = map[lexer.TokenKind]statement_handler{}

func add_led(kind lexer.TokenKind, bp binding_power, handler led_handler) {
	binding_power_lookup[kind] = bp
	led_lookup[kind] = handler
}

func add_nud(kind lexer.TokenKind, bp binding_power, handler nud_handler) {
	binding_power_lookup[kind] = bp
	nud_lookup[kind] = handler
}

func add_statement(kind lexer.TokenKind, handler statement_handler) {
	binding_power_lookup[kind] = default_bp
	statement_lookup[kind] = handler
}

func create_token_lookups() {
	// Assignment

	// Logical
	add_led(lexer.AND, logical, parse_binary_expression)
	add_led(lexer.OR, logical, parse_binary_expression)

	// Relational
	add_led(lexer.LESS, relational, parse_binary_expression)
	add_led(lexer.LESS_EQUALS, relational, parse_binary_expression)
	add_led(lexer.GREATER, relational, parse_binary_expression)
	add_led(lexer.GREATER_EQUALS, relational, parse_binary_expression)
	add_led(lexer.EQUALS, relational, parse_binary_expression)
	add_led(lexer.NOT_EQUALS, relational, parse_binary_expression)

	// Additive & Multiplicitave
	add_led(lexer.PLUS, additive, parse_binary_expression)
	add_led(lexer.DASH, additive, parse_binary_expression)
	add_led(lexer.SLASH, multiplicative, parse_binary_expression)
	add_led(lexer.STAR, multiplicative, parse_binary_expression)
	add_led(lexer.PERCENT, multiplicative, parse_binary_expression)

	// Literals & Symbols
	add_nud(lexer.NUMBER, primary, parse_primary_expression)
	add_nud(lexer.STRING, primary, parse_primary_expression)
	add_nud(lexer.IDENTIFIER, primary, parse_primary_expression)

	// Unary/Prefix

	// Member / Computed // Call

	// Grouping Expression
}
