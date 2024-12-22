package parser

import (
	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/lexer"
)

type parser struct {
	tokens []lexer.Token
	pos    int
}

func Parse(tokens []lexer.Token) ast.BlockStatement {
	body := make([]ast.Statement, 0)

	parser := createParser(tokens)

	for !parser.isEmpty() {
		statement := parseStatement(parser)

		body = append(body, statement)
	}

	return ast.BlockStatement{
		Body: body,
	}
}

func createParser(tokens []lexer.Token) *parser {
	createTokenLookups()

	return &parser{
		tokens: tokens,
		pos:    0,
	}
}

func (parser *parser) currentToken() lexer.Token {
	return parser.tokens[parser.pos]
}

func (parser *parser) advance(n int) {
	parser.pos += n
}

func (parser *parser) isEmpty() bool {
	currentToken := parser.currentToken()

	return parser.pos > len(parser.tokens) ||
		currentToken.Kind == lexer.EOF
}
