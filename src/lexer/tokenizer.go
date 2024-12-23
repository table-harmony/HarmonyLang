package lexer

import "fmt"

func Tokenize(source string) []Token {
	lexer := createLexer(source)

	for !lexer.at_eof() {
		matched := false

		for _, pattern := range lexer.patterns {
			location := pattern.regex.FindStringIndex(lexer.remainder())

			if location != nil && location[0] == 0 {
				pattern.handler(lexer, pattern.regex)
				matched = true
				break
			}
		}

		if !matched {
			panic(fmt.Sprintf("lexer error: unrecognized token near '%v'", lexer.remainder()))
		}
	}

	lexer.push(CreateToken(EOF, "EOF"))
	return lexer.Tokens
}
