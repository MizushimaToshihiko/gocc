//
// tokenizier
//
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// set TokenKind with enum
type TokenKind int

const (
	TK_RESERVED TokenKind = iota // 0: Reserved words, and puncturators
	TK_SIZEOF                    // 1: 'sizeof' operator
	TK_IDENT                     // 2: idenfier such as variables, function names
	TK_STR                       // 3: string literals
	TK_NUM                       // 4: integer
	TK_EOF                       // 5: the end of tokens
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
}

// current token
var token *Token

// inputted program
var userInput []rune

// current index in 'userInput'
var curIdx int

// for error report
// it's arguments are same as printf
func errorAt(errIdx int, formt string, a ...interface{}) string {
	// get the start and end of the line 'errIdx' exists
	line := errIdx
	for 0 < line && userInput[line-1] != rune('\n') {
		line--
	}

	end := errIdx
	for end < len(userInput) && userInput[end] != rune('\n') {
		end++
	}

	// Find out what line the found line is in the whole.
	lineNum := 1
	for i := 0; i < line; i++ {
		if userInput[i] == rune('\n') {
			lineNum++
		}
	}

	// Show found lines along with file name and line number.
	res := fmt.Sprintf("%s:%d: ", filename, lineNum)
	indent := len(res)
	res += fmt.Sprintf("%.*s\n", end-line, string(userInput[line:end]))

	// Point the error location with "^" and display the error message.
	pos := errIdx - line + indent

	return res + fmt.Sprintf("%*s", pos, " ") +
		"^ " +
		fmt.Sprintf(formt, a...) +
		"\n"
}

func errorTok(tok *Token, formt string, ap ...interface{}) string {
	var errStr string
	if tok != nil {
		errStr += errorAt(tok.Loc, formt, ap...)
	}

	return errStr +
		fmt.Sprintf(formt, ap...) +
		"\n"
}

func warnTok(tok *Token, frmt string, ap ...string) {
	var errStr string
	if tok != nil {
		errStr += errorAt(tok.Loc, frmt, ap)
	} else {
		errStr += fmt.Sprintf(frmt, ap) + "\n"
	}
	fmt.Fprint(os.Stderr, errStr)
}

// strNdUp function returns the []rune terminates with '\0'
func strNdUp(b []rune, len int) []rune {
	res := make([]rune, len)
	copy(res, b)
	res = append(res, 0)
	return res
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

// consume returns token(pointer), if the current token is expected
//  symbol, the read position of token exceed one character.
func consume(s string) *Token {
	// defer printCurTok()
	if peek(s) == nil {
		return nil
	}
	t := token
	token = token.Next
	return t
}

// consumeIdent returns the current token if it is an identifier
func consumeIdent() *Token {
	// defer printCurTok()
	if token.Kind != TK_IDENT {
		return nil
	}
	t := token
	token = token.Next
	return t
}

// consumeSizeof returns the token(pointer) and proceed to the next token,
//  if the current token is "sizeof".
func consumeSizeof() *Token {

	if token.Kind != TK_SIZEOF ||
		token.Len != len("sizeof") ||
		token.Str != "sizeof" {
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
		panic("\n" + errorAt(token.Loc, "is not '%s'", string(s)))
	}
	token = token.Next
}

// if next token is integer, the read position of token exceed one
// character or report an error.
func expectNumber() int64 {
	// defer printCurTok()
	if token.Kind != TK_NUM {
		panic("\n" + errorAt(token.Loc, "is not a number"))
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

func startsWithReserved(p string) string {
	// reserved words
	kw := []string{
		"return", "if", "else", "while", "for", "type", "var", "func", "struct",
		"break", "continue", "goto", "switch", "case", "default",
		"true", "false",
		"nil",
		"int", "int64", "uint8", "bool", "rune"}
	// unimplemented:
	// "chan", "const", "defer", "fallthrough", "interface", "map", "package", "range", "select"
	// "int32","bool", "byte", "complex64", "complex128", "error",
	// "float32", "float64", "int8", "int16", "int32", "int64",
	// "string", "uint", "uint16", "uint32", "uint64", "uintptr"
	// "iota"
	// "append", "cap", "close", "complex", "copy", "delete", "imag",
	// "len", "make", "new", "panic", "print", "println", "real", "recover"

	for _, k := range kw {
		if startsWith(p, k) && len(p) >= len(k) && !isAlNum(rune(p[len(k)])) {
			return k
		}
	}

	// Multi-letter punctuator
	ops := []string{"<<=", ">>=", "==", "!=", "<=", ">=", "->", "++", "--",
		"<<", ">>", "+=", "-=", "*=", "/=", "&&", "||"}

	for _, op := range ops {
		if startsWith(p, op) {
			return op
		}
	}
	return ""
}

func isSpace(op rune) bool {
	return strings.Contains("\t\n\v\f\r ", string(op))
}

func isDigit(op rune) bool {
	return '0' <= op && op <= '9'
}

func isAlpha(c rune) bool {
	return ('a' <= c && c <= 'z') ||
		('A' <= c && c <= 'Z') ||
		(c == '_')
}

func isAlNum(c rune) bool {
	return isAlpha(c) || ('0' <= c && c <= '9')
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

// tokenize inputted string 'userInput', and return new tokens.
func tokenize() (*Token, error) {
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
			}
			continue
		}

		// skip block comment
		if startsWith(string(userInput[curIdx:]), "/*") {
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

		// 'sizeof' keyword
		if startsWith(string(userInput[curIdx:]), "sizeof") {
			cur = newToken(TK_SIZEOF, cur, "sizeof", len("sizeof"))
			curIdx += len("sizeof")
			continue
		}

		// keyword or multi-letter punctuator
		kw := startsWithReserved(string(userInput[curIdx:]))
		if kw != "" {
			cur = newToken(TK_RESERVED, cur, kw, len(kw))
			curIdx += len(kw)
			continue
		}

		// single-letter punctuator
		if strings.Contains("+-()*/<>=;{},&[].,!~|^:?", string(userInput[curIdx])) {
			cur = newToken(TK_RESERVED, cur, string(userInput[curIdx]), 1)
			curIdx++
			continue
		}

		// identifier
		// if 'userInput[cutIdx]' is alphabets, it makes a token of TK_IDENT type.
		if isAlpha(userInput[curIdx]) {
			ident := make([]rune, 0, 20)
			for ; curIdx < len(userInput) && isAlNum(userInput[curIdx]); curIdx++ {
				ident = append(ident, userInput[curIdx])
			}
			cur = newToken(TK_IDENT, cur, string(ident), len(string(ident)))
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
			continue
		}

		// number
		if isDigit(userInput[curIdx]) {
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
			continue
		}

		return nil, fmt.Errorf(
			"tokenize(): err:\n%s",
			errorAt(curIdx, "invalid token"),
		)
	}

	newToken(TK_EOF, cur, "", 0)
	return head.Next, nil
}
