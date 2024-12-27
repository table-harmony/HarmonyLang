package lexer

import (
	"regexp"
)

type regex_handler func(lex *lexer, regex *regexp.Regexp)

type regex_pattern struct {
	regex   *regexp.Regexp
	handler regex_handler
}

type lexer struct {
	patterns []regex_pattern
	Tokens   []Token
	source   string
	pos      int
	line     int
}

func create_lexer(source string) *lexer {
	return &lexer{
		patterns: reserved_patterns,
		Tokens:   make([]Token, 0),
		source:   source,
		pos:      0,
		line:     1,
	}
}

func (lex *lexer) advance(n int) {
	lex.pos += n
}

func (lexer *lexer) at() byte {
	return lexer.source[lexer.pos]
}

func (lexer *lexer) remainder() string {
	return lexer.source[lexer.pos:]
}

func (lexer *lexer) push(token Token) {
	token.Line = lexer.line
	lexer.Tokens = append(lexer.Tokens, token)
}

func (lexer *lexer) peek() Token {
	if len(lexer.Tokens) == 0 {
		panic("Lexer hasn't handled any tokens yet")
	}

	return lexer.Tokens[len(lexer.Tokens)-1]
}

func (lexer *lexer) at_eof() bool {
	return lexer.pos >= len(lexer.source)
}
