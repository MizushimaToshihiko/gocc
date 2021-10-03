package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

// set TokenKind with enum
type TokenKind int

const (
	TK_RESERVED TokenKind = iota
	TK_NUM
	TK_EOF
)

type Token struct {
	Kind TokenKind // type of token
	Next *Token    // next
	Val  int       // if 'kind' is TK_NUM, it's integer
	Str  string    // token string
}

// current token
var token *Token

// inputted program
var userInput string
var curIdx int

// for error report
// it's arguments are same as printf
func errorAt(errIdx int, formt string, a ...interface{}) {
	if _, err := fmt.Fprintf(os.Stderr, "%s\n", userInput); err != nil {
		panic(err)
	}
	if _, err := fmt.Fprintf(os.Stderr, "%*s", errIdx, " "); err != nil { // output space
		panic(err)
	}

	if _, err := fmt.Fprint(os.Stderr, "^ "); err != nil {
		panic(err)
	}

	if _, err := fmt.Fprintf(os.Stderr, formt, a...); err != nil {
		panic(err)
	}

	if _, err := fmt.Fprint(os.Stderr, "\n"); err != nil {
		panic(err)
	}
	os.Exit(1)
}

// if the next token is expected symbol, the read position
// of token exceed one character, and returns true.
func consume(op byte) bool {
	if token.Kind != TK_RESERVED || token.Str[0] != op {
		return false
	}
	token = token.Next
	return true
}

// if the next token is expected symbol, the read position
// of token exceed one character
func expect(op byte) {
	if token.Kind != TK_RESERVED || token.Str[0] != op {
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
func newToken(kind TokenKind, cur *Token, str string) *Token {
	tok := &Token{Kind: kind, Str: str}
	cur.Next = tok
	return tok
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

		if userInput[curIdx] == '+' ||
			userInput[curIdx] == '-' ||
			userInput[curIdx] == '*' ||
			userInput[curIdx] == '/' ||
			userInput[curIdx] == '(' ||
			userInput[curIdx] == ')' {
			cur = newToken(TK_RESERVED, cur, string(userInput[curIdx]))
			curIdx++
			continue
		}

		if isDigit(userInput[curIdx]) {
			var sVal string
			for ; curIdx < len(userInput) && isDigit(userInput[curIdx]); curIdx++ {
				sVal += string(userInput[curIdx])
			}
			cur = newToken(TK_NUM, cur, sVal)
			v, err := strconv.Atoi(sVal)
			if err != nil {
				panic(err)
			}
			cur.Val = v
			continue
		}

		errorAt(curIdx, "couldn't tokenize")
	}

	newToken(TK_EOF, cur, "")
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

// expr = mul ("+" mul | "-" mul)*
func expr() *Node {
	node := mul()

	for {
		if consume('+') {
			node = newNode(ND_ADD, node, mul())
		} else if consume('-') {
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
		if consume('*') {
			node = newNode(ND_MUL, node, unary())
		} else if consume('/') {
			node = newNode(ND_DIV, node, unary())
		} else {
			return node
		}
	}
}

// unary   = ("+" | "-")? primary
func unary() *Node {
	if consume('+') {
		return primary()
	}
	if consume('-') {
		return newNode(ND_SUB, newNodeNum(0), primary())
	}
	return primary()
}

// primary = "(" expr ")" | num
func primary() *Node {
	// if the next token is '(', the program must be
	// "(" expr ")"
	if consume('(') {
		node := expr()
		expect(')')
		return node
	}

	// or must be integer
	return newNodeNum(expectNumber())
}

func gen(node *Node, w io.Writer) (err error) {
	if node.Kind == ND_NUM {
		_, err = fmt.Fprintf(w, "	push %d\n", node.Val)
		return
	}

	err = gen(node.Lhs, w)
	if err != nil {
		return
	}
	err = gen(node.Rhs, w)
	if err != nil {
		return
	}

	_, err = fmt.Fprintln(w, "	pop rdi")
	if err != nil {
		return
	}
	_, err = fmt.Fprintln(w, "	pop rax")
	if err != nil {
		return
	}

	switch node.Kind {
	case ND_ADD:
		if _, err = fmt.Fprintln(w, "	add rax, rdi"); err != nil {
			return
		}
	case ND_SUB:
		if _, err = fmt.Fprintln(w, "sub rax, rdi"); err != nil {
			return
		}
	case ND_MUL:
		if _, err = fmt.Fprintln(w, "imul rax, rdi"); err != nil {
			return
		}
	case ND_DIV:
		if _, err = fmt.Fprintln(w, "cqo"); err != nil {
			return
		}
		if _, err = fmt.Fprintln(w, "idiv rdi"); err != nil {
			return
		}
	}

	if _, err = fmt.Fprintln(w, "push rax"); err != nil {
		return
	}

	return
}

func compile(arg string, w io.Writer) error {
	// tokenize and parse
	curIdx = 0 // for test
	userInput = arg
	token = tokenize()
	node := expr()

	printTokens()
	// output the former 3 lines of the assembly
	if _, err := fmt.Fprintln(w, ".intel_syntax noprefix\n.globl main\nmain:"); err != nil {
		return err
	}

	// make the asm code, down on the AST
	if err := gen(node, w); err != nil {
		return err
	}

	// the value of the expression should remain on the top of 'stack',
	// so load this value into rax.
	if _, err := fmt.Fprintln(w, "	pop rax"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "	ret"); err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("invalid number of arguments")
	}

	if err := compile(os.Args[1], os.Stdout); err != nil {
		log.Fatal(err)
	}
}
