package main

import (
	"fmt"
	"strconv"
)

// Scope for local variables, global variables or typedefs
type VarScope struct {
	Next  *VarScope
	Name  string
	Obj   *Obj
	TyDef *Type
}

// Scope for struct tags.
type TagScope struct {
	Next *TagScope
	Name string
	Ty   *Type
}

type Scope struct {
	Next *Scope

	Vars *VarScope
	Tags *TagScope
}

// Variable attributes typedef.
type VarAttr struct {
	IsTydef bool
}

// Variable or function
type Obj struct {
	Next    *Obj
	Name    string // Variable name
	Ty      *Type  // Type
	Tok     *Token // for error message
	IsLocal bool   // local or global

	// Local variables
	Offset int // Offset from RBP

	// Global variable or function
	IsFunc bool
	IsDef  bool

	// Global variables
	InitData string
	Rel      *Relocation

	// Function
	Params  *Obj
	Body    *Node
	Locals  *Obj
	StackSz int
}

type VarList struct {
	Next *VarList
	Obj  *Obj
}

type NodeKind int

const (
	ND_NULL_EXPR NodeKind = iota // Do nothing
	ND_ADD                       // +
	ND_SUB                       // -
	ND_MUL                       // *
	ND_DIV                       // /
	ND_NEG                       // unary -
	ND_MOD                       // %
	ND_BITAND                    // &
	ND_BITOR                     // |
	ND_BITXOR                    // ^
	ND_SHL                       // <<
	ND_SHR                       // >>
	ND_EQ                        // ==
	ND_NE                        // !=
	ND_LT                        // <
	ND_LE                        // <=
	ND_ASSIGN                    // =
	ND_COND                      // ?:
	ND_COMMA                     //
	ND_MEMBER                    // . (struct member access)
	ND_ADDR                      // unary &
	ND_DEREF                     // unary *
	ND_NOT                       // !
	ND_BITNOT                    // ~
	ND_LOGAND                    // &&
	ND_LOGOR                     // ||
	ND_RETURN                    // "return"
	ND_IF                        // "if"
	ND_FOR                       // "for" or "while"
	ND_SWITCH                    // "switch"
	ND_CASE                      // "case"
	ND_BLOCK                     // { ... }
	ND_GOTO                      // "goto"
	ND_LABEL                     // Labeled statement
	ND_FUNCALL                   // Function call
	ND_EXPR_STMT                 // Expression statement
	ND_STMT_EXPR                 // Statement expression
	ND_VAR                       // Variable
	ND_NUM                       // Integer
	ND_CAST                      // Type cast
	ND_MEMZERO                   // Zero-clear a stack variable
	ND_SIZEOF                    // 'Sizeof'
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

	// "break" and "continue" labels
	BrkLabel  string
	ContLabel string

	// Block
	Body *Node

	// Struct member access
	MemName string
	Mem     *Member

	// Function call
	FuncName string
	FuncTy   *Type
	Args     *Node

	// Goto or labeled statement
	Lbl       string
	UniqueLbl string
	GotoNext  *Node

	// Switch-cases
	CaseNext   *Node
	DefCase    *Node
	CaseLbl    int
	CaseEndLbl int

	Obj *Obj  // used if kind == ND_VAR
	Val int64 // used if kind == ND_NUM
}

var locals *Obj
var globals *Obj

var varScope *VarScope

var scope *Scope = &Scope{}

// Points to the function object the parser is currently parsing.
var curFn *Obj

// Lists of all goto statements and labels in the current function.
var gotos *Node
var labels *Node

// Current "goto" and "continue" jump targets.
var brkLabel string
var contLabel string

// Points to a node representing a switch if we are parsing
// a switch statement. Otherwise, nil
var curSwitch *Node

func enterScope() {
	sc := &Scope{Next: scope}
	scope = sc
}

func leaveScope() {
	scope = scope.Next
}

// findVar finds a variable or a typedef by name.
func findVar(tok *Token) *VarScope {
	for sc := scope; sc != nil; sc = sc.Next {
		for sc2 := sc.Vars; sc2 != nil; sc2 = sc2.Next {
			if equal(tok, sc2.Name) {
				return sc2
			}
		}
	}
	return nil
}

func findTag(tok *Token) *Type {
	for sc := scope; sc != nil; sc = sc.Next {
		for sc2 := sc.Tags; sc2 != nil; sc2 = sc2.Next {
			if equal(tok, sc2.Name) {
				return sc2.Ty
			}
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

func newLong(val int64, tok *Token) *Node {
	return &Node{
		Kind: ND_NUM,
		Tok:  tok,
		Val:  val,
		Ty:   ty_long,
	}
}

func newVarNode(v *Obj, tok *Token) *Node {
	return &Node{Kind: ND_VAR, Tok: tok, Obj: v}
}

func newCast(expr *Node, ty *Type) *Node {
	addType(expr)

	return &Node{
		Kind: ND_CAST,
		Tok:  expr.Tok,
		Lhs:  expr,
		Ty:   copyType(ty),
	}
}

func pushScope(name string) *VarScope {
	sc := &VarScope{Name: name, Next: scope.Vars}
	scope.Vars = sc
	return sc
}

// Global variable can be initialized either by a constant expression
// or a pointer to another global variable. This struct represents the
// latter.
type Relocation struct {
	Next   *Relocation
	Offset int
	Lbl    string
	Addend int64
}

// This struct represents a variable initializer. Since initializers
// can be nested (e.g. `var x [2][2]int = [2][2]int{[2]int{1,2},[2]int{3,4}}`),
// this struct is a tree data structure.
type Initializer struct {
	Next *Initializer
	Ty   *Type
	Tok  *Token

	// Constant expression
	Sz  int
	Val int64

	// Reference to another global variable
	Lbl string

	// If it's not an aggregate type and an initializer,
	// `expr` has an initialization expression.
	Expr *Node

	// If it's an initializer for an aggregete type (e.g. array or struct),
	// `children` has initializers for its children.
	Children []*Initializer
}

// For local variable initializer.
type InitDesg struct {
	Next *InitDesg
	Idx  int
	Mem  *Member
	Var  *Obj
}

func newInitializer(ty *Type) *Initializer {
	init := &Initializer{Ty: ty}

	if ty.Kind == TY_ARRAY {
		init.Children = make([]*Initializer, ty.ArrSz)
		for i := 0; i < ty.ArrSz; i++ {
			init.Children[i] = newInitializer(ty.Base)
		}
		return init
	}

	if ty.Kind == TY_STRUCT {
		// Count the number of struct members
		var l int
		for mem := ty.Mems; mem != nil; mem = mem.Next {
			l++
		}

		init.Children = make([]*Initializer, l)

		for mem := ty.Mems; mem != nil; mem = mem.Next {
			init.Children[mem.Idx] = newInitializer(mem.Ty)
		}
		return init
	}

	return init
}

func newVar(name string, ty *Type) *Obj {
	// printCurTok()
	// printCalledFunc()

	v := &Obj{Name: name, Ty: ty}
	pushScope(name).Obj = v
	return v
}

func newLvar(name string, ty *Type) *Obj {
	v := newVar(name, ty)
	v.IsLocal = true
	v.Next = locals
	locals = v
	return v
}

func newGvar(name string, ty *Type) *Obj {
	v := newVar(name, ty)
	v.IsLocal = false
	v.Next = globals
	globals = v
	return v
}

// for newLabel function
var cnt int

func newUniqueName() string {
	res := fmt.Sprintf(".L..%d", cnt)
	cnt++
	return res
}

func newAnonGvar(ty *Type) *Obj {
	return newGvar(newUniqueName(), ty)
}

func newStringLiteral(p string, ty *Type) *Obj {
	v := newAnonGvar(ty)
	v.InitData = p
	return v
}

func getIdent(tok *Token) string {
	if tok.Kind != TK_IDENT {
		errorTok(tok, "expected an identifier")
	}
	return string(strNdUp(tok.Contents, tok.Len))
}

func findTyDef(tok *Token) *Type {
	if tok.Kind == TK_IDENT {
		if sc := findVar(token); sc != nil {
			return sc.TyDef
		}
	}
	return nil
}

func pushTagScope(tok *Token, ty *Type) {
	sc := &TagScope{
		Name: string(strNdUp(tok.Contents, tok.Len)),
		Ty:   ty,
		Next: scope.Tags,
	}
	scope.Tags = sc
}

// typeSpecifier returns a pointer of Type struct.
// If the current tokens represents a typename,
// it returns the Type struct with that typename.
// Otherwise returns the Type struct with TY_VOID.
//
// type-specifier = "*"* builtin-type | struct-decl | typedef-name |
// builtin-type = void | "bool" | "byte"| "int16" | "int" | "int64" |
//                "string"
//
func declSpec(rest **Token, tok *Token) *Type {
	// printCurTok()
	// printCalledFunc()

	nPtr := 0
	for consume(&tok, tok, "*") {
		nPtr++
	}

	var ty *Type
	if equal(tok, "byte") {
		ty = ty_char
	} else if equal(tok, "string") {
		ty = stringType()
	} else if equal(tok, "bool") {
		ty = ty_bool
	} else if equal(tok, "int16") {
		ty = ty_short
	} else if equal(tok, "int") {
		ty = ty_int
	} else if equal(tok, "int64") {
		ty = ty_long
	} else if equal(tok, "struct") { // struct type
		ty = structDecl(&tok, tok.Next)
	}

	// Handle user-defined types.
	ty2 := findTyDef(tok)
	if ty2 != nil {
		ty = ty2
	}

	if ty == nil {
		ty = ty_void
	}

	for i := 0; i < nPtr; i++ {
		ty = pointerTo(ty)
	}

	*rest = tok
	return ty
}

func findBase(rest **Token, tok *Token) (*Type, *Token) {
	// printCurTok()
	// printCalledFunc()

	for !equal(tok, "*") && !isTypename(tok) {
		tok = tok.Next
	}
	ty := declSpec(rest, tok)
	t := tok // どこまでtokenを読んだか
	return ty, t
}

func readArr(tok *Token, base *Type) *Type {
	// printCurTok()
	// printCalledFunc()

	if !consume(&tok, tok, "[") {
		return base
	}
	var sz int64
	if !consume(&tok, tok, "]") {
		sz = constExpr(&tok, tok)
		tok = skip(tok, "]")
	}
	base = readArr(tok, base)
	return arrayOf(base, int(sz))
}

// type-preffix = ("[" const-expr "]")*
func readTypePreffix(rest **Token, tok *Token) *Type {
	// printCurTok()
	// printCalledFunc()

	if !equal(token, "[") {
		return declSpec(&tok, tok)
	}

	base, t := findBase(&tok, tok)
	arrTy := readArr(tok, base)
	*rest = t

	return arrTy
}

// declarator = ident (type-preffix)? declspec
func declarator(rest **Token, tok *Token) *Type {
	if tok.Kind != TK_IDENT {
		panic("\n" + errorTok(tok, "expected a variable name"))
	}
	name := tok
	ty := readTypePreffix(rest, tok.Next)
	ty.Name = name
	return ty
}

// param = declarator
// e.g.
//  x int
//  x *int
//  x **int
//  x [3]int
//  x [3]*int
//  x [2]**int
// params = param ("," param)*
func funcParams(rest **Token, tok *Token, ty *Type) *Type {
	// printCurTok()
	// printCalledFunc()

	head := &Type{}
	cur := head

	for !equal(tok, ")") {
		if cur != head {
			tok = skip(tok, ",")
		}
		ty2 := declarator(&tok, tok)

		// "array of T" is converted tot "pointer to T" only in the parameter
		// context. For example, *argv[] is converted to **argv by this.
		if ty2.Kind == TY_ARRAY {
			name := ty2.Name
			ty2 = pointerTo(ty2.Base)
			ty2.Name = name
		}

		cur.Next = copyType(ty2)
		cur = cur.Next
	}

	ty = funcType(ty)
	return ty
}

func consumeEnd(rest **Token, tok *Token) bool {
	if equal(tok, "}") {
		*rest = tok.Next
		return true
	}

	if equal(tok, ",") && equal(tok.Next, "}") {
		*rest = tok.Next.Next
		return true
	}

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

func skipExcessElement(tok *Token) *Token {
	if equal(tok, "{") {
		tok = skipExcessElement(tok.Next)
		return skip(tok, "}")
	}

	assign(&tok, tok)
	return tok
}

// string-initializer = string-literal
func stringInitializer(rest **Token, tok *Token, init *Initializer) {
	len := min(init.Ty.ArrSz, tok.Ty.ArrSz)
	for i := 0; i < len; i++ {
		init.Children[i].Expr = newNum(int64(tok.Str[i]), tok)
	}
	*rest = tok.Next
}

// array-initializer = (type-preffix)? decl-spec "{" initializer ("," initializer)* ","? "}"
func arrayInitializer1(rest **Token, tok *Token, init *Initializer) {
	tok = skip(tok, "{")

	for i := 0; !consumeEnd(rest, tok); i++ {
		if i > 0 {
			tok = skip(tok, ",")
		}

		if i < init.Ty.ArrSz {
			initializer2(&tok, tok, init.Children[i])
		} else {
			tok = skipExcessElement(tok)
		}
	}
}

// struct-initializer = "{" initializer ("," initializer)* ","? "}"
func structInitializer1(rest **Token, tok *Token, init *Initializer) {
	tok = skip(tok, "{")

	mem := init.Ty.Mems

	for !consumeEnd(rest, tok) {
		if mem != init.Ty.Mems {
			tok = skip(tok, ",")
		}

		if mem != nil {
			initializer2(&tok, tok, init.Children[mem.Idx])
			mem = mem.Next
		} else {
			tok = skipExcessElement(tok)
		}
	}
}

// initializer = string-initializer | array-initializer
//             | struct-initializer
//             | assign
func initializer2(rest **Token, tok *Token, init *Initializer) {
	if init.Ty.Kind == TY_ARRAY && tok.Kind == TK_STR {
		stringInitializer(rest, tok, init)
		return
	}

	if init.Ty.Kind == TY_ARRAY {
		readTypePreffix(rest, tok) // discard the return value
		arrayInitializer1(rest, tok, init)
	}

	if init.Ty.Kind == TY_STRUCT {
		readTypePreffix(rest, tok) // discard the return value
		structInitializer1(rest, tok, init)
		return
	}
}

func initializer(rest **Token, tok *Token, ty *Type, newTy **Type) *Initializer {
	init := newInitializer(ty)
	initializer2(rest, tok, init)

	*newTy = init.Ty
	return init
}

func initDesgExpr(desg *InitDesg, tok *Token) *Node {
	if desg.Var != nil {
		return newVarNode(desg.Var, tok)
	}

	if desg.Mem != nil {
		node := newUnary(ND_MEMBER, initDesgExpr(desg.Next, tok), tok)
		node.Mem = desg.Mem
		return node
	}

	lhs := initDesgExpr(desg.Next, tok)
	rhs := newNum(int64(desg.Idx), tok)
	return newUnary(ND_DEREF, newAdd(lhs, rhs, tok), tok)
}

func createLvarInit(init *Initializer, ty *Type, desg *InitDesg, tok *Token) *Node {
	if ty.Kind == TY_ARRAY {
		node := newNode(ND_NULL_EXPR, tok)
		for i := 0; i < ty.ArrSz; i++ {
			desg2 := &InitDesg{Next: desg, Idx: i}
			rhs := createLvarInit(init.Children[i], ty.Base, desg2, tok)
			node = newBinary(ND_COMMA, node, rhs, tok)
		}
		return node
	}

	if ty.Kind == TY_STRUCT && init.Expr == nil {
		node := newNode(ND_NULL_EXPR, tok)

		for mem := ty.Mems; mem != nil; mem = mem.Next {
			desg2 := &InitDesg{Next: desg, Idx: 0, Mem: mem}
			rhs := createLvarInit(init.Children[mem.Idx], mem.Ty, desg2, tok)
			node = newBinary(ND_COMMA, node, rhs, tok)
		}
		return node
	}

	if init.Expr == nil {
		return newNode(ND_NULL_EXPR, tok)
	}

	lhs := initDesgExpr(desg, tok)
	return newBinary(ND_ASSIGN, lhs, init.Expr, tok)
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
func lvarInitializer(rest **Token, tok *Token, v *Obj) *Node {
	// Initialize a char array with a string literal.
	// => unnecessary

	init := initializer(rest, tok, v.Ty, &v.Ty)
	desg := &InitDesg{nil, 0, nil, v}

	// If a partial initializer list is given, the standard requires
	// that unspecified elements are set to 0. Here, we simply
	// zero-inilialize the entire memory region of a variable defore
	// initializing it with user-supplied values.
	lhs := newNode(ND_MEMZERO, tok)
	lhs.Obj = v

	rhs := createLvarInit(init, v.Ty, desg, tok)
	return newBinary(ND_COMMA, lhs, rhs, tok)
}

func writeBuf(buf *string, val int64, sz int) {
	switch sz {
	case 1:
		*buf = strconv.Itoa(int(int8(val)))
	case 2:
		*buf = strconv.Itoa(int(int16(val)))
	case 4:
		*buf = strconv.Itoa(int(val))
	case 8:
		*buf = strconv.FormatInt(val, 10)
	default:
		panic("internal error")
	}

}

func writeGvarData(cur *Relocation,
	init *Initializer, ty *Type, buf *string, offset int) *Relocation {
	if ty.Kind == TY_ARRAY {
		sz := ty.Base.Sz
		for i := 0; i < ty.ArrSz; i++ {
			cur = writeGvarData(cur, init.Children[i], ty.Base, buf, offset+sz*i)
		}
		return cur
	}

	if ty.Kind == TY_STRUCT {
		for mem := ty.Mems; mem != nil; mem = mem.Next {
			cur = writeGvarData(cur, init.Children[mem.Idx], mem.Ty, buf,
				offset+mem.Offset)
		}
		return cur
	}

	if init.Expr == nil {
		return cur
	}

	var label *string = nil
	val := eval2(init.Expr, label)

	if label == nil {
		writeBuf(buf, val, ty.Sz)
		return cur
	}

	rel := &Relocation{
		Offset: offset,
		Lbl:    *label,
		Addend: int64(val),
	}
	cur.Next = rel
	return cur.Next
}

func gvarInitializer(rest **Token, tok *Token, v *Obj) {
	// printCalledFunc()
	// printCurTok()

	init := initializer(rest, tok, v.Ty, &v.Ty)

	head := &Relocation{}
	var buf string
	writeGvarData(head, init, v.Ty, &buf, 0)
	v.InitData = buf
	v.Rel = head.Next
}

type Designator struct {
	Next *Designator
	Idx  int
	Mem  *Member
}

// Creates a node for an array sccess. For example, if v represents
// a variable x and desg represents indices 3 and 4, this function
// returns a node representing x[3][4]
func newDesgNode2(v *Obj, desg *Designator) *Node {
	tok := v.Tok
	if desg == nil {
		return newVarNode(v, tok)
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

func newDesgNode(v *Obj, desg *Designator, rhs *Node) *Node {
	lhs := newDesgNode2(v, desg)
	node := newBinary(ND_ASSIGN, lhs, rhs, rhs.Tok)
	return newUnary(ND_EXPR_STMT, node, rhs.Tok)
}

func lvarInitZero(cur *Node, v *Obj, ty *Type, desg *Designator) *Node {
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

func stringAssign(cur *Node, v *Obj, ty *Type, desg *Designator, tok *Token) *Node {
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

// declaration = VarDecl | VarSpec(unimplemented) | ShortVarDecl(unimplemented)
// VarDecl = "var" ident type-prefix declspec ("=" expr)
// VarSpec = ident-list (type-preffix type-specifier [ "=" expr-list ] | "=" expr-list)
// ShortVarDecl = "var" ident "=" expr => unimplemented
//              | ident ":=" expr => unimplemented
func declaration(rest **Token, tok *Token) *Node {
	// printCurTok()
	// printCalledFunc()

	head := &Node{}
	cur := head
	var i int

	for !equal(tok, ";") {
		if i > 0 {
			tok = skip(tok, ",")
		}
		i++
		ty := declarator(&tok, tok)
		if ty.Kind == TY_VOID {
			panic("\n" + errorTok(tok, "variable declared void"))
		}

		v := newLvar(getIdent(ty.Name), ty)
		if equal(tok, "=") {
			expr := lvarInitializer(&tok, tok.Next, v)
			cur.Next = newUnary(ND_EXPR_STMT, expr, tok)
			cur = cur.Next
		}

		if v.Ty.Sz < 0 {
			panic("\n" + errorTok(ty.Name, "variable has incomplete type"))
		}
		if v.Ty.Kind == TY_VOID {
			panic("\n" + errorTok(ty.Name, "variable declared void"))
		}
	}

	node := newNode(ND_BLOCK, tok)
	node.Body = head.Next
	*rest = tok.Next
	return node
}

func isTypename(tok *Token) bool {
	// printCurTok()
	// printCalledFunc()

	kw := []string{
		"byte", "bool", "int16", "int", "int64", "struct", "string",
	}

	for i := 0; i < len(kw); i++ {
		if equal(tok, kw[i]) {
			return true
		}
	}
	return findTyDef(token) != nil
}

func isForClause(tok *Token) bool {
	// printCurTok()
	// printCalledFunc()

	for !equal(tok, "{") {
		if equal(tok, ";") {
			return true
		}
		tok = tok.Next
	}
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
func stmt(rest **Token, tok *Token) *Node {
	// printCurTok()
	// printCalledFunc()

	if equal(tok, "return") {
		node := newNode(ND_RETURN, tok)
		exp := expr(&tok, tok.Next)
		*rest = skip(token, ";")

		addType(exp)
		node.Lhs = newCast(exp, curFn.Ty.RetTy)
		return node
	}

	if equal(tok, "if") {
		node := newNode(ND_IF, tok)
		node.Cond = expr(&tok, tok)
		node.Then = stmt(&tok, tok)
		if equal(tok, "else") {
			node.Els = stmt(&tok, tok.Next)
		}
		return node
	}

	if equal(tok, "switch") {
		node := newNode(ND_SWITCH, token)
		node.Cond = expr(&tok, tok)

		sw := curSwitch
		curSwitch = node
		node.Then = stmt(rest, tok)
		curSwitch = sw
		return node
	}

	if equal(tok, "case") {
		if curSwitch == nil {
			panic("\n" + errorTok(token, "stray case"))
		}
		node := newNode(ND_CASE, tok)
		val := constExpr(&tok, tok.Next)
		token = skip(token, ":")
		node.Lbl = newUniqueName()
		node.Lhs = stmt(rest, tok)
		node.Val = val
		node.CaseNext = curSwitch.CaseNext
		curSwitch.CaseNext = node
		return node
	}

	if equal(tok, "default") {
		if curSwitch == nil {
			panic("\n" + errorTok(token, "stray default"))
		}
		node := newNode(ND_CASE, tok)
		token = skip(token, ":")
		node.Lbl = newUniqueName()
		node.Lhs = stmt(rest, tok)
		curSwitch.DefCase = node
		return node
	}

	if equal(tok, "for") {
		if !isForClause(tok) { // for for-stmt like 'while' statement
			node := newNode(ND_FOR, tok)
			if !equal(token, "{") {
				node.Cond = expr(&tok, tok)
			} else {
				node.Cond = newNum(1, token)
			}
			tok = skip(tok, "{")

			brk := brkLabel
			cont := contLabel
			node.BrkLabel = newUniqueName()
			brkLabel = node.BrkLabel
			node.ContLabel = newUniqueName()
			contLabel = node.ContLabel

			node.Then = stmt(rest, tok)

			brkLabel = brk
			contLabel = cont
			return node

		} else { // for for-clause
			node := newNode(ND_FOR, tok)
			enterScope()
			brk := brkLabel
			cont := contLabel
			node.BrkLabel = newUniqueName()
			brkLabel = node.BrkLabel
			node.ContLabel = newUniqueName()
			contLabel = node.ContLabel

			if !equal(tok, ";") {
				node.Init = exprStmt(&tok, tok)
			}
			token = skip(token, ";")
			if !equal(tok, ";") {
				node.Cond = expr(&tok, tok)
			}
			token = skip(token, ";")
			if !equal(token, "{") {
				node.Inc = expr(&tok, tok)
			}
			token = skip(token, "{")

			node.Then = stmt(rest, tok)

			leaveScope()
			brkLabel = brk
			contLabel = cont
			return node
		}
	}

	if equal(tok, "goto") {
		node := newNode(ND_GOTO, tok)
		node.Lbl = getIdent(tok.Next)
		node.GotoNext = gotos
		gotos = node
		*rest = skip(tok.Next.Next, ";")
		return node
	}

	if consume(&token, token, "break") {
		if brkLabel == "" {
			panic("\n" + errorTok(tok, "stray break"))
		}
		node := newNode(ND_GOTO, tok)
		node.UniqueLbl = brkLabel
		*rest = skip(tok.Next, ";")
		return node
	}

	if equal(tok, "continue") {
		if contLabel == "" {
			panic("\n" + errorTok(tok, "stray continue"))
		}
		node := newNode(ND_GOTO, tok)
		node.UniqueLbl = contLabel
		*rest = skip(tok.Next, ";")
		return node
	}

	// Labeled statement
	if tok.Kind == TK_IDENT && equal(tok.Next, ":") {
		node := newNode(ND_LABEL, tok)
		node.Lbl = tok.Str
		node.UniqueLbl = newUniqueName()
		node.Lhs = stmt(rest, tok.Next.Next)
		node.GotoNext = labels
		labels = node
		return node
	}

	if equal(tok, "{") {
		return compoundStmt(rest, tok.Next)
	}

	return exprStmt(rest, tok)
}

// compound-stmt = (typedef | declaration | stmt)* "}"
func compoundStmt(rest **Token, tok *Token) *Node {
	node := newNode(ND_BLOCK, tok)
	head := &Node{}
	cur := head

	enterScope()

	for !equal(tok, "}") {
		if equal(token, "var") {
			cur.Next = declaration(&tok, tok)
			cur = cur.Next
			addType(cur)
			continue
		}

		if consume(&tok, tok, "type") {
			tok = parseTypedef(tok)
			tok = skip(tok, ";")
			continue
		}

		cur.Next = stmt(&tok, tok)
		cur = cur.Next
		addType(cur)
	}
	leaveScope()

	node.Body = head.Next
	*rest = tok.Next
	return node
}

// expr-stmt = expr? ";"
func exprStmt(rest **Token, tok *Token) *Node {
	if equal(tok, ";") {
		*rest = tok.Next
		return newNode(ND_BLOCK, tok)
	}

	node := newNode(ND_EXPR_STMT, tok)
	node.Lhs = expr(&tok, tok)
	*rest = skip(tok, ";")
	return node
}

// expr       = assign ("," assign)*
func expr(rest **Token, tok *Token) *Node {
	// printCurTok()
	// printCalledFunc()

	node := assign(&tok, tok)

	if equal(tok, ",") {
		return newBinary(ND_COMMA, node, expr(rest, tok.Next), tok)
	}
	return node
}

func eval(node *Node) int64 {
	return eval2(node, nil)
}

func eval2(node *Node, label *string) int64 {
	addType(node)

	switch node.Kind {
	case ND_ADD:
		return eval2(node.Lhs, label) + eval(node.Rhs)
	case ND_SUB:
		return eval2(node.Lhs, label) - eval(node.Rhs)
	case ND_MUL:
		return eval(node.Lhs) * eval(node.Rhs)
	case ND_DIV:
		return eval(node.Lhs) / eval(node.Rhs)
	case ND_NEG:
		return -eval(node.Lhs)
	case ND_MOD:
		return eval(node.Lhs) % eval(node.Rhs)
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
	case ND_COND:
		if eval(node.Cond) != 0 {
			return eval2(node.Then, label)
		}
		return eval2(node.Els, label)
	case ND_COMMA:
		return eval2(node.Lhs, label)
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
	case ND_CAST:
		val := eval2(node.Lhs, label)
		if isInteger(node.Ty) {
			switch node.Ty.Sz {
			case 1:
				return int64(uint8(val))
			case 2:
				return int64(uint16(val))
			case 4:
				return int64(uint32(val))
			}
		}
		return val
	case ND_ADDR:
		return evalRval(node.Lhs, label)
	case ND_MEMBER:
		if label == nil {
			panic("\n" + errorTok(node.Tok, "not a compile-time constant"))
		}
		if node.Obj.Ty.Kind != TY_ARRAY && node.Obj.Ty.Kind == TY_FUNC {
			panic("\n" + errorTok(node.Tok, "invalid initializer"))
		}
		*label = node.Obj.Name
		return 0
	case ND_NUM:
		return node.Val
	default:
		panic("\n" + errorTok(node.Tok, "not a compile-time constant"))
	}
}

func evalRval(node *Node, label *string) int64 {
	switch node.Kind {
	case ND_VAR:
		if node.Obj.IsLocal {
			panic("\n" + errorTok(node.Tok, "not a compile-time constant"))
		}
		*label = node.Obj.Name
		return 0
	case ND_DEREF:
		return eval2(node.Lhs, label)
	case ND_MEMBER:
		return evalRval(node.Lhs, label) + int64(node.Mem.Offset)
	default:
		panic("\n" + errorTok(node.Tok, "invalid initializer"))
	}
}

// const-expr
func constExpr(rest **Token, tok *Token) int64 {
	return eval(logor(rest, tok))
}

// Convert `A op= B` to `*tmp = *tmp op B`
// where tmp is a fresh pointer variable.
func toAssign(binary *Node) *Node {
	addType(binary.Lhs)
	addType(binary.Rhs)
	tok := binary.Tok

	v := newLvar("", pointerTo(binary.Lhs.Ty))

	expr1 := newBinary(ND_ASSIGN, newVarNode(v, tok),
		newUnary(ND_ADDR, binary.Lhs, tok), tok)

	expr2 := newBinary(ND_ASSIGN,
		newUnary(ND_DEREF, newVarNode(v, tok), tok),
		newBinary(binary.Kind,
			newUnary(ND_DEREF, newVarNode(v, tok), tok),
			binary.Rhs,
			tok),
		tok)

	return newBinary(ND_COMMA, expr1, expr2, tok)
}

// assign = logor (assign-op assign)?
// assign-op = "=" | "+=" | "-=" | "*=" | "/=" | "<<=" | ">>="
func assign(rest **Token, tok *Token) *Node {
	// printCurTok()
	// printCalledFunc()

	node := logor(&tok, tok)

	if equal(tok, "=") {
		return newBinary(ND_ASSIGN, node, assign(rest, tok.Next), tok)
	}

	if equal(tok, "+=") {
		return toAssign(newAdd(node, assign(rest, tok.Next), tok))
	}

	if equal(tok, "-=") {
		return toAssign(newSub(node, assign(rest, tok.Next), tok))
	}

	if equal(tok, "*=") {
		return toAssign(newBinary(ND_MUL, node, assign(rest, tok.Next), tok))
	}

	if equal(tok, "/=") {
		return toAssign(newBinary(ND_DIV, node, assign(rest, tok.Next), token))
	}

	if equal(tok, "%=") {
		return toAssign(newBinary(ND_BITAND, node, assign(rest, tok.Next), tok))
	}

	if equal(tok, "|=") {
		return toAssign(newBinary(ND_BITOR, node, assign(rest, tok.Next), tok))
	}

	if equal(tok, "^=") {
		return toAssign(newBinary(ND_BITXOR, node, assign(rest, tok.Next), tok))
	}

	if equal(tok, "<<=") {
		return toAssign(newBinary(ND_SHL, node, assign(rest, tok.Next), tok))
	}

	if equal(tok, ">>=") {
		return toAssign(newBinary(ND_SHR, node, assign(rest, tok.Next), tok))
	}

	*rest = tok
	return node
}

// logor = logand ("||" logand)*
func logor(rest **Token, tok *Token) *Node {
	node := logand(&tok, tok)
	for equal(tok, "||") {
		start := tok
		node = newBinary(ND_LOGOR, node, logand(&tok, tok.Next), start)
	}
	*rest = tok
	return node
}

// logand = bitor ("&&" bitor)*
func logand(rest **Token, tok *Token) *Node {
	node := bitor(&tok, tok)
	for equal(tok, "&&") {
		start := tok
		node = newBinary(ND_LOGAND, node, bitor(&tok, tok.Next), start)
	}
	*rest = tok
	return node
}

// bitor = bitxor ("|" bitxor)*
func bitor(rest **Token, tok *Token) *Node {
	node := bitxor(&tok, tok)
	for equal(tok, "|") {
		start := tok
		node = newBinary(ND_BITOR, node, bitxor(&tok, tok.Next), start)
	}
	*rest = tok
	return node
}

// bitxor = bitand ("^" bitand)*
func bitxor(rest **Token, tok *Token) *Node {
	node := bitand(&tok, tok)
	for equal(tok, "^") {
		start := tok
		node = newBinary(ND_BITXOR, node, bitand(&tok, tok), start)
	}
	*rest = tok
	return node
}

// bitand = equality ("&" equality)*
func bitand(rest **Token, tok *Token) *Node {
	node := equality(&tok, tok)
	for equal(tok, "&") {
		start := tok
		node = newBinary(ND_BITAND, node, equality(&tok, tok), start)
	}
	*rest = tok
	return node
}

// equality   = relational ("==" relational | "!=" relational)*
func equality(rest **Token, tok *Token) *Node {
	// printCurTok()
	// printCalledFunc()

	node := relational(&tok, tok)

	for {
		start := tok

		if equal(tok, "==") {
			node = newBinary(ND_EQ, node, relational(&tok, tok.Next), start)
			continue
		}

		if equal(tok, "!=") {
			node = newBinary(ND_NE, node, relational(&tok, tok), start)
			continue
		}

		*rest = tok
		return node
	}
}

// relational = shift ("<" shift | "<=" shift | ">" shift | ">=" shift)*
func relational(rest **Token, tok *Token) *Node {
	// printCurTok()
	// printCalledFunc()

	node := shift(&tok, tok)

	for {
		start := tok

		if equal(tok, "<") {
			node = newBinary(ND_LT, node, shift(&tok, tok.Next), start)
			continue
		}

		if equal(tok, "<=") {
			node = newBinary(ND_LE, node, shift(&tok, tok.Next), start)
			continue
		}

		if equal(tok, ">") {
			node = newBinary(ND_LT, shift(&tok, tok), node, start)
			continue
		}

		if equal(tok, ">=") {
			node = newBinary(ND_LE, shift(&tok, tok), node, start)
			continue
		}

		*rest = tok
		return node
	}
}

// shift = add ("<<" add | ">>" add)*
func shift(rest **Token, tok *Token) *Node {
	node := add(&tok, tok)

	for {
		start := tok

		if consume(&token, token, "<<") {
			node = newBinary(ND_SHL, node, add(&tok, tok.Next), start)
			continue
		}

		if consume(&token, token, ">>") {
			node = newBinary(ND_SHR, node, add(&tok, tok), start)
			continue
		}

		*rest = tok
		return node
	}
}

// `+` operator is overloaded to perform the pointer arithmetic.
// If p is a pointer, p+n add not n but sizeof(*p)*n to the value of p,
// sothat p+n pointes to the location n elements (not bytes) ahead of p.
// In other words, we need to scale an integer value before adding to a
// pointer value. This function takes care of the scaling.
// => that isn't supported in Go.
func newAdd(lhs, rhs *Node, tok *Token) *Node {
	addType(lhs)
	addType(rhs)

	// num + num
	if isInteger(lhs.Ty) && isInteger(rhs.Ty) {
		return newBinary(ND_ADD, lhs, rhs, tok)
	}

	if lhs.Ty.Base != nil && rhs.Ty.Base != nil {
		panic("\n" + errorTok(tok, "invalid operands"))
	}

	// Canonicalize `num + ptr` to `ptr + num`.
	if lhs.Ty.Base == nil && rhs.Ty.Base != nil {
		tmp := lhs
		lhs = rhs
		rhs = tmp
	}

	// ptr + num
	rhs = newBinary(ND_MUL, rhs, newLong(int64(lhs.Ty.Base.Sz), tok), tok)
	return newBinary(ND_ADD, lhs, rhs, tok)
}

// Like `+`, `-` is overloaded for the pointer type.
// => that isn't supported in Go.
func newSub(lhs, rhs *Node, tok *Token) *Node {
	addType(lhs)
	addType(rhs)

	// num - num
	if isInteger(lhs.Ty) && isInteger(rhs.Ty) {
		return newBinary(ND_SUB, lhs, rhs, tok)
	}

	// ptr - num
	if lhs.Ty.Base != nil && isInteger(rhs.Ty) {
		rhs = newBinary(ND_MUL, rhs, newLong(int64(lhs.Ty.Base.Sz), tok), tok)
		addType(rhs)
		node := newBinary(ND_SUB, lhs, rhs, tok)
		node.Ty = lhs.Ty
		return node
	}

	// ptr - ptr, which returns how many elements are between the two.
	if lhs.Ty.Base != nil && rhs.Ty.Base != nil {
		node := newBinary(ND_SUB, lhs, rhs, tok)
		node.Ty = ty_int
		return newBinary(ND_DIV, node, newNum(int64(lhs.Ty.Base.Sz), tok), tok)
	}

	panic("\n" + errorTok(tok, "invalud operands"))
}

// add        = mul ("+" mul | "-" mul)*
func add(rest **Token, tok *Token) *Node {
	// printCurTok()
	// printCalledFunc()

	node := mul(&tok, tok)

	for {
		start := tok

		if equal(tok, "+") {
			node = newAdd(node, mul(&tok, tok.Next), start)
			continue
		}

		if equal(tok, "-") {
			node = newSub(node, mul(&tok, tok), start)
			continue
		}

		*rest = tok
		return node
	}
}

// mul = cast ("*" cast | "/" cast)*
func mul(rest **Token, tok *Token) *Node {
	// printCurTok()
	// printCalledFunc()

	node := cast(&tok, tok)

	for {
		start := tok

		if equal(tok, "*") {
			node = newBinary(ND_MUL, node, cast(&tok, tok.Next), start)
			continue
		}

		if equal(tok, "/") {
			node = newBinary(ND_DIV, node, cast(&tok, tok), start)
			continue
		}

		if equal(tok, "%") {
			node = newBinary(ND_MOD, node, cast(&tok, tok), start)
		}

		*rest = tok
		return node
	}
}

// cast = type-name "(" cast ")" | unary
func cast(rest **Token, tok *Token) *Node {

	if isTypename(tok) {
		start := tok
		ty := readTypePreffix(&tok, tok)
		token = skip(token, "(")
		node := newCast(cast(rest, tok), ty)
		node.Tok = start
		tok = skip(tok, ")")
		return node
	}

	return unary(rest, tok)
}

// unary   = ("+" | "-" | "*" | "&" | "!")? cast
//         | "Sizeof" unary
//         | postfix
func unary(rest **Token, tok *Token) *Node {
	// printCurTok()
	// printCalledFunc()

	if equal(tok, "+") {
		return cast(rest, tok.Next)
	}

	if equal(tok, "-") {
		return newUnary(ND_NEG, cast(rest, tok.Next), tok)
	}

	if equal(tok, "&") {
		return newUnary(ND_ADDR, cast(rest, tok.Next), tok)
	}

	if equal(tok, "*") {
		return newUnary(ND_DEREF, cast(rest, tok.Next), tok)
	}

	if equal(tok, "!") {
		return newUnary(ND_NOT, cast(rest, tok.Next), tok)
	}

	if equal(tok, "^") {
		return newUnary(ND_BITNOT, cast(rest, tok.Next), tok)
	}

	return postfix(rest, tok)
}

// struct-member = ident type-prefix type-specifier
func structMems(rest **Token, tok *Token, ty *Type) *Member {
	// printCurTok()
	// printCalledFunc()

	head := &Member{}
	cur := head
	idx := 0

	for !equal(tok, "}") {
		first := true
		for !consume(&tok, tok, ";") {
			if !first {
				tok = skip(tok, ",")
			}
			first = false

			name := getIdent(tok)
			memTy := readTypePreffix(&tok, tok)
			mem := &Member{
				Name: name,
				Ty:   memTy,
				Idx:  idx,
			}
			idx++
			cur.Next = mem
			cur = cur.Next
		}
	}

	*rest = tok.Next
	return head.Next
}

// struct-decl = "struct" "{" struct-member "}"
func structDecl(rest **Token, tok *Token) *Type {
	// printCurTok()
	// printCalledFunc()

	tok = skip(tok, "{")

	// Construct a struct object.
	ty := structType()
	ty.Mems = structMems(rest, tok, ty)

	// Assign offsers within the struct to members.
	offset := 0
	for mem := ty.Mems; mem != nil; mem = mem.Next {
		offset = alignTo(offset, mem.Ty.Align)
		mem.Offset = offset
		offset += mem.Ty.Sz

		if ty.Align < mem.Ty.Align {
			ty.Align = mem.Ty.Align
		}
	}
	ty.Sz = alignTo(offset, ty.Align)
	return ty
}

func getStructMember(ty *Type, tok *Token) *Member {
	for mem := ty.Mems; mem != nil; mem = mem.Next {
		if mem.Name == tok.Str {
			return mem
		}
	}
	panic("\n" + errorTok(tok, "no such member"))
}

func structRef(lhs *Node, tok *Token) *Node {
	addType(lhs)
	if lhs.Ty.Kind != TY_STRUCT {
		panic("\n" + errorTok(lhs.Tok, "not a struct"))
	}

	node := newUnary(ND_MEMBER, lhs, tok)
	node.Mem = getStructMember(lhs.Ty, tok)
	return node
}

// Convert A++ to `(typeof A)((A += 1) - 1)`
func newIncDec(node *Node, tok *Token, addend int) *Node {
	addType(node)
	return newCast(newAdd(toAssign(newAdd(node, newNum(int64(addend), tok), tok)),
		newNum(int64(addend)*-1, tok), tok),
		node.Ty)
}

// postfix = primary ("[" expr "]" | "." ident | "++" | "--")*
func postfix(rest **Token, tok *Token) *Node {
	// printCurTok()
	// printCalledFunc()

	node := primary(&tok, tok)

	for {
		if equal(tok, "[") {
			// x[y] is short for *(x+y)
			start := tok
			idx := expr(&tok, tok)
			tok = skip(tok, "]")
			node = newUnary(ND_DEREF, newAdd(node, idx, start), start)
			continue
		}

		if equal(tok, ".") {
			node = structRef(node, tok.Next)
			tok = tok.Next.Next
			continue
		}

		if equal(tok, "++") {
			node = newIncDec(node, tok, 1)
			tok = tok.Next
			continue
		}

		if equal(tok, "--") {
			node = newIncDec(node, tok, -1)
			tok = tok.Next
			continue
		}

		*rest = tok
		return node
	}
}

// funcall = ident "(" (assign ("," assign)*)? ")"
func funcall(rest **Token, tok *Token) *Node {
	start := tok
	tok = tok.Next.Next

	sc := findVar(start)
	if sc == nil {
		panic("\n" + errorTok(start, "implicit declaration of a function"))
	}
	if sc.Obj == nil || sc.Obj.Ty.Kind != TY_FUNC {
		panic("\n" + errorTok(start, "not a function"))
	}

	ty := sc.Obj.Ty
	paramTy := ty.Params

	head := &Node{}
	cur := head

	for !equal(tok, ")") {
		if cur != head {
			tok = skip(tok, ",")
		}

		arg := assign(&tok, tok)
		addType(arg)

		if paramTy != nil {
			if paramTy.Kind == TY_STRUCT {
				panic("\n" + errorTok(arg.Tok, "passing struct is not supported yet"))
			}
			arg = newCast(arg, paramTy)
			paramTy = paramTy.Next
		}

		cur.Next = arg
		cur = cur.Next
	}

	*rest = skip(tok, ")")

	node := newNode(ND_FUNCALL, start)
	node.FuncName = start.Str
	node.FuncTy = ty
	node.Ty = ty.RetTy
	node.Args = head.Next
	return node
}

// primary = "(" expr ")" | ident args? | num
// args = "(" ")"
func primary(rest **Token, tok *Token) *Node {
	// printCurTok()
	// printCalledFunc()

	// if the next token is '(', the program must be
	// "(" expr ")"
	if equal(tok, "(") {
		node := expr(&tok, tok.Next)
		*rest = skip(tok, ")")
		return node
	}

	if tok.Kind == TK_IDENT {
		// Function call
		if equal(tok, "(") {
			return funcall(rest, tok)
		}

		sc := findVar(tok)
		if sc == nil {
			panic("\n" + errorTok(tok, "undefined variable"))
		}

		var node *Node
		if sc.Obj != nil {
			node = newVarNode(sc.Obj, tok)
		}

		*rest = tok.Next
		return node
	}

	if tok.Kind == TK_STR {
		v := newStringLiteral(tok.Str, tok.Ty)
		*rest = tok.Next
		return newVarNode(v, tok)
	}

	if tok.Kind != TK_NUM {
		node := newNum(tok.Val, tok)
		*rest = tok.Next
		return node
	}

	panic("\n" + errorTok(tok, "expected expression: %s", tok.Str))
}

// typedef = "type" ident (type-preffix)? decl-spec
func parseTypedef(tok *Token) *Token {
	first := true

	for !consume(&tok, tok, ";") {
		if !first {
			tok = skip(tok, ",")
		}
		first = false

		ty := declarator(&tok, tok)
		pushScope(getIdent(ty.Name)).TyDef = ty
	}
	return tok
}

func createParamLvars(param *Type) {
	if param != nil {
		createParamLvars(param.Next)
		newLvar(getIdent(param.Name), param)
	}
}

// This function matches gotos with labels.
//
// We cannot resolve gotos as we parse a function because gotos
// can refer a label that apears later in the function.
// So, we need to do this after we parse the entire function.
func resolveGotoLabels() {
	for x := gotos; x != nil; x = x.GotoNext {
		for y := labels; y != nil; y = y.GotoNext {
			if x.Lbl == y.Lbl {
				x.UniqueLbl = y.UniqueLbl
				break
			}
		}

		if x.UniqueLbl == "" {
			panic("\n" + errorTok(x.Tok.Next, "use of undeclared label"))
		}
	}

	labels = nil
	gotos = nil
}

// function = "func" ident "(" params? ")" type-prefix type-specifier "{" stmt "}"
func function(tok *Token) *Token {
	// printCurTok()
	// printCalledFunc()

	ty := declarator(&tok, tok)

	fn := newGvar(getIdent(ty.Name), ty)
	fn.IsFunc = true
	fn.IsDef = !consume(&tok, tok, ";")

	if !fn.IsDef {
		return tok
	}

	curFn = fn
	locals = nil
	enterScope()
	createParamLvars(ty.Params)
	fn.Params = locals

	tok = skip(tok, "{")
	fn.Body = compoundStmt(&tok, tok)
	fn.Locals = locals
	leaveScope()
	resolveGotoLabels()
	return tok
}

// global-var = "var" ident type-prefix type-suffix ("=" gvar-initializer)? ";"
//
// For example,
// var x int = 6
// var x *int = &y
// var x string = "abc"
// var x [2]int = [2]int{1,2}
// var x T(typedef) = T{1,2}
func globalVar(tok *Token) *Token {
	// printCurTok()
	// printCalledFunc()

	first := true
	for !consume(&tok, tok, ";") {
		if !first {
			tok = skip(tok, ",")
		}
		first = false

		ty := declarator(&tok, tok)
		v := newGvar(getIdent(ty.Name), ty)
		if equal(tok, "=") {
			gvarInitializer(&tok, tok.Next, v)
		}
	}
	return tok
}

// program = (global-var | function)*
func parse(tok *Token) *Obj {
	// printCurTok()
	// printCalledFunc()

	globals = nil

	for !atEof() {

		if consume(&tok, tok, "func") {
			tok = function(tok)
			continue
		}

		if consume(&tok, tok, "var") {
			tok = globalVar(tok)
			continue
		}

		if consume(&tok, tok, "type") {
			tok = parseTypedef(tok)
			tok = skip(tok, ";")
			continue
		}

		panic("\n" + errorTok(token, "unexpected '%s'", token.Str))

	}

	return globals
}
