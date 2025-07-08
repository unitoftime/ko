package main

import (
	"fmt"
	"strconv"
)

var regTypeMap = make(map[string]Type)

func checkType(ty Type) (Type, bool) {
	ret, ok := regTypeMap[ty.Name()]
	return ret, ok
}

func getType(ty Type) Type {
	ret, ok := regTypeMap[ty.Name()]
	if ok {
		return ret
	}
	regTypeMap[ty.Name()] = ty
	return ty
}


type Type interface {
	Underlying() Type
	Name() string // Returns a unique name for this type
}

type StructType struct {
	name string
	fields []Type
}

func (t *StructType) Underlying() Type {
	return t
}
func (t *StructType) Name() string {
	return t.name
}

type FuncType struct {
	name string
	generics []Type
	args []Type
	returns Type // TODO: only supports one return
}

func (t *FuncType) Underlying() Type {
	return t
}
func (t *FuncType) Name() string {
	gen := ""
	for _, o := range t.generics {
		gen = gen + "_" + o.Name()
	}
	return fmt.Sprintf("func_%s%s", t.name, gen)
	// args := ""
	// for _, o := range t.args {
	// 	args = args + "_" + o.Name()
	// }
	// ret := t.returns.Name()
	// return fmt.Sprintf("func_%s_%s_%s_%s", t.name, gen, args, ret)
}

type PointerType struct {
	base Type
}
func (t *PointerType) Underlying() Type {
	return t.base
}
func (t *PointerType) Name() string {
	return "*"+t.base.Name()
}

type ArrayType struct {
	len int
	base Type
}
func (t *ArrayType) Underlying() Type {
	return t.base
}
func (t *ArrayType) Name() string {
	return fmt.Sprintf("[%d]%s", t.len, t.base.Name())
}

type SliceType struct {
	base Type
}
func (t *SliceType) Underlying() Type {
	return t.base
}
func (t *SliceType) Name() string {
	return fmt.Sprintf("[]%s", t.base.Name())
}


type BasicType struct {
	name string // The name of the type

	// General Type information
	comparable bool
}

func (t *BasicType) Underlying() Type {
	return t
}

func (t *BasicType) Name() string {
	return t.name
}

func (t *BasicType) Default() string {
	return "0" // TODO: need to determine this based on name
}

type GenericType struct {
	name string
	// base Type // This is the type that the generic type is resolving to ????
}

func (t *GenericType) Underlying() Type {
	return t
}

func (t *GenericType) Name() string {
	return t.name
}

// type Type interface {
// 	Underlying()
// 	String()
// }

// Tries to cast type A to type B
func tryCast(a, b Type) bool {
	if a == b {
		return true // Skip: They are the same type
	}

	{
		aa, ok := a.(*GenericType)
		if ok {
			a = genericTypeMap[aa.name]
		}

		bb, ok := b.(*GenericType)
		if ok {
			b = genericTypeMap[bb.name]
		}
	}

	switch aa := a.(type) {
	case *PointerType:
		bb, ok := b.(*PointerType)
		if !ok { return false }

		// Try to cast with the pointer element type
		return tryCast(aa.base, bb.base)

	case *BasicType:
		bb, ok := b.(*BasicType)
		if !ok { return false }
		return tryBasicTypeCast(aa, bb)
	}

	return false
}

// try to cast a to type b
func tryBasicTypeCast(a, b *BasicType) bool {
	switch a {
	case IntLitType:
		_, ok := intLitCast[b.name]
		// Println("IntLitCast:", a, b, ok)
		Println("basictypecast:", a, b, ok)
		return ok
	case FloatLitType:
		_, ok := floatLitCast[b.name]
		// Println("FloatLitCast:", a, b, ok)
		return ok
	}

	return false
}

//TODO: ResolveToConstInt
func castToInt(n Node) (int, bool) {
	switch t := n.(type) {
	case *LitExpr:
		lenVal, err := strconv.Atoi(t.tok.str)
		if err != nil {
			nodeError(n, "array length expression must convert to an int")
		}
		return lenVal, true
	default:
		panic(fmt.Sprintf("ResolveToConstInt: Unknown NodeType: %T", t))
	}
	return 0, false
}

// List of all the types that an untypedInt can cast to
var intLitCast = map[string]bool {
	IntLitName: true,

	"byte": true,
	"rune": true,

	"int": true,
	"uint": true,
	"uintptr": true,

	"u8": true,
	"uint8": true,
	"u16": true,
	"uint16": true,
	"u32": true,
	"uint32": true,
	"u64": true,
	"uint64": true,

	"i8": true,
	"int8": true,
	"i16": true,
	"int16": true,
	"i32": true,
	"int32": true,
	"i64": true,
	"int64": true,
}

var floatLitCast = map[string]bool {
	"byte": true,
	"rune": true,

	"int": true,
	"uint": true,
	"u8": true,
	"uint8": true,
	"u16": true,
	"uint16": true,
	"u32": true,
	"uint32": true,
	"u64": true,
	"uint64": true,

	"i8": true,
	"int8": true,
	"i16": true,
	"int16": true,
	"i32": true,
	"int32": true,
	"i64": true,
	"int64": true,

	//
	"f32": true,
	"float32": true,
	"f64": true,
	"float64": true,
}

const IntLitName = "untypedInt"
const FloatLitName = "untypedFloat"
const StringLitName = "untypedString"
const BoolLitName = "untypedBool"
const UntypedPointerName = "untypedPtr"

var (
	UnknownType Type = nil
	VoidType = &BasicType{"void", false} // TODO: Comparability?

	// These are literal types that can be dynamically resolved to whatever is needed
	// TODO: UntypedBool,Int,Rune,Float,Complex,String,Nil?
	BoolLitType = &BasicType{BoolLitName, true}
	IntLitType = &BasicType{IntLitName, true}
	FloatLitType = &BasicType{FloatLitName, true}
	StringLitType = &BasicType{StringLitName, true}
	PointerLitType = &BasicType{UntypedPointerName, true}

	BoolType = &BasicType{"bool", true}
	IntType = &BasicType{"int", true}
	Float64Type = &BasicType{"f64", true}
	StringType = &BasicType{"string", true}
)

var AppendGenericT = &GenericType{"T"}
var AppendBuiltinType = &FuncType{
	name: "append",

	// args: []Type{
	// 	getType(&SliceType{getType(IntType)}),
	// 	getType(IntType), // TODO: Variadic
	// },
	// returns: VoidType,
	// // returns: getType(&SliceType{getType(IntType)}),

	generics: []Type{AppendGenericT},
	args: []Type{
		&SliceType{AppendGenericT},
		AppendGenericT, // TODO: Variadic
	},
	returns: VoidType,
}
var LenBuiltinType = &FuncType{
	name: "len",
	generics: []Type{AppendGenericT},
	args: []Type{
		AppendGenericT,
	},
	returns: getType(IntType),
}


// const AutoType = "auto"

// Map from ko types to C equivalent types
var koToCMap = map[string]string{
	"nil": "NULL",

	StringLitName: "__ko_string",
	"string": "__ko_string",

	// Default unresolved lits
	IntLitName: "int",
	FloatLitName: "float",

	"byte": "uint8_t",
	"rune": "int32_t",

	"int": "int", // TODO: Correct?
	"uint": "uint", // TODO: Correct?

	"uintptr": "size_t",

	"f32": "float",
	"float32": "float",
	"f64": "double",
	"float64": "double",

	// "complex64": TODO
	// "complex128": TODO

	"u8": "uint8_t",
	"uint8": "uint8_t",
	"u16": "uint16_t",
	"uint16": "uint16_t",
	"u32": "uint32_t",
	"uint32": "uint32_t",
	"u64": "uint64_t",
	"uint64": "uint64_t",

	"i8": "int8_t",
	"int8": "int8_t",
	"i16": "int16_t",
	"int16": "int16_t",
	"i32": "int32_t",
	"int32": "int32_t",
	"i64": "int64_t",
	"int64": "int64_t",

	// "string": TODO
}
func typeStr(in Type) string {
	if in == UnknownType {
		panic("UNKNOWN TYPE")
	}
	if in.Name() == "" {
		panic("BLANK TYPE")
	}

	ret, ok := koToCMap[in.Name()]
	if !ok {
		return string(in.Name()) // If it wasn't a builtin type, then it probably came from a custom type
	}

	// TODO: Might be better to register the type or smth? then look it up later in the LUT
	// if !ok { panic(fmt.Sprintf("Unknown Type: %s", in)) }
	return ret
}

// Resolve literal types into their final type, if they dont resolve manually by a typespec
func resolveLitType(in Type) Type {
	switch t := in.(type) {
	case *BasicType:
		switch t.name {
		case BoolLitName:
			return BoolType
		case IntLitName:
			return IntType
		case FloatLitName:
			return Float64Type
		case StringLitName:
			return StringType
			// TODO: these ones
		// case UntypedPointerName:
		}
	}
	return in
}

// func typeOf(node Node) *Type {
// 	if node == nil {
// 		return UnknownType
// 	}
// 	switch t := node.(type) {
// 	// case *FileNode:
// 	// case PackageNode:
// 	// case CommentNode:
// 	// case Stmt:
// 	// case ForStmt:
// 	// case IfStmt:
// 	// case *CurlyScope:
// 	// case *ArgNode:
// 	// case *ReturnNode:
// 	// case CallExpr:
// 	// case GetExpr:
// 	// case SetExpr:
// 	// case AssignExpr:
// 	// case CompLitExpr:
// 	// case VarExpr:

// 	case *BinaryExpr:
// 		if t.ty != UnknownType {
// 			return t.ty
// 		}
// 		panic("unknown binaryexpr type")
// 	case *CompLitExpr:
// 		if t.ty != UnknownType {
// 			return t.ty
// 		}
// 		panic("unknown complit type")

// 	case *VarStmt:
// 		if t.ty != UnknownType {
// 			return t.ty
// 		}
// 		return t.ty

// 	case *StructNode:
// 		return &Type{t.ident.str, true, true}

// 	case *FuncNode:
// 		panic("YOU CANT CONSTRUCT THESE HERE THY NEED TOBE LOOKED UP FROM SCOPE EVERY TIME!")
// 		// Note: This needs to return a more complicated type that can be given args and can return args
// 		if t.returns == nil || len(t.returns.args) == 0 {
// 			return UnknownType
// 		}
// 		return &Type{t.returns.args[0].kind.str, true}

// 		// Note: Somewhat confusingly, we actually want the type that the funcNode *returns* not the type of the funcnode
// 		// rets := csvJoinArgNode(t.returns)
// 		// return rets

// 		// args := csvJoinArgNode(t.arguments)
// 		// rets := csvJoinArgNode(t.returns)
// 		// return fmt.Sprintf("func(%s) (%s)", args, rets)

// 	case *Arg:
// 		return &Type{t.kind.str, true}
// 	case *LitExpr:
// 		switch t.kind {
// 		case INT:
// 			return IntLitType
// 		case FLOAT:
// 			return FloatLitType
// 		case STRING:
// 			return StringLitType
// 		default:
// 			panic("AAA")
// 		}
// 	case *BuiltinNode:
// 		return t.ty

// 	default:
// 		panic(fmt.Sprintf("typeOf: Unknown NodeType: %T", t))
// 	}
// }

// func csvJoinArgNode(node *ArgNode) Type {
// 	if node == nil { return "" }
// 	return csvJoinArgs(node.args)
// }
// func csvJoinArgs(args []*Arg) Type {
// 	str := Type("")
// 	for i := range args {
// 		str += typeOf(args[i])
// 		if i < len(args)-1 {
// 			str += ","
// 		}
// 	}
// 	return str
// }

// func (r *Resolver) isComparable(n Node) bool {
// 	// Printf("isComparable: %T\n", n)

// 	switch t := n.(type) {
// 	case *StructNode:
// 		// TODO: Check every field to make sure it is also comparable
// 		return true
// 	case *IdentExpr:
// 		nodeDef, ok := r.CheckScope(t.tok.str)
// 		if !ok {
// 			printErr(t.tok, fmt.Sprintf("Unable to find identifier: %s", t.tok.str))
// 			return false
// 		}
// 		return r.isComparable(nodeDef)
// 	case *CallExpr:
// 		// TODO: Check the type that whatever it points to returns
// 		return true
// 	case *GetExpr:
// 		// TODO: Check the type that whatever it points to has at that field
// 		return true
// 	case *VarStmt:
// 		return true
// 	case *LitExpr:
// 		return true // Always comparable
// 	}
// 	return false
// }

