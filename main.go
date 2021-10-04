package main

import (
	"io"
	"log"
	"os"
)

func compile(arg string, w io.Writer) error {
	// tokenize and parse
	curIdx = 0 // for test
	userInput = arg
	token = tokenize()
	// printTokens()

	node := expr()
	// // walk in-order
	// walkInOrder(node)
	// // walk pre order
	// walkPreOrder(node)

	return codeGen(w, node)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("invalid number of arguments")
	}

	if err := compile(os.Args[1], os.Stdout); err != nil {
		log.Fatal(err)
	}
}
