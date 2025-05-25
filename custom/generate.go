package main

import (
	"bytes"
	"fmt"
)

func WalkNode(n Node) {
	switch t := n.(type) {
	case PackageNode:
		fmt.Println("PackageNode:", t.name)
	case CommentNode:
		fmt.Println("CommentNode:")
	case *FileNode:
		fmt.Println("FileNode:", t.filename)
		for _, nn := range t.nodes {
			WalkNode(nn)
		}
	case Stmt:
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
	case BinaryExpr:
		fmt.Println("Binary", t.op)
		WalkNode(t.left)
		WalkNode(t.right)
	case LitExpr:
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

func (b *genBuf) String() string {
	return b.buf.String()
}

// TODO: if a tuple, then return a struct
func returnArgsToString(node Node) string {
	if node == nil { return "void" }

	argNode, ok := node.(ArgNode)
	if !ok { panic(fmt.Sprintf("must be arg node: %T", node)) }
	if len(argNode.args) != 1 {
		panic("only supporting single return types")
	}
	return argNode.args[0].kind
}


func GenerateCode(buf *genBuf, n Node) {
	switch t := n.(type) {
	case *FileNode:
		buf.Add("#include <stdio.h>"). // TODO: hack
			Line().Line()

		for _, nn := range t.nodes {
			GenerateCode(buf, nn)
		}
	case PackageNode:
		buf.Add("// ").Add("package ").Add(t.name).
			Line()
	case CommentNode:
	case *FuncNode:
		buf.Add(returnArgsToString(t.returns)).
			Add(" ").
			Add(t.name).
			Add(" ( ")

		GenerateCode(buf, t.arguments)
		buf.Add(" ) {").Line()

		GenerateCode(buf, t.body)

		buf.Add("}").Line()
	case Stmt:
		GenerateCode(buf, t.node)
		// buf.Add(";").Line()
	case ForStmt:
		buf.Add("for (")
		GenerateCode(buf, t.init)
		buf.Add("; ")
		GenerateCode(buf, t.cond)
		buf.Add("; ")
		GenerateCode(buf, t.inc)
		buf.Add(") {").Line()
		GenerateCode(buf, t.body)
		buf.Add("}")
	case VarStmt:
		typeStr := "int"
		buf.
			Add(typeStr).
			Add(" ").
			Add(t.name.str).
			Add(" = ")
		GenerateCode(buf, t.initExpr)

	case IfStmt:
		buf.Add("if (")
		GenerateCode(buf, t.cond)
		buf.Add(") {").Line()

		GenerateCode(buf, t.thenScope)
		buf.Add("}")
		if t.elseScope == nil {
		} else {
			buf.Add(" else {").Line()
			GenerateCode(buf, t.elseScope)
			buf.Add("}")
		}

	case *ArgNode:
		for i := range t.args {
			buf.Add(t.args[i].kind).
				Add(" ").
				Add(t.args[i].name).
				Add(" ")
			if i < len(t.args)-1 {
				buf.Add(", ")
			}
		}
	case *CurlyScope:
		buf.indent++
		for i := range t.nodes {
			GenerateCode(buf, t.nodes[i])
			buf.Add(";").Line()
		}
		buf.indent--

	case *ReturnNode:
		buf.Add("return (")
		GenerateCode(buf, t.expr)
		buf.Add(")") // buf.Add(");").Line()

	case CallExpr:
		GenerateCode(buf, t.callee)
		buf.Add("( ")
		for i := range t.args {
			GenerateCode(buf, t.args[i])

			if i < len(t.args)-1 {
				buf.Add(", ")
			}
		}
		buf.Add(")")

	case BinaryExpr:
		buf.Add("(")
		GenerateCode(buf, t.left)
		buf.Add(" ").Add(t.op.str).Add(" ")
		GenerateCode(buf, t.right)
		buf.Add(")")
	case AssignExpr:
		buf.Add(t.name.str)
		buf.Add(" = ")
		GenerateCode(buf, t.value)
	case VarExpr:
		buf.Add(t.tok.str)
	case LitExpr:
		buf.Add(t.tok.str)
	default:
		panic(fmt.Sprintf("Unknown NodeType: %T", t))
	}
}
