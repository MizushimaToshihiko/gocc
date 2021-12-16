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

	if err = compile("testdata/tests", asm); err != nil {
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
