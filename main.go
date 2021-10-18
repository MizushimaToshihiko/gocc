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

	// the parsed result is in 'prog'
	var prog *Function = program()
	// add 'Type' to AST
	addType(prog)

	// assign offsets to local variables.
	for fn := prog; fn != nil; fn = fn.Next {
		offset := 0
		for vl := fn.Locals; vl != nil; vl = vl.Next {
			offset += 8
			vl.Var.Offset = offset
		}
		fn.StackSz = offset
	}

	// // walk in-order
	// for _, n := range code {
	// 	walkInOrder(n)
	// }
	// // walk pre order
	// for _, n := range code {
	// 	walkPreOrder(n)
	// }

	return codeGen(w, prog)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("invalid number of arguments")
	}

	if err := compile(os.Args[1], os.Stdout); err != nil {
		log.Fatal(err)
	}
}
