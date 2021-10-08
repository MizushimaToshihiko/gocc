//
// AST parser
//
package main

// the types of AST node
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
	ND_ASSIGN                 // =
	ND_LVAR                   // local variables
	ND_NUM                    // integer
	ND_RETURN                 // 'return'
	ND_IF                     // "if"
	ND_WHILE                  // "while"
	ND_FOR                    // "for"
	ND_BLOCK                  // {...}
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

// primary = num | ident | "(" expr ")"
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

// // walk AST in in-order
// func walkInOrder(node *Node) {
// 	fmt.Print("# Nodes in-order: ")
// 	inOrder(node)
// 	fmt.Println()
// }

// func inOrder(node *Node) {
// 	if node == nil {
// 		return
// 	}
// 	inOrder(node.Lhs)
// 	switch node.Kind {
// 	case ND_NUM:
// 		if isLeaf(node) {
// 			fmt.Printf(" '%s': %d: leaf ", "ND_NUM", node.Val)
// 		} else {
// 			fmt.Printf(" '%s': %d: ", "ND_NUM", node.Val)
// 		}
// 	case ND_ADD:
// 		if isLeaf(node) {
// 			fmt.Printf(" '%s': leaf ", "ND_ADD: +")
// 		} else {
// 			fmt.Printf(" '%s': ", "ND_ADD: +")
// 		}
// 	case ND_SUB:
// 		if isLeaf(node) {
// 			fmt.Printf(" '%s': leaf ", "ND_SUB: -")
// 		} else {
// 			fmt.Printf(" '%s': ", "ND_SUB: -")
// 		}
// 	case ND_MUL:
// 		if isLeaf(node) {
// 			fmt.Printf(" '%s': leaf ", "ND_MUL: *")
// 		} else {
// 			fmt.Printf(" '%s': ", "ND_MUL: *")
// 		}
// 	case ND_DIV:
// 		if isLeaf(node) {
// 			fmt.Printf(" '%s': leaf ", "ND_DIV: /")
// 		} else {
// 			fmt.Printf(" '%s': ", "ND_DIV: /")
// 		}
// 	case ND_EQ:
// 		if isLeaf(node) {
// 			fmt.Printf(" '%s': leaf ", "ND_EQ: ==")
// 		} else {
// 			fmt.Printf(" '%s': ", "ND_EQ: ==")
// 		}
// 	case ND_NE:
// 		if isLeaf(node) {
// 			fmt.Printf(" '%s': leaf ", "ND_NE: !=")
// 		} else {
// 			fmt.Printf(" '%s': ", "ND_NE: !=")
// 		}
// 	case ND_LT:
// 		if isLeaf(node) {
// 			fmt.Printf(" '%s': leaf ", "ND_LT: <")
// 		} else {
// 			fmt.Printf(" '%s': ", "ND_LT: <")
// 		}
// 	case ND_LE:
// 		if isLeaf(node) {
// 			fmt.Printf(" '%s': leaf ", "ND_LE: <=")
// 		} else {
// 			fmt.Printf(" '%s': ", "ND_LE: <=")
// 		}
// 	}
// 	inOrder(node.Rhs)
// }

// // walk AST in pre-order
// func walkPreOrder(node *Node) {
// 	fmt.Print("# Nodes pre-order: ")
// 	preOrder(node)
// 	fmt.Println()
// }

// func preOrder(node *Node) {
// 	if node == nil {
// 		return
// 	}
// 	switch node.Kind {
// 	case ND_NUM:
// 		if isLeaf(node) {
// 			fmt.Printf(" '%s': %d: leaf ", "ND_NUM", node.Val)
// 		} else {
// 			fmt.Printf(" '%s': %d: ", "ND_NUM", node.Val)
// 		}
// 	case ND_ADD:
// 		if isLeaf(node) {
// 			fmt.Printf(" '%s': leaf ", "ND_ADD: +")
// 		} else {
// 			fmt.Printf(" '%s': ", "ND_ADD: +")
// 		}
// 	case ND_SUB:
// 		if isLeaf(node) {
// 			fmt.Printf(" '%s': leaf ", "ND_SUB: -")
// 		} else {
// 			fmt.Printf(" '%s': ", "ND_SUB: -")
// 		}
// 	case ND_MUL:
// 		if isLeaf(node) {
// 			fmt.Printf(" '%s': leaf ", "ND_MUL: *")
// 		} else {
// 			fmt.Printf(" '%s': ", "ND_MUL: *")
// 		}
// 	case ND_DIV:
// 		if isLeaf(node) {
// 			fmt.Printf(" '%s': leaf ", "ND_DIV: /")
// 		} else {
// 			fmt.Printf(" '%s': ", "ND_DIV: /")
// 		}
// 	case ND_EQ:
// 		if isLeaf(node) {
// 			fmt.Printf(" '%s': leaf ", "ND_EQ: ==")
// 		} else {
// 			fmt.Printf(" '%s': ", "ND_EQ: ==")
// 		}
// 	case ND_NE:
// 		if isLeaf(node) {
// 			fmt.Printf(" '%s': leaf ", "ND_NE: !=")
// 		} else {
// 			fmt.Printf(" '%s': ", "ND_NE: !=")
// 		}
// 	case ND_LT:
// 		if isLeaf(node) {
// 			fmt.Printf(" '%s': leaf ", "ND_LT: <")
// 		} else {
// 			fmt.Printf(" '%s': ", "ND_LT: <")
// 		}
// 	case ND_LE:
// 		if isLeaf(node) {
// 			fmt.Printf(" '%s': leaf ", "ND_LE: <=")
// 		} else {
// 			fmt.Printf(" '%s': ", "ND_LE: <=")
// 		}
// 	}
// 	preOrder(node.Lhs)
// 	preOrder(node.Rhs)
// }

// func isLeaf(node *Node) bool {
// 	return node.Lhs == nil && node.Rhs == nil
// }
