//
// tokenizier
//
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
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
	FVal float64   // if 'kind' is TK_NUM, it's value
	Loc  int       // the location in 'userInput'
	Ty   *Type     // Used if TK_NUM or TK_STR
	Str  string    // token string
	Len  int       // length of token

	Contents []int64 // string literal contents including terminating '\0'
	ContLen  int     // string literal length

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
func strNdUp(b []int64, len int) []int64 {
	res := make([]int64, len)
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
var tyName = []string{
	"int8", "int16", "int32", "int64", "int",
	"uint8", "uint16", "uint32", "uint64", "uint",
	"float32", "float64",
	"complex64", "complex128",
	"byte", "rune",
	"string", "bool", "error",
	"struct", "func",
} // <-順番が大事、intとint16ではint16が先

var term = []string{"break", "continue", "fallthrough",
	"return", "++", "--"}

var kw = []string{
	"case", "chan", "const", "default", "defer", "else", "for",
	"func", "go", "goto", "if", "import", "interface", "map",
	"package", "range", "switch", "select", "type",
	"var",
	"true", "false", "iota", "nil",
	"make", "len", "cap", "new", "append", "copy", "close",
	"delete", "complex", "real", "imag", "panic", "recover",
	"Sizeof",
}

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

func contains(str string, r rune) bool {
	s := []rune(str)
	for i := 0; i < len(s); i++ {
		if s[i] == r {
			return true
		}
	}

	return false
}

func strIndex(str string, substr string) int {
	for i := 0; i < len(str); i++ {
		if str[i] == substr[0] {
			flag := true
			for j := 1; j < len(substr); j++ {
				if str[i+j] != substr[j] {
					flag = false
					break
				}
			}
			if flag {
				return i
			}
		}
	}
	return -1
}

func parseInt(str string, base int) int64 {
	var ret int64
	digits := 0
	for i := len(str) - 1; i >= 0; i-- {
		var num int64
		if '0' <= str[i] && str[i] <= '9' {
			num = int64(str[i] - '0')
		} else if 'a' <= str[i] && str[i] <= 'f' {
			num = int64(str[i]-'a'+1) + 9
		} else if 'A' <= str[i] && str[i] <= 'F' {
			num = int64(str[i]-'A'+1) + 9
		} else {
			panic(errorAt(curIdx+i, "couldn't parse"))
		}

		for j := 0; j < digits; j++ {
			num *= int64(base)
		}
		digits++
		ret += num

	}
	return ret
}

func isSpace(op rune) bool {
	return contains("\t\v\f\r ", op)
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

// for integer literal error
func errMustSeparateSuccessiveDigits(idx int) error {
	return errors.New(errorAt(idx, "'_' must separate successive digits"))
}

func readIntLiteral(cur *Token) (*Token, error) {
	var base int = 10

	var sVal string
	var err error
	var startIdx = curIdx

	if startsWith(string(userInput[curIdx:curIdx+2]), "0x") ||
		startsWith(string(userInput[curIdx:curIdx+2]), "0X") {
		base = 16
		curIdx += 2
		sVal, err = readHexDigit()
		if err != nil {
			return nil, err
		}
	} else if startsWith(string(userInput[curIdx:curIdx+2]), "0b") ||
		startsWith(string(userInput[curIdx:curIdx+2]), "0B") {
		base = 2
		curIdx += 2
	} else if startsWith(string(userInput[curIdx:curIdx+2]), "0o") ||
		startsWith(string(userInput[curIdx:curIdx+2]), "0O") ||
		startsWith(string(userInput[curIdx:curIdx+2]), "0_") {
		base = 8
		curIdx += 2
	} else if startsWith(string(userInput[curIdx:curIdx+1]), "0") &&
		isDigit(userInput[curIdx+1]) {
		base = 8
		curIdx += 1
	}

	for ; curIdx < len(userInput) &&
		(isDigit(userInput[curIdx]) || userInput[curIdx] == '_'); curIdx++ {

		if userInput[curIdx-1] == '_' && userInput[curIdx] == '_' {
			return nil, errMustSeparateSuccessiveDigits(startIdx)
		}

		if isDigit(userInput[curIdx]) {
			sVal += string(userInput[curIdx])
		}
	}

	if userInput[curIdx-1] == '_' {
		return nil, errMustSeparateSuccessiveDigits(startIdx)
	}

	cur = newToken(TK_NUM, cur, string(userInput[startIdx:curIdx]), curIdx-startIdx+1)
	var v int64
	if sVal != "" {
		v = parseInt(sVal, base)
	}

	cur.Val = v

	cur.Ty = ty_long
	if v <= 2147483647 {
		cur.Ty = ty_int
	}

	return cur, nil
}

func readNumber(cur *Token) (*Token, error) {
	tok, err := readIntLiteral(cur)
	if err != nil {
		return nil, err
	}

	if !contains(".eEfF", userInput[curIdx]) {
		return tok, nil
	}

	var sVal string = string(userInput[curIdx])
	curIdx++
	for isDigit(userInput[curIdx]) ||
		contains("eEfFpP_", userInput[curIdx]) ||
		(contains("EePp", userInput[curIdx-1]) &&
			contains("+-", userInput[curIdx])) {

		if (userInput[curIdx-1] == '_' && !isDigit(userInput[curIdx])) ||
			(userInput[curIdx] == '_' && !isDigit(userInput[curIdx-1])) {
			return nil, errMustSeparateSuccessiveDigits(curIdx)
		}

		sVal += string(userInput[curIdx])
		curIdx++
	}

	fval, err := strconv.ParseFloat(tok.Str+sVal, 64)
	if err != nil {
		return nil, err
	}

	ty := ty_double

	tok.FVal = fval
	tok.Str += sVal
	tok.Ty = ty
	return tok, nil
}

func readHexDigit() (string, error) {
	var sVal string
	var errIdx = curIdx
	for ; isxdigit(userInput[curIdx]) ||
		userInput[curIdx] == '_'; curIdx++ {

		if userInput[curIdx-1] == '_' && userInput[curIdx] == '_' {
			return "", errMustSeparateSuccessiveDigits(errIdx)
		}

		if isxdigit(userInput[curIdx]) {
			sVal += string(userInput[curIdx])
		}
	}

	if userInput[curIdx-1] == '_' {
		return "", errMustSeparateSuccessiveDigits(errIdx)
	}

	return sVal, nil
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

func mustToString(s []int64) string {
	var ret string
	for i := 0; i < len(s); i++ {
		ret += string(rune(s[i]))
	}
	return ret
}

func readStringLiteral(start int, cur *Token) (*Token, error) {
	var end int
	var err error
	var idx int
	end, err = stringLiteralEnd(start + 1)
	if err != nil {
		return nil, err
	}

	buf := make([]int64, 0, 1024)
	for idx = start + 1; idx < end; idx++ {

		var c rune
		if userInput[idx] == '\\' {
			c, err = getEscapeChar(&idx, idx+1)
			if err != nil {
				return nil, err
			}

			buf = append(buf, int64(c))
			continue
		}

		buf = append(buf, int64(userInput[idx]))
	}
	idx++

	tok := newToken(TK_STR, cur, mustToString(buf), end-start+1)
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
	tok.Ty = ty_int
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
					startsWith(cur.Str, ")") ||
					startsWith(cur.Str, "]") ||
					startsWith(cur.Str, "}")))

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
			tok = tok.Next
			if tok == nil {
				return
			}
			if tok.Str == ";" {
				l := tok.Len
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
			idx := strIndex(string(userInput[curIdx:]), "*/")
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
		if isDigit(userInput[curIdx]) ||
			(userInput[curIdx] == '.' && isDigit(userInput[curIdx+1])) {
			var err error
			cur, err = readNumber(cur)
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
		if contains("+-()*/<>=;{},&[].!|^:?%", userInput[curIdx]) {
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
