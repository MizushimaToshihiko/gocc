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
	ND_ADD      NodeKind = iota // +
	ND_SUB                      // -
	ND_MUL                      // *
	ND_DIV                      // /
	ND_EQ                       // ==
	ND_NE                       // !=
	ND_LT                       // <
	ND_LE                       // <=
	ND_ASSIGN                   // =
	ND_LVAR                     // local variables
	ND_NUM                      // integer
	ND_RETURN                   // 'return'
	ND_IF                       // "if"
	ND_WHILE                    // "while"
	ND_FOR                      // "for"
	ND_BLOCK                    // {...}
	ND_FUNCCALL                 // function call
	ND_ADDR                     // unary &
	ND_DEREF                    // unary *
)

// define AST node
type Node struct {
	Kind NodeKind // the type of node
	Next *Node    // the next node
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

// the type of local variables
type LVar struct {
	Name   string // the name of the variable
	Offset int    // the offset from RBP
}

type VarList struct {
	Next *VarList
	Var  *LVar
}

// local variables
var locals *VarList

func newLVarNode(lvar *LVar, tok *Token) *Node {
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

func pushVar(name string) *LVar {
	lvar := &LVar{
		Name: name,
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

// params = ident ("," ident)*
func readFuncParams() *VarList {
	// printCurTok()
	// printCurFunc()
	if consume(")") != nil { // no argument
		return nil
	}

	head := &VarList{Var: pushVar(expectIdent())}
	cur := head

	for consume(")") == nil {
		expect(",")
		cur.Next = &VarList{Var: pushVar(expectIdent())}
		cur = cur.Next
	}

	return head
}

// function = ident "(" params? ")" "{" stmt* "}"
func function() *Function {
	// printCurTok()
	// printCurFunc()
	locals = nil

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

// stmt = expr ";"
//      | "{" stmt* "}"
//      | "if" "(" expr ")" stmt ("else" stmt)?
//      | "while" "(" expr ")" stmt
//      | "for" "(" expr? ";" expr? ";" expr? ")" stmt
//      | "return" expr ";"
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
			node.Init = expr()
			expect(";")
		}
		if consume(";") == nil {
			node.Cond = expr()
			expect(";")
		}
		if consume(")") == nil {
			node.Inc = expr()
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
		node = expr()
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

//unary = ("+" | "-" | "*" | "&")? unary
//      | primary
func unary() *Node {
	// printCurTok()
	// printCurFunc()
	if t := consume("+"); t != nil {
		return unary()
	}
	if t := consume("-"); t != nil {
		return newNode(ND_SUB, newNodeNum(0, t), unary(), t)
	}
	if t := consume("*"); t != nil {
		return newNode(ND_ADDR, unary(), nil, t)
	}
	if t := consume("&"); t != nil {
		return newNode(ND_DEREF, unary(), nil, t)
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
	for consume(",") != nil {
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
			node := &Node{Kind: ND_FUNCCALL}
			node.FuncName = tok.Str
			node.Args = funcArgs()
			return node
		}

		// local variables
		lvar := findLVar(tok)
		if lvar != nil {
			lvar = pushVar(tok.Str)
		}

		return newLVarNode(lvar, tok)
	}

	tok := token
	if tok.Kind != TK_NUM {
		errorTok(os.Stderr, tok, "it's not a number")
	}
	// otherwise, must be integer
	return newNodeNum(expectNumber(), tok)
}
