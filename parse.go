//
// AST parser
//
package main

import (
	"os"
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
	ND_LVAR                      // 9: local variables
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
	ND_NULL                      // 20: empty statement
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

	// block
	Body *Node

	// for function call
	FuncName string
	Args     *Node

	Val int   // it would be used when 'Kind' is 'ND_NUM'
	Var *LVar // it would be used when 'Kind' is 'ND_LVAR'
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
type LVar struct {
	Name   string // the name of the variable
	Ty     *Type  // the data type
	Offset int    // the offset from RBP
}

type VarList struct {
	Next *VarList
	Var  *LVar
}

// local variables
var locals *VarList

func newVar(lvar *LVar, tok *Token) *Node {
	return &Node{Kind: ND_LVAR, Tok: tok, Var: lvar}
}

// search a local variable by name.
// if it wasn't find, return nil.
func findLVar(tok *Token) *LVar {
	for vl := locals; vl != nil; vl = vl.Next {
		lvar := vl.Var
		if len(lvar.Name) == tok.Len && tok.Str == lvar.Name {
			return lvar
		}
	}
	return nil
}

func pushVar(name string, ty *Type) *LVar {
	lvar := &LVar{
		Name: name,
		Ty:   ty,
	}

	vl := &VarList{Var: lvar, Next: locals}
	locals = vl
	return lvar
}

type Function struct {
	Next   *Function
	Name   string
	Params *VarList

	Node    *Node
	Locals  *VarList
	StackSz int
}

// code is a slice to store prased nodes.
// var code [100]*Node

// program = function*
func program() *Function {
	// printCurTok()
	// printCurFunc()
	cur := &Function{}
	head := cur

	for !atEof() {
		cur.Next = function()
		cur = cur.Next
	}
	return head.Next
}

// basetype = "int" "*"*
func basetype() *Type {
	expect("int")
	ty := intType()
	for consume("*") != nil {
		ty = pointerTo(ty)
	}
	return ty
}

// param = basetype ident
func readFuncParam() *VarList {
	vl := &VarList{}
	ty := basetype() // 'baseType' function will be booted first.
	vl.Var = pushVar(expectIdent(), ty)
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

// declaration = basetype ident ("=" expr) ";"
func declaration() *Node {
	tok := token
	ty := basetype()
	lvar := pushVar(expectIdent(), ty)

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

		if tok := peek("int"); tok != nil {
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

//unary = "+"? primary
//      | "-"? primary
//      | "*" unary
//      | "&" unary
func unary() *Node {
	// printCurTok()
	// printCurFunc()
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
	return primary()
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

// primary = num
//         | ident func-args?
//         | "(" expr ")"
func primary() *Node {
	// printCurTok()
	// printCurFunc()

	// if the next token is '(', the program must be
	// "(" expr ")"
	if t := consume("("); t != nil {
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
			errorTok(os.Stderr, tok, "undefined variable")
		}
		return newVar(lvar, tok)
	}

	tok := token
	if tok.Kind != TK_NUM {
		errorTok(os.Stderr, tok, "it's not a number")
	}
	// otherwise, must be integer
	return newNodeNum(expectNumber(), tok)
}
