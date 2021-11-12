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
	"1":  {0, "0"},
	"2":  {42, "42"},
	"3":  {21, "5+20-4"},
	"4":  {41, " 12 + 34 - 5 "},
	"5":  {47, "5+6*7"},
	"6":  {15, "5*(9-6)"},
	"7":  {4, "(3+5)/2"},
	"8":  {10, "-10+20"},
	"9":  {10, "- -10"},
	"10": {10, "- - +10"},

	"11": {0, "0==1"},
	"12": {1, "42==42"},
	"13": {1, "0!=1"},
	"14": {0, "42!=42"},

	"15": {1, "0<1"},
	"16": {0, "1<1"},
	"17": {0, "2<1"},
	"18": {1, "0<=1"},
	"19": {1, "1<=1"},
	"20": {0, "2<=1"},

	"21": {1, "1>0"},
	"22": {0, "1>1"},
	"23": {0, "1>2"},
	"24": {1, "1>=0"},
	"25": {1, "1>=1"},
	"26": {0, "1>=2"},
}

func TestCompile(t *testing.T) {

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
			b, err := exec.Command("gcc", "-static", "-g", "-o", execN, asmN).CombinedOutput()
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
					t.Logf("%s => %d", c.in, actual)
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
			t.Logf("%s => %d", c.in, actual)
		})
	}
	// asm, err := os.Create("testdata/asm.s")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// // start the test
	// if err := compile("testdata/tests", asm); err != nil {
	// 	t.Fatal(err)
	// }

	// // make a execution file
	// b1, err := exec.Command(
	// 	"gcc",
	// 	"-static",
	// 	"-g",
	// 	"-o",
	// 	"testdata/tmp",
	// 	asm.Name(),
	// ).CombinedOutput()
	// if err != nil {
	// 	t.Fatalf("\noutput:\n%s\n%v", string(b1), err)
	// }

	// // quit this test sequence, if the execution file wasn't made
	// if _, err := os.Stat(asm.Name()); err != nil {
	// 	t.Fatal(err)
	// }

	// b2, err := exec.Command("./testdata/tmp").Output()
	// if err != nil {
	// 	t.Fatalf("\noutput:\n%s\n%v", string(b2), err)
	// }
	// t.Logf("\noutput:\n%s\n", string(b2))
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
