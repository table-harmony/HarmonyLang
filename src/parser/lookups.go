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
	ternary
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

func register_led(kind lexer.TokenKind, bp binding_power, handler led_handler) {
	binding_power_lookup[kind] = bp
	led_lookup[kind] = handler
}

func register_nud(kind lexer.TokenKind, bp binding_power, handler nud_handler) {
	binding_power_lookup[kind] = bp
	nud_lookup[kind] = handler
}

func register_statement(kind lexer.TokenKind, handler statement_handler) {
	binding_power_lookup[kind] = default_bp
	statement_lookup[kind] = handler
}

func create_token_lookups() {
	// Assignment
	register_led(lexer.ASSIGNMENT, assignment, parse_assignment_expression)
	register_led(lexer.PLUS_EQUALS, assignment, parse_assignment_expression)
	register_led(lexer.MINUS_EQUALS, assignment, parse_assignment_expression)
	register_led(lexer.STAR_EQUALS, assignment, parse_assignment_expression)
	register_led(lexer.SLASH_EQUALS, assignment, parse_assignment_expression)
	register_led(lexer.PLUS_PLUS, assignment, parse_assignment_expression)
	register_led(lexer.MINUS_MINUS, assignment, parse_assignment_expression)

	// Logical
	register_led(lexer.AND, logical, parse_binary_expression)
	register_led(lexer.OR, logical, parse_binary_expression)

	// Relational
	register_led(lexer.LESS, relational, parse_binary_expression)
	register_led(lexer.LESS_EQUALS, relational, parse_binary_expression)
	register_led(lexer.GREATER, relational, parse_binary_expression)
	register_led(lexer.GREATER_EQUALS, relational, parse_binary_expression)
	register_led(lexer.EQUALS, relational, parse_binary_expression)
	register_led(lexer.NOT_EQUALS, relational, parse_binary_expression)
	register_led(lexer.IN, relational, parse_binary_expression)

	// Additive
	register_led(lexer.PLUS, additive, parse_binary_expression)
	register_led(lexer.DASH, additive, parse_binary_expression)

	// Multiplicative
	register_led(lexer.SLASH, multiplicative, parse_binary_expression)
	register_led(lexer.STAR, multiplicative, parse_binary_expression)
	register_led(lexer.PERCENT, multiplicative, parse_binary_expression)

	// Literals & Symbols
	register_nud(lexer.NUMBER, primary, parse_primary_expression)
	register_nud(lexer.STRING, primary, parse_primary_expression)
	register_nud(lexer.IDENTIFIER, primary, parse_primary_expression)
	register_nud(lexer.TRUE, primary, parse_primary_expression)
	register_nud(lexer.FALSE, primary, parse_primary_expression)

	// Unary/Prefix
	//TODO: fix this shit -> 10 * 10 - 10 (- has more precidecne because it is unary)
	//register_nud(lexer.DASH, unary, parse_prefix_expression)
	register_nud(lexer.NOT, unary, parse_prefix_expression)
	register_nud(lexer.TYPEOF, unary, parse_prefix_expression)

	// Ternary
	register_led(lexer.QUESTION, ternary, parse_ternary_expression)

	// Member / Computed // Call

	// Grouping Expression
	register_nud(lexer.OPEN_PAREN, default_bp, parse_grouping_expression)
	register_nud(lexer.SWITCH, default_bp, parse_switch_expression)

	// Statements
	register_statement(lexer.LET, parse_variable_declaration_statement)
	register_statement(lexer.CONST, parse_variable_declaration_statement)
	register_statement(lexer.INTERFACE, parse_interface_declaration_statement)
	register_statement(lexer.STRUCT, parse_struct_declaration_statement)
	register_statement(lexer.FUNC, parse_func_declaration_statement)
	register_statement(lexer.IMPORT, parse_import_statement)
	register_statement(lexer.IF, parse_if_statement)
	register_statement(lexer.OPEN_CURLY, parse_block_statement)
	register_statement(lexer.SWITCH, parse_switch_statement)
	register_statement(lexer.FOR, parse_for_statement)
	register_statement(lexer.CONTINUE, parse_loop_control_statement)
	register_statement(lexer.BREAK, parse_loop_control_statement)
}
