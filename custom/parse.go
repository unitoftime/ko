package main

import (
	"fmt"
)

// Adapted from: https://github.com/aaronraff/blog-code/blob/master/how-to-write-a-lexer-in-go/lexer.go

func printErr(tok Token, msg string) {
	fmt.Printf("./%s:%d:%d\t%s\n", "test.txt", tok.pos.line, tok.pos.column, msg)
}

func parseError(expected TokenType, got Token) error {
	return fmt.Errorf("Expected: %s, Got: (%s) %+v", expected, got.token, got)
}

type Node interface {
	// WalkGraphviz(string, *bytes.Buffer)
}

type FileNode struct {
	filename string
	nodes []Node
}
type FuncNode struct {
	name string
	arguments Node
	returns Node
	body Node
}
type CurlyScope struct {
	nodes []Node
}

type PackageNode struct {
	name string
}
type CommentNode struct {
	line string
}

type ReturnNode struct {
	expr Node
}

type Arg struct {
	name string
	kind string
}
type ArgNode struct {
	args []Arg
}

type VarStmt struct {
	name Token
	initExpr Node
}
type IfStmt struct {
	cond Node
	thenScope Node
	elseScope Node
}

type ForStmt struct {
	init, cond, inc Node
	body Node
}


type Stmt struct {
	node Node
}

// type Operator uint8
// const (
// 	OpNone Operator = iota
// 	OpAdd
// 	OpSub
// 	OpMul
// 	OpDiv
// )

// type ExprNode struct {
// 	ops []Node
// }

// type UnaryNode struct {
// 	index int
// 	token Token
// }

type CallExpr struct {
	callee Node
	rparen Token // Just for position data I guess?
	args []Node
}


type LogicalExpr struct {
	left Node
	op Token
	right Node
}

type AssignExpr struct {
	name Token
	value Node
}

type BinaryExpr struct {
	left, right Node
	op Token
}

type UnaryExpr struct {
	right Node
	op Token
}

type LitExpr struct {
	tok Token
}
type VarExpr struct {
	tok Token
}

type GroupingExpr struct {
	Node
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
	panic(parseError(tokType, t.Peek()))
}


func (p *Parser) ParseFile(name string, tokens *Tokens) *FileNode {
	return &FileNode{
		name,
		p.ParseTil(tokens, EOF),
	}
}

type Parser struct {
	tokens *Tokens
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

func (p *Parser) ParseTil(tokens *Tokens, stopToken TokenType) []Node {
	nodes := make([]Node, 0)
	for tokens.Len() > 0 {

		if tokens.Peek().token == stopToken {
			tokens.Next()
			return nodes
		}

		node := p.ParseDecl(tokens)
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

func (p *Parser) ParseDecl(tokens *Tokens) Node {
	next := tokens.Peek()

	switch next.token {
	case PACKAGE:
		tokens.Next()
		next := tokens.Next()
		tokens.Consume(SEMI)

		return PackageNode{
			name: next.str,
		}
	case FUNC:
		tokens.Next()
		return p.ParseFuncNode(tokens)

	case RETURN:
		tokens.Next()
		return p.ParseReturnNode(tokens)
	// case TYPE: // TODO:
	case VAR:
		tokens.Consume(VAR)
		return p.varDecl(tokens)
	case IF:
		tokens.Consume(IF)
		return p.ifStatement(tokens)
	case FOR:
		tokens.Consume(FOR)
		return p.forStatement()
	case IDENT:
		return p.parseStatement(tokens)

	case LINECOMMENT:
		next := tokens.Next() // Discard
		return CommentNode{next.str}
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


func (p *Parser) ParseFuncNode(tokens *Tokens) Node {
	next := tokens.Next()
	if next.token != IDENT {
		panic("MUST BE IDENTIFIER")
	}

	args := p.ParseArgNode(tokens)

	// Try to parse return types
	var returns Node
	{
		next := tokens.Peek()
		switch next.token {
		case LPAREN:
			returns = p.ParseArgNode(tokens)
		case IDENT:
			tokens.Next()
			args := ArgNode{make([]Arg, 0)}
			args.args = append(args.args, Arg{
				name: "", // TODO
				kind: next.str,
			})
			returns = args
		}
	}

	body := p.ParseCurlyScope(tokens)
	f := FuncNode{
		name: next.str,
		arguments: args,
		returns: returns,
		body: body,
	}

	return &f
}

func (p *Parser) ParseCurlyScope(tokens *Tokens) Node {
	next := tokens.Next()
	if next.token != LBRACE {
		panic(parseError(LBRACE, next))
	}

	body := p.ParseTil(tokens, RBRACE)

	return &CurlyScope{body}
}


func (p *Parser) ParseReturnNode(tokens *Tokens) Node {
	r := ReturnNode{
		expr: p.ParseExpression(),
	}
	return &r
}

func (p *Parser) ParseArgNode(tokens *Tokens) Node {
	next := tokens.Next()
	if next.token != LPAREN {
		panic(parseError(LPAREN, next))
	}

	args := ArgNode{make([]Arg, 0)}
	for {
		if tokens.Peek().token == RPAREN { break }

		arg := p.ParseTypedArg(tokens)
		args.args = append(args.args, arg)

		if tokens.Peek().token == COMMA {
			tokens.Next()
		}
	}

	tokens.Next() // Drop the RPAREN

	return &args
}

func (p *Parser) ParseTypedArg(tokens *Tokens) Arg {
	name := tokens.Next()
	if name.token != IDENT {
		panic(fmt.Sprintf("MUST BE IDENT: %s", name.str))
	}

	kind := tokens.Next()
	if kind.token != IDENT {
		panic(fmt.Sprintf("MUST BE IDENT: %s", kind.str))
	}

	return Arg{name.str, kind.str}
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

func (p *Parser) parseStatement(tokens *Tokens) Node {
	// switch tokens.Peek().token {
	// case LPAREN:
	// // 	return Stmt{p.ParseFuncCall(tokens)}
	// }

	return Stmt{p.ParseExpression()}
}

func (p *Parser) varDecl(tokens *Tokens) Node {
	name := tokens.Next()

	var initExpr Node
	if tokens.Peek().token == EQUAL {
		tokens.Next()
		initExpr = p.ParseExpression()
	}

	tokens.Consume(SEMI)
	return VarStmt{name, initExpr}
}

func (p *Parser) ifStatement(tokens *Tokens) Node {
	cond := p.ParseExpression()

	thenScope := p.ParseCurlyScope(tokens)

	var elseScope Node
	if tokens.Peek().token == ELSE {
		tokens.Next()
		elseScope = p.ParseCurlyScope(tokens)
	}

	return IfStmt{cond, thenScope, elseScope}
}

func (p *Parser) forStatement() Node {
	p.PrintNext()

	var init Node
	if (p.tokens.Match(SEMI)) {
		init = nil;
	} else if (p.tokens.Match(VAR)) {
		fmt.Println("vardecl")
		init = p.varDecl(p.tokens);
	} else {
		init = p.ParseExpression();
		p.Consume(SEMI)
	}

	p.PrintNext()

	var cond Node
	if p.Match(SEMI) {
		cond = nil
	} else {
		fmt.Println("cond")
		cond = p.ParseExpression()
		p.Consume(SEMI)
	}

	p.PrintNext()

	var inc Node
	if p.Match(SEMI) {
		inc = nil
	} else {
		fmt.Println("inc")
		inc = p.ParseExpression()
	}

	body := p.ParseCurlyScope(p.tokens)

	return ForStmt{init, cond, inc, body}

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

		varExp, validTarget := expr.(VarExpr)
		if validTarget {
			name := varExp.tok
			return AssignExpr{name, value}
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

		expr = LogicalExpr{expr, op, right}
	}
	return expr
}

func (p *Parser) And() Node {
	expr := p.Equality(p.tokens)
	for p.tokens.Match(AND) {
		op := p.tokens.Prev()
		right := p.Equality(p.tokens)
		expr = LogicalExpr{expr, op, right}
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
			expr = BinaryExpr{expr, right, middle}
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
			expr = BinaryExpr{expr, right, middle}
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
			expr = BinaryExpr{expr, right, middle}
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
			expr = BinaryExpr{expr, right, middle}
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
		return UnaryExpr{right, op}
	}

	return p.ParseExprCall()
}

func (p *Parser) ParseExprCall() Node {
	expr := p.ParseExprPrimary(p.tokens)

	for {
		if p.Match(LPAREN) {
			expr = p.FinishCall(expr)
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) FinishCall(callee Node) Node {
	if p.Match(RPAREN) {
		return CallExpr{callee, p.tokens.Prev(), nil}
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
	return CallExpr{callee, tok, args}
}

func (p *Parser) ParseExprPrimary(tokens *Tokens) Node {
	op := tokens.Peek()
	switch op.token {
	// case FALSE: fallthrough
	// case TRUE: fallthrough
	// case NIL: fallthrough
	case NUMBER:
		tok := tokens.Next()
		return LitExpr{tok}
	case IDENT:
		tok := tokens.Next()
		return VarExpr{tok}
	case STRING:
		tok := tokens.Next()
		return LitExpr{tok}
	case LPAREN:
		expr := GroupingExpr{p.Equality(tokens)}

		tokens.Consume(RPAREN)
		return expr
	}

	panic(parseError(ILLEGAL, op))
}
