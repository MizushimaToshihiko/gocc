package main

type NodeKind int

const (
	ND_ADD    NodeKind = iota // +
	ND_SUB                    // -
	ND_MUL                    // *
	ND_DIV                    // /
	ND_EQ                     // ==
	ND_NE                     // !=
	ND_LT                     // <
	ND_LE                     // <=
	ND_RETURN                 // return
	ND_NUM                    // integer
)

// define AST node
type Node struct {
	Kind NodeKind // type of node
	Next *Node    // Next node
	Lhs  *Node    // left branch
	Rhs  *Node    // right branch
	Val  int64    // it would be used when kind is 'ND_NUM'
}

func newBinary(kind NodeKind, lhs *Node, rhs *Node) *Node {
	return &Node{
		Kind: kind,
		Lhs:  lhs,
		Rhs:  rhs,
	}
}

func newNodeNum(val int64) *Node {
	return &Node{
		Kind: ND_NUM,
		Val:  val,
	}
}

func newUnary(kind NodeKind, expr *Node) *Node {
	node := &Node{Kind: kind, Lhs: expr}
	return node
}

// program = stmt*
func program() *Node {
	// printCurTok()
	// printCalledFunc()

	head := &Node{}
	cur := head

	for !atEof() {
		cur.Next = stmt()
		cur = cur.Next
	}
	return head.Next
}

// stmt = "return" expr (";" | "\n" | EOF)
//      | expr (";" | "\n" | EOF)
func stmt() *Node {
	// printCurTok()
	// printCalledFunc()

	if consume("return") != nil {
		node := newUnary(ND_RETURN, expr())
		expectEnd()
		return node
	}

	node := expr()
	expectEnd()
	return node
}

// expr       = equality
func expr() *Node {
	// printCurTok()
	// printCalledFunc()

	return equality()
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
		return newBinary(ND_SUB, newNodeNum(0), unary())
	}
	return primary()
}

// primary = "(" expr ")" | num
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

	// or must be integer
	return newNodeNum(expectNumber())
}
