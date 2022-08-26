// This file implements the preprocessor.
//
// The preprocessor tekes a list of tokens as an input and returns a
// new list of tokens as an output.
//
// The preprocessing language is designed in such a way that that's
// guaranteed to stop even if there is a recursive macro.
// Informally speaking, a macro is applied only once for each token.
// That is, if a macro token T appears in a result of direct or
// indirect macro expansion of T, T won't be expanded any further.
// For example, if T is defined as U, and U is defined as T, then
// token T is expanded to U and then to T and the macro expansion
// stops at that point.
//
// To archive the above behavior, we attach for each token a set of
// macro names from which the token is expanded. The set is called
// "hideset". Hideset is initially empty, and every time we expanded a
// macro, the macro name is added to the resulting tokens' hidesets.
//
// The above macro expansion algorithm is explained in this document
// written by Dave Prossor, which is used as a basis for the
// standard's wording:
// https://github.com/rui314/chibicc/wiki/cpp.algo.pdf
//
package main

import (
	"fmt"
	"path/filepath"
)

type MacroParam struct {
	Next *MacroParam
	Name string
}

type MacroArg struct {
	Next *MacroArg
	Name string
	Tok  *Token
}

type MacroHandlerFn func(*Token) *Token

type Macro struct {
	Next       *Macro
	Name       string
	IsObjlike  bool // Object-like or function-like
	Params     *MacroParam
	IsVariadic bool
	Body       *Token
	Deleted    bool
	Handler    MacroHandlerFn
}

type Ctx int

const (
	IN_THEN Ctx = iota
	IN_ELIF
	IN_ELSE
)

type CondIncl struct {
	Next     *CondIncl
	Ctx      Ctx
	Tok      *Token
	Included bool
}

type Hideset struct {
	Next *Hideset
	Name string
}

var macros *Macro
var condIncl *CondIncl

func delSemicolonTok(tok *Token) *Token {
	start := tok
	for t := tok; t.Next != nil; t = t.Next {
		if equal(t.Next, ";") {
			t.Next = t.Next.Next
		}
	}
	return start
}

func isHash(tok *Token) bool {
	return tok.AtBol && equal(tok, "#")
}

func skipLine(tok *Token) *Token {
	// Skip the ";" token.
	consume(&tok, tok, ";")

	if tok.AtBol {
		return tok
	}
	if !equal(tok, ";") {
		warnTok(tok, "extra token")
	}
	for tok.AtBol {
		tok = tok.Next
	}
	return tok
}

func copyTok(tok *Token) *Token {
	var t = &Token{}
	*t = *tok
	t.Next = nil
	return t
}

func newEof(tok *Token) *Token {
	t := copyTok(tok)
	t.Kind = TK_EOF
	t.Len = 0
	return t
}

func newHideset(name string) *Hideset {
	return &Hideset{Name: name}
}

func hidesetUnion(hs1, hs2 *Hideset) *Hideset {
	head := &Hideset{}
	cur := head

	for ; hs1 != nil; hs1 = hs1.Next {
		cur.Next = newHideset(hs1.Name)
		cur = cur.Next
	}
	cur.Next = hs2
	return head.Next
}

func hidesetContains(hs *Hideset, s string) bool {
	for ; hs != nil; hs = hs.Next {
		if hs.Name == s {
			return true
		}
	}
	return false
}

func hidesetIntersection(hs1, hs2 *Hideset) *Hideset {
	head := &Hideset{}
	cur := head

	for ; hs1 != nil; hs1 = hs1.Next {
		if hidesetContains(hs2, hs1.Name) {
			cur.Next = newHideset(hs1.Name)
			cur = cur.Next
		}
	}
	return head.Next
}

func addHideset(tok *Token, hs *Hideset) *Token {
	head := &Token{}
	cur := head

	for ; tok != nil; tok = tok.Next {
		t := copyTok(tok)
		t.Hideset = hidesetUnion(t.Hideset, hs)
		cur.Next = t
		cur = cur.Next
	}
	return head.Next
}

// Append tok2 to the end of tok1.
func appendTok(tok1, tok2 *Token) *Token {
	if tok1.Kind == TK_EOF {
		return tok2
	}

	head := &Token{}
	cur := head

	for ; tok1.Kind != TK_EOF; tok1 = tok1.Next {
		cur.Next = copyTok(tok1)
		cur = cur.Next
	}
	cur.Next = tok2
	return head.Next
}

func skipCondIncl2(tok *Token) *Token {
	for tok.Kind != TK_EOF {
		if isHash(tok) &&
			(equal(tok.Next, "if") || equal(tok.Next, "ifdef") ||
				equal(tok.Next, "ifndef")) {
			tok = skipCondIncl2(tok.Next.Next)
			continue
		}
		if isHash(tok) && equal(tok.Next, "endif") {
			return tok.Next.Next
		}
		tok = tok.Next
	}
	return tok
}

// Skip until next `#else`, `#elif` or `#endif`.
// Nested `#if` and `#endif` are skipped.
func skipCondIncl(tok *Token) *Token {
	for tok.Kind != TK_EOF {
		if isHash(tok) &&
			(equal(tok.Next, "if") || equal(tok.Next, "ifdef") ||
				equal(tok.Next, "ifndef")) {
			tok = skipCondIncl2(tok.Next.Next)
			continue
		}

		if isHash(tok) &&
			(equal(tok.Next, "elif") || equal(tok.Next, "else") ||
				equal(tok.Next, "endif")) {
			break
		}
		tok = tok.Next
	}
	return tok
}

// quoteStr adds double-quotes to the front and back of a given string and
// adds 0 at the end, and returns it as a byte array.
func quoteStr(str []byte) []byte {
	var ret = make([]byte, 0)
	ret = append(ret, '"')
	for _, r := range str {
		if r == '\\' || r == '"' {
			ret = append(ret, '\\')
		}
		if r == byte(0) {
			break
		}
		ret = append(ret, r)
	}
	ret = append(ret, '"')
	ret = append(ret, byte(0))

	return ret
}

func newStrTok(str []byte, tmpl *Token) *Token {
	buf := quoteStr(str)
	t, err := tokenize(newFile(tmpl.File.Name, tmpl.File.FileNo, buf))
	if err != nil {
		panic(err)
	}
	return t
}

// Copy all tokens until the next newline, terminate them with
// an EOF token and then return them. This function is used to
// create a new list of tokens for `#if` arguments.
func copyLine(rest **Token, tok *Token) *Token {
	head := &Token{}
	cur := head

	for ; !equal(tok, ";") && !tok.AtBol; tok = tok.Next {
		// `!equal(tok, ";")` -> Stop before ";" appears.
		cur.Next = copyTok(tok)
		cur = cur.Next
	}

	cur.Next = newEof(tok)
	if equal(tok, ";") {
		cur.Next.Str = ""
	}
	*rest = tok
	return head.Next
}

func newNumTok(val int, tmpl *Token) *Token {
	buf := fmt.Sprintf("%d", val)
	tok, err := tokenize(newFile(tmpl.File.Name, tmpl.File.FileNo, []byte(buf)))
	if err != nil {
		panic(err)
	}
	return tok
}

func readConstExpr(rest **Token, tok *Token) *Token {
	tok = copyLine(rest, tok)

	head := &Token{}
	cur := head

	for tok.Kind != TK_EOF {
		// "defined(foo)" or "defined foo" becomes "1" if macro "foo"
		// is defined. Otherwise "0".
		if equal(tok, "defined") {
			start := tok
			hasParen := consume(&tok, tok.Next, "(")

			if tok.Kind != TK_IDENT {
				panic("\n" + errorTok(start, "macro name must be an identifier"))
			}
			m := findMacro(tok)
			tok = tok.Next

			if hasParen {
				tok = skip(tok, ")")
			}

			val := 0
			if m != nil {
				val = 1
			}
			cur.Next = newNumTok(val, start)
			cur = cur.Next
			continue
		}
		cur.Next = tok
		cur = cur.Next
		tok = tok.Next
	}

	cur.Next = tok
	return head.Next
}

// Read and evaluate a contsant expression.
func evalConstExpr(rest **Token, tok *Token) int64 {
	start := tok
	expr := readConstExpr(rest, tok.Next)
	expr = preprocess2(expr)

	if expr.Kind == TK_EOF {
		panic("\n" + errorTok(start, "no expression"))
	}

	// [https://www.sigbus.info/n1570#6.10.1p4] THe standerd requires
	// we replace remaining non-macro identifiers with "0" before
	// evaluateing a constant expression. For example, `#if foo` is
	// equivalent tok `#if 0` if foo is not defined.
	for t := expr; t.Kind != TK_EOF; t = t.Next {
		if t.Kind == TK_IDENT {
			next := t.Next
			*t = *newNumTok(0, t)
			t.Next = next
		}
	}

	// Convert pp-numbers to regular numbers
	convPPTok(expr)

	var rest2 *Token
	val := constExpr(&rest2, expr)
	consume(&rest2, rest2, ";") // If rest2 is ";" token before this, rest2 should be `nil`.
	if rest2 != nil && rest2.Kind != TK_EOF {
		panic("\n" + errorTok(rest2, "extra token"))
	}
	return val
}

func pushCondIncl(tok *Token, included bool) *CondIncl {
	ci := &CondIncl{
		Next:     condIncl,
		Ctx:      IN_THEN,
		Tok:      tok,
		Included: included,
	}
	condIncl = ci
	return ci
}

func findMacro(tok *Token) *Macro {
	if tok.Kind != TK_IDENT {
		return nil
	}

	for m := macros; m != nil; m = m.Next {
		if m.Name == tok.Str {
			if m.Deleted {
				return nil
			} else {
				return m
			}
		}
	}
	return nil
}

func addMacro(name string, isObjlike bool, body *Token) *Macro {
	m := &Macro{
		Next:      macros,
		Name:      name,
		IsObjlike: isObjlike,
		Body:      body,
	}
	macros = m
	return m
}

func readMacroParams(rest **Token, tok *Token, isVariadic *bool) *MacroParam {
	head := &MacroParam{}
	cur := head

	for !equal(tok, ")") {
		if cur != head {
			tok = skip(tok, ",")
		}

		if equal(tok, "...") {
			*isVariadic = true
			*rest = skip(tok.Next, ")")
			return head.Next
		}

		if tok.Kind != TK_IDENT {
			panic("\n" + errorTok(tok, "expected an identifier"))
		}
		name := tok.Str
		m := &MacroParam{Name: name}
		cur.Next = m
		cur = cur.Next
		tok = tok.Next
	}
	*rest = tok.Next
	return head.Next
}

func readMacroDef(rest **Token, tok *Token) {
	if tok.Kind != TK_IDENT {
		panic("\n" + errorTok(tok, "macro name must be an identifier"))
	}
	name := tok.Str
	tok = tok.Next

	if !tok.HasSpace && equal(tok, "(") {
		// Function-like macro
		isVariadic := false
		params := readMacroParams(&tok, tok.Next, &isVariadic)

		m := addMacro(name, false, copyLine(rest, tok))
		m.Params = params
		m.IsVariadic = isVariadic
	} else {
		// Object-like macro
		addMacro(name, true, copyLine(rest, tok))
	}
}

func readMacroArg1(rest **Token, tok *Token, readRest bool) *MacroArg {
	head := &Token{}
	cur := head
	level := 0

	for {
		if level == 0 && equal(tok, ")") {
			break
		}
		if level == 0 && !readRest && equal(tok, ",") {
			break
		}

		if tok.Kind == TK_EOF {
			panic("\n" + errorTok(tok, "premature end of input"))
		}

		if equal(tok, "(") {
			level++
		} else if equal(tok, ")") {
			level--
		}

		cur.Next = copyTok(tok)
		cur = cur.Next
		tok = tok.Next
	}

	cur.Next = newEof(tok)

	arg := &MacroArg{Tok: head.Next}
	*rest = tok
	return arg
}

func readMacroArgs(rest **Token, tok *Token,
	params *MacroParam, isVariadic bool) *MacroArg {
	start := tok
	tok = tok.Next.Next

	head := &MacroArg{}
	cur := head

	pp := params
	for ; pp != nil; pp = pp.Next {
		if cur != head {
			tok = skip(tok, ",")
		}
		cur.Next = readMacroArg1(&tok, tok, false)
		cur = cur.Next
		cur.Name = pp.Name
	}

	if isVariadic {
		var arg *MacroArg
		if equal(tok, ")") {
			arg = &MacroArg{Tok: newEof(tok)}
		} else {
			if pp != params {
				tok = skip(tok, ",")
			}
			arg = readMacroArg1(&tok, tok, true)
		}
		arg.Name = "__VA_ARGS__"
		cur.Next = arg
		cur = cur.Next
	} else if pp != nil {
		panic("\n" + errorTok(start, "too many arguments"))
	}

	skip(tok, ")")
	*rest = tok
	return head.Next
}

func findArg(args *MacroArg, tok *Token) *MacroArg {
	for ap := args; ap != nil; ap = ap.Next {
		if ap.Name == tok.Str {
			return ap
		}
	}
	return nil
}

// Concatenates all tokens in `tok` and returns a new string.
func joinTok(tok, end *Token) []byte {
	// Compute the length of the resulting token.
	len := 1
	for t := tok; t != end && t.Kind != TK_EOF; t = t.Next {
		if t != tok && t.HasSpace {
			len++
		}
		len += t.Len
	}

	var buf []byte

	// Copy token texts.
	for t := tok; t != end && t.Kind != TK_EOF; t = t.Next {
		if t != tok && t.HasSpace {
			buf = append(buf, ' ')
		}

		str := t.Str
		if t.Kind == TK_STR { // add double-quote
			buf = append(buf, '"')
			buf = append(buf, []byte(str)...)
			buf = append(buf, '"')
		} else {
			buf = append(buf, []byte(str)...)
		}
	}

	buf = append(buf, byte(0))
	return buf
}

// Concatenates all tokens in `arg` and returns a new string token.
// This function is used for the stringizing operator (#).
func stringize(hash, arg *Token) *Token {
	// Create a new string token. We need to set some value to its
	// source location for error reporting function, so we use a macro
	// name token as a template.
	s := joinTok(arg, nil)
	return newStrTok(s, hash)
}

// Concatenate two tokens to create a new token.
func paste(lhs, rhs *Token) *Token {
	// Paste the two tokens.
	buf := append([]byte(lhs.Str), []byte(rhs.Str)...)

	// Tokenize the resulting string.
	tok, err := tokenize(newFile(lhs.File.Name, lhs.File.FileNo, buf))
	if err != nil {
		panic(err)
	}
	delSemicolonTok(tok)
	if tok.Next.Kind != TK_EOF {
		panic("\n" + errorTok(lhs, "pasting forms '%s', an invalid token", string(buf)))
	}
	return tok
}

func subst(tok *Token, args *MacroArg) *Token {
	head := &Token{}
	cur := head

	for tok.Kind != TK_EOF {
		// "#" followed by a parameter is replaced with stringized actuals.
		if equal(tok, "#") {
			arg := findArg(args, tok.Next)
			if arg == nil {
				panic("\n" + errorTok(tok.Next, "'#' is not followed by a macro parameter"))
			}
			cur.Next = stringize(tok, arg.Tok)
			cur = cur.Next
			tok = tok.Next.Next
			continue
		}

		if equal(tok, "##") {
			if cur == head {
				panic("\n" + errorTok(tok, "'##' cannot appear at start of macro expansion"))
			}

			if tok.Next.Kind == TK_EOF {
				panic("\n" + errorTok(tok, "'##' cannot appear at end of macro expansion"))
			}

			arg := findArg(args, tok.Next)
			if arg != nil {
				if arg.Tok.Kind != TK_EOF {
					*cur = *paste(cur, arg.Tok)
					for t := arg.Tok.Next; t.Kind != TK_EOF; t = t.Next {
						cur.Next = copyTok(t)
						cur = cur.Next
					}
				}
				tok = tok.Next.Next
				continue
			}

			*cur = *paste(cur, tok.Next)
			tok = tok.Next.Next
			continue
		}

		arg := findArg(args, tok)

		if arg != nil && equal(tok.Next, "##") {
			rhs := tok.Next.Next

			if arg.Tok.Kind == TK_EOF {
				arg2 := findArg(args, rhs)
				if arg2 != nil {
					for t := arg2.Tok; t.Kind != TK_EOF; t = t.Next {
						cur.Next = copyTok(t)
						cur = cur.Next
					}
				} else {
					cur.Next = copyTok(rhs)
					cur = cur.Next
				}
				tok = rhs.Next
				continue
			}

			for t := arg.Tok; t.Kind != TK_EOF; t = t.Next {
				cur.Next = copyTok(t)
				cur = cur.Next
			}
			tok = tok.Next
			continue
		}

		// Handle a macro token. Macro arguments are completely macro-expanded
		// before they are substituted into a macro body.
		if arg != nil {
			t := preprocess2(arg.Tok)
			t.AtBol = tok.AtBol
			t.HasSpace = tok.HasSpace
			for ; t.Kind != TK_EOF; t = t.Next {
				cur.Next = copyTok(t)
				cur = cur.Next
			}

			tok = tok.Next
			continue
		}

		// Handle a non-macro token.
		cur.Next = copyTok(tok)
		cur = cur.Next
		tok = tok.Next
		continue
	}

	cur.Next = tok
	return head.Next
}

// If tok is a macro, expand it and return true.
// Otherwise, do nothing and return false.
func expandMacro(rest **Token, tok *Token) bool {
	if hidesetContains(tok.Hideset, tok.Str) {
		return false
	}

	m := findMacro(tok)
	if m == nil {
		return false
	}

	if m.Handler != nil {
		*rest = m.Handler(tok)
		(*rest).Next = tok.Next
		return true
	}

	// Object-like macro application
	if m.IsObjlike {
		hs := hidesetUnion(tok.Hideset, newHideset(m.Name))
		body := addHideset(m.Body, hs)
		for t := body; t.Kind != TK_EOF; t = t.Next {
			t.Origin = tok
		}
		*rest = appendTok(body, tok.Next)
		(*rest).AtBol = tok.AtBol
		(*rest).HasSpace = tok.HasSpace
		return true
	}

	// If a funclike macro token is not followed by an argument list,
	// treat it as a normal identifier.
	if !equal(tok.Next, "(") {
		return false
	}

	// Function-like macro application
	macroTok := tok
	args := readMacroArgs(&tok, tok, m.Params, m.IsVariadic)
	rparen := tok

	// Tokens that consist a func-like maro invocation may have different
	// hidesets, and if that's the case, it's not clear what the hideset
	// for the new tokens should be. We take the intersection of the
	// macro token and the closing parenthesis and use it as a new hideset
	// as explained in the Dave Prossor's algorithm.
	hs := hidesetIntersection(macroTok.Hideset, rparen.Hideset)
	hs = hidesetUnion(hs, newHideset(m.Name))

	body := subst(m.Body, args)
	body = addHideset(body, hs)
	for t := body; t.Kind != TK_EOF; t = t.Next {
		t.Origin = macroTok
	}
	*rest = appendTok(body, tok.Next)
	(*rest).AtBol = macroTok.AtBol
	(*rest).HasSpace = macroTok.HasSpace
	return true
}

func searchIncludePaths(filename string) string {
	if filename[0] == '/' {
		return filename
	}

	// Search a file from the include paths.
	for i := 0; i < len(includePaths); i++ {
		path := fmt.Sprintf("%s/%s", includePaths[i], filename)
		if exists(path) {
			return path
		}
	}
	return ""
}

// Read an #include argument.
func readIncludeFilename(rest **Token, tok *Token, isDquote *bool) []byte {
	// Pattern 1: #include "foo.h"
	if tok.Kind == TK_STR {
		// A double-quoted filename for #include is a special kind of
		// token, and we don't want to interpret any escape sequences in it.
		// For expample, "\f" in "C:\foo" is not a formfeed character but
		// just two non-control characters, backslash and f.
		// So we don't want to use token.Str.
		// => The above cannot be applied to this compiler.
		*isDquote = true
		*rest = skipLine(tok.Next)
		return tok.File.Contents[tok.Loc+1 : tok.Loc+tok.Len-1]
	}

	// Pattern 2: #include <foo.h>
	if equal(tok, "<") {
		// Reconstruct a filename from a sequence of tokens between
		// "<" and ">".
		start := tok

		// Find closing ">".
		for ; !equal(tok, ">"); tok = tok.Next {
			if tok.AtBol || tok.Kind == TK_EOF {
				panic("\n" + errorTok(tok, "expected '>'"))
			}
		}

		*isDquote = false
		*rest = skipLine(tok.Next)
		return joinTok(start.Next, tok)
	}

	// Pattern 3: #include FOO
	// In this case FOO must be macro-expanded to either
	// a single string token or a sequence of "<" ... ">".
	if tok.Kind == TK_IDENT {
		tok2 := preprocess2(copyLine(rest, tok))
		return readIncludeFilename(&tok2, tok2, isDquote)
	}

	panic("\n" + errorTok(tok, "expected a filename"))
}

func includeFile(tok *Token, path string, filenameTok *Token) *Token {
	tok2, err := tokenizeFile(path)
	if err != nil || tok2 == nil {
		panic("\n" + errorTok(filenameTok, "%s: cannot open file: %s", path, err))
	}
	return appendTok(tok2, tok)
}

// Visit all tokens int `tok` while evaluating preprocessing
// macros and directives.
func preprocess2(tok *Token) *Token {
	head := &Token{}
	cur := head

	for tok.Kind != TK_EOF {
		// If it is a macro, expand it.
		if expandMacro(&tok, tok) {
			continue
		}

		// Pass through if it is not a "#".
		if !isHash(tok) {
			cur.Next = tok
			cur = cur.Next
			tok = tok.Next
			continue
		}

		start := tok
		tok = tok.Next

		if equal(tok, "include") {
			var isDquote bool
			filename := readIncludeFilename(&tok, tok.Next, &isDquote)
			// Delete '\0' in the end of 'filename'.
			if filename[len(filename)-1] == 0 {
				filename = filename[:len(filename)-1]
			}

			var path string
			if filename[0] != '/' && isDquote {
				path = fmt.Sprintf("%s/%s", filepath.Dir(tok.File.Name), string(filename))
				if exists(path) {
					tok = includeFile(tok, path, start.Next.Next)
				}
				continue
			}

			// Search a file from the include paths.
			path = searchIncludePaths(string(filename))
			if path != "" {
				tok = includeFile(tok, path, start.Next.Next)
				continue
			}

			tok = includeFile(tok, string(filename), start.Next.Next)
			continue
		}

		if equal(tok, "define") {
			readMacroDef(&tok, tok.Next)
			continue
		}

		if equal(tok, "undef") {
			tok = tok.Next
			if tok.Kind != TK_IDENT {
				panic("\n" + errorTok(tok, "macro name must be an identifier"))
			}
			name := tok.Str
			tok = skipLine(tok.Next)

			m := addMacro(name, true, nil)
			m.Deleted = true
			continue
		}

		if equal(tok, "if") {
			val := evalConstExpr(&tok, tok)
			pushCondIncl(start, val != 0)
			if val == 0 {
				tok = skipCondIncl(tok)
			}
			continue
		}

		if equal(tok, "ifdef") {
			defined := findMacro(tok.Next)
			pushCondIncl(tok, defined != nil)
			tok = skipLine(tok.Next.Next)
			if defined == nil {
				tok = skipCondIncl(tok)
			}
			continue
		}

		if equal(tok, "ifndef") {
			defined := findMacro(tok.Next)
			pushCondIncl(tok, defined == nil)
			tok = skipLine(tok.Next.Next)
			if defined != nil {
				tok = skipCondIncl(tok)
			}
			continue
		}

		if equal(tok, "elif") {
			if condIncl == nil || condIncl.Ctx == IN_ELSE {
				panic("\n" + errorTok(start, "stray #elif"))
			}
			condIncl.Ctx = IN_ELIF

			if !condIncl.Included && evalConstExpr(&tok, tok) != 0 {
				condIncl.Included = true
			} else {
				tok = skipCondIncl(tok)
			}
			continue
		}

		if equal(tok, "else") {
			if condIncl == nil || condIncl.Ctx == IN_ELSE {
				panic("\n" + errorTok(start, "stray #else"))
			}
			condIncl.Ctx = IN_ELSE
			tok = skipLine(tok.Next)

			if condIncl.Included {
				tok = skipCondIncl(tok)
			}
			continue
		}

		if equal(tok, "endif") {
			if condIncl == nil {
				panic("\n" + errorTok(start, "stray #endif"))
			}
			condIncl = condIncl.Next
			tok = skipLine(tok.Next)
			continue
		}

		if equal(tok, "error") {
			if tok.Next.AtBol {
				panic("\n" + errorTok(tok, "error"))
			}

			start := tok
			errTok := copyLine(&tok, tok)
			errmsg := joinTok(errTok, nil)
			panic("\n" + errorTok(start, string(errmsg)))
		}

		// `#`-only line is legal. It's called a null directive
		if tok.AtBol {
			continue
		}

		panic("\n" + errorTok(tok, "invalid preprocessor directive"))
	}

	cur.Next = tok
	return head.Next
}

func defineMacro(name, buf string) {
	tok, err := tokenize(newFile("<built-in>", 1, []byte(buf)))
	if err != nil {
		panic(err)
	}
	delSemicolonTok(tok)
	addMacro(name, true, tok)
}

func addBuildIn(name string, fn MacroHandlerFn) *Macro {
	m := addMacro(name, true, nil)
	m.Handler = fn
	return m
}

func fileMacro(tmpl *Token) *Token {
	for tmpl.Origin != nil {
		tmpl = tmpl.Origin
	}
	return newStrTok([]byte(tmpl.File.Name), tmpl)
}

func lineMacro(tmpl *Token) *Token {
	for tmpl.Origin != nil {
		tmpl = tmpl.Origin
	}
	return newNumTok(tmpl.LineNo, tmpl)
}

func initMacros() {
	// Define predefined macros
	defineMacro("_LP64", "1")
	defineMacro("__C99_MACRO_WITH_VA_ARGS", "1")
	defineMacro("__ELF__", "1")
	defineMacro("__LP64__", "1")
	defineMacro("__SIZEOF_DOUBLE__", "8")
	defineMacro("__SIZEOF_FLOAT__", "4")
	defineMacro("__SIZEOF_INT__", "4")
	defineMacro("__SIZEOF_LONG_DOUBLE__", "8")
	defineMacro("__SIZEOF_LONG_LONG__", "8")
	defineMacro("__SIZEOF_POINTER__", "8")
	defineMacro("__SIZEOF_PTRDIFF_T__", "8")
	defineMacro("__SIZEOF_SHORT__", "2")
	defineMacro("__SIZEOF_SIZE_T__", "8")
	defineMacro("__SIZE_TYPE__", "unsigned long")
	defineMacro("__STDC_HOSTED__", "1")
	defineMacro("__STDC_NO_ATMICS__", "1")
	defineMacro("__STDC_NO_COMPLEX__", "1")
	defineMacro("__STDC_NO_THREADS__", "1")
	defineMacro("__STDC_NO_VLA__", "1")
	defineMacro("__STDC_VERSION__", "201112L")
	defineMacro("__STDC__", "1")
	defineMacro("__USER_LABEL_PREFIX__", "")
	defineMacro("__alignof__", "_Alignof")
	defineMacro("__amd64", "1")
	defineMacro("__amd64__", "1")
	defineMacro("__const__", "const")
	defineMacro("__gnu_linux__", "1")
	defineMacro("__inline__", "inline")
	defineMacro("__linux", "1")
	defineMacro("__linux__", "1")
	defineMacro("__signed__", "signed")
	defineMacro("__typeof__", "typeof")
	defineMacro("__unix", "1")
	defineMacro("__unix__", "1")
	defineMacro("__volatile__", "volatile")
	defineMacro("__x86_64", "1")
	defineMacro("__x86_64__", "1")
	defineMacro("linux", "1")
	defineMacro("unix", "1")

	addBuildIn("__FILE__", fileMacro)
	addBuildIn("__LINE__", lineMacro)
}

// Entry point function of the preprocessor.
func preprocess(tok *Token) *Token {
	initMacros()
	tok = preprocess2(tok)
	if condIncl != nil {
		panic("\n" + errorTok(condIncl.Tok, "unterminated conditional derective"))
	}
	convPPTok(tok)
	return tok
}
