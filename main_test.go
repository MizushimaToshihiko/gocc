package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestCompile(t *testing.T) {

	var asm *os.File
	var err error
	asm, err = os.Create("testdata/asm.s")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := asm.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	if err = compile("testdata/tests.go", asm); err != nil {
		t.Fatal(err)
	}

	var b []byte
	b, err = exec.Command(
		"gcc",
		"-static",
		"-g",
		"-o",
		"testdata/asm",
		"testdata/asm.s",
	).CombinedOutput()
	if err != nil {
		t.Fatalf("\noutput:\n%s\n%v", string(b), err)
	}

	b, err = exec.Command("testdata/asm").CombinedOutput()
	if err != nil {
		t.Fatalf("\noutput:\n%s\n%v", string(b), err)
	}
	t.Logf("\noutput:\n%s", string(b))
}

func TestGetTypeName(t *testing.T) {
	cases := map[string]struct {
		in   string
		want string
	}{
		"case1": {"[2]int", "[2]int"},
		"case2": {"[2][3]int", "[2][3]int"},
		"case3": {"****int", "pointer"},
		"case4": {"byte", "byte"},
		"case5": {"string", "string"},
		"case6": {"*string", "pointer"},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			userInput = append([]rune(c.in), 0)
			curIdx = 0
			var err error
			token, err = tokenize()
			if err != nil {
				t.Fatal(err)
			}
			ty := readTypePreffix()

			if ty.Name != c.want {
				t.Fatalf("%s expected, but got %s", c.want, ty.Name)
			}
		})
	}
}
