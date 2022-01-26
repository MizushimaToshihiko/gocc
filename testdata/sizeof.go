package test

func main() {
	// assert(1, Sizeof(byte), "Sizeof(byte)")
	// assert(2, Sizeof(int16), "Sizeof(int16)")
	// assert(4, Sizeof(int), "Sizeof(int)")
	// assert(8, Sizeof(int64), "Sizeof(int64)")
	// type T1 struct {
	// 	a int
	// 	b int
	// }
	// assert(8, Sizeof(T1), "Sizeof(type T1 struct {a int; b int;};)")

	// 以下なぜかparseできない
	// var x int = 0
	// assert(1, Sizeof(x+1), "Sizeof(x+1)")
	// assert(8, Sizeof(-10+int64(5)), "Sizeof(-10+int64(5))")
	// assert(8, Sizeof(-10 - int64(5)), "Sizeof(-10 - int64(5))");
	// assert(8, Sizeof(-10 * int64(5), "Sizeof(-10 * int64(5)");
	// assert(8, Sizeof(-10 / int64(5), "Sizeof(-10 / int64(5)");
	// assert(8, Sizeof(int64(-10) + 5), "Sizeof(int64(-10) + 5)");
	// assert(8, Sizeof(int64(-10) - 5), "Sizeof(int64(-10) - 5)");
	// assert(8, Sizeof(int64(-10) * 5), "Sizeof(int64(-10) * 5)");
	// assert(8, Sizeof(int64(-10) / 5), "Sizeof(int64(-10) / 5)");
	var i byte
	assert(1, Sizeof(i++), "var i byte; Sizeof(i++)")

	println("OK")
}
