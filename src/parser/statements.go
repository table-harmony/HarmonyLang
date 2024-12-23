package parser

import (
	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/lexer"
)

func parse_statement(parser *parser) ast.Statement {
	currentToken := parser.currentToken()
	handler, exists := statement_lookup[currentToken.Kind]

	if exists {
		return handler(parser)
	}

	return parse_expression_statement(parser)
}

func parse_expression_statement(parser *parser) ast.ExpressionStatement {
	expression := parse_expression(parser, defalt_bp)
	parser.expect(lexer.SEMI_COLON)

	return ast.ExpressionStatement{
		Expression: expression,
	}
}
