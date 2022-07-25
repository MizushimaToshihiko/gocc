package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// Temporary files for compiling and assembling
var tmpfiles []string

// flags
var inputPaths []string
var isdeb bool // Is debug mode or not.
var optE bool  // -E option

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

	// If 'prtok' option is given, just print tokens for debugging.
	if prtok {
		printTokens2(w, tok)
		return nil
	}

	tok = preprocess(tok)

	// If -E is given, print out preprocessed C code as a result.
	if optE {
		printTokens(w, tok)
		return nil
	}

	if isdeb {
		printTokens2(w, tok)
	}

	// parse the tokens, and make the AST nodes.
	prog := parse(tok)
	if err != nil {
		return err
	}

	return codegen(w, prog) // make the asm code, down on the AST
}

func usage(status int) {
	fmt.Fprintf(os.Stderr, "usage: ./bin/gocc [ -o <path> ] <file>\n")
	os.Exit(status)
}

func cleanup() error {
	for _, tmp := range tmpfiles {
		if err := os.Remove(tmp); err != nil {
			return err
		}
	}
	return nil
}

func createTmpFile() (*os.File, error) {
	tmp, err := os.CreateTemp("/tmp", "gocc-*")
	if err != nil {
		return nil, err
	}

	tmpfiles = append(tmpfiles, tmp.Name())
	return tmp, nil
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
	var err error
	arr := []string{}

	arr = append(arr, "ld")
	arr = append(arr, "-o")
	arr = append(arr, output)
	arr = append(arr, "-m")
	arr = append(arr, "elf_x86_64")
	arr = append(arr, "-dynamic-linker")
	arr = append(arr, "/lib64/ld-linux-x86-64.so.2")

	libPath, err := findLibpath()
	if err != nil {
		return fmt.Errorf("findLibPath: %s", err)
	}
	gccLibPath, err := findGccLibPath()
	if err != nil {
		return fmt.Errorf("findGccLibPath: %s", err)
	}

	arr = append(arr, fmt.Sprintf("%s/crt1.o", libPath))
	arr = append(arr, fmt.Sprintf("%s/crti.o", libPath))
	arr = append(arr, fmt.Sprintf("%s/crtbegin.o", gccLibPath))
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

	cmd := exec.Command(arr[0], arr[1:]...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(out.String())
		return fmt.Errorf(fmt.Sprint(err) + ":\n" + stderr.String())
	}

	return nil
}

func main() {
	defer func() {
		if err := cleanup(); err != nil {
			log.Fatal(err)
		}
	}()

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
	var optS bool
	flag.BoolVar(&optS, "S", false, "compile only or not")
	flag.BoolVar(&optE, "E", false, "stop after the preprocessing stage")
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

		if optS && outpath != "-" {
			outpath = replaceExt(inpath, "s")
		} else if !opto {
			outpath = replaceExt(inpath, "o")
		}

		// Just preprocess
		if optE || prtok {
			var out *os.File
			var err error
			if !opto {
				out = os.Stdout
			} else {
				switch outpath {
				case "-":
					out = os.Stdout
				default:
					out, err = os.Create(outpath)
					if err != nil {
						log.Fatal(err)
					}
				}
			}

			if err := compile(prtok, inpath, out); err != nil {
				log.Fatal(err)
			}
			continue
		}

		// Just compile
		if optS {
			var out *os.File
			var err error
			switch outpath {
			case "-":
				out = os.Stdout
			default:
				out, err = os.Create(outpath)
				if err != nil {
					log.Fatal(err)
				}
			}

			if err := compile(prtok, inpath, out); err != nil {
				log.Fatal(err)
			}
			continue
		}

		if optc {
			temp, err := createTmpFile()
			if err != nil {
				log.Fatal(err)
			}
			// make the gnu-assembly file
			if err := compile(prtok, inpath, temp); err != nil {
				log.Fatal(err)
			}

			// make the object file
			objfile, err := os.Create(outpath)
			if err != nil {
				log.Fatal(err)
			}
			err = assemble(temp.Name(), objfile.Name())
			if err != nil {
				log.Fatal(err)
			}
			continue
		}

		// Compile, assemble and link
		tmp1, err := createTmpFile()
		if err != nil {
			log.Fatal(err)
		}
		tmp2, err := createTmpFile()
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

		if !exists(outpath) {
			_, err := os.Create(outpath)
			if err != nil {
				log.Fatal(err)
			}
		}

		err := runLinker(ldArgs, outpath)
		if err != nil {
			log.Fatal(err)
		}
	}
}
