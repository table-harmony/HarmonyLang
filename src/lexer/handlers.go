package lexer

import (
	"regexp"
)

func default_handler(kind TokenKind, value string) regex_handler {
	return func(lex *lexer, _ *regexp.Regexp) {
		lex.advance(len(value))
		lex.push(CreateToken(kind, value))
	}
}

func string_handler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	stringLiteral := lex.remainder()[match[0]+1 : match[1]-1]

	lex.push(CreateToken(STRING, stringLiteral))
	lex.advance(len(stringLiteral) + 2)
}

func number_handler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.push(CreateToken(NUMBER, match))
	lex.advance(len(match))
}

func symbol_handler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())

	if kind, found := reserved_keywords[match]; found {
		lex.push(CreateToken(kind, match))
	} else {
		lex.push(CreateToken(IDENTIFIER, match))
	}

	lex.advance(len(match))
}

func skip_handler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.advance(len(match))
}

func comment_handler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())

	if match != nil {
		lex.advance(match[1])
	}
}

func newline_handler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())

	if len(lex.Tokens) > 0 {
		last := lex.peek()

		if needs_semi_colon(last) {
			lex.push(CreateToken(SEMI_COLON, ";"))
		}
	}

	lex.line++
	lex.advance(len(match))
}

// TODO: semi colon insertion hurts stuff like switch statements either
// TODO: i need to specify to disregard them on switch statements and more things in the future
// TODO: or just remove it and think of something else

func needs_semi_colon(token Token) bool {
	return token.IsOfKind(
		IDENTIFIER,
		NUMBER,
		STRING,
		BREAK,
		CONTINUE,
		RETURN,
		TRUE,
		FALSE,
		NULL,
		CLOSE_PAREN,
		CLOSE_BRACKET,
		CLOSE_CURLY,
	)
}

// reserved regex patterns
var reserved_patterns = []regex_pattern{
	{regexp.MustCompile(`\r\n|\r|\n`), newline_handler},
	{regexp.MustCompile(`[ \t]+`), skip_handler},
	{regexp.MustCompile(`\/\/.*`), comment_handler},
	{regexp.MustCompile(`"[^"]*"`), string_handler},
	{regexp.MustCompile(`[0-9]+(\.[0-9]+)?`), number_handler},
	{regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`), symbol_handler},

	// Grouping & Braces
	{regexp.MustCompile(`\[`), default_handler(OPEN_BRACKET, "[")},
	{regexp.MustCompile(`\]`), default_handler(CLOSE_BRACKET, "]")},
	{regexp.MustCompile(`\{`), default_handler(OPEN_CURLY, "{")},
	{regexp.MustCompile(`\}`), default_handler(CLOSE_CURLY, "}")},
	{regexp.MustCompile(`\(`), default_handler(OPEN_PAREN, "(")},
	{regexp.MustCompile(`\)`), default_handler(CLOSE_PAREN, ")")},

	// Equivalence
	{regexp.MustCompile(`==`), default_handler(EQUALS, "==")},
	{regexp.MustCompile(`!=`), default_handler(NOT_EQUALS, "!=")},
	{regexp.MustCompile(`=`), default_handler(ASSIGNMENT, "=")},
	{regexp.MustCompile(`!`), default_handler(NOT, "!")},

	// Conditional
	{regexp.MustCompile(`<=`), default_handler(LESS_EQUALS, "<=")},
	{regexp.MustCompile(`<`), default_handler(LESS, "<")},
	{regexp.MustCompile(`>=`), default_handler(GREATER_EQUALS, ">=")},
	{regexp.MustCompile(`>`), default_handler(GREATER, ">")},

	// Logical
	{regexp.MustCompile(`\|\|`), default_handler(OR, "||")},
	{regexp.MustCompile(`&&`), default_handler(AND, "&&")},

	// Symbols
	{regexp.MustCompile(`\.\.`), default_handler(DOT_DOT, "..")},
	{regexp.MustCompile(`\.`), default_handler(DOT, ".")},
	{regexp.MustCompile(`;`), default_handler(SEMI_COLON, ";")},
	{regexp.MustCompile(`:`), default_handler(COLON, ":")},
	{regexp.MustCompile(`\?\?=`), default_handler(NULLISH_ASSIGNMENT, "??=")},
	{regexp.MustCompile(`\?`), default_handler(QUESTION, "?")},
	{regexp.MustCompile(`,`), default_handler(COMMA, ",")},
	{regexp.MustCompile(`->`), default_handler(ARROW, "->")},

	// Shorthand
	{regexp.MustCompile(`\+\+`), default_handler(PLUS_PLUS, "++")},
	{regexp.MustCompile(`--`), default_handler(MINUS_MINUS, "--")},
	{regexp.MustCompile(`\+=`), default_handler(PLUS_EQUALS, "+=")},
	{regexp.MustCompile(`-=`), default_handler(MINUS_EQUALS, "-=")},
	{regexp.MustCompile(`\*=`), default_handler(STAR_EQUALS, "*=")},
	{regexp.MustCompile(`/=`), default_handler(SLASH_EQUALS, "/=")},
	{regexp.MustCompile(`%=`), default_handler(PERCENT_EQUALS, "%=")},

	// Math Operators
	{regexp.MustCompile(`\+`), default_handler(PLUS, "+")},
	{regexp.MustCompile(`-`), default_handler(DASH, "-")},
	{regexp.MustCompile(`/`), default_handler(SLASH, "/")},
	{regexp.MustCompile(`\*`), default_handler(STAR, "*")},
	{regexp.MustCompile(`%`), default_handler(PERCENT, "%")},
}
