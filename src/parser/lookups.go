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
type sed_handler func(parser *parser, left ast.Expression) ast.Statement
type led_handler func(parser *parser, left ast.Expression, bp binding_power) ast.Expression

var binding_power_lookup = map[lexer.TokenKind]binding_power{}
var nud_lookup = map[lexer.TokenKind]nud_handler{}
var led_lookup = map[lexer.TokenKind]led_handler{}
var sed_lookup = map[lexer.TokenKind]sed_handler{}
var statement_lookup = map[lexer.TokenKind]statement_handler{}

func register_led(kind lexer.TokenKind, bp binding_power, handler led_handler) {
	binding_power_lookup[kind] = bp
	led_lookup[kind] = handler
}

func register_nud(kind lexer.TokenKind, bp binding_power, handler nud_handler) {
	binding_power_lookup[kind] = bp
	nud_lookup[kind] = handler
}

func register_sed(kind lexer.TokenKind, handler sed_handler) {
	sed_lookup[kind] = handler
}

func register_statement(kind lexer.TokenKind, handler statement_handler) {
	binding_power_lookup[kind] = default_bp
	statement_lookup[kind] = handler
}

func create_token_lookups() {
	// Assignment
	register_sed(lexer.ASSIGNMENT, parse_assignment_statement)
	register_sed(lexer.PLUS_EQUALS, parse_assignment_statement)
	register_sed(lexer.MINUS_EQUALS, parse_assignment_statement)
	register_sed(lexer.STAR_EQUALS, parse_assignment_statement)
	register_sed(lexer.SLASH_EQUALS, parse_assignment_statement)
	register_sed(lexer.PERCENT_EQUALS, parse_assignment_statement)
	register_sed(lexer.AND_EQUALS, parse_assignment_statement)
	register_sed(lexer.OR_EQUALS, parse_assignment_statement)
	register_sed(lexer.NULLISH_ASSIGNMENT, parse_assignment_statement)
	register_sed(lexer.PLUS_PLUS, parse_assignment_statement)
	register_sed(lexer.MINUS_MINUS, parse_assignment_statement)

	// Logical
	register_led(lexer.AND, logical, parse_binary_expression)
	register_led(lexer.OR, logical, parse_binary_expression)
	register_led(lexer.DOT_DOT, logical, parse_range_expression)

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
	register_nud(lexer.NIL, primary, parse_primary_expression)

	// Data types
	register_nud(lexer.MAP, default_bp, parse_map_instantiation_expression)
	register_nud(lexer.OPEN_BRACKET, default_bp, parse_array_instantiation_expression)

	// Unary / Prefix
	register_nud(lexer.DASH, additive, parse_prefix_expression) // binding power of additive sense a dash or a plus as unary are the same as additive operations
	register_nud(lexer.PLUS, additive, parse_prefix_expression) // making them unary would cause errors because they would have higher precedence than multiplicative
	register_nud(lexer.NOT, unary, parse_prefix_expression)
	register_nud(lexer.AMPERSAND, unary, parse_prefix_expression)
	register_nud(lexer.STAR, unary, parse_prefix_expression)
	register_nud(lexer.TYPEOF, unary, parse_prefix_expression)

	// Ternary
	register_led(lexer.QUESTION, ternary, parse_ternary_expression)

	// Grouping Expression
	register_nud(lexer.OPEN_PAREN, default_bp, parse_grouping_expression)

	// Member / Computed / Call
	register_led(lexer.OPEN_PAREN, call, parse_call_expression)
	register_led(lexer.DOT, member, parse_member_expression)
	register_led(lexer.OPEN_BRACKET, member, parse_computed_member_expression)

	// Block
	register_nud(lexer.OPEN_CURLY, default_bp, parse_block_expression)
	register_nud(lexer.IF, default_bp, parse_if_expression)
	register_nud(lexer.SWITCH, default_bp, parse_switch_expression)
	register_nud(lexer.FN, default_bp, parse_function_declaration_expression)
	register_nud(lexer.TRY, default_bp, parse_try_catch_expression)

	// Struct instantiation
	register_led(lexer.OPEN_CURLY, call, parse_struct_instantiation_expression)

	// Statements
	register_statement(lexer.TYPE, parse_type_declaration_statement)
	register_statement(lexer.IMPORT, parse_import_statement)
	register_statement(lexer.LET, parse_multi_variable_declaration_statement)
	register_statement(lexer.CONST, parse_multi_variable_declaration_statement)
	register_statement(lexer.INTERFACE, parse_interface_declaration_statement)
	register_statement(lexer.STRUCT, parse_struct_declaration_statement)
	register_statement(lexer.FN, parse_function_declaration_statement)
	register_statement(lexer.FOR, parse_for_statement)
	register_statement(lexer.CONTINUE, parse_loop_control_statement)
	register_statement(lexer.BREAK, parse_loop_control_statement)
	register_statement(lexer.RETURN, parse_return_statement)
	register_statement(lexer.THROW, parse_throw_statement)
}
