package main

import "fmt"

type Type = string
const UnknownType = ""

func typeOf(node Node) Type {
	if node == nil {
		return ""
	}
	switch t := node.(type) {
	// case *FileNode:
	// case PackageNode:
	// case CommentNode:
	// case Stmt:
	// case ForStmt:
	// case IfStmt:
	// case *CurlyScope:
	// case *ArgNode:
	// case *ReturnNode:
	// case CallExpr:
	// case GetExpr:
	// case SetExpr:
	// case AssignExpr:
	// case CompLitExpr:
	// case VarExpr:

	case *BinaryExpr:
		if t.ty != "" {
			return t.ty
		}
		panic("unknown binaryexpr type")
	case *CompLitExpr:
		if t.ty != "" {
			return t.ty
		}
		panic("unknown complit type")

	case *VarStmt:
		if t.calcTy != "" {
			return t.calcTy
		}
		return typeOf(t.initExpr)

	case *StructNode:
		return Type(t.ident.str)
	case *FuncNode:
		args := csvJoinArgNode(t.arguments)
		rets := csvJoinArgNode(t.returns)
		return fmt.Sprintf("func(%s) (%s)", args, rets)

	case Arg:
		return Type(t.kind.str)
	case LitExpr:
		switch t.kind {
		case INT:
			return "int"
		case FLOAT:
			return "float"
		case STRING:
			return "string"
		default:
			panic("AAA")
		}

	default:
		panic(fmt.Sprintf("typeOf: Unknown NodeType: %T", t))
	}
}

func csvJoinArgNode(node *ArgNode) string {
	if node == nil { return "" }
	return csvJoinArgs(node.args)
}
func csvJoinArgs(args []Arg) string {
	str := ""
	for i := range args {
		str += typeOf(args[i])
		if i < len(args)-1 {
			str += ","
		}
	}
	return str
}
