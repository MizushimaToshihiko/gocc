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
	TK_RESERVED TokenKind = iota // Reserved words, and puncturators
	TK_IDENT                     // idenfier such as variables, function names
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

// func min(x, y int) int {
// 	if x < y {
// 		return x
// 	}
// 	return y
// }

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

func errorTok(w io.Writer, tok *Token, formt string, a ...interface{}) {
	if tok != nil {
		errorAt(w, curIdx, formt, a...)
	}
	fmt.Fprintf(w, formt, a...)
	fmt.Fprintf(w, "\n")
	os.Exit(1)
}

// peek function returns the token, when the current token matches 's'.
func peek(s string) *Token {
	if token.Kind != TK_RESERVED ||
		len(s) != token.Len ||
		token.Str != s {
		return nil
	}
	return token
}

// if the next token is expected symbol, the read position
// of token exceed one character, and returns token(pointer).
func consume(s string) *Token {
	// defer printCurTok()
	if peek(s) == nil {
		return nil
	}
	t := token
	token = token.Next
	return t
}

// consume the current token if it is an identifier
func consumeIdent() *Token {
	// defer printCurTok()
	if token.Kind != TK_IDENT {
		return nil
	}
	t := token
	token = token.Next
	return t
}

// if the next token is an expected symbol, the read position
// of token exceed one token.
func expect(s string) {
	// defer printCurTok()
	if peek(s) == nil {
		errorAt(os.Stderr, curIdx, "is not '%s'", string(s))
	}
	token = token.Next
}

// if next token is integer, the read position of token exceed one
// character or report an error.
func expectNumber() int {
	// defer printCurTok()
	if token.Kind != TK_NUM {
		errorAt(os.Stderr, curIdx, "is not a number")
	}
	val := token.Val
	token = token.Next
	return val
}

func expectIdent() string {
	// defer printCurTok()
	if token.Kind != TK_IDENT {
		fmt.Println(userInput)
		fmt.Printf("token %d:'%s'\n", token.Kind, token.Str)
		panic("expect an identifier")
	}
	s := token.Str
	token = token.Next
	return s
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

// startsWith compare 'pp' and 'qq' , qq is keyword
func startsWith(pp, qq string) bool {
	p, q := []byte(pp), []byte(qq)
	return len(p) >= len(q) && reflect.DeepEqual(p[:len(q)], q)
}

func startsWithReserved(p string) string {
	// reserved words
	kw := []string{"return", "if", "then", "else", "while", "for",
		"int"}

	for _, k := range kw {
		if startsWith(p, k) && len(p) > len(k) && !isAlNum(p[len(k)]) {
			return k
		}
	}

	// Multi-letter punctuator
	ops := []string{"==", "!=", "<=", ">="}

	for _, op := range ops {
		if startsWith(p, op) {
			return op
		}
	}
	return ""
}

func isDigit(op byte) bool {
	return '0' <= op && op <= '9'
}

func isAlpha(c byte) bool {
	return ('a' <= c && c <= 'z') ||
		('A' <= c && c <= 'Z') ||
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

	// for printToken
	headTok = &head

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

		// single-letter punctuator
		if strings.Contains("+-()*/<>=;{},&", string(userInput[curIdx])) {
			cur = newToken(TK_RESERVED, cur, string(userInput[curIdx]), 1)
			curIdx++
			continue
		}

		// identifier
		// if 'userInput[cutIdx]' is alphabets, it makes a token of TK_IDENT type.
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

		errorAt(os.Stderr, curIdx, "invalid token")
	}

	newToken(TK_EOF, cur, "", 0)
	return head.Next
}
