package lexer

import (
	"regexp"
)

func default_handler(kind TokenKind, value string) regex_handler {
	return func(lex *lexer, _ *regexp.Regexp) {
		lex.advance(len(value))
		lex.push(NewToken(kind, value))
	}
}

func string_handler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	stringLiteral := lex.remainder()[match[0]+1 : match[1]-1]

	lex.push(NewToken(STRING, stringLiteral))
	lex.advance(len(stringLiteral) + 2)
}

func number_handler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.push(NewToken(NUMBER, match))
	lex.advance(len(match))
}

func symbol_handler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())

	if kind, found := reserved_keywords[match]; found {
		lex.push(NewToken(kind, match))
	} else {
		lex.push(NewToken(IDENTIFIER, match))
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
	lex.line++
	lex.advance(len(match))
}

// reserved regex patterns
var reserved_patterns = []regex_pattern{
	{regexp.MustCompile(`\r\n|\r|\n`), newline_handler},
	{regexp.MustCompile(`[ \t]+`), skip_handler},
	{regexp.MustCompile(`\/\/.*`), comment_handler},
	{regexp.MustCompile(`"[^"]*"`), string_handler},
	{regexp.MustCompile(`[0-9]+(\.[0-9]+)?`), number_handler},
	{regexp.MustCompile(`([a-zA-Z_]|[\x{1F600}-\x{1F64F}\x{2700}-\x{27BF}\x{1F680}-\x{1F6FF}\x{1F300}-\x{1F5FF}\x{1F900}-\x{1F9FF}\x{2600}-\x{26FF}\x{2300}-\x{23FF}\x{1F100}-\x{1F1FF}\x{1F200}-\x{1F2FF}\x{3297}\x{3299}\x{1F191}-\x{1F19A}\x{1F170}-\x{1F19A}])([a-zA-Z0-9_]|[\x{1F600}-\x{1F64F}\x{2700}-\x{27BF}\x{1F680}-\x{1F6FF}\x{1F300}-\x{1F5FF}\x{1F900}-\x{1F9FF}\x{2600}-\x{26FF}\x{2300}-\x{23FF}\x{1F100}-\x{1F1FF}\x{1F200}-\x{1F2FF}\x{3297}\x{3299}\x{1F191}-\x{1F19A}\x{1F170}-\x{1F19A}])*`), symbol_handler},

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
	{regexp.MustCompile(`&`), default_handler(AMPERSAND, "&")},

	// Shorthand
	{regexp.MustCompile(`\+\+`), default_handler(PLUS_PLUS, "++")},
	{regexp.MustCompile(`--`), default_handler(MINUS_MINUS, "--")},
	{regexp.MustCompile(`\+=`), default_handler(PLUS_EQUALS, "+=")},
	{regexp.MustCompile(`-=`), default_handler(MINUS_EQUALS, "-=")},
	{regexp.MustCompile(`\*=`), default_handler(STAR_EQUALS, "*=")},
	{regexp.MustCompile(`/=`), default_handler(SLASH_EQUALS, "/=")},
	{regexp.MustCompile(`%=`), default_handler(PERCENT_EQUALS, "%=")},
	{regexp.MustCompile(`\&=`), default_handler(AND_EQUALS, "&=")},
	{regexp.MustCompile(`\|=`), default_handler(OR_EQUALS, "|=")},

	// Math Operators
	{regexp.MustCompile(`\+`), default_handler(PLUS, "+")},
	{regexp.MustCompile(`-`), default_handler(DASH, "-")},
	{regexp.MustCompile(`/`), default_handler(SLASH, "/")},
	{regexp.MustCompile(`\*`), default_handler(STAR, "*")},
	{regexp.MustCompile(`%`), default_handler(PERCENT, "%")},
}
