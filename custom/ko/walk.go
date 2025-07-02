package main

import "fmt"

// func DebugTree(node Node, indent int) {
// 	switch t := node.(type) {
// 	case *PackageNode:
// 		fmt.Println("package:", t.name)
// 	case *FileNode:
// 		fmt.Println("file:", t.filename)
// 		indent++
// 		for _, n := range t.nodes {
// 			DebugTree(n, indent)
// 		}
// 		indent--
// 	case *CommentNode:
// 		// Nothing
// 	case *ForeignScope:
// 		fmt.Println("foreign")
// 		indent++
// 		DebugTree(t.body, indent)
// 		indent--
// 	case *CurlyScope:
// 		fmt.Println(
// 	default:
// 		fmt.Sprintf("Unknown NodeType: %T", t)
// 		panic(fmt.Sprintf("Unknown NodeType: %T", t))
// 	}
// }

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
