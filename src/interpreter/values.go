package interpreter

import (
	"fmt"

	"github.com/table-harmony/HarmonyLang/src/ast"
	"github.com/table-harmony/HarmonyLang/src/helpers"
)

type RuntimeValueType int

const (
	// Primitive Types
	NumberType RuntimeValueType = iota
	StringType
	BooleanType
	NilType

	// Reference Types
	PointerType
	ReferenceType

	// Complex Types
	VariableType
	ArrayType
	SliceType
	MapType
	StructType
	InterfaceType
	FunctionType
	AnonymousFunctionType
	AnyType
)

func (_type RuntimeValueType) ToString() string {
	switch _type {
	case NumberType:
		return "number"
	case StringType:
		return "string"
	case BooleanType:
		return "bool"
	case NilType:
		return "null"
	case PointerType:
		return "pointer"
	case ReferenceType:
		return "reference"
	case VariableType:
		return "variable"
	case ArrayType:
		return "array"
	case SliceType:
		return "slice"
	case MapType:
		return "map"
	case StructType:
		return "struct"
	case InterfaceType:
		return "interface"
	case FunctionType:
		return "function"
	case AnonymousFunctionType:
		return "anonymous_function"
	case AnyType:
		return "any"
	default:
		return fmt.Sprintf("unknown(%d)", _type)
	}
}

type RuntimeValue interface {
	getType() RuntimeValueType
	getValue() RuntimeValue
}

type RuntimeReference struct {
	Value RuntimeValue
}

func (RuntimeReference) getType() RuntimeValueType { return ReferenceType }
func (r RuntimeReference) getValue() RuntimeValue  { return r.Value.getValue() }

type RuntimePointer struct {
	Target *RuntimeValue
}

func (RuntimePointer) getType() RuntimeValueType { return PointerType }
func (p RuntimePointer) getValue() RuntimeValue {
	if p.Target == nil {
		return RuntimeNil{}
	}
	return (*p.Target).getValue()
}

type RuntimeNumber struct {
	Value float64
}

func (RuntimeNumber) getType() RuntimeValueType { return NumberType }
func (n RuntimeNumber) getValue() RuntimeValue  { return n }

type RuntimeString struct {
	Value string
}

func (RuntimeString) getType() RuntimeValueType { return StringType }
func (s RuntimeString) getValue() RuntimeValue  { return s }

type RuntimeBoolean struct {
	Value bool
}

func (RuntimeBoolean) getType() RuntimeValueType { return BooleanType }
func (b RuntimeBoolean) getValue() RuntimeValue  { return b }

type RuntimeNil struct {
	Value RuntimeValue
}

func (RuntimeNil) getType() RuntimeValueType { return NilType }
func (n RuntimeNil) getValue() RuntimeValue  { return n }

func ExpectRuntimeValue[T RuntimeValue](value RuntimeValue) (T, error) {
	if value.getType() == ReferenceType || value.getType() == PointerType {
		value = value.getValue()
	}
	return helpers.ExpectType[T](value)
}

func GetDefaultValue(valueType RuntimeValueType) RuntimeValue {
	switch valueType {
	case NumberType:
		return RuntimeNumber{0}
	case StringType:
		return RuntimeString{""}
	case BooleanType:
		return RuntimeBoolean{false}
	case PointerType:
		return RuntimePointer{nil}
	default:
		return RuntimeNil{}
	}
}

func isEqual(v1, v2 RuntimeValue) bool {
	if v1 == nil || v2 == nil {
		return v1 == v2
	}

	value1 := v1.getValue()
	value2 := v2.getValue()

	if value1 == nil || value2 == nil {
		return value1 == value2
	}

	if value1.getType() != value2.getType() {
		return false
	}

	switch val1 := value1.(type) {
	case RuntimeNumber:
		val2 := value2.(RuntimeNumber)
		return val1.Value == val2.Value
	case RuntimeString:
		val2 := value2.(RuntimeString)
		return val1.Value == val2.Value
	case RuntimeBoolean:
		val2 := value2.(RuntimeBoolean)
		return val1.Value == val2.Value
	case RuntimeNil:
		return true
	default:
		return value1 == value2
	}
}

type RuntimeVariable struct {
	Identifier   string
	IsConstant   bool
	Value        RuntimeValue
	ExplicitType RuntimeValueType
}

func (v RuntimeVariable) getType() RuntimeValueType { return v.Value.getType() }
func (v RuntimeVariable) getValue() RuntimeValue    { return v.Value.getValue() }

type AssignableValue interface {
	RuntimeValue
	assign(value RuntimeValue) error
}

func (v *RuntimeVariable) assign(value RuntimeValue) error {
	if v.IsConstant {
		return fmt.Errorf("cannot reassign constant variable '%s'", v.Identifier)
	}

	if v.ExplicitType != value.getType() && v.ExplicitType != AnyType {
		return fmt.Errorf("type mismatch: variable '%s' explicit type %v but assigned a %v",
			v.Identifier, v.ExplicitType.ToString(), value.getType().ToString())
	}

	v.Value = value
	return nil
}

type RuntimeFunction struct {
	Name       string
	Parameters []ast.Parameter
	Body       []ast.Statement
	ReturnType RuntimeValueType
	//TODO: functions should probably have their closure env (sense the call env is not necciraly the same as the declaration env)
}

func (RuntimeFunction) getType() RuntimeValueType { return FunctionType }
func (f RuntimeFunction) getValue() RuntimeValue  { return f }

type RuntimeAnonymousFunction struct {
	Parameters []ast.Parameter
	Body       []ast.Statement
	ReturnType RuntimeValueType
	//TODO: functions should probably have their closure env (sense the call env is not necciraly the same as the declaration env)
}

func (RuntimeAnonymousFunction) getType() RuntimeValueType { return AnonymousFunctionType }
func (f RuntimeAnonymousFunction) getValue() RuntimeValue  { return f }
