package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func compile(arg string, w io.Writer) error {
	// tokenize and parse
	curIdx = 0 // for test
	userInput = arg
	token = tokenize()
	node := expr()

	printTokens()
	// output the former 3 lines of the assembly
	fmt.Fprintln(w, ".intel_syntax noprefix\n.globl main\nmain:")

	// make the asm code, down on the AST
	if err := gen(node, w); err != nil {
		return err
	}

	// the value of the expression should remain on the top of 'stack',
	// so load this value into rax.
	fmt.Fprintln(w, "	pop rax")
	fmt.Fprintln(w, "	ret")
	return nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("invalid number of arguments")
	}

	if err := compile(os.Args[1], os.Stdout); err != nil {
		log.Fatal(err)
	}
}
