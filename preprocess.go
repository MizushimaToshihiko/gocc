package main

import (
	"fmt"
	"path/filepath"
)

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

func copyTok(tok *Token) *Token {
	var t = &Token{}
	*t = *tok
	t.Next = nil
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

		tok = tok.Next

		if equal(tok, "include") {
			tok = tok.Next

			if tok.Kind != TK_STR {
				panic("\n" + errorTok(tok, "expected a filename"))
			}

			path := fmt.Sprintf("%s/%s", filepath.Dir(tok.File.Name), tok.Str)
			tok2, err := tokenizeFile(path)
			if err != nil || tok2 == nil {
				panic("\n" + errorTok(tok, err.Error()))
			}
			tok = appendTok(tok2, tok.Next)
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

func preprocess(tok *Token) *Token {
	tok = preprocess2(tok)
	convKw(tok)
	delExtraSemicolon(tok)
	return tok
}
