package parser

import (
	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/lexer"
)

func parse_statement(parser *parser) ast.Statement {
	token := parser.currentToken()
	handler, exists := statement_lookup[token.Kind]

	if exists {
		return handler(parser)
	}

	return parse_expression_statement(parser)
}

func parse_expression_statement(parser *parser) ast.Statement {
	expression := parse_expression(parser, default_bp)

	parser.expect(lexer.SEMI_COLON)
	parser.advance(1)

	return ast.ExpressionStatement{
		Expression: expression,
	}
}

func parse_variable_declaration_statement(parser *parser) ast.Statement {
	token := parser.currentToken()
	isConstant := token.Kind == lexer.CONST
	parser.advance(1)

	identifier := parser.expect(lexer.IDENTIFIER).Value
	parser.advance(1)

	var explicitType ast.Type
	if parser.currentToken().Kind == lexer.COLON {
		parser.advance(1)
		explicitType = parse_type(parser, default_bp)
	}

	var value ast.Expression
	if parser.currentToken().Kind == lexer.ASSIGNMENT {
		parser.advance(1)
		value = parse_expression(parser, assignment)
	} else if explicitType == nil {
		panic("Cannot define a variable without an explicit type or default value.")
	}

	parser.expect(lexer.SEMI_COLON)
	parser.advance(1)

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
	for !parser.isEmpty() && parser.currentToken().Kind != lexer.CLOSE_CURLY {
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
	if parser.currentToken().Kind == lexer.ELSE {
		parser.advance(1)

		if parser.currentToken().Kind == lexer.IF {
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

func parse_func_declaration_statement(parser *parser) ast.Statement {
	return ast.FunctionDeclarationStatment{}
}

func parse_interface_declaration_statement(parser *parser) ast.Statement {
	return ast.BlockStatement{}
}

func parse_struct_declaration_statement(parser *parser) ast.Statement {
	return ast.StructDeclarationStatement{}
}

func parse_switch_statement(parser *parser) ast.Statement {
	parser.expect(lexer.SWITCH)
	parser.advance(1)

	value := parse_expression(parser, assignment)

	parser.expect(lexer.OPEN_CURLY)
	parser.advance(1)

	cases := make([]ast.SwitchCaseStatement, 0)
	for !parser.isEmpty() && parser.currentToken().Kind != lexer.CLOSE_CURLY {
		var pattern ast.Expression

		if parser.currentToken().Kind == lexer.DEFAULT {
			parser.advance(1)
		} else {
			parser.expect(lexer.CASE)
			parser.advance(1)
			pattern = parse_expression(parser, assignment)
		}

		parser.expect(lexer.COLON)
		parser.advance(1)

		body := parse_block_statement(parser)

		cases = append(cases, ast.SwitchCaseStatement{
			Pattern: pattern,
			Body:    body,
		})
	}

	parser.expect(lexer.CLOSE_CURLY)
	parser.advance(1)

	return ast.SwitchStatement{
		Value: value,
		Cases: cases,
	}
}
