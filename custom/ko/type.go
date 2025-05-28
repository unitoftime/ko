package main

// TODO: I think every type needs to just be a nil pointer in every node, and then it gets computed as I walk through the resolve. and has information, like what the underlying type is, and if its comparable, and little bits like that. Obviously things like structs and funcs have that all predeclared, then expressions i kinda just walk through to see what matches what
// type Type string
type Type struct {
	name string // The name of the type

	// General Type information
	comparable bool
	isStruct bool
}

var (
	BoolType = &Type{"bool", true, false}
)

// type Type interface {
// 	Underlying()
// 	String()
// }

// Tries to cast type A to type B
func tryCast(a, b *Type) bool {
	switch a {
	case IntLitType:
		_, ok := intLitCast[b.name]
		// fmt.Println("IntLitCast:", a, b, ok)
		return ok
	case FloatLitType:
		_, ok := floatLitCast[b.name]
		// fmt.Println("FloatLitCast:", a, b, ok)
		return ok
	}
	return false
}


// List of types that can be cast to int
var intLitCast = map[string]bool {
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

var (
	UnknownType *Type = nil
	VoidType *Type = &Type{"void", false, false} // TODO: Comparability?

	// These are literal types that can be dynamically resolved to whatever is needed
	// TODO: UntypedBool,Int,Rune,Float,Complex,String,Nil?
	BoolLitType *Type = &Type{BoolLitName, true, false}
	IntLitType *Type = &Type{IntLitName, true, false}
	FloatLitType *Type = &Type{FloatLitName, true, false}
	StringLitType *Type = &Type{StringLitName, true, false}

	// TODO: Should I do compound lits this way too?
)

// const AutoType = "auto"

// Map from ko types to C equivalent types
var typeMap = map[string]string{
	// Default unresolved lits
	IntLitName: "int",
	FloatLitName: "float",
	// StringLitName: "string",

	"byte": "uint8_t",
	"rune": "int32_t",

	"int": "int", // TODO: Correct?
	"uint": "uint", // TODO: Correct?

	// "uintptr": TODO

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
func typeStr(in *Type) string {
	if in == UnknownType {
		panic("UNKNOWN TYPE")
	}
	if in.name == "" {
		panic("BLANK TYPE")
	}

	ret, ok := typeMap[in.name]
	if !ok {
		return string(in.name) // If it wasn't a builtin type, then it probably came from a custom type
	}

	// TODO: Might be better to register the type or smth? then look it up later in the LUT
	// if !ok { panic(fmt.Sprintf("Unknown Type: %s", in)) }
	return ret
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
// 	// fmt.Printf("isComparable: %T\n", n)

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
