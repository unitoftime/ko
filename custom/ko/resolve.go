package main

import (
	"fmt"

	"github.com/unitoftime/flow/ds"
)

type Scope struct {
	ident map[string]Node // A map of identifiers for the scope
	types map[string]Node // A map of types for the scope
}
func NewScope() *Scope{
	return &Scope{
		ident: make(map[string]Node),
		types: make(map[string]Node),
	}
}
func (s *Scope) AddIdent(name string, n Node) {
	_, exists := s.ident[name]
	if exists { panic("identifier already exists " + name) }

	s.ident[name] = n
}
func (s *Scope) CheckIdent(name string) (Node, bool) {
	n, ok := s.ident[name]
	return n, ok
}
// func (s *Scope) AddType(name string, t Node) {
// 	_, exists := s.types[name]
// 	if exists { panic("Type already exists " + name) }
// 	s.types[name] = t
// }
// func (s *Scope) CheckType(name string) (Node, bool) {
// 	t, ok := s.types[name]
// 	return t, ok
// }

type Resolver struct {
	builtin *Scope
	global *Scope // TODO: This could just be 0 in the scope stack?
	scopes *ds.Stack[*Scope]
}

func (r *Resolver) GlobalScope() bool {
	return !r.LocalScope()
}
func (r *Resolver) LocalScope() bool {
	return len(r.scopes.Buffer) > 0
}

func (r *Resolver) Scope() *Scope {
	if len(r.scopes.Buffer) == 0 {
		return r.global
	}
	return r.scopes.Buffer[len(r.scopes.Buffer) - 1]
}

// func (r *Resolver) GetCallExprType(obj CallExpr) (Node, bool) {
// 	CheckScope(
// }

func (r *Resolver) CheckScopeField(obj Node, field string) (Node, bool) {
	switch t := obj.(type) {
	case *IdentExpr:
		n, ok := r.CheckScope(t.tok.str)
		if !ok {
			printErr(t.tok, fmt.Sprintf("Unknown Variable: %s", t.tok.str))
		}
		strType := typeOf(n)
		structNode, ok := r.CheckScope(string(strType))
		if !ok {
			printErr(t.tok, fmt.Sprintf("Unknown type: %s", t.tok.str))
		}
		return getField(structNode, field)

	default:
		panic(fmt.Sprintf("Resolve: Unknown NodeType: %T", t))
	}
}

func getField(n Node, field string) (Node, bool) {
	if n == nil { return nil, false }

	switch t := n.(type) {
	case *StructNode:
		for i := range t.fields {
			if t.fields[i].name.str == field {
				return t.fields[i], true
			}
		}
		return nil, false
	default:
		panic(fmt.Sprintf("Resolve: Unknown NodeType: %T", t))
	}
}


// func (r *Resolver) GetType(name string) (Node, bool) {
// 	for i := len(r.scopes.Buffer) - 1; i >= 0; i-- {
// 		n, ok := r.scopes.Buffer[i].CheckType(name)
// 		if ok {
// 			return n, true
// 		}
// 	}

// 	n, ok := r.global.CheckType(name)
// 	if ok {
// 		return n, true
// 	}

// 	// Fallback to builtin
// 	return r.builtin.CheckType(name)
// }

func (r *Resolver) CheckScope(name string) (Node, bool) {
	for i := len(r.scopes.Buffer) - 1; i >= 0; i-- {
		n, ok := r.scopes.Buffer[i].CheckIdent(name)
		if ok {
			return n, true
		}
	}

	n, ok := r.global.CheckIdent(name)
	if ok {
		return n, true
	}

	// Fallback to builtin
	return r.builtin.CheckIdent(name)
}

func (r *Resolver) PushScope() *Scope {
	r.scopes.Add(NewScope())
	return r.Scope()
}

func (r *Resolver) PopScope() {
	r.scopes.Remove()
}

func NewResolver() *Resolver {
	builtin := NewScope()
	builtin.AddIdent("printf", nil)
	builtin.AddIdent("Assert", nil)

	return &Resolver{
		builtin: builtin,
		global: NewScope(),
		scopes: ds.NewStack[*Scope](),
	}
}

func (r *Resolver) Resolve(result ParseResult) {
	// do this
	// 1. register globally scoped things but dont do any resolving of types
	// 2. do local resolving for *everything* including global, but fallback to global lookups for things as needed. Then if something was already registered globally you dont have to reregister it in global scope (or maybe just ignore registering global scoped things, but just do their type checking)

	fmt.Println("--- Global ---")
	r.RegisterGlobal(result)
	fmt.Println("--- Local ---")
	r.ResolveLocal(result.file)
}

func (r *Resolver) RegisterGlobal(result ParseResult) {
	for i := range result.typeList {
		r.registerGlobal(result.typeList[i])
	}

	for i := range result.fnList {
		r.registerGlobal(result.fnList[i])
	}

	for i := range result.varList {
		r.registerGlobal(result.varList[i])
	}
}

func (r *Resolver) ResolveLocal(n Node) {
	r.resolveLocal(n)
}

func (r *Resolver) registerGlobal(n Node) {
	switch t := n.(type) {
	case *FileNode:
		for _, nn := range t.nodes {
			r.registerGlobal(nn)
		}
	case *PackageNode:
		// Skip
	case *CommentNode:
		// Skip
	case *FuncNode:
		r.global.AddIdent(t.name, t) // Register function name
	case *VarStmt:
		r.global.AddIdent(t.name.str, t) // Register the global identifier
	case *StructNode:
		r.global.AddIdent(t.ident.str, t) // Register struct name
	case *ArgNode:
	case *Stmt:
	case *CompLitExpr:
	case *LitExpr:

	case *ForStmt:
		printErr(t.tok, fmt.Sprintf("For loop not supported in global scope: %s", t.tok.str))

	case *IfStmt:
		tok := Token{}
		printErr(tok, fmt.Sprintf("If statements not supported in global scope: %s", tok.str))

	case *CurlyScope:
		tok := Token{}
		printErr(tok, fmt.Sprintf("Curly Scopes are not supported in global scope: %s", tok.str))

	case *ReturnNode:
		tok := Token{}
		printErr(tok, fmt.Sprintf("Return statements are not supported in global scope: %s", tok.str))

	case *CallExpr:
		tok := Token{}
		printErr(tok, fmt.Sprintf("Call Expressions are not supported in global scope: %s", tok.str))

	case *BinaryExpr:
		tok := Token{}
		printErr(tok, fmt.Sprintf("Binary Expressions are not supported in global scope: %s", tok.str))

	case *AssignExpr:
		tok := Token{}
		printErr(tok, fmt.Sprintf("Assign Expressions are not supported in global scope: %s", tok.str))
	case *IdentExpr:
		printErr(t.tok, fmt.Sprintf("Assign Expressions are not supported in global scope: %s", t.tok.str))

	default:
		panic(fmt.Sprintf("Resolve: Unknown NodeType: %T", t))
	}
}

// Returns the type of that node, if untyped returns ""
// For functions: returns the type that the node expression returns
func (r *Resolver) resolveLocal(node Node) Type {
	switch t := node.(type) {
	case *FileNode:
		for _, nn := range t.nodes {
			r.resolveLocal(nn)
		}
	case *PackageNode:
		// Skip
	case *CommentNode:
		// Skip
	case *StructNode:
		if r.LocalScope() {
			r.Scope().AddIdent(t.ident.str, t)
		}
	case *FuncNode:
		m := r.PushScope()

		if t.arguments != nil {
			for _, arg := range t.arguments.args {
				m.AddIdent(arg.name.str, arg)
			}
		}
		if t.returns != nil {
			for _, ret := range t.returns.args {
				m.AddIdent(ret.name.str, ret)
			}
		}

		r.resolveLocal(t.body)

		r.PopScope()

		return typeOf(t)
	case *Stmt:
		return r.resolveLocal(t.node)
	case *ForStmt:
		// TODO: Invalid unless in function
		r.PushScope()
		r.resolveLocal(t.init)
		r.resolveLocal(t.cond)
		r.resolveLocal(t.inc)

		r.resolveLocal(t.body)
		r.PopScope()

	case *IfStmt:
		// TODO: Invalid unless in function
		r.resolveLocal(t.cond)

		r.PushScope()
		r.resolveLocal(t.thenScope)
		r.PopScope()

		if t.elseScope == nil {
		} else {
			r.PushScope()
			r.resolveLocal(t.elseScope)
			r.PopScope()
		}

	case *CurlyScope:
		for i := range t.nodes {
			r.resolveLocal(t.nodes[i])
		}

	case *VarStmt:
		ty := UnknownType
		ty = r.resolveLocal(t.initExpr)
		if r.LocalScope() {
			r.Scope().AddIdent(t.name.str, t) // For global we register it before
		}
		t.ty = ty
		return ty

	case *Arg:
		return typeOf(t)
	case *ArgNode:
		panic("ARGNODE")
		// TODO: Are these just for func args? Maybe add to global scoping check

		// for i := range t.args {
		// 	buf.Add(t.args[i].kind).
		// 		Add(" ").
		// 		Add(t.args[i].name).
		// 		Add(" ")
		// 	if i < len(t.args)-1 {
		// 		buf.Add(", ")
		// 	}
		// }

	case *ReturnNode:
		// TODO: Check to make sure return type matches func return type or blank if void
		return r.resolveLocal(t.expr)

	case *CallExpr:
		callTy := r.resolveLocal(t.callee) // TODO: How does this work if it returns a call target?

		for i := range t.args {
			r.resolveLocal(t.args[i])
		}

		return callTy

		// fmt.Println("CallExpr:", callTy)
		// // n, ok := r.CheckScope(t.callee)
		// // fmt.Println("n, ok", n, ok)
		// // TODO: This is the callee type, but then need to look it up and find what type it returns
		// return "TODO"

	case *GetExpr:
		r.resolveLocal(t.obj)

		n, ok := r.CheckScopeField(t.obj, t.name.str)
		if !ok {
			printErr(t.name, fmt.Sprintf("Unknown Variable: %s", t.name.str))
		}
		return typeOf(n)
	case *SetExpr:
		r.resolveLocal(t.obj)

		n, ok := r.CheckScopeField(t.obj, t.name.str)
		if !ok {
			printErr(t.name, fmt.Sprintf("Unknown Variable: %s", t.name.str))
		}
		return typeOf(n)
	case *AssignExpr:
		valType := r.resolveLocal(t.value)
		n, ok := r.CheckScope(t.name.str)
		if !ok {
			printErr(t.name, fmt.Sprintf("Missing Variable: %s", t.name.str))
		}

		objType := typeOf(n)
		if valType != objType {
			printErr(t.name, fmt.Sprintf("Mismatched types: %s, %s, %s", objType, valType, "="))
		}

		return ""

	case *BinaryExpr:
		// lType := r.resolveLocal(t.left)
		// rType := r.resolveLocal(t.right)
		// resultType, success := checkBinaryExpr(lType, rType, t.op)
		resultType, success := r.checkBinaryExpr(t)
		if !success {
			// printErr(t.op, fmt.Sprintf("Mismatched types: %s, %s, %s", lType, t.op.str, rType))
			return UnknownType
		}
		t.ty = resultType
		return resultType

	case *UnaryExpr:
		// TODO: Impl

	case *IdentExpr:
		// TODO: Check that we have the needed variable
		node, ok := r.CheckScope(t.tok.str)
		if !ok {
			printErr(t.tok, fmt.Sprintf("Undefined Variable: %s", t.tok.str))
			return ""
		}
		t.ty = typeOf(node)
		return t.ty

	case *CompLitExpr:
		ty := r.resolveLocal(t.callee)
		for i := range t.args {
			r.resolveLocal(t.args[i])
		}
		t.ty = ty
		// fmt.Println("CompLitExpr:", t)
		return ty

	case *LitExpr:
		return typeOf(t)
	default:
		panic(fmt.Sprintf("Resolve: Unknown NodeType: %T", t))
	}

	return ""
}

// Returns the resulting type of the binary expression, and bool if the type check was ok
func (r *Resolver) checkBinaryExpr(t *BinaryExpr) (Type, bool) {
	left := r.resolveLocal(t.left)
	right := r.resolveLocal(t.right)

	typeCheck := (left == right)
	commonType := left
	if !typeCheck {
		var success bool
		commonType, success = checkLitTypeCast(left, right)

		if !success {
			printErr(t.op, fmt.Sprintf("Mismatched types: %s, %s, %s", left, t.op.str, right))
			return UnknownType, false
		}
	}

	switch t.op.token {
	// Comparable
	case BANGEQUAL: fallthrough
	case EQUALEQUAL:
		if !r.isComparable(t.left) {
			printErr(t.op, fmt.Sprintf("Tried to compare incomparable type: %T", t.left))
			panic("Tried to compare incomparable types")
			return UnknownType, false
		}

		return "bool", true

	// Ordered
	case GREATER: fallthrough
	case GREATEREQUAL: fallthrough
	case LESS: fallthrough
	case LESSEQUAL:
		// if !r.isComparable(t.left) {
		// 	printErr(t.op, fmt.Sprintf("Tried to compare incomparable type: %T", t.left))
		// 	panic("Tried to compare incomparable types")
		// 	return UnknownType, false
		// }

		return "bool", true

		// Result matches original
	case SUB: fallthrough
	case ADD: fallthrough
	case DIV: fallthrough
	case MUL:
		return commonType, true
	default:
		printErr(t.op, "checkBinaryExpr: Missing expression type")
		panic("AAAA")
	}
}

// // Returns the resulting type of the binary expression, and bool if the type check was ok
// func checkBinaryExpr(left, right Type, op Token) (Type, bool) {
// 	typeCheck := (left == right)
// 	commonType := left
// 	if !typeCheck {
// 		var success bool
// 		commonType, success = checkLitTypeCast(left, right)

// 		if !success {
// 			return UnknownType, false
// 		}
// 	}


// 	switch op.token {
// 	// Comparable
// 	case BANGEQUAL: fallthrough
// 	case EQUALEQUAL:
// 		return "bool", true

// 	// Ordered
// 	case GREATER: fallthrough
// 	case GREATEREQUAL: fallthrough
// 	case LESS: fallthrough
// 	case LESSEQUAL:
// 		return "bool", true

// 		// Result matches original
// 	case SUB: fallthrough
// 	case ADD: fallthrough
// 	case DIV: fallthrough
// 	case MUL:
// 		return commonType, true
// 	default:
// 		printErr(op, "checkBinaryExpr: Missing expression type")
// 		panic("AAAA")
// 	}
// }

// Check if the left or right side can be type cast to the other from a literal to a concrete
// Returns the resulting concrete type and true if successful, else returns false
func checkLitTypeCast(left, right Type) (Type, bool) {
	ok := tryCast(left, right)
	if ok {
		// fmt.Println("TryLitTypeCast", right, ok)
		return right, ok
	}

	ok = tryCast(right, left)
	if ok {
		// fmt.Println("TryLitTypeCast", left,  ok)
		return left, ok
	}

	return UnknownType, false
}
