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

func led(kind lexer.TokenKind, bp binding_power, handler led_handler) {
	binding_power_lookup[kind] = bp
	led_lookup[kind] = handler
}

func nud(kind lexer.TokenKind, bp binding_power, handler nud_handler) {
	binding_power_lookup[kind] = bp
	nud_lookup[kind] = handler
}

func statement(kind lexer.TokenKind, handler statement_handler) {
	binding_power_lookup[kind] = default_bp
	statement_lookup[kind] = handler
}

func create_token_lookups() {
	// Assignment
	led(lexer.ASSIGNMENT, assignment, parse_assignment_expression)
	led(lexer.PLUS_EQUALS, assignment, parse_assignment_expression)
	led(lexer.MINUS_EQUALS, assignment, parse_assignment_expression)
	led(lexer.STAR_EQUALS, assignment, parse_assignment_expression)
	led(lexer.SLASH_EQUALS, assignment, parse_assignment_expression)

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

	// Additive
	led(lexer.PLUS, additive, parse_binary_expression)
	led(lexer.DASH, additive, parse_binary_expression)

	// Multiplicative
	led(lexer.SLASH, multiplicative, parse_binary_expression)
	led(lexer.STAR, multiplicative, parse_binary_expression)
	led(lexer.PERCENT, multiplicative, parse_binary_expression)

	// Literals & Symbols
	nud(lexer.NUMBER, primary, parse_primary_expression)
	nud(lexer.STRING, primary, parse_primary_expression)
	nud(lexer.IDENTIFIER, primary, parse_primary_expression)

	// Unary/Prefix
	nud(lexer.DASH, unary, parse_prefix_expression)
	nud(lexer.NOT, unary, parse_prefix_expression)
	nud(lexer.TYPEOF, unary, parse_prefix_expression)

	// Member / Computed // Call

	// Grouping Expression
	nud(lexer.OPEN_PAREN, default_bp, parse_grouping_expression)
	nud(lexer.SWITCH, default_bp, parse_switch_expression)

	// Modifiers
	//statement(lexer.STATIC, parse_modifier_statement)

	// Statements
	statement(lexer.LET, parse_variable_declaration_statement)
	statement(lexer.CONST, parse_variable_declaration_statement)
	statement(lexer.INTERFACE, parse_interface_declaration_statement)
	statement(lexer.STRUCT, parse_struct_declaration_statement)
	statement(lexer.FUNC, parse_func_declaration_statement)
	statement(lexer.IMPORT, parse_import_statement)
	statement(lexer.IF, parse_if_statement)
	statement(lexer.OPEN_CURLY, parse_block_statement)
	statement(lexer.SWITCH, parse_switch_statement)
}
