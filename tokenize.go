//
// tokenizier
//
package main

import (
	"fmt"
	"io"
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

// current index in 'userInput'
var curIdx int

// for error report
// it's arguments are same as printf
func errorAt(w io.Writer, errIdx int, formt string, a ...interface{}) {
	fmt.Fprintf(w, "%s\n", userInput)
	fmt.Fprintf(w, "%*s", errIdx, " ")
	fmt.Fprint(w, "^ ")
	fmt.Fprintf(w, formt, a...)
	fmt.Fprint(w, "\n")
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
		errorAt(os.Stderr, curIdx, "is not '%s'", string(op))
	}
	token = token.Next
}

// if next token is integer, the read position of token exceed one
// character or report an error.
func expectNumber() int {
	if token.Kind != TK_NUM {
		errorAt(os.Stderr, curIdx, "is not a number")
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

	// // for printToken
	// headTok = &head

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

		// single-letter punctuator
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

		errorAt(os.Stderr, curIdx, "couldn't tokenize")
	}

	newToken(TK_EOF, cur, "", 0)
	return head.Next
}

// // for printTokens function, the pointer of the head token
// // stored in 'headTok'.
// var headTok *Token

// //
// func printTokens() {
// 	fmt.Print("# Tokens: ")
// 	tok := headTok.Next
// 	for tok.Next != nil {
// 		fmt.Printf(" '%s' ", tok.Str)
// 		tok = tok.Next
// 	}

// 	if tok.Kind == TK_EOF {
// 		fmt.Print(" 'EOF' ")
// 	}

// 	fmt.Println()
// }
