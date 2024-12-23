package parser

import (
	"fmt"
	"strconv"

	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/lexer"
)

func parse_expression(parser *parser, bp binding_power) ast.Expression {
	token := parser.currentToken()
	nud_handler, exists := nud_lookup[token.Kind]

	if !exists {
		panic(fmt.Sprintf("NUD Handler expected for token %s\n", token.Kind.ToString()))
	}

	left := nud_handler(parser)
	token = parser.currentToken()

	for binding_power_lookup[token.Kind] > bp {
		led_handler, exists := led_lookup[token.Kind]
		if !exists {
			panic(fmt.Sprintf("LED Handler expected for token %s\n", token.Kind.ToString()))
		}

		left = led_handler(parser, left, binding_power_lookup[token.Kind])
		token = parser.currentToken()
	}

	return left
}

func parse_binary_expression(parser *parser, left ast.Expression, bp binding_power) ast.Expression {
	operatorToken := parser.currentToken()
	parser.advance(1)

	right := parse_expression(parser, bp)

	return ast.BinaryExpression{
		Left:     left,
		Operator: operatorToken,
		Right:    right,
	}
}

func parse_primary_expression(parser *parser) ast.Expression {
	token := parser.currentToken()
	parser.advance(1)

	switch token.Kind {
	case lexer.NUMBER:
		number, err := strconv.ParseFloat(token.Value, 64)

		if err != nil {
			panic(fmt.Sprintf("Cannot parse token '%s' to float", token.ToString()))
		}

		return ast.NumberExpression{
			Value: number,
		}
	case lexer.STRING:
		return ast.StringExpression{
			Value: token.Value,
		}
	case lexer.IDENTIFIER:
		return ast.SymbolExpression{
			Value: token.Value,
		}
	default:
		panic(fmt.Sprintf("Cannot create primary_expression from %s\n", token.Kind.ToString()))
	}
}

//func parse_assignment_expression(parser *parser, left ast.Expression, bp binding_power) ast.Expression {
//	identifierToken := parser.currentToken()
//	parser.advance(1)
//
//	operatorToken := parser.currentToken()
//	parser.advance(1)
//
//	switch (operatorToken.Kind) {
//	case lexer.ASSIGNMENT:
//
//	}
//
//	return ast.AssignmentExpression{
//
//	}
//}

func parse_assignment_expr(parser *parser, left ast.Expression, bp binding_power) ast.Expression {
	parser.advance(1)
	rhs := parse_expression(parser, bp)

	return ast.AssignmentExpression{
		Assigne:       left,
		AssignedValue: rhs,
	}
}
