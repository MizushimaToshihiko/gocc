//
// AST parser
//
package main

import "os"

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

	Val    int // it would be used when 'Kind' is 'ND_NUM'
	Offset int // it would be used when 'Kind' is 'ND_LVAR'

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

	if consume("return") != nil {

		node = &Node{Kind: ND_RETURN, Lhs: expr()}
		expect(";")

	} else if consume("for") != nil {

		expect("(")
		node = &Node{Kind: ND_FOR}

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

	} else if consume("while") != nil {

		expect("(")
		node = &Node{Kind: ND_WHILE}
		node.Cond = expr()
		expect(")")
		node.Then = stmt()

	} else if consume("if") != nil {

		expect("(")
		node = &Node{Kind: ND_IF}
		node.Cond = expr()
		expect(")")
		node.Then = stmt()

		if consume("else") != nil {
			node.Els = stmt()
		}

	} else if consume("{") != nil {

		head := Node{}
		cur := &head

		for consume("}") == nil {
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
	if t := consume("="); t != nil {
		node = newNode(ND_ASSIGN, node, assign(), t)
	}
	return node
}

// equality   = relational ("==" relational | "!=" relational)*
func equality() *Node {
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
	// if the next token is '(', the program must be
	// "(" expr ")"
	if t := consume("("); t != nil {
		node := expr()
		expect(")")
		return node
	}

	tok := consumeIdent()
	if tok != nil {
		if t := consume("("); t != nil { // function call
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
	tok = token
	if tok.Kind != TK_NUM {
		errorTok(os.Stderr, tok, "expected number")
	}
	return newNodeNum(expectNumber(), tok)
}
