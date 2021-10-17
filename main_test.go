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
	// "1":  {0, "return 0;"},
	// "2":  {42, "return 42;"},
	// "3":  {21, "return 5+20-4;"},
	// "4":  {41, "return  12 + 34 - 5 ;"},
	// "5":  {47, "return 5+6*7;"},
	// "6":  {15, "return 5*(9-6);"},
	// "7":  {4, "return (3+5)/2;"},
	// "8":  {10, "return -10+20;"},
	// "9":  {10, "return - -10;"},
	// "10": {10, "return - - +10;"},

	// "11": {0, "return 0==1;"},
	// "12": {1, "return 42==42;"},
	// "13": {1, "return 0!=1;"},
	// "14": {0, "return 42!=42;"},

	// "15": {1, "return 0<1;"},
	// "16": {0, "return 1<1;"},
	// "17": {0, "return 2<1;"},
	// "18": {1, "return 0<=1;"},
	// "19": {1, "return 1<=1;"},
	// "20": {0, "return 2<=1;"},

	// "21": {1, "return 1>0;"},
	// "22": {0, "return 1>1;"},
	// "23": {0, "return 1>2;"},
	// "24": {1, "return 1>=0;"},
	// "25": {1, "return 1>=1;"},
	// "26": {0, "return 1>=2;"},

	// "27 single-letter variables": {3, "a=3; return  a;"},
	// "28 single-letter variables": {8, "a=3; z=5; return a+z;"},

	// "29 multi-letter variables": {3, "foo=3; return foo;"},
	// "30 multi-letter variables": {8, "foo123=3; bar=5; return foo123+bar;"},

	// "31 multi-letter variables": {8, "foo123=3; returnbar=5; return foo123+returnbar;"},
	// "32 multi-letter variables": {3, "return 3; return 5;"},

	// "33": {1, "1;"},
	// "34": {1, "a_=1; return a_;"},

	// "35 if statement":        {3, "if (0) return 2; return 3;"},
	// "35-1 if else statement": {3, "if (0) return 2; else return 3;"},
	// "36 if statement":        {3, "if (1-1) return 2; return 3;"},
	// "37 if statement":        {2, "if (1) return 2; return 3;"},
	// "38 if statement":        {2, "if (2-1) return 2; return 3;"},

	// "39 while statement": {10, "i=0; while(i<10) i=i+1; return i;"},
	// "40 for statement":   {55, "i=0; j=0; for (i=0; i<=10; i=i+1) j=i+j; return j;"},
	// "41 for statement":   {3, "for (;;) return 3; return 5;"},

	// "42 block statement": {4, "{4;{ 5; } return 4;}"},
	// "43 block statement": {66, "i=0;j=0; while(i<12) {j=i+j; i=i+1;} return j;"},

	// "44 function call": {3, "return ret3();"},
	// "45 function call": {5, "return ret5();"},

	// "46 function call": {8, "return add(3,5);"},
	// "47 function call": {2, "return sub(5,3);"},
	// "48 function call": {21, "return add6(1,2,3,4,5,6);"},
	// "49 function call": {55, "a=11; return fib(11);"},

	"50 unary*&": {3, "main() { x=3; return *&x; }"},
	// 	"51 unary*&": {3, "main() { x=3; y=&x; z=&y; return **z; }"},
	// 	"52 unary*&": {5, "main() { x=3; y=5; return *(&x+8); }"},
	// 	"53 unary*&": {3, "main() { x=3; y=5; return *(&y-8); }"},
	// 	"54 unary*&": {5, "main() { x=3; y=&x; *y=5; return x; }"},
	// 	"55 unary*&": {7, "main() { x=3; y=5; *(&x+8)=7; return y; }"},
	// 	"56 unary*&": {7, "main() { x=3; y=5; *(&y-8)=7; return x; }"},
}

var funcs string = `int ret3() { return 3;}
int ret5() { return 5;}
int add(int x, int y) { return x+y; }
int sub(int x, int y) { return x-y;}

int add6(int a, int b, int c, int d, int e, int f) {
	return a+b+c+d+e+f;
}

int fib(int n) {
	if(n == 1) { return 0;}
	if(n == 2) { return 1;}
	return fib(n-1) + fib(n-2);
}
`

func TestCompile(t *testing.T) {
	var asmName string = "temporary"

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := os.Remove(asmName + ".s")
			if err != nil && !os.IsNotExist(err) {
				t.Fatal(err)
			}
			tmps, err := os.Create(asmName + ".s")
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := tmps.Close(); err != nil {
					t.Fatal(err)
				}
			}()

			if err := compile(c.input, tmps); err != nil {
				t.Fatal(err)
			}

			// make 'funcs' file
			f, err := os.Create("funcs")
			if err != nil {
				return
			}
			defer func() {
				if err = os.Remove(f.Name()); err != nil {
					t.Fatal(err)
				}
			}()

			// write 'funcs' to file
			if _, err = f.WriteString(funcs); err != nil {
				t.Fatal(err)
			}
			if err = f.Sync(); err != nil {
				t.Fatal(err)
			}

			// make a object file from 'funcs'
			_, err = exec.Command("gcc", "-xc", "-c", "-o", asmName+"2.o", f.Name()).Output()
			if err != nil {
				t.Fatal(err)
			}

			// make a execution file with static-link to 'f'
			_, err = exec.Command("gcc", "-static", "-g", "-o", asmName, asmName+".s", asmName+"2.o").Output()
			if err != nil {
				t.Fatal(err)
			}
			// quit this test sequence, if the execution file wasn't made
			if _, err := os.Stat(asmName); err != nil {
				t.Fatal(err)
			}

			_, err = exec.Command("./" + asmName).Output()
			if err != nil {
				if ee, ok := err.(*exec.ExitError); !ok {
					t.Fatal(err)
				} else {
					// the return value of temporary.s is saved in exit status code normally
					actual := ee.ProcessState.ExitCode()
					if c.expected != actual {
						t.Fatalf("\n%s\n%d expected, but got %d", c.input, c.expected, actual)
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
