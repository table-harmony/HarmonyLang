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
	// Primitive types
	type_nud(lexer.IDENTIFIER, primary, parse_symbol_type)

	// Data types
	type_nud(lexer.MAP, primary, parse_map_type)
	type_led(lexer.OPEN_BRACKET, member, parse_array_type_suffix)

	// Function types
	type_nud(lexer.FN, primary, parse_function_type)
}

func parse_type(parser *parser, bp binding_power) ast.Type {
	token := parser.current_token()
	nud_handler, exists := type_nud_lookup[token.Kind]

	if !exists {
		panic(fmt.Sprintf("type: NUD Handler expected for token %s\n", token.Kind.String()))
	}

	left := nud_handler(parser)

	for !parser.is_empty() && type_bp_lookup[parser.current_token().Kind] > bp {
		token = parser.current_token()
		led_handler, exists := type_led_lookup[token.Kind]

		if !exists {
			panic(fmt.Sprintf("type: LED Handler expected for token %s\n", token.Kind.String()))
		}

		left = led_handler(parser, left, type_bp_lookup[token.Kind])
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

func parse_array_type_suffix(parser *parser, underlying ast.Type, bp binding_power) ast.Type {
	parser.expect(lexer.OPEN_BRACKET)
	parser.advance(1)

	if parser.current_token().Kind == lexer.CLOSE_BRACKET {
		parser.advance(1)
		return ast.SliceType{
			Underlying: underlying,
		}
	}

	size := parse_expression(parser, default_bp)

	parser.expect(lexer.CLOSE_BRACKET)
	parser.advance(1)

	return ast.ArrayType{
		Size:       size,
		Underlying: underlying,
	}
}

func parse_map_type(parser *parser) ast.Type {
	parser.expect(lexer.MAP)
	parser.advance(1)

	parser.expect(lexer.OPEN_BRACKET)
	parser.advance(1)

	keyType := parse_type(parser, default_bp)

	parser.expect(lexer.ARROW)
	parser.advance(1)

	valueType := parse_type(parser, default_bp)

	parser.expect(lexer.CLOSE_BRACKET)
	parser.advance(1)

	return ast.MapType{
		Key:   keyType,
		Value: valueType,
	}
}

func parse_function_type(parser *parser) ast.Type {
	parser.expect(lexer.FN)
	parser.advance(1)

	parser.expect(lexer.OPEN_PAREN)
	parser.advance(1)

	var parameters []ast.Parameter
	for parser.current_token().Kind != lexer.CLOSE_PAREN {
		paramIdentifier := parser.expect(lexer.IDENTIFIER).Value
		parser.advance(1)

		var paramType ast.Type
		if parser.current_token().Kind == lexer.COLON {
			parser.expect(lexer.COLON)
			parser.advance(1)

			paramType = parse_type(parser, default_bp)
		}

		parameters = append(parameters, ast.Parameter{
			Name: paramIdentifier,
			Type: paramType,
		})

		if parser.current_token().Kind == lexer.COMMA {
			parser.advance(1)
		}
	}

	parser.expect(lexer.CLOSE_PAREN)
	parser.advance(1)

	var returnType ast.Type = nil
	if parser.current_token().Kind == lexer.ARROW {
		parser.advance(1)
		returnType = parse_type(parser, default_bp)
	}

	return ast.FunctionType{
		Parameters: parameters,
		Return:     returnType,
	}
}
