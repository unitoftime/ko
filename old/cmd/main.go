package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"os/exec"
)

// Ideas:
// - For GDB you can emit C code with #line directives to map the C code to the original code.
// - Basic raylib example: https://github.com/SimonLSchlee/zigraylib
// - Comptime as a way to enforce compiler invariants?
// - Comptime examples
// fn add2(a: anytype, b: anytype) @TypeOf(a) {
//     return a + b;
// }

// fn NewList(x: anytype) List(@TypeOf(x)) {
//     const T = @TypeOf(x);
//     var buffer: [10]T = undefined;

//     return List(T){
//         .items = &buffer,
//         .len = 0,
//     };
// }

// fn List(comptime T: type) type {
//     return struct {
//         items: []T,
//         len: usize,

//         pub fn init() List(T) {
//             var buffer: [10]T = undefined;
//             return .{
//                 .items = &buffer,
//                 .len = 0,
//             };
//         }

//         pub fn size(self: List(T)) usize {
//             return self.len;
//         }
//     };
// }



func main() {
	fmt.Println("Starting")

	fset := token.NewFileSet()
	// TODO: Maybe you should be using ParseFile b/c go only supports .go files
	packages, err := parser.ParseDir(fset, "src/", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	for _, pkg := range packages {
		fmt.Println("Parsing Package:", pkg.Name)

		// cg := CodeGenerator{}
		// ast.Walk(cg, pkg)
		var buf bytes.Buffer
		Generate(&buf, fset, pkg)
		fmt.Println("--------------------------------------------------------------------------------")
		fmt.Println("- Result")
		fmt.Println("--------------------------------------------------------------------------------")
		fmt.Println(buf.String())

		os.WriteFile("build/src/main.zig", buf.Bytes(), 0777)
	}

	cmd := exec.Command("zig", "build", "run")
	cmd.Dir = "./build/"

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		out, _ := io.ReadAll(stdout)
		e, _ := io.ReadAll(stderr)

		fmt.Fprint(os.Stdout, "--------------------------------------------------------------------------------\n")
		fmt.Fprint(os.Stdout, "- stdout\n")
		fmt.Fprint(os.Stdout, "--------------------------------------------------------------------------------\n")
		fmt.Fprint(os.Stdout, string(out))


		fmt.Fprint(os.Stderr, "--------------------------------------------------------------------------------\n")
		fmt.Fprint(os.Stdout, "- stderr\n")
		fmt.Fprint(os.Stderr, "--------------------------------------------------------------------------------\n")
		fmt.Fprint(os.Stderr, string(e))
		fmt.Fprint(os.Stderr, "--------------------------------------------------------------------------------\n")
	}()


	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Finish")
}

// The node type must be *[ast.File], *[CommentedNode], [][ast.Decl], [][ast.Stmt],
// or assignment-compatible to [ast.Expr], [ast.Decl], [ast.Spec], or [ast.Stmt].
func Generate(dst io.Writer, fset *token.FileSet, root ast.Node) {
	if root == nil { return }

	// fmt.Printf("Visit", node)
	p := newPrinter()

	for node := range ast.Preorder(root) {
		switch n := node.(type) {
		case *ast.Package:
			fmt.Println("Package", n)

		case *ast.File:
			fmt.Println("File", n)
			fmt.Println(n.Name.Name)
			p.file(n)
			dst.Write(p.buf.Bytes())

		// for _, importSpec := range file.Imports {
		// 	path := importSpec.Path.Value

		// 	// If there was a custom name, use that
		// 	nameIdent := importSpec.Name
		// 	if nameIdent != nil {
		// 		name = nameIdent.Name
		// 	}
		// }
		default:
			// fmt.Printf("Unhandled Type: %T: %+v\n", node, node)
		}
	}

	fmt.Printf("(len: %d): %v\n", len(p.symbolLut["main.zig"]), p.symbolLut)
	lines := p.symbolLut["main.zig"]
	for i := range lines {
		outputLineNum := i + 1
		fmt.Println(fset.Position(lines[i].Pos), ": ", outputLineNum)
	}
}

func FmtNode(dst io.Writer, fset *token.FileSet, node ast.Node) {
	// format node
	switch node.(type) {
	// case ast.Expr:
	// 	p.expr(n)
	// case ast.Stmt:
	// 	// A labeled statement will un-indent to position the label.
	// 	// Set p.indent to 1 so we don't get indent "underflow".
	// 	if _, ok := n.(*ast.LabeledStmt); ok {
	// 		p.indent = 1
	// 	}
	// 	p.stmt(n, false)
	// case ast.Decl:
	// 	p.decl(n)
	// case ast.Spec:
	// 	p.spec(n, 1, false)
	// case []ast.Stmt:
	// 	// A labeled statement will un-indent to position the label.
	// 	// Set p.indent to 1 so we don't get indent "underflow".
	// 	for _, s := range n {
	// 		if _, ok := s.(*ast.LabeledStmt); ok {
	// 			p.indent = 1
	// 		}
	// 	}
	// 	p.stmtList(n, 0, false)
	// case []ast.Decl:
	// 	p.declList(n)
	// case *ast.File:
	// 	p.file(n)
	// default:
	// 	goto unsupported
	default:
		fmt.Printf("FmtNode: Unhandled Type: %T: %+v\n", node, node)
	}
}

type CodeGenerator struct {
	buf bytes.Buffer
}
func (v CodeGenerator) Visit(node ast.Node) ast.Visitor {
	if node == nil { return nil }

	// fmt.Printf("Visit", node)

	switch n := node.(type) {
	case *ast.Package:
		return v // If we are a package, then just keep searching

	case *ast.File:
		fmt.Println(n.Name.Name)
		// for _, importSpec := range file.Imports {
		// 	path := importSpec.Path.Value

		// 	// If there was a custom name, use that
		// 	nameIdent := importSpec.Name
		// 	if nameIdent != nil {
		// 		name = nameIdent.Name
		// 	}
		// }
	default:
		fmt.Printf("Unhandled Type: %T: %+v\n", node, node)
	}

	// _, ok := node.(*ast.Package)
	// if ok { return v }

	// // If we are a file, then store some data in the visitor so we can use it later
	// file, ok := node.(*ast.File)
	// if ok {
	// 	fmt.Printf(file.Name)

	// 	for _, importSpec := range file.Imports {
	// 		path := importSpec.Path.Value

	// 		// If there was a custom name, use that
	// 		nameIdent := importSpec.Name
	// 		if nameIdent != nil {
	// 			name = nameIdent.Name
	// 		}
	// 	}
	// 	return v
	// }

	// // If we are a function, do the function formatting
	// _, ok = node.(*ast.FuncDecl)
	// if ok {
	// 	return nil // Skip: we don't handle funcs
	// }

	// gen, ok := node.(*ast.GenDecl)
	// if ok {
	// 	cgroups := v.cmap.Filter(gen).Comments()
	// 	sd, ok := v.formatGen(*gen, cgroups)
	// 	if ok {
	// 		v.structs[sd.Name] = sd
	// 	}


	// 	return nil
	// }

	// If all else fails, then keep looking
	return v
}
