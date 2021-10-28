package main

import (
	"io"
	"log"
	"os"
)

var filename string

func readFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	if len(b) == 0 || b[len(b)-1] != '\n' {
		b = append(b, '\n')
	}
	// b = append(b, 0)
	return string(b), nil
}

func alignTo(n, align int) int {
	return (n + align - 1) & ^(align - 1)
}

func compile(arg string, w io.Writer) error {
	// tokenize and parse
	curIdx = 0 // for test
	var err error
	userInput, err = readFile(arg)
	if err != nil {
		return err
	}
	filename = arg

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
		fn.StackSz = alignTo(offset, 8)
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
