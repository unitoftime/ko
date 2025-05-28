package main

import (
	"fmt"
)

// Adapted from: https://github.com/aaronraff/blog-code/blob/master/how-to-write-a-lexer-in-go/lexer.go

func printErr(tok Token, msg string) {
	fmt.Printf("./%s:%d:%d: %s\n", tok.pos.filename, tok.pos.line, tok.pos.column, msg)
}

func parseError(expected TokenType, got Token) error {
	return fmt.Errorf("Expected: %s, Got: (%s) %+v", expected, got.token, got)
}

type Node interface {
	// WalkGraphviz(string, *bytes.Buffer)
	Pos() Position
}

type FileNode struct {
	filename string
	nodes []Node
}
func (n *FileNode) Pos() Position {
	return Position{}
}

type FuncNode struct {
	pos Position
	name string
	arguments *ArgNode
	returns *ArgNode
	body Node
}
func (n *FuncNode) Pos() Position {
	return Position{}
}
type StructNode struct {
	global bool
	ident Token
	fields []*Arg
}
func (n *StructNode) Pos() Position {
	return Position{}
}

type CurlyScope struct {
	nodes []Node
}
func (n *CurlyScope) Pos() Position {
	return Position{}
}

type PackageNode struct {
	name string
}
func (n *PackageNode) Pos() Position {
	return Position{}
}

type CommentNode struct {
	line string
}
func (n *CommentNode) Pos() Position {
	return Position{}
}

type ReturnNode struct {
	expr Node
}
func (n *ReturnNode) Pos() Position {
	return Position{}
}

type Arg struct {
	name Token
	kind Token
}
func (n *Arg) Pos() Position {
	return Position{}
}
func (a *Arg) Type() Type {
	return Type(a.kind.str)
}

type ArgNode struct {
	args []*Arg
}
func (n *ArgNode) Pos() Position {
	return Position{}
}

type VarStmt struct {
	name Token
	global bool
	initExpr Node
	ty Type
}
func (n *VarStmt) Pos() Position {
	return Position{}
}
type IfStmt struct {
	cond Node
	thenScope Node
	elseScope Node
}
func (n *IfStmt) Pos() Position {
	return Position{}
}

type ForStmt struct {
	tok Token
	init, cond, inc Node
	body Node
}
func (n *ForStmt) Pos() Position {
	return Position{}
}


type Stmt struct {
	node Node
}
func (n *Stmt) Pos() Position {
	return Position{}
}

type CallExpr struct {
	callee Node
	rparen Token // Just for position data I guess?
	args []Node
}
func (n *CallExpr) Pos() Position {
	return Position{}
}

type CompLitExpr struct {
	callee Node
	args []Node
	ty Type
}
func (n *CompLitExpr) Pos() Position {
	return Position{}
}

type GetExpr struct {
	obj Node
	name Token
}
func (n *GetExpr) Pos() Position {
	return Position{}
}

type SetExpr struct {
	obj Node
	name Token
	value Node
}
func (n *SetExpr) Pos() Position {
	return Position{}
}

type LogicalExpr struct {
	left Node
	op Token
	right Node
}
func (n *LogicalExpr) Pos() Position {
	return Position{}
}

type AssignExpr struct {
	name Token
	value Node
}
func (n *AssignExpr) Pos() Position {
	return Position{}
}

type BinaryExpr struct {
	left, right Node
	op Token
	ty Type
}
func (n *BinaryExpr) Pos() Position {
	return Position{}
}

type UnaryExpr struct {
	right Node
	op Token
}
func (n *UnaryExpr) Pos() Position {
	return Position{}
}

type LitExpr struct {
	tok Token
	kind TokenType
}
func (n *LitExpr) Pos() Position {
	return Position{}
}
type IdentExpr struct {
	tok Token
	ty Type
}
func (n *IdentExpr) Pos() Position {
	return Position{}
}

type GroupingExpr struct {
	Node
}
func (n *GroupingExpr) Pos() Position {
	return Position{}
}

// --------------------------------------------------------------------------------
// - Parser
// --------------------------------------------------------------------------------
type Tokens struct {
	list []Token
	prev Token
}
func (t *Tokens) Len() int {
	return len(t.list)
}
func (t *Tokens) Peek() Token {
	return t.list[0]
}
func (t *Tokens) Next() Token {
	t.prev = t.list[len(t.list)-1]
	token := t.list[0]
	t.list = t.list[1:]
	return token
}

func (t *Tokens) Prev() Token {
	return t.prev
}

func (t *Tokens) Match(tokType TokenType) bool {
	if t.Peek().token == tokType {
		t.Next()
		return true
	}
	return false
}
func (t *Tokens) Consume(tokType TokenType) Token {
	if t.Peek().token == tokType {
		return t.Next()
	}
	printErr(t.Peek(), "invalid consume")
	panic(parseError(tokType, t.Peek()))
}


type ParseResult struct {
	file *FileNode

	// Globally scoped things
	typeList []Node
	fnList []*FuncNode
	varList []*VarStmt
}

func (p *Parser) Parse(name string) ParseResult {
	file := p.ParseFile(name)
	return ParseResult{
		file: file,
		typeList: p.typeList,
		fnList: p.fnList,
		varList: p.varList,
	}
}

func (p *Parser) ParseFile(name string) *FileNode {
	return &FileNode{
		name,
		p.ParseTil(EOF, true),
	}
}

type Parser struct {
	tokens *Tokens
	blockCompLit bool // If true, we are parsing something like an if or a for with {...} somewhere, so we cant allow composit lits, here
	typeList []Node // A list of every registered type
	fnList []*FuncNode // A list of every registered function
	varList []*VarStmt // list of every global variable
}
func NewParser(tokens *Tokens) *Parser {
	return &Parser{
		tokens: tokens,
		typeList: make([]Node, 0, 32),
		fnList: make([]*FuncNode, 0, 32),
	}
}
func (p *Parser) PrintNext() {
	fmt.Println(p.tokens.Peek())
}

func (p *Parser) Match(tokType TokenType) bool {
	return p.tokens.Match(tokType)
}
func (p *Parser) Consume(tokType TokenType) Token {
	return p.tokens.Consume(tokType)
}

func (p *Parser) ParseTil(stopToken TokenType, globalScope bool) []Node {
	tokens := p.tokens
	nodes := make([]Node, 0)
	for tokens.Len() > 0 {

		if tokens.Peek().token == stopToken {
			tokens.Next()
			return nodes
		}

		node := p.ParseDecl(globalScope)
		if node != nil {
			nodes = append(nodes, node)
		} else {
			next := tokens.Next()
			if next.token != SEMI {
				panic(parseError(stopToken, next))
			} else {
				fmt.Println("-------------------------semi")
			}
		}
	}
	return nodes
}

func (p *Parser) ParseDecl(globalScope bool) Node {
	tokens := p.tokens
	next := tokens.Peek()

	switch next.token {
	case PACKAGE:
		tokens.Next()
		next := tokens.Next()
		tokens.Consume(SEMI)

		return &PackageNode{
			name: next.str,
		}
	case TYPE:
		return p.TypeNode(globalScope)
	case FUNC:
		return p.ParseFuncNode(globalScope)

	case RETURN:
		tokens.Next()
		return p.ParseReturnNode(tokens)
	// case TYPE: // TODO:
	case VAR:
		tokens.Consume(VAR)
		return p.varDecl(globalScope)
	case IF:
		tokens.Consume(IF)
		return p.ifStatement(tokens)
	case FOR:
		return p.forStatement()
	case IDENT:
		return p.parseStatement()

	case LINECOMMENT:
		next := tokens.Next() // Discard
		return &CommentNode{next.str}
	case SEMI:
		// tokens.Next() // Discard
		return nil
	default:
		printErr(next, fmt.Sprintf("unexpected type: %s", next.str))
		panic(fmt.Sprintf("Unknown Type: %s", next.str))
	}

	// if next.str == "func" {
	// 	tokens.Next()
	// 	return p.ParseFuncNode(tokens)
	// } else if next.str == "package" {
	// 	next := tokens.Next()
	// 	return PackageNode{
	// 		name: next.str,
	// 	}
	// } else if next.str == "return" {
	// 	tokens.Next()
	// 	return p.ParseReturnNode(tokens)
	// }

	// return nil // TODO fix
}

// Parsing functions


func (p *Parser) TypeNode(globalScope bool) Node {
	typeTok := p.Consume(TYPE)

	ident := p.Consume(IDENT)

	if p.Match(STRUCT) {
		p.Consume(LBRACE)
		fields := make([]*Arg, 0)
		for {
			if p.Match(RBRACE) {
				break
			}

			field := p.Consume(IDENT)
			kind := p.Consume(IDENT)
			p.Consume(SEMI)
			fields = append(fields, &Arg{field, kind})
		}

		s := &StructNode{
			global: globalScope,
			ident: ident,
			fields: fields,
		}

		if globalScope {
			p.typeList = append(p.typeList, s)
		}

		return s
	}

	printErr(typeTok, "invalid typedef")
	panic("Invalid typedef")
}

func (p *Parser) ParseFuncNode(globalScope bool) Node {
	tokens := p.tokens
	funcToken := tokens.Consume(FUNC)

	next := tokens.Next()
	if next.token != IDENT {
		panic("MUST BE IDENTIFIER")
	}

	args := p.ParseArgNode()

	// Try to parse return types
	var returns *ArgNode
	{
		next := tokens.Peek()
		switch next.token {
		case LPAREN:
			returns = p.ParseArgNode()
		case IDENT:
			tokens.Next()
			args := ArgNode{make([]*Arg, 0)}
			args.args = append(args.args, &Arg{
				name: Token{}, // TODO
				kind: next,
			})
			returns = &args
		}
	}

	body := p.ParseCurlyScope(tokens)
	f := FuncNode{
		pos: funcToken.pos,
		name: next.str,
		arguments: args,
		returns: returns,
		body: body,
	}

	if globalScope {
		p.fnList = append(p.fnList, &f)
	}

	return &f
}

func (p *Parser) ParseCurlyScope(tokens *Tokens) Node {
	next := tokens.Next()
	if next.token != LBRACE {
		panic(parseError(LBRACE, next))
	}

	body := p.ParseTil(RBRACE, false)

	return &CurlyScope{body}
}


func (p *Parser) ParseReturnNode(tokens *Tokens) Node {
	r := &ReturnNode{
		expr: p.ParseExpression(),
	}
	return r
}

func (p *Parser) ParseArgNode() *ArgNode {
	tokens := p.tokens
	next := tokens.Next()
	if next.token != LPAREN {
		printErr(next, "expected left parenthesis")
		panic(parseError(LPAREN, next))
	}

	args := &ArgNode{make([]*Arg, 0)}
	for {
		if tokens.Peek().token == RPAREN { break }

		arg := p.ParseTypedArg(tokens)
		args.args = append(args.args, arg)

		if tokens.Peek().token == COMMA {
			tokens.Next()
		}
	}

	tokens.Next() // Drop the RPAREN

	return args
}

func (p *Parser) ParseTypedArg(tokens *Tokens) *Arg {
	name := tokens.Next()
	if name.token != IDENT {
		panic(fmt.Sprintf("MUST BE IDENT: %s", name.str))
	}

	kind := tokens.Next()
	if kind.token != IDENT {
		panic(fmt.Sprintf("MUST BE IDENT: %s", kind.str))
	}

	return &Arg{name, kind}
}

// func (p *Parser) ParseExprNode(tokens *Tokens) Node {
// 	peek := tokens.Peek()
// 	if peek.token == LPAREN {
// 		tokens.Next()
// 		// Case where we have a subexpression
// 		op := p.ParseExprNode(tokens)
// 		// if tokens.Next().token != RPAREN {
// 		// 	panic("SHOULD BE RPAREN!!!!")
// 		// }
// 		return &ExprNode{
// 			ops: []Node{op},
// 		}
// 	}

// 	expr := ExprNode{
// 		ops: make([]Node, 0),
// 	}
// 	// Case where we have a (potentially long) flat expression
// 	idx := -1
// 	for {
// 		idx++
// 		if tokens.Peek().token == RPAREN {
// 			tokens.Next()
// 			break
// 		}
// 		if tokens.Peek().token == SEMI {
// 			tokens.Next()
// 			break
// 		}
// 		if tokens.Peek().token == LPAREN {
// 			expr.ops = append(expr.ops, p.ParseExprNode(tokens))
// 			continue
// 		}

// 		next := tokens.Next()
// 		expr.ops = append(expr.ops, &UnaryNode{idx, next})
// 	}

// 	return &expr
// }

func (p *Parser) parseStatement() Node {
	// switch tokens.Peek().token {
	// case LPAREN:
	// // 	return Stmt{p.ParseFuncCall(tokens)}
	// }

	return &Stmt{p.ParseExpression()}
}

func (p *Parser) varDecl(globalScope bool) *VarStmt {
	tokens := p.tokens
	name := tokens.Next()

	var initExpr Node
	if tokens.Peek().token == EQUAL {
		tokens.Next()
		initExpr = p.ParseExpression()
	}

	tokens.Consume(SEMI)
	stmt := &VarStmt{name, globalScope, initExpr, UnknownType}

	if globalScope {
		p.varList = append(p.varList, stmt)
	}
	return stmt
}

func (p *Parser) ifStatement(tokens *Tokens) Node {
	p.blockCompLit = true
	defer func() { p.blockCompLit = false }()
	cond := p.ParseExpression()

	thenScope := p.ParseCurlyScope(tokens)

	var elseScope Node
	if tokens.Peek().token == ELSE {
		tokens.Next()
		elseScope = p.ParseCurlyScope(tokens)
	}

	return &IfStmt{cond, thenScope, elseScope}
}

func (p *Parser) forStatement() Node {
	p.blockCompLit = true
	defer func() { p.blockCompLit = false }()

	forTok := p.Consume(FOR)

	var init Node
	if (p.tokens.Match(SEMI)) {
		init = nil;
	} else if (p.tokens.Match(VAR)) {
		init = p.varDecl(false)
	} else {
		init = p.ParseExpression()
		p.Consume(SEMI)
	}

	var cond Node
	if p.Match(SEMI) {
		cond = nil
	} else {
		cond = p.ParseExpression()
		p.Consume(SEMI)
	}

	var inc Node
	if p.Match(SEMI) {
		inc = nil
	} else {
		inc = p.ParseExpression()
	}

	body := p.ParseCurlyScope(p.tokens)

	return &ForStmt{forTok, init, cond, inc, body}

}

// func (p *Parser) ParseFuncCall(tokens *Tokens) Node {

// }

func (p *Parser) ParseExpression() Node {
	return p.Assignment(p.tokens)
}

func (p *Parser) Assignment(tokens *Tokens) Node {
	expr := p.Or()

	if tokens.Peek().token == EQUAL {
		tokens.Next()
		value := p.Assignment(tokens)

		varExp, validTarget := expr.(*IdentExpr)
		if validTarget {
			name := varExp.tok
			return &AssignExpr{name, value}
		}
		getExp, validTarget := expr.(*GetExpr)
		if validTarget {
			name := getExp.name
			return &SetExpr{getExp.obj, name, value}
		}


		panic("INVALID ASSIGNMENT TARGET")
	}

	return expr
}

func (p *Parser) CompLit() Node {
	tokens := p.tokens
	expr := p.Or()

	if tokens.Peek().token == EQUAL {
		tokens.Next()
		value := p.CompLit()

		varExp, validTarget := expr.(*IdentExpr)
		if validTarget {
			name := varExp.tok
			return &AssignExpr{name, value}
		}

		panic("INVALID ASSIGNMENT TARGET")
	}

	return expr
}

func (p *Parser) Or() Node {
	expr := p.And()
	for p.tokens.Match(OR) {
		op := p.tokens.Prev()
		right := p.And()

		expr = &LogicalExpr{expr, op, right}
	}
	return expr
}

func (p *Parser) And() Node {
	expr := p.Equality(p.tokens)
	for p.tokens.Match(AND) {
		op := p.tokens.Prev()
		right := p.Equality(p.tokens)
		expr = &LogicalExpr{expr, op, right}
	}
	return expr
}

func (p *Parser) Equality(tokens *Tokens) Node {
	expr := p.ParseExprComparison(tokens)

	for {
		middle := tokens.Peek()
		switch middle.token {
		case BANGEQUAL: fallthrough
		case EQUALEQUAL:
			middle = tokens.Next()
			right := p.ParseExprComparison(tokens)
			expr = &BinaryExpr{expr, right, middle, UnknownType}
		default:
			return expr
		}
	}
}

func (p *Parser) ParseExprComparison(tokens *Tokens) Node {
	expr := p.ParseExprTerm(tokens)

	for {
		middle := tokens.Peek()
		switch middle.token {
		case GREATER: fallthrough
		case GREATEREQUAL: fallthrough
		case LESS: fallthrough
		case LESSEQUAL:
			middle = tokens.Next()
			right := p.ParseExprTerm(tokens)
			expr = &BinaryExpr{expr, right, middle, UnknownType}
		default:
			return expr
		}
	}
}

func (p *Parser) ParseExprTerm(tokens *Tokens) Node {
	expr := p.ParseExprFactor(tokens)

	for {
		middle := tokens.Peek()
		switch middle.token {
		case SUB: fallthrough
		case ADD:
			middle = tokens.Next()
			right := p.ParseExprFactor(tokens)
			expr = &BinaryExpr{expr, right, middle, UnknownType}
		default:
			return expr
		}
	}
}

func (p *Parser) ParseExprFactor(tokens *Tokens) Node {
	expr := p.Unary(tokens)

	for {
		middle := tokens.Peek()
		switch middle.token {
		case DIV: fallthrough
		case MUL:
			middle = tokens.Next()
			right := p.Unary(tokens)
			expr = &BinaryExpr{expr, right, middle, UnknownType}
		default:
			return expr
		}
	}
}

func (p *Parser) Unary(tokens *Tokens) Node {
	op := tokens.Peek()
	switch op.token {
	case BANG: fallthrough
	case SUB:
		op = tokens.Next()
		right := p.Unary(tokens)
		return &UnaryExpr{right, op}
	}

	return p.ParseExprCall()
}

// Handles function calls and struct instantiation calls ie myStruct{...}
func (p *Parser) ParseExprCall() Node {
	expr := p.ParseExprPrimary(p.tokens)

	for {
		if p.Match(LPAREN) {
			expr = p.FinishCall(expr)
		} else if !p.blockCompLit && p.Match(LBRACE) {
			expr = p.FinishCompLit(expr)
		} else if  p.Match(DOT) {
			name := p.Consume(IDENT)
			expr = &GetExpr{expr, name}
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) FinishCall(callee Node) Node {
	if p.Match(RPAREN) {
		return &CallExpr{callee, p.tokens.Prev(), nil}
	}

	args := make([]Node, 0)
	for {
		args = append(args, p.ParseExpression());
		if !p.Match(COMMA) {
			break
		}

		if len(args) > 255 {
			printErr(p.tokens.Peek(), "can't have more than 255 arguments")
		}
	}
	tok := p.Consume(RPAREN)
	return &CallExpr{callee, tok, args}
}

func (p *Parser) FinishCompLit(callee Node) Node {
	if p.Match(RBRACE) {
		return &CompLitExpr{callee, nil, UnknownType}
	}

	args := make([]Node, 0)
	for {
		args = append(args, p.ParseExpression());
		if !p.Match(COMMA) {
			break
		}

		if len(args) > 255 {
			printErr(p.tokens.Peek(), "can't have more than 255 arguments")
		}
	}
	p.Consume(RBRACE)
	return &CompLitExpr{callee, args, UnknownType}
}

func (p *Parser) ParseExprPrimary(tokens *Tokens) Node {
	op := tokens.Peek()
	switch op.token {
	// case FALSE: fallthrough
	// case TRUE: fallthrough
	// case NIL: fallthrough
	case INT:
		tok := tokens.Next()
		return &LitExpr{tok, INT}
	case FLOAT:
		tok := tokens.Next()
		return &LitExpr{tok, FLOAT}
	case IDENT:
		tok := tokens.Next()
		return &IdentExpr{tok, UnknownType}
	case STRING:
		tok := tokens.Next()
		return &LitExpr{tok, STRING}
	case LPAREN:
		expr := &GroupingExpr{p.Equality(tokens)}

		tokens.Consume(RPAREN)
		return expr
	}

	printErr(op, "illegal expr")
	panic(parseError(ILLEGAL, op))
}
