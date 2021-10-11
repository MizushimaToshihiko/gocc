//
// AST parser
//
package main

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
)

// define AST node
type Node struct {
	Kind NodeKind // the type of node
	Next *Node    // the next node

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

	Val    int // it would be used when 'Kind' is 'ND_NUM'
	Offset int // it would be used when 'Kind' is 'ND_LVAR'

}

func newNode(kind NodeKind, lhs *Node, rhs *Node) *Node {
	return &Node{
		Kind: kind,
		Lhs:  lhs,
		Rhs:  rhs,
	}
}

func newNodeNum(val int) *Node {
	return &Node{
		Kind: ND_NUM,
		Val:  val,
	}
}

// code is a slice to store prased nodes.
var code [100]*Node

// program = stmt*
func program() {
	i := 0
	for !atEof() {
		code[i] = stmt()
		i++
	}
	code[i] = nil
}

// stmt = expr ";"
//      | "{" stmt* "}"
//      | "if" "(" expr ")" stmt ("else" stmt)?
//      | "while" "(" expr ")" stmt
//      | "for" "(" expr? ";" expr? ";" expr? ")" stmt
//      | "return" expr ";"
func stmt() *Node {
	var node *Node

	if consume("return") {

		node = &Node{Kind: ND_RETURN, Lhs: expr()}
		expect(";")

	} else if consume("for") {

		expect("(")
		node = &Node{Kind: ND_FOR}

		if !consume(";") {
			node.Init = expr()
			expect(";")
		}
		if !consume(";") {
			node.Cond = expr()
			expect(";")
		}
		if !consume(")") {
			node.Inc = expr()
			expect(")")
		}
		node.Then = stmt()

	} else if consume("while") {

		expect("(")
		node = &Node{Kind: ND_WHILE}
		node.Cond = expr()
		expect(")")
		node.Then = stmt()

	} else if consume("if") {

		expect("(")
		node = &Node{Kind: ND_IF}
		node.Cond = expr()
		expect(")")
		node.Then = stmt()

		if consume("else") {
			node.Els = stmt()
		}

	} else if consume("{") {

		head := Node{}
		cur := &head

		for !consume("}") {
			cur.Next = stmt()
			cur = cur.Next
		}

		node = &Node{Kind: ND_BLOCK}
		node.Body = head.Next

	} else {
		node = expr()
		expect(";")
	}

	return node
}

// expr       = equality
func expr() *Node {
	return assign()
}

// assign     = equality ("=" assign)?
func assign() *Node {
	node := equality()
	if consume("=") {
		node = newNode(ND_ASSIGN, node, assign())
	}
	return node
}

// equality   = relational ("==" relational | "!=" relational)*
func equality() *Node {
	node := relational()

	for {
		if consume("==") {
			node = newNode(ND_EQ, node, relational())
		} else if consume("!=") {
			node = newNode(ND_NE, node, relational())
		} else {
			return node
		}
	}
}

// relational = add ("<" add | "<=" add | ">" add | ">=" add)*
func relational() *Node {
	node := add()

	for {
		if consume("<") {
			node = newNode(ND_LT, node, add())
		} else if consume("<=") {
			node = newNode(ND_LE, node, add())
		} else if consume(">") {
			node = newNode(ND_LT, add(), node)
		} else if consume(">=") {
			node = newNode(ND_LE, add(), node)
		} else {
			return node
		}
	}
}

// add = mul ("+" mul | "-" mul)*
func add() *Node {
	node := mul()

	for {
		if consume("+") {
			node = newNode(ND_ADD, node, mul())
		} else if consume("-") {
			node = newNode(ND_SUB, node, mul())
		} else {
			return node
		}
	}
}

// mul = unary ("*" unary | "/" unary)*
func mul() *Node {
	node := unary()

	for {
		if consume("*") {
			node = newNode(ND_MUL, node, unary())
		} else if consume("/") {
			node = newNode(ND_DIV, node, unary())
		} else {
			return node
		}
	}
}

// unary   = ("+" | "-")? unary
//         | primary
func unary() *Node {
	if consume("+") {
		return unary()
	}
	if consume("-") {
		return newNode(ND_SUB, newNodeNum(0), unary())
	}
	return primary()
}

// func-args = "(" (assign("," assign)*)? ")"
func funcArgs() *Node {
	if consume(")") {
		return nil
	}

	head := assign()
	cur := head
	for consume(",") {
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
	// if the next token is '(', the program must be
	// "(" expr ")"
	if consume("(") {
		node := expr()
		expect(")")
		return node
	}

	tok := consumeIdent()
	if tok != nil {
		if consume("(") { // function call
			node := &Node{Kind: ND_FUNCCALL}
			node.FuncName = tok.Str
			node.Args = funcArgs()
			return node
		}

		// local variables
		node := &Node{Kind: ND_LVAR}

		lvar := findLVar(tok)
		if lvar != nil {
			node.Offset = lvar.Offset
		} else {
			lvar = &LVar{
				Next:   locals,
				Name:   tok.Str,
				Len:    tok.Len,
				Offset: 8,
			}
			if locals != nil {
				lvar.Offset += locals.Offset
			}
			node.Offset = lvar.Offset
			locals = lvar
		}
		return node

	}

	// otherwise, must be integer
	return newNodeNum(expectNumber())
}
