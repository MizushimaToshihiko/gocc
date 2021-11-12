package main

type NodeKind int

const (
	ND_ADD NodeKind = iota // +
	ND_SUB                 // -
	ND_MUL                 // *
	ND_DIV                 // /
	ND_EQ                  // ==
	ND_NE                  // !=
	ND_LT                  // <
	ND_LE                  // <=
	ND_NUM                 // integer
)

// define AST node
type Node struct {
	Kind NodeKind // type of node
	Lhs  *Node    // left branch
	Rhs  *Node    // right branch
	Val  int64    // it would be used when kind is 'ND_NUM'
}

func newNode(kind NodeKind, lhs *Node, rhs *Node) *Node {
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

// expr       = equality
func expr() *Node {
	return equality()
}

// equality   = relational ("==" relational | "!=" relational)*
func equality() *Node {
	node := relational()

	for {
		if consume("==") != nil {
			node = newNode(ND_EQ, node, relational())
		} else if consume("!=") != nil {
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
		if consume("<") != nil {
			node = newNode(ND_LT, node, add())
		} else if consume("<=") != nil {
			node = newNode(ND_LE, node, add())
		} else if consume(">") != nil {
			node = newNode(ND_LT, add(), node)
		} else if consume(">=") != nil {
			node = newNode(ND_LE, add(), node)
		} else {
			return node
		}
	}
}

// add        = mul ("+" mul | "-" mul)*
func add() *Node {
	node := mul()

	for {
		if consume("+") != nil {
			node = newNode(ND_ADD, node, mul())
		} else if consume("-") != nil {
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
		if consume("*") != nil {
			node = newNode(ND_MUL, node, unary())
		} else if consume("/") != nil {
			node = newNode(ND_DIV, node, unary())
		} else {
			return node
		}
	}
}

// unary   = ("+" | "-")? unary
//         | primary
func unary() *Node {
	if consume("+") != nil {
		return unary()
	}
	if consume("-") != nil {
		return newNode(ND_SUB, newNodeNum(0), unary())
	}
	return primary()
}

// primary = "(" expr ")" | num
func primary() *Node {
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
