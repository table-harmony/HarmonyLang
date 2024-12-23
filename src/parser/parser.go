package parser

import (
	"fmt"

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
		statement := parse_statement(parser)
		body = append(body, statement)
	}

	return ast.BlockStatement{
		Body: body,
	}
}

func createParser(tokens []lexer.Token) *parser {
	create_token_lookups()
	create_type_lookups()

	return &parser{
		tokens: tokens,
		pos:    0,
	}
}

func (parser *parser) currentToken() lexer.Token {
	return parser.tokens[parser.pos]
}

func (parser *parser) previousToken() lexer.Token {
	return parser.tokens[parser.pos-1]
}

func (parser *parser) advance(n int) {
	parser.pos += n
}

func (parser *parser) isEmpty() bool {
	currentToken := parser.currentToken()

	return parser.pos >= len(parser.tokens) ||
		currentToken.Kind == lexer.EOF
}

func (parser *parser) expectError(expectedKind lexer.TokenKind, err any) lexer.Token {
	currentToken := parser.currentToken()

	if currentToken.Kind != expectedKind {
		if err == nil {
			err = fmt.Sprintf("Expected %s but recieved %s instead\n",
				expectedKind.ToString(), currentToken.Kind.ToString())
		}

		panic(err)
	}

	return currentToken
}

func (parser *parser) expect(expectedKind lexer.TokenKind) lexer.Token {
	return parser.expectError(expectedKind, nil)
}
