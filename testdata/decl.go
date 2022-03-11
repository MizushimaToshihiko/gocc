package test_decl

func assert(want int, act int, code string)
func println(format string)

func strcmp(s1 string, s2 string) int

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

	println("OK")
}
