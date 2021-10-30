//
// AST parser
//
package main

import "fmt"

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
	ND_VAR                       // 9: local or global variables
	ND_NUM                       // 10: integer
	ND_RETURN                    // 11: 'return'
	ND_IF                        // 12: "if"
	ND_WHILE                     // 13: "while"
	ND_FOR                       // 14: "for"
	ND_BLOCK                     // 15: {...}
	ND_FUNCCALL                  // 16: function call
	ND_ADDR                      // 17: unary &
	ND_DEREF                     // 18: unary *
	ND_EXPR_STMT                 // 19: expression statement
	ND_STMT_EXPR                 // 20: statement expression
	ND_NULL                      // 21: empty statement
	ND_SIZEOF                    // 22: "sizeof" operator
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

	// for function call
	FuncName string
	Args     *Node

	Val int  // it would be used when 'Kind' is 'ND_NUM'
	Var *Var // it would be used when 'Kind' is 'ND_LVAR'
}

func newNode(kind NodeKind, lhs *Node, rhs *Node, tok *Token) *Node {
	return &Node{
		Kind: kind,
		Lhs:  lhs,
		Rhs:  rhs,
		Tok:  tok,
	}
}

func newNodeNum(val int, tok *Token) *Node {
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
	IsLocal bool   // local or global

	// local variables
	Offset int // the offset from RBP

	// global vaiables
	Contents []byte
	ContLen  int
}

type VarList struct {
	Next *VarList
	Var  *Var
}

// local variables
var locals *VarList
var globals *VarList

func newVar(lvar *Var, tok *Token) *Node {
	return &Node{Kind: ND_VAR, Tok: tok, Var: lvar}
}

// search a local variable by name.
// if it wasn't find, return nil.
func findLVar(tok *Token) *Var {
	for vl := locals; vl != nil; vl = vl.Next {
		lvar := vl.Var
		if len(lvar.Name) == tok.Len && tok.Str == lvar.Name {
			return lvar
		}
	}

	for vl := globals; vl != nil; vl = vl.Next {
		if vl.Var.Name == tok.Str {
			return vl.Var
		}
	}
	return nil
}

func pushVar(name string, ty *Type, isLocal bool) *Var {
	lvar := &Var{
		Name:    name,
		Ty:      ty,
		IsLocal: isLocal,
	}

	vl := &VarList{Var: lvar}

	if isLocal {
		vl.Next = locals
		locals = vl
	} else {
		vl.Next = globals
		globals = vl
	}

	return lvar
}

// for newLabel function
var cnt int

func newLabel() string {
	res := fmt.Sprintf(".L.date.%d", cnt)
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
	basetype()
	t1 := consumeIdent()
	t2 := consume("(")
	isFunc := (t1 != nil) && (t2 != nil)
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
			cur.Next = function()
			cur = cur.Next
		} else {
			globalVar()
		}
	}

	prog := &Program{Globals: globals, Fns: head.Next}
	return prog
}

// basetype = ("int" | "char") "*"*
func basetype() *Type {
	var ty *Type
	if consume("char") != nil {
		ty = charType()
	} else {
		expect("int")
		ty = intType()
	}

	for consume("*") != nil {
		ty = pointerTo(ty)
	}
	return ty
}

func readTypeSuffix(base *Type) *Type {
	if consume("[") == nil {
		return base
	}
	sz := expectNumber()
	expect("]")
	base = readTypeSuffix(base)
	return arrayOf(base, uint16(sz))
}

// param = basetype ident
func readFuncParam() *VarList {
	ty := basetype() // 'baseType' function will be booted first.
	name := expectIdent()
	ty = readTypeSuffix(ty)
	vl := &VarList{Var: pushVar(name, ty, true)}
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

// function = basetype ident "(" params? ")" "{" stmt* "}"
func function() *Function {
	// printCurTok()
	// printCurFunc()
	locals = nil

	basetype()
	fn := &Function{Name: expectIdent()}
	expect("(")
	fn.Params = readFuncParams()
	expect("{")

	cur := &Node{}
	head := cur

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

// global-var = basetype ident ("[" num "]")* ";"
func globalVar() {
	ty := basetype()
	name := expectIdent()
	ty = readTypeSuffix(ty)
	expect(";")
	pushVar(name, ty, false)
}

// declaration = basetype ident ("[" num "]")* ("=" expr) ";"
func declaration() *Node {
	tok := token
	ty := basetype()
	name := expectIdent()
	ty = readTypeSuffix(ty)
	lvar := pushVar(name, ty, true)

	if consume(";") != nil {
		return &Node{Kind: ND_NULL, Tok: tok}
	}

	expect("=")
	lhs := newVar(lvar, tok)
	rhs := expr()
	expect(";")
	node := newNode(ND_ASSIGN, lhs, rhs, tok)
	return newUnary(ND_EXPR_STMT, node, tok)
}

func readExprStmt() *Node {
	tok := token
	return &Node{Kind: ND_EXPR_STMT, Lhs: expr(), Tok: tok}
}

func isTypename() bool {
	return peek("char") != nil || peek("int") != nil
}

// stmt = "return" expr ";"
//      | "if" "(" expr ")" stmt ("else" stmt)?
//      | "while" "(" expr ")" stmt
//      | "for" "(" expr? ";" expr? ";" expr? ")" stmt
//      | "{" stmt* "}"
//      | declaration
//      | expr ";"
func stmt() *Node {
	// printCurTok()
	// printCurFunc()
	var node *Node

	if t := consume("return"); t != nil {

		node = &Node{Kind: ND_RETURN, Lhs: expr(), Tok: t}
		expect(";")

	} else if t := consume("if"); t != nil {

		node = &Node{Kind: ND_IF, Tok: t}
		expect("(")
		node.Cond = expr()
		expect(")")
		node.Then = stmt()

		if consume("else") != nil {
			node.Els = stmt()
		}

	} else if t := consume("while"); t != nil {

		node = &Node{Kind: ND_WHILE, Tok: t}
		expect("(")
		node.Cond = expr()
		expect(")")
		node.Then = stmt()

	} else if t := consume("for"); t != nil {

		node = &Node{Kind: ND_FOR, Tok: t}
		expect("(")

		if consume(";") == nil {
			node.Init = readExprStmt()
			expect(";")
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

	} else if t := consume("{"); t != nil {

		head := Node{}
		cur := &head

		for {
			if consume("}") != nil {
				break
			}
			cur.Next = stmt()
			cur = cur.Next
		}

		node = &Node{Kind: ND_BLOCK, Tok: t}
		node.Body = head.Next

	} else {

		if isTypename() {
			return declaration()
		}

		node = readExprStmt()
		expect(";")
	}

	return node
}

// expr       = assign
func expr() *Node {
	// printCurTok()
	// printCurFunc()
	return assign()
}

// assign     = equality ("=" assign)?
func assign() *Node {
	// printCurTok()
	// printCurFunc()
	node := equality()
	if t := consume("="); t != nil {
		node = newNode(ND_ASSIGN, node, assign(), t)
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

// relational = add ("<" add | "<=" add | ">" add | ">=" add)*
func relational() *Node {
	// printCurTok()
	// printCurFunc()
	node := add()

	for {
		if t := consume("<"); t != nil {
			node = newNode(ND_LT, node, add(), t)
		} else if t := consume("<="); t != nil {
			node = newNode(ND_LE, node, add(), t)
		} else if t := consume(">"); t != nil {
			node = newNode(ND_LT, add(), node, t)
		} else if t := consume(">="); t != nil {
			node = newNode(ND_LE, add(), node, t)
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

// mul = unary ("*" unary | "/" unary)*
func mul() *Node {
	// printCurTok()
	// printCurFunc()
	node := unary()

	for {
		if t := consume("*"); t != nil {
			node = newNode(ND_MUL, node, unary(), t)
		} else if consume("/") != nil {
			node = newNode(ND_DIV, node, unary(), t)
		} else {
			return node
		}
	}
}

// unary = ("+" | "-" | "*" | "&")? unary
//       | "sizeof" unary
//       | postfix
func unary() *Node {
	// printCurTok()
	// printCurFunc()
	if t := consumeSizeof(); t != nil {
		return newUnary(ND_SIZEOF, unary(), t)
	}
	if t := consume("+"); t != nil {
		return unary()
	}
	if t := consume("-"); t != nil {
		return newNode(ND_SUB, newNodeNum(0, t), unary(), t)
	}
	if t := consume("&"); t != nil {
		return newUnary(ND_ADDR, unary(), t)
	}
	if t := consume("*"); t != nil {
		return newUnary(ND_DEREF, unary(), t)
	}
	return postfix()
}

// postfix = primary ("[" expr "]")*
func postfix() *Node {
	node := primary()

	for {
		tok := consume("[")
		if tok == nil {
			break
		}
		// x[y] is short for *(x+y)
		exp := newNode(ND_ADD, node, expr(), tok)
		expect("]")
		node = newUnary(ND_DEREF, exp, tok)
	}
	return node
}

// stmt-expr = "(" "{" stmt stmt* "}" ")"
//
// statement expression is a GNU extension.
func stmtExpr(tok *Token) *Node {
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
		if t := consume("("); t != nil { // function call
			return &Node{
				Kind:     ND_FUNCCALL,
				Tok:      tok,
				FuncName: tok.Str,
				Args:     funcArgs(),
			}
		}

		// local variables
		lvar := findLVar(tok)
		if lvar == nil {
			panic("\n" + errorTok(tok, "undefined variable"))
		}
		return newVar(lvar, tok)
	}

	tok := token
	if tok.Kind == TK_STR {
		token = token.Next

		ty := arrayOf(charType(), uint16(tok.ContLen))
		var_ := pushVar(newLabel(), ty, false)
		var_.Contents = tok.Contents
		var_.ContLen = tok.ContLen
		return newVar(var_, tok)
	}

	if tok.Kind != TK_NUM {
		panic("\n" + errorTok(tok, "expected expression"))
	}
	// otherwise, must be integer
	return newNodeNum(expectNumber(), tok)
}
