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

	// the parsed result is in 'code'
	program()

	// // walk in-order
	// for _, n := range code {
	// 	walkInOrder(n)
	// }
	// // walk pre order
	// for _, n := range code {
	// 	walkPreOrder(n)
	// }

	return codeGen(w)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("invalid number of arguments")
	}

	if err := compile(os.Args[1], os.Stdout); err != nil {
		log.Fatal(err)
	}
}
