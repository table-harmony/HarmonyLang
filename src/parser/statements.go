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

	return ast.ExpressionStatement{
		Expression: expression,
	}
}

// TODO: implement multiple variable declaration (e.g. let a, b int = 1, 2 || let a = 1, b = 2)
func parse_variable_declaration_statement(parser *parser) ast.Statement {
	token := parser.current_token()

	isConstant := token.Kind == lexer.CONST
	parser.advance(1)

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

	return ast.VariableDeclarationStatement{
		Identifier:   identifier,
		IsConstant:   isConstant,
		Value:        value,
		ExplicitType: explicitType,
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

func parse_block_statement(parser *parser) ast.Statement {
	parser.expect(lexer.OPEN_CURLY)
	parser.advance(1)

	body := make([]ast.Statement, 0)
	for !parser.is_empty() && parser.current_token().Kind != lexer.CLOSE_CURLY {
		statement := parse_statement(parser)
		body = append(body, statement)
	}

	parser.expect(lexer.CLOSE_CURLY)
	parser.advance(1)

	return ast.BlockStatement{
		Body: body,
	}
}

func parse_if_statement(parser *parser) ast.Statement {
	parser.expect(lexer.IF)
	parser.advance(1)

	condition := parse_expression(parser, assignment)
	consequent := parse_block_statement(parser)

	var alternate ast.Statement
	if parser.current_token().Kind == lexer.ELSE {
		parser.advance(1)

		if parser.current_token().Kind == lexer.IF {
			alternate = parse_if_statement(parser)
		} else {
			alternate = parse_block_statement(parser)
		}
	}

	return ast.IfStatement{
		Condition:  condition,
		Consequent: consequent,
		Alternate:  alternate,
	}
}

func parse_interface_declaration_statement(parser *parser) ast.Statement {
	panic("Not implemented yet")
}

func parse_struct_declaration_statement(parser *parser) ast.Statement {
	panic("Not implemented yet")
}

func parse_switch_statement(parser *parser) ast.Statement {
	parser.expect(lexer.SWITCH)
	parser.advance(1)

	value := parse_expression(parser, assignment)

	parser.expect(lexer.OPEN_CURLY)
	parser.advance(1)

	cases := make([]ast.SwitchCaseStatement, 0)
	for !parser.is_empty() && parser.current_token().Kind != lexer.CLOSE_CURLY {
		var pattern ast.Expression

		if parser.current_token().Kind == lexer.DEFAULT {
			parser.advance(1)
		} else {
			parser.expect(lexer.CASE)
			parser.advance(1)

			pattern = parse_expression(parser, assignment)
		}

		parser.expect(lexer.COLON)
		parser.advance(1)

		body_statement := parse_block_statement(parser)
		block_statement, err := ast.ExpectStatement[ast.BlockStatement](body_statement)

		if err != nil {
			panic("Switch case body must be a block statement.")
		}

		cases = append(cases, ast.SwitchCaseStatement{
			Pattern: pattern,
			Body:    block_statement.Body,
		})
	}

	parser.expect(lexer.CLOSE_CURLY)
	parser.advance(1)

	return ast.SwitchStatement{
		Value: value,
		Cases: cases,
	}
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

	body_statement := parse_block_statement(parser)
	block_statement, err := ast.ExpectStatement[ast.BlockStatement](body_statement)

	if err != nil {
		panic("For statement body must be a block statement.")
	}

	return ast.ForStatement{
		Initializer: initializer,
		Condition:   condition,
		Post:        post,
		Body:        block_statement.Body,
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
	parser.expect(lexer.FUNC)
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

	body_statement := parse_block_statement(parser)
	block_statement, err := ast.ExpectStatement[ast.BlockStatement](body_statement)

	if err != nil {
		panic(err)
	}

	return ast.FunctionDeclarationStatment{
		Identifier: identifier.Value,
		Parameters: params,
		Body:       block_statement.Body,
		ReturnType: return_type,
	}
}

func parse_return_statement(parser *parser) ast.Statement {
	parser.expect(lexer.RETURN)
	parser.advance(1)

	value := parse_expression(parser, default_bp)

	parser.expect(lexer.SEMI_COLON)
	parser.advance(1)

	return ast.ReturnStatement{
		Value: value,
	}
}
