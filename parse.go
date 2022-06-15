package main

import (
	"fmt"
	"unsafe"
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
	Align   int
}

// Variable or function
type Obj struct {
	Next    *Obj
	Name    string // Variable name
	Ty      *Type  // Type
	Tok     *Token // for error message
	IsLocal bool   // local or global
	Align   int    // alignment

	// Local variables
	Offset int     // Offset from RBP
	Val    int64   // it's integer value
	FVal   float64 // it's floating-point value

	// Global variable or function
	IsFunc   bool
	IsDef    bool
	IsStatic bool

	// Global variables
	InitData []int64
	Rel      *Relocation

	// ret_buffer
	RetNext *Obj

	// Function
	Params  *Obj
	Body    *Node
	Locals  *Obj
	VaArea  *Obj
	StackSz int
	// Global var node for when there are more than 6 return values
	RetValGv *Node
	RetBufGv *Node
}

type NodeKind int

const (
	ND_NULL_EXPR      NodeKind = iota // Do nothing
	ND_ADD                            // +
	ND_SUB                            // -
	ND_MUL                            // *
	ND_DIV                            // /
	ND_NEG                            // unary -
	ND_MOD                            // %
	ND_BITAND                         // &
	ND_BITOR                          // |
	ND_BITXOR                         // ^
	ND_SHL                            // <<
	ND_SHR                            // >>
	ND_EQ                             // ==
	ND_NE                             // !=
	ND_LT                             // <
	ND_LE                             // <=
	ND_ASSIGN                         // =
	ND_COND                           // ?:
	ND_COMMA                          //
	ND_MEMBER                         // . (struct member access)
	ND_ADDR                           // unary &
	ND_DEREF                          // unary *
	ND_NOT                            // !
	ND_BITNOT                         // ~
	ND_LOGAND                         // &&
	ND_LOGOR                          // ||
	ND_RETURN                         // "return"
	ND_IF                             // "if"
	ND_FOR                            // "for" or "while"
	ND_SWITCH                         // "switch"
	ND_CASE                           // "case"
	ND_BLOCK                          // { ... }
	ND_GOTO                           // "goto"
	ND_LABEL                          // Labeled statement
	ND_FUNCALL                        // Function call
	ND_EXPR_STMT                      // Expression statement
	ND_STMT_EXPR                      // Statement expression
	ND_VAR                            // Variable
	ND_NUM                            // Integer
	ND_CAST                           // Type cast
	ND_MEMZERO                        // Zero-clear a stack variable
	ND_SIZEOF                         // 'Sizeof'
	ND_MULTIVALASSIGN                 // Assign multiple values in rhs to maultiple variables, like a,b = b, a.
	ND_MULTIRETASSIGN                 // Assign to multiple variables from functions returning multiple return values
	ND_BLANKIDENT                     // '_'
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
	FuncName    string
	FuncTy      *Type
	Args        *Node
	PassByStack bool
	RetBuf      *Obj

	// Function definition
	RetVals *Node // return values

	// Multi valued assignment
	Lhses *Node
	Rhses *Node

	// Assigning from functions returning multiple return values
	Masg *Node

	// Goto or labeled statement
	Lbl       string
	UniqueLbl string
	GotoNext  *Node

	// Switch-cases
	CaseNext   *Node
	DefCase    *Node
	CaseLbl    int
	CaseEndLbl string
	Expr       *Node

	Obj  *Obj  // used if kind == ND_VAR
	Val  int64 // used if kind == ND_NUM
	FVal float64
}

var locals *Obj
var globals *Obj

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
	printCurTok(tok)
	printCalledFunc()

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
	printCurTok(tok)
	printCalledFunc()

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
	printCurTok(tok)
	printCalledFunc()

	return &Node{Kind: kind, Tok: tok}
}

func newBinary(kind NodeKind, lhs *Node, rhs *Node, tok *Token) *Node {
	printCurTok(tok)
	printCalledFunc()

	return &Node{
		Kind: kind,
		Tok:  tok,
		Lhs:  lhs,
		Rhs:  rhs,
	}
}

func newUnary(kind NodeKind, expr *Node, tok *Token) *Node {
	printCurTok(tok)
	printCalledFunc()

	node := &Node{Kind: kind, Lhs: expr, Tok: tok}
	return node
}

func newNum(val int64, tok *Token) *Node {
	printCurTok(tok)
	printCalledFunc()

	return &Node{
		Kind: ND_NUM,
		Tok:  tok,
		Val:  val,
	}
}

func newLong(val int64, tok *Token) *Node {
	printCurTok(tok)
	printCalledFunc()

	return &Node{
		Kind: ND_NUM,
		Tok:  tok,
		Val:  val,
		Ty:   ty_long,
	}
}
func newUlong(val int64, tok *Token) *Node {
	printCurTok(tok)
	printCalledFunc()

	return &Node{
		Kind: ND_NUM,
		Tok:  tok,
		Val:  val,
		Ty:   ty_ulong,
	}
}

func newVarNode(v *Obj, tok *Token) *Node {
	printCurTok(tok)
	printCalledFunc()

	return &Node{Kind: ND_VAR, Tok: tok, Obj: v}
}

func newCast(expr *Node, ty *Type) *Node {
	printCalledFunc()

	addType(expr)

	return &Node{
		Kind: ND_CAST,
		Tok:  expr.Tok,
		Lhs:  expr,
		Ty:   copyType(ty),
	}
}

func pushScope(name string) *VarScope {
	printCalledFunc()

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
	Next   *Initializer
	Ty     *Type
	Tok    *Token
	IsFlex bool

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

func newInitializer(ty *Type, isflex bool) *Initializer {
	printCalledFunc()

	init := &Initializer{Ty: ty}

	if ty.Kind == TY_ARRAY {
		if isflex { //&& ty.Sz < 0
			init.IsFlex = true
			return init
		}
		init.Children = make([]*Initializer, ty.ArrSz)
		for i := 0; i < ty.ArrSz; i++ {
			init.Children[i] = newInitializer(ty.Base, false)
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
			if isflex && ty.IsFlex && mem.Next == nil {
				child := &Initializer{Ty: mem.Ty, IsFlex: true}
				init.Children[mem.Idx] = child
			} else {
				init.Children[mem.Idx] = newInitializer(mem.Ty, false)
			}
		}
		return init
	}

	return init
}

func newVar(name string, ty *Type) *Obj {
	printCalledFunc()

	v := &Obj{Name: name, Ty: ty, Align: ty.Align}
	if name != "_" {
		pushScope(name).Obj = v
	}
	return v
}

func newLvar(name string, ty *Type) *Obj {
	printCalledFunc()

	v := newVar(name, ty)
	v.IsLocal = true
	v.Next = locals
	locals = v
	return v
}

func newGvar(name string, ty *Type) *Obj {
	printCalledFunc()

	v := newVar(name, ty)
	if ty.Kind == TY_FUNC {
		v.IsFunc = true
	}
	// If name[0] is not uppercase, the 'Obj' can't be exported.
	if ('A' > name[0] || name[0] > 'Z') && name != "main" {
		v.IsStatic = true
	}

	v.IsDef = true
	v.Next = globals
	globals = v
	return v
}

// for newUniqueName function
var cnt int

func newUniqueName() string {
	printCalledFunc()

	res := fmt.Sprintf(".L..%d", cnt)
	cnt++
	return res
}

func newAnonGvar(ty *Type) *Obj {
	printCalledFunc()

	return newGvar(newUniqueName(), ty)
}

func newStringLiteral(p []int64, ty *Type) *Obj {
	printCalledFunc()

	v := newAnonGvar(ty)
	v.InitData = p
	return v
}

func newFavName(s string) string {
	printCalledFunc()

	res := fmt.Sprintf("%s.%d", s, cnt)
	cnt++
	return res
}

func newFavGvar(s string, ty *Type) *Obj {
	printCalledFunc()

	return newGvar(newFavName(s), ty)
}

func getIdent(tok *Token) string {
	printCurTok(tok)
	printCalledFunc()

	printCalledFunc()

	if tok.Kind != TK_IDENT {
		errorTok(tok, "expected an identifier")
	}
	return tok.Str
}

func findTyDef(tok *Token) *Type {
	printCurTok(tok)
	printCalledFunc()

	if tok.Kind == TK_IDENT {
		if sc := findVar(tok); sc != nil {
			return sc.TyDef
		}
	}
	return nil
}

func pushTagScope(tok *Token, ty *Type) {
	printCurTok(tok)
	printCalledFunc()

	sc := &TagScope{
		Name: tok.Str,
		Ty:   ty,
		Next: scope.Tags,
	}
	scope.Tags = sc
}

// declSpec returns a pointer of Type struct.
// If the current tokens represents a typename,
// it returns the Type struct with that typename.
// Otherwise returns the Type struct with TY_VOID.
//
// declspec = "*"* builtin-type | struct-decl | typedef-name |
// builtin-type = void | "bool" | "byte"| "int16" | "int" | "int64" |
//                "string"
//
func declSpec(rest **Token, tok *Token, name *Token) *Type {
	printCurTok(tok)
	printCalledFunc()

	nPtr := 0
	for consume(&tok, tok, "*") {
		nPtr++
	}

	var ty *Type
	if consume(&tok, tok, "byte") {
		ty = ty_uchar
	} else if consume(&tok, tok, "bool") {
		ty = ty_bool
	} else if consume(&tok, tok, "int8") {
		ty = ty_char
	} else if consume(&tok, tok, "int16") {
		ty = ty_short
	} else if consume(&tok, tok, "int") {
		ty = ty_int
	} else if consume(&tok, tok, "int32") {
		ty = ty_int
	} else if consume(&tok, tok, "int64") {
		ty = ty_long
	} else if consume(&tok, tok, "uint8") {
		ty = ty_uchar
	} else if consume(&tok, tok, "uint16") {
		ty = ty_ushort
	} else if consume(&tok, tok, "uint32") {
		ty = ty_uint
	} else if consume(&tok, tok, "uint") {
		ty = ty_uint
	} else if consume(&tok, tok, "uint64") {
		ty = ty_ulong
	} else if consume(&tok, tok, "float32") {
		ty = ty_float
	} else if consume(&tok, tok, "float64") {
		ty = ty_double
	} else if consume(&tok, tok, "string") {
		ty = stringType()
	} else if consume(&tok, tok, "struct") { // struct type
		ty = structDecl(&tok, tok, name)
	} else if consume(&tok, tok, "func") { // func type ,like: "func(int,string) int8"
		ty = funcDecl(&tok, tok, name)
	}

	// Handle user-defined types.
	ty2 := findTyDef(tok)
	if ty2 != nil {
		ty = ty2
		tok = tok.Next
	}

	if ty == nil {
		return ty_void
	}

	for i := 0; i < nPtr; i++ {
		ty = pointerTo(ty)
	}

	*rest = tok
	return ty
}

func findBase(rest **Token, tok *Token, name *Token) *Type {
	printCurTok(tok)
	printCalledFunc()

	for !(isTypename2(tok) && !equal(tok.Next, "(")) && // builtin type-name except for type cast
		!(equal(tok, "*") && isTypename2(tok.Next)) && // pointer to type name
		!(equal(tok, "func") && equal(tok.Next, "(")) { // function type
		tok = tok.Next
	}
	ty := declSpec(&tok, tok, name)
	*rest = tok // どこまでtokenを読んだか
	return ty
}

func readArr(tok *Token, base *Type) *Type {
	printCurTok(tok)
	printCalledFunc()

	if !consume(&tok, tok, "[") {
		return base
	}
	if !consume(&tok, tok, "]") {
		sz := constExpr(&tok, tok)
		tok = skip(tok, "]")
		base = readArr(tok, base)
		return arrayOf(base, int(sz))
	}
	base = readArr(tok, base)
	return sliceType(base, 0, 0)
}

// type-preffix = ("[" const-expr "]")*
func readTypePreffix(rest **Token, tok *Token, name *Token) *Type {
	printCurTok(tok)
	printCalledFunc()

	if consume(&tok, tok, "*") {
		return pointerTo(readTypePreffix(rest, tok, name))
	}

	if !equal(tok, "[") {
		return declSpec(rest, tok, name)
	}

	start := tok

	base := findBase(&tok, tok, name)
	arrTy := readArr(start, base)
	*rest = tok
	return arrTy
}

// declarator = ident (type-preffix)? declspec
//            | ident type-suffix
func declarator(rest **Token, tok *Token) *Type {
	printCurTok(tok)
	printCalledFunc()

	var name *Token
	var namePos *Token = tok

	if tok.Kind == TK_IDENT || tok.Kind == TK_BLANKIDENT {
		name = tok
		tok = tok.Next
	}

	var ty *Type
	if equal(tok, "(") {
		ty = typeSuffix(&tok, tok, nil)
	} else {
		ty = readTypePreffix(&tok, tok, name)
	}
	*rest = tok
	ty.Name = name
	ty.NamePos = namePos
	return ty
}

// func-params = (param ("," param)* ("," "...")? ")"
// param = declarator
// e.g.
//  x int
//  x *int
//  x **int
//  x [3]int
//  x [3]*int
//  x [2]**int
func funcParams(rest **Token, tok *Token, ty *Type) *Type {
	printCurTok(tok)
	printCalledFunc()

	isVariadic := false
	first := true

	paramList := make([]*Type, 0)

	for !equal(tok, ")") {
		if !first {
			tok = skip(tok, ",")
		}
		first = false

		ty2 := declarator(&tok, tok)
		ty2name := ty2.Name
		if ty2.Kind == TY_VOID {
			if equal(tok, "...") {
				isVariadic = true
				ty2 = readTypePreffix(&tok, tok.Next, nil)
				ty2.Name = ty2name
				paramList = append(paramList, copyType(ty2))
				break
			}
		}

		name := ty2.Name

		if ty2.Kind == TY_ARRAY {
			// "array of T" is converted tot "pointer to T" only in the parameter
			// context. For example, *argv[] is converted to **argv by this.
			ty2 = pointerTo(ty2.Base)
			ty2.Name = name
		}

		paramList = append(paramList, copyType(ty2))
	}

	// Handle the cases that typename omittied, like "var a,b int"
	cnt := 0
	for i := 0; i < len(paramList); i++ {
		param := paramList[i]
		if param.Kind == TY_VOID {
			cnt++
		} else {
			ty3 := param
			for j := 0; j < cnt; j++ {
				name := paramList[i-(j+1)].Name
				paramList[i-(j+1)] = copyType(ty3)
				paramList[i-(j+1)].Name = name
			}
		}
	}

	if len(paramList) == 0 {
		isVariadic = true
	} else if cnt == len(paramList) {
		panic(errorTok(tok, "type name expected"))
	}

	// Make a linked list.
	head := &Type{}
	cur := head
	for i := 0; i < len(paramList); i++ {
		cur.Next = paramList[i]
		cur = cur.Next
	}

	ty = funcType(ty, head.Next)
	ty.IsVariadic = isVariadic
	*rest = tok.Next
	return ty
}

// type-suffix = ("(" func-params )?
func typeSuffix(rest **Token, tok *Token, ty *Type) *Type {
	printCurTok(tok)
	printCalledFunc()

	if equal(tok, "(") {
		return funcParams(rest, tok.Next, ty)
	}

	*rest = tok
	return ty
}

func isEnd(tok *Token) bool {
	printCurTok(tok)
	printCalledFunc()

	return equal(tok, "}") || (equal(tok, ",") && equal(tok.Next, "}"))
}

func consumeEnd(rest **Token, tok *Token) bool {
	printCalledFunc()

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

func skipExcessElement(tok *Token) *Token {
	printCurTok(tok)
	printCalledFunc()

	if equal(tok, "{") {
		tok = skipExcessElement(tok.Next)
		return skip(tok, "}")
	}

	assign(&tok, tok)
	return tok
}

// string-initializer = string-literal
func stringInitializer(rest **Token, tok *Token, init *Initializer) {
	printCurTok(tok)
	printCalledFunc()

	if init.IsFlex {
		*init = *newInitializer(arrayOf(init.Ty.Base, tok.Ty.ArrSz), false)
	}

	length := min(init.Ty.ArrSz, tok.Ty.ArrSz)
	for i := 0; i < length; i++ {
		init.Children[i].Expr = newNum(int64(tok.Contents[i]), tok)
	}
	init.Ty.Len = length
	*rest = tok.Next
}

// struct-designator = ident ":"
//
// Use `fieldname :` to move the cursor for a struct initializer. E.g.
//
//   type T struct { a int; b int; c int; };
//   var x T = T{c: 5};
//
// The above initializer sets x.c to 5.
func structDesignator(rest **Token, tok *Token, ty *Type) *Member {
	printCurTok(tok)
	printCalledFunc()

	for mem := ty.Mems; mem != nil; mem = mem.Next {
		if mem.Name.Len == tok.Len && mem.Name.Str == tok.Str {
			tok = skip(tok.Next, ":")
			*rest = tok
			return mem
		}
	}

	panic("\n" + errorTok(tok, "struct has no such member"))
}

// designation = struct-designator initializer
func designation(rest **Token, tok *Token, init *Initializer) {
	printCurTok(tok)
	printCalledFunc()

	if tok.Kind == TK_IDENT && equal(tok.Next, ":") {
		mem := structDesignator(&tok, tok, init.Ty)
		designation(&tok, tok, init.Children[mem.Idx])
		init.Expr = nil
		structInitializer2(rest, tok, init, mem.Next)
	}

	initializer2(rest, tok, init)
}

// array-initializer = (type-preffix)? decl-spec "{" initializer ("," initializer)* ","? "}"
func arrayInitializer(rest **Token, tok *Token, init *Initializer) {
	printCurTok(tok)
	printCalledFunc()

	tok = skip(tok, "{")

	if init.IsFlex {
		len := countArrInitElem(tok, init.Ty)
		*init = *newInitializer(arrayOf(init.Ty.Base, len), false)
	}

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

// array-initializer2 = initializer ("," initializer)*
func arrayInitializer2(rest **Token, tok *Token, init *Initializer) {
	printCurTok(tok)
	printCalledFunc()

	if init.IsFlex {
		len := countArrInitElem(tok, init.Ty)
		*init = *newInitializer(arrayOf(init.Ty.Base, len), false)
	}

	for i := 0; i < init.Ty.ArrSz && !isEnd(tok); i++ {
		if i > 0 {
			tok = skip(tok, ",")
		}
		initializer2(&tok, tok, init.Children[i])
	}
	*rest = tok
}

// struct-initializer = "{" initializer ("," initializer)* ","? "}"
func structInitializer(rest **Token, tok *Token, init *Initializer) {
	printCurTok(tok)
	printCalledFunc()

	tok = skip(tok, "{")

	mem := init.Ty.Mems
	first := true

	for !consumeEnd(rest, tok) {
		if !first {
			tok = skip(tok, ",")
		}
		first = false

		if tok.Kind == TK_IDENT && equal(tok.Next, ":") {
			mem = structDesignator(&tok, tok, init.Ty)
			designation(&tok, tok, init.Children[mem.Idx])
			mem = mem.Next
			continue
		}

		if mem != nil {
			initializer2(&tok, tok, init.Children[mem.Idx])
			mem = mem.Next
		} else {
			tok = skipExcessElement(tok)
		}
	}
}

// struct-initializer2 = initializer ("," initializer)*
func structInitializer2(rest **Token, tok *Token, init *Initializer, mem *Member) {
	printCurTok(tok)
	printCalledFunc()

	first := true

	for ; mem != nil && !isEnd(tok); mem = mem.Next {
		start := tok

		if !first {
			tok = skip(tok, ",")
		}
		first = false

		if tok.Kind == TK_IDENT && equal(tok, ":") {
			*rest = start
			return
		}

		initializer2(&tok, tok, init.Children[mem.Idx])
	}
	*rest = tok
}

func countArrInitElem(tok *Token, ty *Type) int {
	printCurTok(tok)
	printCalledFunc()

	dummy := newInitializer(ty.Base, false)
	i := 0

	for ; !consumeEnd(&tok, tok); i++ {
		if i > 0 {
			tok = skip(tok, ",")
		}
		initializer2(&tok, tok, dummy)
	}
	return i
}

// initializer = string-initializer | array-initializer
//             | struct-initializer
//             | assign
func initializer2(rest **Token, tok *Token, init *Initializer) {
	printCurTok(tok)
	printCalledFunc()

	// If the rhs is string literal.
	if init.Ty.Kind == TY_ARRAY && tok.Kind == TK_STR {
		stringInitializer(rest, tok, init)
		init.Ty.Init = init
		return
	}

	// If the rhs is array literal.
	if init.Ty.Kind == TY_ARRAY {
		readTypePreffix(&tok, tok, nil) // I'll add type checking later
		if equal(tok, "{") {
			if !equal(tok.Next, "}") {
				arrayInitializer(rest, tok, init)
			} else {
				zeroInit2(init, tok)
				*rest = skip(tok.Next, "}")
			}
		} else {
			arrayInitializer2(rest, tok, init)
		}
		init.Ty.Init = init
		return
	}

	// If the rhs is slice.
	if init.Ty.Kind == TY_SLICE {
		sliceTy := readTypePreffix(&tok, tok, nil) // I'll add type checking later
		if sliceTy.Kind == TY_VOID && !equal(tok, "{") {
			// In the case that no typename is written, like: `var x = a[0:1]`
			if init.Expr == nil {
				init.Expr = assign(rest, tok)
			}
			init.Ty.Init = init
			return
		}

		// In the case that any typename is written, like: `var x = []int{1,3}`,
		// make the underlying array.
		uArrTy := arrayOf(init.Ty.Base, 0)
		uArrTy.IsFlex = true
		uArr := newFavGvar("underlying_array", uArrTy)
		cnt++

		gvarInitializer(rest, tok, uArr)

		init.Ty.Len = uArr.Ty.ArrSz
		init.Ty.Cap = uArr.Ty.ArrSz
		uaNode := newVarNode(uArr, tok)
		init.Ty.UArrNode = uaNode

		init.Expr = newUnary(ND_ADDR,
			newUnary(ND_DEREF,
				newAdd(init.Ty.UArrNode, newNum(0, tok), tok), tok), tok)

		init.Ty.Init = init
		return
	}

	if init.Ty.Kind == TY_STRUCT {
		if equal(tok.Next, "{") {
			readTypePreffix(&tok, tok, nil) // I'll add type checking later
		}

		if equal(tok, "{") {
			if !equal(tok.Next, "}") {
				structInitializer(rest, tok, init)
				return
			} else {
				zeroInit2(init, tok)
				*rest = skip(tok.Next, "}")
				return
			}
		}
		// A struct can be initialized with another struct. E.g.
		// `type x y` where y is a another struct.
		// Handle that case first.
		expr := assign(rest, tok)
		addType(expr)
		if expr.Ty.Kind == TY_STRUCT {
			init.Expr = expr
			return
		}

		structInitializer2(rest, tok, init, init.Ty.Mems)
		return
	}

	// If type-name is omitted.
	if init.Ty.Kind == TY_VOID {
		var rhsTy *Type
		if tok.Kind == TK_STR {
			init.Ty = stringType()
			initializer2(rest, tok, init)
			return
		}

		// Get the type from rhs. If type-name is written.
		// like: [2]int{1,2,3}
		rhsTy = readTypePreffix(&tok, tok, nil)

		// If type-name isn't written in the rhs.
		var start *Token = tok
		var startNext *Token = tok.Next
		if rhsTy.Kind == TY_VOID {
			init.Expr = assign(rest, tok)
			addType(init.Expr)

			if init.Expr.Ty.Kind == TY_PTR &&
				init.Expr.Lhs != nil && init.Expr.Lhs.Ty.Kind == TY_ARRAY {
				// the rhs is like "&" and variable.
				rhsTy = pointerTo(init.Expr.Lhs.Ty)

			} else if init.Expr.Ty.Kind == TY_FUNC {
				rhsTy = pointerTo(init.Expr.Ty)

			} else {
				rhsTy = init.Expr.Ty
			}
			// panic(errorTok(tok, "the lhs and rhs both declared void"))
		}

		init.Ty = rhsTy

		// Initialize the lhs.
		if init.Ty.Kind == TY_ARRAY {
			if equal(start, "{") || equal(startNext, "{") {
				init.Children = make([]*Initializer, init.Ty.ArrSz)
				for i := 0; i < init.Ty.ArrSz; i++ {
					init.Children[i] = newInitializer(init.Ty.Base, false)
				}
				initializer2(rest, tok, init)
				init.Ty.Init = init
				return
			}

			// Copy Initializer from rhs, if array can be initialized by other array.
			if rhsTy.Init != nil {
				*init = *rhsTy.Init //
			}
			return
		}

		if (equal(start, "{") || equal(startNext, "{")) &&
			init.Ty.Kind == TY_STRUCT {
			// Count the number of struct members
			var l int
			for mem := init.Ty.Mems; mem != nil; mem = mem.Next {
				l++
			}

			init.Children = make([]*Initializer, l)

			for mem := init.Ty.Mems; mem != nil; mem = mem.Next {
				init.Children[mem.Idx] = newInitializer(mem.Ty, false)
			}
			initializer2(rest, tok, init)
			return
		}
		initializer2(rest, tok, init)
		return
	}

	init.Expr = assign(rest, tok)

}

func copyStructType(ty *Type) *Type {
	ty = copyType(ty)

	head := &Member{}
	cur := head
	for mem := ty.Mems; mem != nil; mem = mem.Next {
		m := mem
		cur.Next = m
		cur = cur.Next
	}

	ty.Mems = head.Next
	return ty
}

func initializer(rest **Token, tok *Token, ty *Type, newTy **Type, v *Obj) *Initializer {
	printCurTok(tok)
	printCalledFunc()

	init := newInitializer(ty, ty.IsFlex)
	initializer2(rest, tok, init)

	if ty.Kind == TY_STRUCT && ty.IsFlex {
		ty = copyStructType(ty)

		mem := ty.Mems
		for mem.Next != nil {
			mem = mem.Next
		}
		mem.Ty = init.Children[mem.Idx].Ty
		ty.Sz += mem.Ty.Sz
	}

	*newTy = init.Ty
	// Change variable's align
	v.Align = init.Ty.Align

	if isInteger(init.Ty) {
		v.Val = eval(init.Expr)
	} else if isFlonum(init.Ty) {
		v.FVal = evalDouble(init.Expr)
	}

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
	printCurTok(tok)
	printCalledFunc()

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

	if desg.Var != nil && desg.Var.Name == "_" {
		return init.Expr
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
	printCurTok(tok)
	printCalledFunc()

	// Initialize a char array with a string literal.
	// => unnecessary for this compiler, I think.

	init := initializer(rest, tok, v.Ty, &v.Ty, v)
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

func writeBuf(buf unsafe.Pointer, val int64, sz int) {
	switch sz {
	case 1:
		*(*uint8)(buf) = uint8(val)
	case 2:
		*(*uint16)(buf) = uint16(val)
	case 4:
		*(*uint32)(buf) = uint32(val)
	case 8:
		*(*uint64)(buf) = uint64(val)
	default:
		panic("writeBuf: internal error")
	}
}

// divFloat32: floatの内部表現を分割する
// 例：1.5
// => 00111111 11000000 00000000 00000000 2進数にする
// => 00000000 00000000 11000000 00111111 リトルエンディアン
// => 0, 0, 192, 63 それぞれintにする
func divFloat32(target int32) []int64 {
	t := fmt.Sprintf("%032b", target)
	ret := make([]int64, 0, 1024)
	for i := len(t) - 8; i >= 0; i -= 8 {
		s := t[i : i+8]
		num := parseInt(s, 2)
		ret = append(ret, num)
	}
	return ret
}

func divFloat64(target int64) []int64 {
	t := fmt.Sprintf("%064b", target)
	ret := make([]int64, 0, 1024)
	for i := len(t) - 8; i >= 0; i -= 8 {
		s := t[i : i+8]
		num := parseInt(s, 2)
		ret = append(ret, num)
	}
	return ret
}

//
//
func writeGvarData(
	cur *Relocation, init *Initializer, ty *Type, buf *[]int64,
	offset int) *Relocation {
	printCalledFunc()

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

	if ty.Kind == TY_FLOAT {
		fval := float32(evalDouble(init.Expr))
		// float32(evalDouble(init.Expr))の内部表現(2進数で取得)をintとして読んだものを取得し
		// 分割してスライスにしてdivedに保存
		dived := divFloat32(*(*int32)(unsafe.Pointer(&fval)))
		for i, j := offset, 0; i < offset+ty.Sz && j < 4; i++ {
			(*buf)[i] = dived[j]
			j++
		}
		return cur
	}

	if ty.Kind == TY_DOUBLE {
		fval := evalDouble(init.Expr)
		dived := divFloat64(*(*int64)(unsafe.Pointer(&fval)))
		for i, j := offset, 0; i < offset+ty.Sz && j < 8; i++ {
			(*buf)[i] = dived[j]
			j++
		}
		return cur
	}

	var label *string = nil
	var val = eval2(init.Expr, &label)

	if label == nil {
		writeBuf(unsafe.Pointer(&((*buf)[offset])), val, ty.Sz)
		return cur
	}

	rel := &Relocation{
		Offset: offset,
		Lbl:    *label,
		Addend: val,
	}
	cur.Next = rel
	return cur.Next
}

func gvarInitializer(rest **Token, tok *Token, v *Obj) {
	printCurTok(tok)
	printCalledFunc()

	init := initializer(rest, tok, v.Ty, &v.Ty, v)
	head := &Relocation{}
	var buf []int64 = make([]int64, v.Ty.Sz)
	writeGvarData(head, init, v.Ty, &buf, 0)
	v.InitData = buf
	v.Rel = head.Next
}

func gvarZeroInit(v *Obj, tok *Token) {
	printCurTok(tok)
	printCalledFunc()

	init := zeroInit(v.Ty, &v.Ty, tok)
	head := &Relocation{}
	var buf []int64 = make([]int64, v.Ty.Sz)
	writeGvarData(head, init, v.Ty, &buf, 0)
	v.InitData = buf
	v.Rel = head.Next
}

// abstruct-declarator = "*"* declspec ("(" abstruct-declarator ")")? type-suffix
func abstructDeclarator(rest **Token, tok *Token, ty *Type) *Type {

	nPtr := 0
	for equal(tok, "*") {
		nPtr++
		tok = tok.Next
	}

	if isTypename(tok.Next) {
		ty = declSpec(&tok, tok, nil)
	}

	for i := 0; i < nPtr; i++ {
		ty = pointerTo(abstructDeclarator(&tok, tok, ty))
	}

	if equal(tok, "(") {
		start := tok
		ty = abstructDeclarator(&tok, start.Next, ty)
		tok = skip(tok, ")")
		ty = typeSuffix(rest, tok, ty)
		return abstructDeclarator(&tok, start.Next, ty)
	}

	return typeSuffix(rest, tok, ty)
}

// type-name = abstruct-declarator
func typename(rest **Token, tok *Token) *Type {
	return abstructDeclarator(rest, tok, &Type{})
}

func strZeroInit(init *Initializer, tok *Token) {
	printCurTok(tok)
	printCalledFunc()

	child := &Initializer{
		Tok:  tok,
		Ty:   init.Ty.Base,
		Expr: newNum(0, tok),
	}
	init.Children = append(init.Children, child)
}

func zeroInit2(init *Initializer, tok *Token) {
	printCurTok(tok)
	printCalledFunc()

	// If init.Ty is string.
	if init.Ty.TyName == "string" {
		strZeroInit(init, tok)
		tokTy := arrayOf(ty_char, 1)
		tokConts := []int64{0}
		v := newStringLiteral(tokConts, tokTy)
		init.Expr = newVarNode(v, tok)
		init.Ty.Init = init
		return
	}

	// If init.Ty is array.
	if init.Ty.Kind == TY_ARRAY {
		for i := 0; i < init.Ty.ArrSz; i++ {
			zeroInit2(init.Children[i], tok)
		}
		init.Ty.Init = init
		return
	}

	// If init.Ty is struct.
	if init.Ty.Kind == TY_STRUCT {
		for mem := init.Ty.Mems; mem != nil; mem = mem.Next {
			zeroInit2(init.Children[mem.Idx], tok)
		}
		init.Ty.Init = init
		return
	}

	// If init.Ty is slice.
	if init.Ty.Kind == TY_SLICE {
		// Make the underlying array.
		uArrTy := arrayOf(init.Ty.Base, init.Ty.Cap)
		uArr := newFavGvar("underlying_array", uArrTy)
		cnt++
		gvarZeroInit(uArr, tok)
		uaNode := newVarNode(uArr, tok)
		init.Expr = newUnary(ND_ADDR,
			newUnary(ND_DEREF, newAdd(uaNode, newNum(0, tok), tok), tok),
			tok)
		init.Ty = sliceType(uArr.Ty.Base, init.Ty.Len, init.Ty.Cap)
		init.Ty.UArrNode = uaNode
		init.Ty.Init = init
		return
	}

	init.Expr = newNum(0, tok)
}

func zeroInit(ty *Type, newTy **Type, tok *Token) *Initializer {
	printCurTok(tok)
	printCalledFunc()

	init := newInitializer(ty, ty.IsFlex)
	zeroInit2(init, tok)

	*newTy = init.Ty
	return init
}

func lvarZeroInit(v *Obj, tok *Token) *Node {
	printCurTok(tok)
	printCalledFunc()

	init := zeroInit(v.Ty, &v.Ty, tok)
	desg := &InitDesg{nil, 0, nil, v}

	// If no initializer list is given, the variable is initialized
	// with 0.
	lhs := newNode(ND_MEMZERO, tok)
	lhs.Obj = v

	rhs := createLvarInit(init, v.Ty, desg, tok)
	addType(rhs)
	return newBinary(ND_COMMA, lhs, rhs, tok)
}

// declaration = VarDecl | VarSpec | ShortVarDecl
// VarDecl = "var" ( VarSpec | "("  { VarSpec ";" } ")" ) .
// VarSpec = IdentifierList (Type [ "=" ExpressionList ] | "=" ExpressionList)
// ShortVarDecl = IdentifierList ":=" ExpressionList .
func declaration(rest **Token, tok *Token, isShort bool) *Node {
	printCurTok(tok)
	printCalledFunc()

	var i int
	identList := make([]*Obj, 0)

	// Read the Lhs
	for !equal(tok, "=") && !equal(tok, ":=") && !equal(tok, ";") {
		if i > 0 {
			tok = skip(tok, ",")
		}
		i++
		ty := declarator(&tok, tok)
		if ty.Name == nil {
			panic(errorTok(ty.NamePos, "variable name omitted"))
		}
		if equal(tok, ",") && ty.Kind != TY_VOID {
			panic(errorTok(ty.NamePos, "expected ';' found ','"))
		}

		v := newLvar(getIdent(ty.Name), ty)
		identList = append(identList, v)
	}

	ty := copyType(identList[len(identList)-1].Ty)
	for j := len(identList) - 2; j >= 0; j-- {
		identList[j].Ty = ty
	}

	head := &Node{}
	cur := head

	// Read the Rhs or initialize the variables with 0.
	if (!isShort && equal(tok, "=")) || (isShort && equal(tok, ":=")) {
		start := tok

		if len(identList) > 1 {
			if rhs := expr(&tok, tok.Next); rhs.Kind == ND_FUNCALL {
				// For funcall that the function returing multiple values.
				// Read the variables in Lhs.
				ty := rhs.Lhs.Obj.Ty.RetTy
				for j := 0; j < len(identList); j++ {
					if identList[j].Ty.Kind == TY_VOID {
						identList[j].Ty = copyType(ty)
					}
					cur.Next = newVarNode(identList[j], identList[j].Tok)
					cur = cur.Next
					addType(cur)
					ty = ty.Next
				}
				numVals := countRetTys(rhs.Lhs.Obj)
				// とりあえずlhsesの長さだけで判定、エラーも適当
				if len(identList) != numVals {
					panic(errorTok(tok, "too many assigns: left:%d, right:%d", i, numVals))
				}
				node := newUnary(ND_MULTIRETASSIGN, rhs, tok)
				node.Masg = head.Next
				*rest = tok.Next
				addType(node)
				return node
			}
		}

		tok = start
		j := 0
		for !equal(tok, ";") {
			v := identList[j]
			expr := lvarInitializer(&tok, tok.Next, v)
			addType(expr)
			cur.Next = newUnary(ND_EXPR_STMT, expr, tok)
			cur = cur.Next
			j++

			if v.Ty.Sz < 0 {
				panic("\n" +
					errorTok(v.Ty.Name, "variable has incomplete type"))
			}
			if v.Ty.Kind == TY_VOID ||
				(v.Ty.Base != nil && v.Ty.Base.Kind == TY_VOID) {
				panic("\n" + errorTok(v.Ty.Name, "variable declared void"))
			}
		}

	} else {

		for j := 0; j < len(identList); j++ {
			v := identList[j]
			// Initialize empty variables.
			expr := lvarZeroInit(v, tok)
			addType(expr)
			cur.Next = newUnary(ND_EXPR_STMT, expr, tok)
			cur = cur.Next
		}

	}

	node := newNode(ND_BLOCK, tok)
	node.Body = head.Next
	*rest = tok.Next
	return node
}

func isTypename2(tok *Token) bool {
	for i := 0; i < len(tyName); i++ {
		if equal(tok, tyName[i]) {
			return true
		}
	}
	return findTyDef(tok) != nil
}

func isTypename(tok *Token) bool {
	printCurTok(tok)
	printCalledFunc()

	for equal(tok, "*") {
		tok = tok.Next
	}

	if equal(tok, "[") {
		for !equal(tok, ";") {
			if equal(tok, "]") && equal(tok.Next, "[") {
				tok = tok.Next.Next
				continue
			}
			if equal(tok, "]") {
				tok = tok.Next
				break
			}
			tok = tok.Next
		}
	}

	for equal(tok, "*") {
		tok = tok.Next
	}

	return isTypename2(tok)
}

// isForClause returns true and exceeds the next token, if ";" will be found
// between "for" and "{".
func hasSimpleStmt(tok *Token) bool {
	printCurTok(tok)
	printCalledFunc()

	for !equal(tok, "{") {
		if equal(tok, ";") {
			return true
		}
		tok = tok.Next
	}
	return false
}

// stmt = "return" expr? ";"
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
//      | ident ":" stmt
//      | assign-list
//      | expr ";"
// for-stmt = "for" [ condition ] block .
// for-clause = "for" (expr? ";" | declaration) condition ";" expr? block
// condition = expr .
// block = "{" stmt-list "};" .
// stmt-list = { stmt ";" } .
func stmt(rest **Token, tok *Token) *Node {
	printCurTok(tok)
	printCalledFunc()

	if equal(tok, "return") {
		node := newNode(ND_RETURN, tok)
		if consume(rest, tok.Next, ";") {
			return node
		}
		tok = skip(tok, "return")

		head := &Node{}
		cur := head
		rvghead := &Node{}
		rvgcur := rvghead
		bufgvhead := &Node{}
		bufgvcur := bufgvhead
		idx := 0
		ty := copyType(curFn.Ty.RetTy)

		for !equal(tok, ";") {
			if idx > 0 {
				tok = skip(tok, ",")
			}

			exp := assign(&tok, tok)
			addType(exp)
			if ty.Kind != TY_STRUCT {
				exp = newCast(exp, ty)
			}

			if idx >= 6 {
				rvgcur.Next = newVarNode(newFavGvar("ret_gv", ty), tok)
				rvgcur = rvgcur.Next
				addType(rvgcur)

				if ty.Kind == TY_STRUCT && 8 < ty.Sz && ty.Sz <= 16 {
					bufgvcur.Next = newVarNode(newFavGvar("buf_gv", ty), tok)
					bufgvcur = bufgvcur.Next
					addType(bufgvcur)
				}
			}

			cur.Next = exp
			cur = cur.Next
			ty = ty.Next
			idx++
		}

		node.RetVals = head.Next
		curFn.RetValGv = rvghead.Next
		curFn.RetBufGv = bufgvhead.Next
		return node
	}

	if equal(tok, "if") {
		node := newNode(ND_IF, tok)
		enterScope()
		// Read 'SimpleStmt'
		if hasSimpleStmt(tok) {
			node.Init = expr(&tok, tok.Next)
		} else {
			tok = tok.Next
		}

		node.Cond = expr(&tok, tok)
		node.Then = stmt(&tok, tok)
		if equal(tok, "else") {
			node.Els = stmt(&tok, tok.Next)
		}
		leaveScope()
		*rest = tok
		return node
	}

	if equal(tok, "switch") {
		node := newNode(ND_SWITCH, tok)
		enterScope()
		// Read 'SimpleStmt'
		if hasSimpleStmt(tok) {
			node.Init = expr(&tok, tok.Next)
		} else {
			tok = tok.Next
		}
		if !equal(tok, "{") {
			node.Cond = expr(&tok, tok)
		} else {
			node.Cond = newNum(1, tok)
		}
		sw := curSwitch
		curSwitch = node

		var brk string = brkLabel
		node.BrkLabel = newUniqueName()
		brkLabel = node.BrkLabel

		node.Then = stmt(rest, tok)

		leaveScope()
		curSwitch = sw
		brkLabel = brk
		return node
	}

	if equal(tok, "case") {
		if curSwitch == nil {
			panic("\n" + errorTok(tok, "stray case"))
		}

		tok = skip(tok, "case")
		var head = &Node{}
		var cur = head
		var first = true

		for !equal(tok, ":") {
			if !first {
				tok = skip(tok, ",")
			}
			first = false

			node := newNode(ND_CASE, tok)
			node.Expr = assign(&tok, tok)
			addType(node.Expr)
			node.Lbl = newUniqueName()
			node.CaseNext = curSwitch.CaseNext
			curSwitch.CaseNext = node

			cur.Next = node
			cur = cur.Next
		}

		tok = skip(tok, ":")
		head2 := &Node{}
		cur2 := head2
		for !equal(tok, "case") && !equal(tok, "default") && !equal(tok, "}") {
			cur2.Next = stmt(&tok, tok)
			cur2 = cur2.Next
		}
		*rest = tok
		lhs := newNode(ND_BLOCK, tok)
		lhs.Body = head2.Next

		cur = head.Next

		for cur != nil {
			cur.Lhs = lhs
			cur = cur.Next
		}

		node := newNode(ND_BLOCK, tok)
		node.Body = head.Next
		return node
	}

	if equal(tok, "default") {
		if curSwitch == nil {
			panic("\n" + errorTok(tok, "stray default"))
		}
		node := newNode(ND_CASE, tok)
		tok = skip(tok.Next, ":")
		node.Lbl = newUniqueName()
		node.Lhs = stmt(rest, tok)
		curSwitch.DefCase = node
		return node
	}

	if equal(tok, "for") {
		if !hasSimpleStmt(tok) { // for-stmt like 'while' statement
			node := newNode(ND_FOR, tok)
			if !equal(tok.Next, "{") {
				node.Cond = expr(&tok, tok.Next)
			} else {
				node.Cond = newNum(1, tok)
				tok = tok.Next
			}

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

		} else { // for-clause
			node := newNode(ND_FOR, tok)
			enterScope()
			brk := brkLabel
			cont := contLabel
			node.BrkLabel = newUniqueName()
			brkLabel = node.BrkLabel
			node.ContLabel = newUniqueName()
			contLabel = node.ContLabel

			if !equal(tok.Next, ";") {
				if tok.Next.Kind == TK_IDENT && equal(tok.Next.Next, ":=") {
					node.Init = declaration(&tok, tok.Next, true)
				} else {
					node.Init = exprStmt(&tok, tok.Next)
				}
			} else {
				tok = skip(tok.Next, ";")
			}

			if !equal(tok, ";") {
				node.Cond = expr(&tok, tok)
			}
			tok = skip(tok, ";")

			if !equal(tok, "{") {
				node.Inc = expr(&tok, tok)
			}

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

	if equal(tok, "break") {
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

	if isMultiAssign(tok) {
		return assignList(rest, tok)
	}

	if equal(tok, "{") {
		return compoundStmt(rest, tok.Next)
	}

	return exprStmt(rest, tok)
}

func isMultiAssign(tok *Token) bool {
	hasComma := false
	hasEqual := false

	for !equal(tok, ";") {
		if equal(tok, "{") {
			return false
		}
		if tok.Kind == TK_STR || tok.Kind == TK_NUM {
			tok = tok.Next
		}
		if equal(tok, ",") {
			hasComma = true
		}
		if equal(tok, "=") {
			hasEqual = true
			break
		}

		tok = tok.Next
	}

	return hasComma && hasEqual
}

func isShortVarSpec(tok *Token) bool {
	for !equal(tok, ";") {
		if equal(tok, ":=") {
			return true
		}

		if equal(tok, ",") && (tok.Next.Kind == TK_IDENT ||
			tok.Next.Kind == TK_BLANKIDENT) {
			tok = tok.Next.Next
			continue
		}
		break
	}
	return false
}

func countRetTys(v *Obj) int {
	i := 0
	for t := v.Ty.RetTy; t != nil; t = t.Next {
		i++
	}
	return i
}

// assign-list = Expr-list "=" Expr-list
func assignList(rest **Token, tok *Token) *Node {
	printCurTok(tok)
	printCalledFunc()

	start := tok
	i := 0

	lhses := &Node{}
	cur := lhses
	for !equal(tok, "=") {
		if i > 0 {
			tok = skip(tok, ",")
		}
		i++
		cur.Next = logor(&tok, tok)
		cur = cur.Next
		addType(cur)
	}

	tok = skip(tok, "=")

	var node *Node
	rhses := &Node{}
	cur = rhses
	lhs := lhses.Next
	valtok := tok
	j := 0
	for ; ; j++ {
		if consume(&tok, tok, ";") {
			break
		}
		if j > 0 {
			tok = skip(tok, ",")
		}

		rhs := logor(&tok, tok)

		if rhs.Kind == ND_FUNCALL {
			numVals := countRetTys(rhs.Lhs.Obj)
			if numVals > 1 {
				// とりあえずlhsesの長さだけで判定、エラーも適当
				if i != numVals {
					panic(errorTok(tok, "too many assigns: left:%d, right:%d", i, numVals))
				}
				node = newUnary(ND_MULTIRETASSIGN, rhs, tok)
				node.Masg = lhses.Next
				*rest = tok.Next
				addType(node)
				return node
			}
		}

		if lhs.Kind == ND_BLANKIDENT {
			lhs = lhs.Next
			continue
		}

		cur.Next = rhs
		addType(rhs)
		if isAppend || isMake {
			lhs.Ty.Len = rhs.Ty.Len
			lhs.Ty.Cap = rhs.Ty.Cap
			lhs.Ty.UArrNode = rhs.Ty.UArrNode
			lhs.Ty.UArrIdx = rhs.Ty.UArrIdx
		}
		if lhs.Obj != nil {
			if isInteger(lhs.Ty) {
				lhs.Obj.Val = eval(rhs)
			} else if isFlonum(lhs.Ty) {
				lhs.Obj.FVal = evalDouble(rhs)
			}
		}

		cur = cur.Next
		lhs = lhs.Next
	}

	if j > i {
		panic("\n" + errorTok(valtok,
			"assignment mismatch: %d variables but %d values", i, j))
	}

	node = newNode(ND_MULTIVALASSIGN, start)
	node.Lhses = lhses.Next
	node.Rhses = rhses.Next

	*rest = tok
	return node
}

// compound-stmt = (typedef | declaration | stmt)* "}"
func compoundStmt(rest **Token, tok *Token) *Node {
	printCurTok(tok)
	printCalledFunc()

	node := newNode(ND_BLOCK, tok)
	head := &Node{}
	cur := head

	enterScope()

	for !equal(tok, "}") {

		if tok.Kind == TK_COMM {
			tok = tok.Next
			continue
		}

		if consume(&tok, tok, "type") {
			tok = parseTypedef(tok)
			continue
		}

		if equal(tok, "var") && equal(tok.Next, "(") {
			tok = tok.Next.Next

			for !equal(tok, ")") {
				if tok.Kind == TK_COMM {
					// skip line comment
					tok = tok.Next
					continue
				}

				if tok.Kind != TK_IDENT && tok.Kind != TK_BLANKIDENT {
					panic("\n" + errorTok(tok, "unexpected expression"))
				}

				cur.Next = declaration(&tok, tok, false)
				cur = cur.Next
			}
			tok = skip(tok, ")")
			continue
		}

		if consume(&tok, tok, "var") {
			cur.Next = declaration(&tok, tok, false)

		} else if (tok.Kind == TK_IDENT || tok.Kind == TK_BLANKIDENT) && isShortVarSpec(tok.Next) {
			cur.Next = declaration(&tok, tok, true)

		} else {
			cur.Next = stmt(&tok, tok)

			if isAppend {
				cur = cur.Next
				addType(cur)
				cur.Next = appendAsg
				isAppend = false
			}

			isMake = false

		}

		cur = cur.Next
		addType(cur) //
	}
	leaveScope()

	node.Body = head.Next
	*rest = tok.Next
	return node
}

// expr-stmt = expr? ";"
func exprStmt(rest **Token, tok *Token) *Node {
	printCurTok(tok)
	printCalledFunc()

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
	printCurTok(tok)
	printCalledFunc()

	node := assign(&tok, tok)

	if equal(tok, ",") {
		return newBinary(ND_COMMA, node, expr(rest, tok.Next), tok)
	}

	*rest = tok
	return node
}

// Evaluate a given node as a constant expression.
//
// A constant expression is either just a number or ptr+n where ptr
// number. The latter form is accept only as an initialization
// expression for a global variable.
func eval(node *Node) int64 {
	printCalledFunc()

	return eval2(node, nil)
}

func eval2(node *Node, label **string) int64 {
	printCalledFunc()

	if node == nil {
		return 0
	}

	addType(node)

	if isFlonum(node.Ty) {
		fval := evalDouble(node)
		return int64(fval)
	}

	switch node.Kind {
	case ND_ADD:
		return eval2(node.Lhs, label) + eval(node.Rhs)
	case ND_SUB:
		return eval2(node.Lhs, label) - eval(node.Rhs)
	case ND_MUL:
		return eval(node.Lhs) * eval(node.Rhs)
	case ND_DIV:
		if node.Ty.IsUnsigned {
			return int64(uint64(eval(node.Lhs))) / eval(node.Rhs)
		}
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
		if node.Ty.IsUnsigned && node.Ty.Sz == 8 {
			return int64(uint64(eval(node.Lhs))) >> eval(node.Rhs)
		}
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
		if node.Lhs.Ty.IsUnsigned &&
			int64(uint64(eval(node.Lhs))) < eval(node.Rhs) {
			return 1
		}
		if eval(node.Lhs) < eval(node.Rhs) {
			return 1
		}
		return 0
	case ND_LE:
		if node.Lhs.Ty.IsUnsigned &&
			int64(uint64(eval(node.Lhs))) <= eval(node.Rhs) {
			return 1
		}
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
		return eval2(node.Rhs, label)
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
				if node.Ty.IsUnsigned {
					return int64(uint8(val))
				}
				return int64(int8(val))
			case 2:
				if node.Ty.IsUnsigned {
					return int64(uint16(val))
				}
				return int64(int16(val))
			case 4:
				if node.Ty.IsUnsigned {
					return int64(uint32(val))
				}
				return int64(int32(val))
			}
		}
		return val // If node.Ty.Sz is 8
	case ND_ADDR:
		return evalRval(node.Lhs, label)
	case ND_MEMBER:
		if label == nil {
			panic("\n" + errorTok(node.Tok, "not a compile-time constant"))
		}
		if node.Ty.Kind != TY_ARRAY {
			panic("\n" + errorTok(node.Tok, "invalid initializer"))
		}
		return evalRval(node.Lhs, label) + int64(node.Mem.Offset)
	case ND_VAR:
		if isInteger(node.Ty) {
			return node.Obj.Val
		}
		if label == nil {
			return 0
			// panic("\n" + errorTok(node.Tok, "not a compile-time constant"))
		}
		if node.Obj.Ty.Kind != TY_ARRAY && node.Obj.Ty.Kind != TY_FUNC {
			panic("\n" + errorTok(node.Tok, "invalid initializer"))
		}
		*label = &node.Obj.Name
		return 0
	case ND_NUM:
		return node.Val
	case ND_FUNCALL:
		switch res := evalFuncall(node).(type) {
		case int64:
			return res
		case int:
			return int64(res)
		default:
			panic("\n" + errorTok(node.Tok, "not a compile-time constant"))
		}
	default:
		return 0
		// panic("\n" + errorTok(node.Tok, "not a compile-time constant"))
	}
}

func evalRval(node *Node, label **string) int64 {
	printCalledFunc()

	switch node.Kind {
	case ND_VAR:
		if node.Obj.IsLocal {
			if isInteger(node.Ty) {
				return node.Obj.Val
			}
			panic("\n" + errorTok(node.Tok, "not a compile-time constant"))
		}
		*label = &node.Obj.Name
		return 0
	case ND_DEREF:
		return eval2(node.Lhs, label)
	case ND_MEMBER:
		return evalRval(node.Lhs, label) + int64(node.Mem.Offset)
	default:
		panic("\n" + errorTok(node.Tok, "invalid initializer"))
	}
}

func evalFuncall(node *Node) interface{} {
	addType(node)

	fn := node.Lhs.Obj
	if fn.IsDef {
		enterScope()
		defer leaveScope()

		createParamLvars(fn.Ty.Params)
		lvars := fn.Locals
		for lv := lvars; lv != nil; lv = lv.Next {
			pushScope(lv.Name).Obj = lv
		}

		// evaluate arguments
		for arg := node.Args; lvars != nil && node.Args != nil; arg = arg.Next {
			findVar(lvars.Ty.Name).Obj.Val = eval(arg)
			lvars = lvars.Next
		}

		node2 := fn.Body
		for n := node2.Body; n != nil; n = n.Next {
			if n.Kind == ND_RETURN {
				if isInteger(n.RetVals.Ty) {
					return eval(n.RetVals)
				} else if isFlonum(n.RetVals.Ty) {
					return evalDouble(n.RetVals)
				}
				panic("\n" + errorTok(node.Tok, "not a compile-time constant"))
			}
		}
	}
	return 0
}

// const-expr
func constExpr(rest **Token, tok *Token) int64 {
	printCurTok(tok)
	printCalledFunc()

	return eval(logor(rest, tok))
}

func evalDouble(node *Node) float64 {
	addType(node)

	if isInteger(node.Ty) {
		if node.Ty.IsUnsigned {
			return float64(uint64(eval(node)))
		}
		return float64(eval(node))
	}

	switch node.Kind {
	case ND_ADD:
		return evalDouble(node.Lhs) + evalDouble(node.Rhs)
	case ND_SUB:
		return evalDouble(node.Lhs) - evalDouble(node.Rhs)
	case ND_MUL:
		return evalDouble(node.Lhs) * evalDouble(node.Rhs)
	case ND_DIV:
		return evalDouble(node.Lhs) / evalDouble(node.Rhs)
	case ND_NEG:
		return -evalDouble(node.Lhs)
	case ND_COND:
		if evalDouble(node.Cond) != 0 {
			return evalDouble(node.Then)
		}
		return evalDouble(node.Els)
	case ND_COMMA:
		return evalDouble(node.Rhs)
	case ND_CAST:
		if isFlonum(node.Lhs.Ty) {
			return evalDouble(node.Lhs)
		}
		return float64(eval(node.Lhs))
	case ND_NUM:
		return node.FVal
	case ND_VAR:
		return node.Obj.FVal
	case ND_FUNCALL:
		switch res := evalFuncall(node).(type) {
		case float64:
			return res
		case float32:
			return float64(res)
		default:
			panic("\n" + errorTok(node.Tok, "not a compile-time constant"))
		}
	default:
		panic("\n" + errorTok(node.Tok, "not a complie-time constant"))
	}
}

// Convert `A op= B` to `tmp = &A, *tmp = *tmp op B`
// where tmp is a fresh pointer variable.
func toAssign(binary *Node) *Node {
	printCalledFunc()

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
	printCurTok(tok)
	printCalledFunc()

	node := logor(&tok, tok)

	if equal(tok, "=") {
		rhs := assign(rest, tok.Next)
		addType(rhs)
		if node.Kind == ND_BLANKIDENT {
			return rhs
		}
		if node.Obj != nil {
			if isInteger(rhs.Ty) {
				node.Obj.Val = eval(rhs)
			} else if isFlonum(rhs.Ty) {
				node.Obj.FVal = evalDouble(rhs)
			}
		}

		if isAppend || isMake {
			addType(node)
			node.Ty.Len = rhs.Ty.Len
			node.Ty.Cap = rhs.Ty.Cap
			node.Ty.UArrNode = rhs.Ty.UArrNode
			node.Ty.UArrIdx = rhs.Ty.UArrIdx
		}

		node = newBinary(ND_ASSIGN, node, rhs, tok)
		addType(node)
		return node
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
		return toAssign(newBinary(ND_DIV, node, assign(rest, tok.Next), tok))
	}

	if equal(tok, "%=") {
		return toAssign(newBinary(ND_MOD, node, assign(rest, tok.Next), tok))
	}

	if equal(tok, "&=") {
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
	printCurTok(tok)
	printCalledFunc()

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
	printCurTok(tok)
	printCalledFunc()

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
	printCurTok(tok)
	printCalledFunc()

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
	printCurTok(tok)
	printCalledFunc()

	node := bitand(&tok, tok)
	for equal(tok, "^") {
		start := tok
		node = newBinary(ND_BITXOR, node, bitand(&tok, tok.Next), start)
	}
	*rest = tok
	return node
}

// bitand = equality ("&" equality)*
func bitand(rest **Token, tok *Token) *Node {
	printCurTok(tok)
	printCalledFunc()

	node := equality(&tok, tok)
	for equal(tok, "&") {
		start := tok
		node = newBinary(ND_BITAND, node, equality(&tok, tok.Next), start)
	}
	*rest = tok
	return node
}

// equality   = relational ("==" relational | "!=" relational)*
// Comparing strings is unimplemented yet.
func equality(rest **Token, tok *Token) *Node {
	printCurTok(tok)
	printCalledFunc()

	node := relational(&tok, tok)

	for {
		start := tok

		if equal(tok, "==") {
			node = newBinary(ND_EQ, node, relational(&tok, tok.Next), start)
			continue
		}

		if equal(tok, "!=") {
			node = newBinary(ND_NE, node, relational(&tok, tok.Next), start)
			continue
		}

		*rest = tok
		return node
	}
}

// relational = shift ("<" shift | "<=" shift | ">" shift | ">=" shift)*
func relational(rest **Token, tok *Token) *Node {
	printCurTok(tok)
	printCalledFunc()

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
			node = newBinary(ND_LT, shift(&tok, tok.Next), node, start)
			continue
		}

		if equal(tok, ">=") {
			node = newBinary(ND_LE, shift(&tok, tok.Next), node, start)
			continue
		}

		*rest = tok
		return node
	}
}

// shift = add ("<<" add | ">>" add)*
func shift(rest **Token, tok *Token) *Node {
	printCurTok(tok)
	printCalledFunc()

	node := add(&tok, tok)

	for {
		start := tok

		if equal(tok, "<<") {
			node = newBinary(ND_SHL, node, add(&tok, tok.Next), start)
			continue
		}

		if equal(tok, ">>") {
			node = newBinary(ND_SHR, node, add(&tok, tok.Next), start)
			continue
		}

		*rest = tok
		return node
	}
}

// newAdd :
// In C, `+` operator is overloaded to perform the pointer arithmetic.
// If p is a pointer, p+n add not n but sizeof(*p)*n to the value of p,
// sothat p+n pointes to the location n elements (not bytes) ahead of p.
// In other words, we need to scale an integer value before adding to a
// pointer value. This function takes care of the scaling.
func newAdd(lhs, rhs *Node, tok *Token) *Node {
	printCurTok(tok)
	printCalledFunc()

	addType(lhs)
	addType(rhs)

	// num + num
	if isNumeric(lhs.Ty) && isNumeric(rhs.Ty) {
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
func newSub(lhs, rhs *Node, tok *Token) *Node {
	printCurTok(tok)
	printCalledFunc()

	addType(lhs)
	addType(rhs)

	// num - num
	if isNumeric(lhs.Ty) && isNumeric(rhs.Ty) {
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
		node.Ty = ty_long
		return newBinary(ND_DIV, node, newNum(int64(lhs.Ty.Base.Sz), tok), tok)
	}

	panic("\n" + errorTok(tok, "invalud operands"))
}

// add        = mul ("+" mul | "-" mul)*
func add(rest **Token, tok *Token) *Node {
	printCurTok(tok)
	printCalledFunc()

	node := mul(&tok, tok)

	for {
		start := tok

		if equal(tok, "+") {
			node = newAdd(node, mul(&tok, tok.Next), start)
			continue
		}

		if equal(tok, "-") {
			node = newSub(node, mul(&tok, tok.Next), start)
			continue
		}

		*rest = tok
		return node
	}
}

// mul = cast ("*" cast | "/" cast)*
func mul(rest **Token, tok *Token) *Node {
	printCurTok(tok)
	printCalledFunc()

	node := cast(&tok, tok)

	for {
		start := tok

		if equal(tok, "*") {
			node = newBinary(ND_MUL, node, cast(&tok, tok.Next), start)
			continue
		}

		if equal(tok, "/") {
			node = newBinary(ND_DIV, node, cast(&tok, tok.Next), start)
			continue
		}

		if equal(tok, "%") {
			node = newBinary(ND_MOD, node, cast(&tok, tok.Next), start)
		}

		*rest = tok
		return node
	}
}

// cast = type-name "(" cast ")"
//      | unary
func cast(rest **Token, tok *Token) *Node {
	printCurTok(tok)
	printCalledFunc()

	start := tok
	if isTypename(tok) {
		ty := readTypePreffix(&tok, tok, nil)

		// conmpound literal
		if equal(tok, "{") {
			return unary(rest, start)
		}

		node := newCast(cast(&tok, tok), ty)
		node.Tok = start
		*rest = tok
		return node
	}

	return unary(rest, tok)
}

// unary   = ("+" | "-" | "*" | "&" | "!")? cast
//         | postfix
func unary(rest **Token, tok *Token) *Node {
	printCurTok(tok)
	printCalledFunc()

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

// func-decl = "func(" param-type* ")" ( return-type | "(" return-type ("," return-type)* ")" )
func funcDecl(rest **Token, tok *Token, name *Token) *Type {
	printCurTok(tok)
	printCalledFunc()

	tok = skip(tok, "(")

	head := &Type{}
	cur := head
	first := true

	// Get parameters
	for !equal(tok, ")") {
		if !first {
			tok = skip(tok, ",")
		}
		first = false

		paramty := readTypePreffix(&tok, tok, name)
		if paramty == nil {
			panic("\n" + errorTok(tok, "is not typename"))
		}
		cur.Next = copyType(paramty)
		cur = cur.Next
	}
	tok = skip(tok, ")")

	var retty *Type
	if equal(tok, "(") {
		tok = skip(tok, "(")
		head := &Type{}
		cur := head
		first := true
		for !equal(tok, ")") {
			if !first {
				tok = skip(tok, ",")
			}
			first = false
			ret := readTypePreffix(&tok, tok, nil)
			cur.Next = copyType(ret)
			cur = cur.Next
		}
		cur.Next = nil
		retty = head.Next
		tok = skip(tok, ")")
	} else {
		retty = readTypePreffix(&tok, tok, nil)
	}

	*rest = tok

	ty := pointerTo(funcType(retty, head.Next))

	if name != nil {
		pushScope(getIdent(name)).TyDef = ty
	} else {
		pushScope(newUniqueName()).TyDef = ty
	}

	return ty
}

// struct-member = (ident | ident-list) type-prefix type-specifier
func structMems(rest **Token, tok *Token, ty *Type) *Member {
	printCurTok(tok)
	printCalledFunc()

	memList := make([]*Type, 0)

	for !equal(tok, "}") {
		first := true
		for !consume(&tok, tok, ";") {
			if !first {
				tok = skip(tok, ",")
			}
			first = false
			memTy := declarator(&tok, tok)
			memList = append(memList, copyType(memTy))
			if memTy.Kind != TY_VOID {
				idx := len(memList) - 2
				for 0 <= idx && memList[idx].Kind == TY_VOID {
					memTy2 := copyType(memTy)
					name := memList[idx].Name
					memList[idx] = memTy2
					memList[idx].Name = name
					idx--
				}
			}
			if equal(tok, "}") {
				break
			}
		}
	}

	head := &Member{}
	cur := head

	for i := 0; i < len(memList); i++ {
		mem := &Member{
			Name:  memList[i].Name,
			Ty:    memList[i],
			Idx:   i,
			Align: memList[i].Align,
		}
		cur.Next = mem
		cur = cur.Next
	}

	// If the last element is an array of imcomlete type, it's
	// called a "flexible array member". It should bahave as if
	// if were a zero-sized array.
	for cur != head && cur.Ty.Kind == TY_ARRAY && cur.Ty.ArrSz <= 0 {
		cur.Ty = arrayOf(cur.Ty.Base, 0)
		ty.IsFlex = true
	}

	*rest = tok.Next
	return head.Next
}

// struct-decl = "struct" "{" struct-member "}"
func structDecl(rest **Token, tok *Token, name *Token) *Type {
	printCurTok(tok)
	printCalledFunc()

	tok = skip(tok, "{")

	// Construct a struct object.
	ty := structType()
	if name != nil {
		pushScope(getIdent(name)).TyDef = ty
	} else {
		pushScope(newUniqueName()).TyDef = ty
	}
	ty.Mems = structMems(rest, tok, ty)

	// Assign offsers within the struct to members.
	offset := 0
	for mem := ty.Mems; mem != nil; mem = mem.Next {
		offset = alignTo(offset, mem.Align)
		mem.Offset = offset
		offset += mem.Ty.Sz

		if ty.Align < mem.Align {
			ty.Align = mem.Align
		}
	}
	ty.Sz = alignTo(offset, ty.Align)
	return ty
}

func getStructMember(ty *Type, tok *Token) *Member {
	printCurTok(tok)
	printCalledFunc()

	if ty.Kind != TY_STRUCT {
		for ty != nil && ty.Base != nil {
			ty = ty.Base
		}
	}

	for mem := ty.Mems; mem != nil; mem = mem.Next {
		if mem.Name.Str == tok.Str {
			return mem
		}
	}
	panic("\n" + errorTok(tok, "no such member"))
}

func structRef(lhs *Node, tok *Token) *Node {
	printCurTok(tok)
	printCalledFunc()

	addType(lhs)

	if lhs.Ty.Kind != TY_STRUCT {
		if lhs.Ty.Base != nil && lhs.Ty.Base.Kind != TY_STRUCT {
			panic("\n" + errorTok(lhs.Tok, "not a struct"))
		}
		// "->" in C.
		lhs = newUnary(ND_DEREF, lhs, tok)
		addType(lhs)
	}

	node := newUnary(ND_MEMBER, lhs, tok.Next)
	node.Mem = getStructMember(lhs.Ty, tok.Next)
	return node
}

// Convert A++ to `(typeof A)((A += 1) - 1)`
func newIncDec(node *Node, tok *Token, addend int) *Node {
	printCurTok(tok)
	printCalledFunc()

	addType(node)
	return newCast(newAdd(toAssign(newAdd(node, newNum(int64(addend), tok), tok)),
		newNum(int64(addend)*-1, tok), tok),
		node.Ty)
}

// slice-expr = primary "[" expr ":" const-expr "]"
func sliceExpr(rest **Token, tok *Token, cur *Node, idx *Node, start *Token) *Node {
	first := eval(idx)

	var end int64
	if equal(tok.Next, "]") {
		switch cur.Obj.Ty.Kind {
		case TY_ARRAY:
			end = int64(cur.Obj.Ty.ArrSz)
		case TY_SLICE:
			end = int64(cur.Obj.Ty.Len)
		}
		*rest = tok.Next
	} else {
		end = constExpr(rest, tok.Next)
	}

	node := newUnary(ND_ADDR,
		newUnary(ND_DEREF, newAdd(cur, idx, start), start), start)
	addType(node)

	len := int(end - first)
	var cap int
	switch cur.Obj.Ty.Kind {
	case TY_ARRAY:
		cap = cur.Obj.Ty.ArrSz - int(first)
	case TY_SLICE:
		cap = cur.Obj.Ty.Cap - int(first)
	default:
		panic(errorTok(start, "is not neither array nor slice"))
	}

	node.Ty = sliceType(node.Ty.Base, len, cap)
	node.Ty.UArrIdx = first
	node.Ty.UArrNode = cur

	return node
}

// postfix = "(" type-name ")" "{" initializer-list "}"
//         | primary postfix-tail*
//
// postfix-tail = "[" expr "]"
//              | "(" func-args ")"
//              | slice-expr
//              | "." ident
//              | "++"
//              | "--"
func postfix(rest **Token, tok *Token) *Node {
	printCurTok(tok)
	printCalledFunc()

	start := tok
	if isTypename(tok) {
		// Compound literal : type-name "{"
		ty := readTypePreffix(&tok, tok, nil)
		if scope.Next == nil {
			v := newAnonGvar(ty)
			gvarInitializer(rest, tok, v)
			return newVarNode(v, start)
		}

		v := newLvar("", ty)
		lhs := lvarInitializer(rest, tok, v)
		rhs := newVarNode(v, tok)
		return newBinary(ND_COMMA, lhs, rhs, start)
	}

	node := primary(&tok, tok)

	for {
		if equal(tok, "(") {
			node = funcall(&tok, tok.Next, node)
			continue
		}

		if equal(tok, "[") {
			// x[y:z] is slice
			start := tok
			idx := expr(&tok, tok.Next)
			if equal(tok, ":") {
				node = sliceExpr(&tok, tok, node, idx, start)
				tok = skip(tok, "]")
				*rest = tok
				return node
			}
			i := eval(idx)
			if i < 0 {
				panic(errorTok(idx.Tok,
					"invalid argument: index %d (constant of type int) must not be negative", i))
			}
			addType(node)
			if node.Ty != nil {
				if node.Ty.Kind == TY_ARRAY && i >= int64(node.Ty.ArrSz) {
					panic(errorTok(idx.Tok, "index out of range [%d] with length %d", i, node.Ty.ArrSz))
				}
				if node.Ty.Kind == TY_SLICE && i >= int64(node.Ty.Len) && !node.Ty.IsFlex {
					panic(errorTok(idx.Tok, "index out of range [%d] with length %d", i, node.Ty.Len))
				}
			}
			tok = skip(tok, "]")
			// x[y] is short for *(x+y)
			node = newUnary(ND_DEREF, newAdd(node, idx, start), start)
			addType(node)
			continue
		}

		if equal(tok, ".") {
			node = structRef(node, tok)
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

// funcall = "(" (assign ("," assign)*)? ")"
//
//
func funcall(rest **Token, tok *Token, fn *Node) *Node {
	printCurTok(tok)
	printCalledFunc()

	addType(fn)

	if fn.Ty.Kind != TY_FUNC &&
		(fn.Ty.Kind != TY_PTR || fn.Ty.Base.Kind != TY_FUNC) {
		panic(errorTok(fn.Tok, "not a function"))
	}

	var ty *Type
	if fn.Ty.Kind == TY_FUNC {
		ty = fn.Ty
	} else {
		ty = fn.Ty.Base
	}
	paramTy := ty.Params

	head := &Node{}
	cur := head

	for !equal(tok, ")") {
		if cur != head {
			tok = skip(tok, ",")
		}

		arg := assign(&tok, tok)
		addType(arg)

		if paramTy == nil && !ty.IsVariadic {
			panic("\n" + errorTok(tok, "too many arguments"))
		}

		if paramTy != nil {
			if paramTy.Kind != TY_STRUCT {
				arg = newCast(arg, paramTy)
			}
			paramTy = paramTy.Next
		} else if arg.Ty.Kind == TY_FLOAT {
			// If parameter type is omitted (e.g. in "..."), float
			// arguments are promoted to double.
			arg = newCast(arg, ty_double)
		}

		cur.Next = arg
		cur = cur.Next
	}

	if paramTy != nil {
		panic("\n" + errorTok(tok, "too few arguments"))
	}

	*rest = skip(tok, ")")

	// fmt.Printf("funcall: ty: %#v\n\n", ty)

	node := newUnary(ND_FUNCALL, fn, tok)
	node.FuncTy = ty
	node.Ty = ty.RetTy
	node.Args = head.Next

	// If a function returns a struct, it is caller's responsibility
	// to allocate a space for the return value.
	vhead := &Obj{}
	vcur := vhead
	for r := ty.RetTy; r != nil; r = r.Next {
		// fmt.Printf("funcall: r: %#v\n\n", r)
		if r.Kind == TY_STRUCT {
			vcur.RetNext = newLvar(newFavName("retbuf"), r)
			vcur = vcur.RetNext
		}
	}
	node.RetBuf = vhead.RetNext
	return node
}

var isMake bool
var isAppend bool
var appendAsg *Node

func countAppElem(tok *Token) int {
	ret := 0
	for !equal(tok, ")") {
		tok = skip(tok, ",")
		assign(&tok, tok)
		ret++
	}
	return ret
}

// primary = "(" expr ")"
//         | "Sizeof" "(" type-name ")"
//         | "Sizeof" unary
//         | "len" unary
//         | "cap" unary
//         | "make" "(" type-name "," const-expr "," const-expr ")"
//         | "append" "(" postfix "," assign ( "," assign)* ")"
//         | "copy" "(" assign "," assign ")"
//         | ident
//         | str
//         | num
func primary(rest **Token, tok *Token) *Node {
	printCurTok(tok)
	printCalledFunc()

	start := tok

	// if the next token is '(', the program must be
	// "(" expr ")"
	if equal(tok, "(") {
		node := expr(&tok, tok.Next)
		*rest = skip(tok, ")")
		return node
	}

	if equal(tok, "Sizeof") && equal(tok.Next, "(") &&
		isTypename(tok.Next.Next) && !equal(tok.Next.Next.Next, "(") {
		ty := readTypePreffix(&tok, tok.Next.Next, nil)
		*rest = skip(tok, ")")
		return newUlong(int64(ty.Sz), start)
	}

	if equal(tok, "Sizeof") && equal(tok.Next, "(") {
		// "(" 以降のtokenをunaryに渡して、
		// unary -> postfix -> primaryの"(" expr ")"でparseする
		node := unary(rest, tok.Next)
		addType(node)
		return newUlong(int64(node.Ty.Sz), tok)
	}

	if equal(tok, "Alignof") {
		node := unary(rest, tok.Next)
		addType(node)
		return newUlong(int64(node.Ty.Align), tok)
	}

	if equal(tok, "len") && equal(tok.Next, "(") {
		node := unary(rest, tok.Next)
		addType(node)
		return newUlong(int64(node.Ty.Len), tok)
	}

	if equal(tok, "cap") && equal(tok.Next, "(") {
		node := unary(rest, tok.Next)
		addType(node)
		return newUlong(int64(node.Ty.Cap), tok)
	}

	// 'make' function for slice only.
	if equal(tok, "make") && equal(tok.Next, "(") {
		isMake = true
		start := tok
		ty := readTypePreffix(&tok, tok.Next.Next, nil)
		tok = skip(tok, ",")
		len := constExpr(&tok, tok)
		var cap int64
		if equal(tok, ")") {
			ty.Len = int(len)
			ty.Cap = int(len)
			cap = len
		} else {
			tok = skip(tok, ",")
			cap = constExpr(&tok, tok)
			ty.Len = int(len)
			ty.Cap = int(cap)
		}
		// Make the underlying array.
		uArr := newFavGvar("underlying_array", arrayOf(ty.Base, int(cap)))
		cnt++
		gvarZeroInit(uArr, tok)
		uaNode := newVarNode(uArr, start)
		ty.UArrNode = uaNode

		node := newUnary(ND_ADDR,
			newUnary(ND_DEREF,
				newAdd(uaNode, newNum(0, start), start), start), start)
		addType(node)
		node.Ty = ty
		*rest = skip(tok, ")")
		return node
	}

	// 'append' function
	if equal(tok, "append") && equal(tok.Next, "(") {
		tok = skip(tok.Next, "(")
		slice := postfix(&tok, tok)
		addType(slice)
		if slice.Ty.Kind != TY_SLICE {
			panic(errorTok(
				tok,
				"first argument to append must be a slice; have a (variable of type %s)",
				slice.Ty.TyName))
		}
		isAppend = true
		cntElem := countAppElem(tok)

		// In the case that the new length is no more than the slice's capacity.
		if slice.Ty.Len+cntElem <= slice.Ty.Cap {
			head := &Node{}
			cur := head
			for !equal(tok, ")") {
				tok = skip(tok, ",")
				elem := assign(&tok, tok)

				// Assign elem to slice[slice.Obj.Ty.Len]
				expr := newBinary(ND_ASSIGN,
					newUnary(ND_DEREF,
						newAdd(slice, newNum(int64(slice.Ty.Len), tok), tok), tok),
					elem, tok)
				slice.Ty.Len++
				cur.Next = newUnary(ND_EXPR_STMT, expr, tok)
				cur = cur.Next
			}
			appendAsg = newNode(ND_BLOCK, tok)
			appendAsg.Body = head.Next

			node := newUnary(ND_ADDR,
				newUnary(ND_DEREF,
					newAdd(slice, newNum(0, tok), tok), tok), tok)
			addType(node)
			node.Ty = slice.Ty
			*rest = skip(tok, ")")
			return node
		}

		// In the case that the new length is more than original slice's capacity,
		// Make a new underlying array.
		uArrTy := arrayOf(slice.Ty.Base, slice.Ty.Cap*2+cntElem)
		uArr := newFavGvar("underlying_array", uArrTy)
		cnt++
		gvarZeroInit(uArr, tok)
		uaNode := newVarNode(uArr, tok)
		addType(uaNode)

		head := &Node{}
		cur := head
		length := slice.Ty.Len
		// Copy to new array from the original underlying array.
		var i int64
		for i = 0; i < int64(length); i++ {
			lhs := newUnary(ND_DEREF,
				newAdd(uaNode, newNum(i, tok), tok), tok)
			addType(lhs)
			// 右辺がnewUnary(ND_DEREF, newAdd(slice, newNum(i, tok), tok), tok)だと上手く代入できない
			rhs := newUnary(ND_DEREF,
				newAdd(slice.Ty.UArrNode, newNum(slice.Ty.UArrIdx+i, tok), tok), tok)
			addType(rhs)
			expr := newBinary(ND_ASSIGN, lhs, rhs, tok)
			cur.Next = newUnary(ND_EXPR_STMT, expr, tok)
			cur = cur.Next
		}

		for !equal(tok, ")") {
			tok = skip(tok, ",")
			elem := assign(&tok, tok)
			addType(elem)

			lhs := newUnary(ND_DEREF,
				newAdd(uaNode, newNum(int64(length), tok), tok), tok)
			addType(lhs)
			// Assign elem to slice[slice.Obj.Ty.Len]
			expr := newBinary(ND_ASSIGN, lhs, elem, tok)
			length++
			addType(expr)
			cur.Next = newUnary(ND_EXPR_STMT, expr, tok)
			cur = cur.Next
		}

		appendAsg = newNode(ND_BLOCK, tok)
		appendAsg.Body = head.Next
		addType(appendAsg)

		node := newUnary(ND_ADDR,
			newUnary(ND_DEREF, newAdd(uaNode, newNum(0, tok), tok), tok), tok)
		node.Ty = sliceType(uArrTy.Base, length, uArrTy.ArrSz)
		node.Ty.UArrNode = uaNode
		*rest = skip(tok, ")")
		return node
	}

	// 'copy' function for slice only.
	if equal(tok, "copy") && equal(tok.Next, "(") {
		start := tok
		tok = skip(tok.Next, "(")

		dst := assign(&tok, tok)
		addType(dst)

		tok = skip(tok, ",")

		src := assign(&tok, tok)
		addType(src)

		// Get the new length of 'dst'.
		newLen := int64(min(dst.Ty.Len, src.Ty.Len))

		// Copy values in 'src'.
		isAppend = true
		head := &Node{}
		cur := head
		var i int64
		for i = 0; i < newLen; i++ {
			lhs := newUnary(ND_DEREF,
				newAdd(dst, newNum(i, tok), tok), tok)
			rhs := newUnary(ND_DEREF,
				newAdd(src, newNum(i, tok), tok), tok)
			expr := newBinary(ND_ASSIGN, lhs, rhs, tok)
			cur.Next = newUnary(ND_EXPR_STMT, expr, tok)
			cur = cur.Next
		}
		appendAsg = newNode(ND_BLOCK, tok)
		appendAsg.Body = head.Next
		addType(appendAsg)

		*rest = tok.Next
		return newNum(newLen, start)
	}

	if tok.Kind == TK_BLANKIDENT {
		*rest = tok.Next
		return newNode(ND_BLANKIDENT, tok)
	}

	if tok.Kind == TK_IDENT {
		// Variable
		sc := findVar(tok)
		*rest = tok.Next

		isSVC := isShortVarSpec(tok.Next)
		if sc != nil && sc.Obj != nil {
			if isSVC {
				panic(errorTok(tok, "no new variables on left side of :="))
			}
			return newVarNode(sc.Obj, tok)
		}

		if isSVC {
			return declaration(rest, tok, true)
		}

		if equal(tok.Next, "(") {
			panic(errorTok(tok, "implicit declaration of a function"))
		}

		panic(errorTok(tok, "undefined variable"))
	}

	if tok.Kind == TK_STR {
		v := newStringLiteral([]int64(tok.Contents), tok.Ty)
		*rest = tok.Next
		return newVarNode(v, tok)
	}

	if tok.Kind == TK_NUM {
		var node *Node
		if isFlonum(tok.Ty) {
			node = newNode(ND_NUM, tok)
			node.FVal = tok.FVal
		} else {
			node = newNum(tok.Val, tok)
		}

		node.Ty = tok.Ty
		*rest = tok.Next
		return node
	}

	panic("\n" + errorTok(tok, "expected expression: %s", tok.Str))
}

// typedef = "type" ident (type-preffix)? decl-spec
func parseTypedef(tok *Token) *Token {
	printCurTok(tok)
	printCalledFunc()

	first := true

	for !consume(&tok, tok, ";") {
		if !first {
			tok = skip(tok, ",")
		}
		first = false

		ty := declarator(&tok, tok)
		if ty.Name == nil {
			panic("\n" + errorTok(ty.NamePos, "typedef name omitted"))
		}
		if ty.Kind != TY_STRUCT {
			pushScope(getIdent(ty.Name)).TyDef = ty
		}
	}
	return tok
}

func createParamLvars(param *Type) {
	if param != nil {
		createParamLvars(param.Next)
		if param.Name == nil {
			panic("\n" + errorTok(param.NamePos, "parameter name omitted"))
		}
		newLvar(getIdent(param.Name), param)
	}
}

// resolveGotoLabels function matches gotos with labels.
//
// We cannot resolve gotos as we parse a function because gotos
// can refer a label that apears later in the function.
// So, we need to do this after we parse the entire function.
func resolveGotoLabels() {
	printCalledFunc()

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
	printCurTok(tok)
	printCalledFunc()

	ty := declarator(&tok, tok)
	if ty.Name == nil {
		panic("\n" + errorTok(ty.NamePos, "function name omitted"))
	}

	var retTy *Type
	if consume(&tok, tok, "(") {
		head := &Type{}
		cur := head
		first := true
		for !consume(&tok, tok, ")") {
			if !first {
				tok = skip(tok, ",")
			}
			first = false
			ret := readTypePreffix(&tok, tok, nil)
			cur.Next = copyType(ret)
			cur = cur.Next
		}
		retTy = head.Next
	} else {
		retTy = readTypePreffix(&tok, tok, nil)
	}

	isvariadic := ty.IsVariadic
	name := ty.Name
	ty = funcType(retTy, ty.Params)
	ty.IsVariadic = isvariadic

	fn := newGvar(getIdent(name), ty)
	fn.IsFunc = true
	fn.IsDef = !consume(&tok, tok, ";")

	if !fn.IsDef {
		return tok
	}

	curFn = fn
	locals = nil
	enterScope()
	createParamLvars(ty.Params)

	// A buffer for a struct return value is passed
	// as the hidden first parameter.
	for rty := ty.RetTy; rty != nil; rty = rty.Next {
		if rty.Kind == TY_STRUCT && rty.Sz > 16 {
			newLvar("", pointerTo(rty))
		}
	}
	fn.Params = locals
	if ty.IsVariadic {
		fn.VaArea = newLvar("__va_area__", arrayOf(ty_char, 136))
	}

	tok = skip(tok, "{")
	fn.Body = compoundStmt(&tok, tok)
	fn.Locals = locals
	leaveScope()
	resolveGotoLabels()
	tok = skip(tok, ";")
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
	printCurTok(tok)
	printCalledFunc()

	var i int

	identList := make([]*Obj, 0)

	for !equal(tok, "=") && !equal(tok, ":=") && !equal(tok, ";") {
		if i > 0 {
			tok = skip(tok, ",")
		}
		i++

		ty := declarator(&tok, tok)
		if ty.Name == nil {
			panic("\n" + errorTok(ty.NamePos, "variable name omitted"))
		}

		v := newGvar(getIdent(ty.Name), ty)
		identList = append(identList, v)
	}

	ty := copyType(identList[len(identList)-1].Ty)
	for j := len(identList) - 2; j >= 0; j-- {
		identList[j].Ty = ty
	}

	if equal(tok, "=") {
		j := 0
		for !equal(tok, ";") {
			v := identList[j]
			gvarInitializer(&tok, tok.Next, v)
			j++
		}

	} else {
		for j := 0; j < len(identList); j++ {
			v := identList[j]
			// Initialize empty variables.
			gvarZeroInit(v, v.Ty.Name)
		}
	}

	tok = skip(tok, ";")
	return tok
}

// program = (global-var | function)*
func parse(tok *Token) *Obj {
	printCurTok(tok)
	printCalledFunc()

	globals = nil

	// package statement 読み飛ばし
	if consume(&tok, tok, "package") {
		tok = tok.Next.Next
	}

	for !atEof(tok) {

		if tok.Kind == TK_COMM {
			tok = tok.Next
			continue
		}

		if consume(&tok, tok, "func") {
			tok = function(tok)
			continue
		}

		if equal(tok, "var") && equal(tok.Next, "(") {
			tok = tok.Next.Next

			for !equal(tok, ")") {
				if tok.Kind == TK_COMM {
					// skip line comment
					tok = tok.Next
					continue
				}

				if tok.Kind != TK_IDENT && tok.Kind == TK_BLANKIDENT {
					panic("\n" + errorTok(tok, "unexpected expression"))
				}
				tok = globalVar(tok)
			}
			tok = skip(tok, ")")
			tok = skip(tok, ";")
			continue
		}

		if consume(&tok, tok, "var") {
			tok = globalVar(tok)
			continue
		}

		if consume(&tok, tok, "type") {
			tok = parseTypedef(tok)
			continue
		}

		panic("\n" + errorTok(tok, "unexpected '%s'", tok.Str))
	}

	return globals
}
