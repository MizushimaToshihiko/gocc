package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var optOut *os.File
var inputPaths []string

var isdeb bool // Is debug mode or not.

func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// baseName returns the filepath without a extension.
func baseName(filePath string) string {
	ext := filepath.Ext(filePath)
	return filePath[:len(filePath)-len(ext)]
}

func replaceExt(tmpl, ext string) string {
	return fmt.Sprintf("%s.%s", baseName(tmpl), ext)
}

func compile(prtok bool, arg string, w io.Writer) error {
	// tokenize and parse
	curIdx = 0 // for test

	if !exists(arg) {
		return fmt.Errorf("compile(): err: %s: %v", arg, os.ErrNotExist)
	}

	// tokenize
	var err error
	var tok *Token
	tok, err = tokenizeFile(arg)
	if err != nil {
		// printTokens(tok)
		return err
	}

	if prtok {
		printTokens(tok)
		return nil
	}

	if isdeb {
		printTokens(tok)
	}

	// parse the tokens, and make the AST nodes.
	prog := parse(tok)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, ".file 1 \"%s\"\n", curFilename)
	return codegen(w, prog) // make the asm code, down on the AST
}

func usage(status int) {
	fmt.Fprintf(os.Stderr, "usage: ./bin/gocc [ -o <path> ] <file>\n")
	os.Exit(status)
}

func assemble(input, output string) error {
	return exec.Command("as", "-c", input, "-o", output).Run()
}

func main() {
	// setting log
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// parse flags.
	var outpath string
	flag.StringVar(&outpath, "o", "", "The output file name")
	var help bool
	flag.BoolVar(&help, "help", false, "Help")
	var prtok bool
	flag.BoolVar(&prtok, "prtok", false, "print tokens only")
	flag.BoolVar(&isdeb, "deb", false, "debug mode or not")
	flag.Parse()

	if help {
		usage(0)
	}

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "no input files")
		usage(1)
	}

	inputPaths = flag.Args()[0:]

	// '-o' option wasn't omitted or not
	flagout := outpath != ""

	// compile
	for _, inpath := range inputPaths {
		var err error
		if !flagout {
			outpath = replaceExt(inpath, "s")
		}

		optOut, err = os.Create(outpath)
		if err != nil {
			fmt.Println(inpath)
			log.Fatal(err)
		}

		if err := compile(prtok, inpath, optOut); err != nil {
			log.Fatal(err)
		}
		objfile, err := os.Create(replaceExt(outpath, "o"))
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(objfile.Name())
		assemble(outpath, objfile.Name())
	}
}
