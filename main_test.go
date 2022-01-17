package main

import (
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
		"case3": {"****int", "pointer", TY_PTR},
		"case4": {"byte", "byte", TY_BYTE},
		"case5": {"string", "string", TY_PTR},
		"case6": {"*string", "pointer", TY_PTR},
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

			// fmt.Printf("tok: %#v\n\n", tok)
			if ty.Kind != c.want2 {
				t.Fatalf("%d expected, but got %d", c.want2, ty.Kind)
			}
			if ty.TyName != c.want1 {
				t.Fatalf("%s expected, but got %s", c.want1, ty.TyName)
			}
		})
	}
}
