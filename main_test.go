package main

// func TestGetTypeName(t *testing.T) {
// 	cases := map[string]struct {
// 		in   string
// 		want string
// 	}{
// 		"case1": {"[2]int", "[2]int"},
// 		"case2": {"[2][3]int", "[2][3]int"},
// 		"case3": {"****int", "pointer"},
// 		"case4": {"byte", "byte"},
// 		"case5": {"string", "string"},
// 		"case6": {"*string", "pointer"},
// 	}

// 	for name, c := range cases {
// 		t.Run(name, func(t *testing.T) {
// 			userInput = append([]rune(c.in), 0)
// 			curIdx = 0
// 			var err error
// 			token, err = tokenize("")
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			ty := readTypePreffix()

// 			if ty.Name != c.want {
// 				t.Fatalf("%s expected, but got %s", c.want, ty.Name)
// 			}
// 		})
// 	}
// }
