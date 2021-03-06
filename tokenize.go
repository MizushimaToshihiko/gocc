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

	TK_BLANKIDENT // 6: '_' identifier

	TK_EOF // 7: the end of tokens
)

type File struct {
	Name     string
	FileNo   int
	Contents []rune
}

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

	File   *File // Source location
	LineNo int   // Line number
	AtBol  bool  // True if this token is at begging of line
}

// Input file
var curFile *File

// A list of all input files.
var inputFiles []*File

// True if the current position is at the beginning of line.
var atBol bool

// current index in 'userInput'
var curIdx int

// for error report
//
// foo.go:
// 10:x = y + 1;
//        ^ <error message here>
func verrorAt(filename string, input []rune,
	lineNum, errIdx int, formt string, a ...interface{}) string {
	// get the start and end of the line 'errIdx' exists
	line := errIdx
	for 0 < line && input[line-1] != rune('\n') {
		line--
	}

	end := errIdx
	for end < len(input) && input[end] != rune('\n') {
		end++
	}

	// Show found lines along with file name and line number.
	res := fmt.Sprintf("\n%s:%d: ", filename, lineNum)
	indent := len(res)
	res += fmt.Sprintf("%.*s\n", end-line, string(input[line:end]))

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
	for i := 0; i < len(curFile.Contents); i++ {
		if curFile.Contents[i] == rune('\n') {
			lineNum++
		}
	}
	return verrorAt(curFile.Name, curFile.Contents, lineNum, errIdx, formt, a...)
}

func errorTok(tok *Token, formt string, ap ...interface{}) string {
	var errStr string

	if tok != nil {
		errStr = fmt.Sprintf("tok: '%s': kind: %d: pos :%d\n",
			tok.Str, tok.Kind, tok.Loc)
		errStr += verrorAt(tok.File.Name, tok.File.Contents,
			tok.LineNo, tok.Loc, formt, ap...)
	}

	return errStr
}

func warnTok(tok *Token, formt string, ap ...interface{}) {
	_, err := fmt.Fprintln(
		os.Stderr,
		verrorAt(tok.File.Name, tok.File.Contents, tok.LineNo, tok.Loc, formt, ap...),
	)
	if err != nil {
		panic(err)
	}
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

	for tok.Next.Kind == TK_COMM {
		tok.Next = tok.Next.Next
	}

	return tok.Next
}

// consume returns token(pointer), if the current token is expected
// symbol, the read position of token exceed one.
func consume(rest **Token, tok *Token, s string) bool {
	// defer printCurTok()
	if equal(tok, s) {
		for tok.Next.Kind == TK_COMM {
			tok.Next = tok.Next.Next
		}
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
	tok := &Token{
		Kind:  kind,
		Str:   str,
		Len:   len,
		Loc:   curIdx,
		File:  curFile,
		AtBol: atBol,
	}
	atBol = false
	cur.Next = tok
	return tok
}

// make new token and append to the end of cur.
func newToken2(kind TokenKind, cur *Token, str string, len int, loc int) *Token {
	tok := &Token{
		Kind:  kind,
		Str:   str,
		Len:   len,
		Loc:   loc,
		File:  curFile,
		AtBol: atBol,
	}
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
} // <-??????????????????int???int16??????int16??????

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
		} else if str[i] == '_' {
			continue
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

func parseFloat(str string) float64 {

	base := 10
	if startsWith(str, "0x") ||
		startsWith(str, "0X") {
		base = 16
		str = string(str[2:])
	}

	// Check literal
	if base == 10 && (contains(str, 'p') || contains(str, 'P')) {
		panic(errorAt(curIdx, "invalid number literal"))
	}

	// Read bofore the floating-point.
	var integer float64
	pos := strIndex(str, ".")
	end := len(str)
	if pos > 0 {
		integer = float64(parseInt(str[:pos], base))
	} else if pos == -1 {
		pos = 0
		if base == 10 {
			if contains(str, 'e') {
				end = strIndex(str, "e")
			} else if contains(str, 'E') {
				end = strIndex(str, "E")
			}
			integer = float64(parseInt(str[:end], base))
		} else if base == 16 {
			if contains(str, 'p') {
				end = strIndex(str, "p")
			} else if contains(str, 'P') {
				end = strIndex(str, "P")
			}
			integer = float64(parseInt(str[:end], base))
		}
	}

	// Read after the floating-point to the end or 'E' or 'P'.
	var float float64
	fbase := base
	i := pos + 1
	for ; i < len(str); i++ {
		if str[i] == '_' {
			if i+1 < len(str) && contains(".eEpP", rune(str[i+1])) {
				panic(errMustSeparateSuccessiveDigits(curIdx + i))
			}
			continue
		}

		if (str[i] == 'e' || str[i] == 'E') && base == 10 {
			break
		}
		if (str[i] == 'p' || str[i] == 'P') && base == 16 {
			break
		}

		var num float64
		if '0' <= str[i] && str[i] <= '9' {
			num = float64(str[i] - '0')
		} else if 'a' <= str[i] && str[i] <= 'f' {
			num = float64(str[i]-'a'+1) + 9
		} else if 'A' <= str[i] && str[i] <= 'F' {
			num = float64(str[i]-'A'+1) + 9
		}
		num /= float64(fbase)
		fbase *= fbase
		float += num
	}

	ret := integer + float
	if i >= len(str)-1 {
		if str[len(str)-1] != '_' {
			return ret
		}
		panic(errMustSeparateSuccessiveDigits(curIdx + i))
	}

	if str[i] == 'e' || str[i] == 'E' {
		if str[i+1] == '_' {
			panic(errMustSeparateSuccessiveDigits(curIdx + i))
		}
		pow := parseInt(string(str[i+2:]), base)
		if isDigit(rune(str[i+1])) {
			pow = parseInt(string(str[i+1:]), base)
		}
		i++
		if str[i] == '+' || isDigit(rune(str[i])) {
			var j int64
			for j = 0; j < pow; j++ {
				ret *= float64(base)
			}
		} else if str[i] == '-' {
			var j int64
			for j = 0; j < pow; j++ {
				ret /= float64(base)
			}
		}
	} else if (str[i] == 'p' || str[i] == 'P') && base == 16 {
		if str[i+1] == '_' {
			panic(errMustSeparateSuccessiveDigits(curIdx + i))
		}
		pow := parseInt(string(str[i+2:]), base)
		if isDigit(rune(str[i+1])) {
			pow = parseInt(string(str[i+1:]), base)
		}
		i++
		if str[i] == '+' || isDigit(rune(str[i])) {
			var j int64
			for j = 0; j < pow; j++ {
				ret *= float64(2)
			}
		} else if str[i] == '-' {
			var j int64
			for j = 0; j < pow; j++ {
				ret /= float64(2)
			}
		}
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

func isIdent3(c rune) bool {
	return ('a' <= c && c <= 'z') ||
		('A' <= c && c <= 'Z')
}

// for integer literal error
func errMustSeparateSuccessiveDigits(idx int) error {
	return fmt.Errorf(errorAt(idx, "'_' must separate successive digits"))
}

func readIntLiteral(cur *Token) (*Token, error) {
	var base int = 10

	var sVal string
	var err error
	var startIdx = curIdx

	if startsWith(string(curFile.Contents[curIdx:curIdx+2]), "0x") ||
		startsWith(string(curFile.Contents[curIdx:curIdx+2]), "0X") {
		base = 16
		curIdx += 2
		sVal, err = readHexDigit()
		if err != nil {
			return nil, err
		}
	} else if startsWith(string(curFile.Contents[curIdx:curIdx+2]), "0b") ||
		startsWith(string(curFile.Contents[curIdx:curIdx+2]), "0B") {
		base = 2
		curIdx += 2
	} else if startsWith(string(curFile.Contents[curIdx:curIdx+2]), "0o") ||
		startsWith(string(curFile.Contents[curIdx:curIdx+2]), "0O") ||
		startsWith(string(curFile.Contents[curIdx:curIdx+2]), "0_") {
		base = 8
		curIdx += 2
	} else if startsWith(string(curFile.Contents[curIdx:curIdx+1]), "0") &&
		isDigit(curFile.Contents[curIdx+1]) {
		base = 8
		curIdx += 1
	}

	for ; curIdx < len(curFile.Contents) &&
		(isDigit(curFile.Contents[curIdx]) || curFile.Contents[curIdx] == '_'); curIdx++ {

		if curFile.Contents[curIdx-1] == '_' && curFile.Contents[curIdx] == '_' {
			return nil, errMustSeparateSuccessiveDigits(startIdx)
		}

		if isDigit(curFile.Contents[curIdx]) {
			sVal += string(curFile.Contents[curIdx])
		}
	}

	if curFile.Contents[curIdx-1] == '_' {
		return nil, errMustSeparateSuccessiveDigits(startIdx)
	}

	cur = newToken(TK_NUM, cur, string(curFile.Contents[startIdx:curIdx]), curIdx-startIdx+1)
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

	if !contains(".eEfF", curFile.Contents[curIdx]) {
		return tok, nil
	}

	var sVal string = string(curFile.Contents[curIdx])
	curIdx++
	for isDigit(curFile.Contents[curIdx]) ||
		contains("eEfFpP_", curFile.Contents[curIdx]) ||
		(contains("EePp", curFile.Contents[curIdx-1]) &&
			contains("+-", curFile.Contents[curIdx])) {

		if (curFile.Contents[curIdx-1] == '_' && !isDigit(curFile.Contents[curIdx])) ||
			(curFile.Contents[curIdx] == '_' && !isDigit(curFile.Contents[curIdx-1])) {
			return nil, errMustSeparateSuccessiveDigits(curIdx)
		}

		sVal += string(curFile.Contents[curIdx])
		curIdx++
	}

	fval := parseFloat(tok.Str + sVal)

	ty := ty_double

	tok.FVal = fval
	tok.Str += sVal
	tok.Ty = ty
	return tok, nil
}

func readHexDigit() (string, error) {
	var sVal string
	var errIdx = curIdx
	for ; isxdigit(curFile.Contents[curIdx]) ||
		curFile.Contents[curIdx] == '_'; curIdx++ {

		if curFile.Contents[curIdx-1] == '_' && curFile.Contents[curIdx] == '_' {
			return "", errMustSeparateSuccessiveDigits(errIdx)
		}

		if isxdigit(curFile.Contents[curIdx]) {
			sVal += string(curFile.Contents[curIdx])
		}
	}

	if curFile.Contents[curIdx-1] == '_' {
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
	if '0' <= curFile.Contents[idx] && curFile.Contents[idx] <= '7' {
		// Read octal number.
		c := curFile.Contents[idx] - '0'
		idx++
		if '0' <= curFile.Contents[idx] && curFile.Contents[idx] <= '7' {
			c = (c << 3) + (curFile.Contents[idx] - '0')
			idx++
			if '0' <= curFile.Contents[idx] && curFile.Contents[idx] <= '7' {
				c = (c << 3) + (curFile.Contents[idx] - '0')
				idx++
			}
		}
		*newPos = idx
		return c, nil
	}

	if curFile.Contents[idx] == 'x' {
		// Read hexadecimal number.
		idx++
		if !isxdigit(curFile.Contents[idx]) {
			return -1, fmt.Errorf(
				"tokenize(): err:\n%s",
				errorAt(idx, "invalid hex escape sequence"))
		}
		var c rune
		for ; isxdigit(curFile.Contents[idx]); idx++ {
			c = (c << 4) + rune(fromHex(int(curFile.Contents[idx])))
		}
		*newPos = idx
		return c, nil
	}

	*newPos = idx

	switch curFile.Contents[idx] {
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
		return curFile.Contents[idx], nil
	}
}

// stringLiteralEnd finds a closing double-quote.
func stringLiteralEnd(idx int) (int, error) {
	var start int = idx
	for ; curFile.Contents[idx] != '"'; idx++ {
		if curFile.Contents[idx] == '\n' || curFile.Contents[idx] == 0 {
			return -1, fmt.Errorf(
				"tokenize(): stringLiteralEnd: err:\n%s",
				errorAt(start, "unclosed string literal"),
			)
		}
		if curFile.Contents[idx] == '\\' {
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
		if curFile.Contents[idx] == '\\' {
			c, err = getEscapeChar(&idx, idx+1)
			if err != nil {
				return nil, err
			}

			buf = append(buf, int64(c))
			continue
		}

		buf = append(buf, int64(curFile.Contents[idx]))
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
	if idx < len(curFile.Contents) && curFile.Contents[idx] == 0 {
		return nil, fmt.Errorf(
			"tokenize(): err:\n%s",
			errorAt(idx, "EOF: unclosed char literal"),
		)
	}

	var c rune
	var err error
	if curFile.Contents[idx] == '\\' {
		c, err = getEscapeChar(&idx, idx+1)
		if err != nil {
			return nil, err
		}
		idx++
	} else {
		c = curFile.Contents[idx]
		idx++
	}

	if curFile.Contents[idx] != '\'' {
		return nil, fmt.Errorf(
			"tokenize(): err:\n%s",
			errorAt(idx, "char literal too long"),
		)
	}
	idx++

	tok := newToken(TK_NUM, cur, string(curFile.Contents[start:idx]), idx-start)
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

	for i := 0; i < len(curFile.Contents); i++ {
		for i == tok.Loc && i < len(curFile.Contents) {
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
		if curFile.Contents[i] == '\n' {
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
func tokenize(file *File) (*Token, error) {
	curFile = file

	var head Token
	head.Next = nil
	cur := &head

	atBol = true

	for curIdx < len(curFile.Contents) && curFile.Contents[curIdx] != 0 {

		// skip space(s)
		if isSpace(curFile.Contents[curIdx]) {
			curIdx++
			continue
		}

		// new line
		if curFile.Contents[curIdx] == '\n' {
			cur = newToken(TK_NL, cur, "", 0)
			curIdx++
			atBol = true
			continue
		}

		// skip line comment
		if startsWith(string(curFile.Contents[curIdx:]), "//") {
			curIdx += 2
			for ; curIdx < len(curFile.Contents) && curFile.Contents[curIdx] != '\n'; curIdx++ {
				// skip to the end of line.
			}
			cur = newToken(TK_COMM, cur, "<line comment>", 0)
			continue
		}

		// skip block comment
		if startsWith(string(curFile.Contents[curIdx:]), "/*") {
			// skip to the first of '*/' in userInput[curIdx:].
			isatBol := atBol // The block comment is at the beginning of line, or not
			idx := strIndex(string(curFile.Contents[curIdx:]), "*/")
			if idx == -1 {
				return nil, fmt.Errorf(
					"tokenize(): err:\n%s",
					errorAt(curIdx, "unclosed block comment"),
				)
			}
			cur = newToken(TK_COMM, cur, "<block comment>", 0)
			curIdx += idx + 2
			atBol = isatBol
			continue
		}

		// blank identifier
		if contains("_", curFile.Contents[curIdx]) && !isIdent3(curFile.Contents[curIdx+1]) {
			cur = newToken(TK_BLANKIDENT, cur, string(curFile.Contents[curIdx]), 1)
			curIdx++
			continue
		}

		// identifier
		// if 'userInput[cutIdx]' is alphabets, it makes a token of TK_IDENT type.
		if isIdent1(curFile.Contents[curIdx]) {
			ident := make([]rune, 0, 20)
			for ; curIdx < len(curFile.Contents) && isIdent2(curFile.Contents[curIdx]); curIdx++ {
				ident = append(ident, curFile.Contents[curIdx])
			}
			cur = newToken(TK_IDENT, cur, string(ident), len(string(ident)))
			continue
		}

		// string literal
		if curFile.Contents[curIdx] == '"' {
			var err error
			cur, err = readStringLiteral(curIdx, cur)
			if err != nil {
				return nil, err
			}
			continue
		}

		// character literal
		if curFile.Contents[curIdx] == '\'' {
			var err error
			cur, err = readCharLiteral(cur)
			if err != nil {
				return nil, err
			}
			continue
		}

		// number
		if isDigit(curFile.Contents[curIdx]) ||
			(curFile.Contents[curIdx] == '.' && isDigit(curFile.Contents[curIdx+1])) {
			var err error
			cur, err = readNumber(cur)
			if err != nil {
				return nil, err
			}
			continue
		}

		// keyword or multi-letter punctuator
		kw := startsWithPunctuator(string(curFile.Contents[curIdx:]))
		if kw != "" {
			cur = newToken(TK_RESERVED, cur, kw, len(kw))
			curIdx += len(kw)
			continue
		}
		// single-letter punctuator
		if contains("+-()*/<>=;{},&[].!|^:?%#", curFile.Contents[curIdx]) {
			cur = newToken(TK_RESERVED, cur, string(curFile.Contents[curIdx]), 1)
			curIdx++
			continue
		}

		return nil, fmt.Errorf(
			"tokenize(): err:\ncurIdx: %s\n%s", string(curFile.Contents[curIdx]),
			errorAt(curIdx, "invalid token"),
		)
	}

	newToken(TK_EOF, cur, "", 0)

	addSemiColn(head.Next)
	delNewLineTok(head.Next)
	addLineNumbers(head.Next)
	curIdx = 0
	return head.Next, nil
}

func readFile(path string) ([]rune, error) {
	var r io.Reader
	switch path {
	case "-":
		r = os.Stdin
	default:
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatal(err)
			}
		}()
		r = f
	}

	br := bufio.NewReader(r)

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

func getInputFiles() []*File {
	return inputFiles
}

func newFile(name string, fileNo int, contents []rune) *File {
	return &File{Name: name, FileNo: fileNo, Contents: contents}
}

// For tokenizeFile function
var fileno int

func tokenizeFile(path string) (*Token, error) {

	p, err := readFile(path)
	if err != nil {
		return nil, err
	}

	file := newFile(path, fileno+1, p)

	// Save the filename for assembler .file directive.
	inputFiles = append(inputFiles, file)
	fileno++

	return tokenize(file)
}
