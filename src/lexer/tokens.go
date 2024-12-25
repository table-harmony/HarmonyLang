package lexer

import "fmt"

type TokenKind int

const (
	EOF TokenKind = iota
	NULL
	TRUE
	FALSE
	NUMBER
	STRING
	IDENTIFIER

	// Grouping & Braces
	OPEN_BRACKET
	CLOSE_BRACKET
	OPEN_CURLY
	CLOSE_CURLY
	OPEN_PAREN
	CLOSE_PAREN

	// Equivilance
	ASSIGNMENT
	EQUALS
	NOT_EQUALS
	NOT

	// Conditional
	LESS
	LESS_EQUALS
	GREATER
	GREATER_EQUALS

	// Logical
	OR
	AND

	// Symbols
	DOT
	DOT_DOT
	SEMI_COLON
	COLON
	QUESTION
	COMMA
	ARROW

	// Shorthand
	PLUS_PLUS
	MINUS_MINUS
	PLUS_EQUALS
	MINUS_EQUALS
	STAR_EQUALS
	SLASH_EQUALS
	PERCENT_EQUALS
	NULLISH_ASSIGNMENT

	//Maths
	PLUS
	DASH
	SLASH
	STAR
	PERCENT

	// Reserved Keywords
	LET
	CONST
	IMPORT
	FROM
	FUNC
	STRUCT
	INTERFACE
	IF
	ELSE
	FOREACH
	WHILE
	FOR
	EXPORT
	TYPEOF
	IN
	RETURN
	STATIC
	SWITCH
	CASE
	DEFAULT
)

// Reserved lookups for keywords
var reservedKeywords map[string]TokenKind = map[string]TokenKind{
	"true":      TRUE,
	"false":     FALSE,
	"null":      NULL,
	"let":       LET,
	"const":     CONST,
	"import":    IMPORT,
	"from":      FROM,
	"func":      FUNC,
	"if":        IF,
	"else":      ELSE,
	"foreach":   FOREACH,
	"while":     WHILE,
	"for":       FOR,
	"export":    EXPORT,
	"typeof":    TYPEOF,
	"in":        IN,
	"return":    RETURN,
	"struct":    STRUCT,
	"interface": INTERFACE,
	"static":    STATIC,
	"switch":    SWITCH,
	"case":      CASE,
	"default":   DEFAULT,
}

type Token struct {
	Kind  TokenKind
	Value string
}

func CreateToken(kind TokenKind, value string) Token {
	return Token{
		kind, value,
	}
}

func (token Token) ToString() string {
	if token.IsOfKind(IDENTIFIER, NUMBER, STRING) {
		return fmt.Sprintf("{ Kind: %s, Value: %s }", token.Kind.ToString(), token.Value)
	}

	return fmt.Sprintf("{ Kind: %s }", token.Kind.ToString())
}

func (token Token) IsOfKind(expectedTokens ...TokenKind) bool {
	for _, expected := range expectedTokens {
		if expected == token.Kind {
			return true
		}
	}

	return false
}

func (kind TokenKind) ToString() string {
	switch kind {
	case EOF:
		return "eof"
	case NULL:
		return "null"
	case NUMBER:
		return "number"
	case STRING:
		return "string"
	case TRUE:
		return "true"
	case FALSE:
		return "false"
	case IDENTIFIER:
		return "identifier"
	case OPEN_BRACKET:
		return "open_bracket"
	case CLOSE_BRACKET:
		return "close_bracket"
	case OPEN_CURLY:
		return "open_curly"
	case CLOSE_CURLY:
		return "close_curly"
	case OPEN_PAREN:
		return "open_paren"
	case CLOSE_PAREN:
		return "close_paren"
	case ASSIGNMENT:
		return "assignment"
	case EQUALS:
		return "equals"
	case NOT_EQUALS:
		return "not_equals"
	case NOT:
		return "not"
	case LESS:
		return "less"
	case LESS_EQUALS:
		return "less_equals"
	case GREATER:
		return "greater"
	case GREATER_EQUALS:
		return "greater_equals"
	case OR:
		return "or"
	case AND:
		return "and"
	case DOT:
		return "dot"
	case DOT_DOT:
		return "dot_dot"
	case SEMI_COLON:
		return "semi_colon"
	case COLON:
		return "colon"
	case QUESTION:
		return "question"
	case COMMA:
		return "comma"
	case PLUS_PLUS:
		return "plus_plus"
	case MINUS_MINUS:
		return "minus_minus"
	case PLUS_EQUALS:
		return "plus_equals"
	case MINUS_EQUALS:
		return "minus_equals"
	case STAR_EQUALS:
		return "star_equals"
	case SLASH_EQUALS:
		return "slash_equals"
	case NULLISH_ASSIGNMENT:
		return "nullish_assignment"
	case PERCENT_EQUALS:
		return "percent_equals"
	case PLUS:
		return "plus"
	case DASH:
		return "dash"
	case SLASH:
		return "slash"
	case STAR:
		return "star"
	case PERCENT:
		return "percent"
	case LET:
		return "let"
	case CONST:
		return "const"
	case IMPORT:
		return "import"
	case FROM:
		return "from"
	case FUNC:
		return "function"
	case IF:
		return "if"
	case ELSE:
		return "else"
	case FOREACH:
		return "foreach"
	case FOR:
		return "for"
	case WHILE:
		return "while"
	case EXPORT:
		return "export"
	case TYPEOF:
		return "typeof"
	case IN:
		return "in"
	case RETURN:
		return "return"
	case STRUCT:
		return "struct"
	case INTERFACE:
		return "interface"
	case STATIC:
		return "static"
	case SWITCH:
		return "switch"
	case CASE:
		return "case"
	case DEFAULT:
		return "default"
	case ARROW:
		return "arrow"
	default:
		return fmt.Sprintf("unknown(%d)", kind)
	}
}
