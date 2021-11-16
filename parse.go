package main

type NodeKind int

type Var struct {
	Name   string // Variable name
	Ty     *Type  // Type
	Offset int    // Offset from RBP
}

type VarList struct {
	Next *VarList
	Var  *Var
}

const (
	ND_ADD       NodeKind = iota // +
	ND_SUB                       // -
	ND_MUL                       // *
	ND_DIV                       // /
	ND_EQ                        // ==
	ND_NE                        // !=
	ND_LT                        // <
	ND_LE                        // <=
	ND_ASSIGN                    // = , ":=" is unimplememted
	ND_ADDR                      // unary &
	ND_DEREF                     // urary *
	ND_RETURN                    // "return"
	ND_IF                        // "if"
	ND_WHILE                     // "for" instead of "while"
	ND_FOR                       // "for"
	ND_BLOCK                     // { ... }
	ND_FUNCALL                   // Function call
	ND_EXPR_STMT                 // Expression statement
	ND_VAR                       // Variables
	ND_NUM                       // Integer
	ND_NULL                      // Empty statement
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

	// Function call
	FuncName string
	Args     *Node

	Var *Var  // used if kind == ND_VAR
	Val int64 // it would be used when kind is 'ND_NUM'
}

var locals *VarList

// Find a local variable by name.
func findVar(tok *Token) *Var {
	for vl := locals; vl != nil; vl = vl.Next {
		if vl.Var.Name == tok.Str {
			return vl.Var
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

func pushVar(name string, ty *Type) *Var {
	v := &Var{Name: name, Ty: ty}

	vl := &VarList{Var: v, Next: locals}
	locals = vl
	return vl.Var
}

type Function struct {
	Next   *Function
	Name   string
	Params *VarList

	Node    *Node
	Locals  *VarList
	StackSz int
}

// program = function*
func program() *Function {
	// printCurTok()
	// printCalledFunc()

	head := &Function{}
	cur := head

	for !atEof() {
		cur.Next = function()
		cur = cur.Next
	}
	return head.Next
}

// basetype = "*"* "int"
func basetype() *Type {
	ty := intType()
	for consume("*") != nil {
		ty = pointerTo(ty)
	}
	expect("int")
	return ty
}

// param = ident basetype
func readFuncParam() *VarList {
	vl := &VarList{}
	name := expectIdent()
	ty := basetype()
	vl.Var = pushVar(name, ty)
	return vl
}

// params = ident ("," ident)*
func readFuncParams() *VarList {
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
	locals = nil

	expect("func")
	fn := &Function{Name: expectIdent()}
	expect("(")
	fn.Params = readFuncParams()
	basetype()
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

// declaration = "var" ident basetype ("=" expr)
func declaration() *Node {
	tok := token
	expect("var")
	name := expectIdent()
	ty := basetype()
	v := pushVar(name, ty)

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
	t := token
	return newUnary(ND_EXPR_STMT, expr(), t)
}

func isForClause() bool {
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

		for consume("}") == nil {
			cur.Next = stmt()
			cur = cur.Next
		}
		expect(";")
		return &Node{Kind: ND_BLOCK, Body: head.Next, Tok: t}
	}

	if t := peek("var"); t != nil {
		return declaration()
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
//         | primary
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
	return primary()
}

// func-args = "(" (assign ("," assign)*)? ")"
func funcArgs() *Node {
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

		v := findVar(t)
		if v == nil {
			panic("\n" + errorTok(t, "undifined variable"))
		}
		return newVar(v, t)
	}

	// or must be integer
	t := token
	if t.Kind != TK_NUM {
		panic("\n" + errorTok(t, "expected expression"))
	}
	return newNum(expectNumber(), t)
}
