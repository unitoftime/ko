package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	fmt.Println("Starting")

	fset := token.NewFileSet()
	// TODO: Maybe you should be using ParseFile b/c go only supports .go files
	packages, err := parser.ParseDir(fset, "../src/", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	for _, pkg := range packages {
		fmt.Println("Parsing Package:", pkg.Name)

		cg := CodeGenerator{}
		ast.Walk(cg, pkg)
	}
}

func Generate(node ast.Node) {
	if node == nil { return }

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
