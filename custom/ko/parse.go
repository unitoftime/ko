package main

import (
	"fmt"
)

// Adapted from: https://github.com/aaronraff/blog-code/blob/master/how-to-write-a-lexer-in-go/lexer.go
func printErr(tok Token, msg string) {
	fmt.Printf("./%s:%d:%d: %s\n", tok.pos.filename, tok.pos.line, tok.pos.column, msg)
	panic("AAA")
}


func errUndefinedIdent(n Node, name string) {
 	nodeError(n, fmt.Sprintf("Undefined Identifier: %s", name))
}
func errUndefinedVar(n Node, name string) {
 	nodeError(n, fmt.Sprintf("Undefined Variable: %s", name))
}
func errUndefinedType(n Node, name string) {
 	nodeError(n, fmt.Sprintf("Undefined Type: %s", name))
}
func errIdentMustBeAType(n Node, name string) {
 	nodeError(n, fmt.Sprintf("Identifier must point to a type: %s", name))
}

func nodeError(n Node, msg string) error {
	p := n.Pos()
	fmt.Printf("./%s:%d:%d: %s\n", p.filename, p.line, p.column, msg)
	panic("AAA")
}

func parseError(expected TokenType, got Token) error {
	return fmt.Errorf("Expected: %s, Got: (%s) %+v", expected, got.token, got)
}

type Node interface {
	// WalkGraphviz(string, *bytes.Buffer)
	Pos() Position
	Type() Type
}

type FileNode struct {
	filename string
	nodes []Node
}
func (n *FileNode) Pos() Position {
	return Position{filename: n.filename}
}
func (n *FileNode) Type() Type {
	return UnknownType
}

type FuncNode struct {
	pos Position
	name string
	arguments *ArgNode
	returns *ArgNode
	body Node
	ty Type
}
func (n *FuncNode) Pos() Position {
	return n.pos
}
func (n *FuncNode) Type() Type {
	return n.ty
}

type StructNode struct {
	global bool
	ident Token
	fields []*Arg
	ty Type
}
func (n *StructNode) Pos() Position {
	return Position{}
}
func (n *StructNode) Type() Type {
	return n.ty
}

type ScopeNode struct {
	Scope *CurlyScope
}
func (n *ScopeNode) Pos() Position {
	return n.Scope.Pos()
}
func (n *ScopeNode) Type() Type {
	return UnknownType
}

type CurlyScope struct {
	pos Position
	nodes []Node
}
func (n *CurlyScope) Pos() Position {
	return n.pos
}
func (n *CurlyScope) Type() Type {
	return UnknownType
}

type PackageNode struct {
	pos Position
	name string
}
func (n *PackageNode) Pos() Position {
	return n.pos
}
func (n *PackageNode) Type() Type {
	return UnknownType
}

type CommentNode struct {
	pos Position
	line string
}
func (n *CommentNode) Pos() Position {
	return Position{}
}
func (n *CommentNode) Type() Type {
	return UnknownType
}

type ReturnNode struct {
	pos Position
	expr Node
	ty Type
}
func (n *ReturnNode) Pos() Position {
	return n.pos
}
func (n *ReturnNode) Type() Type {
	return n.ty
}

type Arg struct {
	name Token
	typeNode *TypeNode
	ty Type
}
func (n *Arg) Pos() Position {
	return n.name.pos
}
func (n *Arg) Type() Type {
	return n.ty
}

type ArgNode struct {
	pos Position
	args []*Arg
}
func (n *ArgNode) Pos() Position {
	return n.pos
}
func (n *ArgNode) Type() Type {
	return UnknownType
}

// For wrapping type grammers to differentiate expressions
type TypeNode struct {
	node Node
	ty Type
}
func NewTypeNode(node Node) *TypeNode {
	return &TypeNode{node, UnknownType}
}

func (n *TypeNode) Name() string {
	return n.ty.Name()
}

func (n *TypeNode) Pos() Position {
	return n.node.Pos()
}
func (n *TypeNode) Type() Type {
	return n.ty
}

type VarStmt struct {
	name Token
	global bool
	typeSpec Node
	initExpr Node
	ty Type
}
func (n *VarStmt) Pos() Position {
	return n.name.pos
}
func (n *VarStmt) Type() Type {
	return n.ty
}

// += -= ... others: +=  -=  *=  /=  %=  &=  |=  ^=  <<=  >>=  &^=
type ShortStmt struct {
	target Node
	op Token
	initExpr Node
}
func (n *ShortStmt) Pos() Position {
	return n.target.Pos()
}
func (n *ShortStmt) Type() Type {
	return UnknownType
}

type IfStmt struct {
	cond Node
	thenScope Node
	elseScope Node
}
func (n *IfStmt) Pos() Position {
	return Position{}
}
func (n *IfStmt) Type() Type {
	return UnknownType
}

type ForStmt struct {
	tok Token
	init, cond, inc Node
	body Node
}
func (n *ForStmt) Pos() Position {
	return Position{}
}
func (n *ForStmt) Type() Type {
	return UnknownType
}

type Stmt struct {
	node Node
}
func (n *Stmt) Pos() Position {
	return n.node.Pos()
}
func (n *Stmt) Type() Type {
	return UnknownType
	// return n.ty
}

type CallExpr struct {
	callee Node
	rparen Token // Just for position data I guess?
	args []Node
	ty Type // Note: This is the type returned by the call
}
func (n *CallExpr) Pos() Position {
	return n.callee.Pos()
}
func (n *CallExpr) Type() Type {
	return n.ty
}

type CompLitExpr struct {
	callee Node
	args []Node
	ty Type
}
func (n *CompLitExpr) Pos() Position {
	return n.callee.Pos()
}
func (n *CompLitExpr) Type() Type {
	return UnknownType
}

type GetExpr struct {
	obj Node
	name Token
	ty Type
}
func (n *GetExpr) Pos() Position {
	return n.obj.Pos()
}
func (n *GetExpr) Type() Type {
	return n.ty
}

type SetExpr struct {
	obj Node
	name Token
	value Node
	ty Type
}
func (n *SetExpr) Pos() Position {
	return n.obj.Pos()
}
func (n *SetExpr) Type() Type {
	return n.ty
}

type LogicalExpr struct {
	left Node
	op Token
	right Node
}
func (n *LogicalExpr) Pos() Position {
	return n.left.Pos()
}
func (n *LogicalExpr) Type() Type {
	return BoolType // TODO: Is this always the case?
}

type AssignExpr struct {
	name Token
	value Node
}
func (n *AssignExpr) Pos() Position {
	return n.name.pos
}
func (n *AssignExpr) Type() Type {
	return UnknownType
}

type BinaryExpr struct {
	left, right Node
	op Token
	ty Type
}
func (n *BinaryExpr) Pos() Position {
	return n.op.pos
}
func (n *BinaryExpr) Type() Type {
	return n.ty
}

type PostfixStmt struct {
	left Node
	op Token
	// ty Type
}
func (n *PostfixStmt) Pos() Position {
	return n.op.pos
}
func (n *PostfixStmt) Type() Type {
	return UnknownType
}

// Note: This is now more of a prefix?
type UnaryExpr struct {
	right Node
	op Token
	ty Type
}
func (n *UnaryExpr) Pos() Position {
	return n.op.pos
}
func (n *UnaryExpr) Type() Type {
	return n.ty
}

type LitExpr struct {
	tok Token
	kind TokenType
	ty Type
}
func (n *LitExpr) Pos() Position {
	return n.tok.pos
}
func (n *LitExpr) Type() Type {
	return n.ty
}

type IdentExpr struct {
	tok Token
	ty Type
}
func (n *IdentExpr) Pos() Position {
	return n.tok.pos
}
func (n *IdentExpr) Type() Type {
	return n.ty
}

type GroupingExpr struct {
	Node
	ty Type
}
func (n *GroupingExpr) Pos() Position {
	return n.Node.Pos()
}
func (n *GroupingExpr) Type() Type {
	return n.ty
}

type BuiltinNode struct {
	ty Type
}
func (n *BuiltinNode) Pos() Position {
	return Position{}
}
func (n *BuiltinNode) Type() Type {
	return n.ty
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
func (t *Tokens) PeekNext() Token {
	return t.list[1]
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
	printErr(t.Peek(), fmt.Sprintf("invalid consume: Expected: %s got: %s", tokType, t.Peek().str))
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

func (p *Parser) Next() Token {
	return p.tokens.Next()
}

func (p *Parser) Peek() Token {
	return p.tokens.Peek()
}
func (p *Parser) PeekNext() Token {
	return p.tokens.PeekNext()
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
				Println("-------------------------semi", p.tokens.Prev())
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
			pos: next.pos,
			name: next.str,
		}
	case TYPE:
		return p.TypeDeclNode(globalScope)
	case FUNC:
		return p.ParseFuncNode(globalScope)

	case RETURN:
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
		if globalScope {
			printErr(next, fmt.Sprintf("Unexpected identifier in global scope: %s", next.str))
		}
		return p.parseStatement()
	case LBRACE:
		return &ScopeNode{p.ParseCurlyScope()}
	case LINECOMMENT:
		next := tokens.Next() // Discard
		return &CommentNode{next.pos, next.str}
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


func (p *Parser) TypeDeclNode(globalScope bool) Node {
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
			// kind := p.Consume(IDENT)
			typeNode := p.ParseTypeNode()
			p.Consume(SEMI)
			fields = append(fields, &Arg{field, typeNode, UnknownType})
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
			// If (..) then parse the full arg node
			returns = p.ParseArgNode()
		case LBRACE:
			// If we see { then there is no return type
		default:
			// Else just parse a single typeNode
			typeNode := p.ParseTypeNode()
			fmt.Println("TYPENODE:", typeNode.node)
			returns = &ArgNode{next.pos, []*Arg{
				{
					name: Token{}, // TODO: Unnamed, nil?
					typeNode: typeNode,
				},
			}}
		}
	}


	body := p.ParseCurlyScope()
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

func (p *Parser) ParseCurlyScope() *CurlyScope {
	tokens := p.tokens
	next := tokens.Next()
	if next.token != LBRACE {
		panic(parseError(LBRACE, next))
	}

	body := p.ParseTil(RBRACE, false)

	return &CurlyScope{next.pos, body}
}


func (p *Parser) ParseReturnNode(tokens *Tokens) Node {
	tok := tokens.Next()
	r := &ReturnNode{
		pos: tok.pos,
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

	args := &ArgNode{next.pos, make([]*Arg, 0)}
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

	// kind := tokens.Next()
	// if kind.token != IDENT {
	// 	panic(fmt.Sprintf("MUST BE IDENT: %s", kind.str))
	// }
	typeNode := p.ParseTypeNode()

	return &Arg{name, typeNode, UnknownType}
}

// func (p *Parser) ParseTypeNode() *TypeNode {
// 	// Loop until we find the identifier, all of the prefix operators help define the type
// 	// TODO if you do generics like myType[int] then you also need to check postfix
// 	tok := p.Next()
// 	switch tok.token {
// 	// TODO: [] Slices, Arrays
// 	case MUL: // Pointer
// 		p.ParseTypeNode()
// 	case IDENT:
// 		return &TypeNode{tok}
// 	default:
// 		parseError(IDENT, tok)
// 	}
// }

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
	switch p.Peek().token {
	case IDENT:
		switch p.PeekNext().token {
		case WALRUS:
			return p.varDecl(false)
		}
	}

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
	var typeSpec Node
	if tokens.Peek().token == EQUAL {
		tokens.Next()
		initExpr = p.ParseExpression()
	} else if tokens.Peek().token == WALRUS {
		tokens.Next()
		initExpr = p.ParseExpression()
	} else {
		// Parse the type
		typeSpec = p.ParseTypeNode()
	}

	fmt.Println(initExpr)

	p.Consume(SEMI)
	stmt := &VarStmt{name, globalScope, typeSpec, initExpr, UnknownType}

	if globalScope {
		p.varList = append(p.varList, stmt)
	}
	return stmt
}

// func (p *Parser) shortStatement() *ShortStmt {
// 	tok := p.Consume(IDENT)
// 	op := p.Next()
// 	target := &IdentExpr{tok, UnknownType} // TODO: What could this be? All sorts of assignment types
// 	initExpr := p.ParseExpression()
// 	return &ShortStmt{target, op, initExpr}
// }


func (p *Parser) ifStatement(tokens *Tokens) Node {
	p.blockCompLit = true
	defer func() { p.blockCompLit = false }()
	cond := p.parseStatement()

	thenScope := p.ParseCurlyScope()

	var elseScope Node
	if tokens.Peek().token == ELSE {
		tokens.Next()
		elseScope = p.ParseCurlyScope()
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
		init = p.parseStatement()
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
		inc = p.parseStatement()
	}

	body := p.ParseCurlyScope()

	return &ForStmt{forTok, init, cond, inc, body}

}

func (p *Parser) ParseTypeNode() *TypeNode {
	p.blockCompLit = true
	defer func() { p.blockCompLit = false }()

	return NewTypeNode(p.ParseExpression())
}

func (p *Parser) ParseExpression() Node {
	return p.Assignment(p.tokens)
}

func (p *Parser) Assignment(tokens *Tokens) Node {
	expr := p.Or()

	switch tokens.Peek().token {
	case EQUAL:
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
			return &SetExpr{getExp.obj, name, value, UnknownType}
		}
		panic(fmt.Sprintf("INVALID ASSIGNMENT TARGET: %+v", tokens.Peek().token))

	case INC: fallthrough
	case DEC:
		// TODO: expr must be assignable
		op := tokens.Next()
		return &PostfixStmt{expr, op}

		// Decided to handle this above because the lhs always has to be an ident?
	// case WALRUS:
	// 	op := tokens.Next()
	// 	return &VarStmt{name, globalScope, initExpr, UnknownType}

	case PLUSEQ: fallthrough
	case SUBEQ:
		// TODO: expr must be assignable
		op := tokens.Next()
		initExpr := p.Or()
		return &ShortStmt{expr, op, initExpr}
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
	for p.tokens.Match(ANDAND) {
		op := p.tokens.Prev()
		right := p.Equality(p.tokens)
		expr = &LogicalExpr{expr, op, right}
	}
	return expr
}

func (p *Parser) Equality(tokens *Tokens) Node {
	expr := p.Comparison(tokens)

	for {
		middle := tokens.Peek()
		switch middle.token {
		case BANGEQUAL: fallthrough
		case EQUALEQUAL:
			middle = tokens.Next()
			right := p.Comparison(tokens)
			expr = &BinaryExpr{expr, right, middle, UnknownType}
		default:
			return expr
		}
	}
}

func (p *Parser) Comparison(tokens *Tokens) Node {
	expr := p.Term(tokens)

	for {
		middle := tokens.Peek()
		switch middle.token {
		case GREATER: fallthrough
		case GREATEREQUAL: fallthrough
		case LESS: fallthrough
		case LESSEQUAL:
			middle = tokens.Next()
			right := p.Term(tokens)
			expr = &BinaryExpr{expr, right, middle, UnknownType}
		default:
			return expr
		}
	}
}

func (p *Parser) Term(tokens *Tokens) Node {
	expr := p.Factor(tokens)

	for {
		middle := tokens.Peek()
		switch middle.token {
		case SUB: fallthrough
		case ADD:
			middle = tokens.Next()
			right := p.Factor(tokens)
			expr = &BinaryExpr{expr, right, middle, UnknownType}

		default:
			return expr
		}
	}
}

func (p *Parser) Factor(tokens *Tokens) Node {
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
	case AND: fallthrough
	case MUL: fallthrough
	case BANG: fallthrough
	case SUB:
		op = tokens.Next()
		right := p.Unary(tokens)
		return &UnaryExpr{right, op, UnknownType}
	}

	return p.Call()
}

// Handles function calls and struct instantiation calls ie myStruct{...}
func (p *Parser) Call() Node {
	expr := p.ParseExprPrimary(p.tokens)

	for {
		if p.Match(LPAREN) {
			expr = p.FinishCall(expr)
		} else if !p.blockCompLit && p.Match(LBRACE) {
			expr = p.FinishCompLit(expr)
		} else if  p.Match(DOT) {
			name := p.Consume(IDENT)
			expr = &GetExpr{expr, name, UnknownType}
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) FinishCall(callee Node) Node {
	if p.Match(RPAREN) {
		return &CallExpr{callee, p.tokens.Prev(), nil, UnknownType}
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
	return &CallExpr{callee, tok, args, UnknownType}
}

func (p *Parser) FinishCompLit(callee Node) Node {
	if p.Match(RBRACE) {
		return &CompLitExpr{callee, nil, UnknownType}
	}

	args := make([]Node, 0)
	for {
		args = append(args, p.ParseExpression())
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
	// TODO: NIL literal?

	case TRUE:
		tok := tokens.Next()
		return &LitExpr{tok, TRUE, BoolLitType}
	case FALSE:
		tok := tokens.Next()
		return &LitExpr{tok, FALSE, BoolLitType}

	case INT:
		tok := tokens.Next()
		return &LitExpr{tok, INT, IntLitType}
	case FLOAT:
		tok := tokens.Next()
		return &LitExpr{tok, FLOAT, FloatLitType}
	case IDENT:
		tok := tokens.Next()
		return &IdentExpr{tok, UnknownType}
	case STRING:
		tok := tokens.Next()
		return &LitExpr{tok, STRING, StringLitType}
	case LPAREN:
		tokens.Consume(LPAREN)
		// TODO: Shoudl this be Or?
		expr := &GroupingExpr{p.Equality(p.tokens), UnknownType}

		tokens.Consume(RPAREN)
		return expr
	}

	printErr(op, "illegal expr")
	panic(parseError(ILLEGAL, op))
}
