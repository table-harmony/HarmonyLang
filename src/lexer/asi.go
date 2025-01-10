package lexer

func asi(lex *lexer) {
	for i := 0; i < len(lex.Tokens)-1; i++ {
		currentToken := lex.Tokens[i]
		nextToken := lex.Tokens[i+1]

		// Special handling for return statements followed by curly braces
		if currentToken.Kind == RETURN {
			// If there's a line break between return and {, treat it as expression
			if nextToken.Kind == OPEN_CURLY && nextToken.Line > currentToken.Line {
				continue // Skip semicolon insertion
			}
			// If return and { are on same line, proceed with normal semicolon rules
		}

		// Insert semicolon before closing curly brace
		if nextToken.Kind == CLOSE_CURLY {
			if needs_semi_colon(currentToken) {
				lex.insert_semi_colon(i)
				i++ // Skip the inserted semicolon
			}
		}

		// Don't insert semicolon if next token is a dot (method chaining)
		if nextToken.Kind == DOT {
			continue
		}

		// Don't insert semicolon if current token is a dot
		if currentToken.Kind == DOT {
			continue
		}

		// Insert semicolon if tokens are on different lines and current token needs semicolon
		if nextToken.Line > currentToken.Line && needs_semi_colon(currentToken) {
			// Check if next token is one that shouldn't have semicolon before it
			if !is_continuation_token(nextToken) {
				lex.insert_semi_colon(i)
				i++ // Skip the inserted semicolon
			}
		}
	}

	// Handle last token before EOF
	if len(lex.Tokens) > 1 {
		lastToken := lex.Tokens[len(lex.Tokens)-1] // -2 because -1 is EOF
		if needs_semi_colon(lastToken) {
			lex.insert_semi_colon(len(lex.Tokens) - 1)
		}
	}
}

func is_continuation_token(token Token) bool {
	return token.IsOfKind(
		DOT,
		PLUS,
		DASH,
		SLASH,
		PERCENT,
		OR,
		AND,
		EQUALS,
		NOT_EQUALS,
		LESS,
		LESS_EQUALS,
		GREATER,
		GREATER_EQUALS,
		ARROW,
	)
}
