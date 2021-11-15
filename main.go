package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

var filename string

func readFile(path string) ([]rune, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	br := bufio.NewReader(f)

	ret := make([]rune, 0, 1064)
	for {
		ru, sz, err := br.ReadRune()
		if sz == 0 || err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		ret = append(ret, ru)
	}
	ret = append(ret, 0)
	return ret, nil
}

func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func compile(arg string, w io.Writer) error {
	// tokenize and parse
	curIdx = 0 // for test

	// if !exists(arg) {
	// 	return fmt.Errorf("compile(): err: %s: %v", arg, os.ErrNotExist)
	// }

	var err error
	userInput = []rune(arg)
	if err != nil {
		return err
	}
	filename = arg

	token, err = tokenize()
	if err != nil {
		printTokens()
		return err
	}

	printTokens()
	return nil
	// prog := program()

	// Assign offsets to local variables
	// offset := 0
	// for v := prog.Locals; v != nil; v = v.Next {
	// 	offset += 8
	// 	v.Offset = offset
	// }
	// prog.StackSz = offset

	// for n := node; n != nil; n = n.Next {
	// 	walkInOrder(n)
	// }

	// return codegen(prog, w) // make the asm code, down on the AST
	/*
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
				offset = alignTo(offset, vl.Var.Ty.Align)
				offset += sizeOf(vl.Var.Ty, vl.Var.Tok)
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

		return codeGen(w, prog) */
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("invalid number of arguments")
	}

	if err := compile(os.Args[1], os.Stdout); err != nil {
		log.Fatal(err)
	}
}
