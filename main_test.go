package main

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

type testcase struct {
	want int
	in   string
}

var cases = map[string]testcase{
	"1":  {0, "return 0\n"},
	"2":  {42, "return 42\n"},
	"3":  {21, "return 5+20-4\n"},
	"4":  {41, "return  12 + 34 - 5\n"},
	"5":  {47, "return 5+6*7\n"},
	"6":  {15, "return 5*(9-6)\n"},
	"7":  {4, "return (3+5)/2\n"},
	"8":  {10, "return -10+20;"},
	"9":  {10, "return - -10\n"},
	"10": {10, "return - - +10\n"},

	"11": {0, "return 0==1\n"},
	"12": {1, "return 42==42\n"},
	"13": {1, "return 0!=1\n"},
	"14": {0, "return 42!=42\n"},

	"15": {1, "return 0<1\n"},
	"16": {0, "return 1<1\n"},
	"17": {0, "return 2<1\n"},
	"18": {1, "return 0<=1\n"},
	"19": {1, "return 1<=1\n"},
	"20": {0, "return 2<=1\n"},

	"21": {1, "return 1>0\n"},
	"22": {0, "return 1>1\n"},
	"23": {0, "return 1>2\n"},
	"24": {1, "return 1>=0\n"},
	"25": {1, "return 1>=1\n"},
	"26": {0, "return 1>=2\n"},

	"28": {0, "return 0==1\n42==42\n12 + 34 - 5\n0"},

	"27": {1, "return 1\n2\n3"},
	"29": {2, "1\nreturn 2\n3\n"},
	"30": {3, "1\n2\nreturn 3"},

	"31": {3, "a=3\nreturn a"},
	"32": {8, "a=3\nz=5\nreturn a+z\n"},

	"33": {3, "foo=3\nreturn foo\n"},
	"34": {8, "foo123=3\nbar=5\nreturn foo123+bar"},

	"35": {3, "if 0 {\nreturn 2\n}\nreturn 3\n"},
	"36": {3, "if 1-1{\nreturn 2\n}\nreturn 3\n"},
	"37": {2, "if 1 {\nreturn 2\n}\nreturn 3\n"},
	"38": {2, "if 2-1{\nreturn 2\n}\nreturn 3\n"},
	"39": {2, "if 2 - 1 {\nreturn 2\n}\nreturn 3\n"},

	"35-1": {3, "if 0 return 2\nreturn 3\n"},

	"40": {10, "i=0\nfor i<10 {\n\ti=i+1\n}\nreturn i\n"},
	"41": {6, "i=0\nfor {\n\ti=i+1\n\tif i>5 {\n\t\treturn i\n\t}\n}\nreturn 0\n"},

	"42": {55, "i=0\nj=0\nfor i=0; i<=10; i=i+1 {\n\tj=i+j\n}\nreturn j\n"},
	"43": {3, "for ;; {\n\treturn 3\n\treturn 5\n}"},

	"44": {3, "return ret3()"},
	"45": {5, "return ret5()"},
}

func TestCompile(t *testing.T) {

	b, err := exec.Command(
		"/bin/bash", "-c",
		"echo \"int ret3() { return 3; }\nint ret5() { return 5; }\" | gcc -xc -c -o testdata/tmp2.o -",
	).CombinedOutput()
	if err != nil {
		t.Fatalf("\noutput:\n%s\n%c", string(b), err)
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			asmN := fmt.Sprintf("testdata/asm%s.s", name)
			asm, err := os.Create(asmN)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := asm.Close(); err != nil {
					t.Fatal(err)
				}
			}()

			if err := compile(c.in, asm); err != nil {
				t.Fatal(err)
			}

			execN := fmt.Sprintf("testdata/asm%s", name)
			b, err := exec.Command("gcc", "-static", "-g", "-o", execN, asmN, "testdata/tmp2.o").CombinedOutput()
			if err != nil {
				t.Fatalf("\noutput:\n%s\n%v", string(b), err)
			}

			b, err = exec.Command(execN).CombinedOutput()
			if err != nil {
				if ee, ok := err.(*exec.ExitError); !ok {
					t.Fatalf("\noutput:\n%s\n%v", string(b), err)
				} else {
					// the return value of temporary.s is saved in exit status code normally
					actual := ee.ProcessState.ExitCode()
					if c.want != actual {
						t.Fatalf("%d expected, but got %d", c.want, actual)
					}
					t.Logf("\n%s => %d", c.in, actual)
					return
				}
			}

			// the return value of temporary.s is saved in exit status code,
			// so the below will be used only when the return value is 0.
			ans, err := exec.Command("sh", "-c", "echo $?").Output()
			if err != nil {
				t.Fatal(err)
			}

			actual, err := strconv.Atoi(strings.Trim(string(ans), "\n"))
			if err != nil {
				t.Fatal(err)
			}

			if c.want != actual {
				t.Fatalf("%d expected, but got %d", c.want, actual)
			}
			t.Logf("\n%s => %d", c.in, actual)
		})
	}
}

func TestIsSpace(t *testing.T) {
	cases := map[string]struct {
		in   rune
		want bool
	}{
		"1": {'\n', true},
		"2": {'\t', true},
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

func TestExpect(t *testing.T) {
	cases := map[string]struct {
		in   string
		want string
	}{
		"case 1": {
			in:   "\n",
			want: "",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			next := &Token{Kind: TK_EOF}
			token = &Token{
				Kind: TK_RESERVED,
				Len:  len(c.in),
				Str:  c.in,
				Next: next,
			}
			expect(c.in)
			if token.Kind != TK_EOF {
				t.Fatal("unexpected token")
			}
		})
	}
}
