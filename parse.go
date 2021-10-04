//
// parser
//
package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// set TokenKind with enum
type TokenKind int

const (
	TK_RESERVED TokenKind = iota // symbol
	TK_NUM                       // integer
	TK_EOF                       // the end of tokens
)

type Token struct {
	Kind TokenKind // type of token
	Next *Token    // next
	Val  int       // if 'kind' is TK_NUM, it's integer
	Str  string    // token string
	Len  int       // length of token
}

// current token
var token *Token

// inputted program
var userInput string
var curIdx int

// for error report
// it's arguments are same as printf
func errorAt(errIdx int, formt string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s\n", userInput)
	fmt.Fprintf(os.Stderr, "%*s", errIdx, " ")
	fmt.Fprint(os.Stderr, "^ ")
	fmt.Fprintf(os.Stderr, formt, a...)
	fmt.Fprint(os.Stderr, "\n")
	os.Exit(1)
}

// if the next token is expected symbol, the read position
// of token exceed one character, and returns true.
func consume(op string) bool {
	if token.Kind != TK_RESERVED ||
		len(op) != token.Len ||
		token.Str != op {
		return false
	}
	token = token.Next
	return true
}

// if the next token is expected symbol, the read position
// of token exceed one character
func expect(op string) {
	if token.Kind != TK_RESERVED ||
		len(op) != token.Len ||
		token.Str != op {
		errorAt(curIdx, "is not '%s'", string(op))
	}
	token = token.Next
}

// if next token is integer, the read position of token exceed one
// character or report an error.
func expectNumber() int {
	if token.Kind != TK_NUM {
		errorAt(curIdx, "is not an integer")
	}
	val := token.Val
	token = token.Next
	return val
}

func atEof() bool {
	return token.Kind == TK_EOF
}

// make new token and append to the end of cur.
func newToken(kind TokenKind, cur *Token, str string, len int) *Token {
	tok := &Token{Kind: kind, Str: str, Len: len}
	cur.Next = tok
	return tok
}

func startsWith(pp, qq string) bool {
	p := []byte(pp)
	q := []byte(qq)
	return reflect.DeepEqual(p[:len(q)], q)
}

func isDigit(op byte) bool {
	return '0' <= op && op <= '9'
}

// tokenize inputted string 'p', and return this.
func tokenize() *Token {
	var head Token
	head.Next = nil
	cur := &head

	// for printToken
	headTok = &head

	for curIdx < len(userInput) {
		// skip space(s)
		if userInput[curIdx] == ' ' {
			curIdx++
			continue
		}

		// multi-letter punctuator
		if startsWith(userInput[curIdx:], "==") ||
			startsWith(userInput[curIdx:], "!=") ||
			startsWith(userInput[curIdx:], "<=") ||
			startsWith(userInput[curIdx:], ">=") {
			cur = newToken(TK_RESERVED, cur, userInput[curIdx:curIdx+2], 2)
			curIdx += 2
			continue
		}

		if strings.Contains("+-()*/<>=", string(userInput[curIdx])) {
			cur = newToken(TK_RESERVED, cur, string(userInput[curIdx]), 1)
			curIdx++
			continue
		}

		if isDigit(userInput[curIdx]) {
			var sVal string
			for ; curIdx < len(userInput) && isDigit(userInput[curIdx]); curIdx++ {
				sVal += string(userInput[curIdx])
			}
			cur = newToken(TK_NUM, cur, sVal, len(sVal))
			v, err := strconv.Atoi(sVal)
			if err != nil {
				panic(err)
			}
			cur.Val = v
			continue
		}

		errorAt(curIdx, "couldn't tokenize")
	}

	newToken(TK_EOF, cur, "", 0)
	return head.Next
}

// for printTokens function, the pointer of the head token
// stored in 'headTok'.
var headTok *Token

//
func printTokens() {
	fmt.Print("# Tokens: ")
	tok := headTok.Next
	for tok.Next != nil {
		fmt.Printf(" '%s' ", tok.Str)
		tok = tok.Next
	}
	fmt.Println()
}

// the types of AST node
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
	Val  int      // it would be used when kind is 'ND_NUM'
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

// expr       = equality
func expr() *Node {
	return equality()
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

// add        = mul ("+" mul | "-" mul)*
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

// primary = "(" expr ")" | num
func primary() *Node {
	// if the next token is '(', the program must be
	// "(" expr ")"
	if consume("(") {
		node := expr()
		expect(")")
		return node
	}

	// or must be integer
	return newNodeNum(expectNumber())
}
