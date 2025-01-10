package interpreter

import (
	"reflect"

	"github.com/table-harmony/HarmonyLang/src/ast"
)

type statement_handler func(statement ast.Statement, scope *Scope)
type expression_handler func(expression ast.Expression, scope *Scope) Value

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
	register_statement_handler[ast.TraditionalForStatement](evaluate_traditional_for_statement)
	register_statement_handler[ast.IteratorForStatement](evaluate_iterator_for_statement)
	register_statement_handler[ast.ContinueStatement](evaluate_continue_statement)
	register_statement_handler[ast.BreakStatement](evaluate_break_statement)
	register_statement_handler[ast.ReturnStatement](evaluate_return_statement)
	register_statement_handler[ast.FunctionDeclarationStatment](evaluate_function_declaration_statement)
	register_statement_handler[ast.AssignmentStatement](evaluate_assignment_statement)
	register_statement_handler[ast.ThrowStatement](evaluate_throw_statement)
	register_statement_handler[ast.TypeDeclarationStatement](evaluate_type_declaration_statement)
	register_statement_handler[ast.ImportStatement](evaluate_import_statement)
	register_statement_handler[ast.StructDeclarationStatement](evaluate_struct_declaration_statement)

	// Expressions
	register_expression_handler[ast.PrefixExpression](evaluate_prefix_expression)
	register_expression_handler[ast.BinaryExpression](evaluate_binary_expression)
	register_expression_handler[ast.SymbolExpression](evaluate_symbol_expression)
	register_expression_handler[ast.TernaryExpression](evaluate_ternary_expression)
	register_expression_handler[ast.CallExpression](evaluate_call_expression)
	register_expression_handler[ast.FunctionDeclarationExpression](evaluate_function_declaration_expression)
	register_expression_handler[ast.ComputedMemberExpression](evaluate_computed_member_expression)
	register_expression_handler[ast.MemberExpression](evaluate_member_expression)
	register_expression_handler[ast.RangeExpression](evaluate_range_expression)
	register_expression_handler[ast.StructLiteralExpression](evaluate_struct_instantiation_expression)

	// Block expressions
	register_expression_handler[ast.BlockExpression](evaluate_block_expression)
	register_expression_handler[ast.IfExpression](evaluate_if_expression)
	register_expression_handler[ast.SwitchExpression](evaluate_switch_expression)
	register_expression_handler[ast.TryCatchExpression](evaluate_try_catch_expression)

	// Data types expressions
	register_expression_handler[ast.ArrayInstantiationExpression](evaluate_array_instantiation_expression)
	register_expression_handler[ast.SliceInstantiationExpression](evaluate_slice_instantiation_expression)
	register_expression_handler[ast.MapInstantiationExpression](evaluate_map_instantiation_expression)

	// Primary expressions
	register_expression_handler[ast.BooleanExpression](evaluate_primary_expression)
	register_expression_handler[ast.NumberExpression](evaluate_primary_expression)
	register_expression_handler[ast.StringExpression](evaluate_primary_expression)
	register_expression_handler[ast.NilExpression](evaluate_primary_expression)
}
