//
// tokenizier
//
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// set TokenKind with enum
type TokenKind int

const (
	TK_RESERVED TokenKind = iota // 0: Reserved words(key words), operators, and puncturators
	TK_IDENT                     // 1: idenfier such as variables, function names

	// literals
	TK_STR // 2: string literals
	TK_NUM // 3: integer

	TK_EOF // 4: the end of tokens
)

type Token struct {
	Kind TokenKind // type of token
	Next *Token    // next
	Val  int64     // if 'kind' is TK_NUM, it's integer
	Loc  int       // the location in 'userInput'
	Str  string    // token string
	Len  int       // length of token

	Contents []rune // string literal contents including terminating '\0'
	ContLen  int    // string literal length

	LineNo int // Line number
}

var curFilename string

// current token
var token *Token

// inputted program
var userInput []rune

// current index in 'userInput'
var curIdx int

// for error report
//
// foo.go:
// 10:x = y + 1;
//        ^ <error message here>
func verrorAt(lineNum, errIdx int, formt string, a ...interface{}) string {
	// get the start and end of the line 'errIdx' exists
	line := errIdx
	for 0 < line && userInput[line-1] != rune('\n') {
		line--
	}

	end := errIdx
	for end < len(userInput) && userInput[end] != rune('\n') {
		end++
	}

	// Show found lines along with file name and line number.
	res := fmt.Sprintf("\n%s:%d: ", curFilename, lineNum)
	indent := len(res)
	res += fmt.Sprintf("%.*s\n", end-line, string(userInput[line:end]))

	// Point the error location with "^" and display the error message.
	pos := errIdx - line + indent

	return res + fmt.Sprintf("%*s", pos, " ") +
		"^ " +
		fmt.Sprintf(formt, a...) +
		"\n\n"
}

// it's arguments are same as printf
func errorAt(errIdx int, formt string, a ...interface{}) string {

	// Find out what line the found line is in the whole.
	lineNum := 1
	for i := 0; i < len(userInput); i++ {
		if userInput[i] == rune('\n') {
			lineNum++
		}
	}
	return verrorAt(lineNum, errIdx, formt, a...)
}

func errorTok(tok *Token, formt string, ap ...interface{}) string {
	var errStr string
	if tok != nil {
		errStr += verrorAt(tok.LineNo, tok.Loc, formt, ap...)
	}

	return errStr
}

// strNdUp function returns the []rune terminates with '\0'
func strNdUp(b []rune, len int) []rune {
	res := make([]rune, len)
	copy(res, b)
	res = append(res, 0)
	return res
}

// Consumes the current token if it matches 's'.
func equal(tok *Token, s string) bool {
	return tok.Str == s
}

// if the current token is an expected symbol, the read position
// of token exceed one token.
func skip(tok *Token, s string) *Token {
	// defer printCurTok()
	if !equal(tok, s) {
		panic("\n" + errorTok(tok, "'%s' expected", string(s)))
	}
	return tok.Next
}

// consume returns token(pointer), if the current token is expected
// symbol, the read position of token exceed one.
func consume(rest **Token, tok *Token, s string) bool {
	// defer printCurTok()
	if equal(tok, s) {
		*rest = tok.Next
		return true
	}
	*rest = tok
	return false
}

// consumeIdent returns the current token if it is an identifier,
// and the read position of token exceed one.
func consumeIdent() *Token {
	// defer printCurTok()
	if token.Kind != TK_IDENT {
		return nil
	}
	t := token
	token = token.Next
	return t
}

// if next token is integer, the read position of token exceed one
// character or report an error.
func expectNumber() int64 {
	// defer printCurTok()
	if token.Kind != TK_NUM {
		panic("\n" + errorTok(token, "is not a number"))
	}
	val := token.Val
	token = token.Next
	return val
}

func expectIdent() string {
	// defer printCurTok()
	if token.Kind != TK_IDENT {
		panic("\n" + errorTok(token, "expect an identifier"))
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
	tok := &Token{Kind: kind, Str: str, Len: len, Loc: curIdx}
	cur.Next = tok
	return tok
}

// startsWith compare 'p' and 'q' , qq is keyword
func startsWith(p, q string) bool {
	return len(p) >= len(q) && p[:len(q)] == q
}

func startsWithTypeName(p string) string {
	var tyName = []string{"int16", "int64", "int",
		"uint8", "byte", "bool", "rune", "string",
	} // <-順番が大事、intとint16ではint16が先

	for _, k := range tyName {
		if startsWith(p, k) {
			return k
		}
	}
	return ""
}

func startsWithTermKw(p string) string {
	var term = []string{"break", "continue", "fallthrough", "return", "++", "--"}

	for _, k := range term {
		if startsWith(p, k) {
			return k
		}
	}
	return ""
}

func startsWithReserved(p string) string {
	// reserved words
	if k := startsWithTypeName(p); k != "" {
		return k
	}

	if k := startsWithTermKw(p); k != "" {
		return k
	}

	kw := []string{
		"if", "else", "for", "type", "var", "func", "struct",
		"goto", "switch", "case", "default",
		"true", "false",
		"nil", "Sizeof"}
	// unimplemented:
	// "chan", "const", "defer", "fallthrough", "interface", "map", "package", "range", "select"
	// "complex64", "complex128", "error",
	// "float32", "float64", "int8",
	// "uint", "uint16", "uint32", "uint64", "uintptr"
	// "iota"
	// "append", "cap", "close", "complex", "copy", "delete", "imag",
	// "len", "make", "new", "panic", "print", "println", "real", "recover"

	for _, k := range kw {
		if startsWith(p, k) && len(p) >= len(k) && !isIdent2(rune(p[len(k)])) {
			return k
		}
	}

	// Multi-letter punctuator
	ops := []string{"<<=", ">>=", "==", "!=", "<=", ">=", "->", "++", "--",
		"<<", ">>", "+=", "-=", "*=", "/=", ":=", "&&", "||"}

	for _, op := range ops {
		if startsWith(p, op) {
			return op
		}
	}
	return ""
}

func isSpace(op rune) bool {
	return strings.Contains("\n\t\v\f\r ", string(op))
}

func isDigit(op rune) bool {
	return '0' <= op && op <= '9'
}

func isIdent1(c rune) bool {
	return ('a' <= c && c <= 'z') ||
		('A' <= c && c <= 'Z') ||
		(c == '_')
}

func isIdent2(c rune) bool {
	return isIdent1(c) || isDigit(c)
}

func readDigit(cur *Token) (*Token, error) {
	var sVal string
	for ; curIdx < len(userInput) && isDigit(userInput[curIdx]); curIdx++ {
		sVal += string(userInput[curIdx])
	}
	cur = newToken(TK_NUM, cur, sVal, len(sVal))
	v, err := strconv.ParseInt(sVal, 10, 64)
	if err != nil {
		return nil, err
	}
	cur.Val = v
	return cur, nil
}

func getEscapeChar(c rune) rune {

	switch c {
	case 'a':
		return '\a'
	case 'b':
		return '\b'
	case 't':
		return '\t'
	case 'n':
		return '\n'
	case 'v':
		return '\v'
	case 'f':
		return '\f'
	case 'r':
		return '\r'
	case 'e':
		return 27
	case '0':
		return 0
	default:
		return c
	}
}

func readStringLiteral(cur *Token) (*Token, error) {
	p := 0

	buf := make([]rune, 0, 1024)
	for ; curIdx < len(userInput); curIdx++ {
		if userInput[curIdx] == 0 {
			return nil, fmt.Errorf(
				"tokenize(): err:\n%s",
				errorAt(curIdx+p, "unclosed string literal"),
			)
		}
		if userInput[curIdx] == '"' {
			break
		}

		if userInput[curIdx] == '\\' {
			curIdx++
			buf = append(buf, getEscapeChar(userInput[curIdx]))
		} else {
			buf = append(buf, userInput[curIdx])
		}
	}

	tok := newToken(TK_STR, cur, string(buf), len(buf))
	tok.Contents = strNdUp(buf, len(buf))
	tok.ContLen = len(buf) + 1
	curIdx++
	return tok, nil
}

func readCharLiteral(cur *Token, start int) (*Token, error) {
	p := start + 1
	if p < len(userInput) && userInput[p] == 0 {
		return nil, fmt.Errorf(
			"tokenize(): err:\n%s",
			errorAt(curIdx, "unclosed char literal"),
		)
	}

	var c rune
	if userInput[p] == '\\' {
		p++
		c = getEscapeChar(userInput[p])
		p++
	} else {
		c = userInput[p]
		p++
	}

	if userInput[p] != '\'' {
		return nil, fmt.Errorf(
			"tokenize(): err:\n%s",
			errorAt(curIdx, "char literal too long"),
		)
	}
	p++

	tok := newToken(TK_NUM, cur, string(userInput[start:p]), p-start)
	tok.Val = int64(c)
	return tok, nil
}

func isTermOfProd(cur *Token) bool {
	if curIdx == len(userInput) || userInput[curIdx] == '\n' {
		return cur.Kind == TK_IDENT ||
			cur.Kind == TK_NUM ||
			cur.Kind == TK_STR ||
			(cur.Kind == TK_RESERVED &&
				(startsWithTermKw(cur.Str) != "" ||
					startsWithTypeName(cur.Str) != "" ||
					strings.Contains(")]}", cur.Str)))
	}
	return false
}

// addSemiColn adds ";" token as terminators
// Reference: https://golang.org/ref/spec#Semicolons
func addSemiColn(cur *Token) *Token {
	if isTermOfProd(cur) {
		return newToken(TK_RESERVED, cur, ";", 0)
	}
	return cur
}

// Initialize lineinfo for all tokens.
func addLineNumbers(tok *Token) {
	var n int = 1

	for i := 0; i < len(userInput); i++ {
		if i == tok.Loc {
			tok.LineNo = n
			tok = tok.Next
		}
		if userInput[i] == '\n' {
			n++
		}
	}
}

// tokenize inputted string 'userInput', and return new tokens.
func tokenize(filename string) (*Token, error) {
	var head Token
	head.Next = nil
	cur := &head

	// for printToken
	headTok = &head

	for curIdx < len(userInput) && userInput[curIdx] != 0 {

		// skip space(s)
		if isSpace(userInput[curIdx]) {
			curIdx++
			continue
		}

		// skip line comment
		if startsWith(string(userInput[curIdx:]), "//") {
			curIdx += 2
			for ; curIdx < len(userInput) && userInput[curIdx] != '\n'; curIdx++ {
				// skip to the end of line.
			}
			continue
		}

		// skip block comment
		if startsWith(string(userInput[curIdx:]), "/*") {
			// skip to the first of '*/' in userInput[curIdx:].
			idx := strings.Index(string(userInput[curIdx:]), "*/")
			if idx == -1 {
				return nil, fmt.Errorf(
					"tokenize(): err:\n%s",
					errorAt(curIdx, "unclosed block comment"),
				)
			}
			curIdx += idx + 2
			continue
		}

		// keyword or multi-letter punctuator
		kw := startsWithReserved(string(userInput[curIdx:]))
		if kw != "" {
			cur = newToken(TK_RESERVED, cur, kw, len(kw))
			curIdx += len(kw)
			cur = addSemiColn(cur)
			continue
		}

		// single-letter punctuator
		if strings.Contains("+-()*/<>=;{},&[].,!|^:?", string(userInput[curIdx])) {
			cur = newToken(TK_RESERVED, cur, string(userInput[curIdx]), 1)
			curIdx++
			cur = addSemiColn(cur)
			continue
		}

		// identifier
		// if 'userInput[cutIdx]' is alphabets, it makes a token of TK_IDENT type.
		if isIdent1(userInput[curIdx]) {
			ident := make([]rune, 0, 20)
			for ; curIdx < len(userInput) && isIdent2(userInput[curIdx]); curIdx++ {
				ident = append(ident, userInput[curIdx])
			}
			cur = newToken(TK_IDENT, cur, string(ident), len(string(ident)))
			cur = addSemiColn(cur)
			continue
		}

		// string literal
		if userInput[curIdx] == '"' {
			curIdx++
			var err error
			cur, err = readStringLiteral(cur)
			if err != nil {
				return nil, err
			}
			cur = addSemiColn(cur)
			continue
		}

		// character literal
		if userInput[curIdx] == '\'' {
			var err error
			cur, err = readCharLiteral(cur, curIdx)
			if err != nil {
				return nil, err
			}
			curIdx += cur.Len
			cur = addSemiColn(cur)
			continue
		}

		// number
		if isDigit(userInput[curIdx]) {
			var err error
			cur, err = readDigit(cur)
			if err != nil {
				return nil, err
			}
			cur = addSemiColn(cur)
			continue
		}

		return nil, fmt.Errorf(
			"tokenize(): err:\n%s",
			errorAt(curIdx, "invalid token"),
		)
	}

	newToken(TK_EOF, cur, "", 0)
	addLineNumbers(head.Next)
	return head.Next, nil
}

func readFile(path string) ([]rune, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	br := bufio.NewReader(f)

	ret := make([]rune, 0, 1064)
	for {
		ru, sz, err := br.ReadRune()
		if sz == 0 || err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		ret = append(ret, ru)
	}
	ret = append(ret, 0)
	return ret, nil
}

func tokenizeFile(path string) (*Token, error) {
	var err error
	userInput, err = readFile(path)
	if err != nil {
		return nil, err
	}
	return tokenize(path)
}
