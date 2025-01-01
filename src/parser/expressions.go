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
		panic(fmt.Sprintf("NUD Handler expected for token %s\n", token.Kind.String()))
	}

	left := nud_handler(parser)
	token = parser.current_token()
	for binding_power_lookup[token.Kind] > bp {
		led_handler, exists := led_lookup[token.Kind]
		if !exists {
			panic(fmt.Sprintf("LED Handler expected for token %s\n", token.Kind.String()))
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
	case lexer.NIL:
		return ast.NilExpression{}
	case lexer.TRUE:
		return ast.BooleanExpression{Value: true}
	case lexer.FALSE:
		return ast.BooleanExpression{Value: false}
	case lexer.NUMBER:
		number, err := strconv.ParseFloat(token.Value, 64)
		if err != nil {
			panic(fmt.Sprintf("Cannot parse token '%s' to float", token.String()))
		}

		return ast.NumberExpression{Value: number}
	case lexer.STRING:
		return ast.StringExpression{Value: token.Value}
	case lexer.IDENTIFIER:
		return ast.SymbolExpression{Value: token.Value}
	default:
		panic(fmt.Sprintf("Cannot create primary_expression from %s\n", token.Kind.String()))
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

func parse_call_expression(parser *parser, left ast.Expression, bp binding_power) ast.Expression {
	params := make([]ast.Expression, 0)

	parser.expect(lexer.OPEN_PAREN)
	parser.advance(1)

	for !parser.is_empty() && parser.current_token().Kind != lexer.CLOSE_PAREN {
		param := parse_expression(parser, default_bp)
		params = append(params, param)

		if !parser.is_empty() && parser.current_token().Kind != lexer.CLOSE_PAREN {
			parser.expect(lexer.COMMA)
			parser.advance(1)
		}
	}

	parser.expect(lexer.CLOSE_PAREN)
	parser.advance(1)

	return ast.CallExpression{
		Caller: left,
		Params: params,
	}
}

func parse_member_expression(parser *parser, left ast.Expression, bp binding_power) ast.Expression {
	parser.expect(lexer.DOT)
	parser.advance(1)

	property, err := ast.ExpectExpression[ast.SymbolExpression](parse_primary_expression(parser))

	if err != nil {
		panic(err)
	}

	return ast.MemberExpression{
		Owner:    left,
		Property: property,
	}
}

func parse_computed_member_expression(parser *parser, left ast.Expression, bp binding_power) ast.Expression {
	parser.expect(lexer.OPEN_BRACKET)
	parser.advance(1)

	property := parse_expression(parser, default_bp)

	parser.expect(lexer.CLOSE_BRACKET)
	parser.advance(1)

	return ast.ComputedMemberExpression{
		Owner:    left,
		Property: property,
	}
}

func parse_block_expression(parser *parser) ast.Expression {
	parser.expect(lexer.OPEN_CURLY)
	parser.advance(1)

	statements := make([]ast.Statement, 0)
	for !parser.is_empty() && parser.current_token().Kind != lexer.CLOSE_CURLY {
		statement := parse_statement(parser)
		statements = append(statements, statement)
	}

	parser.expect(lexer.CLOSE_CURLY)
	parser.advance(1)

	return ast.BlockExpression{
		Statements: statements,
	}
}

func parse_if_expression(parser *parser) ast.Expression {
	parser.expect(lexer.IF)
	parser.advance(1)

	condition := parse_expression(parser, assignment)
	consequent := parse_block_expression(parser).(ast.BlockExpression)

	var alternate ast.Expression
	if parser.current_token().Kind == lexer.ELSE {
		parser.advance(1)

		if parser.current_token().Kind == lexer.IF {
			alternate = parse_if_expression(parser).(ast.IfExpression)
		} else {
			alternate = parse_block_expression(parser).(ast.BlockExpression)
		}
	}

	return ast.IfExpression{
		Condition:  condition,
		Consequent: consequent,
		Alternate:  alternate,
	}
}

func parse_switch_expression(parser *parser) ast.Expression {
	parser.expect(lexer.SWITCH)
	parser.advance(1)

	value := parse_expression(parser, assignment)

	parser.expect(lexer.OPEN_CURLY)
	parser.advance(1)

	cases := make([]ast.SwitchCaseStatement, 0)
	for !parser.is_empty() && parser.current_token().Kind != lexer.CLOSE_CURLY {
		if parser.current_token().Kind == lexer.DEFAULT {
			parser.advance(1)

			body := parse_block_expression(parser).(ast.BlockExpression)

			cases = append(cases, ast.SwitchCaseStatement{
				Body:      body,
				IsDefault: true,
			})
			for !parser.is_empty() && parser.current_token().Kind == lexer.SEMI_COLON {
				parser.advance(1)
			}
		} else {
			parser.expect(lexer.CASE)
			parser.advance(1)

			var patterns []ast.Expression
			for !parser.is_empty() && parser.current_token().Kind != lexer.OPEN_CURLY {
				pattern := parse_expression(parser, assignment)

				patterns = append(patterns, pattern)

				if parser.current_token().Kind != lexer.OPEN_CURLY {
					parser.expect(lexer.COMMA)
					parser.advance(1)
				}
			}

			body := parse_block_expression(parser).(ast.BlockExpression)

			cases = append(cases, ast.SwitchCaseStatement{
				Patterns: patterns,
				Body:     body,
			})

			for !parser.is_empty() && parser.current_token().Kind == lexer.SEMI_COLON {
				parser.advance(1)
			}
		}
	}

	parser.expect(lexer.CLOSE_CURLY)
	parser.advance(1)

	return ast.SwitchExpression{
		Value: value,
		Cases: cases,
	}
}

func parse_function_declaration_expression(parser *parser) ast.Expression {
	parser.expect(lexer.FN)
	parser.advance(1)

	parser.expect(lexer.OPEN_PAREN)
	parser.advance(1)

	params := make([]ast.Parameter, 0)
	for !parser.is_empty() && parser.current_token().Kind != lexer.CLOSE_PAREN {
		param_name := parser.expect(lexer.IDENTIFIER).Value
		parser.advance(1)

		var param_type ast.Type
		if parser.current_token().Kind == lexer.COLON {
			parser.expect(lexer.COLON)
			parser.advance(1)

			param_type = parse_type(parser, default_bp)
		}

		var param_default_value ast.Expression
		if parser.current_token().Kind == lexer.ASSIGNMENT {
			parser.expect(lexer.ASSIGNMENT)
			parser.advance(1)

			param_default_value = parse_expression(parser, default_bp)
		}

		params = append(params, ast.Parameter{
			Name:         param_name,
			Type:         param_type,
			DefaultValue: param_default_value,
		})

		if !parser.current_token().IsOfKind(lexer.CLOSE_PAREN, lexer.EOF) {
			parser.expect(lexer.COMMA)
			parser.advance(1)
		}
	}

	parser.expect(lexer.CLOSE_PAREN)
	parser.advance(1)

	var returnType ast.Type
	if parser.current_token().Kind == lexer.ARROW {
		parser.advance(1)
		returnType = parse_type(parser, default_bp)
	}

	parser.expect(lexer.OPEN_CURLY)
	parser.advance(1)

	body := make([]ast.Statement, 0)
	for !parser.is_empty() && parser.current_token().Kind != lexer.CLOSE_CURLY {
		statement := parse_statement(parser)
		body = append(body, statement)
	}

	parser.expect(lexer.CLOSE_CURLY)
	parser.advance(1)

	return ast.FunctionDeclarationExpression{
		Parameters: params,
		Body:       body,
		ReturnType: returnType,
	}
}

func parse_try_catch_expression(parser *parser) ast.Expression {
	parser.expect(lexer.TRY)
	parser.advance(1)

	tryBlock := parse_block_expression(parser)

	parser.expect(lexer.CATCH)
	parser.advance(1)

	//TODO: parse error as well

	catchBlock := parse_block_expression(parser)

	return ast.TryCatchExpression{
		TryBlock:   tryBlock,
		CatchBlock: catchBlock,
	}
}

func parse_map_instantiation_expression(parser *parser) ast.Expression {
	parser.expect(lexer.MAP)
	parser.advance(1)

	parser.expect(lexer.OPEN_BRACKET)
	parser.advance(1)

	keyType := parse_type(parser, default_bp)

	parser.expect(lexer.ARROW)
	parser.advance(1)

	valueType := parse_type(parser, default_bp)

	parser.expect(lexer.CLOSE_BRACKET)
	parser.advance(1)

	parser.expect(lexer.OPEN_CURLY)
	parser.advance(1)

	entries := make([]ast.MapEntry, 0)
	for !parser.is_empty() && parser.current_token().Kind != lexer.CLOSE_CURLY {
		key := parse_expression(parser, comma)

		parser.expect(lexer.ARROW)
		parser.advance(1)

		value := parse_expression(parser, comma)

		entries = append(entries, ast.MapEntry{
			Key:   key,
			Value: value,
		})

		if parser.current_token().Kind != lexer.CLOSE_CURLY {
			parser.expect(lexer.COMMA)
			parser.advance(1)
		}
	}

	parser.expect(lexer.CLOSE_CURLY)
	parser.advance(1)

	return ast.MapInstantiationExpression{
		KeyType:   keyType,
		ValueType: valueType,
		Entries:   entries,
	}
}

func parse_array_instantiation_expression(parser *parser) ast.Expression {
	parser.expect(lexer.OPEN_BRACKET)
	parser.advance(1)

	var size ast.Expression
	if parser.current_token().Kind != lexer.CLOSE_BRACKET {
		size = parse_expression(parser, comma)
	}

	parser.expect(lexer.CLOSE_BRACKET)
	parser.advance(1)

	elementType := parse_type(parser, default_bp)

	parser.expect(lexer.OPEN_CURLY)
	parser.advance(1)

	elements := make([]ast.Expression, 0)
	for !parser.is_empty() && parser.current_token().Kind != lexer.CLOSE_CURLY {
		element := parse_expression(parser, comma)
		elements = append(elements, element)

		if parser.current_token().Kind == lexer.SEMI_COLON || parser.current_token().Kind == lexer.COMMA {
			parser.advance(1)
		}
	}

	parser.expect(lexer.CLOSE_CURLY)
	parser.advance(1)

	if size != nil {
		return ast.ArrayInstantiationExpression{
			Size:        size,
			ElementType: elementType,
			Elements:    elements,
		}
	}

	return ast.SliceInstantiationExpression{
		ElementType: elementType,
		Elements:    elements,
	}
}
