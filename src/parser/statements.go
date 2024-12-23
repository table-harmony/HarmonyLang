package parser

import (
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
	isConstant := token.Kind == lexer.CONST
	parser.advance(1)

	identifier := parser.expect(lexer.IDENTIFIER).Value
	parser.advance(1)

	var explicitType ast.Type
	if parser.currentToken().Kind == lexer.COLON {
		parser.advance(1)
		explicitType = parse_type(parser)
	}

	var value ast.Expression
	if parser.currentToken().Kind == lexer.ASSIGNMENT {
		parser.advance(1)
		value = parse_expression(parser, assignment)
	} else if explicitType == nil {
		panic("Cannot define a variable without an explicit type or default value.")
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
