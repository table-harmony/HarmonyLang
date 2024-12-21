package lexer

import (
	"regexp"
)

// A function which handles a given regex and deals with tokenization for the lexer
type regexHandler func(lex *lexer, regex *regexp.Regexp)

// A regex pattern including a pattern and a handler function
type regexPattern struct {
	regex   *regexp.Regexp
	handler regexHandler
}

// The lexer
type lexer struct {
	patterns []regexPattern
	Tokens   []Token
	source   string
	pos      int
}

func createLexer(source string) *lexer {
	return &lexer{
		pos:      0,
		source:   source,
		Tokens:   make([]Token, 0),
		patterns: reservedPatterns,
	}
}

func (lex *lexer) advance(n int) {
	lex.pos += n
}

func (lex *lexer) at() byte {
	return lex.source[lex.pos]
}

func (lex *lexer) remainder() string {
	return lex.source[lex.pos:]
}

func (lex *lexer) push(token Token) {
	lex.Tokens = append(lex.Tokens, token)
}

func (lex *lexer) at_eof() bool {
	return lex.pos >= len(lex.source)
}
