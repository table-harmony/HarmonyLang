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

func Parse(tokens []lexer.Token) []ast.Statement {
	statements := make([]ast.Statement, 0)
	parser := create_parser(tokens)

	for !parser.is_empty() {
		statement := parse_statement(parser)
		statements = append(statements, statement)
	}

	return statements
}

func create_parser(tokens []lexer.Token) *parser {
	create_token_lookups()
	create_type_token_lookups()

	return &parser{
		tokens: tokens,
		pos:    0,
	}
}

func (parser *parser) current_token() lexer.Token {
	return parser.tokens[parser.pos]
}

func (parser *parser) previous_token() lexer.Token {
	return parser.tokens[parser.pos-1]
}

func (parser *parser) advance(n int) {
	parser.pos += n
}

func (parser *parser) is_empty() bool {
	currentToken := parser.current_token()

	return parser.pos >= len(parser.tokens) ||
		currentToken.Kind == lexer.EOF
}

func (parser *parser) expect_error(expectedKind lexer.TokenKind, err any) lexer.Token {
	currentToken := parser.current_token()

	if currentToken.Kind != expectedKind {
		if err == nil {
			err = fmt.Sprintf("Expected %s but recieved %s instead\n",
				expectedKind.ToString(), currentToken.Kind.ToString())
		}

		panic(err)
	}

	return currentToken
}

func (parser *parser) expect(expected lexer.TokenKind) lexer.Token {
	return parser.expect_error(expected, nil)
}
