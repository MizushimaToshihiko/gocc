package main

func isHash(tok *Token) bool {
	return tok.AtBol && equal(tok, "#")
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
	return tok
}
