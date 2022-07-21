package main

func preprocess(tok *Token) *Token {
	convKw(tok)
	return tok
}
