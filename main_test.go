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
	"1":  {0, "0;"},
	"2":  {42, "42;"},
	"3":  {21, "5+20-4;"},
	"4":  {41, " 12 + 34 - 5 ;"},
	"5":  {47, "5+6*7;"},
	"6":  {15, "5*(9-6);"},
	"7":  {4, "(3+5)/2;"},
	"8":  {10, "-10+20;"},
	"9":  {10, "- -10;"},
	"10": {10, "- - +10;"},

	"11": {0, "0==1;"},
	"12": {1, "42==42;"},
	"13": {1, "0!=1;"},
	"14": {0, "42!=42;"},

	"15": {1, "0<1;"},
	"16": {0, "1<1;"},
	"17": {0, "2<1;"},
	"18": {1, "0<=1;"},
	"19": {1, "1<=1;"},
	"20": {0, "2<=1;"},

	"21": {1, "1>0;"},
	"22": {0, "1>1;"},
	"23": {0, "1>2;"},
	"24": {1, "1>=0;"},
	"25": {1, "1>=1;"},
	"26": {0, "1>=2;"},

	"27": {3, "a=3; a;"},
	"28": {8, "a=3; z=5; a+z;"},

	"29": {3, "foo=3; foo;"},
	"30": {8, "foo123=3; bar=5; foo123+bar;"},
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
