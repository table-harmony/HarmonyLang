package interpreter

import (
	"reflect"

	"github.com/table-harmony/HarmonyLang/src/ast"
)

type statement_handler func(statement ast.Statement, env *Environment)
type expression_handler func(expression ast.Expression, env *Environment) RuntimeValue

var statement_lookup = make(map[reflect.Type]statement_handler)
var expression_lookup = make(map[reflect.Type]expression_handler)

func register_statement_handler[T ast.Statement](handler statement_handler) {
	statement_type := reflect.TypeOf((*T)(nil)).Elem()
	statement_lookup[statement_type] = handler
}

func register_expression_handler[T ast.Expression](handler expression_handler) {
	expression_type := reflect.TypeOf((*T)(nil)).Elem()
	expression_lookup[expression_type] = handler
}

func create_lookups() {
	register_statement_handler[ast.ExpressionStatement](evaluate_expression_statement)
	register_statement_handler[ast.VariableDeclarationStatement](evaluate_variable_declaration_statement)
	register_statement_handler[ast.IfStatement](evaluate_if_statement)
	register_statement_handler[ast.BlockStatement](evaluate_block_statement)
	register_statement_handler[ast.ForStatement](evaluate_for_statement)
	register_statement_handler[ast.SwitchStatement](evaluate_switch_statement)

	//
	register_expression_handler[ast.PrefixExpression](evaluate_prefix_expression)
	register_expression_handler[ast.BinaryExpression](evalute_binary_expression)
	register_expression_handler[ast.AssignmentExpression](evaluate_assignment_expression)
	register_expression_handler[ast.SymbolExpression](evaluate_symbol_expression)
	register_expression_handler[ast.SwitchExpression](evaluate_switch_expression)

	//
	register_expression_handler[ast.BooleanExpression](evaluate_primary_statement)
	register_expression_handler[ast.NumberExpression](evaluate_primary_statement)
	register_expression_handler[ast.StringExpression](evaluate_primary_statement)
	register_expression_handler[ast.TernaryExpression](evaluate_ternary_expression)
}
