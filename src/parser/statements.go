package parser

import (
	"fmt"

	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/lexer"
)

func parse_statement(parser *parser) ast.Statement {
	token := parser.currentToken()
	handler, exists := statement_lookup[token.Kind]

	if exists {
		return handler(parser)
	}

	return parse_expression_statement(parser)
}

func parse_expression_statement(parser *parser) ast.ExpressionStatement {
	expression := parse_expression(parser, default_bp)

	parser.expect(lexer.SEMI_COLON)
	parser.advance(1)

	return ast.ExpressionStatement{
		Expression: expression,
	}
}

func parse_variable_declaration_statement(parser *parser) ast.Statement {
	token := parser.currentToken()

	var explicitType ast.Type = nil

	isConstant := token.Kind == lexer.CONST
	parser.advance(1)
	token = parser.currentToken()

	parser.expectError(lexer.IDENTIFIER,
		fmt.Sprintf("Following %s expected variable name however instead recieved %s instead\n",
			token.Kind.ToString(), token.Kind.ToString()))

	identifier := parser.currentToken().Value
	parser.advance(1)
	token = parser.currentToken()

	if token.Kind == lexer.COLON {
		parser.expect(lexer.COLON)
		explicitType = nil //TODO: parse type
	}

	var value ast.Expression
	if token.Kind != lexer.SEMI_COLON {
		parser.expect(lexer.ASSIGNMENT)
		parser.advance(1)

		value = parse_expression(parser, assignment)
	}

	parser.expect(lexer.SEMI_COLON)
	parser.advance(1)

	if isConstant && value == nil {
		panic("Cannot define constant variable without providing default value.")
	}

	return ast.VariableDeclarationStatement{
		Identifier:   identifier,
		IsConstant:   isConstant,
		Value:        value,
		ExplicitType: explicitType,
	}
}
