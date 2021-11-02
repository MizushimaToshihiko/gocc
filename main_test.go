package main

import (
	"os"
	"os/exec"
	"reflect"
	"testing"
)

// "error 1": {0, "return a;"},
// "error 2": {0, "int return a;"},
// "error 3": {0, "int main(){ return 1}"},
// "error 4": {0, "int main() {int return a;"},
// "error 5": {0, "int main() { x = y + + 5;}"},
// "error 6": {0, "int main() { int x; int y; y = 1; x = y + + 5;}"},
// "error 7": {0, "int main() { /* return 2;} "},

func TestCompile(t *testing.T) {

	asm, err := os.Create("testdata/asm.s")
	if err != nil {
		t.Fatal(err)
	}

	// start the test
	if err := compile("testdata/tests.c", asm); err != nil {
		t.Fatal(err)
	}

	// make a execution file with static-link to 'f'
	b1, err := exec.Command("gcc", "-static", "-o", "testdata/tmp", asm.Name()).CombinedOutput()
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

// func TestFindLVar(t *testing.T) {
// 	cases := map[string]struct {
// 		lvar *LVar
// 		tok  *Token
// 	}{
// 		"case 1": {
// 			&LVar{Name: "x"},
// 			&Token{Str: "x", Len: 1},
// 		},
// 	}

// 	for name, c := range cases {
// 		t.Run(name, func(t *testing.T) {
// 			locals = &VarList{Var: c.lvar, Next: nil}
// 			lv := findLVar(c.tok)
// 			fmt.Printf("%#v\n", lv)
// 		})
// 	}
// }

func TestStartsWithReserved(t *testing.T) {
	cases := map[string]struct {
		kw string
		in string
	}{
		"case ==": {
			kw: "==",
			in: "==0;",
		},
		// "case //": {
		// 	kw: "//",
		// 	in: "// aaa",
		// },
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
