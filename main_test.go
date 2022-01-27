package main

import (
	"fmt"
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
		"case4": {"byte", "byte", TY_BYTE},
		"case5": {"string", "string", TY_PTR},
		"case6": {"*string", "*string", TY_PTR},
		"case7": {"[1 + 1]int", "[2]int", TY_ARRAY},
		"case8": {"[1 + 1]*int", "[2]*int", TY_ARRAY},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			userInput = append([]rune(c.in), 0)
			curIdx = 0
			var err error
			var tok *Token
			tok, err = tokenize("main_test")
			if err != nil {
				t.Fatal(err)
			}
			ty := readTypePreffix(&tok, tok)

			fmt.Printf("tok: %#v\n\n", tok)
			if ty.Kind != c.want2 {
				t.Fatalf("%s: %d expected, but got %d", c.in, c.want2, ty.Kind)
			}
			if ty.TyName != c.want1 {
				t.Fatalf("%s: %s expected, but got %s", c.in, c.want1, ty.TyName)
			}
			if tok.Kind != TK_EOF {
				t.Fatalf("the token position: EOF expected, but %s", tok.Str)
			}
		})
	}
}
