package lexer

import (
	"regexp"
)

// A function which handles a given regex and deals with tokenization for the lexer
type regex_handler func(lex *lexer, regex *regexp.Regexp)

// A regex pattern including a pattern and a handler function
type regex_pattern struct {
	regex   *regexp.Regexp
	handler regex_handler
}

type lexer struct {
	patterns []regex_pattern
	Tokens   []Token
	source   string
	pos      int
}

func createLexer(source string) *lexer {
	return &lexer{
		pos:      0,
		source:   source,
		Tokens:   make([]Token, 0),
		patterns: reserved_patterns,
	}
}

func (lexer *lexer) advance(n int) {
	lexer.pos += n
}

func (lexer *lexer) at() byte {
	return lexer.source[lexer.pos]
}

func (lexer *lexer) remainder() string {
	return lexer.source[lexer.pos:]
}

func (lexer *lexer) push(token Token) {
	lexer.Tokens = append(lexer.Tokens, token)
}

func (lexer *lexer) at_eof() bool {
	return lexer.pos >= len(lexer.source)
}
