package main

import (
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func TestCompile(t *testing.T) {

	asm, err := os.Create("testdata/asm.s")
	if err != nil {
		t.Fatal(err)
	}

	// start the test
	if err := compile("testdata/tests.c", asm); err != nil {
		t.Fatal(err)
	}

	// make tmp2.o
	const shell = "/bin/bash"
	b0, err := exec.Command(shell,
		"-c", "echo 'int char_fn() { return 257; }' | gcc -xc -c -o testdata/tmp2.o -",
	).CombinedOutput()
	if err != nil {
		t.Fatalf("\noutput:\n%s\n%v", string(b0), err)
	}

	// make a execution file
	b1, err := exec.Command(
		"gcc",
		"-static",
		"-g",
		"-o",
		"testdata/tmp",
		asm.Name(),
		"testdata/tmp2.o",
	).CombinedOutput()
	if err != nil {
		t.Fatalf("\noutput:\n%s\n%v", string(b1), err)
	}
	// quit this test sequence, if the execution file wasn't made
	if _, err := os.Stat(asm.Name()); err != nil {
		t.Fatal(err)
	}

	b2, err := exec.Command("./testdata/tmp").Output()
	if err != nil {
		t.Fatalf("\noutput:\n%s\n%v", string(b2), err)
	}
	t.Logf("\noutput:\n%s\n", string(b2))
}

func TestIsSpace(t *testing.T) {
	cases := map[string]struct {
		in   rune
		want bool
	}{
		"1": {'\t', true},
		"2": {'\n', true},
		"3": {'\v', true},
		"4": {'\f', true},
		"5": {'\r', true},
		"6": {' ', true},
		"7": {'a', false},
		"8": {'"', false},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			act := isSpace(c.in)
			if act != c.want {
				t.Fatalf("%t expected, but got %t", c.want, act)
			}
		})
	}
}

func TestStartsWithReserved(t *testing.T) {
	cases := map[string]struct {
		kw string
		in string
	}{
		"case ==": {
			kw: "==",
			in: "==0;",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			acb := startsWith(c.in, c.kw)
			if !acb {
				t.Fatal("actual is not expected")
			}
			t.Log("startsWith OK")

			ac := startsWithReserved(c.in)
			if startsWith(c.in, c.kw) && len(c.in) >= len(c.kw) && !isAlNum(rune(c.in[len(c.kw)])) {
				t.Log("true, ac: ", ac)
			} else {
				t.Log("false, ac: ", ac)
			}

			if ac != c.kw {
				t.Fatalf("%s expected, but got %s", c.kw, ac)
			}
		})
	}
}

func TestStartsWith(t *testing.T) {
	cases := map[string]struct {
		kw string
		in string
	}{
		"case ==": {
			kw: "==",
			in: "==0;",
		},
		"case //": {
			kw: "//",
			in: "// aaa",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			acb := startsWith(c.in, c.kw)
			if !acb {
				t.Fatal("actual is not expected")
			}
			t.Log("startsWith OK")
		})
	}
}

func TestSkipLineComment(t *testing.T) {
	cases := map[string]struct {
		in string
	}{
		"case //": {
			in: "// aaa",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			userInput = []rune(c.in)
			curIdx = 0
			_, err := tokenize()
			if err != nil {
				t.Fatal(err)
			}

			// printTokens()
			// fmt.Printf("%#v\n", headTok.Next)
			if headTok.Next.Kind != TK_EOF {
				t.Fatal("failed tokenize comments")
			}
		})
	}
}

func TestReadStringLiteral(t *testing.T) {
	cases := map[string]struct {
		in   string
		want []rune
	}{
		"case 1": {
			in:   `"\"\\j\"[0]"`,
			want: append([]rune("\"\\j\"[0]"), 0),
		},
		"case 2": {
			in:   `"\"\""`,
			want: append([]rune("\"\""), 0),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			tok := &Token{}
			userInput = []rune(c.in)
			curIdx = 1
			var err error
			tok, err = readStringLiteral(tok)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tok.Contents, c.want) {
				t.Fatalf("%s expected, but got %s", string(c.want), string(tok.Contents))
			}
		})
	}
}

func TestReadCharLiteral(t *testing.T) {
	cases := map[string]struct {
		in    string
		want1 int64
		want2 string
	}{
		"case 'a'": {
			in:    "'a'",
			want1: int64('a'),
			want2: "'a'",
		},
		"case '\n'": {
			in:    "'\n'",
			want1: int64('\n'),
			want2: "'\n'",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			tok := &Token{}
			userInput = []rune(c.in)
			curIdx = 0
			var err error
			tok, err = readCharLiteral(tok, curIdx)
			if err != nil {
				t.Fatal(err)
			}

			if tok.Val != int64(c.want1) {
				t.Fatalf("tok.Val: %d expected, but got %d", c.want1, tok.Val)
			}
			if tok.Str != c.want2 {
				t.Fatalf("tok.Str: %s expected, but got %s", c.want2, tok.Str)
			}
		})
	}
}

func TestTokenize(t *testing.T) {
	cases := map[string]struct {
		in   []rune
		kind []TokenKind
		str  []string
		val  []int64
	}{
		"case 'a',": {
			in:   append([]rune("3, 'a',"), 0),
			kind: []TokenKind{TK_NUM, TK_RESERVED, TK_NUM, TK_RESERVED, TK_EOF},
			str:  []string{"3", ",", "'a'", ","},
			val:  []int64{3, 0, 97, 0},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			userInput = c.in
			curIdx = 0
			tokenize()
			i := 0
			for tok := headTok.Next; tok.Next != nil; tok = tok.Next {
				if tok.Kind != c.kind[i] {
					t.Fatalf("tok.Kind: %d expected, but got %d", c.kind[i], tok.Kind)
				}
				if tok.Str != c.str[i] {
					t.Fatalf("tok.Str: %s expected, but got %s", c.str[i], tok.Str)
				}
				if tok.Val != c.val[i] {
					t.Fatalf("tok.Val: %d expected, but got %d", c.val[i], tok.Val)
				}
				i++
			}
		})
	}
}
