package main

type NodeKind int

type Var struct {
	Next   *Var
	Name   string // Variable name
	Offset int    // Offset from RBP
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
	ND_RETURN                    // "return"
	ND_IF                        // "if"
	ND_WHILE                     // "for" instead of "while"
	ND_FOR                       // "for"
	ND_BLOCK                     // { ... }
	ND_FUNCALL                   // Function call
	ND_EXPR_STMT                 // Expression statement
	ND_VAR                       // Variables
	ND_NUM                       // Integer
)

// define AST node
type Node struct {
	Kind NodeKind // type of node
	Next *Node    // Next node

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

	Var *Var  // used if kind == ND_VAR
	Val int64 // it would be used when kind is 'ND_NUM'
}

var locals *Var

// Find a local variable by name.
func findVar(tok *Token) *Var {
	for v := locals; v != nil; v = v.Next {
		if v.Name == tok.Str {
			return v
		}
	}
	return nil
}

func newBinary(kind NodeKind, lhs *Node, rhs *Node) *Node {
	return &Node{
		Kind: kind,
		Lhs:  lhs,
		Rhs:  rhs,
	}
}

func newUnary(kind NodeKind, expr *Node) *Node {
	node := &Node{Kind: kind, Lhs: expr}
	return node
}

func newNum(val int64) *Node {
	return &Node{
		Kind: ND_NUM,
		Val:  val,
	}
}

func newVar(v *Var) *Node {
	return &Node{Kind: ND_VAR, Var: v}
}

func pushVar(name string) *Var {
	v := &Var{Next: locals, Name: name}
	locals = v
	return v
}

type Program struct {
	Node    *Node
	Locals  *Var
	StackSz int
}

// program = stmt*
func program() *Program {
	// printCurTok()
	// printCalledFunc()

	locals = nil

	head := &Node{}
	cur := head

	for !atEof() {
		cur.Next = stmt()
		cur = cur.Next
	}
	return &Program{Node: head.Next, Locals: locals}
}

func readExprStmt() *Node {
	return newUnary(ND_EXPR_STMT, expr())
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
//      | expr ";"
// for-stmt = "for" [ condition ] block .
// condition = expr .
// block = "{" stmt-list "};" .
// stmt-list = { stmt ";" } .
func stmt() *Node {
	// printCurTok()
	// printCalledFunc()

	if consume("return") != nil {
		node := newUnary(ND_RETURN, expr())
		expect(";")
		return node
	}

	if consume("if") != nil {
		node := &Node{Kind: ND_IF}
		node.Cond = expr()
		node.Then = stmt()
		if consume("else") != nil {
			node.Els = stmt()
		}
		return node
	}

	if consume("for") != nil {
		if !isForClause() {
			node := &Node{Kind: ND_WHILE}
			if peek("{") == nil {
				node.Cond = expr()
			} else {
				node.Cond = newNum(1)
			}

			node.Then = stmt()
			return node

		} else {
			node := &Node{Kind: ND_FOR}
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

	if consume("{") != nil {

		head := Node{}
		cur := &head

		for consume("}") == nil {
			cur.Next = stmt()
			cur = cur.Next
		}
		expect(";")
		return &Node{Kind: ND_BLOCK, Body: head.Next}
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
	if consume("=") != nil {
		node = newBinary(ND_ASSIGN, node, assign())
	}
	return node
}

// equality   = relational ("==" relational | "!=" relational)*
func equality() *Node {
	// printCurTok()
	// printCalledFunc()

	node := relational()

	for {
		if consume("==") != nil {
			node = newBinary(ND_EQ, node, relational())
		} else if consume("!=") != nil {
			node = newBinary(ND_NE, node, relational())
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
		if consume("<") != nil {
			node = newBinary(ND_LT, node, add())
		} else if consume("<=") != nil {
			node = newBinary(ND_LE, node, add())
		} else if consume(">") != nil {
			node = newBinary(ND_LT, add(), node)
		} else if consume(">=") != nil {
			node = newBinary(ND_LE, add(), node)
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
		if consume("+") != nil {
			node = newBinary(ND_ADD, node, mul())
		} else if consume("-") != nil {
			node = newBinary(ND_SUB, node, mul())
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
		if consume("*") != nil {
			node = newBinary(ND_MUL, node, unary())
		} else if consume("/") != nil {
			node = newBinary(ND_DIV, node, unary())
		} else {
			return node
		}
	}
}

// unary   = ("+" | "-")? unary
//         | primary
func unary() *Node {
	// printCurTok()
	// printCalledFunc()

	if consume("+") != nil {
		return unary()
	}
	if consume("-") != nil {
		return newBinary(ND_SUB, newNum(0), unary())
	}
	return primary()
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

	if tok := consumeIdent(); tok != nil {
		if consume("(") != nil {
			expect(")")
			return &Node{Kind: ND_FUNCALL, FuncName: tok.Str}
		}

		v := findVar(tok)
		if v == nil {
			v = pushVar(tok.Str)
		}
		return newVar(v)
	}

	// or must be integer
	return newNum(expectNumber())
}
