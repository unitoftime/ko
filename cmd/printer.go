package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"path"
	"strings"
)

// Zig
const (
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
	FUNC = "fn"
	LINE = "\n"
	SPACE = " "
	TUPLESEP = ", "
	ELLIPSIS = "..."
	SEMICOLON = ";"
	COLON = ":"

	// TYPE = "const"
	CONST = "const"
	VAR = "var"
	ASSIGN = "="
	STRUCT = "struct"
	INTERFACE = "interface"
	IMPORT = "@import"
)

// Go
// const (
// 	LPAREN = "("
// 	RPAREN = ")"
// 	LBRACE = "{"
// 	RBRACE = "}"
// 	FUNC = "func"
// 	LINE = "\n"
// 	SPACE = " "
// 	TUPLESEP = ", "
// 	ELLIPSIS = "..."
// 	SEMICOLON = ";"
	// CONST = "const"
	// VAR = "var"
	// ASSIGN = "="
	// STRUCT = "struct"
	// INTERFACE = "interface"
// )

type printer struct {
	buf bytes.Buffer
}

func (p *printer) file(src *ast.File) {
	// p.print("package ", src.Name.Name, LINE)
	// fmt.Fprintf(&p.buf, "package %s", src.Name.Name) // Expr technically

	p.declList(src.Decls)
}

func (p *printer) declList(list []ast.Decl) {
	for _, d := range list {
		p.decl(d)
	}
}

func (p *printer) decl(decl ast.Decl) {
	switch d := decl.(type) {
	// case *ast.BadDecl:
		// p.setPos(d.Pos())
		// p.print("BadDecl") // TODO: Unsure
	case *ast.GenDecl:
		p.genDecl(d)
		p.print(LINE)

	case *ast.FuncDecl:
		p.funcDecl(d)
	default:
		panic(fmt.Sprintf("decl: Unknown type: %T", decl))
	}
}

func (p *printer) genDecl(d *ast.GenDecl) {
	switch d.Tok {
	case token.IMPORT:
	case token.TYPE:
	case token.VAR:
		p.print(VAR, SPACE)
	case token.CONST:
		p.print(CONST, SPACE)
	default:
		panic(fmt.Sprintf("unhandled genDecl: %s", d.Tok.String()))
	}

	for i := range d.Specs {
		p.spec(d.Specs[i])
		p.print(LINE)
	}


	// // p.setComment(d.Doc)
	// // p.setPos(d.Pos())
	// p.print(d.Tok, " ")

	// if d.Lparen.IsValid() || len(d.Specs) != 1 {
	// 	// group of parenthesized declarations
	// 	// p.setPos(d.Lparen)
	// 	p.print(LPAREN)
	// 	if n := len(d.Specs); n > 0 {
	// 		p.print(LINE)
	// 		// p.print(indent, formfeed)
	// 		if n > 1 && (d.Tok == token.CONST || d.Tok == token.VAR) {
	// 			// // two or more grouped const/var declarations:
	// 			// // determine if the type column must be kept
	// 			// keepType := keepTypeColumn(d.Specs)
	// 			// var line int
	// 			for i, s := range d.Specs {
	// 				if i > 0 {
	// 					// p.linebreak(p.lineFor(s.Pos()), 1, ignore, p.linesFrom(line) > 0)
	// 					p.print(LINE)
	// 				}
	// 				// p.recordLine(&line)
	// 				p.valueSpec(s.(*ast.ValueSpec))
	// 			}
	// 		} else {
	// 			// var line int
	// 			for i, s := range d.Specs {
	// 				if i > 0 {
	// 				// 	p.linebreak(p.lineFor(s.Pos()), 1, ignore, p.linesFrom(line) > 0)
	// 					p.print(LINE)
	// 				}
	// 				// p.recordLine(&line)
	// 				p.spec(s)
	// 			}
	// 		}
	// 		// p.print(unindent, formfeed)
	// 	}
	// 	// p.setPos(d.Rparen)
	// 	p.print(LINE, RPAREN, LINE)

	// } else if len(d.Specs) > 0 {
	// 	// single declaration
	// 	p.spec(d.Specs[0])
	// }
}

func (p *printer) funcDecl(d *ast.FuncDecl) {
	// p.setComment(d.Doc)
	// p.setPos(d.Pos())
	p.print(LINE)

	if d.Name.Name == "main" {
		// If public
		// TODO: also handle capital named functions
		p.print("pub", SPACE, FUNC, SPACE)
	} else {
		p.print(FUNC, SPACE)
	}
	// We have to save startCol only after emitting FUNC; otherwise it can be on a
	// different line (all whitespace preceding the FUNC is emitted only when the
	// FUNC is emitted).
	// startCol := p.out.Column - len("func ")
	if d.Recv != nil {
		p.parameters(d.Recv)
		p.print(SPACE)
	}

	p.expr(d.Name)

	p.signature(d.Type)
	p.funcBody(d.Body)
}

func (p *printer) returnParams(fields *ast.FieldList) {
	if fields == nil { return }
	// p.print(LPAREN)
	for _, f := range fields.List {
		// for i := range f.Names {
		// 	p.expr(f.Names[i])
		// 	p.print(COLON, SPACE)
		// 	p.expr(f.Type)
		// 	if i < len(f.Names)-1 {
		// 		p.print(TUPLESEP)
		// 	}
		// }
		p.expr(f.Type)
	}
	// p.print(RPAREN)
}

func (p *printer) parameters(fields *ast.FieldList) {
	if fields == nil { return }
	p.print(LPAREN)
	for _, f := range fields.List {
		for i := range f.Names {
			p.expr(f.Names[i])
			p.print(COLON, SPACE)
			p.expr(f.Type)
			if i < len(f.Names)-1 {
				p.print(TUPLESEP)
			}
		}
	}
	p.print(RPAREN)


	// if fields == nil { return }
	// p.print(LPAREN)
	// for _, f := range fields.List {
	// 	p.identList(f.Names)
	// 	p.print(SPACE)
	// 	p.expr(f.Type)
	// 	p.print(SPACE)
	// 	// TODO: Tags?
	// }
	// p.print(RPAREN)
}

func (p *printer) signature(sig *ast.FuncType) {
	p.parameters(sig.TypeParams)
	p.print(SPACE)
	p.parameters(sig.Params)
	p.print(SPACE)
	if sig.Results != nil {
		p.returnParams(sig.Results)
		// TODO: If the error type gets returned then !err or whatever
	} else {
		p.print("void")
	}
	p.print(SPACE)
}

func (p *printer) funcBody(block *ast.BlockStmt) {
	if block == nil { return }

	p.print(LBRACE, LINE)
	defer p.print(RBRACE, LINE)

	for _, s := range block.List {
		p.print(LINE)
		p.stmt(s)
	}
	p.print(LINE)
}

func (p *printer) fieldList(fields *ast.FieldList) {
	if fields == nil { return }
	p.print(LBRACE)
	p.print(LINE)
	for _, f := range fields.List {
		p.identList(f.Names)
		p.print(COLON, SPACE)
		p.expr(f.Type)
		p.print(LINE)
		// TODO: Tags?
	}
	p.print(RBRACE, SEMICOLON)
}

func (p *printer) stmt(stmt ast.Stmt) {
	switch s := stmt.(type) {
	case *ast.ExprStmt:
		p.expr(s.X)
		p.print(SEMICOLON)
	case *ast.DeclStmt:
		p.genDecl(s.Decl.(*ast.GenDecl))
		p.print(SEMICOLON)

	case *ast.AssignStmt:
		if len(s.Lhs) != 1 || len(s.Rhs) != 1 {
			// TODO: Could: https://ziggit.dev/t/return-multiple-values-from-a-function/460/2
			panic("currently only supporting single assignemtnts")
		}
		tok := VAR
		if s.Tok == token.CONST {
			tok = CONST
		}
		p.print(tok, SPACE)
		p.expr(s.Lhs[0])
		p.print(SPACE, ASSIGN, SPACE)
		p.expr(s.Rhs[0])
		p.print(SEMICOLON)
	case *ast.ReturnStmt:
		p.print("return ")
		p.exprList(s.Results)
		p.print(SEMICOLON)
	case *ast.DeferStmt:
		p.print("defer ")
		p.expr(s.Call)
		p.print(SEMICOLON)
	case *ast.ForStmt:
		p.forStmt(s)

	default:
		panic(fmt.Sprintf("stmt: unknown type: %T", stmt))
	}
}

func (p *printer) valueSpec(s *ast.ValueSpec) {
	p.identList(s.Names) // always present
	if s.Type != nil {
		p.print(SPACE)
		p.expr(s.Type)
	}
	if s.Values != nil {
		p.print(SPACE, ASSIGN, SPACE)
		p.exprList(s.Values)
	}

	// // p.setComment(s.Doc)
	// p.identList(s.Names) // always present
	// // extraTabs := 3
	// // if s.Type != nil || keepType {
	// // 	p.print(vtab)
	// // 	extraTabs--
	// // }
	// if s.Type != nil {
	// 	p.expr(s.Type)
	// }
	// // if s.Values != nil {
	// // 	p.print(vtab, token.ASSIGN, blank)
	// // 	p.exprList(token.NoPos, s.Values, 1, 0, token.NoPos, false)
	// // 	extraTabs--
	// // }
	// // if s.Comment != nil {
	// // 	for ; extraTabs > 0; extraTabs-- {
	// // 		p.print(vtab)
	// // 	}
	// // 	p.setComment(s.Comment)
	// // }
}

func (p *printer) ident(list ...*ast.Ident) {
	p.identList(list)
}

func (p *printer) identList(list []*ast.Ident) {
	// convert into an expression list so we can re-use exprList formatting
	xlist := make([]ast.Expr, len(list))
	for i, x := range list {
		xlist[i] = x
	}
	// var mode exprListMode
	// if !indent {
	// 	mode = noIndent
	// }
	p.exprList(xlist)
}

func (p *printer) print(args ...any) {
	for _, a := range args {
		switch v := a.(type) {
		case string:
			fmt.Fprint(&p.buf, v)
		case *ast.Ident:
			fmt.Fprint(&p.buf, v.Name)
		case *ast.BasicLit:
			fmt.Fprint(&p.buf, v.Value)
		case token.Token:
			p.token(v)
		default:
			panic(fmt.Sprintf("print: Unknown Type: %T", a))
		}
	}
}

func (p *printer) token(t token.Token) {
	switch t {
	case token.ADD: fallthrough
	case token.SUB: fallthrough
	case token.QUO: fallthrough
	case token.REM: fallthrough
	case token.AND: fallthrough
	case token.OR: fallthrough
	case token.XOR: fallthrough
	case token.SHL: fallthrough
	case token.SHR: fallthrough
	case token.AND_NOT: fallthrough
	case token.MUL: fallthrough
	case token.NOT:
		p.print(t.String())
	// case token.TYPE:
	// 	p.print("TYPE")
	// case token.IMPORT:
	// 	p.print(IMPORT)
	default:
		panic(fmt.Sprintf("token: Unknown Token: %s", t))
	}
	// fmt.Fprint(&p.buf, v.String())
}

//--------------------------------------------------------------------------------

func (p *printer) exprList(list []ast.Expr) {
	if len(list) == 0 {
		return
	}

	for i, xx := range list {
		switch x := xx.(type) {
		case *ast.KeyValueExpr:
			// fmt.Println("ast.KeyValueExpr", x)
			p.print(".", x.Key, " = ", x.Value)
		default:
			p.expr(xx)
			if i < len(list)-1 {
				p.print(TUPLESEP)
			}
		}
	}
}

func (p *printer) selectorExpr(x *ast.SelectorExpr) {
	p.expr(x.X)
	p.print(".")
	p.print(x.Sel)
}

//--------------------------------------------------------------------------------

func (p *printer) expr(expr ast.Expr) {
	switch x := expr.(type) {
	case *ast.BadExpr:
		p.print("BadExpr")

	case *ast.Ident:
		p.print(x)

	case *ast.BinaryExpr:
		// if depth < 1 {
		// 	p.internalError("depth < 1:", depth)
		// 	depth = 1
		// }
		// p.binaryExpr(x, prec1, cutoff(x, depth), depth)
		p.expr(x.X)
		p.print(SPACE, x.Op, SPACE)
		p.expr(x.Y)

	// case *ast.KeyValueExpr:
	// 	p.expr(x.Key)
	// 	p.setPos(x.Colon)
	// 	p.print(token.COLON, blank)
	// 	p.expr(x.Value)

	// case *ast.StarExpr:
	// 	const prec = token.UnaryPrec
	// 	if prec < prec1 {
	// 		// parenthesis needed
	// 		p.print(token.LPAREN)
	// 		p.print(token.MUL)
	// 		p.expr(x.X)
	// 		p.print(token.RPAREN)
	// 	} else {
	// 		// no parenthesis needed
	// 		p.print(token.MUL)
	// 		p.expr(x.X)
	// 	}

	case *ast.UnaryExpr:
		p.token(x.Op)
		p.expr(x.X)
		// const prec = token.UnaryPrec
		// if prec < prec1 {
		// 	// parenthesis needed
		// 	p.print(token.LPAREN)
		// 	p.expr(x)
		// 	p.print(token.RPAREN)
		// } else {
		// 	// no parenthesis needed
		// 	p.print(x.Op)
		// 	if x.Op == token.RANGE {
		// 		// TODO(gri) Remove this code if it cannot be reached.
		// 		p.print(blank)
		// 	}
		// 	p.expr1(x.X, prec, depth)
		// }

	case *ast.BasicLit:
		// if p.Config.Mode&normalizeNumbers != 0 {
		// 	x = normalizedNumber(x)
		// }
		p.print(x)

	// case *ast.FuncLit:
	// 	p.setPos(x.Type.Pos())
	// 	p.print(token.FUNC)
	// 	// See the comment in funcDecl about how the header size is computed.
	// 	startCol := p.out.Column - len("func")
	// 	p.signature(x.Type)
	// 	p.funcBody(p.distanceFrom(x.Type.Pos(), startCol), blank, x.Body)

	// case *ast.ParenExpr:
	// 	if _, hasParens := x.X.(*ast.ParenExpr); hasParens {
	// 		// don't print parentheses around an already parenthesized expression
	// 		// TODO(gri) consider making this more general and incorporate precedence levels
	// 		p.expr0(x.X, depth)
	// 	} else {
	// 		p.print(token.LPAREN)
	// 		p.expr0(x.X, reduceDepth(depth)) // parentheses undo one level of depth
	// 		p.setPos(x.Rparen)
	// 		p.print(token.RPAREN)
	// 	}

	case *ast.SelectorExpr:
		p.selectorExpr(x)

	// case *ast.TypeAssertExpr:
	// 	p.expr1(x.X, token.HighestPrec, depth)
	// 	p.print(token.PERIOD)
	// 	p.setPos(x.Lparen)
	// 	p.print(token.LPAREN)
	// 	if x.Type != nil {
	// 		p.expr(x.Type)
	// 	} else {
	// 		p.print(token.TYPE)
	// 	}
	// 	p.setPos(x.Rparen)
	// 	p.print(token.RPAREN)

	// case *ast.IndexExpr:
	// 	// TODO(gri): should treat[] like parentheses and undo one level of depth
	// 	p.expr1(x.X, token.HighestPrec, 1)
	// 	p.setPos(x.Lbrack)
	// 	p.print(token.LBRACK)
	// 	p.expr0(x.Index, depth+1)
	// 	p.setPos(x.Rbrack)
	// 	p.print(token.RBRACK)

	// case *ast.IndexListExpr:
	// 	// TODO(gri): as for IndexExpr, should treat [] like parentheses and undo
	// 	// one level of depth
	// 	p.expr1(x.X, token.HighestPrec, 1)
	// 	p.setPos(x.Lbrack)
	// 	p.print(token.LBRACK)
	// 	p.exprList(x.Lbrack, x.Indices, depth+1, commaTerm, x.Rbrack, false)
	// 	p.setPos(x.Rbrack)
	// 	p.print(token.RBRACK)

	// case *ast.SliceExpr:
	// 	// TODO(gri): should treat[] like parentheses and undo one level of depth
	// 	p.expr1(x.X, token.HighestPrec, 1)
	// 	p.setPos(x.Lbrack)
	// 	p.print(token.LBRACK)
	// 	indices := []ast.Expr{x.Low, x.High}
	// 	if x.Max != nil {
	// 		indices = append(indices, x.Max)
	// 	}
	// 	// determine if we need extra blanks around ':'
	// 	var needsBlanks bool
	// 	if depth <= 1 {
	// 		var indexCount int
	// 		var hasBinaries bool
	// 		for _, x := range indices {
	// 			if x != nil {
	// 				indexCount++
	// 				if isBinary(x) {
	// 					hasBinaries = true
	// 				}
	// 			}
	// 		}
	// 		if indexCount > 1 && hasBinaries {
	// 			needsBlanks = true
	// 		}
	// 	}
	// 	for i, x := range indices {
	// 		if i > 0 {
	// 			if indices[i-1] != nil && needsBlanks {
	// 				p.print(blank)
	// 			}
	// 			p.print(token.COLON)
	// 			if x != nil && needsBlanks {
	// 				p.print(blank)
	// 			}
	// 		}
	// 		if x != nil {
	// 			p.expr0(x, depth+1)
	// 		}
	// 	}
	// 	p.setPos(x.Rbrack)
	// 	p.print(token.RBRACK)

	case *ast.CallExpr:
		// if len(x.Args) > 1 {
		// 	depth++
		// }

		// Conversions to literal function types or <-chan
		// types require parentheses around the type.
		paren := false
		switch t := x.Fun.(type) {
		case *ast.FuncType:
			paren = true
		// case *ast.ChanType:
		// 	paren = t.Dir == ast.RECV
		case *ast.SelectorExpr:
			p.selectorExpr(t)
		case *ast.Ident:
			p.ident(t)
		default:
			panic(fmt.Sprintf("expr: ast.CallExpr unsupported type: %T", t))
		}
		if paren {
			p.print(LPAREN)
		}

		// wasIndented := p.possibleSelectorExpr(x.Fun, token.HighestPrec, depth)
		if paren {
			p.print(RPAREN)
		}

		p.print(LPAREN)
		if x.Ellipsis.IsValid() {
			p.exprList(x.Args)
			p.print(ELLIPSIS)
			if x.Rparen.IsValid() { // && p.lineFor(x.Ellipsis) < p.lineFor(x.Rparen) {
				p.print(TUPLESEP)
			}
		} else {
			p.exprList(x.Args)
		}
		p.print(RPAREN)

	case *ast.CompositeLit:
		// composite literal elements that are composite literals themselves may have the type omitted
		if x.Type != nil {
			_, ok := x.Type.(*ast.StructType)
			if ok {
				p.print(".")
			} else {
				p.expr(x.Type)
			}
		}
		fmt.Printf("AAAA: %+v\n", x.Elts)
		p.print(LBRACE)
		p.exprList(x.Elts)
		// // do not insert extra line break following a /*-style comment
		// // before the closing '}' as it might break the code if there
		// // is no trailing ','
		// mode := noExtraLinebreak

		// // do not insert extra blank following a /*-style comment
		// // before the closing '}' unless the literal is empty
		// if len(x.Elts) > 0 {
		// 	mode |= noExtraBlank
		// }
		// need the initial indent to print lone comments with
		// the proper level of indentation
		// p.print(indent, unindent, mode)
		// p.setPos(x.Rbrace)
		p.print(RBRACE)
		// p.level--

	// case *ast.Ellipsis:
	// 	p.print(token.ELLIPSIS)
	// 	if x.Elt != nil {
	// 		p.expr(x.Elt)
	// 	}

	// case *ast.ArrayType:
	// 	p.print(token.LBRACK)
	// 	if x.Len != nil {
	// 		p.expr(x.Len)
	// 	}
	// 	p.print(token.RBRACK)
	// 	p.expr(x.Elt)

	case *ast.StructType:
		p.print(STRUCT)
		p.fieldList(x.Fields)

	case *ast.FuncType:
		p.print(FUNC)
		p.signature(x)

	case *ast.InterfaceType:
		p.print(INTERFACE)
		p.fieldList(x.Methods)

	// case *ast.MapType:
	// 	p.print(token.MAP, token.LBRACK)
	// 	p.expr(x.Key)
	// 	p.print(token.RBRACK)
	// 	p.expr(x.Value)

	// case *ast.ChanType:
	// 	switch x.Dir {
	// 	case ast.SEND | ast.RECV:
	// 		p.print(token.CHAN)
	// 	case ast.RECV:
	// 		p.print(token.ARROW, token.CHAN) // x.Arrow and x.Pos() are the same
	// 	case ast.SEND:
	// 		p.print(token.CHAN)
	// 		p.setPos(x.Arrow)
	// 		p.print(token.ARROW)
	// 	}
	// 	p.print(blank)
	// 	p.expr(x.Value)

	default:
		panic(fmt.Sprintf("expr: Missing Type: %T", expr))
		// panic("unreachable")
	}
}

//--------------------------------------------------------------------------------

// 1. use shim layer for fmt
// 2. pub fn main()

var stdShims = map[string]string{
	"fmt": "@import(\"lib/fmt.zig\")",
	"github.com/raysan5/raylib": `@cImport({
    @cInclude("raylib.h");
    @cInclude("raymath.h");
    @cInclude("rlgl.h");
})`,
}

func getPackagePath(p string) string {
	str := strings.Trim(p, "\"")
	pkg, ok := stdShims[str]
	if ok {
		return pkg
	}
	return fmt.Sprintf("@import(%s)", p)
}

func nameFromPath(p string) string {
	str := strings.Trim(p, "\"")

	return path.Base(str)
}

func (p *printer) spec(spec ast.Spec) {
	switch s := spec.(type) {
	case *ast.ImportSpec:
		pkgPath := getPackagePath(s.Path.Value)
		if s.Name != nil {
			p.print(CONST, SPACE)
			p.expr(s.Name)
			p.print(ASSIGN, SPACE)
			// p.expr(s.Path) // TODO: Sanitize?
			p.print(pkgPath)
			p.print(SEMICOLON)
		} else {
			p.print(CONST, SPACE)
			name := nameFromPath(s.Path.Value)
			p.print(name, SPACE)
			p.print(ASSIGN, SPACE)
			// p.expr(s.Path) // TODO: Sanitize?
			p.print(pkgPath)
			p.print(SEMICOLON)
		}

		// if s.Name != nil {
		// 	p.expr(s.Name)
		// }

		// p.expr(s.Path) // TODO: Sanitize?

		// p.setComment(s.Doc)
		// if s.Name != nil {
		// 	p.expr(s.Name)
		// 	p.print(blank)
		// }
		// p.expr(sanitizeImportPath(s.Path))
		// p.setComment(s.Comment)
		// p.setPos(s.EndPos)

	case *ast.ValueSpec:
		p.valueSpec(s)
		// // if n != 1 {
		// // 	p.internalError("expected n = 1; got", n)
		// // }
		// // p.setComment(s.Doc)
		// p.identList(s.Names) // always present
		// if s.Type != nil {
		// 	p.print(SPACE)
		// 	p.expr(s.Type)
		// }
		// if s.Values != nil {
		// 	p.print(SPACE, token.ASSIGN, SPACE)
		// 	p.exprList(s.Values)
		// }
		// // p.setComment(s.Comment)

	case *ast.TypeSpec:
		p.typeSpec(s)
		// fmt.Println("TYPESPEC:", s.Name.Name)
		// fmt.Printf("%T: %+v\n", s.Type, s.Type)
		// p.expr(s.Name)
		// p.print(SPACE)
		// if s.TypeParams != nil {
		// 	p.parameters(s.TypeParams)
		// }
		// if s.Assign.IsValid() {
		// 	p.print(ASSIGN, SPACE)
		// }
		// p.expr(s.Type)

		// p.setComment(s.Doc)
		// p.expr(s.Name)
		// if s.TypeParams != nil {
		// 	p.parameters(s.TypeParams, typeTParam)
		// }
		// if n == 1 {
		// 	p.print(blank)
		// } else {
		// 	p.print(vtab)
		// }
		// if s.Assign.IsValid() {
		// 	p.print(token.ASSIGN, blank)
		// }
		// p.expr(s.Type)
		// p.setComment(s.Comment)

	default:
		// panic("unreachable")
		panic(fmt.Sprintf("spec: Missing Type: %T", spec))
	}
}

func (p *printer) typeSpec(s *ast.TypeSpec) {
	// fmt.Println("TYPESPEC:", s.Name.Name)
	// fmt.Printf("%T: %+v\n", s.Type, s.Type)

	// const mystruct = struct
	p.print(CONST, SPACE)
	p.expr(s.Name)
	p.print(SPACE, ASSIGN, SPACE)
	if s.TypeParams != nil {
		p.parameters(s.TypeParams)
	}
	if s.Assign.IsValid() {
		p.print(ASSIGN, SPACE)
	}
	p.expr(s.Type)
}

func (p *printer) forStmt(s *ast.ForStmt) {
	isWhile := s.Init == nil && s.Post == nil

	if isWhile {
		p.print("while ", LPAREN)
		if s.Cond != nil {
			p.expr(s.Cond)
		}
		p.print(RPAREN)
	} else {
		if s.Init != nil {
			p.stmt(s.Init)
		}
		if s.Cond != nil {
			p.expr(s.Cond)
		}
		if s.Post != nil {
			p.stmt(s.Post)
		}
	}

	p.funcBody(s.Body)
}
