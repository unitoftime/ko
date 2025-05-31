package main

import (
	"bytes"
	"fmt"

	_ "embed"
)

type genBuf struct {
	buf *bytes.Buffer
	indent int
	newline bool
}
func (b *genBuf) Add(str string) *genBuf {
	if b.newline {
		b.newline = false
		for range b.indent {
			b.buf.WriteString("\t")
		}
	}

	b.buf.WriteString(str)
	return b
}

func (b *genBuf) Line() *genBuf {
	b.Add("\n")
	b.newline = true
	return b
}
func (b *genBuf) LineDirective(pos Position) *genBuf {
	// #line 31 "test.txt"
	// b.Add(fmt.Sprintf("#line %d \"%s\"", pos.line, pos.filename))
	b.Line()

	return b
}

func (b *genBuf) String() string {
	return b.buf.String()
}

// TODO: if a tuple, then return a struct
func returnArgsToString(node Node) string {
	if node == nil { return "void" }

	argNode, ok := node.(*ArgNode)
	if !ok { return "void" }
	if argNode == nil { return "void" }

	if !ok { panic(fmt.Sprintf("must be arg node: %T", node)) }
	if len(argNode.args) != 1 {
		panic("only supporting single return types")
	}
	Printf("returnArgsToString: %+v\n", argNode)
	Printf("returnArgsToString: %T\n", argNode.args[0])
	return typeNameC(argNode.args[0].ty)
}

//go:embed runtime.h
var runtimeFile string

//go:embed slice.tmpl
var sliceTemplate string
type SliceTemplateDef struct {
	Name string
	Type string
}

func (buf *genBuf) Generate(result ParseResult) {
	// -- Includes --
	// buf.
	// 	Add("#include <stdio.h>"). // TODO: hack. runtime.h?
	// 	Line()
	// 	Add("#include <stdint.h>").// TODO: hack. runtime.h?
	// 	Line()
	// // buf.Add("#include \"raylib.h\""). // TODO: hack. runtime.h?
	// // Line().Line()

	// -- Runtime --
	buf.Add(runtimeFile).Line()


	// -- Forward declarations --
	buf.Add("int __mainRet__ = 0;").Line()

	// Forward Declare all types
	for _, node := range result.typeList {
		buf.PrintForwardDecl(node)
		buf.Add(";").Line()

		structNode, isStruct := node.(*StructNode)
		if isStruct {
			buf.printEqualityPrototype(structNode.Type())
			buf.Add(";").Line()
		}
	}

	// Forward declare all arrays
	for _, ty := range regTypeMap {
		// TODO: You could probably register *special* types you find during typechecking, rather than looping everything
		buf.PrintGeneratedType(ty)
	}
	buf.Line()

	// buf.PrintStructForwardDecl("__ko_int_slice")
	// buf.Add(";").Line()
	// tmpl := template.Must(template.New("cslice").Parse(sliceTemplate))
	// err := tmpl.Execute(buf.buf, SliceTemplateDef{
	// 	Name: "__ko_int_slice",
	// 	Type: "int",
	// })
	// if err != nil {
	// 	panic(err)
	// }


	// Forward Declare all functions
	for i := range result.fnList {
		buf.PrintForwardDecl(result.fnList[i])
		buf.Add(";").Line()
	}

	// Complete all types
	for _, node := range result.typeList {
		buf.PrintCompleteType(node)
		buf.Add(";").Line()

		structNode, isStruct := node.(*StructNode)
		if isStruct {
			buf.printStructEqualityFunction(structNode)
		}
	}

	// Declare all global variables
	for i := range result.varList {
		buf.LineDirective(result.varList[i].name.pos)
		buf.PrintForwardDecl(result.varList[i])
		buf.Add(";").Line()
	}

	buf.Print(result.file)
}

func (buf *genBuf) PrintFuncDef(t *FuncNode) {
	isMain := t.name == "main"
	retArgs := "int" // Note: Default main return type
	if !isMain {
		retArgs = returnArgsToString(t.returns)
	}
	buf.Add(retArgs).
		Add(" ").
		Add(t.name).
		Add(" (")
	buf.Print(t.arguments)
	buf.Add(")")
}

func (buf *genBuf) PrintGeneratedType(ty Type) {
	switch t := ty.(type) {
	case *ArrayType:
		name := typeNameC(ty)

		// Forward Declaration
		// TODO: Would be good to hoist this all up so nested arrays arent problematic
		buf.Add("typedef struct ").
			Add(name).Add(" ").Add(name)
		buf.Add(";").Line()

		// Type definition
		elemName := typeNameC(t.base)
		buf.Add("struct ").Add(name).Add(" {").Line()
		buf.indent++
		buf.Add(elemName).Add(" ").Add("a")
		buf.Add(fmt.Sprintf("[%d]", t.len))
		buf.Add(";").Line()
		buf.indent--
		buf.Add("};")
	}
}

func (buf *genBuf) PrintStructForwardDecl(name string) {
	buf.Add("typedef struct ").Add(name).Add(" ").Add(name)
}

func (buf *genBuf) PrintForwardDecl(n Node) {
	switch t := n.(type) {
	case *StructNode:
		buf.PrintStructForwardDecl(t.ident.str)
	case *VarStmt:
		typeStr := typeNameC(t.ty)
		buf.
			Add(typeStr).
			Add(" ").
			Add(t.name.str).
			Add(" = ")
		buf.Print(t.initExpr)
	case *FuncNode:
		buf.LineDirective(t.pos)
		buf.PrintFuncDef(t)
	default:
		panic(fmt.Sprintf("PrintForwardDecl: Unknown NodeType: %T", t))
	}
}

func (buf *genBuf) PrintCompleteType(n Node) {
	switch t := n.(type) {
	case *StructNode:
		buf.PrintStructNode(t)
	default:
		panic(fmt.Sprintf("PrintForwardDecl: Unknown NodeType: %T", t))
	}
}

func (buf *genBuf) PrintArgList(args []Node) {
	for i := range args {
		buf.Print(args[i])

		if i < len(args)-1 {
			buf.Add(", ")
		}
	}
}

// Print the default value for the type
func (buf *genBuf) PrintDefault(ty Type) {
	switch t := ty.(type) {
	case *StructType:
		buf.PrintCompLit(&CompLitExpr{
			// callee: // TODO: Not needed I dont think. could maybe do a fake ident node with just field name
			args: nil, // So that it does the default for this lit too
			ty: ty,
		})

		// for _, field := range t.fields {
		// 	buf.PrintCompLit(&CompLitExpr{
		// 		// callee: // TODO: Not needed I dont think. could maybe do a fake ident node with just field name
		// 		args: nil, // So that it does the default for this lit too
		// 		ty: field,
		// 	})
		// }
	case *ArrayType:
		// TODO: Technically this can be written: {0}
		buf.Add("{0}")
		// buf.Add("{{")
		// buf.PrintDefault(t.base)
		// buf.Add("}}")
	case *PointerType:
		buf.Add("NULL")
	case *BasicType:
		buf.Add(t.Default())
		// TODO: Lookup
	default:
		panic(fmt.Sprintf("Unhandled Type: %T", ty))
	}
}

func (buf *genBuf) PrintCompLit(c *CompLitExpr) {
	ty := c.Type()
	if buf.indent > 0 {
		// Global variables use a different composit lit syntax
		buf.Add("(").
			Add(typeNameC(c.ty)).
			Add(")")
	}

	switch t := ty.(type) {
	case *StructType:
		buf.Add("{ ")
		for i := range t.fields {
			arg := c.GetArg(i)
			if arg != nil {
				buf.Print(arg)
			} else {
				buf.PrintDefault(t.fields[i])
			}

			if i < len(t.fields)-1 {
				buf.Add(", ")
			}
		}
		buf.Add(" }")
	case *ArrayType:
		if c.args == nil {
			buf.PrintDefault(t)
		} else {
			buf.Add("{{ ")
			for i := range t.len {
				arg := c.GetArg(i)
				if arg != nil {
					buf.Print(arg)
				} else {
					buf.PrintDefault(t.base)
				}

				if i < t.len-1 {
					buf.Add(", ")
				}
			}
			buf.Add(" }}")
		}
	default:
		panic(fmt.Sprintf("Unhandled Type: %T", ty))
	}
}

func (buf *genBuf) PrintStructNode(t *StructNode) {
	buf.Add("struct ").
		Add(t.ident.str).
		Add(" {").
		Line()

	buf.indent++
	for _, field := range t.fields {
		buf.Add(typeNameC(field.typeNode.Type())).
			Add(" ").
			Add(field.name.str).
			Add(";").
			Line()
	}
	buf.indent--

	buf.Add("}")
}

func (buf *genBuf) PrintBinaryExpr(t *BinaryExpr) {
	fmt.Println("HERERERE:", t.left.Type(), t.right.Type())
	if useCustomEqualityFunc(t.left.Type()) {
		ty := t.left.Type()
		buf.Add("(")
		buf.Add(equalityFunctionName(ty)).Add("(")
		buf.Print(t.left)
		buf.Add(", ")
		buf.Print(t.right)
		buf.Add(")")
		buf.Add(" ").Add(t.op.str).Add(" ")
		buf.Add("true")
		buf.Add(")")
	} else {
		// Simple binary expression
		buf.Add("(")
		buf.Print(t.left)
		buf.Add(" ").Add(t.op.str).Add(" ")
		buf.Print(t.right)
		buf.Add(")")
	}
}

func (buf *genBuf) Print(n Node) {
	switch t := n.(type) {
	case *FileNode:
		for _, nn := range t.nodes {
			buf.Print(nn)
		}
	case *PackageNode:
		buf.Add("// ").Add("package ").Add(t.name).
			Line()
	case *CommentNode:
	case *StructNode:
		if !t.global {
			buf.Add("typedef ")
			buf.PrintStructNode(t)
			buf.Add(" ").Add(t.ident.str)
			buf.Add(";").Line()
		}
	case *FuncNode:
		buf.LineDirective(t.pos)
		buf.PrintFuncDef(t)
		buf.Add(" {").Line()

		buf.Print(t.body)

		// Print Return
		if t.name == "main" {
			buf.Add("return __mainRet__;").Line()
		}
		buf.Add("}").Line()
	case *Stmt:
		buf.Print(t.node)
		// buf.Add(";").Line()
	case *ForStmt:
		buf.Add("for (")
		buf.Print(t.init)
		buf.Add("; ")
		buf.Print(t.cond)
		buf.Add("; ")
		buf.Print(t.inc)
		buf.Add(") {").Line()
		buf.Print(t.body)
		buf.Add("}")
	case *VarStmt:
		if !t.global {
			Println("VarStmt:", *t)
			buf.PrintVarDecl(t.name.str, t.Type(), t.initExpr)
		}
	case *ShortStmt:
		buf.Print(t.target)
		buf.Add(" ").Add(t.op.str).Add(" ")
		buf.Print(t.initExpr)

	case *IfStmt:
		buf.Add("if (")
		buf.Print(t.cond)
		buf.Add(") {").Line()

		buf.Print(t.thenScope)
		buf.Add("}")
		if t.elseScope == nil {
		} else {
			buf.Add(" else {").Line()
			buf.Print(t.elseScope)
			buf.Add("}")
		}

	case *ArgNode:
		if len(t.args) <= 0 {
			buf.Add("void")
		} else {
			for i := range t.args {
				buf.Add(typeNameC(t.args[i].ty)).
					Add(" ").
					Add(t.args[i].name.str).
					Add(" ")
				if i < len(t.args)-1 {
					buf.Add(", ")
				}
			}
		}
	case *ScopeNode:
		buf.Add("{").Line()
		buf.Print(t.Scope)
		buf.Add("}").Line()
	case *CurlyScope:
		buf.indent++
		for i := range t.nodes {
			buf.Print(t.nodes[i])
			buf.Add(";").Line()
		}
		buf.indent--

	case *ReturnNode:
		buf.Add("return (")
		buf.Print(t.expr)
		buf.Add(")") // buf.Add(");").Line()

	case *CallExpr:
		buf.Print(t.callee)
		buf.Add("(")
		buf.PrintArgList(t.args)
		buf.Add(")")
	case *GetExpr:
		buf.Print(t.obj)
		buf.Add(".")
		buf.Add(t.name.str)
	case *SetExpr:
		buf.Print(t.obj)
		buf.Add(".")
		buf.Add(t.name.str)
		buf.Add(" = ")
		buf.Print(t.value)
	case *BinaryExpr:
		Println(t.left, t.right)
		buf.PrintBinaryExpr(t)

	case *IndexExpr:
		buf.Print(t.callee)
		buf.Add(".a[")
		buf.Print(t.index)
		buf.Add("]")

	case *PostfixStmt:
		buf.Add("(")
		buf.Print(t.left)
		buf.Add(t.op.str)
		buf.Add(")")
	case *UnaryExpr:
		buf.Add("(").Add(t.op.str)
		buf.Print(t.right)
		buf.Add(")")
	case *AssignExpr:
		buf.Print(t.name)
		buf.Add(" = ")
		buf.Print(t.value)
	case *IdentExpr:
		buf.Add(t.tok.str)
	case *CompLitExpr:
		buf.PrintCompLit(t)

	case *GroupingExpr:
		buf.Add("(")
		buf.Print(t.Node)
		buf.Add(")")
	case *LitExpr:
		buf.Add(t.tok.str)
	default:
		panic(fmt.Sprintf("Print: Unknown NodeType: %T", t))
	}
}

func equalityFunctionName(ty Type) string {
	return "__ko_"+typeNameC(ty)+"_equality"
}

func (buf *genBuf) printEqualityPrototype(ty Type)  {
	buf.Add("bool ").Add(equalityFunctionName(ty)).
		Add("(").
		Add(typeNameC(ty)).Add(" a").
		Add(", ").
		Add(typeNameC(ty)).Add(" b").
		Add(")")

}
func (buf *genBuf) printStructEqualityFunction(t *StructNode) {
	ty := t.Type()
	buf.printEqualityPrototype(ty)
	buf.Add("{").Line()
	buf.indent++

	// buf.Add("return a == b;").Line()
	buf.Add("return (")
	for i, field := range t.fields {
		fname := field.name.str
		// buf.Add("(a.").Add(fname).Add(" == ").Add("b.").Add(fname).Add(")")
		expr := &BinaryExpr{
			// left: &IdentExpr{tok: Token{token: IDENT, str: "a"}, ty: field.ty},
			// right: &IdentExpr{tok: Token{token: IDENT, str: "b"}, ty: field.ty},
			left: &GetExpr{
				obj: &IdentExpr{tok: Token{token: IDENT, str: "a"}, ty: t.ty},
				name: Token{str: fname},
				ty: field.ty,
			},
			right: &GetExpr{
				obj: &IdentExpr{tok: Token{token: IDENT, str: "b"}, ty: t.ty},
				name: Token{str: fname},
				ty: field.ty,
			},
			op: Token{token: EQUALEQUAL, str: "=="},
			ty: UnknownType, // TODO: Should come from operator
		}
		buf.PrintBinaryExpr(expr)
		if i < len(t.fields)-1 {
			buf.Add(" && ")
		}
	}
	buf.Add(");").Line()

	buf.indent--
	buf.Add("}").Line()
}

func useCustomEqualityFunc(ty Type) bool {
	switch ty.(type) {
	case *BasicType:
		return false
	case *StructType:
		return true // Technically you need to ensure all fields are comparable
	case *ArrayType:
		return true // Technically you need to ensure all fields are comparable
	default:
		panic(fmt.Sprintf("Unknown Type: %T", ty))
	}

}

// Emits a variable declaration with name type and init expression
// If init is nil, emits the default value of the type
func (buf *genBuf) PrintVarDecl(name string, ty Type, init Node) {
	buf.Add(varDeclLHS(name, ty))
	buf.Add(" = ")

	// RHS
	if init == nil {
		buf.PrintDefault(ty)
	} else {
		buf.Print(init)
	}
}

// Returns the left hand side of a variable declaration
func varDeclLHS(name string, ty Type) string {
	typeStr := typeNameC(ty)
	return typeStr + " " + name

	// typeStr := typeNameC(ty)
	// switch t := ty.(type) {
	// // Arrays have a weird syntax in C
	// case *ArrayType:
	// 	base := varDeclLHS(name, t.base)
	// 	return fmt.Sprintf("%s[%d]", base, t.len)
	// default:
	// 	return typeStr + " " + name
	// }
}

// Returns the type name in C
// For arrays we return the type name of the base type
func typeNameC(ty Type) string {
	switch t := ty.(type) {
	case *BasicType:
		return typeStr(t)
	case *PointerType:
		return typeNameC(t.base)+"*"
	case *StructType:
		return t.Name()
	case *ArrayType:
		return fmt.Sprintf("__ko_%d%s_arr", t.len, typeNameC(t.base))
		// return typeNameC(t.base)
	default:
		panic(fmt.Sprintf("Unknown Type: %T", ty))
	}
	return ""
}
