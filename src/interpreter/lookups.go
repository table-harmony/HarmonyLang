package interpreter

import (
	"reflect"

	"github.com/table-harmony/HarmonyLang/src/ast"
)

type statement_handler func(statement ast.Statement, env *Scope)
type expression_handler func(expression ast.Expression, env *Scope) Value

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
	// Statements
	register_statement_handler[ast.ExpressionStatement](evaluate_expression_statement)
	register_statement_handler[ast.VariableDeclarationStatement](evaluate_variable_declaration_statement)
	register_statement_handler[ast.MultiVariableDeclarationStatement](evaluate_multi_variable_declaration_statement)
	register_statement_handler[ast.ForStatement](evaluate_for_statement)
	register_statement_handler[ast.ContinueStatement](evaluate_continue_statement)
	register_statement_handler[ast.BreakStatement](evaluate_break_statement)
	register_statement_handler[ast.ReturnStatement](evaluate_return_statement)
	register_statement_handler[ast.FunctionDeclarationStatment](evaluate_function_declaration_statement)
	register_statement_handler[ast.AssignmentStatement](evaluate_assignment_statement)

	// Expressions
	register_expression_handler[ast.PrefixExpression](evaluate_prefix_expression)
	register_expression_handler[ast.BinaryExpression](evaluate_binary_expression)
	register_expression_handler[ast.SymbolExpression](evaluate_symbol_expression)
	register_expression_handler[ast.TernaryExpression](evaluate_ternary_expression)
	register_expression_handler[ast.CallExpression](evaluate_call_expression)
	register_expression_handler[ast.FunctionDeclarationExpression](evaluate_function_declaration_expression)

	// Block expressions
	register_expression_handler[ast.BlockExpression](evaluate_block_expression)
	register_expression_handler[ast.IfExpression](evaluate_if_expression)
	register_expression_handler[ast.SwitchExpression](evaluate_switch_expression)
	register_expression_handler[ast.TryCatchExpression](evaluate_try_catch_expression)

	// Primary expressions
	register_expression_handler[ast.BooleanExpression](evaluate_primary_expression)
	register_expression_handler[ast.NumberExpression](evaluate_primary_expression)
	register_expression_handler[ast.StringExpression](evaluate_primary_expression)
	register_expression_handler[ast.NilExpression](evaluate_primary_expression)
}
