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

	TK_COMM // 4: comment
	TK_NL   // 5: new line

	TK_EOF // 6: the end of tokens
)

type Token struct {
	Kind TokenKind // type of token
	Next *Token    // next
	Val  int64     // if 'kind' is TK_NUM, it's integer
	Loc  int       // the location in 'userInput'
	Ty   *Type     // Used if TK_STR
	Str  string    // token string
	Len  int       // length of token

	Contents []rune // string literal contents including terminating '\0'
	ContLen  int    // string literal length

	LineNo int // Line number
}

var curFilename string

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
		errStr = fmt.Sprintf("tok: '%s': kind: %d: pos :%d\n", tok.Str, tok.Kind, tok.Loc)
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
	printCurTok(tok)
	printCalledFunc()

	return tok.Str == s
}

// if the current token is an expected symbol, the read position
// of token exceed one token.
func skip(tok *Token, s string) *Token {
	printCurTok(tok)
	printCalledFunc()

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

func atEof(tok *Token) bool {
	return tok.Kind == TK_EOF
}

// make new token and append to the end of cur.
func newToken(kind TokenKind, cur *Token, str string, len int) *Token {
	tok := &Token{Kind: kind, Str: str, Len: len, Loc: curIdx}
	cur.Next = tok
	return tok
}

// make new token and append to the end of cur.
func newToken2(kind TokenKind, cur *Token, str string, len int, loc int) *Token {
	tok := &Token{Kind: kind, Str: str, Len: len, Loc: loc}
	cur.Next = tok
	return tok
}

// startsWith compare 'p' and 'q' , q is keyword
func startsWith(p, q string) bool {
	return len(p) >= len(q) && p[:len(q)] == q
}

// reserved words
var tyName = []string{"int16", "int64", "int",
	"uint8", "byte", "bool", "rune", "string",
} // <-順番が大事、intとint16ではint16が先

var term = []string{"break", "continue", "fallthrough",
	"return", "++", "--"}

var kw = []string{
	"if", "else", "for", "type", "var", "func", "struct",
	"goto", "switch", "case", "default", "package", "import",
	"true", "false", "nil", "Sizeof"}

// unimplemented:
// "chan", "const", "defer", "fallthrough", "interface", "map", "package", "range", "select"
// "complex64", "complex128", "error",
// "float32", "float64", "int8",
// "uint", "uint16", "uint32", "uint64", "uintptr"
// "iota"
// "append", "cap", "close", "complex", "copy", "delete", "imag",
// "len", "make", "new", "panic", "print", "println", "real", "recover"

func startsWithTypeName(p string) string {
	for _, k := range tyName {
		if startsWith(p, k) {
			return k
		}
	}
	return ""
}

func startsWithTermKw(p string) string {
	for _, k := range term {
		if startsWith(p, k) {
			return k
		}
	}
	return ""
}

func startsWithPunctuator(p string) string {

	// Multi-letter punctuator
	ops := []string{
		"<<=", ">>=", "==", "!=", "<=", ">=", "->", "++", "--",
		"<<", ">>", "+=", "-=", "*=", "/=", "%=", "&=", "|=", "^=",
		":=", "&&", "||", "...",
	}

	for _, op := range ops {
		if startsWith(p, op) {
			return op
		}
	}
	return ""
}

func isKw(tok *Token) bool {
	for _, k := range tyName {
		if tok.Str == k {
			return true
		}
	}

	for _, k := range term {
		if tok.Str == k {
			return true
		}
	}

	for _, k := range kw {
		if tok.Str == k {
			return true
		}
	}

	return false
}

func isSpace(op rune) bool {
	return strings.Contains("\t\v\f\r ", string(op))
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
	var base int = 10

	if startsWith(string(userInput)[curIdx:curIdx+2], "0x") ||
		startsWith(string(userInput)[curIdx:curIdx+2], "0X") {
		curIdx += 2
		return readHexDigit(cur)
	}

	var sVal string

	if startsWith(string(userInput)[curIdx:curIdx+2], "0b") ||
		startsWith(string(userInput)[curIdx:curIdx+2], "0B") {
		base = 2
		curIdx += 2
	}

	if startsWith(string(userInput)[curIdx:curIdx+2], "0o") ||
		startsWith(string(userInput)[curIdx:curIdx+2], "0O") {
		base = 8
		curIdx += 2
	}

	for ; curIdx < len(userInput) && isDigit(userInput[curIdx]); curIdx++ {
		sVal += string(userInput[curIdx])
	}

	cur = newToken(TK_NUM, cur, sVal, len(sVal))
	v, err := strconv.ParseInt(sVal, base, 64)
	if err != nil {
		return nil, err
	}

	cur.Val = v
	return cur, nil
}

func readHexDigit(cur *Token) (*Token, error) {
	var sVal string
	for ('0' <= userInput[curIdx] && userInput[curIdx] <= '9') ||
		('A' <= userInput[curIdx] && userInput[curIdx] <= 'F') ||
		('a' <= userInput[curIdx] && userInput[curIdx] <= 'f') {
		sVal += string(userInput[curIdx])
		curIdx++
	}
	cur = newToken(TK_NUM, cur, sVal, len(sVal))
	v, err := strconv.ParseInt(sVal, 16, 64)
	if err != nil {
		return nil, err
	}
	cur.Val = v
	return cur, nil
}

func isxdigit(p rune) bool {
	return ('0' <= p && p <= '9') ||
		('A' <= p && p <= 'F') ||
		('a' <= p && p <= 'f')
}

func fromHex(c int) int {
	if '0' <= c && c <= '9' {
		return c - '0'
	}
	if 'a' <= c && c <= 'f' {
		return c - 'a' + 10
	}
	return c - 'A' + 10
}

func getEscapeChar(newPos *int, idx int) (rune, error) {
	if '0' <= userInput[idx] && userInput[idx] <= '7' {
		// Read octal number.
		c := userInput[idx] - '0'
		idx++
		if '0' <= userInput[idx] && userInput[idx] <= '7' {
			c = (c << 3) + (userInput[idx] - '0')
			idx++
			if '0' <= userInput[idx] && userInput[idx] <= '7' {
				c = (c << 3) + (userInput[idx] - '0')
				idx++
			}
		}
		*newPos = idx
		return c, nil
	}

	if userInput[idx] == 'x' {
		// Read hexadecimal number.
		idx++
		if !isxdigit(userInput[idx]) {
			return -1, fmt.Errorf(
				"tokenize(): err:\n%s",
				errorAt(idx, "invalid hex escape sequence"))
		}
		var c rune
		for ; isxdigit(userInput[idx]); idx++ {
			c = (c << 4) + rune(fromHex(int(userInput[idx])))
		}
		*newPos = idx
		return c, nil
	}

	*newPos = idx

	switch userInput[idx] {
	case 'a':
		return '\a', nil
	case 'b':
		return '\b', nil
	case 't':
		return '\t', nil
	case 'n':
		return '\n', nil
	case 'v':
		return '\v', nil
	case 'f':
		return '\f', nil
	case 'r':
		return '\r', nil
	case 'e':
		return 27, nil
	default:
		return userInput[idx], nil
	}
}

// stringLiteralEnd finds a closing double-quote.
// strings.Indexを使えば簡単なんだけど
func stringLiteralEnd(idx int) (int, error) {
	var start int = idx
	for ; userInput[idx] != '"'; idx++ {
		if userInput[idx] == '\n' || userInput[idx] == 0 {
			return -1, fmt.Errorf(
				"tokenize(): stringLiteralEnd: err:\n%s",
				errorAt(start, "unclosed string literal"),
			)
		}
		if userInput[idx] == '\\' {
			idx++
		}
	}
	return idx, nil
}

func readStringLiteral(start int, cur *Token) (*Token, error) {
	var end int
	var err error
	var idx int
	end, err = stringLiteralEnd(start + 1)
	if err != nil {
		return nil, err
	}

	buf := make([]rune, 0, 1024)
	for idx = start + 1; idx < end; idx++ {

		var c rune
		if userInput[idx] == '\\' {
			c, err = getEscapeChar(&idx, idx+1)
			if err != nil {
				return nil, err
			}

			buf = append(buf, c)
			continue
		}

		buf = append(buf, userInput[idx])
	}
	idx++

	tok := newToken(TK_STR, cur, string(buf), end-start+1)
	tok.Contents = strNdUp(buf, len(buf))
	tok.ContLen = len(buf) + 1
	tok.Ty = arrayOf(ty_char, len(buf)+1)
	curIdx += tok.Len
	return tok, nil
}

func readCharLiteral(cur *Token) (*Token, error) {
	start := curIdx
	idx := start + 1
	if idx < len(userInput) && userInput[idx] == 0 {
		return nil, fmt.Errorf(
			"tokenize(): err:\n%s",
			errorAt(idx, "EOF: unclosed char literal"),
		)
	}

	var c rune
	var err error
	if userInput[idx] == '\\' {
		c, err = getEscapeChar(&idx, idx+1)
		if err != nil {
			return nil, err
		}
		idx++
	} else {
		c = userInput[idx]
		idx++
	}

	if userInput[idx] != '\'' {
		return nil, fmt.Errorf(
			"tokenize(): err:\n%s",
			errorAt(idx, "char literal too long"),
		)
	}
	idx++

	tok := newToken(TK_NUM, cur, string(userInput[start:idx]), idx-start)
	tok.Val = int64(c)
	curIdx += tok.Len
	return tok, nil
}

func isTermOfProd(cur *Token) bool {
	if cur.Next != nil && cur.Next.Str != ";" &&
		(cur.Next.Kind == TK_COMM ||
			cur.Next.Kind == TK_NL ||
			cur.Next.Kind == TK_EOF) {
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
func addSemiColn(head *Token) {
	var cur *Token = head
	for cur != nil && cur.Kind != TK_EOF {
		if isTermOfProd(cur) {
			tmp := cur.Next
			cur.Next = newToken2(TK_RESERVED, cur, ";", 0, cur.Loc+cur.Len)
			cur.Next.Next = tmp
			cur = cur.Next.Next
		}
		cur = cur.Next
	}
}

//
func delNewLineTok(head *Token) {
	var cur *Token = head
	prev := cur
	for cur != nil && cur.Kind != TK_EOF {
		if cur.Kind == TK_NL {
			prev.Next = cur.Next
			cur = cur.Next
			continue
		}
		prev = cur
		cur = cur.Next
	}
}

// Initialize line info for all tokens.
func addLineNumbers(head *Token) {
	var tok *Token = head
	var n int = 1

	for i := 0; i < len(userInput); i++ {
		for i == tok.Loc && i < len(userInput) {
			tok.LineNo = n
			l := tok.Len
			tok = tok.Next
			if tok == nil {
				return
			}
			if tok.Str == ";" {
				tok.Loc = i + l
			}
		}
		if userInput[i] == '\n' {
			n++
		}
	}
}

func convKw(tok *Token) {
	for t := tok; t != nil; t = t.Next {
		if isKw(t) {
			t.Kind = TK_RESERVED
		}
	}
}

// tokenize inputted string 'userInput', and return new tokens.
func tokenize(filename string) (*Token, error) {
	curFilename = filename
	var head Token
	head.Next = nil
	cur := &head

	for curIdx < len(userInput) && userInput[curIdx] != 0 {

		// skip space(s)
		if isSpace(userInput[curIdx]) {
			curIdx++
			continue
		}

		// new line
		if userInput[curIdx] == '\n' {
			cur = newToken(TK_NL, cur, "", 0)
			curIdx++
			continue
		}

		// skip line comment
		if startsWith(string(userInput[curIdx:]), "//") {
			curIdx += 2
			for ; curIdx < len(userInput) && userInput[curIdx] != '\n'; curIdx++ {
				// skip to the end of line.
			}
			cur = newToken(TK_COMM, cur, "<line comment>", 0)
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
			cur = newToken(TK_COMM, cur, "<block comment>", 0)
			curIdx += idx + 2
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
			continue
		}

		// string literal
		if userInput[curIdx] == '"' {
			var err error
			cur, err = readStringLiteral(curIdx, cur)
			if err != nil {
				return nil, err
			}
			continue
		}

		// character literal
		if userInput[curIdx] == '\'' {
			var err error
			cur, err = readCharLiteral(cur)
			if err != nil {
				return nil, err
			}
			continue
		}

		// number
		if isDigit(userInput[curIdx]) {
			var err error
			cur, err = readDigit(cur)
			if err != nil {
				return nil, err
			}
			continue
		}

		// keyword or multi-letter punctuator
		kw := startsWithPunctuator(string(userInput[curIdx:]))
		if kw != "" {
			cur = newToken(TK_RESERVED, cur, kw, len(kw))
			curIdx += len(kw)
			continue
		}

		// single-letter punctuator
		if strings.Contains("+-()*/<>=;{},&[].,!|^:?%", string(userInput[curIdx])) {
			cur = newToken(TK_RESERVED, cur, string(userInput[curIdx]), 1)
			curIdx++
			continue
		}

		return nil, fmt.Errorf(
			"tokenize(): err:\ncurIdx: %s\n%s", string(userInput[curIdx]),
			errorAt(curIdx, "invalid token"),
		)
	}

	newToken(TK_EOF, cur, "", 0)

	addSemiColn(head.Next)
	delNewLineTok(head.Next)
	addLineNumbers(head.Next)
	convKw(head.Next)
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
