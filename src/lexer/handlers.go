package lexer

import (
	"regexp"
)

// The default regexHandler method
func defaultHandler(kind TokenKind, value string) regexHandler {
	return func(lex *lexer, _ *regexp.Regexp) {
		lex.advance(len(value))
		lex.push(CreateToken(kind, value))
	}
}

// The string regexHandler method
func stringHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	stringLiteral := lex.remainder()[match[0]+1 : match[1]-1]

	lex.push(CreateToken(STRING, stringLiteral))
	lex.advance(len(stringLiteral) + 2)
}

// The number regexHandler method
func numberHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.push(CreateToken(NUMBER, match))
	lex.advance(len(match))
}

// The symbol regexHandler method
func symbolHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())

	if kind, found := reservedKeywords[match]; found {
		lex.push(CreateToken(kind, match))
	} else {
		lex.push(CreateToken(IDENTIFIER, match))
	}

	lex.advance(len(match))
}

// The skip regexHandler method for blank spaces e.t.c
func skipHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	lex.advance(match[1])
}

// The commetn regexHandler ignores comments and advances in the lexer
func commentHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())

	if match != nil {
		lex.advance(match[1])
	}
}

// reserved regex patterns
var reservedPatterns = []regexPattern{
	{regexp.MustCompile(`\s+`), skipHandler},
	{regexp.MustCompile(`\/\/.*`), commentHandler},
	{regexp.MustCompile(`"[^"]*"`), stringHandler},
	{regexp.MustCompile(`[0-9]+(\.[0-9]+)?`), numberHandler},
	{regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`), symbolHandler},

	// Grouping & Braces
	{regexp.MustCompile(`\[`), defaultHandler(OPEN_BRACKET, "[")},
	{regexp.MustCompile(`\]`), defaultHandler(CLOSE_BRACKET, "]")},
	{regexp.MustCompile(`\{`), defaultHandler(OPEN_CURLY, "{")},
	{regexp.MustCompile(`\}`), defaultHandler(CLOSE_CURLY, "}")},
	{regexp.MustCompile(`\(`), defaultHandler(OPEN_PAREN, "(")},
	{regexp.MustCompile(`\)`), defaultHandler(CLOSE_PAREN, ")")},

	// Equivalence
	{regexp.MustCompile(`==`), defaultHandler(EQUALS, "==")},
	{regexp.MustCompile(`!=`), defaultHandler(NOT_EQUALS, "!=")},
	{regexp.MustCompile(`=`), defaultHandler(ASSIGNMENT, "=")},
	{regexp.MustCompile(`!`), defaultHandler(NOT, "!")},

	// Conditional
	{regexp.MustCompile(`<=`), defaultHandler(LESS_EQUALS, "<=")},
	{regexp.MustCompile(`<`), defaultHandler(LESS, "<")},
	{regexp.MustCompile(`>=`), defaultHandler(GREATER_EQUALS, ">=")},
	{regexp.MustCompile(`>`), defaultHandler(GREATER, ">")},

	// Logical
	{regexp.MustCompile(`\|\|`), defaultHandler(OR, "||")},
	{regexp.MustCompile(`&&`), defaultHandler(AND, "&&")},

	// Symbols
	{regexp.MustCompile(`\.\.`), defaultHandler(DOT_DOT, "..")},
	{regexp.MustCompile(`\.`), defaultHandler(DOT, ".")},
	{regexp.MustCompile(`;`), defaultHandler(SEMI_COLON, ";")},
	{regexp.MustCompile(`:`), defaultHandler(COLON, ":")},
	{regexp.MustCompile(`\?\?=`), defaultHandler(NULLISH_ASSIGNMENT, "??=")},
	{regexp.MustCompile(`\?`), defaultHandler(QUESTION, "?")},
	{regexp.MustCompile(`,`), defaultHandler(COMMA, ",")},

	// Shorthand
	{regexp.MustCompile(`\+\+`), defaultHandler(PLUS_PLUS, "++")},
	{regexp.MustCompile(`--`), defaultHandler(MINUS_MINUS, "--")},
	{regexp.MustCompile(`\+=`), defaultHandler(PLUS_EQUALS, "+=")},
	{regexp.MustCompile(`-=`), defaultHandler(MINUS_EQUALS, "-=")},
	{regexp.MustCompile(`\*=`), defaultHandler(STAR_EQUALS, "*=")},
	{regexp.MustCompile(`/=`), defaultHandler(SLASH_EQUALS, "/=")},

	// Math Operators
	{regexp.MustCompile(`\+`), defaultHandler(PLUS, "+")},
	{regexp.MustCompile(`-`), defaultHandler(DASH, "-")},
	{regexp.MustCompile(`/`), defaultHandler(SLASH, "/")},
	{regexp.MustCompile(`\*`), defaultHandler(STAR, "*")},
	{regexp.MustCompile(`%`), defaultHandler(PERCENT, "%")},
}
