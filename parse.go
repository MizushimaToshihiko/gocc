package main

import (
	"fmt"
	"reflect"
)

// Scope for local variables, global variables or typedefs
type VarScope struct {
	Next  *VarScope
	Name  string
	Var   *Var
	TyDef *Type
}

type Var struct {
	Name    string // Variable name
	Ty      *Type  // Type
	Tok     *Token // for error message
	IsLocal bool   // local or global

	// Local variables
	Offset int // Offset from RBP

	// Global variables
	Init *Initializer
}

type VarList struct {
	Next *VarList
	Var  *Var
}

type Program struct {
	Globs *VarList  // global variables
	Fns   *Function // functions
}

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
	ND_MEMBER                    // 9: . (struct menber access)
	ND_VAR                       // 10: variables
	ND_NUM                       // 11: integer
	ND_RETURN                    // 12: 'return'
	ND_IF                        // 13: "if"
	ND_WHILE                     // 14: "while"
	ND_FOR                       // 15: "for"
	ND_BLOCK                     // 16: {...}
	ND_FUNCALL                   // 17: function call
	ND_ADDR                      // 18: unary &
	ND_DEREF                     // 19: unary *
	ND_EXPR_STMT                 // 20: expression statement
	ND_CAST                      // 21: type cast
	ND_NULL                      // 22: empty statement
	ND_SIZEOF                    // 23: "sizeof"
	ND_COMMA                     // 24: comma
	ND_INC                       // 25: post ++
	ND_DEC                       // 26: post --
	ND_A_ADD                     // 27: +=
	ND_A_SUB                     // 28: -=
	ND_A_MUL                     // 29: *=
	ND_A_DIV                     // 30: /=
	ND_NOT                       // 31: !
	ND_BITNOT                    // 32: unary ^
	ND_BITAND                    // 33: &
	ND_BITOR                     // 34: |
	ND_BITXOR                    // 35: ^
	ND_LOGAND                    // 36: &&
	ND_LOGOR                     // 37: ||
	ND_BREAK                     // 38: "break"
	ND_CONTINUE                  // 39: "continue"
	ND_GOTO                      // 40: "goto"
	ND_LABEL                     // 41: Labeled statement
	ND_SWITCH                    // 42: "switch"
	ND_CASE                      // 43: "case"
	ND_SHL                       // 44: <<
	ND_SHR                       // 45: >>
	ND_A_SHL                     // 46: <<=
	ND_A_SHR                     // 47: >>=
)

// define AST node
type Node struct {
	Kind NodeKind // type of node
	Next *Node    // Next node
	Ty   *Type    // Type e.g. int or pointer to int
	Tok  *Token   // Representive token

	Lhs *Node // left branch
	Rhs *Node // right branch

	// "if" or "for" statement
	Cond *Node
	Then *Node
	Els  *Node
	Init *Node
	Inc  *Node

	// Block
	Body *Node

	// Struct member access
	MemName string
	Mem     *Member

	// Function call
	FuncName string
	Args     *Node

	// Goto or labeled statement
	LblName string

	// Switch-cases
	CaseNext   *Node
	DefCase    *Node
	CaseLbl    int
	CaseEndLbl int

	Var *Var  // used if kind == ND_VAR
	Val int64 // it would be used when kind is 'ND_NUM'
}

var locals *VarList
var globals *VarList

var varScope *VarScope

var curSwitch *Node

// findVar finds a variable or a typedef by name.
func findVar(tok *Token) *VarScope {
	for sc := varScope; sc != nil; sc = sc.Next {
		if len(sc.Name) == tok.Len && tok.Str == sc.Name {
			return sc
		}
	}
	return nil
}

func newNode(kind NodeKind, tok *Token) *Node {
	return &Node{Kind: kind, Tok: tok}
}

func newBinary(kind NodeKind, lhs *Node, rhs *Node, tok *Token) *Node {
	return &Node{
		Kind: kind,
		Tok:  tok,
		Lhs:  lhs,
		Rhs:  rhs,
	}
}

func newUnary(kind NodeKind, expr *Node, tok *Token) *Node {
	node := &Node{Kind: kind, Lhs: expr, Tok: tok}
	return node
}

func newNum(val int64, tok *Token) *Node {
	return &Node{
		Kind: ND_NUM,
		Tok:  tok,
		Val:  val,
	}
}

func newVar(v *Var, tok *Token) *Node {
	return &Node{Kind: ND_VAR, Tok: tok, Var: v}
}

func pushScope(name string) *VarScope {
	sc := &VarScope{Name: name, Next: varScope}
	varScope = sc
	return sc
}

func pushVar(name string, ty *Type, isLocal bool, tok *Token) *Var {
	// printCurTok()
	// printCalledFunc()

	v := &Var{Name: name, Ty: ty, IsLocal: isLocal, Tok: tok}

	var vl *VarList
	if isLocal {
		vl = &VarList{Var: v, Next: locals}
		locals = vl
	} else if ty.Kind != TY_FUNC {
		vl = &VarList{Var: v, Next: globals}
		globals = vl
	}

	pushScope(name).Var = v
	return v
}

func findTyDef(tok *Token) *Type {
	if tok.Kind == TK_IDENT {
		if sc := findVar(token); sc != nil {
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

// Global variable initializer. Global variables can be initialized
// either by a constant expression or a pointer to another global
// variable.
type Initializer struct {
	Next *Initializer

	// Constant expression
	Sz  int
	Val int64

	// Reference to another global variable
	Lbl string
}

type Function struct {
	Next   *Function
	Name   string
	Params *VarList

	Node    *Node
	Locals  *VarList
	StackSz int
}

// program = (global-var | function)*
func program() *Program {
	// printCurTok()
	// printCalledFunc()

	head := &Function{}
	cur := head
	globals = nil

	for !atEof() {
		if consume("func") != nil {
			cur.Next = function()
			cur = cur.Next
		} else if consume("var") != nil {
			globalVar()
		} else if consume("type") != nil {
			name := expectIdent()
			ty := readTypePreffix()
			pushScope(name).TyDef = ty
			expect(";")
		} else {
			panic("\n" + errorTok(token, "unexpected '%s'", token.Str))
		}
	}

	return &Program{Globs: globals, Fns: head.Next}
}

func stringDecl() *Type {
	return pointerTo(charType())
}

// type-specifier = "*"* builtin-type | struct-decl | typedef-name |
// builtin-type = void | "bool" | "byte"| "int16" | "int" | "int64"
func typeSpecifier() *Type {
	// printCurTok()
	// printCalledFunc()

	nPtr := 0
	for consume("*") != nil {
		nPtr++
	}

	var ty *Type
	if consume("byte") != nil {
		ty = charType()
	} else if consume("string") != nil {
		ty = stringDecl()
	} else if consume("bool") != nil {
		ty = boolType()
	} else if consume("int16") != nil {
		ty = shortType()
	} else if consume("int") != nil {
		ty = intType()
	} else if consume("int64") != nil {
		ty = longType()
	} else if peek("struct") != nil { // struct type
		ty = structDecl()
	} else if t := consumeIdent(); t != nil {
		ty = findVar(t).TyDef
	}

	if ty == nil {
		ty = voidType()
	}

	for i := 0; i < nPtr; i++ {
		ty = pointerTo(ty)
	}

	return ty
}

func findBase() (*Type, *Token) {
	// printCurTok()
	// printCalledFunc()

	tok := token
	for peek("*") == nil && !isTypename() {
		token = token.Next
	}
	ty := typeSpecifier()
	t := token // どこまでtokenを読んだか
	token = tok
	return ty, t
}

func readArr(base *Type) *Type {
	// printCurTok()
	// printCalledFunc()

	if consume("[") == nil {
		return base
	}
	sz := constExpr()
	expect("]")
	base = readArr(base)
	return arrayOf(base, int(sz))
}

// type-preffix = ("[" const-expr "]")*
func readTypePreffix() *Type {
	// printCurTok()
	// printCalledFunc()

	if peek("[") == nil {
		return typeSpecifier()
	}

	base, t := findBase()
	arrTy := readArr(base)
	token = t
	return arrTy
}

// struct-decl = "struct" "{" struct-member "}"
func structDecl() *Type {
	// printCurTok()
	// printCalledFunc()

	expect("struct")
	expect("{")

	head := &Member{}
	cur := head

	for consume("}") == nil {
		cur.Next = structMem()
		cur = cur.Next
	}

	ty := &Type{Kind: TY_STRUCT, Mems: head.Next}

	// Assign offsets within the struct to members.
	offset := 0
	for mem := ty.Mems; mem != nil; mem = mem.Next {
		offset = alignTo(offset, mem.Ty.Align)
		mem.Offset = offset
		offset += sizeOf(mem.Ty, mem.Tok)

		if ty.Align < mem.Ty.Align {
			ty.Align = mem.Ty.Align
		}
	}

	return ty
}

// struct-member = ident type-prefix type-specifier
func structMem() *Member {
	// printCurTok()
	// printCalledFunc()

	tok := token
	mem := &Member{
		Name: expectIdent(),
		Ty:   readTypePreffix(),
		Tok:  tok,
	}
	expect(";")
	return mem
}

// param = ident type-prefix type-specifier
// e.g.
//  x int
//  x *int
//  x **int
//  x [3]int
//  x [3]*int
//  x [2]**int
func readFuncParam() *VarList {
	// printCurTok()
	// printCalledFunc()

	tok := token
	name := expectIdent()
	ty := readTypePreffix()
	vl := &VarList{}
	vl.Var = pushVar(name, ty, true, tok)
	return vl
}

// params = param ("," param)*
func readFuncParams() *VarList {
	// printCurTok()
	// printCalledFunc()

	if consume(")") != nil {
		return nil
	}

	head := readFuncParam()
	cur := head

	for consume(")") == nil {
		expect(",")
		cur.Next = readFuncParam()
		cur = cur.Next
	}

	return head
}

// function = "func" ident "(" params? ")" type-prefix type-specifier "{" stmt "}"
func function() *Function {
	// printCurTok()
	// printCalledFunc()

	locals = nil

	// Construct a function object
	tok := token
	fn := &Function{Name: expectIdent()}
	expect("(")
	fn.Params = readFuncParams()
	ty := typeSpecifier()

	// Add a function type to the scope
	pushVar(fn.Name, funcType(ty), false, tok)
	expect("{")

	// Read function body
	head := &Node{}
	cur := head
	for consume("}") == nil {
		cur.Next = stmt()
		cur = cur.Next
	}
	expect(";")
	fn.Node = head.Next
	fn.Locals = locals
	return fn
}

// Initializer list can end with "}".
// This function returns true if it looks like
// we are at the end of an initializer list.
func peekEnd() bool {
	tok := token
	ret := consume("}") != nil
	token = tok
	return ret
}

func consumeEnd() bool {
	tok := token
	if consume("}") != nil {
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
	init := &Initializer{Lbl: label}
	cur.Next = init
	return init
}

func newInitZero(cur *Initializer, nbytes int) *Initializer {
	for i := 0; i < nbytes; i++ {
		cur = newInitVal(cur, 1, 0)
	}
	return cur
}

func gvarInitString(p []rune, len int) *Initializer {
	head := &Initializer{}
	cur := head
	for i := 0; i < len; i++ {
		cur = newInitVal(cur, 1, int(p[i]))
	}
	return head.Next
}

func emitStructPadding(cur *Initializer, parent *Type, mem *Member) *Initializer {
	end := mem.Offset + sizeOf(mem.Ty, token)

	padding := sizeOf(parent, token) - end
	if mem.Next != nil {
		padding = mem.Next.Offset - end
	}

	if padding != 0 {
		cur = newInitZero(cur, padding)
	}
	return cur
}

func gvarInitializer(cur *Initializer, ty *Type) *Initializer {
	// printCalledFunc()
	// printCurTok()

	tok := token

	ty2 := readTypePreffix()
	// if neither type-preffix nor ty-specifier, and "tok" is string literal
	if ty2.Kind == TY_VOID && tok.Kind == TK_STR {
		ty2 = stringDecl()
	}

	if ty2.Kind != TY_VOID && !isSameTy(ty, ty2) {
		panic("\n" + errorTok(tok,
			"connot use \"%s\" (type %s) as type %s in assignment", tok.Str, ty2.Kind, ty.Kind))
	}

	if ty.Kind == TY_ARRAY {

		if consume("{") == nil {
			panic("\n" + errorTok(tok, "invalid initializer"))
		}

		var i int
		limit := ty.ArrSz

		for {
			cur = gvarInitializer(cur, ty.Base)
			i++
			if i >= limit || peekEnd() || consume(",") == nil {
				break
			}
		}

		if !consumeEnd() {
			panic("\n" + errorTok(token, "array index %d out of bounds [0:%d]", i, limit))
		}

		// Set excess array elements to zero.
		if i < ty.ArrSz {
			cur = newInitZero(cur, sizeOf(ty.Base, tok)*(ty.ArrSz-i))
		}

		return cur
	}

	if ty.Kind == TY_STRUCT {
		if consume("{") == nil {
			panic("\n" + errorTok(tok, "invalid initializer"))
		}

		mem := ty.Mems

		for {
			cur = gvarInitializer(cur, mem.Ty)
			cur = emitStructPadding(cur, ty, mem)
			mem = mem.Next
			if mem == nil || peekEnd() || consume(",") == nil {
				break
			}
		}

		if !consumeEnd() {
			panic("\n" + errorTok(token, "too many values"))
		}

		// Set excess struct elements to zero.
		if mem != nil {
			sz := sizeOf(ty, tok) - mem.Offset
			if sz != 0 {
				cur = newInitZero(cur, sz)
			}
		}
		return cur
	}

	expr := logor()

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

// global-var = "var" ident type-prefix type-suffix ("=" gvar-initializer)? ";"
//
// For example,
// var x int = 6
// var x *int = &y
// var x string = "abc"
// var x [2]int = [2]int{1,2}
// var x T(typedef) = T{1,2}
func globalVar() {
	// printCurTok()
	// printCalledFunc()

	tok := token
	name := expectIdent()
	ty := readTypePreffix()

	v := pushVar(name, ty, false, tok)

	if consume("=") != nil {
		head := &Initializer{}
		gvarInitializer(head, ty)
		v.Init = head.Next
	}

	expect(";")
}

type Designator struct {
	Next *Designator
	Idx  int
	Mem  *Member
}

// Creates a node for an array sccess. For example, if v represents
// a variable x and desg represents indices 3 and 4, this function
// returns a node representing x[3][4]
func newDesgNode2(v *Var, desg *Designator) *Node {
	tok := v.Tok
	if desg == nil {
		return newVar(v, tok)
	}

	node := newDesgNode2(v, desg.Next)

	if desg.Mem != nil {
		node = newUnary(ND_MEMBER, node, desg.Mem.Tok)
		node.MemName = desg.Mem.Name
		return node
	}

	node = newBinary(ND_ADD, node, newNum(int64(desg.Idx), tok), tok)
	return newUnary(ND_DEREF, node, tok)
}

func newDesgNode(v *Var, desg *Designator, rhs *Node) *Node {
	lhs := newDesgNode2(v, desg)
	node := newBinary(ND_ASSIGN, lhs, rhs, rhs.Tok)
	return newUnary(ND_EXPR_STMT, node, rhs.Tok)
}

func lvarInitZero(cur *Node, v *Var, ty *Type, desg *Designator) *Node {
	if ty.Kind == TY_ARRAY {
		for i := 0; i < ty.ArrSz; i++ {
			desg2 := &Designator{desg, i, nil}
			i++
			cur = lvarInitZero(cur, v, ty.Base, desg2)
		}
		return cur
	}

	cur.Next = newDesgNode(v, desg, newNum(0, token))
	return cur.Next
}

func stringAssign(cur *Node, v *Var, ty *Type, desg *Designator, tok *Token) *Node {
	var length int = tok.ContLen
	if ty.ArrSz != tok.ContLen {
		ty.ArrSz = tok.ContLen
	}
	var i int

	for i = 0; i < length; i++ {
		desg2 := &Designator{desg, i, nil}
		rhs := newNum(int64(tok.Contents[i]), tok)
		cur.Next = newDesgNode(v, desg2, rhs)
		cur = cur.Next
	}

	for ; i < ty.ArrSz; i++ {
		desg2 := &Designator{desg, i, nil}
		cur = lvarInitZero(cur, v, ty.Base, desg2)
	}
	return cur
}

// lvar-initializer = assign
//                  | "{" lvar-initializer ("," lvar-initializer)* "}"
//
// An initializer for a local variable is expanded to multiple
// assignments. For example, this function creates the following
// nodes for var x [2][3]int=[2][3]int{{1,2,3},{4,5,6}}.
//
// x[0][0]=1
// x[0][1]=2
// x[0][2]=3
// x[1][0]=4
// x[1][1]=5
// x[1][2]=6
//
// Struct members are initialized in declaration order. For example,
// 'type x struct {
// 	a int
// 	b int
// }
// var x T = T{1, 2}'
// sets x.a to 1 and x.b to 2.
//
// If an initializer list is shorter than an array, excess array
// elements are initialized with 0.
//
// A string(char array) can be initialized by a string literal. For example,
// `var x string="abc"`
func lvarInitializer(cur *Node, v *Var, ty *Type, desg *Designator) *Node {
	// Initialize a char array with a string literal.
	// => unnecessary

	// Initialize an array or a struct
	tok2 := token
	ty2 := readTypePreffix()
	// if neither type-preffix nor ty-specifier, and "tok" is string literal
	if ty2.Kind == TY_VOID && tok2.Kind == TK_STR {
		ty2 = stringDecl()
	}

	if ty2.Kind != TY_VOID && !isSameTy(ty, ty2) {
		panic("\n" + errorTok(tok2,
			"connot use \"%s\" (type %s) as type %s in assignment", tok2.Str, ty2.Kind, ty.Kind))
	}

	if ty2.Kind != TY_STRUCT && ty2.Kind != TY_ARRAY {
		cur.Next = newDesgNode(v, desg, assign())
		return cur.Next
	}

	expect("{")
	tok := consume("{")
	if ty.Kind == TY_ARRAY {
		i := 0
		limit := ty.ArrSz

		for {
			desg2 := &Designator{desg, i, nil}
			i++
			cur = lvarInitializer(cur, v, ty.Base, desg2)
			if i >= limit || peekEnd() || consume(",") == nil {
				break
			}
		}

		if !consumeEnd() {
			panic("\n" + errorTok(token, "array index %d out of bounds [0:%d]", i, limit))
		}

		// Set excess array elements to zero.
		for i < ty.ArrSz {
			desg2 := &Designator{desg, i, nil}
			i++
			cur = lvarInitZero(cur, v, ty.Base, desg2)
		}
		return cur
	}

	if ty.Kind == TY_STRUCT {
		mem := ty.Mems

		for {
			desg2 := &Designator{desg, 0, mem}
			cur = lvarInitializer(cur, v, mem.Ty, desg2)
			mem = mem.Next
			if mem == nil || peekEnd() || consume(",") == nil {
				break
			}
		}

		if !consumeEnd() {
			panic("\n" + errorTok(token, "too many values"))
		}

		// Set excess struct elements to zero.
		for ; mem != nil; mem = mem.Next {
			desg2 := &Designator{desg, 0, mem}
			cur = lvarInitZero(cur, v, mem.Ty, desg2)
		}
		return cur
	}

	panic("\n" + errorTok(tok, "invalid initializer"))
}

func isSameTy(ty1, ty2 *Type) bool {
	if reflect.DeepEqual(ty1, ty2) {
		if ty1 == nil && ty2 == nil {
			return true
		}
		return isSameTy(ty1.Base, ty2.Base)
	}
	return false
}

// declaration = VarDecl | VarSpec(unimplemented) | ShortVarDecl(unimplemented)
// VarDecl = "var" ident type-prefix type-specifier ("=" expr)
// VarSpec = ident-list (type-preffix type-specifier [ "=" expr-list ] | "=" expr-list)
// ShortVarDecl = "var" ident "=" expr => unimplemented
//              | ident ":=" expr => unimplemented
func declaration() *Node {
	// printCurTok()
	// printCalledFunc()

	expect("var")
	tok := token

	name := expectIdent()
	ty := readTypePreffix()
	assert(ty.Kind != TY_VOID, "\n"+errorTok(tok, "variable declared void"))

	v := pushVar(name, ty, true, tok)
	if consume(";") != nil {
		return newNode(ND_NULL, tok)
	}
	// ここでShortVarDecl("var" ident = expr)の場合はty==nilでvがpushVarされていない状態 => unimplemented

	expect("=")

	// cannot assign array variables to array variables now.
	head := &Node{}
	lvarInitializer(head, v, v.Ty, nil)
	expect(";")

	node := newNode(ND_BLOCK, tok)
	node.Body = head.Next
	return node
}

func readExprStmt() *Node {
	// printCurTok()
	// printCalledFunc()

	// t := token
	return expr()
}

func isTypename() bool {
	// printCurTok()
	// printCalledFunc()

	return peek("byte") != nil || peek("bool") != nil ||
		peek("int16") != nil || peek("int") != nil ||
		peek("int64") != nil || peek("struct") != nil ||
		peek("string") != nil ||
		findTyDef(token) != nil
}

func isForClause() bool {
	// printCurTok()
	// printCalledFunc()

	tok := token

	for peek("{") == nil {
		if peek(";") != nil {
			token = tok
			return true
		}
		token = token.Next
	}
	token = tok
	return false
}

// stmt = "return" expr ";"
//      | "if" expr "{" stmt "};" ("else" "{" stmt "};" )?
//      | "switch" "{" expr "}" stmt
//      | "case" const-expr ":" stmt
//      | "default" ":" stmt
//      | for-stmt
//      | for-clause
//      | "{" stmt* "}"
//      | "break" ";"
//      | "continue" ";"
//      | "goto" ident ";"
//      | "ident ":" stmt
//      | "type" ident type-prefix basetype ";"
//      | declaration
//      | expr ";"
// for-stmt = "for" [ condition ] block .
// for-clause = "for" (expr? ";" | declaration) condition ";" expr? block
// condition = expr .
// block = "{" stmt-list "};" .
// stmt-list = { stmt ";" } .
func stmt() *Node {
	// printCurTok()
	// printCalledFunc()

	if t := consume("return"); t != nil {
		node := newUnary(ND_RETURN, expr(), t)
		expect(";")
		return node
	}

	if t := consume("if"); t != nil {
		node := newNode(ND_IF, t)
		node.Cond = expr()
		node.Then = stmt()
		if consume("else") != nil {
			node.Els = stmt()
		}
		return node
	}

	if t := consume("switch"); t != nil {
		node := newNode(ND_SWITCH, t)
		node.Cond = expr()

		sw := curSwitch
		curSwitch = node
		node.Then = stmt()
		curSwitch = sw
		return node
	}

	if t := consume("case"); t != nil {
		if curSwitch == nil {
			panic("\n" + errorTok(t, "stray case"))
		}
		val := constExpr()
		expect(":")

		node := newUnary(ND_CASE, stmt(), t)
		node.Val = val
		node.CaseNext = curSwitch.CaseNext
		curSwitch.CaseNext = node
		return node
	}

	if t := consume("default"); t != nil {
		if curSwitch == nil {
			panic("\n" + errorTok(t, "stray default"))
		}
		expect(":")
		node := newUnary(ND_CASE, stmt(), t)
		curSwitch.DefCase = node
		return node
	}

	if t := consume("for"); t != nil {
		if !isForClause() { // for for-stmt
			node := newNode(ND_WHILE, t)
			if peek("{") == nil {
				node.Cond = expr()
			} else {
				node.Cond = newNum(1, t)
			}

			node.Then = stmt()
			return node

		} else { // for for-clause
			node := newNode(ND_FOR, t)
			if consume(";") == nil {
				node.Init = readExprStmt()
				expect(";")
			}
			if consume(";") == nil {
				node.Cond = expr()
				expect(";")
			}
			if peek("{") == nil {
				node.Inc = readExprStmt()
			}
			node.Then = stmt()
			return node
		}
	}

	if t := consume("{"); t != nil {

		head := Node{}
		cur := &head

		sc := varScope
		for consume("}") == nil {
			cur.Next = stmt()
			cur = cur.Next
		}
		varScope = sc

		consume(";")
		return &Node{Kind: ND_BLOCK, Body: head.Next, Tok: t}
	}

	if t := consume("break"); t != nil {
		expect(";")
		return newNode(ND_BREAK, t)
	}

	if t := consume("continue"); t != nil {
		expect(";")
		return newNode(ND_CONTINUE, t)
	}

	if t := consume("goto"); t != nil {
		node := newNode(ND_GOTO, t)
		node.LblName = expectIdent()
		expect(";")
		return node
	}

	if t := consumeIdent(); t != nil {
		if consume(":") != nil {
			node := newUnary(ND_LABEL, stmt(), t)
			node.LblName = t.Str
			return node
		}
		token = t
	}

	if peek("var") != nil {
		return declaration()
	}

	if t := consume("type"); t != nil {
		name := expectIdent()
		ty := readTypePreffix()
		expect(";")
		pushScope(name).TyDef = ty
		return newNode(ND_NULL, t)
	}

	node := readExprStmt()
	expect(";")
	return node
}

// expr       = assign
func expr() *Node {
	// printCurTok()
	// printCalledFunc()

	return assign()
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
		if eval(node.Lhs) == eval(node.Rhs) {
			return 1
		}
		return 0
	case ND_NE:
		if eval(node.Lhs) != eval(node.Rhs) {
			return 1
		}
		return 0
	case ND_LT:
		if eval(node.Lhs) < eval(node.Rhs) {
			return 1
		}
		return 0
	case ND_LE:
		if eval(node.Lhs) <= eval(node.Rhs) {
			return 1
		}
		return 0
	case ND_NOT:
		if eval(node.Lhs) == 0 {
			return 1
		}
		return 0
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

// const-expr
func constExpr() int64 {
	return eval(logor())
}

// assign = logor (assign-op assign)?
// assign-op = "=" | "+=" | "-=" | "*=" | "/=" | "<<=" | ">>="
func assign() *Node {
	// printCurTok()
	// printCalledFunc()

	node := logor()
	if t := consume("="); t != nil {
		if token.Kind == TK_STR && node.Var.Ty.Kind == TY_ARRAY &&
			node.Var.Ty.Base.Kind == TY_BYTE {
			tok := token
			token = token.Next
			head := &Node{}
			stringAssign(head, node.Var, node.Var.Ty, nil, tok)
			n := newNode(ND_BLOCK, tok)
			n.Body = head.Next
			return n
		} else {
			node = newBinary(ND_ASSIGN, node, assign(), t)
		}
	} else if t := consume("+="); t != nil {
		node = newBinary(ND_A_ADD, node, assign(), t)
	} else if t := consume("-="); t != nil {
		node = newBinary(ND_A_SUB, node, assign(), t)
	} else if t := consume("*="); t != nil {
		node = newBinary(ND_A_MUL, node, assign(), t)
	} else if t := consume("/="); t != nil {
		node = newBinary(ND_A_DIV, node, assign(), t)
	} else if t := consume("<<="); t != nil {
		node = newBinary(ND_A_SHL, node, assign(), t)
	} else if t := consume(">>="); t != nil {
		node = newBinary(ND_A_SHR, node, assign(), t)
	}
	return node
}

// logor = logand ("||" logand)*
func logor() *Node {
	node := logand()
	t := consume("||")
	for t != nil {
		node = newBinary(ND_LOGOR, node, logand(), t)
		t = consume("||")
	}
	return node
}

// logand = bitor ("&&" bitor)*
func logand() *Node {
	node := bitor()
	t := consume("&&")
	for t != nil {
		node = newBinary(ND_LOGAND, node, bitor(), t)
		t = consume("&&")
	}
	return node
}

// bitor = bitxor ("|" bitxor)*
func bitor() *Node {
	node := bitxor()
	t := consume("|")
	for t != nil {
		node = newBinary(ND_BITOR, node, bitxor(), t)
		t = consume("|")
	}
	return node
}

// bitxor = bitand ("^" bitand)*
func bitxor() *Node {
	node := bitand()
	t := consume("^")
	for t != nil {
		node = newBinary(ND_BITXOR, node, bitxor(), t)
		t = consume("^")
	}
	return node
}

// bitand = equality ("&" equality)*
func bitand() *Node {
	node := equality()
	t := consume("&")
	for t != nil {
		node = newBinary(ND_BITAND, node, equality(), t)
		t = consume("&")
	}
	return node
}

// equality   = relational ("==" relational | "!=" relational)*
func equality() *Node {
	// printCurTok()
	// printCalledFunc()

	node := relational()

	for {
		if t := consume("=="); t != nil {
			node = newBinary(ND_EQ, node, relational(), t)
		} else if t := consume("!="); t != nil {
			node = newBinary(ND_NE, node, relational(), t)
		} else {
			return node
		}
	}
}

// relational = shift ("<" shift | "<=" shift | ">" shift | ">=" shift)*
func relational() *Node {
	// printCurTok()
	// printCalledFunc()

	node := shift()

	for {
		if t := consume("<"); t != nil {
			node = newBinary(ND_LT, node, shift(), t)
		} else if t := consume("<="); t != nil {
			node = newBinary(ND_LE, node, shift(), t)
		} else if t := consume(">"); t != nil {
			node = newBinary(ND_LT, shift(), node, t)
		} else if t := consume(">="); t != nil {
			node = newBinary(ND_LE, shift(), node, t)
		} else {
			return node
		}
	}
}

// shift = add ("<<" add | ">>" add)*
func shift() *Node {
	node := add()

	for {
		if t := consume("<<"); t != nil {
			node = newBinary(ND_SHL, node, add(), t)
		} else if t := consume(">>"); t != nil {
			node = newBinary(ND_SHR, node, add(), t)
		} else {
			return node
		}
	}
}

// add        = mul ("+" mul | "-" mul)*
func add() *Node {
	// printCurTok()
	// printCalledFunc()

	node := mul()

	for {
		if t := consume("+"); t != nil {
			node = newBinary(ND_ADD, node, mul(), t)
		} else if t := consume("-"); t != nil {
			node = newBinary(ND_SUB, node, mul(), t)
		} else {
			return node
		}
	}
}

// mul = cast ("*" cast | "/" cast)*
func mul() *Node {
	// printCurTok()
	// printCalledFunc()

	node := cast()

	for {
		if t := consume("*"); t != nil {
			node = newBinary(ND_MUL, node, cast(), t)
		} else if t := consume("/"); t != nil {
			node = newBinary(ND_DIV, node, cast(), t)
		} else {
			return node
		}
	}
}

// cast = type-name "(" cast ")" | unary
func cast() *Node {

	if isTypename() {
		ty := readTypePreffix()
		expect("(")
		node := newUnary(ND_CAST, cast(), token)
		node.Ty = ty
		expect(")")
		return node
	}

	return unary()
}

// unary   = ("+" | "-" | "*" | "&" | "!")? cast
//         | "sizeof" unary
//         | postfix
func unary() *Node {
	// printCurTok()
	// printCalledFunc()

	if t := consumeSizeof(); t != nil {
		return newUnary(ND_SIZEOF, cast(), t)
	}
	if t := consume("+"); t != nil {
		return cast()
	}
	if t := consume("-"); t != nil {
		return newBinary(ND_SUB, newNum(0, t), cast(), t)
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
	if t := consume("^"); t != nil {
		return newUnary(ND_BITNOT, cast(), t)
	}
	return postfix()
}

// postfix = primary ("[" expr "]" | "." ident | "++" | "--")*
func postfix() *Node {
	// printCurTok()
	// printCalledFunc()

	node := primary()

	for {
		if t := consume("["); t != nil {
			// x[y] is short for *(x+y)
			exp := newBinary(ND_ADD, node, expr(), t)
			expect("]")
			node = newUnary(ND_DEREF, exp, t)
			continue
		}

		if t := consume("."); t != nil {
			node = newUnary(ND_MEMBER, node, t)
			node.MemName = expectIdent()
			continue
		}

		if t := consume("++"); t != nil {
			node = newUnary(ND_INC, node, t)
			continue
		}

		if t := consume("--"); t != nil {
			node = newUnary(ND_DEC, node, t)
			continue
		}

		return node
	}
}

// func-args = "(" (assign ("," assign)*)? ")"
func funcArgs() *Node {
	// printCurTok()
	// printCalledFunc()

	if consume(")") != nil {
		return nil
	}

	head := assign()
	cur := head

	for consume(",") != nil {
		cur.Next = assign()
		cur = cur.Next
	}
	expect(")")
	return head
}

// primary = "(" expr ")" | ident args? | num
// args = "(" ")"
func primary() *Node {
	// printCurTok()
	// printCalledFunc()

	// if the next token is '(', the program must be
	// "(" expr ")"
	if consume("(") != nil {
		node := expr()
		expect(")")
		return node
	}

	if t := consumeIdent(); t != nil {
		if consume("(") != nil {
			node := &Node{
				Kind:     ND_FUNCALL,
				Tok:      t,
				FuncName: t.Str,
				Args:     funcArgs(),
			}

			sc := findVar(t)
			if sc != nil {
				if sc.Var == nil || sc.Var.Ty.Kind != TY_FUNC {
					panic("\n" + errorTok(t, "not a function"))
				}
				node.Ty = sc.Var.Ty.RetTy
			} else {
				node.Ty = intType()
			}
			return node
		}

		sc := findVar(t)
		if sc != nil && sc.Var != nil {
			return newVar(sc.Var, t)
		}
		panic("\n" + errorTok(t, "undifined variable"))
	}

	t := token
	if t.Kind == TK_STR {
		token = token.Next

		ty := arrayOf(charType(), t.ContLen)
		v := pushVar(newLabel(), ty, false, nil)
		v.Init = gvarInitString(t.Contents, t.ContLen)
		return newVar(v, t)
	}

	if t.Kind != TK_NUM {
		panic("\n" + errorTok(t, "expected expression: %s", t.Str))
	}
	return newNum(expectNumber(), t)
}
