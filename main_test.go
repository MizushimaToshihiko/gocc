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
			userInput = append([]rune(c.in), 0)
			curIdx = 0
			var err error
			var tok *Token
			tok, err = tokenize("main_test")
			if err != nil {
				t.Fatal(err)
			}
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

			userInput = append([]rune(c.in), 0)
			curIdx = 0
			var err error
			var tok *Token
			tok, err = tokenize("main_test")
			if err != nil {
				t.Fatal(err)
			}

			printTokens(tok)
			act := isTypename(tok)
			if act != c.want {
				t.Fatalf("%s: %t expected, but got %t", c.in, c.want, act)
			}
		})
	}
}

// func TestTokenize(t *testing.T) {
// 	cases := map[string]struct {
// 		in   string
// 		want string
// 	}{
// 		"case 1": {in: "1,2,3", want: "1"},
// 	}

// 	for name, c := range cases {
// 		t.Run(name, func(t *testing.T) {

// 			userInput = append([]rune(c.in), 0)
// 			curIdx = 0
// 		})
// 	}

// }
