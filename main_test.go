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
	"1":  {0, "0"},
	"2":  {42, "42"},
	"3":  {41, " 12 + 34 - 5 "},
	"4":  {47, "5+6*7"},
	"5":  {15, "5*(9-6)"},
	"6":  {4, "(3+5)/2"},
	"7":  {10, "-10+20"},
	"8":  {10, "- -10"},
	"9":  {10, "- - +10"},
	"10": {0, "0==1"},
	"11": {1, "42==42"},
	"12": {1, "0!=1"},
	"13": {0, "42!=42"},
	"14": {1, "0<1"},
	"15": {0, "1<1"},
	"16": {0, "2<1"},
	"17": {1, "0<=1"},
	"18": {1, "1<=1"},
	"19": {0, "2<=1"},
	"20": {1, "1>0"},
	"21": {0, "1>1"},
	"22": {0, "1>2"},
	"23": {1, "1>=0"},
	"24": {1, "1>=1"},
	"25": {0, "1>=2"},
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
