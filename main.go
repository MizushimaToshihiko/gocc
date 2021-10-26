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

	var err error
	token, err = tokenize()
	if err != nil {
		return err
	}
	// printTokens()

	// the parsed result is in 'prog'
	prog := program()
	// add 'Type' to ASTs
	err = addType(prog)
	if err != nil {
		return err
	}

	// assign offsets to local variables.
	for fn := prog.Fns; fn != nil; fn = fn.Next {
		offset := 0
		for vl := fn.Locals; vl != nil; vl = vl.Next {
			offset += sizeOf(vl.Var.Ty)
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
