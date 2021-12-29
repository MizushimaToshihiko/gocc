package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var optOut *os.File
var inputPath string

func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func compile(arg string, w io.Writer) error {
	// tokenize and parse
	curIdx = 0 // for test

	if !exists(arg) {
		return fmt.Errorf("compile(): err: %s: %v", arg, os.ErrNotExist)
	}

	var err error
	token, err = tokenizeFile(arg)
	if err != nil {
		// printTokens()
		return err
	}

	// printTokens()
	prog := program()
	err = addType(prog)
	if err != nil {
		return err
	}

	// for n := node; n != nil; n = n.Next {
	// 	walkInOrder(n)
	// }
	// }

	fmt.Fprintf(w, ".file 1 \"%s\"\n", curFilename)
	return codegen(w, prog) // make the asm code, down on the AST
}

func usage(status int) {
	fmt.Fprintf(os.Stderr, "gocc [ -o <path> ] <file>\n")
	os.Exit(status)
}

func main() {
	// setting log
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	var outpath string
	flag.StringVar(&outpath, "o", "", "The output file name")
	var help bool
	flag.BoolVar(&help, "help", false, "Help")
	flag.Parse()

	if help {
		usage(0)
	}

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "no input files")
		usage(1)
	}

	inputPath = flag.Args()[0]

	var err error
	if outpath == "" {
		optOut = os.Stdout
	} else {
		optOut, err = os.Create(outpath)
		if err != nil {
			fmt.Println(inputPath)
			log.Fatal(err)
		}
	}

	if err := compile(inputPath, optOut); err != nil {
		log.Fatal(err)
	}
}
