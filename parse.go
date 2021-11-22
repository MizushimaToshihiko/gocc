package main

import "fmt"

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
	IsLocal bool   // local or global

	// Local variables
	Offset int // Offset from RBP

	// Global variables
	Conts   []rune
	ContLen int
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
	ND_VAR                       // 10: local variables
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
	ND_NULL                      // 21: empty statement
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

	Var *Var  // used if kind == ND_VAR
	Val int64 // it would be used when kind is 'ND_NUM'
}

var locals *VarList
var globals *VarList

var varScope *VarScope

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

func pushVar(name string, ty *Type, isLocal bool) *Var {
	// printCurTok()
	// printCalledFunc()

	v := &Var{Name: name, Ty: ty, IsLocal: isLocal}

	var vl *VarList
	if isLocal {
		vl = &VarList{Var: v, Next: locals}
		locals = vl
	} else {
		vl = &VarList{Var: v, Next: globals}
		globals = vl
	}

	pushScope(name).Var = v
	return vl.Var
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

type Function struct {
	Next   *Function
	Name   string
	Params *VarList

	Node    *Node
	Locals  *VarList
	StackSz int
}

func isFunction() bool {
	// printCurTok()
	// printCalledFunc()

	return peek("func") != nil
}

// program = (global-var | function)*
func program() *Program {
	// printCurTok()
	// printCalledFunc()

	head := &Function{}
	cur := head
	globals = nil

	for !atEof() {
		if isFunction() {
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

// basetype = "*"* ("byte" | "int" | struct-decl | typedef-name)
func basetype() *Type {
	// printCurTok()
	// printCalledFunc()

	nPtr := 0
	for consume("*") != nil {
		nPtr++
	}

	if !isTypename() {
		panic(errorTok(token, "typename expected"))
	}

	var ty *Type
	if consume("byte") != nil {
		ty = charType()
	} else if consume("int") != nil {
		ty = intType()
	} else if peek("struct") != nil { // struct type
		ty = structDecl()
	} else {
		ty = findVar(consumeIdent()).TyDef
	}
	if ty == nil {
		panic("\n" + errorTok(token, "'ty' is nil"))
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
	ty := basetype()
	t := token
	token = tok
	return ty, t
}

func readArr(base *Type) *Type {
	// printCurTok()
	// printCalledFunc()

	if consume("[") == nil {
		return base
	}
	sz := expectNumber()
	expect("]")
	base = readArr(base)
	return arrayOf(base, int(sz))
}

func readTypePreffix() *Type {
	// printCurTok()
	// printCalledFunc()

	if peek("[") == nil {
		return basetype()
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
		offset += sizeOf(mem.Ty)

		if ty.Align < mem.Ty.Align {
			ty.Align = mem.Ty.Align
		}
	}

	return ty
}

// struct-member = ident basetype
func structMem() *Member {
	// printCurTok()
	// printCalledFunc()

	mem := &Member{Name: expectIdent(), Ty: readTypePreffix()}
	expect(";")
	return mem
}

// param = ident basetype
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

	name := expectIdent()
	ty := readTypePreffix()
	vl := &VarList{}
	vl.Var = pushVar(name, ty, true)
	return vl
}

// params = ident ("," ident)*
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

// function = "func" ident basetype "(" params? ")" "{" stmt "}"
func function() *Function {
	// printCurTok()
	// printCalledFunc()

	locals = nil

	expect("func")
	fn := &Function{Name: expectIdent()}
	expect("(")
	fn.Params = readFuncParams()
	if isTypename() {
		basetype()
	}
	expect("{")

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

// global-var = "var" ident ("[" num "]")* basetype
//
func globalVar() {
	// printCurTok()
	// printCalledFunc()

	name := expectIdent()
	ty := readTypePreffix()
	expect(";")
	pushVar(name, ty, false)
}

// declaration = "var" ident ("[" num "]")* basetype ("=" expr)
func declaration() *Node {
	// printCurTok()
	// printCalledFunc()

	tok := token
	name := expectIdent()
	ty := readTypePreffix()
	v := pushVar(name, ty, true)

	if consume(";") != nil {
		return newNode(ND_NULL, tok)
	}

	expect("=")
	lhs := newVar(v, tok)
	rhs := expr()
	expect(";")
	node := newBinary(ND_ASSIGN, lhs, rhs, tok)
	return newUnary(ND_EXPR_STMT, node, tok)
}

func readExprStmt() *Node {
	// printCurTok()
	// printCalledFunc()

	t := token
	return newUnary(ND_EXPR_STMT, expr(), t)
}

func isTypename() bool {
	// printCurTok()
	// printCalledFunc()

	return peek("byte") != nil ||
		peek("int") != nil ||
		peek("struct") != nil ||
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
//      | for-stmt
//      | "{" stmt* "}"
//      | "type" ident ("[" num "]")* basetype ";"
//      | declaration
//      | expr ";"
// for-stmt = "for" [ condition ] block .
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

	if t := consume("for"); t != nil {
		if !isForClause() {
			node := newNode(ND_WHILE, t)
			if peek("{") == nil {
				node.Cond = expr()
			} else {
				node.Cond = newNum(1, t)
			}

			node.Then = stmt()
			return node

		} else {
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

	if consume("var") != nil {
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

func assign() *Node {
	// printCurTok()
	// printCalledFunc()

	node := equality()
	if t := consume("="); t != nil {
		node = newBinary(ND_ASSIGN, node, assign(), t)
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
		} else if consume("!=") != nil {
			node = newBinary(ND_NE, node, relational(), t)
		} else {
			return node
		}
	}
}

// relational = add ("<" add | "<=" add | ">" add | ">=" add)*
func relational() *Node {
	// printCurTok()
	// printCalledFunc()

	node := add()

	for {
		if t := consume("<"); t != nil {
			node = newBinary(ND_LT, node, add(), t)
		} else if t := consume("<="); t != nil {
			node = newBinary(ND_LE, node, add(), t)
		} else if t := consume(">"); t != nil {
			node = newBinary(ND_LT, add(), node, t)
		} else if t := consume(">="); t != nil {
			node = newBinary(ND_LE, add(), node, t)
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

// mul = unary ("*" unary | "/" unary)*
func mul() *Node {
	// printCurTok()
	// printCalledFunc()

	node := unary()

	for {
		if t := consume("*"); t != nil {
			node = newBinary(ND_MUL, node, unary(), t)
		} else if t := consume("/"); t != nil {
			node = newBinary(ND_DIV, node, unary(), t)
		} else {
			return node
		}
	}
}

// unary   = ("+" | "-" | "*" | "&")? unary
//         | postfix
func unary() *Node {
	// printCurTok()
	// printCalledFunc()

	if t := consume("+"); t != nil {
		return unary()
	}
	if t := consume("-"); t != nil {
		return newBinary(ND_SUB, newNum(0, t), unary(), t)
	}
	if t := consume("&"); t != nil {
		return newUnary(ND_ADDR, unary(), t)
	}
	if t := consume("*"); t != nil {
		return newUnary(ND_DEREF, unary(), t)
	}
	return postfix()
}

// postfix = primary ("[" expr "]" | "." ident)*
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
			// fmt.Printf("postfix(): node: %#v\n%s\n\n", node, node.Tok.Str)
			// fmt.Printf("postfix(): node.Lhs: %#v\n%s\n\n", node.Lhs, node.Lhs.Tok.Str)
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
			return &Node{Kind: ND_FUNCALL, Tok: t, FuncName: t.Str, Args: funcArgs()}
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
		v := pushVar(newLabel(), ty, false)
		v.Conts = t.Contents
		v.ContLen = t.ContLen
		return newVar(v, t)
	}

	if t.Kind != TK_NUM {
		panic("\n" + errorTok(t, "expected expression"))
	}
	return newNum(expectNumber(), t)
}
