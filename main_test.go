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
	"1":  {0, "return 0;"},
	"2":  {42, "return 42;"},
	"3":  {21, "return 5+20-4;"},
	"4":  {41, "return  12 + 34 - 5 ;"},
	"5":  {47, "return 5+6*7;"},
	"6":  {15, "return 5*(9-6);"},
	"7":  {4, "return (3+5)/2;"},
	"8":  {10, "return -10+20;"},
	"9":  {10, "return - -10;"},
	"10": {10, "return - - +10;"},

	"11": {0, "return 0==1;"},
	"12": {1, "return 42==42;"},
	"13": {1, "return 0!=1;"},
	"14": {0, "return 42!=42;"},

	"15": {1, "return 0<1;"},
	"16": {0, "return 1<1;"},
	"17": {0, "return 2<1;"},
	"18": {1, "return 0<=1;"},
	"19": {1, "return 1<=1;"},
	"20": {0, "return 2<=1;"},

	"21": {1, "return 1>0;"},
	"22": {0, "return 1>1;"},
	"23": {0, "return 1>2;"},
	"24": {1, "return 1>=0;"},
	"25": {1, "return 1>=1;"},
	"26": {0, "return 1>=2;"},

	"27 single-letter variables": {3, "a=3; return  a;"},
	"28 single-letter variables": {8, "a=3; z=5; return a+z;"},

	"29 multi-letter variables": {3, "foo=3; return foo;"},
	"30 multi-letter variables": {8, "foo123=3; bar=5; return foo123+bar;"},

	"31 multi-letter variables": {8, "foo123=3; returnbar=5; return foo123+returnbar;"},
	"32 multi-letter variables": {3, "return 3; return 5;"},

	"33": {1, "1;"},
	"34": {1, "a_=1; return a_;"},

	"35 if statement": {3, "if (0) return 2; return 3;"},
	"36 if statement": {3, "if (1-1) return 2; return 3;"},
	"37 if statement": {2, "if (1) return 2; return 3;"},
	"38 if statement": {2, "if (2-1) return 2; return 3;"},

	"39 while statement": {10, "i=0; while(i<10) i=i+1; return i;"},
}

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

			_, err = exec.Command("gcc", "-o", asmName, asmName+".s").Output()
			if err != nil {
				t.Fatal(err)
			}
			// 実行ファイルができていなかったら落とす
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
						t.Fatalf("%d expected, but got %d", c.expected, actual)
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
