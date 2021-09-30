package main

// set TokenKind with enum
type TokenKind int

const (
	TK_RESERVED TokenKind = iota
	TK_NUM
	TK_EOF
)

type Token struct {
	kind TokenKind
	next *Token
	val  int
	str  string
}

// current token
var token *Token

// inputted program
var user_input *string

// for error report
// it's arguments are same as printf
func error_at(loc string, fmt ...interface{})

func main() {

}
