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

// // the types of AST node
// type NodeKind int

// const (
// 	ND_ADD NodeKind = iota // +
// 	ND_SUB                 // -
// 	ND_MUL                 // *
// 	ND_DIV                 // /
// 	ND_NUM                 // integer
// )

// // define AST node
// type Node struct {
// 	Kind NodeKind // type of node
// 	Lhs  *Node    // left
// 	Rhs  *Node    // right
// 	Val  int      // it would be used when kind is 'ND_NUM'
// }

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
		errorAt(curIdx, "is not integer")
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

	fmt.Print("#")

	for curIdx < len(userInput) {
		// skip space
		if userInput[curIdx] == ' ' {
			curIdx++
			continue
		}

		if userInput[curIdx] == '+' || userInput[curIdx] == '-' {
			cur = newToken(TK_RESERVED, cur, string(userInput[curIdx]))
			fmt.Printf(" '%s' ", cur.Str)
			curIdx++
			continue
		}

		if isDigit(userInput[curIdx]) {
			var sVal string = string(userInput[curIdx])
			curIdx++
			for ; curIdx < len(userInput) && isDigit(userInput[curIdx]); curIdx++ {
				sVal += string(userInput[curIdx])
			}
			cur = newToken(TK_NUM, cur, sVal)
			fmt.Printf(" '%s' ", cur.Str)
			v, err := strconv.Atoi(sVal)
			if err != nil {
				panic(err)
			}
			cur.Val = v
			curIdx++
			continue
		}

		errorAt(curIdx, "couldn't tokenize")
	}

	newToken(TK_EOF, cur, "")
	fmt.Println()
	return head.Next
}

func compile(arg string, w io.Writer) {
	// tokenize
	userInput = arg
	token = tokenize()

	// output the former of the assembly
	fmt.Fprintln(w, ".intel_syntax noprefix")
	fmt.Fprintln(w, ".global main")
	fmt.Fprintln(w, "main:")

	// check the beginning of expression is interger,
	// and output the first 'mov' command.
	fmt.Fprintf(w, "	mov rax, %d\n", expectNumber())

	// '+ <NUM>' or '- <NUM>'
	for !atEof() {
		if consume('+') {
			fmt.Fprintf(w, "	add rax, %d\n", expectNumber())
			continue
		}

		expect('-')
		fmt.Fprintf(w, "	sub rax, %d\n", expectNumber())
	}

	fmt.Fprintln(w, "	ret")
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("invalid number of arguments")
	}
	compile(os.Args[1], os.Stdout)
}
