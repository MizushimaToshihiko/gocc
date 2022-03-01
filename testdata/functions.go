package test_function

func assert(want int, act int, code string)
func println(frmt ...string)

func ret3() int {
	return 3
	return 5
}

func add2(x int, y int) int {
	return x + y
}

func sub2(x int, y int) int {
	return x - y
}

func add6(a int, b int, c int, d int, e int, f int) int {
	return a + b + c + d + e + f
}

func addx(x *int, y int) int {
	return *x + y
}

func subChar(a byte, b byte, c byte) int {
	return a - b - c
}

func fib(x int) int {
	if x <= 1 {
		return 1
	}
	return fib(x-1) + fib(x-2)
}

func subLong(a int64, b int64, c int64) int {
	return a - b - c
}

func subShort(a int16, b int16, c int16) int {
	return a - b - c
}

var g1 int

func g1Ptr() *int {
	return &g1
}

func intToChar(x int) byte {
	return x
}

func divLong(a int64, b int64) int {
	return a / b
}

func boolFnAdd(x bool) bool {
	return x + 1
}

func boolFnSub(x bool) bool {
	return x - 1
}

func paramDecay(x []int) int {
	return x[0]
}

func retNone() {
	return
}

func falseFn() bool

func trueFn() bool

func charFn() byte

func shortFn() int16

// sliceの追加後
// func addAll(n ...int) int
// func printAll(s ...string) {
// 	println(s)
// }

func add_double(x float64, y float64) float64
func add_float(x float32, y float32) float32

func add_float3(x float32, y float32, z float32) float32 {
	return x + y + z
}

func add_double3(x float64, y float64, z float64) float64 {
	return x + y + z
}

func sprintf(buf string, format ...string) string
func strcmp(s1 string, s2 string) int

func main() {
	assert(3, ret3(), "ret3()")
	assert(8, add2(3, 5), "add2(3, 5)")
	assert(2, sub2(5, 3), "sub2(5, 3)")
	assert(21, add6(1, 2, 3, 4, 5, 6), "add6(1,2,3,4,5,6)")
	assert(66, add6(1, 2, add6(3, 4, 5, 6, 7, 8), 9, 10, 11), "add6(1,2,add6(3,4,5,6,7,8),9,10,11)")
	assert(136, add6(1, 2, add6(3, add6(4, 5, 6, 7, 8, 9), 10, 11, 12, 13), 14, 15, 16), "add6(1,2,add6(3,add6(4,5,6,7,8,9),10,11,12,13),14,15,16)")

	assert(7, add2(3, 4), "add2(3,4)")
	assert(1, sub2(4, 3), "sub2(4,3)")
	assert(55, fib(9), "fib(9)")

	assert(1, subChar(7, 3, 3), "subChar(7, 3, 3)")

	assert(1, subLong(7, 3, 3), "subLong(7, 3, 3)")
	assert(1, subShort(7, 3, 3), "subShort(7, 3, 3)")

	g1 = 3

	assert(3, *g1Ptr(), "*g1Ptr()")
	assert(5, intToChar(261), "intToChar(261)")
	assert(5, intToChar(261), "intToChar(261)")
	assert(-5, divLong(-10, 2), "divLong(-10, 2)")

	assert(1, boolFnAdd(3), "boolFnAdd(3)")
	assert(0, boolFnSub(3), "boolFnSub(3)")
	assert(1, boolFnAdd(-3), "boolFnAdd(-3)")
	assert(0, boolFnSub(-3), "boolFnSub(-3)")
	assert(1, boolFnAdd(0), "boolFnAdd(0)")
	assert(1, boolFnSub(0), "boolFnSub(0)")
	var x [2]int
	x[0] = 3
	assert(3, paramDecay(x), "var x [2]int ; x[0]=3; paramDecay(x)")

	retNone()

	assert(1, trueFn(), "trueFn()")
	assert(0, falseFn(), "falseFn()")
	assert(3, charFn(), "charFn()")
	assert(5, shortFn(), "shortFn()")

	// sliceの追加後
	// assert(6, addAll(3, 1, 2, 3), "addAll(3,1,2,3)")
	// assert(5, addAll(4, 1, 2, 3, -1), "addAll(4,1,2,3,-1)")
	// printAll("1", "2", "3", "4")
	// printAll("1", "2", "3", "4", "5", "6")

	assert(6, int(add_float(2.3, 3.8)), "int(add_float(2.3, 3.8))")
	assert(6, int(add_double(2.3, 3.8)), "int(add_double(2.3, 3.8))")

	assert(7, int(add_float3(2.5, 2.5, 2.5)), "int(add_float3(2.5, 2.5, 2.5))")
	assert(7, int(add_double3(2.5, 2.5, 2.5)), "int(add_double3(2.5, 2.5, 2.5))")

	var buf string
	sprintf(buf, "%.1f", float32(3.5))
	assert(0, strcmp(buf, "3.5"), "var buf string;sprintf(buf,\"%.1f\",float32(3.5));strcmp(buf,\"3.5\")")

	println("OK")
}
