package parser

import (
	"fmt"

	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/lexer"
)

func parse_statement(parser *parser) ast.Statement {
	token := parser.current_token()
	statement_handler, exists := statement_lookup[token.Kind]

	var ast ast.Statement
	if exists {
		ast = statement_handler(parser)
	} else {
		ast = parse_expression_statement(parser)
	}

	if !parser.is_empty() {
		parser.expect(lexer.SEMI_COLON)
		parser.advance(1)
	}

	return ast
}

func parse_expression_statement(parser *parser) ast.Statement {
	expression := parse_expression(parser, default_bp)

	if handler, exists := sed_lookup[parser.current_token().Kind]; exists {
		return handler(parser, expression)
	}

	return ast.ExpressionStatement{
		Expression: expression,
	}
}

func parse_assignment_statement(parser *parser, left ast.Expression) ast.Statement {
	operator := parser.current_token()
	parser.advance(1)

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
		valueExpression.Right = ast.NumberExpression{Value: 1}
	case lexer.NULLISH_ASSIGNMENT, lexer.ASSIGNMENT:
		return ast.AssignmentStatement{
			Assigne:  left,
			Value:    parse_expression(parser, default_bp),
			Operator: lexer.CreateToken(getBinaryOperator(), ""),
		}
	default:
		valueExpression.Operator = lexer.CreateToken(getBinaryOperator(), "")
		valueExpression.Right = parse_expression(parser, default_bp)
	}

	return ast.AssignmentStatement{
		Assigne:  left,
		Value:    valueExpression,
		Operator: valueExpression.Operator,
	}
}

func parse_multi_variable_declaration_statement(parser *parser) ast.Statement {
	declarations := make([]ast.VariableDeclarationStatement, 0)

	isConstant := parser.current_token().Kind == lexer.CONST
	parser.advance(1)

	for !parser.is_empty() && parser.current_token().Kind != lexer.SEMI_COLON {
		identifier := parser.expect(lexer.IDENTIFIER).Value
		parser.advance(1)

		var explicitType ast.Type
		if parser.current_token().Kind == lexer.COLON {
			parser.advance(1)
			explicitType = parse_type(parser, default_bp)
		}

		var value ast.Expression
		if parser.current_token().Kind == lexer.ASSIGNMENT {
			parser.advance(1)
			value = parse_expression(parser, assignment)
		} else if explicitType == nil {
			panic("Cannot define a variable without an explicit type or default value.")
		}

		if isConstant && value == nil {
			panic("Cannot define constant variable without providing default value.")
		}

		declarations = append(declarations, ast.VariableDeclarationStatement{
			Identifier:   identifier,
			IsConstant:   isConstant,
			Value:        value,
			ExplicitType: explicitType,
		})

		if !parser.is_empty() && parser.current_token().Kind != lexer.SEMI_COLON {
			parser.expect(lexer.COMMA)
			parser.advance(1)
		}
	}

	return ast.MultiVariableDeclarationStatement{
		Declarations: declarations,
	}
}

func parse_import_statement(parser *parser) ast.Statement {
	parser.expect(lexer.IMPORT)
	parser.advance(1)

	name := parser.expect(lexer.IDENTIFIER).Value
	parser.advance(1)

	parser.expect(lexer.FROM)
	parser.advance(1)

	from := parser.expect(lexer.STRING).Value
	parser.advance(1)

	parser.expect(lexer.SEMI_COLON)
	parser.advance(1)

	return ast.ImportStatement{
		Name: name,
		From: from,
	}
}

func parse_interface_declaration_statement(parser *parser) ast.Statement {
	panic("Not implemented yet")
}

func parse_struct_declaration_statement(parser *parser) ast.Statement {
	panic("Not implemented yet")
}

func parse_for_statement(parser *parser) ast.Statement {
	parser.expect(lexer.FOR)
	parser.advance(1)

	initializer := parse_statement(parser)
	condition := parse_expression(parser, assignment)

	parser.expect(lexer.SEMI_COLON)
	parser.advance(1)

	var post []ast.Expression
	for !parser.is_empty() && parser.current_token().Kind != lexer.CLOSE_PAREN {
		post = append(post, parse_expression(parser, default_bp))

		if parser.current_token().Kind != lexer.CLOSE_PAREN {
			parser.expect(lexer.COMMA)
			parser.advance(1)
		}
	}

	body := make([]ast.Statement, 0)
	for !parser.is_empty() && parser.current_token().Kind != lexer.CLOSE_CURLY {
		statement := parse_statement(parser)
		body = append(body, statement)
	}

	return ast.ForStatement{
		Initializer: initializer,
		Condition:   condition,
		Post:        post,
		Body:        body,
	}
}

func parse_loop_control_statement(parser *parser) ast.Statement {
	token := parser.current_token()
	parser.advance(2)

	switch token.Kind {
	case lexer.CONTINUE:
		return ast.ContinueStatement{}
	case lexer.BREAK:
		return ast.BreakStatement{}
	default:
		panic(fmt.Sprintf("Cannot parse from token '%s' kind to loop_control_statement", token.ToString()))
	}
}

func parse_function_declaration_statement(parser *parser) ast.Statement {
	parser.expect(lexer.FN)
	parser.advance(1)

	identifier := parser.expect(lexer.IDENTIFIER)
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

	var return_type ast.Type
	if parser.current_token().Kind == lexer.ARROW {
		parser.advance(1)
		return_type = parse_type(parser, default_bp)
	}

	body := make([]ast.Statement, 0)
	for !parser.is_empty() && parser.current_token().Kind != lexer.CLOSE_CURLY {
		statement := parse_statement(parser)
		body = append(body, statement)
	}

	return ast.FunctionDeclarationStatment{
		Identifier: identifier.Value,
		Parameters: params,
		Body:       body,
		ReturnType: return_type,
	}
}

func parse_return_statement(parser *parser) ast.Statement {
	parser.expect(lexer.RETURN)
	parser.advance(1)

	value := parse_expression(parser, default_bp)

	return ast.ReturnStatement{
		Value: value,
	}
}
