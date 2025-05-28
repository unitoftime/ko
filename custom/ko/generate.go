package main

import (
	"bytes"
	"fmt"

	_ "embed"
)

func WalkNode(n Node) {
	switch t := n.(type) {
	case *PackageNode:
		fmt.Println("PackageNode:", t.name)
	case *CommentNode:
		fmt.Println("CommentNode:")
	case *FileNode:
		fmt.Println("FileNode:", t.filename)
		for _, nn := range t.nodes {
			WalkNode(nn)
		}
	case *Stmt:
		fmt.Println(t.node)
	case *FuncNode:
		fmt.Println("FuncNode:", t.name)
		WalkNode(t.arguments)
		WalkNode(t.body)
	case *ArgNode:
		for i := range t.args {
			fmt.Println("Arg:", t.args[i])
		}
	case *CurlyScope:
		fmt.Println("CurlyScope")
		for i := range t.nodes {
			WalkNode(t.nodes[i])
		}
	// case *ExprNode:
	// 	for i := range t.ops {
	// 		WalkNode(t.ops[i])
	// 	}
	// case *UnaryNode:
	// 	fmt.Println("Unary", t.token)
	case *ReturnNode:
		fmt.Println("Return")
		WalkNode(t.expr)
	case *BinaryExpr:
		fmt.Println("Binary", t.op)
		WalkNode(t.left)
		WalkNode(t.right)
	case *LitExpr:
		fmt.Println("Lit:", t.tok)
	default:
		fmt.Sprintf("Unknown NodeType: %T", t)
		panic(fmt.Sprintf("Unknown NodeType: %T", t))
	}
}

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
	b.Add(fmt.Sprintf("#line %d \"%s\"", pos.line, pos.filename))
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
	if argNode == nil { return "void" }

	if !ok { panic(fmt.Sprintf("must be arg node: %T", node)) }
	if len(argNode.args) != 1 {
		panic("only supporting single return types")
	}
	fmt.Printf("returnArgsToString: %+v\n", argNode)
	fmt.Printf("returnArgsToString: %T\n", argNode.args[0])
	return typeStr(argNode.args[0].ty)
}

//go:embed runtime.h
var runtimeFile string


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
	}

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

func (buf *genBuf) PrintForwardDecl(n Node) {
	switch t := n.(type) {
	case *StructNode:
		buf.Add("typedef struct ").
			Add(t.ident.str).
			Add(" ").
			Add(t.ident.str)
	case *VarStmt:
		typeStr := typeStr(t.ty)
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
		buf.Add(";").Line()
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

func (buf *genBuf) PrintStructNode(t *StructNode) {
	buf.Add("struct ").
		Add(t.ident.str).
		Add(" {").
		Line()

	buf.indent++
	for _, field := range t.fields {
		buf.Add(field.kind.str).
			Add(" ").
			Add(field.name.str).
			Add(";").
			Line()
	}
	buf.indent--

	buf.Add("}")
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
			buf.PrintStructNode(t)
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
			fmt.Println("VarStmt:", *t)
			typeStr := typeStr(t.ty)

			buf.
				Add(typeStr).
				Add(" ").
				Add(t.name.str).
				Add(" = ")
			buf.Print(t.initExpr)
		}

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
		for i := range t.args {
			buf.Add(typeStr(t.args[i].ty)).
				Add(" ").
				Add(t.args[i].name.str).
				Add(" ")
			if i < len(t.args)-1 {
				buf.Add(", ")
			}
		}
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
		fmt.Println(t.left, t.right)
		if t.left.Type().isStruct || t.right.Type().isStruct {
			ty := t.left.Type()
			buf.Add("(")
			buf.Add("__ko_" + ty.name +"_equality(")
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

	case *UnaryExpr:
		buf.Add("(").Add(t.op.str)
		buf.Print(t.right)
		buf.Add(")")
	case *AssignExpr:
		buf.Add(t.name.str)
		buf.Add(" = ")
		buf.Print(t.value)
	case *IdentExpr:
		buf.Add(t.tok.str)
	case *CompLitExpr:
		buf.Add("(").
			Add(typeStr(t.ty)).
			Add(")")
		buf.Add("{ ")
		buf.PrintArgList(t.args)
		buf.Add(" }")

	case *LitExpr:
		buf.Add(t.tok.str)
	default:
		panic(fmt.Sprintf("Print: Unknown NodeType: %T", t))
	}
}

func equalityFunctionName(name string) string {
	return "__ko_"+name+"_equality"
}

func (buf *genBuf) printStructEqualityFunction(t *StructNode) {
	ty := t.Type()
	buf.Add("bool ").Add(equalityFunctionName(ty.name)).
		Add("(").
		Add(ty.name).Add(" a").
		Add(", ").
		Add(ty.name).Add(" b").
		Add(") {").Line()
	buf.indent++

	// buf.Add("return a == b;").Line()
	buf.Add("return (")
	for i, field := range t.fields {
		fname := field.name.str
		buf.Add("(a.").Add(fname).Add(" == ").Add("b.").Add(fname).Add(")")
		if i < len(t.fields)-1 {
			buf.Add(" && ")
		}
	}
	buf.Add(");").Line()

	buf.indent--
	buf.Add("}").Line()
}
