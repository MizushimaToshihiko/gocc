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
	"1": {
		expected: 0,
		input:    "0",
	},
	"2": {
		expected: 42,
		input:    "42",
	},
	"3": {
		expected: 41,
		input:    " 12 + 34 - 5 ",
	},
	"4": {
		expected: 47,
		input:    "5+6*7",
	},
	"5": {
		expected: 15,
		input:    "5*(9-6)",
	},
	"6": {
		expected: 4,
		input:    "(3+5)/2",
	},
	"7": {
		expected: 8,
		input:    "5--3",
	},
	"8": {
		expected: 1,
		input:    "9-(3+5)",
	},
	"9": {
		expected: 10,
		input:    "-10+20",
	},
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
					// the return value of temporary.s is saved in exit status code
					actual := ee.ProcessState.ExitCode()
					if c.expected != actual {
						t.Fatalf("%d expected, but got %d", c.expected, actual)
					}
					t.Logf("%s => %d", c.input, actual)
					return
				}
			}

			ans, err := exec.Command("sh", "-c", "echo $?").Output()
			if err != nil {
				t.Fatal(err)
			}

			// the return value of temporary.s is saved in exit status code,
			// so the below will be used only when the return value is 0.
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
