package main

import (
	"os"
	"reflect"
	"testing"
)

func TestGetTypeName(t *testing.T) {
	cases := map[string]struct {
		in    string
		want1 string
		want2 TypeKind
	}{
		"case1": {"[2]int", "[2]int", TY_ARRAY},
		"case2": {"[2][3]int", "[2][3]int", TY_ARRAY},
		"case3": {"****int", "****int", TY_PTR},
		"case4": {"byte", "uint8", TY_BYTE},
		"case5": {"string", "string", TY_PTR},
		"case6": {"*string", "*string", TY_PTR},
		"case7": {"[1 + 1]int", "[2]int", TY_ARRAY},
		"case8": {"[1 + 1]*int", "[2]*int", TY_ARRAY},
		"case9": {"*[1 % 2]int", "*[1]int", TY_PTR},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			testfile := makeTestFile(t, c.in)
			curIdx = 0
			var err error
			var tok *Token
			tok, err = tokenizeFile(testfile.Name())
			if err != nil {
				t.Fatal(err)
			}
			tok = preprocess(tok)
			ty := readTypePreffix(&tok, tok, nil)

			// fmt.Printf("tok: %#v\n\n", tok)
			if ty.Kind != c.want2 {
				t.Fatalf("%s: %d expected, but got %d", c.in, c.want2, ty.Kind)
			}
			if ty.TyName != c.want1 {
				t.Fatalf("%s: %s expected, but got %s", c.in, c.want1, ty.TyName)
			}
			if tok.Kind != TK_RESERVED &&
				tok.Str == ";" &&
				tok.Next.Kind == TK_EOF {
				t.Fatalf("the token position: EOF expected, but %s", tok.Str)
			}
		})
	}
}

func TestIsTypename(t *testing.T) {
	cases := map[string]struct {
		in   string
		want bool
	}{
		"case 1": {"int", true},
		"case 2": {"string", true},
		"case 3": {"[2]int", true},
		"case 4": {"[2][2]int", true},
		"case 5": {"[int]", false},
		"case 6": {"[[[[1]int", true},
		"case 7": {"[[[[1int", false},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			testfile := makeTestFile(t, c.in)
			curIdx = 0
			var err error
			var tok *Token
			tok, err = tokenizeFile(testfile.Name())
			if err != nil {
				t.Fatal(err)
			}

			act := isTypename(tok)
			if act != c.want {
				t.Fatalf("%s: %t expected, but got %t", c.in, c.want, act)
			}
		})
	}
}

func makeTestFile(t *testing.T, input string) *os.File {
	testfile, err := createTmpFile()
	if err != nil {
		t.Fatalf("makeTestFile: creating testfile failed: %s", err)
	}
	defer func() {
		if err := testfile.Close(); err != nil {
			t.Fatalf("makeTestFile: closing testfile failed: %s", err)
		}
	}()

	b := []byte(input)
	b = append(b, 0)
	if _, err := testfile.Write(b); err != nil {
		t.Fatalf("makeTestFile: writing to testfile failed: %s", err)
	}
	return testfile
}

func TestReadUniversalChar(t *testing.T) {

	cases := map[string]struct {
		in   []byte
		want []byte
	}{
		"case1": {in: []byte(`\u3042`), want: []byte("あ")},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			chara := readUniversalChar16(c.in[2:], 4)
			idx := encodeUft8(&c.in, int(chara), 0)
			c.in = append(c.in[:idx], c.in[6:]...)
			if !reflect.DeepEqual(c.in, c.want) {
				t.Fatalf("expected %v, got %v", c.want, c.in)
			}
		})
	}
}

func TestConvUniversalChars(t *testing.T) {

	cases := map[string]struct {
		in   string
		want []byte
	}{
		"case1": {in: `\u03B1\u03B2\u03B3`, want: []byte("αβγ")},
		"case2": {in: `\u3042`, want: []byte("あ")},
		"case3": {in: `\U000065E5\U0000672C\U00008A9E`, want: []byte("日本語")},
		"case4": {in: `\xff`, want: []byte{255}},
		"case5": {in: `\378`, want: []byte{255}},
		"case6": {in: "\343\201\202", want: []byte("あ")},
		"case7": {in: "\343\201\204", want: []byte("い")},
		"case8": {in: "\xc3\xbf", want: []byte("ÿ")},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			s := []byte(c.in)
			convUniversalChars(&s)
			if !reflect.DeepEqual(s, c.want) {
				t.Fatalf("expected %v: %s, got %v: %s", c.want, string(c.want), s, string(s))
			}
		})
	}
}
