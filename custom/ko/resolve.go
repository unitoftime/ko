package main

import (
	"fmt"

	"github.com/unitoftime/flow/ds"
)

type Scope struct {
	funcNode *FuncNode
	ident map[string]Node // A map of identifiers for the scope
	// types map[string]Node // A map of types for the scope
}
func NewScope() *Scope{
	return &Scope{
		ident: make(map[string]Node),
		// types: make(map[string]Node),
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

	// Output data
	genericInstantiations []GenericInstance
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

func (r *Resolver) AddIdent(name string, n Node) {
	Println("AddIdent", name)

	// Check Shadowing
	_, exists := r.CheckScope(name)
	if exists {
		nodeError(n, fmt.Sprintf("declartion shadows a previous variable: %s", name))
		panic("AA")
	}

	// Add
	r.Scope().AddIdent(name, n)
}

// func (r *Resolver) GetCallExprType(obj CallExpr) (Node, bool) {
// 	CheckScope(
// }

// TODO: Resolve selector path?
func (r *Resolver) CheckScopeField(obj Node, field string) (Node, bool) {
	Println("---")
	Println("CheckScopeField:", obj, field)

	switch t := obj.(type) {
	case *GetExpr:
		Println("CheckField:.GetExpr:", t.obj, t.name.str, field)
		// return r.CheckScopeField(t.obj, field)
		left, ok := r.CheckScopeField(t.obj, t.name.str)
		if !ok {
			panic("AAA")
		}
		Printf("CheckField.Left: %T, %v\n", left, left)
		ret, ok := r.CheckScopeField(left, field)
		Println("Returning:", ret)
		return ret, ok
	case *IdentExpr:
		Println("CheckField.IdentExpr:", t, field)
		n, ok := r.CheckScope(t.tok.str)
		if !ok {
			errUndefinedVar(n, t.tok.str)
		}
		strType := n.Type()

		underlyingName := strType.Underlying().Name()
		structNode, ok := r.CheckScope(underlyingName)
		if !ok {
			errUndefinedType(t, underlyingName)
		}

		Println("StructNode: ", t, field)
		return getField(structNode, field)
	case *Arg:
		Println("CheckField:.Arg:", t, t.Type().Name(), field)
		argTypeName := t.Type().Name()
		structNode, ok := r.CheckScope(argTypeName)
		if !ok {
			errUndefinedType(structNode, argTypeName)
		}
		ret, ok := getField(structNode, field)
		Println("CheckField:.Arg.Return:", ret, ok)
		return ret, ok
	default:
		panic(fmt.Sprintf("Resolve: Unknown NodeType: %T", t))
	}
}

func getField(n Node, field string) (Node, bool) {
	if n == nil { return nil, false }

	switch t := n.(type) {
	case *StructNode:
		Println("getField:", t, field)
		for i := range t.fields {
			if t.fields[i].name.str == field {
				Println("Found Field:", t.fields[i])
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

func (r *Resolver) GetFuncScope() (*FuncNode, bool) {
	for i := len(r.scopes.Buffer) - 1; i >= 0; i-- {
		if r.scopes.Buffer[i].funcNode != nil {
			return r.scopes.Buffer[i].funcNode, true
		}
	}

	return nil, false
}

func (r *Resolver) PushScope() *Scope {
	Println("--->", r.scopes.Len())
	r.scopes.Add(NewScope())
	return r.Scope()
}

func (r *Resolver) PopScope() {
	Println("<---", r.scopes.Len())
	r.scopes.Remove()
}

func NewResolver() *Resolver {
	builtin := NewScope()
	builtin.AddIdent("printf", &BuiltinNode{getType(&FuncType{returns: VoidType})}) // TODO: ARGS
	builtin.AddIdent("ko_printf", &BuiltinNode{getType(&FuncType{returns: VoidType})}) // TODO: ARGS
	builtin.AddIdent("Assert", &BuiltinNode{getType(&FuncType{returns: VoidType})}) // TODO: macro?
	// builtin.AddIdent("ko_byte_malloc", &BuiltinNode{getType(&BasicType{"uint8_t*", false})})
	// builtin.AddIdent("sizeof", &BuiltinNode{getType(&BasicType{"size_t", false})})

	// builtin.AddIdent("printf", &BuiltinNode{getType(&BasicType{"void", false})})
	// builtin.AddIdent("Assert", &BuiltinNode{getType(&BasicType{"void", false})})
	// builtin.AddIdent("ko_byte_malloc", &BuiltinNode{getType(&BasicType{"uint8_t*", false})})
	// builtin.AddIdent("sizeof", &BuiltinNode{getType(&BasicType{"size_t", false})})

	builtin.AddIdent("append", &BuiltinNode{AppendBuiltinType})
	builtin.AddIdent("len", &BuiltinNode{LenBuiltinType})

	// Add builtin types
	builtin.AddIdent("nil", &BuiltinNode{PointerLitType})

	builtin.AddIdent("u8", &BuiltinNode{getType(&BasicType{"u8", true})})
	builtin.AddIdent("u16", &BuiltinNode{getType(&BasicType{"u16", true})})
	builtin.AddIdent("u32", &BuiltinNode{getType(&BasicType{"u32", true})})
	builtin.AddIdent("u64", &BuiltinNode{getType(&BasicType{"u64", true})})
	builtin.AddIdent("i8", &BuiltinNode{getType(&BasicType{"i8", true})})
	builtin.AddIdent("i16", &BuiltinNode{getType(&BasicType{"i16", true})})
	builtin.AddIdent("i32", &BuiltinNode{getType(&BasicType{"i32", true})})
	builtin.AddIdent("i64", &BuiltinNode{getType(&BasicType{"i64", true})})

	builtin.AddIdent("f64", &BuiltinNode{getType(Float64Type)})
	builtin.AddIdent("string", &BuiltinNode{getType(StringType)})

	builtin.AddIdent("int", &BuiltinNode{getType(IntType)})
	builtin.AddIdent("uintptr", &BuiltinNode{getType(IntType)})
	builtin.AddIdent("usize", &BuiltinNode{getType(IntType)})


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

	Println("--- Global ---")
	r.RegisterGlobal(result)
	Println("--- Local ---")
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


	// Resolve
	for i := range result.typeList {
		Println("Resolve:", result.typeList[i])
		r.resolveLocal(result.typeList[i])
	}

	for i := range result.fnList {
		Println("Resolve:", result.fnList[i])
		// r.resolveLocal(result.fnList[i])
		r.resolveFuncNodePrototype(result.fnList[i])
	}

	for i := range result.varList {
		Println("ResolveVarList:", result.varList[i])
		ty := r.resolveLocal(result.varList[i])
		Println("ResolveVarList.Finish:", ty)
		Println(result.varList[i])
	}
}

func (r *Resolver) resolveFuncNodePrototype(t *FuncNode) {
	// TODO: Build a Func type that has all the filds and returns
	funcType := &FuncType{
		name: t.name,
	}

	r.PushScope()
	defer r.PopScope()

	r.tryPushGenericArgs(t)

	// Generics
	genTy := make([]Type, 0)
	if t.generic != nil {
		for i := range t.generic.Args {
			genTy = append(genTy, r.resolveLocal(t.generic.Args[i]))
		}
	}
	funcType.generics = genTy

	// Args
	argTy := make([]Type, 0)
	for i := range t.arguments.args {
		at := r.resolveLocal(t.arguments.args[i])
		argTy = append(argTy, at)
		Printf("FuncArgs %s: %d, %T\n", t.name, i, at)
	}
	funcType.args = argTy

	// Returns
	if t.returns == nil || len(t.returns.args) == 0 {
		funcType.returns = VoidType
	} else {
		funcType.returns = r.resolveLocal(t.returns.args[0])
	}

	t.ty = funcType
}

func (r *Resolver) tryPushGenericArgs(t *FuncNode) {
	if t.generic == nil { return }

	for _, genArg := range t.generic.Args {
		genArg.ty = r.ResolveTypeNodeExpr(genArg)
		r.AddIdent(genArg.name.str, genArg)
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
		// t.ty = typeOf(t)
		r.global.AddIdent(t.ident.str, t) // Register struct name
	case *ArgNode:
	case *Stmt:
	case *CompLitExpr:
	case *LitExpr:

	case *ForStmt:
		nodeError(t, "For loop not supported in global scope")

	case *IfStmt:
		nodeError(t, "If statements not supported in global scope")

	case *CurlyScope:
		nodeError(t, "Curly Scopes are not supported in global scope")

	case *ReturnNode:
		nodeError(t, "Return statements are not supported in global scope")

	case *CallExpr:
		nodeError(t, "Call Expressions are not supported in global scope")

	case *BinaryExpr:
		nodeError(t, "Binary Expressions are not supported in global scope")

	case *AssignExpr:
		nodeError(t, "Assign Expressions are not supported in global scope")
	case *IdentExpr:
		nodeError(t, "Identifier Expressions are not supported in global scope")

	default:
		panic(fmt.Sprintf("Resolve: Unknown NodeType: %T", t))
	}
}

// Returns the type of that node, if untyped returns ""
// For functions: returns the type that the node expression returns
func (r *Resolver) resolveLocal(node Node) Type {
	// if node.Type() != nil {
	// 	return node.Type()
	// }

	switch t := node.(type) {
	case *FileNode:
		for _, nn := range t.nodes {
			r.resolveLocal(nn)
		}
	case *PackageNode:
		// Skip
	case *CommentNode:
		// Skip
	case *ForeignScope:
		r.resolveLocal(t.body)
	case *StructNode:
		if t.ty != UnknownType {
			return t.ty
		}

		Println("StructNode:", t)
		fields := make([]Type, len(t.fields))
		for i, field := range t.fields {
			t.ty = r.resolveLocal(field)
			fields[i] = t.ty
		}

		t.ty = getType(&StructType{t.ident.str, fields})

		if r.LocalScope() {
			r.AddIdent(t.ident.str, t)
		}
		return t.ty
	case *FuncNode:
		Println("FuncNode:", t.name)

		// Note: This only handles function body, the function type gets resolved earlier
		m := r.PushScope()
		m.funcNode = t

		r.tryPushGenericArgs(t)

		if t.arguments != nil {
			for _, arg := range t.arguments.args {
				at := r.resolveLocal(arg)
				Printf("t.arguments.args: %T, %v\n", at, arg)
				r.AddIdent(arg.name.str, arg)
			}
		}
		if t.returns != nil {
			for _, ret := range t.returns.args {
				Println("t.returns.args", ret.name.str)
				r.resolveLocal(ret)
				if ret.name.str != "" {
					r.AddIdent(ret.name.str, ret)
				}
			}
		}

		Println("t.body", t.body)

		if t.Generic() {
			// If it is generic, we will do the type checking when we detect a new instantiation
		} else {
			if t.body != nil {
				r.resolveLocal(t.body)
			}
		}

		r.PopScope()

		return t.ty
	case *Stmt:
		Println("Stmt:", t)
		r.resolveLocal(t.node)
		return UnknownType
	case *ForStmt:
		Println("ForStmt:", t)
		// TODO: Invalid unless in function
		r.PushScope()
		r.resolveLocal(t.init)
		r.resolveLocal(t.cond)
		r.resolveLocal(t.inc)

		r.resolveLocal(t.body)
		r.PopScope()

	case *IfStmt:
		Println("IfStmt:", t)
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

	case *SwitchStmt:
		r.resolveLocal(t.cond)

		for i := range t.cases {
			r.PushScope()
			r.resolveLocal(t.cases[i])
			r.PopScope()
		}

	case *CaseStmt:
		if t.expr != nil {
			r.resolveLocal(t.expr)
		}
		r.resolveLocal(t.body)

	case *ScopeNode:
		Println("ScopeNode")
		r.PushScope()
		ty := r.resolveLocal(t.Scope)
		r.PopScope()

		return ty

	case *CurlyScope:
		Println("CurlyScope")
		for i := range t.nodes {
			Printf("CurlyScopeNode: %T\n", t.nodes[i])
			r.resolveLocal(t.nodes[i])
		}

	case *VarStmt:
		Printf("VarStmt: %s, %T +%v\n", t.name.str, t.initExpr, t.initExpr)
		if t.typeSpec != nil {
			t.ty = r.resolveLocal(t.typeSpec)
		} else if t.initExpr != nil {
			initTy := r.resolveLocal(t.initExpr)

			t.ty = resolveLitType(initTy)
		}

		if r.LocalScope() {
			r.AddIdent(t.name.str, t) // For global we register it before
		}

		// Fold Constant expressions as needed
		if t.constant {
			t.folded = r.FoldConstExpression(t.initExpr)
		}

		Printf("VarStmt.Typed: %+v\n", t)
		return t.ty

	case *Arg:
		Println("Resolve Arg:", t)
		t.ty = r.resolveLocal(t.typeNode)
		Println("Resolve Arg... type:", t.ty)

		return t.ty
	case *ArgNode:
		// TODO: Are these just for func args? Maybe add to global scoping check
		for _, a := range t.args {
			r.resolveLocal(a)
		}

	case *ReturnNode:
		Println("ReturnNode:", t)
		t.ty = r.resolveLocal(t.expr)
		// TODO: Check to make sure return type matches func return type or blank if void
		currentFuncScope, ok := r.GetFuncScope()
		if !ok {
			// TODO: Positioning
			nodeError(t, "Return statement must be inside a function")
		}

		//----------------------------------------
		// TODO: We only support 1 argument currently
		//----------------------------------------

		// Match return args with func args
		if len(currentFuncScope.returns.args) != 1 {
			// TODO: Positioning
			nodeError(t, "Mismatched return arguments")
		}

		if currentFuncScope.returns.args[0].Type() != t.ty {
			// TODO: Positioning
			nodeError(t, "Incorrect return type")
		}

		return t.ty

	case *CallExpr:
		Println("CallExpr:", t)
		callTy := r.resolveLocal(t.callee)

		// TODO: infer generic types

		// TODO: validate argument types match
		for i := range t.args {
			r.resolveLocal(t.args[i])
		}

		switch tt := callTy.(type) {
		case *FuncType:
			t.ty = tt.returns
			return t.ty
		case *BasicType:
			// TODO: validate typecast argument can cast to output type
			// TODO: Validate there is only one thing passed into cast operation
			t.ty = tt
			return t.ty
		default:
			nodeError(t, fmt.Sprintf("Unexpected call expressions type: %T", tt))
		}

	case *GetExpr:
		r.resolveLocal(t.obj)

		n, ok := r.CheckScopeField(t.obj, t.name.str)
		if !ok {
			errUndefinedVar(t, t.name.str)
		}
		t.ty = n.Type()
		Println("GetExpr:", t)

		return t.ty

	case *SetExpr:
		t.ty = r.resolveLocal(t.obj)
		r.resolveLocal(t.value)

		n, ok := r.CheckScopeField(t.obj, t.name.str)
		if !ok {
			errUndefinedVar(t, t.name.str)
		}
		Println("SetExpr:", t)
		return n.Type()
	case *IndexExpr:
		Println("IndexExpr:", t)
		objType := r.resolveLocal(t.callee)

		switch ot := objType.(type) {
		case *ArrayType:
			t.ty = ot.base
		case *SliceType:
			t.ty = ot.base

			// gt, ok := t.ty.(*GenericType)
			// if ok {
			// 	ct, ok := genericTypeMap[gt.name]
			// 	if !ok { panic("Unknown generic impl") }
			// 	t.ty = ct
			// }
		case *FuncType:
			t.ty = r.InstantiateGenericFunc(ot, t)
			it, ok := t.callee.(*IdentExpr)
			if ok {
				it.ty = t.ty // resolve the generic for the ident type
			}

			return t.ty
		default:
			nodeError(t, fmt.Sprintf("type: %s doesn't support indexing", objType.Name()))
		}


		// If we are here then it must be either an array or slice, so check casting to an int
		idxType := r.resolveLocal(t.index)
		// TODO: Technically only for array/slices: Ensure index type is castable to an int
		supportedIndexType := IntType

		Println("idxType", idxType, supportedIndexType)
		if !tryCast(idxType, supportedIndexType) {
			nodeError(t, "array index must be castable to an int")
		}

		return t.ty

	case *AssignExpr:
		Println("AssignExpr:", t)
		valType := r.resolveLocal(t.value)
		objType := r.resolveLocal(t.name)
		Println("AssignTypes:", valType, objType)

		if valType == UnknownType || objType == UnknownType {
			nodeError(t, fmt.Sprintf("UnknownAssignmentType: %s, %s", objType.Name(), valType.Name()))
		}

		// objType := n.Type()
		if !tryCast(valType, objType) {
			nodeError(t, fmt.Sprintf("Mismatched assignment types: %s, %s", objType.Name(), valType.Name()))
		}

		return UnknownType

	case *BinaryExpr:
		Println("BinaryExpr:", t)
		resultType, success := r.checkBinaryExpr(t)
		if !success {
			// printErr(t.op, fmt.Sprintf("Mismatched types: %s, %s, %s", lType, t.op.str, rType))
			return UnknownType
		}
		t.ty = resultType
		return resultType

	case *ShortStmt:
		lType := r.resolveLocal(t.target)
		rType := r.resolveLocal(t.initExpr)

		if !tryCast(rType, lType) {
			nodeError(t, fmt.Sprintf("Mismatched types: %s, %s", lType.Name(), rType.Name()))
		}

		return UnknownType

	case *PostfixStmt:
		// TODO: Currently you only support ++ and -- which both returh r-values and cant be used for other stuff
		r.resolveLocal(t.left)
		return UnknownType
	case *UnaryExpr:
		Println("UnaryExpr:", t)

		// TODO: Some unary operators may modify the type
		t.ty = r.resolveLocal(t.right)

		switch t.op.token {
		case MUL:
			// Dereferencing a pointer
			ptr, ok := t.ty.(*PointerType)
			if !ok {
				nodeError(t, "must be a pointer to dereference")
			}
			t.ty = ptr.base
		case AND:
			// getting an address of an object
			// TODO: Check t.ty must be addressable
			t.ty = getType(&PointerType{t.ty})
		}

		return t.ty

	case *IdentExpr:
		Println("IdentExpr:", t)
		// TODO: Check that we have the needed variable
		node, ok := r.CheckScope(t.tok.str)
		if !ok {
			errUndefinedIdent(t, t.tok.str)
			return UnknownType
		}
		t.ty = node.Type()

		switch tt := node.(type) {
		case *VarStmt:
			// For certain Indentifier types, lookup the folded constant if available
			if tt.constant {
				t.folded = tt.folded
			}
		}

		Printf("IdentExpr.Type: %T %v", t.ty, t)

		return t.ty

	case *CompLitExpr:
		Println("Resolve.CompLitExpr:", t)
		// t.ty = r.resolveLocal(t.callee)
		t.ty = r.ResolveTypeNodeExpr(t.callee)

		for i := range t.args {
			r.resolveLocal(t.args[i])
		}
		Println("CompLitExpr.Typed:", t.ty)
		return t.ty

	case *GroupingExpr:
		t.ty = r.resolveLocal(t.Node)
		return t.ty
	case *LitExpr:
		Println("LitExpr:", t)
		return t.Type()
		// t.ty = typeOf(t)
		// return t.ty
	case *BuiltinNode:
		return t.ty
	case *GenArg:
		t.ty = r.ResolveTypeNodeExpr(t)
		return t.ty
	case *TypeNode:
		Println("TypeNode:", t)
		t.ty = r.ResolveTypeNodeExpr(t.node)
		if t.ty == nil {
			panic("FAILED TO RESOLVE")
		}
		return t.ty

	default:
		panic(fmt.Sprintf("Resolve: Unknown NodeType: %T", t))
	}

	return UnknownType
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
			nodeError(t, fmt.Sprintf("Mismatched operator types: %s, %s", left.Name(), right.Name()))
			return UnknownType, false
		}
	}

	switch t.op.token {
	// Comparable
	case BANGEQUAL: fallthrough
	case EQUALEQUAL:
		if !isComparable(commonType) {
			nodeError(t, fmt.Sprintf("Type is not comparable: %s", commonType.Name()))
			panic("Tried to compare incomparable types")
			return UnknownType, false
		}

		return BoolType, true

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

		return BoolType, true

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
		// Println("TryLitTypeCast", right, ok)
		return right, ok
	}

	ok = tryCast(right, left)
	if ok {
		// Println("TryLitTypeCast", left,  ok)
		return left, ok
	}

	return UnknownType, false
}

func (r *Resolver) ResolveTypeNodeExpr(n Node) Type {
	if n.Type() != nil {
		return n.Type()
	}

	switch t := n.(type) {
	case *UnaryExpr:
		t.ty = r.ResolveTypeNodeExpr(t.right)

		// Note: Some unary operators modify the type
		switch t.op.token {
		case MUL:
			// It is a pointer
			t.ty = getType(&PointerType{t.ty})
		// case AND:
			// TODO: If you ever do reference types like c++
		default:
			panic("ASLFKJSAFL")
		}

		return t.ty

	case *IdentExpr:
		node, ok := r.CheckScope(t.tok.str)
		if !ok {
			errUndefinedIdent(t, t.tok.str)
			panic("AAA")
			return UnknownType
		}

		if !r.isTypedef(node) {
			errIdentMustBeAType(node, t.tok.str)
			panic("Must be a type!") // TODO: Improve error msg
		}

		t.ty = node.Type()
		Println("TYPETYPETYPE: ", node.Type())

		return t.ty
	case *ArrayNode:
		if t.len == nil {
			// TODO: Slice type
			elemType := r.ResolveTypeNodeExpr(t.elem)
			t.ty = getType(&SliceType{elemType})
		} else {
			lenExpr := r.resolveLocal(t.len)
			if !tryCast(lenExpr, IntType) {
				nodeError(t, "array length expression must be castable to an int")
			}
			r.resolveLocal(t.len)

			lenVal, _ := castToInt(t.len)

			elemType := r.ResolveTypeNodeExpr(t.elem)
			t.ty = getType(&ArrayType{lenVal, elemType})
		}
		return t.ty
	case *GenArg:
		t.ty = getType(&GenericType{t.name.str})
		return t.ty
	default:
		panic(fmt.Sprintf("ResolveTypeNodeExpr: Unknown NodeType: %T", t))
	}
}

// Returns true if the node is a typedef of some sort
func (r *Resolver) isTypedef(n Node) bool {
	// TODO: Add others if needed
	switch  n.(type) {
	case *StructNode:
		return true
	case *GenArg:
		return true
	case *BuiltinNode:
		// TODO: This isn't right, we technically need to see if the builtin type returned is a builtin type like an int or a u64
		return true
	}
	return false
}

// Returns true if the type is comparable
func isComparable(ty Type) bool {
	// TODO: Add others if needed
	switch t := ty.(type) {
	case *BasicType:
		return t.comparable
	case *StructType:
		return true
	}
	return false
}


type GenericInstance struct {
	node *IndexExpr
	funcNode *FuncNode
	ty Type
}

func (r *Resolver) InstantiateGenericFunc(t *FuncType, index *IndexExpr) Type {
	funcNodeNode, ok := r.CheckScope(t.name)
	if !ok {
		panic("Could not find identifier")
	}

	var funcNode *FuncNode
	switch tt := funcNodeNode.(type) {
	case *BuiltinNode:
	case *FuncNode:
		funcNode = tt
	default:
		panic("Only support funcNode generics")
	}

	genericMap := make(map[string]Type)

	// TODO: Check length of each, make sure they match
	for _, g := range t.generics {
		indexVal := r.resolveLocal(index.index) // TODO: index by i (the order params are passed to the index matches the order they are defined in the generics list)
		genericMap[g.Name()] = indexVal
	}

	// This is the concrete function
	finalFunc := &FuncType{
		name: t.name,
		generics: make([]Type, 0),
		args: make([]Type, 0),
	}

	for _, g := range t.generics {
		concreteArg, ok := genericMap[g.Name()]
		if !ok {
			nodeError(index, "undefined generic type")
			panic("AAA")
		}
		finalFunc.generics = append(finalFunc.generics, concreteArg)
	}

	for _, a := range t.args {
		genArg, isGenType := a.(*GenericType)
		if isGenType {
			// Convert to a concrete type
			concreteArg, ok := genericMap[genArg.Name()]
			if !ok {
				nodeError(index, "undefined generic type")
				panic("AAA")
			}

			finalFunc.args = append(finalFunc.args, concreteArg)
		} else {
			finalFunc.args = append(finalFunc.args, a)
		}
	}

	{
		genArg, isGenType := t.returns.(*GenericType)
		if isGenType {
			// Convert to a concrete type
			concreteArg, ok := genericMap[genArg.Name()]
			if !ok {
				nodeError(index, "undefined generic type")
				panic("AAA")
			}

			finalFunc.returns = concreteArg
		} else {
			finalFunc.returns = t.returns
		}
	}

	// Check if we've already instantiated this type/generic
	ty, ok := checkType(finalFunc)
	if ok {
		return ty
	}

	// if the funcNode is nil, then it is a builtin, so just return the type, no instantiation needed
	if funcNode == nil {
		return finalFunc
	}

	// Because we detected a new instantiation, we need to execute type checking
	// TODO: Need a way to proces the body. may need a generic mapping map like I do for code generation portion so the type checker can validate based on the correct types. You also need to repush all of the identifiers provided by the function declaration (like you do in resolveLocal)
	// if funcNode.body != nil {
	// 	r.resolveLocal(funcNode.body)
	// }
	{
		t := funcNode
		// Note: This only handles function body, the function type gets resolved earlier
		m := r.PushScope()
		m.funcNode = t

		r.tryPushGenericArgs(t)


		// TODO: This wont work for nested geneics. you need to make the map a nested scope thingy
		if t.generic != nil {
			for _, genArg := range t.generic.Args {
				concreteType, ok := genericMap[genArg.name.str]
				if !ok { continue }
				addGenericMap(genArg.name.str, concreteType)
			}
		}
		defer clearGenericMap()

		if t.arguments != nil {
			for _, arg := range t.arguments.args {
				Println("t.arguments.args", arg)
				r.resolveLocal(arg)
				r.AddIdent(arg.name.str, arg)
			}
		}
		if t.returns != nil {
			for _, ret := range t.returns.args {
				Println("t.returns.args", ret.name.str)
				r.resolveLocal(ret)
				if ret.name.str != "" {
					r.AddIdent(ret.name.str, ret)
				}
			}
		}

		Println("t.body", t.body)

		if t.body != nil {
			r.resolveLocal(t.body)
		}

		r.PopScope()
	}


	r.genericInstantiations = append(r.genericInstantiations, GenericInstance{
		node: index,
		funcNode: funcNode,
		ty: ty,
	})

	return getType(finalFunc)
}

func (r *Resolver) FoldConstExpression(initExpr Node) Node {
	switch t := initExpr.(type) {
	case *LitExpr:
		return t
	default:
		nodeError(t, fmt.Sprintf("FoldConstExpression: Unhandled type: %T", t))
		panic("AAAA")
	}
}
