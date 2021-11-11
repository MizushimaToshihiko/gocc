//
// AST parser
//
package main

import (
	"fmt"
	"math"
)

// the types of AST node
type NodeKind int

const (
	ND_ADD       NodeKind = iota // 0: +
	ND_SUB                       // 1: -
	ND_MUL                       // 2: *
	ND_DIV                       // 3: /
	ND_EQ                        // 4: ==
	ND_NE                        // 5: !=
	ND_LT                        // 6: <
	ND_LE                        // 7: <=
	ND_ASSIGN                    // 8: =
	ND_PRE_INC                   // 9: pre ++
	ND_PRE_DEC                   // 10: pre --
	ND_POST_INC                  // 11: post ++
	ND_POST_DEC                  // 12: post --
	ND_A_ADD                     // 13: +=
	ND_A_SUB                     // 14: -=
	ND_A_MUL                     // 15: *=
	ND_A_DIV                     // 16: /=
	ND_COMMA                     // 17: ,
	ND_VAR                       // 18: local or global variables
	ND_NUM                       // 19: integer
	ND_RETURN                    // 20: 'return'
	ND_IF                        // 21: "if"
	ND_WHILE                     // 22: "while"
	ND_FOR                       // 23: "for"
	ND_BLOCK                     // 24: {...}
	ND_FUNCCALL                  // 25: function call
	ND_MEMBER                    // 26: . (struct member access)
	ND_ADDR                      // 27: unary &
	ND_DEREF                     // 28: unary *
	ND_NOT                       // 29: !
	ND_EXPR_STMT                 // 30: expression statement
	ND_STMT_EXPR                 // 31: statement expression
	ND_CAST                      // 32: type cast
	ND_NULL                      // 33: empty statement
	ND_SIZEOF                    // 34: "sizeof" operator
	ND_BITNOT                    // 35: ~
	ND_BITAND                    // 36: &
	ND_BITOR                     // 37: |
	ND_BITXOR                    // 38: ^
	ND_LOGAND                    // 40: &&
	ND_LOGOR                     // 41: ||
	ND_BREAK                     // 42: "break"
	ND_CONTINUE                  // 43: "continue"
	ND_GOTO                      // 44: "goto"
	ND_LABEL                     // 45: label
	ND_SWITCH                    // 46: "switch"
	ND_CASE                      // 47: "break"
	ND_SHL                       // 48: <<
	ND_SHR                       // 49: >>
	ND_A_SHL                     // 50: <<=
	ND_A_SHR                     // 51: >>=
	ND_TERNARY                   // 52: "?:"
)

// define AST node
type Node struct {
	Kind NodeKind // the type of node
	Next *Node    // the next node
	Ty   *Type    // the data type
	Tok  *Token   // current token

	Lhs *Node // the left branch
	Rhs *Node // the right branch

	// "if" or "while" of "for" statement
	Cond *Node
	Then *Node
	Els  *Node
	Init *Node
	Inc  *Node

	// block or statement expression
	Body *Node

	// struct member access
	MemName string
	Mem     *Member

	// for function call
	FuncName string
	Args     *Node

	// goto or labeld statement
	LabelName string

	// Switch-cases
	CaseNext   *Node
	DefCase    *Node
	CaseLbl    int
	CaseEndLbl int

	Val int64 // it would be used when 'Kind' is 'ND_NUM'
	Var *Var  // it would be used when 'Kind' is 'ND_VAR'
}

type Initializer struct {
	Next *Initializer

	// constant expression
	Sz  int
	Val int64

	// reference to another global variable
	Label string
}

func newNode(kind NodeKind, lhs *Node, rhs *Node, tok *Token) *Node {
	return &Node{
		Kind: kind,
		Lhs:  lhs,
		Rhs:  rhs,
		Tok:  tok,
	}
}

func newNodeNum(val int64, tok *Token) *Node {
	return &Node{
		Kind: ND_NUM,
		Val:  val,
		Tok:  tok,
	}
}

func newUnary(kind NodeKind, expr *Node, tok *Token) *Node {
	return &Node{Kind: kind, Tok: tok, Lhs: expr}
}

// the type of local variables
type Var struct {
	Name    string // the name of the variable
	Ty      *Type  // the data type
	Tok     *Token // for error message
	IsLocal bool   // local or global

	// local variables
	Offset int // the offset from RBP

	// global vaiables
	Contents    []rune
	ContLen     int
	Initializer *Initializer
}

type VarList struct {
	Next *VarList
	Var  *Var
}

// scope for local variables, global variables or typedefs
type VarScope struct {
	Next    *VarScope
	Name    string
	Depth   int
	Var     *Var
	TyDef   *Type
	EnumTy  *Type
	EnumVal int
}

// scope for struct tags
type TagScope struct {
	Next  *TagScope
	Name  string
	Depth int
	Ty    *Type
}

type Scope struct {
	VarScope *VarScope
	TagScope *TagScope
}

// local variables
var locals *VarList
var globals *VarList

var varScope *VarScope
var tagScope *TagScope
var scopeDepth int

var currentSwitch *Node

func enterScope() *Scope {
	sc := &Scope{
		VarScope: varScope,
		TagScope: tagScope,
	}
	scopeDepth++
	return sc
}

func leaveScope(sc *Scope) {
	varScope = sc.VarScope
	tagScope = sc.TagScope
	scopeDepth--
}

// findVar searchs a variable by name.
// if it wasn't find, return nil.
func findVar(tok *Token) *VarScope {
	for sc := varScope; sc != nil; sc = sc.Next {
		if len(sc.Name) == tok.Len && tok.Str == sc.Name {
			return sc
		}
	}
	return nil
}

func findTag(tok *Token) *TagScope {
	for sc := tagScope; sc != nil; sc = sc.Next {
		if len(sc.Name) == tok.Len && tok.Str == sc.Name {
			return sc
		}
	}
	return nil
}

func newVar(lvar *Var, tok *Token) *Node {
	return &Node{Kind: ND_VAR, Tok: tok, Var: lvar}
}

func pushScope(name string) *VarScope {
	sc := &VarScope{
		Name:  name,
		Next:  varScope,
		Depth: scopeDepth,
	}
	varScope = sc
	return sc
}

func pushVar(name string, ty *Type, isLocal bool, tok *Token) *Var {
	lvar := &Var{
		Name:    name,
		Ty:      ty,
		IsLocal: isLocal,
		Tok:     tok,
	}

	vl := &VarList{Var: lvar}

	if isLocal {
		vl.Next = locals
		locals = vl
	} else if ty.Kind != TY_FUNC {
		vl.Next = globals
		globals = vl
	}

	return lvar
}

func findTypedef(tok *Token) *Type {
	if tok.Kind == TK_IDENT {
		sc := findVar(token)
		if sc != nil {
			return sc.TyDef
		}
	}
	return nil
}

// for newLabel function
var cnt int

func newLabel() string {
	res := fmt.Sprintf(".L.data.%d", cnt)
	cnt++
	return res
}

type Function struct {
	Next   *Function
	Name   string
	Params *VarList

	Node    *Node
	Locals  *VarList
	StackSz int
}

type Program struct {
	Globals *VarList
	Fns     *Function
}

func isFunction() bool {
	tok := token

	ty := typeSpecifier()
	var name string
	declarator(ty, &name)
	isFunc := name != "" && consume("(") != nil

	token = tok
	return isFunc
}

// program = (global-var | function*)
func program() *Program {
	// printCurTok()
	// printCurFunc()
	cur := &Function{}
	head := cur
	globals = nil

	for !atEof() {
		if isFunction() {
			fn := function()
			if fn == nil {
				continue
			}
			cur.Next = fn
			cur = cur.Next
			continue
		}

		globalVar()
	}

	prog := &Program{Globals: globals, Fns: head.Next}
	return prog
}

// type-specifier = builtin-type | struct-decl | typedef-name | enum-specifier
//
// builtin-type   = "void"
//                | "_Bool"
//                | "char"
//                | "short" | "short" "int" | "int" "short"
//                | "int"
//                | "long" | "long" "int" | "int" "long"
//
// node that "typedef" and "static" can appear anywhere in a type-specifier
func typeSpecifier() *Type {
	if !isTypename() {
		panic("\n" + errorTok(token, "typename expected"))
	}

	var ty *Type = nil

	const (
		VOID  = 1 << 1
		BOOL  = 1 << 3
		CHAR  = 1 << 5
		SHORT = 1 << 7
		INT   = 1 << 9
		LONG  = 1 << 11
	)

	baseTy := 0
	var userTy *Type

	isTypedef := false
	isStatic := false

	for {
		// read one token at a time.
		tok := token
		if consume("typedef") != nil {
			isTypedef = true
		} else if consume("static") != nil {
			isStatic = true
		} else if consume("void") != nil {
			baseTy += VOID
		} else if consume("_Bool") != nil {
			baseTy += BOOL
		} else if consume("char") != nil {
			baseTy += CHAR
		} else if consume("short") != nil {
			baseTy += SHORT
		} else if consume("int") != nil {
			baseTy += INT
		} else if consume("long") != nil {
			baseTy += LONG
		} else if peek("struct") != nil {
			if baseTy != 0 || userTy != nil {
				break
			}
			userTy = structDecl()
		} else if peek("enum") != nil {
			if baseTy != 0 || userTy != nil {
				break
			}
			userTy = enumSpecifier()
		} else {
			if baseTy != 0 || userTy != nil {
				break
			}
			ty_ := findTypedef(token)
			if ty_ == nil {
				break
			}
			token = token.Next
			userTy = ty_
		}

		switch baseTy {
		case VOID:
			ty = voidType()
		case BOOL:
			ty = boolType()
		case CHAR:
			ty = charType()
		case SHORT, SHORT + INT:
			ty = shortType()
		case INT:
			ty = intType()
		case LONG, LONG + INT:
			ty = longType()
		case 0:
			// if there's no type specifier, it becomes int.
			// for expample, 'typedef x' defines x as an alias for int.
			if userTy != nil {
				ty = userTy
			} else {
				ty = intType()
			}
		default:
			panic("\n" + errorTok(tok, "invalid type"))
		}
	}
	ty.IsTypedef = isTypedef
	ty.IsStatic = isStatic
	return ty
}

// declarator = "*" ("(" declarator ")") | ident) type-suffix
func declarator(ty *Type, name *string) *Type {
	for consume("*") != nil {
		ty = pointerTo(ty)
	}

	if consume("(") != nil {
		placeholder := &Type{}
		newTy := declarator(placeholder, name)
		expect(")")
		*placeholder = *typeSuffix(ty)
		return newTy
	}

	*name = expectIdent()
	return typeSuffix(ty)
}

// abstract-declarator = "*"* ("(" abstract-declarator ")")? type-suffix
func abstractDeclarator(ty *Type) *Type {
	for consume("*") != nil {
		ty = pointerTo(ty)
	}

	if consume("(") != nil {
		placeholder := &Type{}
		newTy := abstractDeclarator(placeholder)
		expect(")")
		*placeholder = *typeSuffix(ty)
		return newTy
	}
	return typeSuffix(ty)
}

// type-suffix = ("[" const-expr? "]" type-suffix)?
func typeSuffix(ty *Type) *Type {
	if consume("[") == nil {
		return ty
	}

	var sz int64
	var isIncomp bool = true
	if consume("]") == nil {
		sz = constExpr()
		isIncomp = false
		expect("]")
	}

	ty = typeSuffix(ty)
	ty = arrayOf(ty, uint16(sz))
	ty.IsIncomp = isIncomp
	return ty
}

func typeName() *Type {
	ty := typeSpecifier()
	ty = abstractDeclarator(ty)
	return typeSuffix(ty)
}

func pushTagScope(tok *Token, ty *Type) {
	sc := &TagScope{
		Next:  tagScope,
		Name:  tok.Str,
		Ty:    ty,
		Depth: scopeDepth,
	}
	tagScope = sc
}

// struct-decl = "struct" ident? ("{" struct-member "}")?
func structDecl() *Type {

	// read struct tag.
	expect("struct")
	tag := consumeIdent()
	if tag != nil && peek("{") == nil {
		sc := findTag(tag)

		if sc == nil {
			ty := structType()
			pushTagScope(tag, ty)
			return ty
		}

		if sc.Ty.Kind != TY_STRUCT {
			panic("\n" + errorTok(tag, "not a struct tag"))
		}
		return sc.Ty
	}

	// Although it looks weird, "struct *foo" is legal C that defines
	// foo as a pointer to an unnamed incomplete struct type.
	if consume("{") == nil {
		return structType()
	}

	sc := findTag(tag)
	var ty *Type

	if sc != nil && sc.Depth == scopeDepth {
		// If there's an existing struct type having the same tag name in
		// the same block scope, this is a redefinition.
		if sc.Ty.Kind != TY_STRUCT {
			panic("\n" + errorTok(tag, "not a struct tag"))
		}
		ty = sc.Ty
	} else {
		// Register a struct type as an incomplete type early, so that you
		// can write recursive structs such as
		// "struct T { struct T *next; }".
		ty = structType()
		if tag != nil {
			pushTagScope(tag, ty)
		}
	}

	// read struct members.
	head := &Member{}
	cur := head

	for consume("}") == nil {
		cur.Next = structMember()
		cur = cur.Next
	}

	ty.Mems = head.Next

	// assign offsets within the struct to members.
	offset := 0
	for mem := ty.Mems; mem != nil; mem = mem.Next {
		offset = alignTo(offset, mem.Ty.Align)
		mem.Offset = offset
		offset += sizeOf(mem.Ty, mem.Tok)

		if ty.Align < mem.Ty.Align {
			ty.Align = mem.Ty.Align
		}
	}

	// register the struct type if a name was given.
	ty.IsIncomp = false
	return ty
}

// enum-specifier = "enum" ident
//                | "enum" ident? "{" enum-list? "}"
//
// enum-list = enum-elem ("," enum-elem)* ","?
// enum-elem = ident ("=" const-expr)?
func enumSpecifier() *Type {
	expect("enum")
	ty := enumType()

	// read an enum tag
	tag := consumeIdent()
	if tag != nil && peek("{") == nil {
		sc := findTag(tag)
		if sc == nil {
			panic("\n" + errorTok(tag, "unknown enum type"))
		}
		if sc.Ty.Kind != TY_ENUM {
			panic("\n" + errorTok(tag, "not an enum tag"))
		}
		return sc.Ty
	}

	expect("{")

	// read enum-list
	var cnt int64 = 0
	for {
		name := expectIdent()
		if consume("=") != nil {
			cnt = constExpr()
		}

		sc := pushScope(name)
		sc.EnumTy = ty
		sc.EnumVal = int(cnt)
		cnt++

		if consume(",") != nil {
			if consume("}") != nil {
				break
			}
			continue
		}
		expect("}")
		break
	}

	if tag != nil {
		pushTagScope(tag, ty)
	}
	return ty
}

// struct-member = type-specifier declarator type-suffix ";"
func structMember() *Member {
	var ty *Type = typeSpecifier()
	tok := token
	var name string
	ty = declarator(ty, &name)
	ty = typeSuffix(ty)
	expect(";")

	mem := &Member{Ty: ty, Name: name, Tok: tok}
	return mem
}

// param = type-specifier declarator type-suffix
func readFuncParam() *VarList {
	ty := typeSpecifier()
	var name string
	tok := token
	ty = declarator(ty, &name)
	ty = typeSuffix(ty)

	// "array of T" is converted to "pointer to T" only in the parameter
	// context. For examplem *argv[] is converted to **argv by this.
	if ty.Kind == TY_ARRAY {
		ty = pointerTo(ty.PtrTo)
	}

	var_ := pushVar(name, ty, true, tok)
	pushScope(name).Var = var_

	vl := &VarList{Var: var_}
	return vl
}

// params = param ("," param)*
func readFuncParams() *VarList {
	// printCurTok()
	// printCurFunc()
	if consume(")") != nil { // no argument
		return nil
	}

	head := readFuncParam()
	cur := head

	for {
		if consume(")") != nil {
			break
		}
		expect(",")
		cur.Next = readFuncParam()
		cur = cur.Next
	}

	return head
}

// function = type-specifier declarator "(" params? ")" ("{" stmt* "}" | ";")
func function() *Function {
	// printCurTok()
	// printCurFunc()
	locals = nil

	ty := typeSpecifier()
	var name string
	tok := token
	ty = declarator(ty, &name)

	// add a function type to the scope
	var_ := pushVar(name, funcType(ty), false, tok)
	pushScope(name).Var = var_

	// construct a function object
	fn := &Function{Name: name}
	expect("(")
	fn.Params = readFuncParams()

	if consume(";") != nil {
		return nil
	}

	// read function body
	cur := &Node{}
	head := cur
	expect("{")

	for {
		if t := consume("}"); t != nil {
			break
		}
		cur.Next = stmt()
		cur = cur.Next
	}

	fn.Node = head.Next
	fn.Locals = locals
	return fn
}

// Initializer list can end either with "}" or "," followed by "}"
// to allow a trailing comma. This function returns true if it looks
// like we are at the end of an initializer list.
func peekEnd() bool {
	tok := token
	ret := consume("}") != nil || (consume(",") != nil && consume("}") != nil)
	token = tok
	return ret
}

func consumeEnd() bool {
	tok := token
	if (consume("}") != nil || consume(",") != nil) && consume("}") != nil {
		return true
	}
	token = tok
	return false
}

func newInitVal(cur *Initializer, sz int, val int) *Initializer {
	init := &Initializer{Sz: sz, Val: int64(val)}
	cur.Next = init
	return init
}

func newInitLabel(cur *Initializer, label string) *Initializer {
	init := &Initializer{Label: label}
	cur.Next = init
	return init
}

func newInitZero(cur *Initializer, nbytes int) *Initializer {
	for i := 0; i < nbytes; i++ {
		cur = newInitVal(cur, 1, 0)
	}
	return cur
}

func gvarInitString(p string, len int) *Initializer {
	head := &Initializer{}
	head.Next = nil
	cur := head
	for i := 0; i < len; i++ {
		cur = newInitVal(cur, 1, int(p[i]))
	}
	return head.Next
}

func emitStructPadding(cur *Initializer, parent *Type, mem *Member) *Initializer {
	end := mem.Offset + sizeOf(mem.Ty, token)

	var padding int
	if mem.Next != nil {
		padding = mem.Next.Offset - end
	} else {
		padding = sizeOf(parent, token) - end
	}

	if padding != 0 {
		cur = newInitZero(cur, padding)
	}
	return cur
}

func skipExcessElements2() {
	for {
		if consume("{") != nil {
			skipExcessElements2()
		} else {
			assign()
		}

		if consumeEnd() {
			return
		}
		expect(",")
	}
}

func skipExcessElements() {
	expect(",")
	warnTok(token, "excess elements in initializer")
	skipExcessElements2()
}

func gvarInitializer(cur *Initializer, ty *Type) *Initializer {
	tok := token

	if ty.Kind == TY_ARRAY {
		open := consume("{") != nil
		i := 0
		var limit int
		if ty.IsIncomp {
			limit = math.MaxInt32
		} else {
			limit = int(ty.ArraySize)
		}

		for {
			cur = gvarInitializer(cur, ty.PtrTo)
			i++
			if i < limit && !peekEnd() && consume(",") != nil {
				continue
			}
			break
		}

		if open && !consumeEnd() {
			skipExcessElements()
		}

		// set excess arra elements to zero.
		cur = newInitZero(cur,
			sizeOf(ty.PtrTo, tok)*(int(ty.ArraySize)-i))

		if ty.IsIncomp {
			ty.ArraySize = uint16(i)
			ty.IsIncomp = false
		}

		return cur
	}

	if ty.Kind == TY_STRUCT {
		open := consume("{") != nil
		mem := ty.Mems

		for {
			cur = gvarInitializer(cur, mem.Ty)
			cur = emitStructPadding(cur, ty, mem)
			mem = mem.Next
			if mem != nil && !peekEnd() && consume(",") != nil {
				continue
			}
			break
		}

		if open && !consumeEnd() {
			skipExcessElements()
		}

		// set excess struct elements to zero.
		if mem != nil {
			sz := sizeOf(ty, tok) - mem.Offset
			if sz != 0 {
				cur = newInitZero(cur, sz)
			}
		}
		return cur
	}

	open := consume("{") != nil
	expr := conditional()
	if open {
		expect("}")
	}

	if expr.Kind == ND_ADDR {
		if expr.Lhs.Kind != ND_VAR {
			panic("\n" + errorTok(tok, "invalid initializer"))
		}
		return newInitLabel(cur, expr.Lhs.Var.Name)
	}

	if expr.Kind == ND_VAR && expr.Var.Ty.Kind == TY_ARRAY {
		return newInitLabel(cur, expr.Var.Name)
	}

	return newInitVal(cur, sizeOf(ty, token), int(eval(expr)))
}

// global-var = type-specifier declarator type-suffix ("=" gvar-initializer)? ";"
func globalVar() {
	ty := typeSpecifier()
	var name string
	tok := token
	ty = declarator(ty, &name)
	ty = typeSuffix(ty)

	var_ := pushVar(name, ty, false, tok)
	pushScope(name).Var = var_

	if consume("=") != nil {
		head := &Initializer{}
		gvarInitializer(head, ty)
		var_.Initializer = head.Next
	}
	expect(";")
}

type Designator struct {
	Next *Designator
	Idx  int     // array
	Mem  *Member // struct
}

// Creates a node for an array access. For example, if var represents
// a variable x and desg represents indices 3 and 4, this function
// returns a node representing x[3][4].
func newDesgNode2(va *Var, desg *Designator) *Node {
	tok := va.Tok
	if desg == nil {
		return newVar(va, tok)
	}

	node := newDesgNode2(va, desg.Next)

	if desg.Mem != nil {
		node = newUnary(ND_MEMBER, node, desg.Mem.Tok)
		node.MemName = desg.Mem.Name
		return node
	}

	node = newNode(ND_ADD, node, newNodeNum(int64(desg.Idx), tok), tok)
	return newUnary(ND_DEREF, node, tok)
}

func newDesgNode(va *Var, desg *Designator, rhs *Node) *Node {
	lhs := newDesgNode2(va, desg)
	node := newNode(ND_ASSIGN, lhs, rhs, rhs.Tok)
	return newUnary(ND_EXPR_STMT, node, rhs.Tok)
}

func lvarInitZero(cur *Node, va *Var, ty *Type, desg *Designator) *Node {
	if ty.Kind == TY_ARRAY {
		for i := 0; i < int(ty.ArraySize); i++ {
			desg2 := &Designator{Next: desg, Idx: i, Mem: nil}
			i++
			cur = lvarInitZero(cur, va, ty.PtrTo, desg2)
		}
		return cur
	}

	cur.Next = newDesgNode(va, desg, newNodeNum(0, token))
	return cur.Next
}

// lvar-initializer = assign
//                  | "{" lvar-initializer ("," lvar-initializer)* ","? "}"
//
// An initializer for a local variable is expanded to multiple
// assignments. For expample, this function creates the following
// nodes for x[2][3]={{1,2,3},{4,5,6}}.
//
//   x[0][0]=1;
//   x[0][1]=2;
//   x[0][2]=3;
//   x[1][0]=4;
//   x[1][1]=5;
//   x[1][2]=6;
//
// Struct members are initialized in declaration order. For example,
// `struct { int a; int b;} x = {1, 2}` sets x.a to 1 and x.b to 2.
//
// There are a few special rules for ambiguous initializer and
// shorthand notations:
//
// - If an initilizer list is shorter than an array, excess array
//   elements are initilized with 0.
//
// - A char array can be initializer by a string literal, For example,
//   `char x[4] = "foo"` is equivalent to `char x[4] = {'f', 'o', 'o',
//   '\0'}`.
//
// - If a rhs is an incomplete array, its size is set by counting the
//   number of items on the rhs. For example, `x` in `int x[]={1,2,3}`
//   has type `int[3]`.
func lvarInitializer(cur *Node, va *Var, ty *Type, desg *Designator) *Node {
	if ty.Kind == TY_ARRAY && ty.PtrTo.Kind == TY_CHAR &&
		token.Kind == TK_STR {
		// Initialize a char array with a string literal.
		tok := token
		token = token.Next

		if ty.IsIncomp {
			ty.ArraySize = uint16(tok.ContLen)
			ty.IsIncomp = false
		}

		var len int
		if ty.ArraySize < uint16(tok.ContLen) {
			len = int(ty.ArraySize)
		} else {
			len = tok.ContLen
		}

		var i int
		for i = 0; i < len; i++ {
			desg2 := &Designator{Next: desg, Idx: i, Mem: nil}
			rhs := newNodeNum(int64(tok.Contents[i]), tok)
			cur.Next = newDesgNode(va, desg2, rhs)
			cur = cur.Next
		}

		for ; i < int(ty.ArraySize); i++ {
			desg2 := &Designator{Next: desg, Idx: i, Mem: nil}
			cur = lvarInitZero(cur, va, ty.PtrTo, desg2)
		}
		return cur
	}

	if ty.Kind == TY_ARRAY {
		open := consume("{") != nil
		i := 0
		var limit int
		if ty.IsIncomp {
			limit = math.MaxInt32
		} else {
			limit = int(ty.ArraySize)
		}

		for {
			desg2 := &Designator{Next: desg, Idx: i, Mem: nil}
			i++
			cur = lvarInitializer(cur, va, ty.PtrTo, desg2)
			if i < limit && !peekEnd() && consume(",") != nil {
				continue
			}
			break
		}

		if open && !consumeEnd() {
			skipExcessElements()
		}

		// Set excess array elements to zero.
		for i < int(ty.ArraySize) {
			desg2 := &Designator{Next: desg, Idx: i}
			i++
			cur = lvarInitZero(cur, va, ty.PtrTo, desg2)
		}

		if ty.IsIncomp {
			ty.ArraySize = uint16(i)
			ty.IsIncomp = false
		}

		return cur
	}

	if ty.Kind == TY_STRUCT {
		open := consume("{") != nil
		mem := ty.Mems

		for {
			desg2 := &Designator{Next: desg, Idx: 0, Mem: mem}
			cur = lvarInitializer(cur, va, mem.Ty, desg2)
			mem = mem.Next
			if mem != nil && !peekEnd() && consume(",") != nil {
				continue
			}
			break
		}

		if open && !consumeEnd() {
			skipExcessElements()
		}

		// set excess struct elements to zero.
		for ; mem != nil; mem = mem.Next {
			desg2 := &Designator{Next: desg, Idx: 0, Mem: mem}
			cur = lvarInitZero(cur, va, mem.Ty, desg2)
		}
		return cur
	}

	open := consume("{") != nil
	cur.Next = newDesgNode(va, desg, assign())
	if open {
		expect("}")
	}
	return cur.Next
}

// declaration = type-specifier declarator type-suffix ("=" lvar-initializer)? ";"
//             | type-specifier ";"
func declaration() *Node {
	ty := typeSpecifier()
	if tok := consume(";"); tok != nil {
		return &Node{Kind: ND_NULL, Tok: tok}
	}

	tok := token
	var name string
	ty = declarator(ty, &name)
	ty = typeSuffix(ty)

	if ty.IsTypedef {
		expect(";")
		ty.IsTypedef = false
		pushScope(name).TyDef = ty
		return &Node{Kind: ND_NULL, Tok: tok}
	}

	if ty.Kind == TY_VOID {
		panic("\n" + errorTok(tok, "variable declared void"))
	}

	var var_ *Var
	if ty.IsStatic {
		var_ = pushVar(newLabel(), ty, false, tok)
	} else {
		var_ = pushVar(name, ty, true, tok)
	}
	pushScope(name).Var = var_

	if consume(";") != nil {
		return &Node{Kind: ND_NULL, Tok: tok}
	}

	expect("=")

	var head Node
	head.Next = nil
	lvarInitializer(&head, var_, var_.Ty, nil)
	expect(";")

	return &Node{Kind: ND_BLOCK, Tok: tok, Body: head.Next}
}

func readExprStmt() *Node {
	tok := token
	return &Node{Kind: ND_EXPR_STMT, Lhs: expr(), Tok: tok}
}

func isTypename() bool {
	return peek("void") != nil || peek("_Bool") != nil || peek("char") != nil ||
		peek("short") != nil || peek("int") != nil || peek("long") != nil ||
		peek("enum") != nil || peek("struct") != nil || peek("typedef") != nil ||
		peek("static") != nil || findTypedef(token) != nil
}

// stmt = "return" expr ";"
//      | "if" "(" expr ")" stmt ("else" stmt)?
//      | "switch" "(" expr ")" stmt
//      | "case" const-expr ":" stmt
//      | "default" ":" stmt
//      | "while" "(" expr ")" stmt
//      | "for" "(" (expr? ";" | declaration) expr? ";" expr? ")" stmt
//      | "{" stmt* "}"
//      | "break" ";"
//      | "continue" ";"
//      | "goto" ident ";"
//      | ident ":" stmt
//      | declaration
//      | expr ";"
func stmt() *Node {
	// printCurTok()
	// printCurFunc()
	var node *Node

	if t := consume("return"); t != nil {

		node = &Node{Kind: ND_RETURN, Lhs: expr(), Tok: t}
		expect(";")
		return node
	}

	if t := consume("if"); t != nil {
		node = &Node{Kind: ND_IF, Tok: t}
		expect("(")
		node.Cond = expr()
		expect(")")
		node.Then = stmt()

		if consume("else") != nil {
			node.Els = stmt()
		}
		return node
	}

	if t := consume("switch"); t != nil {
		node := &Node{Kind: ND_SWITCH, Tok: t}
		expect("(")
		node.Cond = expr()
		expect(")")

		sw := currentSwitch
		currentSwitch = node
		node.Then = stmt()
		currentSwitch = sw
		return node
	}

	if t := consume("case"); t != nil {
		if currentSwitch == nil {
			panic("\n" + errorTok(t, "stray case"))
		}
		val := constExpr()
		expect(":")

		node := newUnary(ND_CASE, stmt(), t)
		node.Val = val
		node.CaseNext = currentSwitch.CaseNext
		currentSwitch.CaseNext = node
		return node
	}

	if t := consume("default"); t != nil {
		if currentSwitch == nil {
			panic("\n" + errorTok(t, "stray default"))
		}
		expect(":")

		node = newUnary(ND_CASE, stmt(), t)
		currentSwitch.DefCase = node
		return node
	}

	if t := consume("while"); t != nil {
		node = &Node{Kind: ND_WHILE, Tok: t}
		expect("(")
		node.Cond = expr()
		expect(")")
		node.Then = stmt()
		return node
	}

	if t := consume("for"); t != nil {
		node = &Node{Kind: ND_FOR, Tok: t}
		expect("(")

		sc := enterScope()

		if consume(";") == nil {
			if isTypename() {
				node.Init = declaration()
			} else {
				node.Init = readExprStmt()
				expect(";")
			}
		}
		if consume(";") == nil {
			node.Cond = expr()
			expect(";")
		}
		if consume(")") == nil {
			node.Inc = readExprStmt()
			expect(")")
		}
		node.Then = stmt()

		leaveScope(sc)
		return node
	}

	if t := consume("{"); t != nil {

		head := Node{}
		cur := &head

		sc := enterScope()
		for {
			if consume("}") != nil {
				break
			}
			cur.Next = stmt()
			cur = cur.Next
		}
		leaveScope(sc)

		node = &Node{Kind: ND_BLOCK, Tok: t}
		node.Body = head.Next
		return node
	}

	if t := consume("break"); t != nil {
		expect(";")
		return &Node{Kind: ND_BREAK, Tok: t}
	}

	if t := consume("continue"); t != nil {
		expect(";")
		return &Node{Kind: ND_CONTINUE, Tok: t}
	}

	if t := consume("goto"); t != nil {
		node = &Node{Kind: ND_GOTO, Tok: t, LabelName: expectIdent()}
		expect(";")
		return node
	}

	if t := consumeIdent(); t != nil {
		if consume(":") != nil {
			node = newUnary(ND_LABEL, stmt(), t)
			node.LabelName = t.Str
			return node
		}
		token = t
	}

	if isTypename() {
		return declaration()
	}

	node = readExprStmt()
	expect(";")
	return node
}

// expr       = assign ("," assign)*
func expr() *Node {
	// printCurTok()
	// printCurFunc()
	node := assign()
	for {
		tok := consume(",")
		if tok == nil {
			break
		}
		node = newUnary(ND_EXPR_STMT, node, node.Tok)
		node = newNode(ND_COMMA, node, assign(), tok)
	}
	return node
}

func eval(node *Node) int64 {
	switch node.Kind {
	case ND_ADD:
		return eval(node.Lhs) + eval(node.Rhs)
	case ND_SUB:
		return eval(node.Lhs) - eval(node.Rhs)
	case ND_MUL:
		return eval(node.Lhs) * eval(node.Rhs)
	case ND_DIV:
		return eval(node.Lhs) / eval(node.Rhs)
	case ND_BITAND:
		return eval(node.Lhs) & eval(node.Rhs)
	case ND_BITOR:
		return eval(node.Lhs) | eval(node.Rhs)
	case ND_BITXOR:
		return eval(node.Lhs) ^ eval(node.Rhs)
	case ND_SHL:
		return eval(node.Lhs) << eval(node.Rhs)
	case ND_SHR:
		return eval(node.Lhs) >> eval(node.Rhs)
	case ND_EQ:
		b := eval(node.Lhs) == eval(node.Rhs)
		if b {
			return 1
		}
		return 0
	case ND_NE:
		b := eval(node.Lhs) != eval(node.Rhs)
		if b {
			return 1
		}
		return 0
	case ND_LT:
		b := eval(node.Lhs) < eval(node.Rhs)
		if b {
			return 1
		}
		return 0
	case ND_LE:
		b := eval(node.Lhs) <= eval(node.Rhs)
		if b {
			return 1
		}
		return 0
	case ND_TERNARY:
		if eval(node.Cond) != 0 {
			return eval(node.Then)
		}
		return eval(node.Els)
	case ND_COMMA:
		return eval(node.Rhs)
	case ND_NOT:
		if eval(node.Lhs) != 0 {
			return 0
		}
		return 1
	case ND_BITNOT:
		return ^eval(node.Lhs)
	case ND_LOGAND:
		if eval(node.Lhs) != 0 && eval(node.Rhs) != 0 {
			return 1
		}
		return 0
	case ND_LOGOR:
		if eval(node.Lhs) != 0 || eval(node.Rhs) != 0 {
			return 1
		}
		return 0
	case ND_NUM:
		return node.Val
	default:
		panic("\n" + errorTok(node.Tok, "not a constant expression"))
	}
}

func constExpr() int64 {
	return eval(conditional())
}

// assign     = conditional (assign-op assign)?
// assign-op  = "=" | "+=" | "-=" | "*=" | "/=" | "<<=" | ">>="
func assign() *Node {
	// printCurTok()
	// printCurFunc()
	node := conditional()

	if t := consume("="); t != nil {
		node = newNode(ND_ASSIGN, node, assign(), t)
	}
	if t := consume("+="); t != nil {
		node = newNode(ND_A_ADD, node, assign(), t)
	}
	if t := consume("-="); t != nil {
		node = newNode(ND_A_SUB, node, assign(), t)
	}
	if t := consume("*="); t != nil {
		node = newNode(ND_A_MUL, node, assign(), t)
	}
	if t := consume("/="); t != nil {
		node = newNode(ND_A_DIV, node, assign(), t)
	}
	if t := consume("<<="); t != nil {
		node = newNode(ND_A_SHL, node, assign(), t)
	}
	if t := consume(">>="); t != nil {
		node = newNode(ND_A_SHR, node, assign(), t)
	}

	return node
}

// conditional = logor ("?" expr ":" conditional)?
func conditional() *Node {
	node := logor()
	t := consume("?")
	if t == nil {
		return node
	}

	ternary := &Node{Kind: ND_TERNARY, Tok: t, Cond: node, Then: expr()}
	expect(":")
	ternary.Els = conditional()
	return ternary
}

// logor = logand ("||" logand)*
func logor() *Node {
	node := logand()
	for {
		tok := consume("||")
		if tok == nil {
			break
		}
		node = newNode(ND_LOGOR, node, logand(), tok)
	}
	return node
}

// logand = bitor ("&&" bitor)*
func logand() *Node {
	node := bitor()
	for {
		tok := consume("&&")
		if tok == nil {
			break
		}
		node = newNode(ND_LOGAND, node, bitor(), tok)
	}
	return node
}

// bitor = bitxor ("|" bitxor)*
func bitor() *Node {
	node := bitxor()
	for {
		tok := consume("|")
		if tok == nil {
			break
		}
		node = newNode(ND_BITOR, node, bitxor(), tok)
	}
	return node
}

// bitxor = bitand ("^" bitand)*
func bitxor() *Node {
	node := bitand()
	for {
		tok := consume("^")
		if tok == nil {
			break
		}
		node = newNode(ND_BITXOR, node, bitxor(), tok)
	}
	return node
}

// bitand = equality ("&" equality)*
func bitand() *Node {
	node := equality()
	for {
		tok := consume("&")
		if tok == nil {
			break
		}
		node = newNode(ND_BITAND, node, equality(), tok)
	}
	return node
}

// equality   = relational ("==" relational | "!=" relational)*
func equality() *Node {
	// printCurTok()
	// printCurFunc()
	node := relational()

	for {
		if t := consume("=="); t != nil {
			node = newNode(ND_EQ, node, relational(), t)
		} else if consume("!=") != nil {
			node = newNode(ND_NE, node, relational(), t)
		} else {
			return node
		}
	}
}

// relational = shift ("<" shift | "<=" shift | ">" shift | ">=" shift)*
func relational() *Node {
	// printCurTok()
	// printCurFunc()
	node := shift()

	for {
		if t := consume("<"); t != nil {
			node = newNode(ND_LT, node, shift(), t)
		} else if t := consume("<="); t != nil {
			node = newNode(ND_LE, node, shift(), t)
		} else if t := consume(">"); t != nil {
			node = newNode(ND_LT, shift(), node, t)
		} else if t := consume(">="); t != nil {
			node = newNode(ND_LE, shift(), node, t)
		} else {
			return node
		}
	}
}

// shift  = add ("<<" add | ">>" add)*
func shift() *Node {
	node := add()

	for {
		if t := consume("<<"); t != nil {
			node = newNode(ND_SHL, node, add(), t)
		} else if t := consume(">>"); t != nil {
			node = newNode(ND_SHR, node, add(), t)
		} else {
			return node
		}
	}
}

// add = mul ("+" mul | "-" mul)*
func add() *Node {
	// printCurTok()
	// printCurFunc()
	node := mul()

	for {
		if t := consume("+"); t != nil {
			node = newNode(ND_ADD, node, mul(), t)
		} else if t := consume("-"); t != nil {
			node = newNode(ND_SUB, node, mul(), t)
		} else {
			return node
		}
	}
}

// mul = cast ("*" cast | "/" cast)*
func mul() *Node {
	// printCurTok()
	// printCurFunc()
	node := cast()

	for {
		if t := consume("*"); t != nil {
			node = newNode(ND_MUL, node, cast(), t)
		} else if consume("/") != nil {
			node = newNode(ND_DIV, node, cast(), t)
		} else {
			return node
		}
	}
}

// cast = "(" type-name ")" cast | unary
func cast() *Node {
	tok := token

	if consume("(") != nil {
		if isTypename() {
			ty := typeName()
			expect(")")
			node := newUnary(ND_CAST, cast(), tok)
			node.Ty = ty
			return node
		}
		token = tok
	}

	return unary()
}

// unary = ("+" | "-" | "*" | "&" | "!" | "~")? cast
//       | ("++" | "--") unary
//       | "sizeof" "(" type-name ")"
//       | "sizeof" unary
//       | postfix
func unary() *Node {
	// printCurTok()
	// printCurFunc()
	if t := consumeSizeof(); t != nil {
		if consume("(") != nil {
			if isTypename() {
				ty := typeName()
				expect(")")
				return newNodeNum(int64(sizeOf(ty, t)), t)
			}
			token = t.Next
		}
		return newUnary(ND_SIZEOF, unary(), t)
	}

	if t := consume("+"); t != nil {
		return cast()
	}
	if t := consume("-"); t != nil {
		return newNode(ND_SUB, newNodeNum(0, t), cast(), t)
	}
	if t := consume("&"); t != nil {
		return newUnary(ND_ADDR, cast(), t)
	}
	if t := consume("*"); t != nil {
		return newUnary(ND_DEREF, cast(), t)
	}
	if t := consume("!"); t != nil {
		return newUnary(ND_NOT, cast(), t)
	}
	if t := consume("~"); t != nil {
		return newUnary(ND_BITNOT, cast(), t)
	}
	if t := consume("++"); t != nil {
		return newUnary(ND_PRE_INC, unary(), t)
	}
	if t := consume("--"); t != nil {
		return newUnary(ND_PRE_DEC, unary(), t)
	}
	return postfix()
}

// postfix = primary ("[" expr "]" | "." ident | "->" ident | "++" | "--")*
func postfix() *Node {
	node := primary()

	for {
		if tok := consume("["); tok != nil {
			// x[y] is short for *(x+y)
			exp := newNode(ND_ADD, node, expr(), tok)
			expect("]")
			node = newUnary(ND_DEREF, exp, tok)
			continue
		}

		if tok := consume("."); tok != nil {
			node = newUnary(ND_MEMBER, node, tok)
			node.MemName = expectIdent()
			continue
		}

		if tok := consume("->"); tok != nil {
			// x->y is shrot for (*x).y
			node = newUnary(ND_DEREF, node, tok)
			node = newUnary(ND_MEMBER, node, tok)
			node.MemName = expectIdent()
			continue
		}

		if tok := consume("++"); tok != nil {
			node = newUnary(ND_POST_INC, node, tok)
			continue
		}

		if tok := consume("--"); tok != nil {
			node = newUnary(ND_POST_DEC, node, tok)
			continue
		}

		return node
	}
}

// stmt-expr = "(" "{" stmt stmt* "}" ")"
//
// statement expression is a GNU extension.
func stmtExpr(tok *Token) *Node {
	sc := enterScope()

	node := &Node{
		Kind: ND_STMT_EXPR,
		Tok:  tok,
		Body: stmt(),
	}
	cur := node.Body

	for {
		if consume("}") != nil {
			break
		}
		cur.Next = stmt()
		cur = cur.Next
	}
	expect(")")

	leaveScope(sc)

	if cur.Kind != ND_EXPR_STMT {
		panic("\n" +
			errorTok(cur.Tok, "stmt expr returning void is not supported"))
	}
	*cur = *cur.Lhs
	return node
}

// func-args = "(" (assign("," assign)*)? ")"
func funcArgs() *Node {
	// printCurTok()
	// printCurFunc()
	if consume(")") != nil {
		return nil
	}

	head := assign()
	cur := head
	for {
		if consume(",") == nil {
			break
		}
		cur.Next = assign()
		cur = cur.Next
	}
	expect(")")
	return head
}

// primary = "(" "{" stmt-expr-tail
//         | ident func-args?
//         | "(" expr ")"
//         | num
//         | str
func primary() *Node {
	// printCurTok()
	// printCurFunc()

	// if the next token is '(', the program must be
	// "(" expr ")"
	if t := consume("("); t != nil {
		if consume("{") != nil {
			return stmtExpr(t)
		}

		node := expr()
		expect(")")
		return node
	}

	if tok := consumeIdent(); tok != nil {
		var node *Node
		if t := consume("("); t != nil { // function call
			node = &Node{
				Kind:     ND_FUNCCALL,
				Tok:      tok,
				FuncName: tok.Str,
				Args:     funcArgs(),
			}

			sc := findVar(tok)
			if sc != nil {
				if sc.Var == nil || sc.Var.Ty.Kind != TY_FUNC {
					panic("\n" + errorTok(tok, "not a function"))
				}
				node.Ty = sc.Var.Ty.RetTy
			} else {
				node.Ty = intType()
			}
			return node
		}

		// local variables
		sc := findVar(tok)
		if sc != nil {
			if sc.Var != nil {
				return newVar(sc.Var, tok)
			}
			if sc.EnumTy != nil {
				return newNodeNum(int64(sc.EnumVal), tok)
			}
		}
		panic("\n" + errorTok(tok, "undefined variable"))
	}

	tok := token
	if tok.Kind == TK_STR {
		token = token.Next

		ty := arrayOf(charType(), uint16(tok.ContLen))
		var_ := pushVar(newLabel(), ty, false, nil)
		var_.Initializer = gvarInitString(string(tok.Contents), tok.ContLen)
		return newVar(var_, tok)
	}

	if tok.Kind != TK_NUM {
		panic("\n" + errorTok(tok, "expected expression"))
	}
	// otherwise, must be integer
	return newNodeNum(expectNumber(), tok)
}
