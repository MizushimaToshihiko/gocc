package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// flags
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

	if !exists(arg) && arg != "" && arg != "-" {
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

func findLibpath() (string, error) {
	if exists("/usr/lib/x86_64-linux-gnu/crti.o") {
		return "/usr/lib/x86_64-linux-gnu", nil
	}
	if exists("/usr/lib64/crti.o") {
		return "/usr/lib64", nil
	}

	return "", errors.New("library path is not found")
}

func findGccLibPath() (string, error) {
	paths := []string{
		"/usr/lib/gcc/x86_64-linux-gnu/*/crtbegin.o",
		"/usr/lib/gcc/x86_64-pc-linux-gnu/*/crtbegin.o", // For Gentoo
		"/usr/lib/gcc/x86_64-redhat-linux/*/crtbegin.o", // For Fedora
	}

	for _, path := range paths {
		p, err := filepath.Glob(path)
		if p != nil {
			return filepath.Dir(p[len(p)-1]), err
		}
	}

	return "", errors.New("gcc library path is not found")
}

func runLinker(inputs []string, output string) error {
	arr := []string{}

	arr = append(arr, "ld")
	arr = append(arr, "-o")
	arr = append(arr, output)
	arr = append(arr, "-m")
	arr = append(arr, "elf_x86_64")
	arr = append(arr, "-dynamic-linker")
	arr = append(arr, "/lib64/ld-linux-x86-64.so.2")

	libPath, err := findLibpath()
	// fmt.Println("libPath:", libPath)
	if err != nil {
		return fmt.Errorf("findLibPath: %s", err)
	}
	gccLibPath, err := findGccLibPath()
	// fmt.Println("gccLibPath:", gccLibPath)
	if err != nil {
		return fmt.Errorf("findGccLibPath: %s", err)
	}

	arr = append(arr, fmt.Sprintf("%s/crt1.o", libPath))
	arr = append(arr, fmt.Sprintf("%s/crti.o", libPath))
	arr = append(arr, fmt.Sprintf("%s/crtbegin.o", libPath))
	arr = append(arr, fmt.Sprintf("-L%s", gccLibPath))
	arr = append(arr, fmt.Sprintf("-L%s", libPath))
	arr = append(arr, fmt.Sprintf("-L%s/..", libPath))
	arr = append(arr, "-L/usr/lib64")
	arr = append(arr, "-L/lib64")
	arr = append(arr, "-L/usr/lib/x86_64-linux-gnu")
	arr = append(arr, "-L/usr/lib/x86_64-pc-linux-gnu")
	arr = append(arr, "-L/usr/lib/x86_64-redhat-linux")
	arr = append(arr, "-L/usr/lib")
	arr = append(arr, "-L/lib")

	arr = append(arr, inputs...)

	arr = append(arr, "-lc")
	arr = append(arr, "-lgcc")
	arr = append(arr, "--as-needed")
	arr = append(arr, "-lgcc_s")
	arr = append(arr, "--no-as-needed")
	arr = append(arr, fmt.Sprintf("%s/crtend.o", gccLibPath))
	arr = append(arr, fmt.Sprintf("%s/crtn.o", libPath))

	return exec.Command(arr[0], arr[1:]...).Run()
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
	var optc bool
	flag.BoolVar(&optc, "c", false, "compile and assemble only, or not")
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
	opto := outpath != ""

	ldArgs := []string{}

	// compile
	for _, inpath := range inputPaths {

		var err error
		if !opto {
			outpath = replaceExt(inpath, "s")
		}

		optOut, err = os.Create(outpath)
		if err != nil {
			fmt.Println(inpath)
			log.Fatal(err)
		}

		fmt.Println("outpath:", outpath)
		fmt.Println("optOut.Name():", optOut.Name())

		if optc {
			// make the gnu-assembly file
			if err := compile(prtok, inpath, optOut); err != nil {
				log.Fatal(err)
			}

			// make the object file
			objfile, err := os.Create(replaceExt(outpath, "o"))
			if err != nil {
				log.Fatal(err)
			}
			err = assemble(outpath, objfile.Name())
			if err != nil {
				log.Fatal(err)
			}
			continue
		}

		if !exists("./testdata/tmp") {
			err := os.MkdirAll("./testdata/tmp", 0755)
			if err != nil {
				log.Fatal(err)
			}
		}
		// Compile, assemble and link
		tmp1, err := os.CreateTemp("./testdata/tmp", "temp1_*.s")
		if err != nil {
			log.Fatal(err)
		}
		tmp2, err := os.CreateTemp("./testdata/tmp", "temp2_*.o")
		if err != nil {
			log.Fatal(err)
		}
		if err := compile(prtok, inpath, tmp1); err != nil {
			log.Fatal(err)
		}
		if err := assemble(tmp1.Name(), tmp2.Name()); err != nil {
			log.Fatal(err)
		}
		ldArgs = append(ldArgs, tmp2.Name())
	}

	if len(ldArgs) > 0 {
		if !opto {
			outpath = "a.out"
		}
		err := runLinker(ldArgs, outpath)
		if err != nil {
			log.Fatal(err)
		}
		if err := os.RemoveAll("./testdata/tmp"); err != nil {
			log.Fatal(err)
		}
	}
}
