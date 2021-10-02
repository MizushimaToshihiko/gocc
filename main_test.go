package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
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
}

func TestE2E(t *testing.T) {
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := os.Remove("tmp.s")
			if err != nil && !os.IsNotExist(err) {
				t.Fatal(err)
			}
			tmps, err := os.OpenFile("tmp.s", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				t.Fatal(err)
			}
			compile(c.input, tmps)

			cc, err := exec.Command("cc", "-o", "tmp", "tmp.s").Output()
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println("cc:", string(cc))

			// tmpファイルができていなかったら落とす
			if _, err := os.Stat("tmp"); err != nil {
				t.Fatal(err)
			}

			tmp, err := exec.Command("./tmp").Output()
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println("tmp:", string(tmp))

			ans, err := exec.Command("echo", "$?").Output()
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println("actual:", string(ans))

			actual, err := strconv.Atoi(string(ans))
			if err != nil {
				t.Fatal(err)
			}

			if c.expected != actual {
				t.Fatalf("%d expected, but got %d", c.expected, actual)
			}
		})
	}
}
