package main

import (
	"fmt"
	"path/filepath"
)

type CondIncl struct {
	Next *CondIncl
	Tok  *Token
}

var condIncl *CondIncl

// delete extra semicolons
func delExtraSemicolon(tok *Token) {
	for t := tok; t.Kind != TK_EOF && t != nil; t = t.Next {
		if t.Str == ";" && t.Next.Str == ";" {
			t.Next = t.Next.Next
		}
	}
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
	warnTok(tok, "extra token")
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

// Append tok2 to the end of tok1.
func appendTok(tok1, tok2 *Token) *Token {
	if tok1 == nil || tok1.Kind == TK_EOF {
		return tok2
	}

	head := &Token{}
	cur := head

	for ; tok1 != nil && tok1.Kind != TK_EOF; tok1 = tok1.Next {
		cur.Next = copyTok(tok1)
		cur = cur.Next
	}
	cur.Next = tok2
	return head.Next
}

// Skip until next `#endif`
// Nested `#if` and `#endif` are skipped.
func skipCondIncl(tok *Token) *Token {
	for tok.Kind != TK_EOF {
		if isHash(tok) && equal(tok.Next, "if") {
			tok = skipCondIncl(tok.Next.Next)
			tok = tok.Next
			continue
		}
		if isHash(tok) && equal(tok.Next, "endif") {
			break
		}
		tok = tok.Next
	}
	return tok
}

// Copy all tokens until the next newline, terminate them with
// an EOF token and then return them. This function is used to
// create a new list of tokens for `#if` arguments.
func copyLine(rest **Token, tok *Token) *Token {
	head := &Token{}
	cur := head

	for ; !tok.AtBol; tok = tok.Next {
		cur.Next = copyTok(tok)
		cur = cur.Next
	}
	cur.Next = newEof(tok)
	*rest = tok
	return head.Next
}

// Read and evaluate a contsant expression.
func evalConstExpr(rest **Token, tok *Token) int64 {
	start := tok
	expr := copyLine(rest, tok.Next)

	if expr.Kind == TK_EOF {
		panic("\n" + errorTok(start, "no expression"))
	}

	var rest2 *Token
	val := constExpr(&rest2, expr)
	consume(&rest2, rest2, ";")
	if rest2.Kind != TK_EOF {
		panic("\n" + errorTok(rest2, "extra token"))
	}
	return val
}

func pushCondIncl(tok *Token) *CondIncl {
	ci := &CondIncl{Next: condIncl, Tok: tok}
	condIncl = ci
	return ci
}

// Visit all tokens int `tok` while evaluating preprocessing
// macros and directives.
func preprocess2(tok *Token) *Token {
	head := &Token{}
	cur := head

	for tok.Kind != TK_EOF {
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
			tok = tok.Next

			if tok.Kind != TK_STR {
				panic("\n" + errorTok(tok, "expected a filename"))
			}

			var path string
			if tok.Str[0] == '/' {
				path = tok.Str
			} else {
				path = fmt.Sprintf("%s/%s", filepath.Dir(tok.File.Name), tok.Str)
			}

			tok2, err := tokenizeFile(path)
			if err != nil || tok2 == nil {
				panic("\n" + errorTok(tok, err.Error()))
			}
			tok = skipLine(tok.Next)
			tok = appendTok(tok2, tok)
			continue
		}

		if equal(tok, "if") {
			val := evalConstExpr(&tok, tok)
			pushCondIncl(start)
			if val == 0 {
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

		// `#`-only line is legal. It's called a null directive
		if tok.AtBol {
			continue
		}

		panic("\n" + errorTok(tok, "invalid preprocessor directive"))
	}

	cur.Next = tok
	return head.Next
}

// Entry point function of the preprocessor.
func preprocess(tok *Token) *Token {
	tok = preprocess2(tok)
	if condIncl != nil {
		panic("\n" + errorTok(condIncl.Tok, "unterminated conditional derective"))
	}
	convKw(tok)
	delExtraSemicolon(tok)
	return tok
}
