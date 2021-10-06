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
	TK_IDENT                     // idenfier such as variables, function names
	TK_NUM                       // integer
	TK_RETURN                    // 'return' statement
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

// the type of local variables
type LVar struct {
	Next   *LVar
	Name   string
	Len    int
	Offset int
}

// local variables
var locals *LVar

// search a local variable by name.
// if it wasn't find, return nil.
func findLVar(tok *Token) *LVar {
	for lvar := locals; lvar != nil; lvar = lvar.Next {
		if lvar.Len == tok.Len && startsWith(tok.Str, lvar.Name) {
			return lvar
		}
	}
	return nil
}

// inputted program
var userInput string

// current index in 'userInput'
var curIdx int

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

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
	if (token.Kind != TK_RESERVED &&
		token.Kind != TK_RETURN) ||
		len(op) != token.Len ||
		token.Str != op {
		return false
	}
	token = token.Next
	return true
}

// consume the current token if it is an identifier
func consumeIdent() *Token {
	if token.Kind != TK_IDENT {
		return nil
	}
	t := token
	token = token.Next
	return t
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

// startsWith compare 'pp' and 'qq' , pp is keyword
func startsWith(pp, qq string) bool {
	p, q := []byte(pp), []byte(qq)
	return reflect.DeepEqual(p[:len(q)], q)
}

func startsWithReserved(p string) string {
	// keyword
	kw := []string{"return", "if", "then"}

	for _, k := range kw {
		if startsWith(k, p) && !isAlNum(p[min(len(k), len(p))]) {
			return k
		}
	}

	// Multi-letter punctuator
	ops := []string{"==", "!=", "<=", ">="}

	for _, op := range ops {
		if startsWith(op, p) {
			return op
		}
	}
	return ""
}

func isDigit(op byte) bool {
	return '0' <= op && op <= '9'
}

func isAlpha(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') ||
		(c == '_')
}

func isAlNum(c byte) bool {
	return isAlpha(c) || ('0' <= c && c <= '9')
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

		// keyword or multi-letter punctuator
		kw := startsWithReserved(userInput[curIdx:])
		if kw != "" {
			cur = newToken(TK_RESERVED, cur, kw, len(kw))
			curIdx += len(kw)
			continue
		}
		// if curIdx+2 <= len(userInput) &&
		// 	(startsWith(userInput[curIdx:curIdx+2], "==") ||
		// 		startsWith(userInput[curIdx:curIdx+2], "!=") ||
		// 		startsWith(userInput[curIdx:curIdx+2], "<=") ||
		// 		startsWith(userInput[curIdx:curIdx+2], ">=")) {
		// 	cur = newToken(TK_RESERV,ED, cur, userInput[curIdx:curIdx+2], 2)
		// 	curIdx += 2
		// 	continue
		// }

		// single-letter punctuator
		if strings.Contains("+-()*/<>=;", string(userInput[curIdx])) {
			cur = newToken(TK_RESERVED, cur, string(userInput[curIdx]), 1)
			curIdx++
			continue
		}

		// // reserved words
		// if curIdx+6 <= len(userInput) &&
		// 	startsWith(userInput[curIdx:curIdx+6], "return") &&
		// 	!isAlNum(userInput[curIdx+6]) {
		// 	cur = newToken(TK_RETURN, cur, userInput[curIdx:curIdx+6], 6)
		// 	curIdx += 6
		// 	continue
		// }

		// identifier
		if isAlpha(userInput[curIdx]) {
			ident := make([]byte, 0, 20)
			for ; curIdx < len(userInput) && isAlNum(userInput[curIdx]); curIdx++ {
				ident = append(ident, userInput[curIdx])
			}
			cur = newToken(TK_IDENT, cur, string(ident), len(string(ident)))
			continue
		}

		// number
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
// 	var kind string
// 	for tok.Next != nil {
// 		switch tok.Kind {
// 		case TK_IDENT:
// 			kind = "IDENT"
// 		case TK_NUM:
// 			kind = "NUM"
// 		case TK_RESERVED:
// 			kind = "RESERVED"
// 		case TK_RETURN:
// 			kind = "RETURN"
// 		default:
// 			log.Fatal("unknown token kind")
// 		}
// 		fmt.Printf(" %s:'%s' ", kind, tok.Str)
// 		tok = tok.Next
// 	}

// 	if tok.Kind == TK_EOF {
// 		fmt.Print(" EOF ")
// 	}

// 	fmt.Println()
// }
