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
	"1":  {0, "func main() int {\n\treturn 0\n}"},
	"2":  {42, "func main() int {\n\treturn 42\n}"},
	"3":  {21, "func main() int {\n\treturn 5+20-4\n}"},
	"4":  {41, "func main() int {\n\treturn  12 + 34 - 5\n}"},
	"5":  {47, "func main() int {\n\treturn 5+6*7\n}"},
	"6":  {15, "func main() int {\n\treturn 5*(9-6)\n}"},
	"7":  {4, "func main() int {\n\treturn (3+5)/2\n}"},
	"8":  {10, "func main() int {\n\treturn -10+20;\n}"},
	"9":  {10, "func main() int {\n\treturn - -10\n}"},
	"10": {10, "func main() int {\n\treturn - - +10\n}"},

	"11": {0, "func main() int {\n\treturn 0==1\n}"},
	"12": {1, "func main() int {\n\treturn 42==42\n}"},
	"13": {1, "func main() int {\n\treturn 0!=1\n}"},
	"14": {0, "func main() int {\n\treturn 42!=42\n}"},

	"15": {1, "func main() int {\n\treturn 0<1\n}"},
	"16": {0, "func main() int {\n\treturn 1<1\n}"},
	"17": {0, "func main() int {\n\treturn 2<1\n}"},
	"18": {1, "func main() int {\n\treturn 0<=1\n}"},
	"19": {1, "func main() int {\n\treturn 1<=1\n}"},
	"20": {0, "func main() int {\n\treturn 2<=1\n}"},

	"21": {1, "func main() int {\n\treturn 1>0\n}"},
	"22": {0, "func main() int {\n\treturn 1>1\n}"},
	"23": {0, "func main() int {\n\treturn 1>2\n}"},
	"24": {1, "func main() int {\n\treturn 1>=0\n}"},
	"25": {1, "func main() int {\n\treturn 1>=1\n}"},
	"26": {0, "func main() int {\n\treturn 1>=2\n}"},

	"28": {0, "func main() int {\n\treturn 0==1\n\t42==42\n\t12 + 34 - 5\n\t0\n}"},

	"27": {1, "func main() int {\n\treturn 1\n\t2\n\t3\n}"},
	"29": {2, "func main() int {\n\t1\n\treturn 2\n\t3\n}"},
	"30": {3, "func main() int {\n\t1\n\t2\n\treturn 3\n}"},

	"31": {3, "func main() int {\n\tvar a int=3\n\treturn a\n}"},
	"32": {8, "func main() int {\n\tvar a int=3\n\tvar z int=5\n\treturn a+z\n}"},

	"33": {3, "func main() int {\n\tvar foo int=3\n\treturn foo\n}"},
	"34": {8, "func main() int {\n\tvar foo123 int=3\n\tvar bar int=5\n\treturn foo123+bar\n}"},

	"35": {3, "func main() int {\n\tif 0 {\n\t\treturn 2\n\t}\n\treturn 3\n}"},
	"36": {3, "func main() int {\n\tif 1-1{\n\t\treturn 2\n\t}\n\treturn 3\n}"},
	"37": {2, "func main() int {\n\tif 1 {\n\t\treturn 2\n\t}\n\treturn 3\n}"},
	"38": {2, "func main() int {\n\tif 2-1{\n\t\treturn 2\n\t}\n\treturn 3\n}"},
	"39": {2, "func main() int {\nif 2 - 1 {\n\t\treturn 2\n\t}\n\treturn 3\n}"},

	"35-1": {3, "func main() int {\n\tif 0 return 2\n\treturn 3\n}"},

	"40": {10, "func main() int {\n\tvar i int=0\n\tfor i<10 {\n\t\ti=i+1\n\t}\n\treturn i\n}"},
	"41": {6, "func main() int {\n\tvar i int=0\n\tfor {\n\t\ti=i+1\n\t\tif i>5 {\n\t\t\treturn i\n\t\t}\n\t}\n\treturn 0\n}"},

	"42": {55, "func main() int {\n\tvar i int=0\n\tvar j int=0\n\tfor i=0; i<=10; i=i+1 {\n\t\tj=i+j\n\t}\n\treturn j\n}"},
	"43": {3, "func main() int {\nfor ;; {\n\treturn 3\n\treturn 5\n}\n}"},

	"44": {3, "func main() int {\n\treturn ret3()\n}"},
	"45": {5, "func main() int {\n\treturn ret5()\n}"},
	"46": {8, "func main() int {\n\treturn add(3, 5)\n}"},
	"47": {2, "func main() int {\n\treturn sub(5, 3)\n}"},
	"48": {21, "func main() int {\n\treturn add6(1,2,3,4,5,6)\n}"},
	"49": {7, "func main() int {\n\treturn add2(3,4)\n}\nfunc add2(x int,y int) int {\n\treturn x+y\n}"},
	"50": {1, "func main() int {\n\treturn sub2(4,3)\n}\nfunc sub2(x int,y int) int {\n\treturn x-y\n}"},
	"51": {55, "func main() int {\n\treturn fib(9)\n}\nfunc fib(x int) int {\n\tif x<=1 {\n\t\treturn 1\n\t}\n\treturn fib(x-1) + fib(x-2)\n}"},

	"52": {3, "func main() int {\n\tvar x int=3\n\treturn *&x\n}"},
	"53": {3, "func main() int {\n\tvar x int=3\n\tvar y *int=&x\n\tvar z **int=&y\n\treturn **z\n}"},
	"54": {5, "func main() int {\n\tvar x int=3\n\tvar y int=5\n\treturn *(&x+1)\n}"},
	"55": {3, "func main() int {\n\tvar x int=3\n\tvar y int=5\n\treturn *(&y-1)\n}"},
	"56": {5, "func main() int {\n\tvar x int=3\n\tvar y *int=&x\n\t*y=5\n\treturn x\n}"},
	"57": {7, "func main() int {\n\tvar x int=3\n\tvar y int=5\n\t*(&x+1)=7\n\treturn y\n}"},
	"58": {7, "func main() int {\n\tvar x int=3\n\tvar y int=5\n\t*(&y-1)=7\n\treturn x\n}"},
	"59": {8, "func main() int {\n\tvar x int=3\n\tvar y int=5\n\treturn foo(&x, y)\n}\nfunc foo(x *int, y int) int {\n\treturn *x + y\n}"},

	"60": {3, "func main() int {\n\tvar x [2]int\n\tvar y *int=&x\n\t*y=3\n\treturn *x\n}"},

	"61": {3, "func main() int {\n\tvar x [3]int\n\t*x=3\n\t*(x+1)=4\n\t*(x+2)=5\n\treturn *x\n}"},
	"62": {4, "func main() int {\n\tvar x [3]int\n\t*x=3\n\t*(x+1)=4\n\t*(x+2)=5\n\treturn *(x+1)\n}"},
	"63": {5, "func main() int {\n\tvar x [3]int\n\t*x=3\n\t*(x+1)=4\n\t*(x+2)=5\n\treturn *(x+2)\n}"},

	"64": {0, "func main() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\t*y=0\n\treturn **x\n}"},
	"65": {1, "func main() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\t*(y+1)=1\n\treturn *(*x+1)\n}"},
	"66": {2, "func main() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\t*(y+2)=2\n\treturn *(*x+2)\n}"},
	"67": {3, "func main() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\t*(y+3)=3\n\treturn **(x+1)\n}"},
	"68": {4, "func main() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\t*(y+4)=4\n\treturn *(*(x+1)+1)\n}"},
	"69": {5, "func main() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\t*(y+5)=5\n\treturn *(*(x+1)+2)\n}"},
	"70": {6, "func main() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\t*(y+6)=6\n\treturn **(x+2)\n}"},

	"71": {3, "func main() int {\n\tvar x [3]int\n\t*x=3\n\tx[1]=4\n\tx[2]=5\n\treturn *x\n}"},
	"72": {4, "func main() int {\n\tvar x [3]int\n\t*x=3\n\tx[1]=4\n\tx[2]=5\n\treturn *(x+1)\n}"},
	"73": {5, "func main() int {\n\tvar x [3]int\n\t*x=3\n\tx[1]=4\n\tx[2]=5\n\treturn *(x+2)\n}"},
	"74": {5, "func main() int {\n\tvar x [3]int\n\t*x=3\n\tx[1]=4\n\tx[2]=5\n\treturn *(x+2)\n}"},
	"75": {5, "func main() int {\n\tvar x [3]int\n\t*x=3\n\tx[1]=4\n\t2[x]=5\n\treturn *(x+2)\n}"},

	"76": {0, "func main() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\ty[0]=0\n\treturn x[0][0]\n}"},
	"77": {1, "func main() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\ty[1]=1\n\treturn x[0][1]\n}"},
	"78": {2, "func main() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\ty[2]=2\n\treturn x[0][2]\n}"},
	"79": {3, "func main() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\ty[3]=3\n\treturn x[1][0]\n}"},
	"80": {4, "func main() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\ty[4]=4\n\treturn x[1][1]\n}"},
	"81": {5, "func main() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\ty[5]=5\n\treturn x[1][2]\n}"},
	"82": {6, "func main() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\ty[6]=6\n\treturn x[2][0]\n}"},

	"83": {0, "var x int\nfunc main() int {\n\treturn x\n}"},
	"84": {3, "var x int\nfunc main() int {\n\tx=3\n\treturn x\n}"},
	"85": {0, "var x [4]int\nfunc main() int {\n\tx[0]=0\n\tx[1]=1\n\tx[2]=2\n\tx[3]=3\n\treturn x[0]\n}"},
	"86": {1, "var x [4]int\nfunc main() int {\n\tx[0]=0\n\tx[1]=1\n\tx[2]=2\n\tx[3]=3\n\treturn x[1]\n}"},
	"87": {2, "var x [4]int\nfunc main() int {\n\tx[0]=0\n\tx[1]=1\n\tx[2]=2\n\tx[3]=3\n\treturn x[2]\n}"},
	"88": {3, "var x [4]int\nfunc main() int {\n\tx[0]=0\n\tx[1]=1\n\tx[2]=2\n\tx[3]=3\n\treturn x[3]\n}"},

	"89": {1, "func main() int {\n\tvar x byte=1\n\treturn x\n}"},
	"90": {1, "func main() int {\n\tvar x byte=1\n\tvar y byte=2\n\treturn x\n}"},
	"91": {2, "func main() int {\n\tvar x byte=1\n\tvar y byte=2\n\treturn y\n}"},

	"92": {1, "func main() int {\n\treturn sub_char(7, 3, 3)\n}\nfunc sub_char(a byte, b byte, c byte) int {\n\treturn a-b-c\n}"},
}

var tmp2 string = `int ret3() { return 3; }
int ret5() { return 5; }
int add(int x, int y) { return x+y; }
int sub(int x, int y) { return x-y; }

int add6(int a, int b, int c, int d, int e, int f) {
  return a+b+c+d+e+f;
}`

func TestCompile(t *testing.T) {

	b, err := exec.Command(
		"/bin/bash", "-c",
		"echo \""+tmp2+"\" | gcc -xc -c -o testdata/tmp2.o -",
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
			b, err := exec.Command(
				"gcc",
				"-static",
				"-g",
				"-o",
				execN,
				asmN,
				"testdata/tmp2.o",
			).CombinedOutput()
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

// func TestGlobalVar(t *testing.T) {
// 	userInput = []rune("var x [2]int\n")
// 	var err error
// 	token, err = tokenize()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	globalVar()
// 	t.Logf("globals: %#v\n", globals)
// 	t.Logf("globals.Var: %#v\n", globals.Var)
// 	t.Logf("globals.Var.Ty: %#v\n", globals.Var.Ty)
// 	t.Logf("globals.Var.Ty.Base: %#v\n", globals.Var.Ty.Base)
// }

// func TestStmt(t *testing.T) {
// 	filename = "test"
// 	userInput = []rune("var z *int\n**z")
// 	var err error
// 	token, err = tokenize()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	printTokens()

// 	head := &Node{}
// 	cur := head
// 	for !atEof() {
// 		cur.Next = stmt()
// 		cur = cur.Next
// 	}
// 	n := head.Next

// 	t.Logf("node: %#v\n\nlocals.Var: %#v\n\n%#v\n", n, locals.Var, locals.Var.Ty)
// 	walkInOrder(n)
// }

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
		"case func": {
			kw: "func",
			in: "func ",
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
