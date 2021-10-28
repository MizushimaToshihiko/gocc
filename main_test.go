package main

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

type testcase struct {
	expected int
	input    string
}

var cases = map[string]testcase{
	"1":  {0, "int main() { return 0; }"},
	"2":  {42, "int main() { return 42; }"},
	"3":  {21, "int main() { return 5+20-4; }"},
	"4":  {41, "int main() { return  12 + 34 - 5 ; }"},
	"5":  {47, "int main() { return 5+6*7; }"},
	"6":  {15, "int main() { return 5*(9-6); }"},
	"7":  {4, "int main() { return (3+5)/2; }"},
	"8":  {10, "int main() { return -10+20; }"},
	"9":  {10, "int main() { return - -10; }"},
	"10": {10, "int main() { return - - +10; }"},

	"11": {0, "int main() { return 0==1; }"},
	"12": {1, "int main() { return 42==42; }"},
	"13": {1, "int main() { return 0!=1; }"},
	"14": {0, "int main() { return 42!=42; }"},

	"15": {1, "int main() { return 0<1; }"},
	"16": {0, "int main() { return 1<1; }"},
	"17": {0, "int main() { return 2<1; }"},
	"18": {1, "int main() { return 0<=1; }"},
	"19": {1, "int main() { return 1<=1; }"},
	"20": {0, "int main() { return 2<=1; }"},

	"21": {1, "int main() { return 1>0; }"},
	"22": {0, "int main() { return 1>1; }"},
	"23": {0, "int main() { return 1>2; }"},
	"24": {1, "int main() { return 1>=0; }"},
	"25": {1, "int main() { return 1>=1; }"},
	"26": {0, "int main() { return 1>=2; }"},

	"28": {3, "int main() { int a=3; return a; }"},
	"29": {8, "int main() { int a=3; int z=5; return a+z; }"},

	"30": {1, "int main() { return 1; 2; 3; }"},
	"31": {2, "int main() { 1; return 2; 3; }"},
	"32": {3, "int main() { 1; 2; return 3; }"},

	"33": {3, "int main() { int foo=3; return foo; }"},
	"34": {8, "int main() { int foo123=3; int bar=5; return foo123+bar; }"},

	"35": {3, "int main() { if (0) return 2; return 3; }"},
	"36": {3, "int main() { if (1-1) return 2; return 3; }"},
	"37": {2, "int main() { if (1) return 2; return 3; }"},
	"38": {2, "int main() { if (2-1) return 2; return 3; }"},

	"39": {3, "int main() { {1; {2;} return 3;} }"},

	"40": {10, "int main() { int i=0; i=0; while(i<10) i=i+1; return i; }"},
	"41": {55, "int main() { int i=0; int j=0; while(i<=10) {j=i+j; i=i+1;} return j; }"},

	"42": {55, "int main() { int i=0; int j=0; for (i=0; i<=10; i=i+1) j=i+j; return j; }"},
	"43": {3, "int main() { for (;;) return 3; return 5; }"},

	"44": {3, "int main() { return ret3(); }"},
	"45": {5, "int main() { return ret5(); }"},
	"46": {8, "int main() { return add(3, 5); }"},
	"47": {2, "int main() { return sub(5, 3); }"},
	"48": {21, "int main() { return add6(1,2,3,4,5,6); }"},

	"49": {32, "int main() { return ret32(); } int ret32() { return 32; }"},
	"50": {7, "int main() { return add2(3,4); } int add2(int x,int y) { return x+y; }"},
	"51": {1, "int main() { return sub2(4,3); } int sub2(int x,int y) { return x-y; }"},
	"52": {55, "int main() { return fib(9); } int fib(int x) { if (x<=1) return 1; return fib(x-1) + fib(x-2); }"},

	"53": {3, "int main() { int x=3; return *&x; }"},
	"54": {3, "int main() { int x=3; int *y=&x; int **z=&y; return **z; }"},
	"55": {5, "int main() { int x=3; int y=5; return *(&x+1); }"},
	"56": {5, "int main() { int x=3; int y=5; return *(1+&x); }"},
	"57": {3, "int main() { int x=3; int y=5; return *(&y-1); }"},
	"58": {5, "int main() { int x=3; int y=5; int *z=&x; return *(z+1); }"},
	"59": {3, "int main() { int x=3; int y=5; int *z=&y; return *(z-1); }"},
	"60": {5, "int main() { int x=3; int *y=&x; *y=5; return x; }"},
	"61": {7, "int main() { int x=3; int y=5; *(&x+1)=7; return y; }"},
	"62": {7, "int main() { int x=3; int y=5; *(&y-1)=7; return x; }"},
	"63": {8, "int main() { int x=3; int y=5; return foo(&x, y); } int foo(int *x, int y) { return *x + y; }"},

	"64": {8, "int main() {int a; return sizeof(a);}"},
	"65": {8, "int main() {int *a; return sizeof(a);}"},

	"66": {3, "int main() { int x[3]; *x=3; *(x+1)=4; *(x+2)=5; return *x; }"},
	"67": {4, "int main() { int x[3]; *x=3; *(x+1)=4; *(x+2)=5; return *(x+1); }"},
	"68": {5, "int main() { int x[3]; *x=3; *(x+1)=4; *(x+2)=5; return *(x+2); }"},

	"69": {0, "int main() { int x[2][3]; int *y=x; *y=0; return **x; }"},
	"70": {1, "int main() { int x[2][3]; int *y=x; *(y+1)=1; return *(*x+1); }"},
	"71": {2, "int main() { int x[2][3]; int *y=x; *(y+2)=2; return *(*x+2); }"},
	"72": {3, "int main() { int x[2][3]; int *y=x; *(y+3)=3; return **(x+1); }"},
	"73": {4, "int main() { int x[2][3]; int *y=x; *(y+4)=4; return *(*(x+1)+1); }"},
	"74": {5, "int main() { int x[2][3]; int *y=x; *(y+5)=5; return *(*(x+1)+2); }"},
	"75": {6, "int main() { int x[2][3]; int *y=x; *(y+6)=6; return **(x+2); }"},

	"76": {3, "int main() { int x[3]; *x=3; x[1]=4; x[2]=5; return *x; }"},
	"77": {4, "int main() { int x[3]; *x=3; x[1]=4; x[2]=5; return *(x+1); }"},
	"78": {5, "int main() { int x[3]; *x=3; x[1]=4; x[2]=5; return *(x+2); }"},
	"79": {5, "int main() { int x[3]; *x=3; x[1]=4; x[2]=5; return *(x+2); }"},
	"80": {5, "int main() { int x[3]; *x=3; x[1]=4; 2[x]=5; return *(x+2); }"},

	"81": {0, "int main() { int x[2][3]; int *y=x; y[0]=0; return x[0][0]; }"},
	"82": {1, "int main() { int x[2][3]; int *y=x; y[1]=1; return x[0][1]; }"},
	"83": {2, "int main() { int x[2][3]; int *y=x; y[2]=2; return x[0][2]; }"},
	"84": {3, "int main() { int x[2][3]; int *y=x; y[3]=3; return x[1][0]; }"},
	"85": {4, "int main() { int x[2][3]; int *y=x; y[4]=4; return x[1][1]; }"},
	"86": {5, "int main() { int x[2][3]; int *y=x; y[5]=5; return x[1][2]; }"},
	"87": {6, "int main() { int x[2][3]; int *y=x; y[6]=6; return x[2][0]; }"},

	"88": {0, "int x; int main() { return x; }"},
	"89": {3, "int x; int main() { x=3; return x; }"},
	"90": {0, "int x[4]; int main() { x[0]=0; x[1]=1; x[2]=2; x[3]=3; return x[0]; }"},
	"91": {1, "int x[4]; int main() { x[0]=0; x[1]=1; x[2]=2; x[3]=3; return x[1]; }"},
	"92": {2, "int x[4]; int main() { x[0]=0; x[1]=1; x[2]=2; x[3]=3; return x[2]; }"},
	"93": {3, "int x[4]; int main() { x[0]=0; x[1]=1; x[2]=2; x[3]=3; return x[3]; }"},

	"94": {8, "int x; int main() { return sizeof(x); }"},
	"95": {32, "int x[4]; int main() { return sizeof(x); }"},

	"96": {1, "int main() { char x=1; return x; }"},
	"97": {1, "int main() { char x=1; char y=2; return x; }"},
	"98": {2, "int main() { char x=1; char y=2; return y; }"},

	"99":  {1, "int main() { char x; return sizeof(x); }"},
	"100": {10, "int main() { char x[10]; return sizeof(x); }"},
	"101": {1, "int main() { return sub_char(7, 3, 3); } int sub_char(char a, char b, char c) { return a-b-c; }"},

	"102": {97, "int main() { return \"abc\"[0]; }"},
	"103": {98, "int main() { return \"abc\"[1]; }"},
	"104": {99, "int main() { return \"abc\"[2]; }"},
	"105": {0, "int main() { return \"abc\"[3]; }"},
	"106": {4, "int main() { return sizeof(\"abc\"); }"},
	// "error 1": {0, "return a;"},
	// "error 2": {0, "int return a;"},
	// "error 3": {0, "int main(){ return 1}"},
	// "error 4": {0, "int main() {int return a;"},
	// "error 5": {0, "int main() { x = y + + 5;}"},
	// "error 6": {0, "int main() { int x; int y; y = 1; x = y + + 5;}"},
}

var funcs string = `int ret3() { return 3;}
int ret5() { return 5;}
int add(int x, int y) { return x+y; }
int sub(int x, int y) { return x-y;}

int add6(int a, int b, int c, int d, int e, int f) {
	return a+b+c+d+e+f;
}
`

func TestCompile(t *testing.T) {
	var asmName string = "temp"

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			// make a asm file
			asm, err := os.Create(asmName + ".s")
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := os.Remove(asm.Name()); err != nil {
					t.Fatal(err)
				}
			}()

			// make a input file
			in, err := os.Create("in.c")
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := os.Remove(in.Name()); err != nil {
					t.Fatal(err)
				}
			}()

			// write 'c.input' to 'in' file
			if _, err = in.WriteString(c.input); err != nil {
				t.Fatal(err)
			}
			if err = in.Sync(); err != nil {
				t.Fatal(err)
			}

			// start the test
			if err := compile(in.Name(), asm); err != nil {
				t.Fatal(err)
			}

			// make 'funcs_file' file
			f, err := os.Create("funcs_file")
			if err != nil {
				return
			}
			defer func() {
				if err = os.Remove(f.Name()); err != nil {
					t.Fatal(err)
				}
			}()

			// write 'funcs' to 'funcs_file'
			if _, err = f.WriteString(funcs); err != nil {
				t.Fatal(err)
			}
			if err = f.Sync(); err != nil {
				t.Fatal(err)
			}

			// make a object file from 'funcs_file'
			_, err = exec.Command("gcc", "-xc", "-c", "-o", asmName+"2.o", f.Name()).Output()
			if err != nil {
				t.Fatal(err)
			}

			// make a execution file with static-link to 'f'
			b, err := exec.Command("gcc", "-static", "-g", "-o", asmName, asm.Name(), asmName+"2.o").Output()
			if err != nil {
				t.Fatalf("\noutput: %s\n%v", string(b), err)
			}
			// quit this test sequence, if the execution file wasn't made
			if _, err := os.Stat(asm.Name()); err != nil {
				t.Fatal(err)
			}

			b, err = exec.Command("./" + asmName).Output()
			if err != nil {
				if ee, ok := err.(*exec.ExitError); !ok {
					t.Fatal(err)
				} else {
					// the return value of temporary.s is saved in exit status code normally,
					// except for the return value is 0.
					actual := ee.ProcessState.ExitCode()
					if c.expected != actual {
						t.Fatalf("\noutput: %s\n%s\n%d expected, but got %d", string(b), c.input, c.expected, actual)
					}
					t.Logf("%s => %d", c.input, actual)
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

			if c.expected != actual {
				t.Fatalf("%d expected, but got %d", c.expected, actual)
			}
			t.Logf("%s => %d", c.input, actual)
		})
	}
}

// func TestIsSpace(t *testing.T) {
// 	cases := map[string]struct {
// 		in   byte
// 		want bool
// 	}{
// 		"1": {'\t', true},
// 		"2": {'\n', true},
// 		"3": {'\v', true},
// 		"4": {'\f', true},
// 		"5": {'\r', true},
// 		"6": {' ', true},
// 		"7": {'a', false},
// 		"8": {'"', false},
// 	}

// 	for name, c := range cases {
// 		t.Run(name, func(t *testing.T) {
// 			act := isSpace(c.in)
// 			if act != c.want {
// 				t.Fatalf("%t expected, but got %t", c.want, act)
// 			}
// 		})
// 	}
// }

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

// func TestStartWith(t *testing.T) {
// 	kw := "return"
// 	in := "return return;"

// 	acb := startsWith(in, kw)
// 	if !acb {
// 		t.Fatal("actual is not expected")
// 	}
// 	t.Log("startsWith OK")

// 	ac := startsWithReserved(in)
// 	if startsWith(in, kw) && len(in) > len(kw) && !isAlNum(in[len(kw)]) {
// 		t.Log("true")
// 	} else {
// 		t.Log("false")
// 	}

// 	if ac != kw {
// 		t.Fatalf("%s expected, but got %s", kw, ac)
// 	}
// }
