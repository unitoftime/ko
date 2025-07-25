package main

import (
	"fmt"
	"html/template"
	"io"

	_ "embed"
)

type genBuf struct {
	// buf *bytes.Buffer
	buf io.Writer
	indent int
	newline bool
	lineHasContent bool
}
func (b *genBuf) Add(str string) *genBuf {
	if b.newline {
		b.newline = false
		b.lineHasContent = true
		for range b.indent {
			io.WriteString(b.buf, "\t")
		}
	}

	io.WriteString(b.buf, str)
	return b
}

func (b *genBuf) SemiLine() *genBuf {
	if b.lineHasContent {
		b.Add(";")
	}
	b.Line()
	return b
}

func (b *genBuf) Line() *genBuf {
	b.Add("\n")
	b.newline = true
	b.lineHasContent = false
	return b
}
func (b *genBuf) LineDirective(pos Position) *genBuf {
	// #line 31 "test.txt"
	// b.Add(fmt.Sprintf("#line %d \"%s\"", pos.line, pos.filename))
	b.Line()

	return b
}

// func (b *genBuf) String() string {
// 	return b.buf.String()
// }

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

//go:embed array.tmpl
var arrayTemplate string

//go:embed slice.tmpl
var sliceTemplate string
type SliceTemplateDef struct {
	Name string
	Type string
	Size int
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
		buf.SemiLine()

		// structNode, isStruct := node.(*StructNode)
		// if isStruct {
		// 	buf.printEqualityPrototype(structNode.Type())
		// 	buf.Add(";").Line()
		// }

		switch t := node.(type) {
		case *StructNode:
			buf.printEqualityPrototype(t.Type())
			buf.SemiLine()
		}
	}

	// Forward declare all arrays
	for _, ty := range regTypeMap {
		// TODO: You could probably register *special* types you find during typechecking, rather than looping everything
		Printf("RegTypeMap: %T\n", ty)
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
		buf.SemiLine()
	}
	for i := range result.genericInstantiations {
		buf.PrintGenericForwardDecl(result.genericInstantiations[i])
		buf.SemiLine()
	}


	// Complete all types
	for _, node := range result.typeList {
		buf.PrintCompleteType(node)
		buf.SemiLine()

		structNode, isStruct := node.(*StructNode)
		if isStruct {
			buf.printStructEqualityFunction(structNode)
		}
	}

	for i := range result.genericInstantiations {
		buf.PrintCompleteGenericType(result.genericInstantiations[i])
	}

	// Declare all global variables
	for i := range result.varList {
		buf.LineDirective(result.varList[i].name.pos)
		buf.PrintForwardDecl(result.varList[i])
		buf.SemiLine()
	}

	buf.Print(result.file)
}

func (buf *genBuf) PrintFuncDef(t *FuncNode) {
	isMain := t.name == "main"
	if isMain {
		buf.Add("#include \"args.h\"").Line()
		buf.Add("int main(int argc, const char* argv[])")
	} else {
		retArgs := returnArgsToString(t.returns)

		buf.Add(retArgs).
			Add(" ").
			Add(t.name).
			Add(" (")
		buf.Print(t.arguments)
		buf.Add(")")
	}
}

func (buf *genBuf) PrintGeneratedType(ty Type) {
	switch t := ty.(type) {
	case *SliceType:
		name := typeNameC(ty)
		elemName := typeNameC(t.base)

		buf.PrintStructForwardDecl(name)
		buf.SemiLine()
		tmpl := template.Must(template.New("cslice").Parse(sliceTemplate))
		err := tmpl.Execute(buf.buf, SliceTemplateDef{
			Name: name,
			Type: elemName,
		})
		if err != nil {
			panic(err)
		}
		buf.Line()
	case *ArrayType:
		name := typeNameC(ty)
		elemName := typeNameC(t.base)

		// Forward Declaration
		// TODO: Would be good to hoist this all up so nested arrays arent problematic
		buf.Add("typedef struct ").
			Add(name).Add(" ").Add(name)
		buf.SemiLine()

		tmpl := template.Must(template.New("carray").Parse(arrayTemplate))
		err := tmpl.Execute(buf.buf, SliceTemplateDef{
			Name: name,
			Type: elemName,
			Size: t.len,
		})
		if err != nil {
			panic(err)
		}
		buf.Line()
	}
}

func (buf *genBuf) PrintStructForwardDecl(name string) {
	buf.Add("typedef struct ").Add(name).Add(" ").Add(name)
}

func (buf *genBuf) PrintEnumForwardDecl(name string) {
	// TODO: enum type/size from node contents
	buf.Add("typedef uint8_t ").Add(name)
}

func (buf *genBuf) PrintForwardDecl(n Node) bool {
	switch t := n.(type) {
	case *StructNode:
		if t.foreign { return false }
		buf.PrintStructForwardDecl(t.ident.str)
	case *EnumNode:
		if t.foreign { return false }
		buf.PrintEnumForwardDecl(t.ident.str) // TODO
	case *VarStmt:
		if t.foreign { return false }
		typeStr := typeNameC(t.ty)
		buf.
			Add(typeStr).
			Add(" ").
			Add(t.name.str)

		if t.initExpr != nil {
			buf.Add(" = ")
			buf.Print(t.initExpr)
		} else {
			// TODO: possible composit literal: Assign all default values
		}
	case *FuncNode:
		if t.body == nil { return false }
		if t.Generic() { return false } // TODO: Eventually handle these
		if t.name == "main" { return false }

		buf.LineDirective(t.pos)
		buf.PrintFuncDef(t)

	default:
		panic(fmt.Sprintf("PrintForwardDecl: Unknown NodeType: %T", t))
	}

	return true
}

func (buf *genBuf) PrintGenericForwardDecl(g GenericInstance) {
	ty := g.node.Type()
	switch t := ty.(type) {
	case *FuncType:
		defer clearGenericMap()
		for i, genArg := range g.funcNode.generic.Args {
			addGenericMap(genArg.name.str, t.args[i])
		}

		newFuncNode := g.funcNode.ToConcrete(t)

		// buf.LineDirective(g.funcNode.pos)
		buf.PrintForwardDecl(newFuncNode)
	default:
		panic(fmt.Sprintf("PrintGenericForwardDecl: Unknown Type: %T", ty))
	}
}

func (buf *genBuf) PrintCompleteGenericType(g GenericInstance) {
	ty := g.node.Type()
	switch t := ty.(type) {
	case *FuncType:
		defer clearGenericMap()
		for i, genArg := range g.funcNode.generic.Args {
			addGenericMap(genArg.name.str, t.args[i])
		}

		newFuncNode := g.funcNode.ToConcrete(t)

		buf.Print(newFuncNode)
	default:
		panic(fmt.Sprintf("PrintCompleteGenericType: Unknown Type: %T", ty))
	}
}


func (buf *genBuf) PrintCompleteType(n Node) {
	switch t := n.(type) {
	case *StructNode:
		if t.foreign { return }
		buf.PrintStructNode(t)
	case *EnumNode:
		if t.foreign { return }
		buf.PrintEnumNode(t)
	default:
		panic(fmt.Sprintf("PrintForwardDecl: Unknown NodeType: %T", t))
	}
}

func (buf *genBuf) PrintCallExpr(ce *CallExpr) {
	switch t := ce.callee.Type().(type) {
	case *FuncType:
		buf.Print(ce.callee)
		buf.Add("(")
		buf.PrintArgList(ce.args)
		buf.Add(")")
	case *BasicType:
		buf.Add("(")
		buf.Add(typeNameC(t))
		buf.Add(")")
		buf.Add("(")
		buf.PrintArgList(ce.args)
		buf.Add(")")
	case *PointerType:
		buf.Add("(")
		buf.Add(typeNameC(t))
		buf.Add(")")
		buf.Add("(")
		buf.PrintArgList(ce.args)
		buf.Add(")")
	default:

		panic(fmt.Sprintf("PrintCallExpr: Unknown NodeType: %T", t))
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
	case *SliceType:
		buf.Add("{0}") // TODO: I guess this would end up being ptr=nil, len=0, cap=0
	case *PointerType:
		buf.Add("NULL")
	case *BasicType:
		buf.Add(t.Default())
		// TODO: Lookup
	case *GenericType:
		concType, ok := genericTypeMap[t.name]
		if !ok { panic("Unknown generic impl") }
		buf.PrintDefault(concType)
	default:
		panic(fmt.Sprintf("Unhandled Type: %T", ty))
	}
}

func (buf *genBuf) PrintCompLit(c *CompLitExpr) {
	ty := c.Type()
	switch t := ty.(type) {
	case *StructType:
		if buf.indent > 0 {
			// Global variables use a different composit lit syntax
			buf.Add("(").
				Add(typeNameC(c.ty)).
				Add(")")
		}
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
		if buf.indent > 0 {
			// Global variables use a different composit lit syntax
			buf.Add("(").
				Add(typeNameC(c.ty)).
				Add(")")
		}
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
	case *SliceType:
		if c.args == nil {
			buf.PrintDefault(t)
		} else {
			// buf.Add("__ko_int_slice_init((int[]){1, 2, 3, 4}, 4)")
			buf.Add("__ko_int_slice_init(")
			// arrayCL := *c
			// arrayCL.ty = &ArrayType{len: len(c.args), base: t.base}
			buf.Add("(").Add(typeNameC(t.base)).Add("[]){")
			for i := range len(c.args) {
				arg := c.GetArg(i)
				if arg != nil {
					buf.Print(arg)
				} else {
					buf.PrintDefault(t.base)
				}

				if i < len(c.args)-1 {
					buf.Add(", ")
				}
			}
			buf.Add("}").Add(", ").Add(fmt.Sprintf("%d", len(c.args))).Add(")")
		}
	default:
		panic(fmt.Sprintf("Unhandled Type: %T", ty))
	}
}

func (buf *genBuf) PrintEnumNode(t *EnumNode) {
	buf.Add("enum ").
		Add(t.ident.str).
		Add(" {")
	for i := range t.fields {
		// TODO: Might need to handle fields with set values
		buf.Add(t.fields[i].str)

		if i < len(t.fields)-1 {
			buf.Add(", ")
		}
	}
	buf.Add("}").SemiLine()
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

func (buf *genBuf) PrintIndexExpr(t *IndexExpr) {
	switch tt := t.callee.Type().(type) {
	case *ArrayType:
		buf.Print(t.callee)
		buf.Add(".a[")
		buf.Print(t.index)
		buf.Add("]")
	case *SliceType:
		buf.Print(t.callee)
		buf.Add(".a[")
		buf.Print(t.index)
		buf.Add("]")
	case *FuncType:
		// buf.Print(t.callee)
		// TODO: I think this is technically wrong, I don't want to decide the callee name based on the type, I want to decide it based on the node

		// If we are here it means that we are doing a compile time instantiation of the functype
		// The name is established and will be generated elsewhere, so we just need to emit that
		if tt.name == "append" {
			buf.Add("__ko_").
				Add(typeNameC(tt.generics[0])). // TODO: a bit hacky
				Add("_slice_append")
		} else if tt.name == "slice" {
			buf.
				Add(typeNameC(tt.generics[0])). // TODO: a bit hacky
				Add("_slice")
		} else if tt.name == "len" {
			buf.Add("__ko_").
				Add(typeNameC(tt.generics[0])). // TODO: a bit hacky
				Add("_slice_len")
		} else if tt.name == "make" {
			buf.
				Add(typeNameC(tt.generics[0])). // TODO: a bit hacky
				Add("_new")
		} else {
			buf.Add(t.Type().Name())
		}
	default:
		panic(fmt.Sprintf("Unknown Type: %T", t.callee.Type()))
	}
}

func (buf *genBuf) PrintBinaryExpr(t *BinaryExpr) {
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
	case *ForeignScope:
		// Skip: Externally defined
	case *StructNode:
		if !t.global {
			buf.Add("typedef ")
			buf.PrintStructNode(t)
			buf.Add(" ").Add(t.ident.str)
			buf.SemiLine()
		}
	case *EnumNode:
		if !t.global {
			buf.PrintEnumNode(t)
		}
	case *FuncNode:
		if t.Generic() { return } // TODO: Eventually handle these

		buf.LineDirective(t.pos)
		buf.PrintFuncDef(t)
		buf.Add(" {").Line()

		if t.name == "main" {
			// Handle argc/argv
			buf.Add("\t").Add("__ko_init_args(argc, argv);").Line()
		}

		buf.Print(t.body)

		// Print Return
		if t.name == "main" {
			buf.Add("\t").Add("return __mainRet__;").Line()
		}
		buf.Add("}").Line()
	case *Stmt:
		buf.Print(t.node)
		// buf.Add(";").Line()
	case *ForStmt:
		buf.Add("for (")
		if t.init != nil {
			buf.Print(t.init)
		}
		buf.Add("; ")
		if t.cond != nil {
			buf.Print(t.cond)
		}
		buf.Add("; ")
		if t.inc != nil {
			buf.Print(t.inc)
		}
		buf.Add(") {").Line()
		buf.Print(t.body)
		buf.Add("}")
	case *VarStmt:
		if !t.global && !t.constant {
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
	case *SwitchStmt:
		buf.Add("switch (")
		buf.Print(t.cond)
		buf.Add(") {")
		buf.Line()

		for i := range t.cases {
			buf.Print(t.cases[i])
		}

		buf.Add("}")
	case *CaseStmt:
		if t.expr == nil {
			buf.Add("default: {").Line()
		} else {
			buf.Add("case ")
			buf.Print(t.expr)
			buf.Add(": {")
			buf.Line()
		}

		buf.Print(t.body)
		buf.Add("\t").Add("break;")
		buf.Line().Add("}")
		buf.Line()

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
			buf.SemiLine()
		}
		buf.indent--

	case *ReturnNode:
		buf.Add("return (")
		buf.Print(t.expr)
		buf.Add(")") // buf.Add(");").Line()

	case *CallExpr:
		buf.PrintCallExpr(t)
	case *GetExpr:
		buf.selectorExpr(t.obj, t.name.str)
	case *SetExpr:
		buf.selectorExpr(t.obj, t.name.str)
		buf.Add(" = ")
		buf.Print(t.value)
	case *BinaryExpr:
		Println(t.left, t.right)
		buf.PrintBinaryExpr(t)

	case *IndexExpr:
		buf.PrintIndexExpr(t)

	case *PostfixStmt:
		buf.Add("(")
		buf.Print(t.left)
		buf.Add(t.op.str)
		buf.Add(")")
	case *UnaryExpr:
		buf.Add("(").Add(t.op.str).Add("(")
		buf.Print(t.right)
		buf.Add("))")
	case *AssignExpr:
		buf.Print(t.name)
		buf.Add(" = ")
		buf.Print(t.value)
	case *IdentExpr:
		if t.folded != nil {
			buf.Print(t.folded)
			break
		}

		// TODO: cleanup how we check for builtin func calls
		if t.tok.str == "append" {
			ft, ok := t.ty.(*FuncType)
			if ok {
				buf.Add("__ko_").
					Add(typeNameC(ft.generics[0])). // TODO: a bit hacky
					Add("_slice_append")
			}
			break
		}

		buf.Add(t.tok.str)
	case *CompLitExpr:
		buf.PrintCompLit(t)

	case *GroupingExpr:
		buf.Add("(")
		buf.Print(t.Node)
		buf.Add(")")
	case *LogicalExpr:
		buf.Add("(")
		buf.Print(t.left)
		buf.Add(")")

		buf.Add(t.op.str)

		buf.Add("(")
		buf.Print(t.right)
		buf.Add(")")
	case *LitExpr:
		if t.tok.token == NIL {
			buf.Add("NULL")
		} else if t.tok.token == STRING {
			buf.Add("__ko_string_make(").
				Add(t.tok.str).
				Add(")")
		} else {
			buf.Add(t.tok.str)
		}
	case *BreakStmt:
		buf.Add("break;").Line()
	default:
		panic(fmt.Sprintf("Print: Unknown NodeType: %T", t))
	}
}

func (buf *genBuf) selectorExpr(obj Node, field string) {
	buf.Print(obj)

	if obj.Type() == nil {
		Printf("SelectorType nil: %s: %T (in %T)\n", field, obj.Type(), obj)

		panic("Missing TYPE")
	}

	switch obj.Type().(type) {
	case *PointerType:
		buf.Add("->")
	default:
		buf.Add(".")
	}

	buf.Add(field)
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
	if t.foreign { return }

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

// TODO: This is also used by resolver
// TODO: shouldn't be global
var genericTypeMap = make(map[string]Type)
func addGenericMap(genericName string, concreteType Type) {
	// TODO: validate to make sure it hasn't been added before
	genericTypeMap[genericName] = concreteType
}

func clearGenericMap() {
	genericTypeMap = make(map[string]Type)
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
	case *SliceType:
		return fmt.Sprintf("__ko_%s_slice", typeNameC(t.base))
	case *GenericType:
		concreteType, ok := genericTypeMap[t.name]
		if !ok {
			panic(fmt.Sprintf("Missing generic mapping for: %s", t.name))
		}
		return typeNameC(concreteType)
	default:
		panic(fmt.Sprintf("Unknown Type: %T", ty))
	}
	return ""
}

func useCustomEqualityFunc(ty Type) bool {
	switch t := ty.(type) {
	case *BasicType:
		return false
	case *StructType:
		return true // Technically you need to ensure all fields are comparable
	case *ArrayType:
		return true // Technically you need to ensure all fields are comparable
	case *PointerType:
		return false // TODO: Should I find the base type and compare those? or just compare addresses?
	case *SliceType:
		return true // TODO: Should I find the base type and compare those? or just compare addresses?
	case *GenericType:
		concreteType, ok := genericTypeMap[t.name]
		if !ok {
			panic(fmt.Sprintf("Missing generic mapping for: %s", t.name))
		}
		return useCustomEqualityFunc(concreteType)
	default:
		panic(fmt.Sprintf("Unknown Type: %T", ty))
	}
}
