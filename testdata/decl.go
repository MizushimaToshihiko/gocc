package test_decl

func assert(want int, act int, code string)
func println(format string)

func strcmp(s1 string, s2 string) int

var g1, g2, g3 bool
var g4, g5, g6 = 2.0, 8, "foo"

var (
	g7          int
	g8, g9, g10 = 2.0, 3.0, "bar"
)

func main() {
	var x1 byte
	assert(1, Sizeof(x1), "var x1 byte; Sizeof(x1)")
	var x2 int16
	assert(2, Sizeof(x2), "var x2 int16; Sizeof(x2)")
	var x3 int
	assert(4, Sizeof(x3), "var x3 int; Sizeof(x3)")
	var x4 int64
	assert(8, Sizeof(x4), "var x4 int64; Sizeof(x4)")

	var x5 bool = 0
	assert(0, x5, "var x5 bool=0; x5")
	var x6 bool = 1
	assert(1, x6, "var x6 bool=1; x6")
	var x7 bool = 2
	assert(1, x7, "var x7 bool=2; x7")
	assert(1, bool(1), "bool(1)")
	assert(1, bool(2), "bool(2)")
	assert(0, bool(byte(256)), "bool(byte(256))")

	var x8, x9 int
	assert(0, x8, "x8")
	assert(4, Sizeof(x8), "Sizeof(x8)")
	assert(0, x9, "x9")

	var x10, x11, x12, x13 int = 1, 2, 3, 4
	assert(1, x10, "x10")
	assert(4, Sizeof(x10), "Sizeof(x10)")
	assert(2, x11, "x11")
	assert(3, x12, "x12")
	assert(4, x13, "x13")

	var x14, x15, x16, x17 string = "1", "2", "3", "4"
	assert(8, Sizeof(x14), "Sizeof(x14)")
	assert(0, strcmp(x14, "1"), "strcmp(x14, \"1\")")
	assert(0, strcmp(x15, "2"), "strcmp(x15, \"2\")")
	assert(0, strcmp(x16, "3"), "strcmp(x16, \"3\")")
	assert(0, strcmp(x17, "4"), "strcmp(x17, \"4\")")

	x18, x19, x20 := 1, 2, 3
	assert(1, x18, "x18")
	assert(2, x19, "x19")
	assert(3, x20, "x20")

	x21, x22, x23, x24 := "1", "2", "3", "4"
	assert(8, Sizeof(x21), "Sizeof(x21)")
	assert(0, strcmp(x21, "1"), "strcmp(x21, \"1\")")
	assert(0, strcmp(x22, "2"), "strcmp(x22, \"2\")")
	assert(0, strcmp(x23, "3"), "strcmp(x23, \"3\")")
	assert(0, strcmp(x24, "4"), "strcmp(x24, \"4\")")

	var (
		i25           int
		u25, v25, s25 = 2.0, 3.0, "bar"
	)
	assert(0, i25, "i25")
	assert(2.0, u25, "u25")
	assert(3.0, v25, "v25")
	assert(0, strcmp(s25, "bar"), "strcmp(s25, \"bar\")")

	assert(0, g1, "g1")
	assert(0, g2, "g2")
	assert(0, g3, "g3")
	assert(2.0, g4, "g4")
	assert(8, g5, "g5")
	assert(0, strcmp(g6, "foo"), "strcmp(g6, \"foo\")")

	assert(0, g7, "g7")
	assert(2.0, g8, "g8")
	assert(3.0, g9, "g9")
	assert(0, strcmp(g10, "bar"), "strcmp(g10, \"bar\")")

	println("OK")
}
