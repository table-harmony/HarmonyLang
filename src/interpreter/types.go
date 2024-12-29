package interpreter

import "github.com/table-harmony/HarmonyLang/src/ast"

//TODO: lookups
func evaluate_type(_type ast.Type) RuntimeValueType {
	if _type == nil {
		return AnyType
	}

	if _, ok := _type.(ast.StringType); ok {
		return StringType
	}

	if _, ok := _type.(ast.NumberType); ok {
		return NumberType
	}

	if _, ok := _type.(ast.BooleanType); ok {
		return BooleanType
	}

	return AnyType
}
