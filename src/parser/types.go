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

func register_type_led(kind lexer.TokenKind, bp binding_power, handler type_led_handler) {
	type_bp_lookup[kind] = bp
	type_led_lookup[kind] = handler
}

func register_type_nud(kind lexer.TokenKind, bp binding_power, handler type_nud_handler) {
	type_bp_lookup[kind] = bp
	type_nud_lookup[kind] = handler
}

func create_type_token_lookups() {
	// Primitive types
	register_type_nud(lexer.IDENTIFIER, primary, parse_symbol_type)
	register_type_nud(lexer.NIL, primary, parse_nil_type)

	// Pointer
	register_type_nud(lexer.STAR, unary, parse_pointer_type)

	// Data types
	register_type_nud(lexer.MAP, primary, parse_map_type)
	register_type_nud(lexer.OPEN_BRACKET, member, parse_array_type)

	// Function types
	register_type_nud(lexer.FN, primary, parse_function_type)
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
	case "error":
		return ast.ErrorType{}
	case "any":
		return ast.AnyType{}
	default:
		return ast.SymbolType{
			Value: token.Value,
		}
	}
}

func parse_nil_type(parser *parser) ast.Type {
	parser.expect(lexer.NIL)
	parser.advance(1)

	return ast.NilType{}
}

func parse_array_type(parser *parser) ast.Type {
	parser.expect(lexer.OPEN_BRACKET)
	parser.advance(1)

	if parser.current_token().Kind == lexer.CLOSE_BRACKET {
		parser.advance(1)
		return ast.SliceType{
			Underlying: parse_type(parser, default_bp),
		}
	}

	size := parse_expression(parser, default_bp)

	parser.expect(lexer.CLOSE_BRACKET)
	parser.advance(1)

	return ast.ArrayType{
		Size:       size,
		Underlying: parse_type(parser, default_bp),
	}
}

func parse_map_type(parser *parser) ast.Type {
	parser.expect(lexer.MAP)
	parser.advance(1)

	if parser.current_token().Kind != lexer.OPEN_BRACKET {
		return ast.MapType{
			Key:   ast.AnyType{},
			Value: ast.AnyType{},
		}
	}

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

func parse_pointer_type(parser *parser) ast.Type {
	parser.expect(lexer.STAR)
	parser.advance(1)

	return ast.PointerType{
		Target: parse_type(parser, unary),
	}
}
