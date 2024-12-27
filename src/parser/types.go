package parser

import (
	"fmt"

	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/lexer"
)

type type_nud_handler func(p *parser) ast.Type
type type_led_handler func(p *parser, left ast.Type, bp binding_power) ast.Type

var type_bp_lookup = map[lexer.TokenKind]binding_power{}
var type_nud_lookup = map[lexer.TokenKind]type_nud_handler{}
var type_led_lookup = map[lexer.TokenKind]type_led_handler{}

func type_led(kind lexer.TokenKind, bp binding_power, handler type_led_handler) {
	type_bp_lookup[kind] = bp
	type_led_lookup[kind] = handler
}

func type_nud(kind lexer.TokenKind, bp binding_power, handler type_nud_handler) {
	type_bp_lookup[kind] = bp
	type_nud_lookup[kind] = handler
}

func create_type_token_lookups() {
	type_nud(lexer.IDENTIFIER, primary, parse_symbol_type)

	type_led(lexer.OPEN_BRACKET, call, parse_array_type)
}

func parse_type(parser *parser, bp binding_power) ast.Type {
	token := parser.current_token()
	nud_handler, exists := type_nud_lookup[token.Kind]

	if !exists {
		panic(fmt.Sprintf("type: NUD Handler expected for token %s\n", token.Kind.ToString()))
	}

	left := nud_handler(parser)

	for type_bp_lookup[parser.current_token().Kind] > bp {
		token = parser.current_token()
		led_handler, exists := type_led_lookup[token.Kind]

		if !exists {
			panic(fmt.Sprintf("type: LED Handler expected for token %s\n", token.Kind.ToString()))
		}

		left = led_handler(parser, left, bp)
	}

	return left
}

func parse_symbol_type(parser *parser) ast.Type {
	token := parser.expect(lexer.IDENTIFIER)
	parser.advance(1)

	switch token.Value {
	case "number":
		return ast.NumberType{}
	case "bool":
		return ast.BooleanType{}
	case "string":
		return ast.StringType{}
	default:
		return ast.SymbolType{
			Value: token.Value,
		}
	}
}

func parse_array_type(parser *parser, left ast.Type, bp binding_power) ast.Type {
	parser.expect(lexer.OPEN_BRACKET)
	parser.advance(1)
	parser.expect(lexer.CLOSE_BRACKET)
	parser.advance(1)

	return ast.ArrayType{
		Underlying: left,
	}
}
