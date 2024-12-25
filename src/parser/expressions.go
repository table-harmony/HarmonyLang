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

	right := parse_expression(parser, bp+1)

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
	case lexer.TRUE:
		return ast.BooleanExpression{
			Value: true,
		}
	case lexer.FALSE:
		return ast.BooleanExpression{
			Value: false,
		}
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

func parse_assignment_expression(parser *parser, left ast.Expression, bp binding_power) ast.Expression {
	parser.advance(1)

	operator := parser.previousToken()

	var right ast.Expression
	switch operator.Kind {
	case lexer.PLUS_PLUS:
		right = ast.NumberExpression{Value: 1}
	case lexer.MINUS_MINUS:
		right = ast.NumberExpression{Value: -1}
	default:
		right = parse_expression(parser, bp)
	}

	return ast.AssignmentExpression{
		Assigne:  left,
		Value:    right,
		Operator: operator,
	}
}

func parse_grouping_expression(parser *parser) ast.Expression {
	parser.expect(lexer.OPEN_PAREN)
	parser.advance(1)

	expression := parse_expression(parser, default_bp)

	parser.expect(lexer.CLOSE_PAREN)
	parser.advance(1)

	return expression
}

func parse_prefix_expression(parser *parser) ast.Expression {
	operatorToken := parser.currentToken()
	parser.advance(1)

	right := parse_expression(parser, unary)

	return ast.PrefixExpression{
		Operator: operatorToken,
		Right:    right,
	}
}

func parse_switch_expression(parser *parser) ast.Expression {
	parser.expect(lexer.SWITCH)
	parser.advance(1)

	value := parse_expression(parser, assignment)

	parser.expect(lexer.OPEN_CURLY)
	parser.advance(1)

	cases := make([]ast.SwitchCase, 0)
	for !parser.isEmpty() && parser.currentToken().Kind != lexer.CLOSE_CURLY {
		pattern := parse_expression(parser, assignment)

		parser.expect(lexer.ARROW)
		parser.advance(1)

		value := parse_expression(parser, assignment)

		if parser.currentToken().Kind != lexer.CLOSE_CURLY {
			parser.expect(lexer.COMMA)
			parser.advance(1)
		}

		cases = append(cases, ast.SwitchCase{
			Pattern: pattern,
			Value:   value,
		})
	}

	parser.expect(lexer.CLOSE_CURLY)
	parser.advance(1)

	return ast.SwitchExpression{
		Value: value,
		Cases: cases,
	}
}
