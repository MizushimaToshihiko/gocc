package test

func assert(want int, act int, code string)
func println(format string)

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

	println("OK")
}
