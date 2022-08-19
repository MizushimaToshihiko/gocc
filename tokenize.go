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

	TK_PP_NUM // 7: Preprocessing numbers

	TK_EOF // 8: the end of tokens
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

	File     *File    // Source location
	LineNo   int      // Line number
	AtBol    bool     // True if this token is at begging of line
	HasSpace bool     // True if this token follows a space character
	Hideset  *Hideset // For macro expansion
	Origin   *Token   // If this is expanded from a macro, the original token
}

// Input file
var curFile *File

// A list of all input files.
var inputFiles []*File

// True if the current position is at the beginning of line.
var atBol bool

// True if the current position follows a space character
var hasSpace bool

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

	for tok.Next != nil && tok.Next.Kind == TK_COMM {
		tok.Next = tok.Next.Next
	}

	return tok.Next
}

// consume returns token(pointer), if the current token is expected
// symbol, the read position of token exceed one.
func consume(rest **Token, tok *Token, s string) bool {
	// defer printCurTok()
	if equal(tok, s) {
		for tok.Next != nil && tok.Next.Kind == TK_COMM {
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
		Kind:     kind,
		Str:      str,
		Len:      len,
		Loc:      curIdx,
		File:     curFile,
		AtBol:    atBol,
		HasSpace: hasSpace,
	}
	atBol = false
	hasSpace = false
	cur.Next = tok
	return tok
}

// make new token and append to the end of cur, with the location.
func newToken2(kind TokenKind, cur *Token, str string, len int, loc int) *Token {
	tok := &Token{
		Kind:     kind,
		Str:      str,
		Len:      len,
		Loc:      loc,
		File:     curFile,
		AtBol:    atBol,
		HasSpace: hasSpace,
	}
	atBol = false
	hasSpace = false
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
		":=", "&&", "||", "...", "##",
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

func parseInt(str string, base, errIdx int) int64 {
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
			panic(errorAt(errIdx+i, "couldn't parse"))
		}

		for j := 0; j < digits; j++ {
			num *= int64(base)
		}
		digits++
		ret += num

	}
	return ret
}

func parseFloat(str string, errIdx int) float64 {

	base := 10
	if startsWith(str, "0x") ||
		startsWith(str, "0X") {
		base = 16
		str = string(str[2:])
	}

	// Check literal
	if base == 10 && (contains(str, 'p') || contains(str, 'P')) {
		panic(errorAt(errIdx, "invalid number literal"))
	}

	// Read bofore the floating-point.
	var integer float64
	pos := strIndex(str, ".")
	end := len(str)
	if pos > 0 {
		integer = float64(parseInt(str[:pos], base, errIdx))
	} else if pos == -1 {
		pos = 0
		if base == 10 {
			if contains(str, 'e') {
				end = strIndex(str, "e")
			} else if contains(str, 'E') {
				end = strIndex(str, "E")
			}
			integer = float64(parseInt(str[:end], base, errIdx))
		} else if base == 16 {
			if contains(str, 'p') {
				end = strIndex(str, "p")
			} else if contains(str, 'P') {
				end = strIndex(str, "P")
			}
			integer = float64(parseInt(str[:end], base, errIdx))
		}
	}

	// Read after the floating-point to the end or 'E' or 'P'.
	var float float64
	fbase := base
	i := pos + 1
	for ; i < len(str); i++ {
		if str[i] == '_' {
			if i+1 < len(str) && contains(".eEpP", rune(str[i+1])) {
				panic(errMustSeparateSuccessiveDigits(errIdx + i))
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
		pow := parseInt(string(str[i+2:]), base, errIdx)
		if isDigit(rune(str[i+1])) {
			pow = parseInt(string(str[i+1:]), base, errIdx)
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
		pow := parseInt(string(str[i+2:]), base, errIdx)
		if isDigit(rune(str[i+1])) {
			pow = parseInt(string(str[i+1:]), base, errIdx)
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

// for integer literal error
func errMustSeparateSuccessiveDigits(idx int) error {
	return fmt.Errorf(errorAt(idx, "'_' must separate successive digits"))
}

func convPPInt(tok *Token) bool {
	var base int = 10

	var sVal string
	var err error
	var startIdx = tok.Loc

	r := []rune(tok.Str)
	idx := 0

	if idx+2 <= len(r) &&
		(startsWith(string(r[idx:idx+2]), "0x") ||
			startsWith(string(r[idx:idx+2]), "0X")) {
		base = 16
		idx += 2
		sVal, err = readHexDigit(tok, &idx)
		if err != nil {
			panic(errorTok(tok, "invalid integer literal: %v", err))
		}
	} else if idx+2 <= len(r) &&
		(startsWith(string(r[idx:idx+2]), "0b") ||
			startsWith(string(r[idx:idx+2]), "0B")) {
		base = 2
		idx += 2
	} else if idx+2 <= len(r) &&
		(startsWith(string(r[idx:idx+2]), "0o") ||
			startsWith(string(r[idx:idx+2]), "0O") ||
			startsWith(string(r[idx:idx+2]), "0_")) {
		base = 8
		idx += 2
	} else if idx+1 < len(r) &&
		startsWith(string(r[idx:idx+1]), "0") &&
		isDigit(r[idx+1]) {
		base = 8
		idx += 1
	}

	for ; idx < len(r) &&
		(isDigit(r[idx]) || r[idx] == '_'); idx++ {

		if idx > 0 && r[idx-1] == '_' && r[idx] == '_' {
			panic(errorTok(tok, "invalid integer literal: %v", errMustSeparateSuccessiveDigits(startIdx)))
		}

		if isDigit(r[idx]) {
			sVal += string(r[idx])
		}
	}

	if idx < len(r) && idx > 0 && r[idx-1] == '_' {
		panic(errorTok(tok, "invalid integer literal: %v", errMustSeparateSuccessiveDigits(startIdx)))
	}

	if idx != tok.Len-1 {
		// fmt.Println("ここ")
		return false
	}

	// cur = newToken(TK_NUM, cur, string(curFile.Contents[startIdx:curIdx]), curIdx-startIdx+1)
	tok.Kind = TK_NUM
	var v int64
	if sVal != "" {
		v = parseInt(sVal, base, tok.Loc)
	}

	tok.Val = v

	tok.Ty = ty_long
	if v <= 2147483647 {
		tok.Ty = ty_int
	}

	return true
}

// The definition of the numeric literal at the preprocessing stage
// is more relaxed than the definition of that at the later stages.
// In order to handle that, a numeric literal is tokenized as a
// "pp-number" token first and then converted to a regular number
// token after preprocessing.
//
// This function converts a pp-number token to a regular number token.
func convPPNum(tok *Token) {
	curFile = tok.File
	// Try to parse as an integer constant.
	if convPPInt(tok) {
		return
	}

	idx := 0
	r := []rune(tok.Str)

	if idx+2 <= len(r) &&
		(startsWith(string(r[idx:idx+2]), "0x") ||
			startsWith(string(r[idx:idx+2]), "0X")) {
		idx += 2
	} else if idx+2 <= len(r) &&
		(startsWith(string(r[idx:idx+2]), "0b") ||
			startsWith(string(r[idx:idx+2]), "0B")) {
		idx += 2
	} else if idx+2 <= len(r) &&
		(startsWith(string(r[idx:idx+2]), "0o") ||
			startsWith(string(r[idx:idx+2]), "0O") ||
			startsWith(string(r[idx:idx+2]), "0_")) {
		idx += 2
	} else if idx+1 < len(r) &&
		startsWith(string(r[idx:idx+1]), "0") &&
		isDigit(r[idx+1]) {
		idx++
	} else if isDigit(r[idx]) || contains("eEfFpP.", r[idx]) {
		idx++
	} else {
		panic(errorTok(tok, "invalid numeric constant"))
	}

	for isDigit(r[idx]) ||
		contains("eEfFpP_.", r[idx]) ||
		(contains("EePp", r[idx-1]) &&
			contains("+-", r[idx])) {

		if idx+1 < len(r) &&
			((r[idx] == '_' && !isDigit(r[idx+1])) ||
				(r[idx+1] == '_' && !isDigit(r[idx]))) {
			panic(errorTok(tok, "invalid numeric constant: %v", errMustSeparateSuccessiveDigits(tok.Loc)))
		}

		// sVal += string(r[idx])
		idx++
		if idx >= len(r) {
			break
		}
	}

	if idx != tok.Len-1 {
		fmt.Println("idx:", idx)
		fmt.Printf("tok: %#v\n\n", tok)
		panic(errorTok(tok, "invalid numeric constant"))
	}

	fval := parseFloat(tok.Str, tok.Loc)

	ty := ty_double

	tok.Kind = TK_NUM
	tok.FVal = fval
	tok.Ty = ty
}

func convPPTok(tok *Token) {
	for t := tok; t.Kind != TK_EOF; t = t.Next {
		if isKw(t) {
			t.Kind = TK_RESERVED
		} else if t.Kind == TK_PP_NUM {
			convPPNum(t)
		}
	}
}

func readHexDigit(tok *Token, idx *int) (string, error) {
	var sVal string
	var errIdx = tok.Loc

	r := []rune(tok.Str)
	for ; *idx < len(r) && (isxdigit(r[*idx]) ||
		r[*idx] == '_'); *idx++ {

		if *idx > 0 && r[*idx-1] == '_' && r[*idx] == '_' {
			return "", errMustSeparateSuccessiveDigits(errIdx)
		}

		if isxdigit(r[*idx]) {
			sVal += string(r[*idx])
		}
	}

	if *idx > 0 && r[*idx-1] == '_' {
		return "", errMustSeparateSuccessiveDigits(errIdx)
	}

	return sVal, nil
}

func isxdigit(p rune) bool {
	return ('0' <= p && p <= '9') ||
		('A' <= p && p <= 'F') ||
		('a' <= p && p <= 'f')
}

func isxdigitb(p byte) bool {
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

		if curFile.Contents[idx] == 0 { // curFile.Contents[idx] == '\n' ||
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
			cur.Kind == TK_PP_NUM ||
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
	hasSpace = false

	for curIdx < len(curFile.Contents) {
		if curFile.Contents[curIdx] == 0 {
			break
		}

		// skip space(s)
		if isSpace(curFile.Contents[curIdx]) {
			curIdx++
			hasSpace = true
			continue
		}

		// new line
		if curFile.Contents[curIdx] == '\n' {
			cur = newToken(TK_NL, cur, "", 0)
			curIdx++
			atBol = true
			hasSpace = false
			continue
		}

		// skip line comment
		if startsWith(string(curFile.Contents[curIdx:]), "//") {
			curIdx += 2
			for ; curIdx < len(curFile.Contents) && curFile.Contents[curIdx] != '\n'; curIdx++ {
				// skip to the end of line.
			}
			cur = newToken(TK_COMM, cur, "<line comment>", 0)
			hasSpace = true
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
			hasSpace = true
			continue
		}

		// blank identifier
		if contains("_", curFile.Contents[curIdx]) && !isIdent2(curFile.Contents[curIdx+1]) {
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
			startIdx := curIdx
			if curIdx+1 < len(curFile.Contents) &&
				(startsWith(string(curFile.Contents[curIdx:curIdx+2]), "0X") ||
					startsWith(string(curFile.Contents[curIdx:curIdx+2]), "0x")) {
				for {
					if curIdx+1 < len(curFile.Contents) &&
						contains("pP", curFile.Contents[curIdx]) &&
						contains("+-", curFile.Contents[curIdx+1]) {
						curIdx += 2
					} else if curIdx < len(curFile.Contents) &&
						(isIdent2(curFile.Contents[curIdx]) ||
							curFile.Contents[curIdx] == '.') {
						curIdx++
					} else {
						break
					}
				}
			} else {
				for {
					if curIdx+1 < len(curFile.Contents) &&
						contains("eEpP", curFile.Contents[curIdx]) &&
						contains("+-", curFile.Contents[curIdx+1]) {
						curIdx += 2
					} else if curIdx < len(curFile.Contents) &&
						(isIdent2(curFile.Contents[curIdx]) ||
							curFile.Contents[curIdx] == '.') {
						curIdx++
					} else {
						break
					}
				}
			}
			cur = newToken2(TK_PP_NUM, cur, string(curFile.Contents[startIdx:curIdx]),
				curIdx-startIdx+1, startIdx)
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
		if contains("+-()*/<>=;{},&[].!|^:?%#`", curFile.Contents[curIdx]) {
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
			return nil, fmt.Errorf("readFile: os.Open: %v", err)
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

func readUniversalChar(p []rune, len int) rune {
	c := 0

	for i := 0; i < len; i++ {
		if !isxdigit(p[i]) {
			return 0
		}
		c = (c << 4) | fromHex(int(p[i]))
	}
	return rune(c)
}

// Replace \u or \U escape sequances with corresponding UTF-8 bytes.
func convUniversalChars(p *[]rune) {
	i := 0
	for i < len(*p) {
		if i+2 <= len(*p) && startsWith(string((*p)[i:i+2]), "\\u") {
			c := readUniversalChar((*p)[i+2:], 4)
			if c != 0 {
				(*p)[i] = c
				*p = append((*p)[:i+1], (*p)[i+6:]...)
				i++
			} else {
				i++
			}
		} else if i+2 <= len(*p) && startsWith(string((*p)[i:i+2]), "\\U") {
			c := readUniversalChar((*p)[i+2:], 8)
			if c != 0 {
				(*p)[i] = c
				*p = append((*p)[:i+1], (*p)[i+10:]...)
				i++
			} else {
				i++
			}
		} else {
			i++
		}
	}
}

// For tokenizeFile function
var fileno int

func tokenizeFile(path string) (*Token, error) {

	p, err := readFile(path)
	if err != nil {
		return nil, err
	}

	convUniversalChars(&p)

	file := newFile(path, fileno+1, p)

	// Save the filename for assembler .file directive.
	inputFiles = append(inputFiles, file)
	fileno++

	return tokenize(file)
}
