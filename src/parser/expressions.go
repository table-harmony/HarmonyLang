package parser

import (
	"fmt"
	"strconv"

	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/lexer"
)

func parse_expression(parser *parser, bp binding_power) ast.Expression {
	token := parser.current_token()
	nud_handler, exists := nud_lookup[token.Kind]

	if !exists {
		panic(fmt.Sprintf("NUD Handler expected for token %s\n", token.Kind.ToString()))
	}

	left := nud_handler(parser)
	token = parser.current_token()

	for binding_power_lookup[token.Kind] > bp {
		led_handler, exists := led_lookup[token.Kind]
		if !exists {
			panic(fmt.Sprintf("LED Handler expected for token %s\n", token.Kind.ToString()))
		}

		left = led_handler(parser, left, binding_power_lookup[token.Kind])
		token = parser.current_token()
	}

	return left
}

func parse_binary_expression(parser *parser, left ast.Expression, bp binding_power) ast.Expression {
	operatorToken := parser.current_token()
	parser.advance(1)

	right := parse_expression(parser, bp)

	return ast.BinaryExpression{
		Left:     left,
		Operator: operatorToken,
		Right:    right,
	}
}

func parse_primary_expression(parser *parser) ast.Expression {
	token := parser.current_token()
	parser.advance(1)

	switch token.Kind {
	case lexer.NULL:
		return ast.NullExpression{}
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
	operator := parser.previous_token()

	valueExpression := ast.BinaryExpression{
		Left: left,
	}

	binaryOperators := map[lexer.TokenKind]lexer.TokenKind{
		lexer.PLUS_PLUS:      lexer.PLUS,
		lexer.MINUS_MINUS:    lexer.DASH,
		lexer.PLUS_EQUALS:    lexer.PLUS,
		lexer.MINUS_EQUALS:   lexer.DASH,
		lexer.STAR_EQUALS:    lexer.STAR,
		lexer.PERCENT_EQUALS: lexer.PERCENT,
		lexer.AND_EQUALS:     lexer.AND,
		lexer.OR_EQUALS:      lexer.OR,
	}

	getBinaryOperator := func() lexer.TokenKind {
		if op, exists := binaryOperators[operator.Kind]; exists {
			return op
		}
		return operator.Kind
	}

	switch operator.Kind {
	case lexer.PLUS_PLUS:
		valueExpression.Operator = lexer.CreateToken(lexer.PLUS, "++")
		valueExpression.Right = ast.NumberExpression{Value: 1}
	case lexer.MINUS_MINUS:
		valueExpression.Operator = lexer.CreateToken(lexer.DASH, "--")
		valueExpression.Right = ast.NumberExpression{Value: -1}
	case lexer.NULLISH_ASSIGNMENT, lexer.ASSIGNMENT:
		return ast.AssignmentExpression{
			Assigne:  left,
			Value:    parse_expression(parser, bp),
			Operator: lexer.CreateToken(getBinaryOperator(), ""),
		}
	default:
		valueExpression.Operator = lexer.CreateToken(getBinaryOperator(), "")
		valueExpression.Right = parse_expression(parser, bp)
	}

	return ast.AssignmentExpression{
		Assigne:  left,
		Value:    valueExpression,
		Operator: valueExpression.Operator,
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
	operatorToken := parser.current_token()
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

	switchCases := make([]ast.SwitchCase, 0)
	for !parser.is_empty() && parser.current_token().Kind != lexer.CLOSE_CURLY {
		if parser.current_token().Kind == lexer.IDENTIFIER && parser.current_token().Value == "_" {
			parser.advance(1)

			parser.expect_error(lexer.ARROW, fmt.Errorf("expected arrow after 'default' case in switch expression"))
			parser.advance(1)

			value := parse_expression(parser, assignment)

			switchCases = append(switchCases, ast.DefaultSwitchCase{
				Value: value,
			})

			if parser.current_token().Kind != lexer.CLOSE_CURLY {
				parser.expect(lexer.COMMA)
				parser.advance(1)
			}

			continue
		}

		var patterns []ast.Expression

		for !parser.is_empty() && parser.current_token().Kind != lexer.ARROW {
			pattern := parse_expression(parser, assignment)

			patterns = append(patterns, pattern)

			if parser.current_token().Kind != lexer.ARROW {
				parser.expect(lexer.COMMA)
				parser.advance(1)
			}
		}

		parser.expect(lexer.ARROW)
		parser.advance(1)

		value := parse_expression(parser, assignment)

		if parser.current_token().Kind != lexer.CLOSE_CURLY {
			parser.expect(lexer.COMMA)
			parser.advance(1)
		}

		switchCases = append(switchCases, ast.NormalSwitchCase{
			Patterns: patterns,
			Value:    value,
		})
	}

	parser.expect(lexer.CLOSE_CURLY)
	parser.advance(1)

	return ast.SwitchExpression{
		Value: value,
		Cases: switchCases,
	}
}

func parse_ternary_expression(parser *parser, left ast.Expression, bp binding_power) ast.Expression {
	parser.expect(lexer.QUESTION)
	parser.advance(1)

	consequent := parse_expression(parser, bp)

	parser.expect(lexer.COLON)
	parser.advance(1)

	alternate := parse_expression(parser, bp)

	return ast.TernaryExpression{
		Condition:  left,
		Consequent: consequent,
		Alternate:  alternate,
	}
}
