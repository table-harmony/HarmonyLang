package parser

import (
	"fmt"
	"strconv"

	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/lexer"
)

func parse_expression(parser *parser, bp binding_power) ast.Expression {
	var currentToken = parser.currentToken()
	nud_handler, exists := nud_lookup[currentToken.Kind]

	if !exists {
		panic(fmt.Sprintf("NUD Handler expected for token %s\n", currentToken.Kind.ToString()))
	}

	left := nud_handler(parser)

	for binding_power_lookup[parser.currentToken().Kind] > bp {
		currentToken = parser.currentToken()
		led_handler, exists := led_lookup[currentToken.Kind]

		if !exists {
			panic(fmt.Sprintf("LED Handler expected for token %s\n", currentToken.Kind.ToString()))
		}

		left = led_handler(parser, left, bp)
	}

	return left
}

func parse_binary_expression(parser *parser, left ast.Expression, bp binding_power) ast.Expression {
	operatorToken := parser.currentToken()
	parser.advance(1)

	right := parse_expression(parser, defalt_bp)

	return ast.BinaryExpression{
		Left:     left,
		Operator: operatorToken,
		Right:    right,
	}
}

func parse_primary_expression(parser *parser) ast.Expression {
	currentToken := parser.currentToken()
	parser.advance(1)

	switch currentToken.Kind {
	case lexer.NUMBER:
		number, err := strconv.ParseFloat(currentToken.Value, 64)

		if err != nil {
			panic(fmt.Sprintf("Cannot parse token '%s' to float", currentToken.ToString()))
		}

		return ast.NumberExpression{
			Value: number,
		}
	case lexer.STRING:
		return ast.StringExpression{
			Value: currentToken.Value,
		}
	case lexer.IDENTIFIER:
		return ast.SymbolExpression{
			Value: currentToken.Value,
		}
	default:
		panic(fmt.Sprintf("Cannot create primary_expression from %s\n", currentToken.Kind.ToString()))
	}
}
