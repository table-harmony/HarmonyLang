package parser

import (
	"fmt"

	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/lexer"
)

type type_handler func(parser *parser) ast.Type

var type_lookup = map[lexer.TokenKind]type_handler{}

func add_type(kind lexer.TokenKind, handler type_handler) {
	type_lookup[kind] = handler
}

func create_type_lookups() {
	add_type(lexer.IDENTIFIER, parse_primitive_type)
	add_type(lexer.OPEN_BRACKET, parse_array_type)
}

func parse_type(parser *parser) ast.Type {
	token := parser.currentToken()

	handler, exists := type_lookup[token.Kind]
	if !exists {
		panic(fmt.Sprintf("No type handler for token %s", token.Kind.ToString()))
	}

	return handler(parser)
}

func parse_primitive_type(parser *parser) ast.Type {
	token := parser.currentToken()
	parser.advance(1)

	switch token.Value {
	case "number", "string", "boolean":
		return ast.PrimitiveType{Name: token.Value}
	default:
		panic(fmt.Sprintf("Unknown type: %s", token.Value))
	}
}

func parse_array_type(parser *parser) ast.Type {
	parser.advance(1)
	parser.expect(lexer.CLOSE_BRACKET)
	parser.advance(1)

	elementType := parse_type(parser)
	return ast.ArrayType{Underlying: elementType}
}
