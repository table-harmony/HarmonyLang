package lexer

import "fmt"

type TokenKind int

const (
	EOF TokenKind = iota
	NIL
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
	AMPERSAND

	// Shorthand
	PLUS_PLUS
	MINUS_MINUS
	PLUS_EQUALS
	MINUS_EQUALS
	STAR_EQUALS
	SLASH_EQUALS
	PERCENT_EQUALS
	AND_EQUALS
	OR_EQUALS
	NULLISH_ASSIGNMENT

	// Maths
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
	FN
	STRUCT
	INTERFACE
	IF
	ELSE
	WHILE
	FOR
	IN
	RETURN
	STATIC
	SWITCH
	CASE
	DEFAULT
	CONTINUE
	BREAK
	MAP
	TRY
	CATCH
	THROW
	TYPEOF
	TYPE
	AS
)

var reserved_keywords map[string]TokenKind = map[string]TokenKind{
	"true":      TRUE,
	"false":     FALSE,
	"nil":       NIL,
	"let":       LET,
	"const":     CONST,
	"import":    IMPORT,
	"from":      FROM,
	"fn":        FN,
	"if":        IF,
	"else":      ELSE,
	"while":     WHILE,
	"for":       FOR,
	"in":        IN,
	"return":    RETURN,
	"struct":    STRUCT,
	"interface": INTERFACE,
	"static":    STATIC,
	"switch":    SWITCH,
	"case":      CASE,
	"default":   DEFAULT,
	"continue":  CONTINUE,
	"break":     BREAK,
	"map":       MAP,
	"try":       TRY,
	"catch":     CATCH,
	"throw":     THROW,
	"typeof":    TYPEOF,
	"type":      TYPE,
	"as":        AS,
}

type Token struct {
	Kind  TokenKind
	Value string
	Line  int
}

func NewToken(kind TokenKind, value string) Token {
	return Token{
		Kind:  kind,
		Value: value,
	}
}

func (token Token) String() string {
	if token.IsOfKind(IDENTIFIER, NUMBER, STRING) {
		return fmt.Sprintf("{ Kind: %s, Value: %s }", token.Kind.String(), token.Value)
	}

	return fmt.Sprintf("{ Kind: %s }", token.Kind.String())
}

func (token Token) IsOfKind(expectedTokens ...TokenKind) bool {
	for _, expected := range expectedTokens {
		if expected == token.Kind {
			return true
		}
	}

	return false
}

func (kind TokenKind) String() string {
	switch kind {
	case EOF:
		return "eof"
	case NIL:
		return "nil"
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
	case OR_EQUALS:
		return "or_equals"
	case AND_EQUALS:
		return "and_equals"
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
	case FN:
		return "function"
	case IF:
		return "if"
	case ELSE:
		return "else"
	case FOR:
		return "for"
	case WHILE:
		return "while"
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
	case CONTINUE:
		return "continue"
	case BREAK:
		return "break"
	case AMPERSAND:
		return "ampersand"
	case MAP:
		return "map"
	case TRY:
		return "try"
	case CATCH:
		return "catch"
	case THROW:
		return "throw"
	case TYPEOF:
		return "typeof"
	case TYPE:
		return "type"
	case AS:
		return "as"
	default:
		return fmt.Sprintf("unknown(%d)", kind)
	}
}
