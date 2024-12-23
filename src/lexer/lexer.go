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

// The lexer
type lexer struct {
	patterns []regex_pattern
	Tokens   []Token
	source   string
	pos      int
}

// Returns a reference to a new lexer
func createLexer(source string) *lexer {
	return &lexer{
		pos:      0,
		source:   source,
		Tokens:   make([]Token, 0),
		patterns: reserved_patterns,
	}
}

// Advances the lexer position by n
func (lexer *lexer) advance(n int) {
	lexer.pos += n
}

// Returns the lexer current value at the source
func (lexer *lexer) at() byte {
	return lexer.source[lexer.pos]
}

// Returns the remainder of the source
func (lexer *lexer) remainder() string {
	return lexer.source[lexer.pos:]
}

// Pushes a token onto the lexer
func (lexer *lexer) push(token Token) {
	lexer.Tokens = append(lexer.Tokens, token)
}

// Returns whether the lexer is at the end of the source
func (lexer *lexer) at_eof() bool {
	return lexer.pos >= len(lexer.source)
}
